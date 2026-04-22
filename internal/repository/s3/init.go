package s3repo

import (
	"context"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

type S3Repository struct {
	bucketName string
	publicURL  *url.URL
	client     *s3.Client
	presigner  *s3.PresignClient
}

func NewS3Repo() *S3Repository {
	cfg, _ := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(
					os.Getenv("S3_ACCESS_KEY"),
					os.Getenv("S3_SECRET_KEY"),
					"",
				),
			),
		),
		config.WithRegion(os.Getenv("S3_REGION")),
	)

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
	})

	publicURL, err := url.Parse(os.Getenv("S3_PUBLIC_URL"))
	if err != nil {
		logrus.Info("s3: without public url", err)
	}

	return &S3Repository{
		bucketName: os.Getenv("S3_BUCKET_NAME"),
		publicURL:  publicURL,
		client:     client,
		presigner:  s3.NewPresignClient(client),
	}
}
