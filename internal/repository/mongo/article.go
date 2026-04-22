package mongorepo

import (
	"app/domain"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArticleRepository struct {
	conn              *mongo.Database
	articleCollection string
}

func NewArticleRepository(Conn *mongo.Database) *ArticleRepository {
	return &ArticleRepository{
		conn:              Conn,
		articleCollection: "articles",
	}
}

// for default query
func defaultArticleQuery() map[string]any {
	return map[string]any{
		"deleted_at": map[string]any{
			"$eq": nil,
		},
	}
}

func (r *ArticleRepository) FetchArticle(ctx context.Context, options domain.ArticleFilter) (cur *mongo.Cursor, err error) {
	// generate query
	query := options.Query(defaultArticleQuery())

	// query options
	findOptions := options.FindOptions()

	cur, err = r.conn.Collection(r.articleCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchArticle Find:", err)
		return
	}

	return
}

func (r *ArticleRepository) FetchOneArticle(ctx context.Context, options domain.ArticleFilter) (row *domain.Article, err error) {
	// generate query
	query := options.Query(defaultArticleQuery())

	err = r.conn.Collection(r.articleCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneArticle FindOne:", err)
		return
	}

	return
}

func (r *ArticleRepository) CountArticle(ctx context.Context, options domain.ArticleFilter) (total int64) {
	// generate query
	query := options.Query(defaultArticleQuery())

	total, err := r.conn.Collection(r.articleCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountArticle", err)
		return 0
	}

	return
}

func (r *ArticleRepository) CreateArticle(ctx context.Context, row *domain.Article) (err error) {
	_, err = r.conn.Collection(r.articleCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateArticle InsertOne:", err)
		return
	}
	return
}
