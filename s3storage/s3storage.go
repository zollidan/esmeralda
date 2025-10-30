package s3storage

import (
	"bytes"
	"context"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zollidan/esmeralda/config"
)

type S3Storage struct {
	client *s3.Client
	Bucket string
}

func New(cfg *config.Config) *S3Storage {
	sdkConfig, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.S3AccessKey,
				cfg.S3SecretKey,
				"",
			),
		),
		awsConfig.WithBaseEndpoint(cfg.S3URL),
		awsConfig.WithRegion(cfg.S3Region),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to load S3 config: %s", err.Error()))
	}

	return &S3Storage{
		client: s3.NewFromConfig(sdkConfig),
		Bucket: cfg.BucketName,
	}
}

func (c *S3Storage) UploadItem(ctx context.Context, key string, body []byte) error {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &c.Bucket,
		Key:    &key,
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		return fmt.Errorf("failed to upload object %s: %w", key, err)
	}

	return nil
}

func (c *S3Storage) ListItems(ctx context.Context) (*s3.ListObjectsV2Output, error) {
	result, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &c.Bucket,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in bucket %s: %w", c.Bucket, err)
	}
	return result, nil
}

func (c *S3Storage) DeleteItem(ctx context.Context, key string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &c.Bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", key, err)
	}
	return nil
}

func (c *S3Storage) GetItem(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	result, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object %s: %w", key, err)
	}
	return result, nil
}
