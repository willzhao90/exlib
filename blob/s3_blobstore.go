package blob

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	awsext "gitlab.com/sdce/exlib/aws"
)

type AwsBlobStore struct {
	S          *session.Session
	Expiration time.Duration
}

//NewAwsBlobStore construct a new AwsBlobStore
func NewAwsBlobStore(c *awsext.AwsConfig, urlExpirationTime time.Duration) (BlobStore, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.Region),
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")})
	if err != nil {
		return nil, fmt.Errorf("Failed to create a new session %v", err)
	}
	return &AwsBlobStore{sess, urlExpirationTime}, nil
}

//PullObjectURL offer a url for pulling file from s3
func (s *AwsBlobStore) PullObjectURL(bucket, itemKey string) (string, error) {
	// Create S3 service client
	svc := s3.New(s.S)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(itemKey),
	})

	//for otc chat items hardcode a longer expiration
	var expire time.Duration
	if strings.Contains(itemKey, "otc") {
		expire = time.Hour * 100
	} else {
		expire = s.Expiration
	}

	urlStr, err := req.Presign(expire)
	if err != nil {
		log.Println("Failed to sign request", err)
		return "", err
	}
	log.Println("The URL is", urlStr)
	return urlStr, nil
}

//PutObjectURL offer a url for putting file to s3
func (s *AwsBlobStore) PutObjectURL(bucket, key string) (string, error) {
	svc := s3.New(s.S)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	str, err := req.Presign(s.Expiration)
	if err != nil {
		log.Println("Failed to sign request", err)
		return "", err
	}
	log.Println("The URL is:", str)
	return str, nil
}

//PullObject full a file from s3
func (s *AwsBlobStore) PullObject(bucket, key string, writer io.WriterAt) error {
	downloader := s3manager.NewDownloader(s.S)

	numBytes, err := downloader.Download(writer,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		log.Printf("Unable to download item %q, %v", key, err)
		return err
	}
	log.Println("Downloaded", numBytes, "bytes")
	return nil
}

//PutObject put a file to s3
func (s *AwsBlobStore) PutObject(bucket, key string, file io.Reader) error {
	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploader(s.S)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Printf("Unable to upload %q to %q, %v", key, bucket, err)
	}
	log.Printf("Successfully uploaded %q to %q\n", key, bucket)
	return err
}
