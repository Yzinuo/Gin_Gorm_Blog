package upload

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	PublicBaseURL   string
}

type R2Store struct {
	client  *s3.Client
	bucket  string
	baseURL string
}

func NewR2Store(_ context.Context, raw struct {
	AccountID       string `mapstructure:"account_id"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Bucket          string `mapstructure:"bucket"`
	PublicBaseURL   string `mapstructure:"public_base_url"`
}) (*R2Store, error) {
	if raw.AccountID == "" || raw.AccessKeyID == "" || raw.SecretAccessKey == "" || raw.Bucket == "" || raw.PublicBaseURL == "" {
		return nil, fmt.Errorf("incomplete R2 configuration")
	}
	client := s3.New(s3.Options{
		Region:       "auto",
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(raw.AccessKeyID, raw.SecretAccessKey, "")),
		BaseEndpoint: aws.String("https://" + raw.AccountID + ".r2.cloudflarestorage.com"),
		UsePathStyle: true,
	})
	return &R2Store{client: client, bucket: raw.Bucket, baseURL: strings.TrimRight(raw.PublicBaseURL, "/")}, nil
}

func (s *R2Store) Put(ctx context.Context, key string, body io.Reader, size int64, opts PutOptions) (ObjectInfo, error) {
	key, err := cleanObjectKey(key)
	if err != nil {
		return ObjectInfo{}, err
	}
	metadata := map[string]string{}
	if opts.SHA256 != "" {
		metadata["sha256"] = opts.SHA256
	}
	output, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket), Key: aws.String(key), Body: body,
		ContentLength: aws.Int64(size), ContentType: aws.String(opts.ContentType),
		CacheControl: aws.String(opts.CacheControl), Metadata: metadata,
	})
	if err != nil {
		return ObjectInfo{}, fmt.Errorf("put R2 object: %w", err)
	}
	return ObjectInfo{Key: key, Size: size, ETag: strings.Trim(aws.ToString(output.ETag), "\""), ContentType: opts.ContentType, CacheControl: opts.CacheControl, SHA256: opts.SHA256}, nil
}

func (s *R2Store) Head(ctx context.Context, key string) (ObjectInfo, error) {
	key, err := cleanObjectKey(key)
	if err != nil {
		return ObjectInfo{}, err
	}
	output, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)})
	if err != nil {
		return ObjectInfo{}, fmt.Errorf("head R2 object: %w", err)
	}
	return ObjectInfo{Key: key, Size: aws.ToInt64(output.ContentLength), ETag: strings.Trim(aws.ToString(output.ETag), "\""), ContentType: aws.ToString(output.ContentType), CacheControl: aws.ToString(output.CacheControl), SHA256: output.Metadata["sha256"]}, nil
}

func (s *R2Store) Delete(ctx context.Context, key string) error {
	key, err := cleanObjectKey(key)
	if err != nil {
		return err
	}
	_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)})
	return err
}

func (s *R2Store) PublicURL(key string) string { return s.baseURL + "/" + strings.TrimLeft(key, "/") }
