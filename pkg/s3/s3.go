package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

// S3API is the subset of the S3 client API used by this plugin.
type S3API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type Plugin struct {
	bucketName string
	configPath string
	client     S3API
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

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(bucketRegion),
	)
	if err != nil {
		exitErrorf("error loading AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	plugin := &Plugin{
		client:     client,
		bucketName: bucketName,
		configPath: configPath,
	}
	return config.NewPlugin(pluginName, plugin.GenerateConfigs)
}

func (p *Plugin) GenerateConfigs(ctx context.Context) (map[string]string, error) {
	input := &s3.GetObjectInput{
		Bucket: &p.bucketName,
		Key:    &p.configPath,
	}

	result, err := p.client.GetObject(ctx, input)
	if err != nil {
		if _, ok := errors.AsType[*types.NoSuchBucket](err); ok {
			return nil, fmt.Errorf("bucket %s does not exist", p.bucketName)
		}
		if _, ok := errors.AsType[*types.NoSuchKey](err); ok {
			return nil, fmt.Errorf("object with key %s does not exist in bucket %s", p.configPath, p.bucketName)
		}
		return nil, fmt.Errorf("error fetching S3 object: %w", err)
	}
	defer func() { _ = result.Body.Close() }()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("reading S3 object body: %w", err)
	}

	return map[string]string{pluginName: string(data)}, nil
}
