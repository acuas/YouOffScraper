package lib

import (
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/gofiber/fiber"
	"github.com/minio/minio-go/v6"
	"github.com/rylio/ytdl"
	"github.com/youoffcrawler/config"
)

///////////////////////////////////////////////////////////////////////////////

// Application represents the heart of the crawler. It contains
// all important third parties that is used within this project.
type Application struct {
	Config        *config.Config
	YouTubeClient *ytdl.Client
	MinioClient   *minio.Client
	ES            *elasticsearch7.Client
	Srv           *fiber.App
}

// Exported pointer to be accessible from entire project
var App *Application

///////////////////////////////////////////////////////////////////////////////
