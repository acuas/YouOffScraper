package minio

type minioConfig struct {
	Endpoint        string `config:"endpoint"`
	AccessKeyID     string `config:"access_key_id"`
	SecretAccessKey string `config:"secret_access_key"`
	UseSSL          bool   `config:"use_ssl"`
	BucketName      string `config:"bucket_name"`
	BucketRegion    string `config:"bucket_region"`
}
