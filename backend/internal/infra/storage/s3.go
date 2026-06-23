package storage

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3         *s3.Client
	presigner  *s3.PresignClient
	bucketName string
	expiry     time.Duration
}

func NewClient(ctx context.Context, region, bucketName string, expirySec int, endpoint, presignEndpoint, accessKeyID, secretAccessKey string) (*Client, error) {
	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	if endpoint != "" {
		ak := accessKeyID
		sk := secretAccessKey
		if ak == "" || sk == "" {
			// Fallback for local development (LocalStack/MinIO)
			ak = "test"
			sk = "test"
		}
		loadOpts = append(loadOpts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(ak, sk, ""),
		))
	}

	cfg, err := config.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		return nil, err
	}

	apiClient := s3.NewFromConfig(cfg, s3EndpointOpts(endpoint))

	presignBase := endpoint
	if presignEndpoint != "" {
		presignBase = presignEndpoint
	}
	presignClient := s3.NewFromConfig(cfg, s3EndpointOpts(presignBase))

	return &Client{
		s3:         apiClient,
		presigner:  s3.NewPresignClient(presignClient),
		bucketName: bucketName,
		expiry:     time.Duration(expirySec) * time.Second,
	}, nil
}

func s3EndpointOpts(endpoint string) func(*s3.Options) {
	if endpoint == "" {
		return func(_ *s3.Options) {}
	}
	return func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true // required for LocalStack
	}
}

func (c *Client) ListObjectKeys(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	paginator := s3.NewListObjectsV2Paginator(c.s3, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucketName),
		Prefix: aws.String(prefix),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, obj := range page.Contents {
			if obj.Key != nil {
				keys = append(keys, *obj.Key)
			}
		}
	}
	return keys, nil
}

func (c *Client) GetPresignedURL(ctx context.Context, key string) (string, error) {
	req, err := c.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(c.expiry))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (c *Client) GetObjectBody(ctx context.Context, key string) ([]byte, error) {
	resp, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
