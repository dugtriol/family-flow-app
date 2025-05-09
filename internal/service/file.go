package service

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	`github.com/aws/aws-sdk-go-v2/config`
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type FileService struct {
	Client           *s3.Client
	BucketName       string
	Region           string
	EndpointResolver string
}

func NewFileService(
	ctx context.Context, backetName, region, endpointResolverURL string,
) *FileService {
	client, err := initS3Client(ctx, region, endpointResolverURL)
	if err != nil {
		return nil
	}

	service := &FileService{
		Client:           client,
		BucketName:       backetName,
		Region:           region,
		EndpointResolver: endpointResolverURL,
	}
	return service
}

func initS3Client(ctx context.Context, regionInput, url string) (*s3.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID && region == regionInput {
				return aws.Endpoint{
					PartitionID:   "yc",
					URL:           url,
					SigningRegion: regionInput,
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		},
	)

	cfg, err := config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(customResolver))

	if err != nil {
		log.Error("FileService - initS3Client - ", err)
	}

	client := s3.NewFromConfig(cfg)
	//cfg, err := config.LoadDefaultConfig(context.TODO())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Создаем клиента для доступа к хранилищу S3
	//client := s3.NewFromConfig(cfg)

	log.Info("Initializing S3 client", "region", regionInput, "url", url)
	log.Info(fmt.Sprintf("Config: %v", cfg))

	//client := s3.NewFromConfig(cfg)
	return client, nil
}

type FileUploadInput struct {
	FileName string
	FileBody []byte
}

func (f *FileService) Test(ctx context.Context) {
	result, err := f.Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Error("Failed to list buckets", "error", err)
	}

	log.Info("Successfully listed buckets")
	for _, bucket := range result.Buckets {
		log.Info(
			"Bucket details",
			"name", aws.ToString(bucket.Name),
			"creation_time", bucket.CreationDate.Format("2006-01-02 15:04:05 Monday"),
		)
	}
}

func (f *FileService) Upload(ctx context.Context, log *slog.Logger, file FileUploadInput) (string, error) {
	log.Info("Starting file upload process")

	//result, err := f.Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	//if err != nil {
	//	log.Error("Failed to list buckets", "error", err)
	//	return "", err
	//}
	//
	//log.Info("Successfully listed buckets")
	//for _, bucket := range result.Buckets {
	//	log.Info(
	//		"Bucket details",
	//		"name", aws.ToString(bucket.Name),
	//		"creation_time", bucket.CreationDate.Format("2006-01-02 15:04:05 Monday"),
	//	)
	//}

	// Запрашиваем список бакетов
	output, err := f.Client.ListObjectsV2(
		context.TODO(), &s3.ListObjectsV2Input{
			Bucket: aws.String(f.BucketName),
		},
	)
	if err != nil {
		log.Error(err.Error())
	}

	log.Info("first page results")
	for _, object := range output.Contents {
		log.Info("key=%s size=%d", aws.ToString(object.Key), *object.Size)
	}
	//result, err := f.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	//if err != nil {
	//	log.Error(err.Error())
	//	return "", err
	//}
	//
	//for _, bucket := range result.Buckets {
	//	log.Info(
	//		"bucket=%s creation time=%s",
	//		aws.ToString(bucket.Name),
	//		bucket.CreationDate.Local().Format("2006-01-02 15:04:05 Monday"),
	//	)
	//}

	name := uuid.New().String()
	ext := filepath.Ext(file.FileName)
	if ext == "" {
		ext = ".jpg"
	}
	pathName := strings.Join([]string{name, ext}, "")
	log.Info("Generated file path", "pathName", pathName)

	_, err = f.Client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(f.BucketName), Key: aws.String(pathName), Body: strings.NewReader(string(file.FileBody)),
		},
	)
	if err != nil {
		log.Error("Failed to upload file", "pathName", pathName, "error", err)
		return "", err
	}

	log.Info("File uploaded successfully", "pathName", pathName)
	return pathName, nil
}

func (f *FileService) Delete(ctx context.Context, path string) (bool, error) {
	_, err := f.Client.DeleteObject(
		ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(f.BucketName),
			Key:    aws.String(path),
		},
	)

	if err != nil {
		log.Error("FileService - Delete - f.Client.DeleteObject - ", err)
		return false, err
	}
	return true, nil
}

func (f *FileService) BuildImageURL(pathName string) string {
	elems := []string{f.EndpointResolver, f.BucketName, pathName}
	path := strings.Join(elems, "/")
	return path
}
