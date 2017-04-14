package storage

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/minodisk/resizer/option"
)

type Logger struct{}

func (l Logger) Print(values ...interface{}) {
	log.Println(values...)
}

type Storage struct {
	*gorm.DB
}

func New(o option.Option) (*Storage, error) {
	db, err := gorm.Open("mysql", o.MysqlDataSourceName)
	if err != nil {
		return nil, err
	}
	db.LogMode(false)
	// db.LogMode(true)
	// db.SetLogger(&Logger{})
	if os.Getenv("ENVIRONMENT") == "development" {
		db.DropTable(&Image{})
	}
	db.CreateTable(&Image{})
	db.AutoMigrate(&Image{})

	return &Storage{db}, nil
}

func (self *Storage) Close() error {
	return self.DB.DB().Close()
}
