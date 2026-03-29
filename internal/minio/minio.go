package minio

import (
	"context"
	"fmt"
	"log"

	"LunoFuns-Server/configs"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

func InitMinIO(cfg *configs.MinIOConfig) error {
	var err error
	Client, err = minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	// 检查桶是否存在，不存在则创建
	ctx := context.Background()
	exists, err := Client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		err = Client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	
	// 无论桶是否刚创建，都强制设置桶策略为公开读
	policy := fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"]}]}`, cfg.BucketName)
	err = Client.SetBucketPolicy(ctx, cfg.BucketName, policy)
	if err != nil {
		log.Printf("Warning: failed to set bucket policy: %v", err)
	}

	log.Printf("MinIO initialized successfully, bucket: %s", cfg.BucketName)
	return nil
}
