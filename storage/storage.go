package storage

import (
	"fmt"
	"os"

	"github.com/go-microservices/resizer/log"
	"github.com/go-microservices/resizer/option"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Storage struct {
	*gorm.DB
}

func New(o option.Options) (*Storage, error) {
	t := log.Start()
	defer log.End(t)

	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&parseTime=True", o.DBUser, o.DBPassword, o.DBProtocol, o.DBAddress, o.DBName)

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
