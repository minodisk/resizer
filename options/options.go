package options

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	EnvGoogleAuthJSON = "GOOGLE_AUTH_JSON"

	EnvGoogleApplicationCredentials = "GOOGLE_APPLICATION_CREDENTIALS"
	EnvAccount                      = "RESIZER_ACCOUNT"
	EnvBucket                       = "RESIZER_BUCKET"
	EnvConnections                  = "RESIZER_CONNECTIONS"
	EnvDSN                          = "RESIZER_DSN"
	EnvHost                         = "RESIZER_HOST"
	EnvPort                         = "RESIZER_PORT"
	EnvPrefix                       = "RESIZER_PREFIX"
	EnvVerbose                      = "RESIZER_VERBOSE"

	FlagAccount     = "account"
	FlagBucket      = "bucket"
	FlagConnections = "connections"
	FlagDSN         = "dsn"
	FlagHost        = "host"
	FlagPort        = "port"
	FlagPrefix      = "prefix"
	FlagVerbose     = "verbose"
)

var (
	Envs = []string{
		EnvGoogleApplicationCredentials,
		EnvAccount,
		EnvBucket,
		EnvConnections,
		EnvDSN,
		EnvHost,
		EnvPort,
		EnvPrefix,
		EnvVerbose,
	}
	Flags = []string{
		FlagAccount,
		FlagAccount,
		FlagBucket,
		FlagConnections,
		FlagDSN,
		FlagHost,
		FlagPort,
		FlagPrefix,
		FlagVerbose,
	}
	EnvFlagMap = map[string]string{}
)

func init() {
	for i, env := range Envs {
		EnvFlagMap[env] = Flags[i]
	}
}

type Options struct {
	ServiceAccount     ServiceAccount
	Bucket             string
	MaxHTTPConnections int
	DataSourceName     string
	AllowedHosts       Hosts
	Port               int
	ObjectPrefix       string
	Verbose            bool
}

func (o *Options) Parse(args []string) error {
	if v := os.Getenv(EnvGoogleAuthJSON); v != "" {
		b := []byte(v)
		if err := json.Unmarshal(b, &o.ServiceAccount); err != nil {
			return err
		}
		o.ServiceAccount.Path = filepath.Join(os.TempDir(), "resizer-google-auth.json")
		if err := ioutil.WriteFile(o.ServiceAccount.Path, b, 0644); err != nil {
			return err
		}
	}

	fs := flag.NewFlagSet("resizer", flag.ContinueOnError)
	fs.Var(&o.ServiceAccount, "account", `Path to the file of Google service account JSON.`)
	fs.StringVar(&o.Bucket, "bucket", "", `Bucket name of Google Cloud Storage to upload the resized image.`)
	fs.IntVar(&o.MaxHTTPConnections, "connections", 0, `Max simultaneous connections to be accepted by server.
         When 0 or less is specified, the number of connections isn't limited.
         `)
	fs.StringVar(&o.DataSourceName, "dsn", "", `Data source name of database to store resizing information.`)
	fs.Var(&o.AllowedHosts, "host", `Hosts of the image that is allowed to resize.
         When this value isn't specified, all hosts are allowed.
         Multiple hosts can be specified with:
             $ resizer -host a.com,b.com
             $ resizer -host a.com -host b.com`)
	fs.IntVar(&o.Port, "port", 80, `Port to be listened.
         `)
	fs.StringVar(&o.ObjectPrefix, "prefix", "", ``)
	fs.BoolVar(&o.Verbose, "verbose", false, `Verbose output.
         `)
	for _, env := range Envs {
		flag := EnvFlagMap[env]
		if v := os.Getenv(env); v != "" {
			fs.Set(flag, v)
		}
	}
	return fs.Parse(args)
}
