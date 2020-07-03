package lib

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/minio/minio-go/v6"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////

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

///////////////////////////////////////////////////////////////////////////////

// YouTubeVideo is the structure where scraper store details about a YouTube video
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
func NewVideoFromUrl(urlStr string) (*YouTubeVideo, error) {
	// Parse the video url
	video := &YouTubeVideo{}
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
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

// Download is the function which check if a video exists in Minio. If it does
// not exists in the Minio the function will download and upload it to the
// Minio bucket specified in the configuration file
func (video *YouTubeVideo) Download(finished chan bool) error {
	if video.VideoId == "" || video.ChannelId == "" {
		finished <- false
		return errors.New("Scraper doesn't have details about the video!")
	}

	var err error
	// Check if the object exists in youtube bucket
	_, err = App.MinioClient.StatObject(
		App.Config.MinioBucketName,
		video.ChannelId + "/" + video.VideoId + ".mp4",
		minio.StatObjectOptions{},
	)

	if err == nil {
		log.Printf("Video %v is already in the bucket. Skip.", video.VideoId)
		finished <- true
		return nil
	}

	// Attempt to create a temporary directory and ignore any issues
	usr, _ := user.Current()
	currentDir := fmt.Sprintf("%v/Movies/youtubedr", usr.HomeDir)
	_ = os.MkdirAll(currentDir+"/"+video.ChannelId, os.ModeDir)

	// Download the video
	ctx := context.Background()
	vid, err := App.YouTubeClient.GetVideoInfo(ctx, fmt.Sprintf("https://www.youtube.com/watch?v=%v", video.VideoId))
	if err != nil {
		finished <- false
		return err
	}
	objectName := video.ChannelId + "/" + video.VideoId + ".mp4"
	pathToVideo := filepath.Join(currentDir, objectName)
	file, _ := os.Create(pathToVideo)
	err = App.YouTubeClient.Download(ctx, vid, vid.Formats[0], file)
	if err != nil {
		return err
	}

	// Upload the video to s3
	_, err = App.MinioClient.FPutObject(App.Config.MinioBucketName, objectName, pathToVideo, minio.PutObjectOptions{ContentType: "video/mp4"})
	if err != nil {
		finished <- false
		log.Fatalln(err)
	}
	finished <- true
	log.Printf("Succes upload of video %v to S3", video.VideoId)

	// Remove the video from temporary directory
	file.Close()
	os.Remove(pathToVideo)
	return nil
}

///////////////////////////////////////////////////////////////////////////////

// YouTubeChannel is the structure which is used by scraper to maintain details
// about a YouTube channel
type YouTubeChannel struct {
	Id          string `json:"channelId"`
	Title       string `json:"channelTitle"`
	Description string `json:"description"`
	Country     string `json:"country"`
	PublishedAt string `json:"publishedAt"`
}

// NewChannelFromUrl create a YouTubeChannel from a giver url. The function
// takes the details about the channel using the YouTube Data API
func (channel *YouTubeChannel) NewChannelFromUrl(url string) error {
	urlComponents := strings.Split(url, "/")
	channel.Id = urlComponents[4]
	call := YouTubeSvc.Channels.List([]string{"snippet"}).Id(channel.Id)
	response, err := call.Do()
	if err != nil {
		return err
	}

	for _, item := range response.Items {
		if item.Kind == "youtube#channel" {
			channel.Title = item.Snippet.Title
			channel.Description = item.Snippet.Description
			channel.Country = item.Snippet.Country
			channel.PublishedAt = item.Snippet.PublishedAt
		}
	}
	return nil
}

// A buffered channel through we can send video to be downloaded by the workers
var VideoQueue = make(chan YouTubeVideo, 4)

// Given a YouTube channel the scraper can extract all video from it using this
// function. The function will iterate over all videos extracting batch of 50
// videos which represents the maximum results that can be returned by the
// YouTube Data API
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

	totalNumberOfVideos := response.PageInfo.TotalResults
	for totalNumberOfVideos > 0 {
		totalNumberOfVideos -= 50
		for _, item := range response.Items {
			// Decode item data
			video := YouTubeVideo{
				VideoId:      item.Id.VideoId,
				ChannelId:    item.Snippet.ChannelId,
				PublishedAt:  item.Snippet.PublishedAt,
				Title:        item.Snippet.Title,
				Description:  item.Snippet.Description,
				ChannelTitle: item.Snippet.ChannelTitle,
			}

			VideoQueue <- video
		}

		log.Printf("Go to next page %v", response.NextPageToken)
		// Go to next page
		response, err = call.PageToken(response.NextPageToken).Do()
		if err != nil {
			log.Fatalln(err)
		}
	}
}

///////////////////////////////////////////////////////////////////////////////