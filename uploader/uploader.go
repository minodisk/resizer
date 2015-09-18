package uploader

import (
	"bytes"
	"fmt"
	"os"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/s3"
	"github.com/go-microservices/resizer/log"
	"github.com/go-microservices/resizer/storage"
)

var (
	EnvRegion = "RESIZER_S3_REGION"
	EnvBucket = "RESIZER_S3_BUCKET"
)

type Uploader struct {
	s3.S3
	region string
	bucket string
}

// New はアップローダーを作成する。
func New() (*Uploader, error) {
	region := os.Getenv(EnvRegion)
	if region == "" {
		return nil, fmt.Errorf("requires environment variable: %s", EnvRegion)
	}
	bucket := os.Getenv(EnvBucket)
	if bucket == "" {
		return nil, fmt.Errorf("requires environment variable: %s", EnvBucket)
	}
	return &Uploader{*s3.New(&aws.Config{Region: region}), region, bucket}, nil
}

// Write はデータ buf を filename という名前で contentType と共にアップロードする。
func (self *Uploader) Upload(buf *bytes.Buffer, f storage.Image) (string, error) {
	t := log.Start()
	defer log.End(t)

	data := buf.Bytes()
	input := &s3.PutObjectInput{
		ACL:           aws.String("public-read"),
		Bucket:        aws.String(self.bucket),
		Key:           aws.String(f.Filename),
		Body:          bytes.NewReader(data),
		ContentType:   aws.String(f.ContentType),
		ContentLength: aws.Long(int64(buf.Len())),
	}
	output, err := self.PutObject(input)
	if err != nil {
		log.Printf("aws error: %v", err)
		return "", err
	}

	if expected := fmt.Sprintf("\"%s\"", f.ETag); *output.ETag != expected {
		return "", fmt.Errorf("wrong etag: expected=%s actual=%s", expected, *output.ETag)
	}

	url := self.CreateURL(f.Filename)
	log.Printf("ok: url=%s", url)

	return url, nil
}

func (self *Uploader) CreateURL(path string) string {
	return fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", self.region, self.bucket, path)
}
