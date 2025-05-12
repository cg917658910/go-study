package models

import (
	"context"

	"github.com/cg917658910/go-study/lib/db"
	"github.com/cg917658910/go-study/searchs"
)

type Article struct {
	ID    uint `gorm:"primaryKey"`
	Title string
}

func FindArticle(ctx context.Context, id uint) (*Article, error) {
	data := &Article{}

	err := db.DB().First(data, id).Error
	if err != nil {
		return nil, err
	}

	return data, preloadArticle(ctx, data)
}

func ListArticle(ctx context.Context, search *searchs.ArticleSearch) ([]*Article, error) {

	data := make([]*Article, 0)

	return data, db.DB().Scopes(search.BuildSearch()).Find(data).Error
}

func preloadArticle(ctx context.Context, articles ...*Article) error {

	return nil
}
