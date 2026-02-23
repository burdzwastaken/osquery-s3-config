package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithy "github.com/aws/smithy-go"
)

type mockS3Client struct {
	output *s3.GetObjectOutput
	err    error
}

func (m *mockS3Client) GetObject(_ context.Context, _ *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m.output, m.err
}

func TestGenerateConfigsSuccess(t *testing.T) {
	payload := `{"options":{}}`
	mock := &mockS3Client{
		output: &s3.GetObjectOutput{
			Body: io.NopCloser(bytes.NewBufferString(payload)),
		},
	}

	p := &Plugin{
		client:     mock,
		bucketName: "test-bucket",
		configPath: "osquery.conf",
	}

	configs, err := p.GenerateConfigs(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, ok := configs[pluginName]
	if !ok {
		t.Fatalf("expected key %q in configs", pluginName)
	}
	if got != payload {
		t.Errorf("expected %q, got %q", payload, got)
	}
}

func TestGenerateConfigsNoSuchKey(t *testing.T) {
	mock := &mockS3Client{
		err: &types.NoSuchKey{
			Message: new(string),
		},
	}

	p := &Plugin{
		client:     mock,
		bucketName: "test-bucket",
		configPath: "osquery.conf",
	}

	_, err := p.GenerateConfigs(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := fmt.Sprintf("object with key %s does not exist in bucket %s", p.configPath, p.bucketName)
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestGenerateConfigsNoSuchBucket(t *testing.T) {
	mock := &mockS3Client{
		err: &types.NoSuchBucket{
			Message: new(string),
		},
	}

	p := &Plugin{
		client:     mock,
		bucketName: "test-bucket",
		configPath: "osquery.conf",
	}

	_, err := p.GenerateConfigs(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := fmt.Sprintf("bucket %s does not exist", p.bucketName)
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestGenerateConfigsUnknownError(t *testing.T) {
	mock := &mockS3Client{
		err: &smithy.GenericAPIError{
			Code:    "InternalError",
			Message: "something went wrong",
		},
	}

	p := &Plugin{
		client:     mock,
		bucketName: "test-bucket",
		configPath: "osquery.conf",
	}

	_, err := p.GenerateConfigs(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
