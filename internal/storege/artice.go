package storege

import (
	"context"
	"myapp/model"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
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
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(
		ctx,
		`INSERT INTO articles (source_id, title, link, summary, published_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING`,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var articles []dbArticle
	err = conn.SelectContext(
		ctx,
		&articles,
		`SELECT * FROM articles WHERE posted_at IS NULL AND published_at >= $1 ORDER BY published_at DESC LIMIT $2`,
		since.UTC().Format(time.RFC3339),
		limit,
	)
	if err != nil {
		return nil, err
	}

	return lo.Map(articles, func(article dbArticle, _ int) model.Article {
		return model.Article{
			ID:        article.ID,
			Name:      article.Name,
			FeedURL:   article.FeedURL,
			CreatedAt: article.CreatedAt,
		}
	}), nil
}

func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(
		ctx,
		`UPDATE articles SET posted_at = $1 WHERE id = $2`,
		time.Now().UTC(),
		id,
	)
	if err != nil {
		return err
	}

	return nil
}
