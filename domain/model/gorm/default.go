package gorm_model

import (
	"app/domain/model"
	"time"

	"gorm.io/gorm"
)

type DefaultFilter struct {
	ID  int64
	IDs []int64

	CreatedAtGt    *time.Time
	CreatedAtGte   *time.Time
	CreatedAtLt    *time.Time
	CreatedAtLte   *time.Time
	CreatedAtRange *model.DatetimeRange

	UpdatedAtGt    *time.Time
	UpdatedAtGte   *time.Time
	UpdatedAtLt    *time.Time
	UpdatedAtLte   *time.Time
	UpdatedAtRange *model.DatetimeRange

	Limit  *int64
	Offset *int64
	Sorts  []map[string]string
}

func (f *DefaultFilter) Query(q *gorm.DB) {

	if f.ID > 0 {
		q.Where("id = ?", f.ID)
	} else if len(f.IDs) > 0 {
		q.Where("id IN ?", f.IDs)
	}

	// created at
	if f.CreatedAtGt != nil {
		q.Where("created_at > ?", f.CreatedAtGt)
	} else if f.CreatedAtGte != nil {
		q.Where("created_at >= ?", f.CreatedAtGt)
	}

	if f.CreatedAtLt != nil {
		q.Where("created_at < ?", f.CreatedAtLt)
	} else if f.CreatedAtLte != nil {
		q.Where("created_at <= ?", f.CreatedAtLt)
	}

	if f.CreatedAtRange != nil {
		q.Where("created_at BETWEEN ? AND ?", f.CreatedAtRange.Start, f.CreatedAtRange.End)
	}

	// updated at
	if f.UpdatedAtGt != nil {
		q.Where("updated_at > ?", f.UpdatedAtGt)
	} else if f.UpdatedAtGte != nil {
		q.Where("updated_at >= ?", f.UpdatedAtGt)
	}

	if f.UpdatedAtLt != nil {
		q.Where("updated_at < ?", f.UpdatedAtLt)
	} else if f.UpdatedAtLte != nil {
		q.Where("updated_at <= ?", f.UpdatedAtLt)
	}

	if f.UpdatedAtRange != nil {
		q.Where("updated_at BETWEEN ? AND ?", f.UpdatedAtRange.Start, f.UpdatedAtRange.End)
	}

	// Limit & Offset
	if f.Limit != nil {
		q.Limit(int(*f.Limit))
	}

	if f.Offset != nil {
		q.Offset(int(*f.Offset))
	}

	// Sorts
	if len(f.Sorts) > 0 {
		for _, sort := range f.Sorts {
			for key, value := range sort {
				q.Order(key + " " + value)
			}
		}
	}

}
