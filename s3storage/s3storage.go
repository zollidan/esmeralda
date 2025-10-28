package s3storage

import (
	"context"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zollidan/esmeralda/config"
)

type S3Storage struct {
	client *s3.Client
}


func New(cfg *config.Config) (*S3Storage, error) {
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
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load S3 config: %w", err)
	}

	return &S3Storage{
		client: s3.NewFromConfig(sdkConfig),
	}, nil
}

func (c *S3Storage) ListItems(ctx context.Context, cfg *config.Config) (*s3.ListObjectsV2Output, error) {
	result, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &cfg.BucketName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in bucket %s: %w", cfg.BucketName, err)
	}
	return result, nil
}
