package cloud_storage

import (
	"mime/multipart"
)

// CloudStorage defines the interface for cloud storage operations
type CloudStorage interface {
	Upload(bucketName string, blob multipart.FileHeader, filePath string) error
	Download(bucketName string, filePath string) string
}

// GetCloudStorageByProvider returns the cloud storage object based on the provider
func GetCloudStorageByProvider(provider string) CloudStorage {
	switch provider {
	case "s3":
		s3 := &S3Storage{
			//SecretKey:  cfg.GetConfig().AwsSecretAccessKey,
			//AccessKey:  cfg.GetConfig().AwsAccessKeyId,
			//RegionName: cfg.GetConfig().AwsRegion
		}
		return s3.GetInstance()
	default:
		panic("Unsupported cloud storage provider")
	}
}
