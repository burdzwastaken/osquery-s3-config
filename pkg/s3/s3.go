package s3

import (
	"bytes"
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kolide/osquery-go/plugin/config"
	"github.com/pkg/errors"
)

type Plugin struct {
	bucketName        string
	bucketRegion      string
	configurationPath string
	client            *s3.S3
}

func New() *config.Plugin {
	bucketName := os.Getenv("OSQUERY_S3_CONFIG_BUCKET_NAME")
	bucketRegion := os.Getenv("OSQUERY_S3_CONFIG_BUCKET_REGION")
	configurationPath := os.Getenv("OSQUERY_S3_CONFIG_PATH")
	if len(configurationPath) == 0 {
		configurationPath = "osquery.conf"
	}

	client := s3.New(session.New(), &aws.Config{Region: aws.String(bucketRegion)})
	plugin := &Plugin{client: client, bucketName: bucketName, configurationPath: configurationPath}
	return config.NewPlugin("s3", plugin.GenerateConfigs)
}

func (p *Plugin) GenerateConfigs(ctx context.Context) (map[string]string, error) {
	result, err := p.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(p.configurationPath),
	})
	if err != nil {
		errors.Wrap(err, "get config s3")
	}
	defer result.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)

	return map[string]string{"s3": buf.String()}, nil
}
