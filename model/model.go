package model

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=localhost user=postgres password=mysecretpassword dbname=grupal " +
	"port=5432 sslmode=disable TimeZone=America/Los_Angeles"

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err = DB.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
}

// Get a given model from the database with a given field/value.
func Get[T any](field string, value string) (T, error) {
	item := new(T)
	err := DB.First(item, field+" = ?", value).Error
	return *item, err
}

// Create a given model from a JSON blob, returning the created T.
func Create[T any](row T) (*T, error) {
	err := DB.Create(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}
