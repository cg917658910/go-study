package test

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

// 存储已发布文章的哈希表（避免重复）
var publishedArticles = make(map[string]bool)

// 订阅的 RSS 资源
var rssFeeds = map[string]string{
	"Tech":   "https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml",
	"News":   "http://feeds.bbci.co.uk/news/rss.xml",
	"Hacker": "https://news.ycombinator.com/rss",
}

// 定期拉取 RSS 并发布
func FetchAndPublishRSS() {
	parser := gofeed.NewParser()

	for {
		for channel, url := range rssFeeds {
			feed, err := parser.ParseURL(url)
			if err != nil {
				fmt.Println("❌ RSS 拉取失败:", err)
				continue
			}

			for _, item := range feed.Items {
				if _, exists := publishedArticles[item.GUID]; !exists {
					publishedArticles[item.GUID] = true
					message := fmt.Sprintf("[%s] %s - %s", channel, item.Title, item.Link)
					//PublishKafka(channel, message)
					fmt.Println("✅ 发布文章:", message)
				}
			}
		}
		time.Sleep(10 * time.Minute) // 每 10 分钟更新
	}
}
