package minio

import (
	"testing"

	yaml "github.com/elastic/go-ucfg/yaml"
	min "github.com/minio/minio-go/v6"
)

func TestMinioConfiguration(t *testing.T) {
	conf, err := yaml.NewConfig([]byte(`
endpoint: youoffminio:9000
access_key_id: acces_key_development
secret_access_key: secret_key_development
use_ssl: false
bucket_name: testbucket
bucket_region: ro-south
`))
	if err != nil {
		t.Errorf("yaml invalid format error: %v", err.Error())
	}
	_, err = NewMinio(conf)
	if err != nil {
		t.Errorf("Minio error: %v", err.Error())
	}

	t.Cleanup(func() {
		minioClient, _ := min.New(
			"youoffminio:9000",
			"acces_key_development",
			"secret_key_development",
			false,
		)
		minioClient.RemoveBucket("testbucket")
	})
}

func TestMinioInvalidConfiguration(t *testing.T) {
	conf, err := yaml.NewConfig([]byte(`endpint: testendpoint`))
	if err != nil {
		t.Errorf("yaml invalid format error: %v", err.Error())
	}
	_, err = NewMinio(conf)
	if err == nil {
		t.Errorf("Something wrong in minio configuration!")
	}
}

func TestFileExists(t *testing.T) {
	conf, _ := yaml.NewConfig([]byte(`
endpoint: youoffminio:9000
access_key_id: acces_key_development
secret_access_key: secret_key_development
use_ssl: false
bucket_name: testbucket
bucket_region: ro-south
`))
	m, _ := NewMinio(conf)
	if m.FileExists("test.mp4") != false {
		t.Errorf("Error for checking file absence!")
	}

	// Upload a file and check the existance
	err := m.Upload("config.go", "/go/src/github.com/youoffcrawler/storage/minio/config.go")
	if err != nil {
		t.Errorf("Error uploading file : %s", err.Error())
	}

	if m.FileExists("config.go") != true {
		t.Errorf("Error for checking file existence!")
	}

	t.Cleanup(func() {
		minioClient, _ := min.New(
			"youoffminio:9000",
			"acces_key_development",
			"secret_key_development",
			false,
		)
		minioClient.RemoveObject("testbucket", "config.go")
		minioClient.RemoveBucket("testbucket")
	})
}
