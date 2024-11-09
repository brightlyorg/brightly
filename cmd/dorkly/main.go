package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/brightlyorg/brightly/internal/brightly"
	"os"
)

const (
	brightlyYamlEnvVar           = "BRIGHTLY_YAML"
	brightlyEndpointEnvVar       = "BRIGHTLY_ENDPOINT"
	s3BucketEnvVar               = "BRIGHTLY_S3_BUCKET"
	defaultBrightlyYamlInputPath = "project"
)

var logger = brightly.GetLogger()

func main() {
	ctx := context.Background()
	brightlyYamlInputPath := os.Getenv(brightlyYamlEnvVar)
	if brightlyYamlInputPath == "" {
		logger.Debugf("Env var [%s] not set. Using default: %s", brightlyYamlEnvVar, defaultBrightlyYamlInputPath)
		brightlyYamlInputPath = defaultBrightlyYamlInputPath
	}

	brightlyEndpoint := os.Getenv(brightlyEndpointEnvVar)
	if brightlyEndpoint == "" {
		logger.Fatalf("Required env var [%s] not set.", brightlyEndpointEnvVar)
	}

	s3Bucket := os.Getenv(s3BucketEnvVar)
	if s3Bucket == "" {
		logger.Fatalf("Required env var [%s] not set.", s3BucketEnvVar)
	}

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Fatalf("Couldn't load default aws configuration. Have you set up your AWS account? %v", err)
		return
	}

	secretsService := brightly.NewAwsSecretsService(awsConfig)

	logAwsCallerIdentity(awsConfig, ctx)

	s3Client := s3.NewFromConfig(awsConfig)
	s3ArchiveService, err := brightly.NewS3RelayArchiveService(s3Client, s3Bucket)
	if err != nil {
		logger.Fatal(err)
	}
	reconciler := brightly.NewReconciler(s3ArchiveService, secretsService, brightlyYamlInputPath, brightlyEndpoint)

	err = reconciler.Reconcile(ctx)
	if err != nil {
		logger.Fatal(err)
	}
}

func logAwsCallerIdentity(awsConfig aws.Config, ctx context.Context) {
	svc := sts.NewFromConfig(awsConfig)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(ctx, input)
	if err != nil {
		logger.Fatal(err)
	}
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debugf("AWS Identity: %v", string(jsonBytes))
}
