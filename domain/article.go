package domain

import (
	mongo_model "app/domain/model/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type Article struct {
	ID        primitive.ObjectID `json:"id"         bson:"_id"`
	Title     string             `json:"title"      bson:"title"`
	Content   string             `json:"content"    bson:"content"`
	Author    ArticleAuthor      `json:"author"     bson:"author"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt gorm.DeletedAt     `json:"deleted_at" bson:"deleted_at"`
}

type ArticleAuthor struct {
	Name  string `json:"name"  bson:"name"`
	Email string `json:"email" bson:"email"`
}

var ArticleAllowedSort = []string{"title", "content", "author.name", "created_at", "updated_at"}

type ArticleFilter struct {
	mongo_model.DefaultFilter
	AuthorName *string
}

// ArticleFilter for MongoDB
func (f *ArticleFilter) Query(defaultQuery map[string]any) map[string]any {
	// default query
	f.DefaultFilter.Query(defaultQuery)

	if f.AuthorName != nil {
		defaultQuery["author.name"] = f.AuthorName
	}

	return defaultQuery
}
