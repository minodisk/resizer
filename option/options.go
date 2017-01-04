package option

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
)

const (
	AllowedHosts        = "ALLOWED_HOSTS"
	Environment         = "ENVIRONMENT"
	GCProjectID         = "GC_PROJECT_ID"
	GCStorageBucket     = "GC_STORAGE_BUCKET"
	GCServiceAccount    = "GC_SERVICE_ACCOUNT"
	MysqlDataSourceName = "MYSQL_DATA_SOURCE_NAME"
	MaxHTTPConnections  = "MAX_HTTP_CONNECTIONS"
)

type Options map[string]NewOption

func Load(filepath string) (Options, error) {
	o := Options{}

	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return o, err
	}

	if err := yaml.Unmarshal(buf, &o); err != nil {
		return o, err
	}

	fmt.Printf("%+v", o)

	return o, nil
}

func (os Options) Find(rawurl string) (NewOption, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return NewOption{}, err
	}

	if o, ok := os[u.Host]; ok {
		return o, nil
	}
	if o, ok := os["*"]; ok {
		return o, nil
	}

	return NewOption{}, fmt.Errorf("configuration for '%s' is not found", u.Host)
}

type GoogleCloud struct {
	ProjectID         string `json:"project_id",yaml:"project_id"`
	ServiceAccount    string `json:"service_account",yaml:"service_account"`
	StorageBucketName string `json:"storage_bucket_name",yaml:"storage_bucket_name"`
}

type MySQL struct {
	DataSourceName string `json:"data_source_name",yaml:"data_source_name"`
}

type NewOption struct {
	MaxHTTPConnections int         `json:"max_http_connections",yaml:"max_http_connections"`
	GoogleCloud        GoogleCloud `json:"google_cloud",yaml:"google_cloud"`
	MySQL              MySQL       `json:"mysql",yaml:"mysql"`
}

type Option struct {
	AllowedHosts        []string
	GCProjectID         string
	GCStorageBucket     string
	GCServiceAccount    string
	MysqlDataSourceName string
	MaxHTTPConnections  int
}

func New(args []string) (Option, error) {
	env := os.Getenv(Environment)
	if env == "" {
		env = "production"
	}
	maxConn, err := strconv.Atoi(os.Getenv(MaxHTTPConnections))
	return Option{
		AllowedHosts:        strings.Split(os.Getenv(AllowedHosts), ","),
		GCProjectID:         os.Getenv(GCProjectID),
		GCStorageBucket:     os.Getenv(GCStorageBucket),
		GCServiceAccount:    os.Getenv(GCServiceAccount),
		MysqlDataSourceName: os.Getenv(MysqlDataSourceName),
		MaxHTTPConnections:  maxConn,
	}, err
}
