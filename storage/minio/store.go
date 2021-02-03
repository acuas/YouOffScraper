package minio

import (
	"os"

	"github.com/acuas/YouOffScraper/storage"
	ucfg "github.com/elastic/go-ucfg"
	min "github.com/minio/minio-go/v6"
)

func init() {
	err := storage.Register("minio", NewMinio)
	if err != nil {
		panic(err)
	}
}

type minio struct {
	config minioConfig
	client *min.Client
}

func NewMinio(config *ucfg.Config) (storage.Storage, error) {
	minio := &minio{}
	err := config.Unpack(&minio.config)
	if err != nil {
		return nil, err
	}

	minio.client, err = min.New(
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
	_, err := m.client.StatObject(
		m.config.BucketName,
		path,
		min.StatObjectOptions{},
	)
	if err != nil {
		return false
	}
	return true
}

func (m *minio) Upload(name, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}
	_, err = m.client.PutObject(
		m.config.BucketName,
		name,
		file,
		fileStat.Size(),
		min.PutObjectOptions{},
	)
	return err
}
