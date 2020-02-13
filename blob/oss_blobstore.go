package blob

import (
	"io"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	aliext "gitlab.com/sdce/exlib/alicloud"
)

type OssBlobStore struct {
	client     *oss.Client
	Expiration time.Duration
}

func NewOssBlobStore(cfg *aliext.AliConfig, Expiration time.Duration) (BlobStore, error) {
	// New client
	client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	return &OssBlobStore{client, Expiration}, nil
}

func (s *OssBlobStore) PullObjectURL(bucket, itemKey string) (string, error) {

	ossBucket, err := s.client.Bucket(bucket)
	if err != nil {
		return "", err
	}

	options := []oss.Option{}
	signedURL, err := ossBucket.SignURL(itemKey, oss.HTTPGet, int64(s.Expiration.Seconds()), options...)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

func (s *OssBlobStore) PutObjectURL(bucket, itemKey string) (string, error) {
	ossBucket, err := s.client.Bucket(bucket)
	if err != nil {
		return "", err
	}

	options := []oss.Option{}
	signedURL, err := ossBucket.SignURL(itemKey, oss.HTTPPut, int64(s.Expiration.Seconds()), options...)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

// to be implemented
func (s *OssBlobStore) PullObject(string, string, io.WriterAt) error {
	return nil
}

func (s *OssBlobStore) PutObject(bucket, itemKey string, file io.Reader) error {
	ossBucket, err := s.client.Bucket(bucket)
	if err != nil {
		return err
	}
	err = ossBucket.PutObject(itemKey, file)
	if err != nil {
		return err
	}
	return nil
}
