package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	s3Client *s3.Client
	uploader *manager.Uploader
	bucket   string
)

func init() {
	bucket = os.Getenv("S3_BUCKET_NAME")
	if bucket == "" {
		log.Fatal("S3_BUCKET_NAME is not set")
	}

	var cfg aws.Config
	var err error

	// Use AWS_PROFILE locally if set, otherwise fallback to env variables
	profile := os.Getenv("AWS_PROFILE")
	if profile != "" {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithSharedConfigProfile(profile))
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}

	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	uploader = manager.NewUploader(s3Client)

	fmt.Println("S3 initialized using bucket:", bucket)
}

func UploadFileToS3(ctx context.Context, file multipart.File, fileName string) (string, error) {

	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String("application/octet-stream"),
		ACL:         types.ObjectCannedACLPrivate,
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, fileName)
	return url, nil
}

func GeneratePresignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s3Client)

	req, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(duration))

	if err != nil {
		return "", fmt.Errorf("failed to presign url: %v", err)
	}

	return req.URL, nil
}
