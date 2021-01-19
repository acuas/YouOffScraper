package minio

import "github.com/youoffcrawler/storage"

func init() {
	minio := &Minio{}
	err := storage.Register("minio", minio)
	if err != nil {
		panic(err)
	}
}

type Minio struct {
}

func (m *Minio) FileExists(path string) bool {
	return true
}

func (m *Minio) Upload(path string) error {
	return nil
}
