package minio

import (
	"github.com/acuas/YouOffScraper/storage"
	ucfg "github.com/elastic/go-ucfg"
	v6 "github.com/minio/minio-go/v6"
)

func init() {
	err := storage.Register("minio", NewMinio)
	if err != nil {
		panic(err)
	}
}

type minio struct {
	config minioConfig
	client *v6.Client
}

func NewMinio(config *ucfg.Config) (storage.Storage, error) {
	minio := &minio{}
	err := config.Unpack(&minio.config)
	if err != nil {
		return nil, err
	}

	minio.client, err = v6.New(
		minio.config.Endpoint,
		minio.config.AccessKeyID,
		minio.config.SecretAccessKey,
		minio.config.UseSSL,
	)
	if err != nil {
		return nil, err
	}

	found, err := minio.client.BucketExists(minio.config.BucketName)
	if err != nil {
		return nil, err
	}
	if !found {
		err = minio.client.MakeBucket(minio.config.BucketName, minio.config.BucketRegion)
		if err != nil {
			return nil, err
		}
	}

	return minio, nil
}

func (m *minio) FileExists(path string) bool {
	return true
}

func (m *minio) Upload(path string) error {
	return nil
}
