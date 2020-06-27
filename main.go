package main

import (
	"github.com/joho/godotenv"
	zerolog "github.com/rs/zerolog/log"
	"github.com/rylio/ytdl"
	"github.com/youoffcrawler/config"
	"github.com/youoffcrawler/lib"
	"log"
	"net/http"
)


// init is invoked before main()
func init() {
	// loads valued from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	// Load env variables by config package
	lib.App = &lib.Application{
		Config: config.New(),
		YouTubeClient: &ytdl.Client{
			HTTPClient: http.DefaultClient,
			Logger: zerolog.Logger,
		},
	}

	lib.SetupYouTubeSvc()
	youTubeC := &lib.YouTubeChannel{}
	youTubeC.NewChannelFromUrl("https://www.youtube.com/channel/UC9WayAVqWKIoyg1eN28n9Ug")
	youTubeC.ScrapeVideos()
}