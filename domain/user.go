package domain

import (
	gorm_model "app/domain/model/gorm"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"id"         gorm:"primarykey"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"          gorm:"index"`
}

var UserAllowedSort = []string{"name", "email", "created_at", "updated_at"}

type UserFilter struct {
	gorm_model.DefaultFilter
	ID    *uuid.UUID
	Email *string
}

func (f *UserFilter) Query(q *gorm.DB) {
	// base query
	f.DefaultFilter.Query(q)

	if f.ID != nil {
		q.Where("id = ?", *f.ID)
	}

	if f.Email != nil {
		q.Where("email = ?", *f.Email)
	}
}
