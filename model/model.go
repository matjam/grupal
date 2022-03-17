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

func Update[T any](row T) (*T, error) {
	tt := new(T) // just used to tell gorm what model we're using, hopefully can go away with generics
	err := DB.Model(tt).Updates(row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func Search[T any](filter map[string]any, limit int, offset int, order []string) ([]T, error) {
	var results []T

	// may not work for where IN [list of things]
	q := DB.Where(filter)

	if len(order) == 2 {
		q = q.Order(order[0] + " " + order[1])
	}

	// not sure how a limit of 0 could ever be what you want.
	if limit > 0 {
		q = q.Limit(limit)
	}

	if offset > 0 {
		q = q.Offset(offset)
	}

	err := q.Find(&results).Error
	return results, err
}
