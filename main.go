package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
	"go.uber.org/zap"

	"github.com/sdn0303/golog"
	"github.com/sdn0303/transcoder/infra/mediaconverter"
	"github.com/sdn0303/transcoder/infra/storage"
)

var logger *golog.Logger

const (
	jobSettingBucket = "job-settings"
	jobSettingKey    = "job.json"
)

func generateS3Path(bucketName, key string) string {
	return fmt.Sprintf("s3://%s/%s", bucketName, key)
}

func loadJobSettings(inputKey, optKey string) (*mediaconvert.JobSettings, error) {

	s3 := storage.New()
	wb, err := s3.Downloader(jobSettingBucket, jobSettingKey)
	if err != nil {
		logger.Error("downloader failed", zap.Error(err))
		return nil, err
	}

	logger.Info("Start to Create Job settings")

	var js *mediaconvert.JobSettings
	if err := json.Unmarshal(wb.Bytes(), &js); err != nil {
		logger.Error("failed to unmarshal job settings", zap.Error(err))
		return nil, err
	}

	js.Inputs[0].FileInput = aws.String(inputKey)
	js.OutputGroups[0].OutputGroupSettings.HlsGroupSettings.Destination = aws.String(optKey)
	js.OutputGroups[1].OutputGroupSettings.FileGroupSettings.Destination = aws.String(optKey)

	return js, nil
}

func handler(ctx context.Context, event events.S3Event) error {

	bucketName := event.Records[0].S3.Bucket.Name
	key := event.Records[0].S3.Object.Key

	inputKey := generateS3Path(bucketName, key)
	optKey := generateS3Path(os.Getenv("HLS_BUCKET"), key)

	js, err := loadJobSettings(inputKey, optKey)
	if err != nil {
		logger.Error("", zap.Error(err))
		return err
	}

	// Create MediaConvert job
	mc := mediaconverter.New()
	opt, err := mc.CreateJob(js)
	if err != nil {
		logger.Error("failed to create mediaconvert job", zap.Error(err))
		return err
	}

	logger.Info("Job Created")
	logger.Info(opt.String())

	return nil
}

func main() {
	logger = golog.GetInstance()
	lambda.Start(handler)
}
