package s3

import (
	"bytes"
	"context"
	"fmt"
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
	bucketName        string
	bucketRegion      string
	configurationPath string
	client            *s3.S3
}

func exitErrorf(msg string, args ...interface{}) {
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

	client := s3.New(session.New(), &aws.Config{Region: aws.String(bucketRegion)})
	plugin := &Plugin{client: client, bucketName: bucketName, configurationPath: configPath}
	return config.NewPlugin(pluginName, plugin.GenerateConfigs)
}

func (p *Plugin) GenerateConfigs(ctx context.Context) (map[string]string, error) {
	result, err := p.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(p.configurationPath),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				exitErrorf("Bucket %s does not exist", p.bucketName)
			case s3.ErrCodeNoSuchKey:
				exitErrorf("Object with key %s does not exist in bucket %s", p.configurationPath, p.bucketName)
			}
		}
		exitErrorf("Unknown error occurred, %v", err)
	}
	defer result.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)

	return map[string]string{pluginName: buf.String()}, nil
}
