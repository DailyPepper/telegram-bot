package model

import "time"

type Item struct {
	Title      string
	Categories []string
	Link       string
	Date       time.Time
	Summary    string
	SourseName string
}

type Sourse struct {
	ID        int64
	Name      string
	FeedURL   string
	CreatedAd time.Time
}

type Article struct {
	ID          int64
	SourseId    int64
	Title       string
	Link        string
	Summary     string
	PublishedAt time.Time
	PostedAt    time.Time
	CreatedAt   time.Time
}
