// Package catalog provides S3 storage for AEP files
package catalog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// StorageInterface defines the contract for file storage implementations
type StorageInterface interface {
	// Upload stores a file and returns the storage key
	Upload(ctx context.Context, key string, data io.Reader, metadata map[string]string) (*StorageInfo, error)
	
	// Download retrieves a file by key
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	
	// Delete removes a file
	Delete(ctx context.Context, key string) error
	
	// List returns files matching a prefix
	List(ctx context.Context, prefix string) ([]*StorageInfo, error)
	
	// GetURL returns a presigned URL for download
	GetURL(ctx context.Context, key string, expiry time.Duration) (string, error)
	
	// HealthCheck verifies storage connectivity
	HealthCheck(ctx context.Context) error
}

// StorageInfo contains metadata about stored files
type StorageInfo struct {
	Key         string
	Bucket      string
	Size        int64
	ContentType string
	VersionID   string
	ETag        string
	LastModified time.Time
	Metadata    map[string]string
}

// S3Storage implements StorageInterface using AWS S3
type S3Storage struct {
	client *s3.Client
	bucket string
	prefix string
}

// S3Config holds S3 configuration
type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
	Prefix          string // Optional prefix for all keys
	Endpoint        string // Optional custom endpoint
}

// NewS3Storage creates a new S3 storage instance
func NewS3Storage(cfg S3Config) (*S3Storage, error) {
	// Create AWS config
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	storage := &S3Storage{
		client: s3Client,
		bucket: cfg.Bucket,
		prefix: cfg.Prefix,
	}

	// Verify bucket access
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := storage.HealthCheck(ctx); err != nil {
		return nil, fmt.Errorf("failed to verify S3 access: %w", err)
	}

	log.Printf("Connected to S3 bucket: %s", cfg.Bucket)
	return storage, nil
}

// Upload stores a file in S3
func (s *S3Storage) Upload(ctx context.Context, key string, data io.Reader, metadata map[string]string) (*StorageInfo, error) {
	// Add prefix if configured
	fullKey := s.getFullKey(key)

	// Read data into buffer to get size
	buf := &bytes.Buffer{}
	size, err := io.Copy(buf, data)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	// Prepare metadata
	s3Metadata := make(map[string]string)
	for k, v := range metadata {
		s3Metadata[k] = v
	}
	s3Metadata["uploaded-by"] = "mobot2025"
	s3Metadata["upload-time"] = time.Now().UTC().Format(time.RFC3339)

	// Determine content type
	contentType := "application/octet-stream"
	if strings.HasSuffix(key, ".aep") {
		contentType = "application/vnd.adobe.aftereffects.project"
	}

	// Upload to S3
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(fullKey),
		Body:        buf,
		ContentType: aws.String(contentType),
		Metadata:    s3Metadata,
		// Enable server-side encryption
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	}

	result, err := s.client.PutObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	info := &StorageInfo{
		Key:          key,
		Bucket:       s.bucket,
		Size:         size,
		ContentType:  contentType,
		ETag:         aws.ToString(result.ETag),
		VersionID:    aws.ToString(result.VersionId),
		LastModified: time.Now(),
		Metadata:     metadata,
	}

	log.Printf("Uploaded %s to S3 (size: %d bytes, etag: %s)", fullKey, size, info.ETag)
	return info, nil
}

// Download retrieves a file from S3
func (s *S3Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	fullKey := s.getFullKey(key)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fullKey),
	}

	result, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}

	return result.Body, nil
}

// Delete removes a file from S3
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	fullKey := s.getFullKey(key)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fullKey),
	}

	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	log.Printf("Deleted %s from S3", fullKey)
	return nil
}

