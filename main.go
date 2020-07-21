package main

import (
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v6"
	zerolog "github.com/rs/zerolog/log"
	"github.com/rylio/ytdl"
	"github.com/youoffcrawler/config"
	"github.com/youoffcrawler/lib"
	. "github.com/youoffcrawler/api/v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
    . "gitlab.com/c0b/go-ordered-json"
)

///////////////////////////////////////////////////////////////////////////////

// init is invoked before main()
func init() {
	// loads valued from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Load env variables by config package
	lib.App = &lib.Application{
		Config: config.New(),
		YouTubeClient: &ytdl.Client{
			HTTPClient: http.DefaultClient,
			Logger:     zerolog.Logger,
		},
	}


	// Initialize MinIO client
	var err error
	lib.App.MinioClient, err = minio.New(
		lib.App.Config.MinioEndpoint,
		lib.App.Config.MinioAccessKeyID,
		lib.App.Config.MinioSecretAccessKey,
		lib.App.Config.MinioUseSSL,
	)
	if err != nil {
		log.Fatal("Error in initializing MinIO client!")
	}

	// Check if the bucket where the crawler is going to store the videos exists
	found, err := lib.App.MinioClient.BucketExists(lib.App.Config.MinioBucketName)
	if err != nil {
		log.Fatal("Error in checking if the bucket exists in MinIO!")
	}

	log.Printf("Checked connectivity to MinIO at %v !", time.Now().Format(time.RFC3339))
	// If the bucket doesn't exist the crawler is going to create it
	if !found {
		log.Println("Bucket doesn't exist, so the crawler will create it, according to your env variable MINIO_BUCKET")
		err = lib.App.MinioClient.MakeBucket(lib.App.Config.MinioBucketName, lib.App.Config.MinioBucketRegion)
		if err != nil {
			log.Fatal("The crawler couldn't create the new bucket! Check MinIO instance for more details!")
		} else {
			log.Printf("The crawler created the bucket %v with succes!", lib.App.Config.MinioBucketName)
		}
	} else {
		log.Printf("The bucket %v already exists!", lib.App.Config.MinioBucketName)
	}


	// Check connectivity to elasticsearch
	lib.App.ES, err = elasticsearch7.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	res, err := lib.App.ES.Info()
	if err != nil {
		log.Fatalf("Error getting response from elasticsearch: %s", err)
	}
	log.Printf("Checked connectivity to elasticsearch at %v !", time.Now().Format(time.RFC3339))
	io.Copy(ioutil.Discard, res.Body)
	defer res.Body.Close()

	// Try to create an index if it's not already present
	lib.App.ES.Indices.Create(
		lib.App.Config.EsIndex,
		lib.App.ES.Indices.Create.WithBody(strings.NewReader(`{
			"mappings": {
				"properties": {
					"channel_id": { "type": "keyword" },
					"published_at": { "type": "date", "format": "epoch_second"},
					"title": { "type": "text" },
					"description": { "type": "text" },
					"channel_title": { "type": "text" }
				}
			}
		}`)),
	)


	// Start a number of workers to process the workload with concurrency that enables parallelism
	lib.StartDispatcher(1)

	lib.SetupYouTubeSvc()

	// Setup fiber server
	lib.App.Srv = fiber.New()

	lib.App.Srv.Get("/", func(c *fiber.Ctx) {
		c.JSON(NewOrderedMapFromKVPairs([]*KVPair{
			{Key: "ok", Value: 1},
			{Key: "data", Value: NewOrderedMapFromKVPairs([]*KVPair{
				{Key: "name", Value: "api"},
				{Key: "path", Value: "/api"},
			})},
		}))
		c.Status(200)
	})

	lib.App.Srv.Get("/api", func(c *fiber.Ctx) {
		c.JSON(NewOrderedMapFromKVPairs([]*KVPair{
			{Key: "ok", Value: 1},
			{Key: "data", Value: NewOrderedMapFromKVPairs([]*KVPair{
				{Key: "name", Value: "v1"},
				{Key: "path", Value: "/v1"},
			})},
		}))
		c.Status(200)
	})

	// Organize the API using a Group
	api := lib.App.Srv.Group("/api")

	// Mount V1 api
	V1(api)

	// TODO: Make this configurable via env
	lib.App.Srv.Listen(8000)
	//youTubeC := &lib.YouTubeChannel{}
	//youTubeC.NewChannelFromUrl("https://www.youtube.com/channel/UC9WayAVqWKIoyg1eN28n9Ug")
	//youTubeC.ScrapeChannel()
}

///////////////////////////////////////////////////////////////////////////////