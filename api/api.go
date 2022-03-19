package api

// package api defines the public interface as well as all the structures used by the service.

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Name     string     `json:"name"`
	Email    *string    `json:"email" gorm:"uniqueIndex"`
	Password *string    `json:"password"`
	Birthday *time.Time `json:"birthday"`
}
