package services

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"LunoFuns-Server/configs"
	"LunoFuns-Server/internal/minio"

	"github.com/google/uuid"
	minioGo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

// GeneratePresignedURL 获取简单上传凭证（适用于小文件/封面）
func (s *UploadService) GeneratePresignedURL(filename string, contentType string) (string, string, error) {
	ext := filepath.Ext(filename)
	objectName := fmt.Sprintf("%s/%s%s", time.Now().Format("2006/01/02"), uuid.New().String(), ext)

	reqParams := make(url.Values)
	if contentType != "" {
		reqParams.Set("response-content-type", contentType)
	}

	presignedURL, err := minio.Client.PresignedPutObject(context.Background(), configs.GlobalConfig.MinIO.BucketName, objectName, time.Hour*24)
	if err != nil {
		return "", "", err
	}

	fileURL := fmt.Sprintf("%s/%s/%s", configs.GlobalConfig.MinIO.URLPrefix, configs.GlobalConfig.MinIO.BucketName, objectName)
	return presignedURL.String(), fileURL, nil
}

// 实际上对于分片上传，标准的做法是利用MinIO的core API或者通过后端中转
// 但为了简单和直接的客户端直传，MinIO支持通过Presigned URL进行上传。
// 不过如果是"分片上传"需要初始化uploadId，获取多个part的上传链接，最后合并。

// InitMultipartUpload 初始化分片上传
func (s *UploadService) InitMultipartUpload(filename string, contentType string) (string, string, error) {
    ext := filepath.Ext(filename)
	objectName := fmt.Sprintf("videos/%s/%s%s", time.Now().Format("2006/01/02"), uuid.New().String(), ext)

	core, err := minioGo.NewCore(
		configs.GlobalConfig.MinIO.Endpoint,
		&minioGo.Options{
			Creds:  credentials.NewStaticV4(configs.GlobalConfig.MinIO.AccessKey, configs.GlobalConfig.MinIO.SecretKey, ""),
			Secure: configs.GlobalConfig.MinIO.UseSSL,
		},
	)
	if err != nil {
		return "", "", err
	}

	uploadID, err := core.NewMultipartUpload(context.Background(), configs.GlobalConfig.MinIO.BucketName, objectName, minioGo.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", "", err
	}

	return uploadID, objectName, nil
}

// GetMultipartUploadURLs 获取分片上传链接列表
func (s *UploadService) GetMultipartUploadURLs(objectName, uploadID string, partCount int) ([]string, error) {
	var urls []string
	
	for i := 1; i <= partCount; i++ {
		reqParams := make(url.Values)
		reqParams.Set("partNumber", fmt.Sprintf("%d", i))
		reqParams.Set("uploadId", uploadID)

		presignedURL, err := minio.Client.Presign(context.Background(), "PUT", configs.GlobalConfig.MinIO.BucketName, objectName, time.Hour*24, reqParams)
		if err != nil {
			return nil, err
		}
		urls = append(urls, presignedURL.String())
	}

	return urls, nil
}

// CompleteMultipartUpload 完成分片上传合并
func (s *UploadService) CompleteMultipartUpload(objectName, uploadID string, parts []minioGo.CompletePart) (string, error) {
	core, err := minioGo.NewCore(
		configs.GlobalConfig.MinIO.Endpoint,
		&minioGo.Options{
			Creds:  credentials.NewStaticV4(configs.GlobalConfig.MinIO.AccessKey, configs.GlobalConfig.MinIO.SecretKey, ""),
			Secure: configs.GlobalConfig.MinIO.UseSSL,
		},
	)
	if err != nil {
		return "", err
	}

	_, err = core.CompleteMultipartUpload(context.Background(), configs.GlobalConfig.MinIO.BucketName, objectName, uploadID, parts, minioGo.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("%s/%s/%s", configs.GlobalConfig.MinIO.URLPrefix, configs.GlobalConfig.MinIO.BucketName, objectName)
	return fileURL, nil
}
