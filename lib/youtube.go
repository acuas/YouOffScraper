package lib

import (
	"context"
	"flag"
	"fmt"
	ydr "github.com/kkdai/youtube"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"os/user"
	"path/filepath"
	"strings"
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
	y := ydr.NewYoutube(true)
	for _, item := range response.Items {
		video := YouTubeVideo{
			VideoId: item.Id.VideoId,
			ChannelId: item.Snippet.ChannelId,
			PublishedAt: item.Snippet.PublishedAt,
			Title: item.Snippet.Title,
			Description: item.Snippet.Description,
			ChannelTitle: item.Snippet.ChannelTitle,
		}

		y.DecodeURL(fmt.Sprintf("https://www.youtube.com/watch?v=%v", video.VideoId))
		if //noinspection GoFunctionCall
		err := y.StartDownload(filepath.Join(currentDir, video.VideoId + ".mp4")); err != nil {
			log.Fatalln("err: ", err)
		}
		log.Println(video)
	}
}