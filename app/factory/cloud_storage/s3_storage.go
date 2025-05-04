package cloud_storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"mime/multipart"
	"sync"
	"time"
)

// S3Storage implements CloudStorage for S3
type S3Storage struct {
	singletonInstance *S3Storage

	SecretKey  string
	AccessKey  string
	RegionName string

	// Additional fields for AWS SDK session and S3 client
	svc *s3.Client

	// Mutex for thread safety
	mu sync.Mutex
}

func (s *S3Storage) GetInstance() *S3Storage {

	// Check if the instance already exists
	if s.singletonInstance == nil {
		// Use a mutex to ensure thread safety during instance creation
		s.mu.Lock()
		defer s.mu.Unlock()

		// Check again inside the critical section to prevent race conditions
		if s.singletonInstance == nil {
			// Create the singleton instance
			s.singletonInstance = s
			// Initialize AWS session and S3 client
			s.singletonInstance.init()
		}
	}

	return s.singletonInstance
}

// initAWS initializes the AWS session and S3 client.
func (s *S3Storage) init() {
	// Initialize AWS session
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s.RegionName),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s.AccessKey, s.SecretKey, "")),
	)
	if err != nil {
		fmt.Printf("Error initializing AWS session: %v\n", err)
		return
	}

	fmt.Println("AWS Configuration:", cfg)

	// Create S3 client
	s.svc = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.ClientLogMode = aws.LogSigning | aws.LogRequest | aws.LogResponseWithBody
	})
}

func (s *S3Storage) Upload(bucketName string, blob multipart.FileHeader, filePath string) (err error) {
	_file, err := blob.Open()
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return err
	}
	defer _file.Close()

	// Read the contents of the file into a buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, _file); err != nil {
		fmt.Printf("Error reading file: %v", err)
		return err
	}

	// This uploads the contents of the buffer to S3
	result, err := s.svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return err
	}

	fmt.Println("Uploading", result, "to S3")
	return nil
}

func (s *S3Storage) Download(bucketName string, fileName string) string {
	fmt.Println("aws.String(bucketName) :: %v", aws.String(bucketName))
	fmt.Println("aws.String(fileName) :: %v", aws.String(fileName))
	preSignClient := s3.NewPresignClient(s.svc)
	preSignedUrl, err := preSignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fileName),
		},
		s3.WithPresignExpires(time.Hour*12))
	if err != nil {
		fmt.Println(err)
	}

	return preSignedUrl.URL
}
