-- +goose Up
-- +goose StatementBegin
CREATE TABLE articles (
    id SERIAL PRIMARY key,
    source_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    link VARCHAR(255) NOT NULL,
    summery TEXT NOT NULL
    published_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    posted_at TIMESTAMP,
    CONSTRAINT fk_articles_source_id
)-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
