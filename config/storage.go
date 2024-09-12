package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	BucketName string
	Region     string
}

var S3Client *s3.Client
var S3Bucket string

func InitS3() {
	awsRegion := os.Getenv("AWS_REGION")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsBucket := os.Getenv("S3_BUCKET_NAME")

	if awsRegion == "" || awsAccessKey == "" || awsSecretKey == "" || awsBucket == "" {
		log.Fatal("AWS credentials missing")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, "")),
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	S3Client = s3.NewFromConfig(cfg)
	S3Bucket = awsBucket
	fmt.Println("S3 Client initialized successfully")
}
