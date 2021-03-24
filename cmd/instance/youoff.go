package instance

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	youtube "github.com/acuas/YouOffScraper/YouTube"
	"github.com/acuas/YouOffScraper/storage"
	"github.com/acuas/YouOffScraper/utils"
	"go.uber.org/zap"

	ucfg "github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
)

var log *zap.SugaredLogger

type YouOffConfig struct {
	Storage  *ucfg.Config `config:"storage"`
	Youtube  *ucfg.Config `config:"youtube"`
	Loggging *ucfg.Config `config:"logging"`
}

type YouOff struct {
	Ctx context.Context

	Storage storage.Storage
	URL     string
}

func Run(cpath, URL string) error {
	// load the configuration
	youOffConfig, err := loadConfig(cpath)
	if err != nil {
		return err
	}

	// if are multiple storage used in config just first read it will be used
	storageTypes := youOffConfig.Storage.GetFields()
	if len(storageTypes) < 1 {
		return fmt.Errorf("storage must be included in your configuration file")
	}
	// check if storage exists and it's registered
	storageFactory, err := storage.GetFactory(storageTypes[0])
	if err != nil {
		return err
	}
	storageConfig, err := youOffConfig.Storage.Child(storageTypes[0], -1)
	if err != nil {
		return fmt.Errorf("something is wrong in storage configuration: %s", err.Error())
	}
	storage, err := storageFactory(storageConfig)
	if err != nil {
		return fmt.Errorf("cannot create storage %v because an error occured: %v", storageTypes[0], err.Error())
	}

	// logging configuration
	utils.ConfigureLogger(youOffConfig.Loggging)
	log, err = utils.NewLogger("instance")
	if err != nil {
		return fmt.Errorf("error creating a new logger: %v", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the youtube service
	err = youtube.NewService(ctx, youOffConfig.Youtube)
	if err != nil {
		return err
	}

	youoff := &YouOff{
		Ctx:     ctx,
		Storage: storage,
		URL:     URL,
	}

	return launch(youoff)
}

func launch(youoff *YouOff) error {
	log.Info("launching youoff...")
	channel, err := youtube.NewChannel(youoff.URL)
	if err != nil {
		return err
	}

	playlist, err := youtube.NewPlaylist(*channel)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(youoff.Ctx)
	defer cancel()

	// launch pipeline
	videoStream := retrieveVideos(ctx, playlist)
	numWorkers := runtime.NumCPU()
	log.Debugf("youoff will use %v workers!", numWorkers)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		log.Debugf("starting worker %d...", i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			uploadVideo(youoff, ctx, downloadVideo(ctx, videoStream))
		}()
	}
	wg.Wait()
	log.Info("stopping youoff...")

	return nil
}

func retrieveVideos(ctx context.Context, playlist *youtube.Playlist) <-chan youtube.Video {
	videosStream := make(chan youtube.Video)
	go func() {
		defer close(videosStream)
		for playlist.HasNextPage() {
			videos, err := playlist.GetNextVideos()
			if err != nil {
				log.Errorf("error getting next videos from a playlist: %v", err.Error())
			}

			for _, v := range videos {
				select {
				case <-ctx.Done():
					return
				case videosStream <- v:
				}
			}
		}
	}()
	return videosStream
}

func downloadVideo(ctx context.Context, videoStream <-chan youtube.Video) <-chan string {
	pathStream := make(chan string)
	go func() {
		defer close(pathStream)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-videoStream:
				if !ok {
					return
				}
				log.Infof("Downloading video: %v\n", v.Title)
				path, err := v.Download()
				if err != nil {
					log.Errorf("Error downloading video: %v", err.Error())
				}
				select {
				case <-ctx.Done():
					return
				case pathStream <- path:
				}
			}
		}
	}()
	return pathStream
}

func uploadVideo(youoff *YouOff, ctx context.Context, pathStream <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case path, ok := <-pathStream:
			if !ok {
				return
			}
			youoff.Storage.Upload(path, path)
			log.Infof("Uploaded video: %v\n", path)
		}
	}
}

func loadConfig(path string) (*YouOffConfig, error) {
	cfg, err := yaml.NewConfigWithFile(path, ucfg.PathSep("."))
	if err != nil {
		return nil, fmt.Errorf("error loading configuration file: %s", err.Error())
	}

	youOffConfig := &YouOffConfig{}
	err = cfg.Unpack(youOffConfig)
	if err != nil {
		return nil, fmt.Errorf("error unpacking configuration file: %s", err.Error())
	}

	return youOffConfig, nil
}
