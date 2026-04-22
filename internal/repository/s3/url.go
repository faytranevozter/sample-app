package s3repo

import (
	"context"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

func (r *S3Repository) GetPresignedLink(objectKey string, expires *time.Duration) string {
	resSigned, err := r.presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		if expires != nil {
			opts.Expires = *expires
		}
	})

	if err != nil {
		logrus.Error("GetPresignedLink error: ", err)
		return ""
	}

	return resSigned.URL
}

func (r *S3Repository) GetPublicLink(objectKey string) string {
	url := &url.URL{}
	if r.publicURL == nil {
		endpoint := r.client.Options().BaseEndpoint
		url, _ = url.Parse(*endpoint)
	} else {
		url = r.publicURL
	}

	// add path with object key
	url.Path = objectKey

	return url.String()
}
