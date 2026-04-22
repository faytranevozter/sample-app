package mongo_model

import (
	"app/domain/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

type DefaultFilter struct {
	ID     primitive.ObjectID
	IDStr  *string
	IDs    []primitive.ObjectID
	IDsStr []string

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

	Raw map[string]any

	Limit  *int64
	Offset *int64
	Sorts  bson.D
}

func (f *DefaultFilter) Query(query map[string]any) {
	if !f.ID.IsZero() {
		query["_id"] = f.ID
	} else if f.IDStr != nil {
		obj, _ := primitive.ObjectIDFromHex(*f.IDStr)
		query["_id"] = obj
	}

	if len(f.IDs) > 0 {
		query["_id"] = map[string]any{
			"$in": f.IDs,
		}
	} else if len(f.IDsStr) > 0 {
		objIDs := make([]primitive.ObjectID, 0)
		for _, id := range f.IDsStr {
			if obID, err := primitive.ObjectIDFromHex(id); err == nil {
				objIDs = append(objIDs, obID)
			}
		}
		query["_id"] = objIDs
	}

	// created at
	createdAt := make(map[string]any)
	if f.CreatedAtGt != nil {
		createdAt["$gt"] = f.CreatedAtGt
	} else if f.CreatedAtGte != nil {
		createdAt["$gte"] = f.CreatedAtGt
	}

	if f.CreatedAtLt != nil {
		createdAt["$lt"] = f.CreatedAtLt
	} else if f.CreatedAtLte != nil {
		createdAt["$lte"] = f.CreatedAtLte
	}

	if f.CreatedAtRange != nil {
		createdAt["$gte"] = f.CreatedAtRange.Start
		createdAt["$lte"] = f.CreatedAtRange.End
	}

	if len(createdAt) > 0 {
		query["createdAt"] = createdAt
	}

	// updated at
	updatedAt := make(map[string]any)
	if f.UpdatedAtGt != nil {
		updatedAt["$gt"] = f.UpdatedAtGt
	} else if f.UpdatedAtGte != nil {
		updatedAt["$gte"] = f.UpdatedAtGt
	}

	if f.UpdatedAtLt != nil {
		updatedAt["$lt"] = f.UpdatedAtLt
	} else if f.UpdatedAtLte != nil {
		updatedAt["$lte"] = f.UpdatedAtLte
	}

	if f.UpdatedAtRange != nil {
		updatedAt["$gte"] = f.UpdatedAtRange.Start
		updatedAt["$lte"] = f.UpdatedAtRange.End
	}

	if len(updatedAt) > 0 {
		query["updatedAt"] = updatedAt
	}

	// raw data
	if len(f.Raw) > 0 {
		for k, v := range f.Raw {
			query[k] = v
		}
	}
}

func (f *DefaultFilter) FindOptions() *moptions.FindOptions {
	findOptions := &moptions.FindOptions{}

	if f.Limit != nil {
		findOptions.SetLimit(*f.Limit)
	}

	if f.Offset != nil {
		findOptions.SetSkip(*f.Offset)
	}

	if len(f.Sorts) > 0 {
		findOptions.SetSort(f.Sorts)
	}

	return findOptions
}
