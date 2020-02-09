package mediaconverter

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
)

type MediaConverter struct {
	Client *mediaconvert.MediaConvert
}

func New() *MediaConverter {
	return &MediaConverter{
		Client: mediaconvert.New(session.Must(session.NewSession(&aws.Config{
			Region:   aws.String(os.Getenv("REGION")),
			Endpoint: aws.String(os.Getenv("MEDIA_CONVERT_ENDPOINT")),
		}))),
	}
}

func (mc *MediaConverter) CreateJob(settings *mediaconvert.JobSettings) (*mediaconvert.CreateJobOutput, error) {
	return mc.Client.CreateJob(&mediaconvert.CreateJobInput{
		JobTemplate: aws.String(os.Getenv("JOB_TEMPLATE")),
		Queue:       aws.String(os.Getenv("QUEUE_ARN")),
		Role:        aws.String(os.Getenv("ROLE_ARN")),
		Settings:    settings,
	})
}
