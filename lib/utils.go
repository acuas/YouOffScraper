package lib

import "github.com/youoffcrawler/config"

type Application struct {
	Config *config.Config
}

var App *Application