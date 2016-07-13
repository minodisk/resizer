package option

import (
	"github.com/alecthomas/kingpin"
	"os"
	"strings"
)

type Options struct {
	ProjectID  string
	Bucket     string
	JSON       string
	DBUser     string
	DBPassword string
	DBProtocol string
	DBAddress  string
	DBName     string
	Hosts      []string
}

func New(args []string) (o Options, err error) {
	app := kingpin.New("resizer", "Image resizing processor.")

	ProjectID := app.Flag("id", "Project ID of Google Cloud Platform").OverrideDefaultFromEnvar("RESIZER_PROJECT_ID").Required().String()
	Bucket := app.Flag("bucket", "Bucket of Google Cloud Storage").OverrideDefaultFromEnvar("RESIZER_BUCKET").Required().String()
	JSON := app.Flag("json", "Path to json of service account for Google Cloud Platform").OverrideDefaultFromEnvar("RESIZER_JSON").Required().String()
	DBUser := app.Flag("dbuser", "Database user name").OverrideDefaultFromEnvar("RESIZER_DB_USER").Required().String()
	DBPassword := app.Flag("dbpassword", "Database password").OverrideDefaultFromEnvar("RESIZER_DB_PASSWORD").Default("").String()
	DBProtocol := app.Flag("dbprotocol", "Database access protocol").OverrideDefaultFromEnvar("RESIZER_DB_PROTOCOL").Required().String()
	DBAddress := app.Flag("dbaddress", "Database address").OverrideDefaultFromEnvar("RESIZER_DB_ADDRESS").Required().String()
	DBName := app.Flag("dbname", "Database name").OverrideDefaultFromEnvar("RESIZER_DB_NAME").Required().String()
	FlagHosts := app.Flag("host", "Allowed host").Strings()

	_, err = app.Parse(args)
	if err != nil {
		return
	}

	var Hosts *[]string
	if *FlagHosts == nil {
		SplitedHosts := strings.Split(os.Getenv("RESIZER_HOSTS"), ",")
		Hosts = &SplitedHosts
	} else {
		Hosts = FlagHosts
	}

	return Options{
		ProjectID:  *ProjectID,
		Bucket:     *Bucket,
		JSON:       *JSON,
		DBUser:     *DBUser,
		DBPassword: *DBPassword,
		DBProtocol: *DBProtocol,
		DBAddress:  *DBAddress,
		DBName:     *DBName,
		Hosts:      *Hosts,
	}, nil
}
