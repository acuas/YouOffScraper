package v1

import (
	"fmt"

	"github.com/gofiber/fiber"
	"github.com/youoffcrawler/lib"
	. "gitlab.com/c0b/go-ordered-json"
)

func V1(api *fiber.Group) {
	v1 := api.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) {
		c.Status(fiber.StatusOK).JSON(NewOrderedMapFromKVPairs([]*KVPair{
			{Key: "ok", Value: 1},
			{Key: "data", Value: NewOrderedMapFromKVPairs([]*KVPair{
				{Key: "name", Value: "scrape"},
				{Key: "path", Value: "/scrape"},
				{Key: "info", Value: "Scrape an entire channel or video from YouTube"},
				{Key: "types", Value: NewOrderedMapFromKVPairs([]*KVPair{
					{Key: "channel", Value: fiber.Map{}},
					{Key: "video", Value: fiber.Map{}},
				})},
			})},
		}))
	})

	v1.Post("/scrape", func(c *fiber.Ctx) {
		t := c.Query("type")
		id := c.Query("id")
		if id == "" {
			panic("Parameter id must be not null!")
		}

		switch t {
		case "channel":
			{
				youTubeChan := &lib.YouTubeChannel{}
				err := youTubeChan.NewChannelFromUrl(fmt.Sprintf("https://www.youtube.com/channel/%v", id))
				if err != nil {
					panic(err)
				}
				go youTubeChan.ScrapeChannel()
			}
		case "video":
			{
				youTubeVid, err := lib.NewVideoFromUrl(fmt.Sprintf("https://www.youtube.com/watch?v=%v", id))
				if err != nil {
					panic(err)
				}

				go youTubeVid.Download(make(chan bool))
			}
		default:
			{
				panic("Type must be channel or video!")
			}
		}

		c.JSON(NewOrderedMapFromKVPairs([]*KVPair{
			{Key: "ok", Value: 1},
			{Key: "data", Value: NewOrderedMapFromKVPairs([]*KVPair{
				{Key: "bucket", Value: lib.App.Config.MinioBucketName},
				{Key: "minio_path", Value: fmt.Sprintf("%v/%v", t, id)},
				{Key: "info", Value: "Your channel/video will be soon available in Minio."},
			})},
		}))
		c.SendStatus(fiber.StatusCreated)
	})
}
