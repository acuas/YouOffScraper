package minio

import (
	"testing"

	yaml "github.com/elastic/go-ucfg/yaml"
)

func TestConfiguration(t *testing.T) {
	conf, err := yaml.NewConfig([]byte(`
endpoint: 'localhost:9000'
access_key_id: acces_key_development
secret_access_key: secret_key_development
use_ssl: false
bucket_name: youtube
bucket_region: ro-south
`))
	if err != nil {
		t.Errorf("Error configuration of minio!")
	}
	_, err = NewMinio(conf)
	if err != nil {
		t.Errorf("Minio error: %v", err.Error())
	}
}
