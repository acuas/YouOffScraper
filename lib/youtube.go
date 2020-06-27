package lib

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"github.com/minio/minio-go/v6"
	"time"
)

// YouTubeSvc is the service which it is used by the scraper
// to get url of videos from a chanel
var YouTubeSvc *youtube.Service

func SetupYouTubeSvc () {
	ctx := context.Background()
	var err error
	YouTubeSvc, err = youtube.NewService(ctx, option.WithAPIKey(App.Config.YTApiKey))
	if err != nil {
		log.Fatalln(err)
	}
}


type YouTubeVideo struct {
	VideoId      string `json:"videoId"`
	ChannelId    string `json:"channelId"`
	PublishedAt  string `json:"publishedAt"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ChannelTitle string `json:"channelTitle"`
}

type YouTubeChannel struct {
	Id          string `json:"channelId"`
	Title       string `json:"channelTitle"`
	Description string `json:"description"`
	Country     string `json:"country"`
	PublishedAt string `json:"publishedAt"`
}

func (channel *YouTubeChannel) NewChannelFromUrl(url string) {
	urlComponents := strings.Split(url, "/")
	channel.Id = urlComponents[4]
	call := YouTubeSvc.Channels.List([]string{"snippet"}).Id(channel.Id)
	response, err := call.Do()
	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range response.Items {
		if item.Kind == "youtube#channel" {
			channel.Title = item.Snippet.Title
			channel.Description = item.Snippet.Description
			channel.Country = item.Snippet.Country
			channel.PublishedAt = item.Snippet.PublishedAt
		}
	}
}

func (channel *YouTubeChannel) ScrapeVideos() {
	maxResults := flag.Int64("max-results", 50, "Max Youtube results")
	videoSyndicated := flag.String("videoSyndicated", "true", "Search to only videos that can be played outside youtube.com")
	call := YouTubeSvc.Search.List([]string{"snippet"}).
		ChannelId(channel.Id).
		Order("date").
		Type("video").
		VideoSyndicated(*videoSyndicated).
		MaxResults(*maxResults)

	response, err := call.Do()
	if err != nil {
		log.Fatalln(err)
	}

	usr, _ := user.Current()
	currentDir := fmt.Sprintf("%v/Movies/youtubedr", usr.HomeDir)
	log.Println("download to dir=", currentDir)
	client := App.YouTubeClient

	endpoint := "youoffminio:9000"
	accessKeyID := "acces_key_development"
	secretAccessKey := "secret_key_development"
	useSSL := false

	// Initialize minio client object
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}


	log.Printf("%v\n", minioClient)
	for _, item := range response.Items {
		// Decode item data
		video := YouTubeVideo{
			VideoId: item.Id.VideoId,
			ChannelId: item.Snippet.ChannelId,
			PublishedAt: item.Snippet.PublishedAt,
			Title: item.Snippet.Title,
			Description: item.Snippet.Description,
			ChannelTitle: item.Snippet.ChannelTitle,
		}

		// Download video
		ctx := context.Background()
		vid, err := client.GetVideoInfo(ctx, fmt.Sprintf("https://www.youtube.com/watch?v=%v", video.VideoId))
		if err != nil {
			log.Println("Failed to get video info")
			return
		}
		objectName := video.ChannelId + "/" + video.VideoId + ".mp4"
		pathToVideo := filepath.Join(currentDir, objectName)
		file, _ := os.Create(pathToVideo)
		client.Download(ctx, vid, vid.Formats[0], file)

		// Upload the video to s3
		_, err = minioClient.FPutObject("youtube", objectName, pathToVideo, minio.PutObjectOptions{ContentType: "video/mp4"})
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Succes upload")
		file.Close()
		time.Sleep(120000)
	}
}