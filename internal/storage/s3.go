package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	"github.com/aminbenmansour/terrakube-go/internal/config"
)

// S3 stores objects in any S3-compatible bucket (AWS, MinIO, Ceph, ...).
type S3 struct {
	client *s3.Client
	bucket string
}

// NewS3 builds an S3 client. When explicit credentials are not configured the
// default AWS credential chain is used, which covers IAM roles and workload
// identity.
func NewS3(ctx context.Context, cfg config.S3) (*S3, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("storage: s3 bucket must not be empty")
	}

	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.Region),
	}
	if cfg.AccessKeyID != "" || cfg.SecretAccessKey != "" {
		opts = append(opts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("storage: load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = cfg.UsePathStyle
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	return &S3{client: client, bucket: cfg.Bucket}, nil
}

func (s *S3) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	clean, err := cleanKey(key)
	if err != nil {
		return err
	}

	// The AWS SDK retries by rewinding the body, so spool to a temporary file
	// instead of holding the whole object in memory.
	tmp, err := os.CreateTemp("", "terrakube-s3-*")
	if err != nil {
		return fmt.Errorf("storage: create s3 spool file: %w", err)
	}
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}()

	written, err := io.Copy(tmp, r)
	if err != nil {
		return fmt.Errorf("storage: spool s3 body: %w", err)
	}
	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("storage: rewind s3 body: %w", err)
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(clean),
		Body:   tmp,
	}
	if size >= 0 && size == written {
		input.ContentLength = aws.Int64(written)
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	if _, err := s.client.PutObject(ctx, input); err != nil {
		return fmt.Errorf("storage: put s3 object %s: %w", clean, err)
	}
	return nil
}

func (s *S3) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	clean, err := cleanKey(key)
	if err != nil {
		return nil, err
	}
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(clean),
	})
	if err != nil {
		if isS3NotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("storage: get s3 object %s: %w", clean, err)
	}
	return out.Body, nil
}

func (s *S3) Delete(ctx context.Context, key string) error {
	clean, err := cleanKey(key)
	if err != nil {
		return err
	}
	if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(clean),
	}); err != nil {
		return fmt.Errorf("storage: delete s3 object %s: %w", clean, err)
	}
	return nil
}

func (s *S3) Stat(ctx context.Context, key string) (Object, error) {
	clean, err := cleanKey(key)
	if err != nil {
		return Object{}, err
	}
	out, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(clean),
	})
	if err != nil {
		if isS3NotFound(err) {
			return Object{}, ErrNotFound
		}
		return Object{}, fmt.Errorf("storage: stat s3 object %s: %w", clean, err)
	}

	obj := Object{Key: clean}
	if out.ContentLength != nil {
		obj.Size = *out.ContentLength
	}
	if out.ContentType != nil {
		obj.ContentType = *out.ContentType
	}
	if out.LastModified != nil {
		obj.ModifiedAt = *out.LastModified
	}
	return obj, nil
}

func (s *S3) Health(ctx context.Context) error {
	if _, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(s.bucket)}); err != nil {
		return fmt.Errorf("storage: s3 bucket %s: %w", s.bucket, err)
	}
	return nil
}

func isS3NotFound(err error) bool {
	var noKey *types.NoSuchKey
	if errors.As(err, &noKey) {
		return true
	}
	var notFound *types.NotFound
	if errors.As(err, &notFound) {
		return true
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "NoSuchKey", "NotFound", "NoSuchBucket":
			return true
		}
	}
	return false
}

var _ Backend = (*S3)(nil)
