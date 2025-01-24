package storege

import (
	"context"
	"myapp/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

type dbArticle struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	FeedURL   string `db:"feed_url"`
	CreatedAt string `db:"created_at"`
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		query: `INSERT INTO articles (source_id, title, link, summary, published_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING`,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt,
	); err != nil{
		return err
	}
	return nil
}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Arcticle, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil
	defer conn.Close()

	var articles []dbArticle
	if err := conn.SelectContext(
		ctx,
		&articles,
		`SELECT * FROM articles WHERE posted_at IS NULL AND published_at >= $1::timestamp ORDER BY published_at DESC LIMIT $2`,
		since.UTC().Format(time.RFC3339),
		limit
	); err != nil {
		return nil, err
	}
	return lo.Map(articles, func(article dbArticle, _ int)model.Arcticle{
		return model.Arcticle
	}),nil
}
func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {}
