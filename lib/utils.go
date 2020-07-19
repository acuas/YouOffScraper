package lib

import (
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/minio/minio-go/v6"
	"github.com/rylio/ytdl"
	"github.com/youoffcrawler/config"
)

type Application struct {
	Config        *config.Config
	YouTubeClient *ytdl.Client
	MinioClient   *minio.Client
	ES            *elasticsearch7.Client
}

var App *Application
