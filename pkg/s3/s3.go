package s3

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/osquery/osquery-go/plugin/config"
)

const (
	pluginName = "s3"

	bucketNameEnvVar       = "OSQUERY_S3_CONFIG_BUCKET_NAME"
	bucketRegionEnvVar     = "OSQUERY_S3_CONFIG_BUCKET_REGION"
	bucketConfigPathEnvVar = "OSQUERY_S3_CONFIG_PATH"

	defaultBucketRegion = "us-east-1"
	defaultConfigPath   = "osquery.conf"
)

type Plugin struct {
	bucketName   string
	bucketRegion string
	configPath   string
	client       *s3.S3
}

func exitErrorf(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func New() *config.Plugin {
	bucketName := os.Getenv(bucketNameEnvVar)
	if len(bucketName) == 0 {
		exitErrorf("%s is required", bucketNameEnvVar)
	}

	bucketRegion := getEnv(bucketRegionEnvVar, defaultBucketRegion)
	configPath := getEnv(bucketConfigPathEnvVar, defaultConfigPath)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(bucketRegion),
	})
	if err != nil {
		exitErrorf("error creating AWS session: %v", err)
	}

	client := s3.New(sess)
	plugin := &Plugin{
		client:       client,
		bucketName:   bucketName,
		bucketRegion: bucketRegion,
		configPath:   configPath,
	}
	return config.NewPlugin(pluginName, plugin.GenerateConfigs)
}

func (p *Plugin) GenerateConfigs(ctx context.Context) (map[string]string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(p.configPath),
	}

	result, err := p.client.GetObjectWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, fmt.Errorf("bucket %s does not exist", p.bucketName)
			case s3.ErrCodeNoSuchKey:
				return nil, fmt.Errorf("object with key %s does not exist in bucket %s", p.configPath, p.bucketName)
			default:
				return nil, fmt.Errorf("AWS S3 error: %s - %s", aerr.Code(), aerr.Message())
			}
		}
		return nil, fmt.Errorf("error fetching S3 object: %w", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("reading S3 object body: %w", err)
	}

	return map[string]string{pluginName: string(data)}, nil
}
