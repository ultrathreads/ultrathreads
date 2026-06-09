package uploader

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"

	"ultrathreads/util/log"
	"ultrathreads/util/urls"
)

type aliyunOssUploader struct {
	once   sync.Once
	bucket *oss.Bucket
}

func NewAliyun() *aliyunOssUploader {
	return &aliyunOssUploader{once: sync.Once{}}
}

func (a *aliyunOssUploader) PutImage(data []byte) (string, error) {
	key := generateImageKey(data)
	return a.PutObject(key, data)
}

func (a *aliyunOssUploader) PutObject(key string, data []byte) (string, error) {
	bucket := a.getBucket()
	if bucket == nil {
		return "", fmt.Errorf("oss bucket not initialized, check config")
	}
	if err := bucket.PutObject(key, bytes.NewReader(data)); err != nil {
		return "", fmt.Errorf("oss put object failed: %w", err)
	}
	return urls.UrlJoin(viper.GetString("uploader.oss.host"), key), nil
}

func (a *aliyunOssUploader) CopyImage(originUrl string) (string, error) {
	data, err := safeDownload(originUrl)
	if err != nil {
		return "", fmt.Errorf("copy image download failed: %w", err)
	}
	return a.PutImage(data)
}

func (a *aliyunOssUploader) getBucket() *oss.Bucket {
	a.once.Do(func() {
		endpoint := viper.GetString("uploader.oss.endpoint")
		accessID := viper.GetString("uploader.oss.access_id")
		accessSecret := viper.GetString("uploader.oss.access_secret")
		bucketName := viper.GetString("uploader.oss.bucket")

		client, err := oss.New(endpoint, accessID, accessSecret)
		if err != nil {
			log.Error("OSS client init failed: %v", err)
			return
		}
		bucket, err := client.Bucket(bucketName)
		if err != nil {
			log.Error("OSS bucket init failed: %v", err)
			return
		}
		a.bucket = bucket
		log.Info("OSS bucket [%s] initialized successfully", bucketName)
	})
	return a.bucket
}