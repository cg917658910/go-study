package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cg917658910/go-study/lib/db"
	"github.com/cg917658910/go-study/models"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	fmt.Println("setup mysql...")
	db.SetupMysql(ctx)

	article, err := models.FindArticle(ctx, 10)
	if err != nil {
		fmt.Println("Find Article error: ", err)
		return
	}

	fmt.Printf("Find Article: id[%d] title[%s]\n", article.ID, article.Title)
}
