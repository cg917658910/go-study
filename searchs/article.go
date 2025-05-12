package searchs

import (
	"fmt"

	"gorm.io/gorm"
)

type ArticleSearch struct {
	IDs   []uint
	Title string
}

func (s *ArticleSearch) BuildSearch() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", s.Title))
		return db
	}
}
