package storage

import (
	"fmt"
	"os"

	"github.com/go-microservices/resizer/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	EnvUsername = "RESIZER_DB_USERNAME"
	EnvPassword = "RESIZER_DB_PASSWORD"
	EnvProtocol = "RESIZER_DB_PROTOCOL"
	EnvAddress  = "RESIZER_DB_ADDRESS"
	EnvName     = "RESIZER_DB_NAME"
)

type Storage struct {
	*gorm.DB
}

func New() (*Storage, error) {
	t := log.Start()
	defer log.End(t)

	username := os.Getenv(EnvUsername)
	if username == "" {
		return nil, fmt.Errorf("requires environment variable: %s", EnvUsername)
	}
	password := os.Getenv(EnvPassword)
	protocol := os.Getenv(EnvProtocol)
	if protocol == "" {
		return nil, fmt.Errorf("requires environment variable: %s", EnvProtocol)
	}
	address := os.Getenv(EnvAddress)
	if address == "" {
		return nil, fmt.Errorf("requires environment variable: %s", EnvAddress)
	}
	name := os.Getenv(EnvName)
	if name == "" {
		return nil, fmt.Errorf("requires environment variable: %s", EnvName)
	}
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&parseTime=True", username, password, protocol, address, name)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	db.SetLogger(&log.Gorm{})
	if os.Getenv("ENVIRONMENT") == "develop" {
		db.DropTable(&Image{})
	}
	db.CreateTable(&Image{})
	db.AutoMigrate(&Image{})

	return &Storage{db}, nil
}

func (self *Storage) Close() error {
	t := log.Start()
	defer log.End(t)

	return self.DB.DB().Close()
}
