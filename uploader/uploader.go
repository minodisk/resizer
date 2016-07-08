package uploader

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-microservices/resizer/log"
	"github.com/go-microservices/resizer/storage"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	gcs "google.golang.org/api/storage/v1"
)

var (
	EnvGCSProjectID = "RESIZER_GCS_PROJECT_ID"
	EnvGCSBucket    = "RESIZER_GCS_BUCKET"
	EnvGCSJSON      = "RESIZER_GCS_SERVICE_ACCOUNT"
)

const scope = gcs.DevstorageFullControlScope

type Uploader struct {
	service   *gcs.Service
	projectID string
	bucket    string
}

// New はアップローダーを作成する。
func New() (*Uploader, error) {
	bucket := os.Getenv(EnvGCSBucket)
	if bucket == "" {
		return nil, fmt.Errorf("requires environment variable: %s", EnvGCSBucket)
	}
	projectID := os.Getenv(EnvGCSProjectID)
	if projectID == "" {
		return nil, fmt.Errorf("requires environment variable: %s", EnvGCSProjectID)
	}
	jsonPath := os.Getenv(EnvGCSJSON)
	if jsonPath == "" {
		return nil, fmt.Errorf("requires environment variable: %s", EnvGCSJSON)
	}
	jsonFile, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf("Could not open json: %v", err)
	}
	config, err := google.JWTConfigFromJSON(jsonFile, scope)
	if err != nil {
		log.Fatalf("Could not parse json: %v", err)
	}
	client := config.Client(context.Background())
	service, err := gcs.New(client)
	if err != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}
	return &Uploader{service, projectID, bucket}, nil
}

func (self *Uploader) Upload(buf *bytes.Buffer, f storage.Image) (string, error) {
	t := log.Start()
	defer log.End(t)

	object := &gcs.Object{Name: f.Filename}
	if res, err := self.service.Objects.Insert(self.bucket, object).Media(buf).Do(); err == nil {
		fmt.Printf("Created object %v at location %v\n\n", res.Name, res.SelfLink)
	} else {
		log.Fatalf("Objects.Insert failed: %v", err)
	}

	url := self.CreateURL(f.Filename)
	log.Printf("ok: url=%s", url)

	return url, nil
}

func (self *Uploader) CreateURL(path string) string {
	return fmt.Sprintf("https://%s.storage.googleapis.com/%s", self.bucket, path)
}
