package database

import (
	"fmt"
	"github.com/matjam/grupal/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 	dsn := fmt.Sprintf"host=localhost user=postgres password=mysecretpassword dbname=grupal " +
//		"port=5432 sslmode=disable TimeZone=America/Los_Angeles"

// Model is an opaque struct that is used to mount CRUD functions on
type Model[T any] struct {
	db    *gorm.DB
	model *T
}

type DB struct {
	*gorm.DB

	User *Model[api.User]
}

func NewDB(host, user, password, dbname string, port int, models []any) DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=America/Los_Angeles",
		host, user, password, dbname, port)

	var err error
	var newDB DB
	newDB.User = new(Model[api.User])

	newDB.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	for _, model := range models {
		if err = newDB.AutoMigrate(model); err != nil {
			panic(err)
		}
	}

	return newDB
}

// Create a given database from a JSON blob, returning the created T.
func (m *Model[T]) Create(row T) (*T, error) {
	err := m.db.Create(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// Read from the database
func (m *Model[T]) Read(filter map[string]any, limit int, offset int, order []string) ([]T, error) {
	var results []T

	if order == nil {
		order = []string{}
	}

	// may not work for where IN [list of things]
	q := m.db.Where(filter)

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

func (m *Model[T]) Update(id string, row map[string]any) (*T, error) {
	updatedRow := new(T)

	m.db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Model(updatedRow).Where("id = ?", id).Updates(row).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", id).First(updatedRow).Error; err != nil {
			return err
		}
		return nil
	})

	return updatedRow, nil
}

func (m *Model[T]) Delete(id string) error {
	deletedRow := new(T)
	err := m.db.Where("id = ?", id).Delete(deletedRow).Error
	return err
}
