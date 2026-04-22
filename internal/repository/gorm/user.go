package gormrepo

import (
	"app/domain"
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	db.AutoMigrate(&domain.User{})

	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) StructScan(rows *sql.Rows, dest any) error {
	return r.db.ScanRows(rows, dest)
}

func (r *UserRepository) FetchUser(ctx context.Context, options domain.UserFilter) (cur *sql.Rows, err error) {
	// generate query
	q := r.db.Model(&domain.User{})
	options.Query(q)

	cur, err = q.Rows()
	if err != nil {
		logrus.Error("FetchUser Find:", err)
		return
	}

	return
}

func (r *UserRepository) FetchOneUser(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
	// generate query
	q := r.db.Model(&domain.User{})
	options.Query(q)

	// set row
	row := domain.User{}

	err := q.First(&row).Error
	if err != nil {
		logrus.Error("FetchOneUser Get:", err)
		return nil, err
	}

	return &row, nil
}

func (r *UserRepository) CountUser(ctx context.Context, options domain.UserFilter) (total int64) {
	// generate query
	q := r.db.Model(&domain.User{})
	options.Query(q)

	err := q.Count(&total).Error
	if err != nil {
		logrus.Error("CountUser", err)
		return 0
	}

	return
}

func (r *UserRepository) CreateUser(ctx context.Context, row *domain.User) (err error) {
	// set created at and updated at
	row.CreatedAt = time.Now()
	row.UpdatedAt = time.Now()

	err = r.db.Create(row).Error
	if err != nil {
		logrus.Error("CreateUser Exec:", err)
		return
	}

	return
}
