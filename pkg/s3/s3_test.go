package s3

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/awstesting/unit"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestGenerateConfigsSuccess(t *testing.T) {
	svc := s3.New(unit.Session)
	svc.Handlers.Send.Clear()
	const payload = `{"schedule": {"foo": "bar"}}`
	svc.Handlers.Send.PushBack(func(r *request.Request) {
		r.HTTPResponse = &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(payload)),
			Header:     http.Header{},
		}
		out := r.Data.(*s3.GetObjectOutput)
		out.Body = io.NopCloser(strings.NewReader(payload))
	})

	plugin := &Plugin{
		client:     svc,
		bucketName: "bucket",
		configPath: "path",
	}

	got, err := plugin.GenerateConfigs(context.Background())
	if err != nil {
		t.Fatalf("GenerateConfigs returned error: %v", err)
	}
	if got[pluginName] != payload {
		t.Fatalf("config mismatch: got %q, want %q", got[pluginName], payload)
	}
}

func TestGenerateConfigsNoSuchKey(t *testing.T) {
	svc := s3.New(unit.Session)
	svc.Handlers.Send.Clear()
	svc.Handlers.Send.PushBack(func(r *request.Request) {
		r.HTTPResponse = &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(strings.NewReader("not found")),
			Header:     http.Header{},
		}
		r.Error = awserr.New(s3.ErrCodeNoSuchKey, "not found", nil)
	})

	plugin := &Plugin{
		client:     svc,
		bucketName: "bucket",
		configPath: "path",
	}

	_, err := plugin.GenerateConfigs(context.Background())
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	want := "object with key path does not exist in bucket bucket"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("unexpected error: got %q, want substring %q", err.Error(), want)
	}
}
