package lib

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v6"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// YouTubeSvc is the service which it is used by the scraper
// to get url of videos from a chanel
var YouTubeSvc *youtube.Service

func SetupYouTubeSvc() {
	ctx := context.Background()
	var err error
	YouTubeSvc, err = youtube.NewService(ctx, option.WithAPIKey(App.Config.YTApiKey))
	if err != nil {
		log.Fatalln(err)
	}
}

// YouTubeVideo is the structure where we store details about a YouTube video
type YouTubeVideo struct {
	VideoId      string `json:"videoId"`
	ChannelId    string `json:"channelId"`
	PublishedAt  string `json:"publishedAt"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ChannelTitle string `json:"channelTitle"`
}

// NewVideoFromUrl create a YouTubeVideo filling its field calling
// the YouTube API
func NewVideoFromUrl(urlStr string) (*YouTubeVideo, error){
	// Parse the video url
	video := &YouTubeVideo{}
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
		return video, err
	}
	q := parsedUrl.Query()
	video.VideoId = q.Get("v")
	if video.VideoId == "" {
		return video, errors.New("The url isn't valid!")
	}

	// Query the YouTube API
	res, err := YouTubeSvc.Videos.List([]string{"snippet"}).Id(video.VideoId).Do()
	if err != nil {
		log.Println(err)
		return video, err
	}

	if len(res.Items) > 1 {
		return video, errors.New("Too many videos! Something is wrong!")
	}

	item := res.Items[0]
	video.VideoId = item.Id
	video.ChannelTitle = item.Snippet.ChannelTitle
	video.ChannelId = item.Snippet.ChannelId
	video.Description = item.Snippet.Description
	video.Title = item.Snippet.Title
	video.PublishedAt = item.Snippet.PublishedAt

	return video, nil
}

func (video *YouTubeVideo) Download() error {
	if video.VideoId == "" || video.ChannelId == "" {
		return errors.New("Scraper doesn't have details about the video!")
	}

	usr, _ := user.Current()
	currentDir := fmt.Sprintf("%v/Movies/youtubedr", usr.HomeDir)
	// Attempt to create the directory and ignore any issues
	_ = os.Mkdir(currentDir+"/"+video.ChannelId, os.ModeDir)

	var err error
	ctx := context.Background()
	vid, err := App.YouTubeClient.GetVideoInfo(ctx, fmt.Sprintf("https://www.youtube.com/watch?v=%v", video.VideoId))
	if err != nil {
		log.Println("Failed to get video info")
		return err
	}

	objectName := video.ChannelId + "/" + video.VideoId + ".mp4"
	pathToVideo := filepath.Join(currentDir, objectName)
	file, _ := os.Create(pathToVideo)
	defer file.Close()
	err = App.YouTubeClient.Download(ctx, vid, vid.Formats[0], file)
	if err != nil {
		log.Println("wow" + err.Error())
		return err
	}

	// Upload the video to s3
	_, err = App.MinioClient.FPutObject(App.Config.MinioBucketName, objectName, pathToVideo, minio.PutObjectOptions{ContentType: "video/mp4"})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Succes upload of video %v to S3", video.VideoId)
	return nil
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

func (channel *YouTubeChannel) ScrapeChannel() {
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

	for _, item := range response.Items {
		// Decode item data
		video := &YouTubeVideo{
			VideoId:      item.Id.VideoId,
			ChannelId:    item.Snippet.ChannelId,
			PublishedAt:  item.Snippet.PublishedAt,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			ChannelTitle: item.Snippet.ChannelTitle,
		}

		// Download video
		err = video.Download()
		if err != nil {
			log.Printf("Skip video with id %v", video.VideoId)
		}
		time.Sleep(time.Hour)
	}
}
