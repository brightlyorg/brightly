package brightly

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsService interface {
	fmt.Stringer
	getSdkKey(ctx context.Context, project, env string) (string, error)
	getMobileKey(ctx context.Context, project, env string) (string, error)
}

var _ SecretsService = &awsSecretsService{}

type awsSecretsService struct {
	awsConfig aws.Config
}

func (s awsSecretsService) String() string {
	return "awsSecretsService"
}

func NewAwsSecretsService(awsConfig aws.Config) SecretsService {
	return &awsSecretsService{
		awsConfig: awsConfig,
	}
}

func (s awsSecretsService) getSdkKey(ctx context.Context, project, env string) (string, error) {
	return s.getAwsSecret(ctx, sdkKeySecretName(project, env))
}

func (s awsSecretsService) getMobileKey(ctx context.Context, project, env string) (string, error) {
	return s.getAwsSecret(ctx, mobileKeySecretName(project, env))
}

func (s awsSecretsService) getAwsSecret(ctx context.Context, secretId string) (string, error) {
	logger.Infof("Fetching AWS secret: %s", secretId)
	svc := secretsmanager.NewFromConfig(s.awsConfig)
	input := &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretId)}
	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		return "", err
	}
	return *result.SecretString, nil
}

// Secrets naming convention must be kept in sync with the terraform bits here:
// https://github.com/brightlyorg/terraform-aws-brightly-flags/blob/main/brightly_environment/main.tf
func sdkKeySecretName(project, env string) string {
	return "brightly-" + project + "-" + env + "-sdk-key"
}

// Secrets naming convention must be kept in sync with the terraform bits here:
// https://github.com/brightlyorg/terraform-aws-brightly-flags/blob/main/brightly_environment/main.tf
func mobileKeySecretName(project, env string) string {
	return "brightly-" + project + "-" + env + "-mob-key"
}
