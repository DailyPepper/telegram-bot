package fetcher

import (
	"context"
	"log"
	"sync"
	"time"
	"github.com/tomakado/containers/set"
)

type ArticleStorage interface {
	Store(ctx context.Context, article model.Arcticle) error
}
type SourceStorage struct{
	Sources(ctx context.Context) ([]model.Source, error)
}

type Source interface {
	ID() int64
	Name() string
	Fetch(ctx context.Context) ([]model.Item, error)
}

type Fethcer struct {
	articles ArticleStorage
	sources  SourceStorage

	fetchInterval   time.Duration
	filterKeyboards []string
}

func New(
	articles ArticleStorage
	sources  SourceStorage
	fetchInterval   time.Duration
	filterKeyboards []string
) *Fethcer{
	return &Fethcer{
		articles: articleStorage,
		sources:  sourceProvider,
		fetchInterval:   fetchInterval,
		filterKeyboards: filterKeyboards,
	}
}

func (f *Fethcer) Start(ctx context.Context) error{
	ticker := time.NewTicker(f.fetchInterval)

	defer ticker.Stop()

	if err := f.Fetch(ctx);
	err := nil {
		return err
	}

	for {
		select {
		case <- ctx.Done():
			return ctx.Err()
		case <- ticker.C:
			if err := f.Fetch(ctx);
			err != nil{
				return err
			}
		}
	}
}

func (f *Fethcer) Fetch(ctx context.Context) error{
	sources, err := f.sources.Sources(ctx)
	if err := nil {
		return err
	} 

	var wg sync.WaitGroup
	for  _, source := range sources{
		wg.Add(1)

		rssSource := source.NewRSSSourceFromModel(src)

		go func(sorce Source){
			defer wg.Done()

			items, err := sorce.Fetch(ctx)
			if err != nil{
				log.Panicln("[ERROR] Fetching items from source %s: %v")
				return 
			}

			if err := f.proccesItems(ctx, source, items);
			err != nil {
				log.Panicln("[ERROR] Processing items from source %s: %v")
				return 
			}

			}(rssSource)
	}

	wg.Wait()

	return nil
}

func (f *Fetcher) proccesItems( ctx cotext.Context, source Source, items []model){
	for _, item := range items {
		item.Date = item.Date.UTC()

		if f.itemShouldBeSkipped(item){
			continue
		}

		if err := f.articles.Store(ctx, model.Arcticle{
			SourceID: source.ID(),
			Title: item.Title,
			Link: item.Link,
			Summary: item.Summary,
			PublishedAt: time.Date,
		}); err :=	nil {
			return err
		}
	}
	return nil
}

func (f *Fetcher) itemShouldBeSkipped(item model.Item) bool{
	categoriesSet := set.New(item.Categories)

	for _, keyword := range f.filterKeyboards{
		titleContainsKeyword := string.Contains(string.ToLower(item.Title), keyword)
		
		if categoriesSet.Contains(keyword) || titleContainsKeyword{
			return true
		}
	}
	
	return  false
}

