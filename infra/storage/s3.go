package storage

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	Session *session.Session
}

func New() *S3 {
	return &S3{
		Session: session.Must(session.NewSession(&aws.Config{
			Region: aws.String(os.Getenv("REGION")),
		})),
	}
}

func (storage *S3) Downloader(bucket, key string) (*aws.WriteAtBuffer, error) {
	downloader := s3manager.NewDownloader(storage.Session)

	wb := &aws.WriteAtBuffer{}
	_, err := downloader.Download(wb, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	return wb, nil
}