// List returns files matching a prefix
func (s *S3Storage) List(ctx context.Context, prefix string) ([]*StorageInfo, error) {
	fullPrefix := s.getFullKey(prefix)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(fullPrefix),
	}

	var files []*StorageInfo
	paginator := s3.NewListObjectsV2Paginator(s.client, input)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		for _, obj := range output.Contents {
			// Remove our prefix from the key
			key := aws.ToString(obj.Key)
			if s.prefix != "" {
				key = strings.TrimPrefix(key, s.prefix+"/")
			}

			files = append(files, &StorageInfo{
				Key:          key,
				Bucket:       s.bucket,
				Size:         aws.ToInt64(obj.Size),
				ETag:         aws.ToString(obj.ETag),
				LastModified: aws.ToTime(obj.LastModified),
			})
		}
	}

	return files, nil
}

// GetURL returns a presigned URL for download
func (s *S3Storage) GetURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	fullKey := s.getFullKey(key)

	presignClient := s3.NewPresignClient(s.client)
	
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fullKey),
	}

	presignResult, err := presignClient.PresignGetObject(ctx, input, func(po *s3.PresignOptions) {
		po.Expires = expiry
	})
	if err != nil {
		return "", fmt.Errorf("failed to create presigned URL: %w", err)
	}

	return presignResult.URL, nil
}

// HealthCheck verifies S3 connectivity
func (s *S3Storage) HealthCheck(ctx context.Context) error {
	// Try to list with a limit of 1 to verify access
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.bucket),
		MaxKeys: aws.Int32(1),
	}

	_, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return fmt.Errorf("S3 health check failed: %w", err)
	}

	return nil
}

// getFullKey adds the configured prefix to a key
func (s *S3Storage) getFullKey(key string) string {
	if s.prefix == "" {
		return key
	}
	return filepath.Join(s.prefix, key)
}

// GetDefaultS3Config returns S3 configuration from environment
func GetDefaultS3Config() S3Config {
	return S3Config{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Region:          os.Getenv("AWS_DEFAULT_REGION"),
		Bucket:          os.Getenv("AWS_BUCKET"),
		Prefix:          "mobot2025", // Default prefix for organization
	}
}

// LocalStorage implements StorageInterface using local filesystem (for development)
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
	}, nil
}

// Upload stores a file locally
func (l *LocalStorage) Upload(ctx context.Context, key string, data io.Reader, metadata map[string]string) (*StorageInfo, error) {
	fullPath := filepath.Join(l.basePath, key)
	
	// Create directory if needed
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	size, err := io.Copy(file, data)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// Metadata storage can be implemented later if needed
	// metaPath := fullPath + ".meta.json"

	info := &StorageInfo{
		Key:          key,
		Bucket:       "local",
		Size:         size,
		LastModified: time.Now(),
		Metadata:     metadata,
	}

	return info, nil
}

// Download retrieves a file locally
func (l *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	fullPath := filepath.Join(l.basePath, key)
	return os.Open(fullPath)
}

// Delete removes a file locally
func (l *LocalStorage) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(l.basePath, key)
	return os.Remove(fullPath)
}

// List returns files matching a prefix
func (l *LocalStorage) List(ctx context.Context, prefix string) ([]*StorageInfo, error) {
	searchPath := filepath.Join(l.basePath, prefix)
	
	var files []*StorageInfo
	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && !strings.HasSuffix(path, ".meta.json") {
			relPath, _ := filepath.Rel(l.basePath, path)
			files = append(files, &StorageInfo{
				Key:          relPath,
				Bucket:       "local",
				Size:         info.Size(),
				LastModified: info.ModTime(),
			})
		}
		
		return nil
	})
	
	return files, err
}

// GetURL returns a file:// URL for local storage
func (l *LocalStorage) GetURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	fullPath := filepath.Join(l.basePath, key)
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}
	return "file://" + absPath, nil
}

// HealthCheck verifies local storage is accessible
func (l *LocalStorage) HealthCheck(ctx context.Context) error {
	testFile := filepath.Join(l.basePath, ".health-check")
	
	// Try to write a test file
	if err := os.WriteFile(testFile, []byte("ok"), 0644); err != nil {
		return fmt.Errorf("failed to write test file: %w", err)
	}
	
	// Clean up
	os.Remove(testFile)
	return nil
}