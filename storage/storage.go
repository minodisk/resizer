package storage

import (
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/minodisk/resizer/options"
	"github.com/pkg/errors"
)

type Logger struct{}

func (l Logger) Print(values ...interface{}) {
	log.Println(values...)
}

type Storage struct {
	*gorm.DB
}

func New(o *options.Options) (*Storage, error) {
	var db *gorm.DB
	for {
		var err error
		db, err = gorm.Open("mysql", o.DataSourceName)
		if err == nil {
			break
		}
		log.Println(errors.Wrap(err, "wait for connection"))
		time.Sleep(time.Second)
	}
	db.LogMode(false)
	// db.LogMode(true)
	// db.SetLogger(&Logger{})
	if os.Getenv("ENVIRONMENT") == "development" {
		db.DropTable(&Image{})
	}
	db.CreateTable(&Image{})
	db.AutoMigrate(&Image{})

	for {
		err := db.DB().Ping()
		if err == nil {
			break
		}
		log.Println(errors.Wrap(err, "wait for communication"))
		time.Sleep(time.Second)
	}

	return &Storage{db}, nil
}

func (self *Storage) Close() error {
	return self.DB.DB().Close()
}
