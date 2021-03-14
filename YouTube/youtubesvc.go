package youtube

import (
	"context"

	"github.com/elastic/go-ucfg"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type config struct {
	APIKey string `config:"api_key"`
}

var svc *youtube.Service

// NewService must be called before the youtube package is going to be used.
// It receives a context and a configuration which contains the api token
// which can be found (or can be created) here:
// https://console.developers.google.com/apis/credentials
func NewService(ctx context.Context, cfg *ucfg.Config) error {
	config := config{}
	err := cfg.Unpack(&config)
	if err != nil {
		return err
	}
	svc, err = youtube.NewService(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		return err
	}
	return nil
}
