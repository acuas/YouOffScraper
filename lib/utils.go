package lib

import (
	"github.com/rylio/ytdl"
	"github.com/youoffcrawler/config"
)

type Application struct {
	Config *config.Config
	YouTubeClient *ytdl.Client
}

var App *Application