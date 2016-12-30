package option

import (
	"os"
	"strconv"
	"strings"
)

const (
	Environment         = "ENVIRONMENT"
	GCProjectID         = "GC_PROJECT_ID"
	GCStorageBucket     = "GC_STORAGE_BUCKET"
	GCServiceAccount    = "GC_SERVICE_ACCOUNT"
	MysqlDataSourceName = "MYSQL_DATA_SOURCE_NAME"
	AllowedHosts        = "ALLOWED_HOSTS"
	MaxHTTPConnections  = "MAX_HTTP_CONNECTIONS"
)

type Options struct {
	// JSON       string
	// DBUser     string
	// DBPassword string
	// DBProtocol string
	// DBAddress  string
	// DBName     string

	Environment         string
	GCProjectID         string
	GCStorageBucket     string
	GCServiceAccount    string
	MysqlDataSourceName string
	AllowedHosts        []string
	MaxHTTPConnections  int
}

func New(args []string) (Options, error) {
	env := os.Getenv(Environment)
	if env == "" {
		env = "production"
	}
	maxConn, err := strconv.Atoi(os.Getenv(MaxHTTPConnections))
	return Options{
		Environment:         env,
		GCProjectID:         os.Getenv(GCProjectID),
		GCStorageBucket:     os.Getenv(GCStorageBucket),
		GCServiceAccount:    os.Getenv(GCServiceAccount),
		MysqlDataSourceName: os.Getenv(MysqlDataSourceName),
		AllowedHosts:        strings.Split(os.Getenv(AllowedHosts), ","),
		MaxHTTPConnections:  maxConn,
	}, err

	// app := kingpin.New("resizer", "Image resizing processor.")
	// ProjectID := app.Flag("id", "Project ID of Google Cloud Platform").OverrideDefaultFromEnvar("RESIZER_PROJECT_ID").Required().String()
	// Bucket := app.Flag("bucket", "Bucket of Google Cloud Storage").OverrideDefaultFromEnvar("RESIZER_BUCKET").Required().String()
	// JSON := app.Flag("json", "Path to json of service account for Google Cloud Platform").OverrideDefaultFromEnvar("RESIZER_JSON").Required().String()
	// DBUser := app.Flag("dbuser", "Database user name").OverrideDefaultFromEnvar("RESIZER_DB_USER").Required().String()
	// DBPassword := app.Flag("dbpassword", "Database password").OverrideDefaultFromEnvar("RESIZER_DB_PASSWORD").Default("").String()
	// DBProtocol := app.Flag("dbprotocol", "Database access protocol").OverrideDefaultFromEnvar("RESIZER_DB_PROTOCOL").Required().String()
	// DBAddress := app.Flag("dbaddress", "Database address").OverrideDefaultFromEnvar("RESIZER_DB_ADDRESS").Required().String()
	// DBName := app.Flag("dbname", "Database name").OverrideDefaultFromEnvar("RESIZER_DB_NAME").Required().String()
	// FlagHosts := app.Flag("host", "Allowed host").Strings()
	// MaxConn := app.Flag("maxconn", "Max number of current connections").OverrideDefaultFromEnvar("RESIZER_MAX_CONNECTION").Default("10").Int()
	// _, err = app.Parse(args)
	// if err != nil {
	// 	return
	// }
	//
	// var Hosts *[]string
	// if *FlagHosts == nil {
	// 	SplitedHosts := strings.Split(os.Getenv("RESIZER_HOSTS"), ",")
	// 	Hosts = &SplitedHosts
	// } else {
	// 	Hosts = FlagHosts
	// }
	//
	// return Options{
	// 	GCProjectID:     *ProjectID,
	// 	GCStorageBucket: *Bucket,
	// 	JSON:            *JSON,
	// 	DBUser:          *DBUser,
	// 	DBPassword:      *DBPassword,
	// 	DBProtocol:      *DBProtocol,
	// 	DBAddress:       *DBAddress,
	// 	DBName:          *DBName,
	// 	AllowedHosts:    *Hosts,
	// 	MaxConn:         *MaxConn,
	// }, nil
}
