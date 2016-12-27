package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/go-microservices/resizer/option"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Logger struct{}

func (l Logger) Print(values ...interface{}) {
	log.Println(values...)
}

type Storage struct {
	*gorm.DB
}

func New(o option.Options) (*Storage, error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&parseTime=True", o.DBUser, o.DBPassword, o.DBProtocol, o.DBAddress, o.DBName)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// db.LogMode(true)
	db.SetLogger(&Logger{})
	if os.Getenv("ENVIRONMENT") == "develop" {
		db.DropTable(&Image{})
	}
	db.CreateTable(&Image{})
	db.AutoMigrate(&Image{})

	return &Storage{db}, nil
}

func (self *Storage) Close() error {
	return self.DB.DB().Close()
}
