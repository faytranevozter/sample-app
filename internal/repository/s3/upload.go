package s3repo

import (
	storage_model "app/domain/model/storage"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sirupsen/logrus"
)

func (r *S3Repository) UploadFilePublic(objectKey string, body io.Reader, contentType string) (uploadData *storage_model.UploadResponse, err error) {
	_, err = r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(objectKey),
		Body:        body,
		ACL:         types.ObjectCannedACLPublicRead,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		logrus.Errorf(
			"UploadFilePublic: Couldn't upload file to %v:%v. Here's why: %v\n",
			r.bucketName, objectKey, err,
		)
		return
	}

	url := r.GetPublicLink(objectKey)
	if url == "" {
		logrus.Error("UploadFilePublic: GetPublicLink error")
		return
	}

	uploadData = &storage_model.UploadResponse{
		Key:         objectKey,
		ContentType: contentType,
		URL:         url,
	}

	return
}

func (r *S3Repository) UploadFilePrivate(objectKey string, body io.Reader, contentType string, expires *time.Duration) (uploadData *storage_model.UploadResponse, err error) {
	_, err = r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(objectKey),
		Body:        body,
		ACL:         types.ObjectCannedACLPrivate,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		logrus.Errorf(
			"UploadFilePrivate: Couldn't upload file to %v:%v. Here's why: %v\n",
			r.bucketName, objectKey, err,
		)
		return
	}

	url := r.GetPresignedLink(objectKey, expires)
	if url == "" {
		logrus.Error("UploadFilePrivate: GetPresignedLink error")
		return
	}

	uploadData = &storage_model.UploadResponse{
		Key:         objectKey,
		ContentType: contentType,
		URL:         url,
	}

	return
}
