package lib

import (
	"github.com/minio/minio-go/v6"
	"github.com/rylio/ytdl"
	"github.com/youoffcrawler/config"
)

type Application struct {
	Config        *config.Config
	YouTubeClient *ytdl.Client
	MinioClient   *minio.Client
}

var App *Application
