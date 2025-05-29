package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cg917658910/go-study/lib/cache"
	"github.com/redis/go-redis/v9"
)

type Cg map[string]map[int]string

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	fmt.Println("setup mysql...")
	//db.SetupMysql(ctx)
	cache.SetupRedis(ctx)

	testCg := make(Cg)
	testCg["ss"] = map[int]string{
		1: "ss",
	}

	playgroundRedis()
	/* article, err := models.FindArticle(ctx, 10)
	if err != nil {
		fmt.Println("Find Article error: ", err)
		return
	}

	fmt.Printf("Find Article: id[%d] title[%s]\n", article.ID, article.Title) */
}

type NotifyResult struct {
	MsgId  string `json:"msg_id" redis:"msg_id"`
	Result string `json:"result" redis:"result"` //通知返回内容
	Code   int    `json:"status" redis:"status"` //通知状态 1 成功 2 重试多次无响应 3 无效请求地址
	Desc   string `json:"msg" redis:"msg"`       //通知结果描述
	//RequestTime string                  `json:"request_time"` //最后通知时间
}

type NotifyResultRepo struct {
	rdb       *redis.Client
	keyPrefix string
	ttl       time.Duration
	ctx       context.Context
}

func NewNotifyResultRepo(ctx context.Context) *NotifyResultRepo {
	return &NotifyResultRepo{
		rdb:       cache.RedisClient(),
		keyPrefix: "Cg:Notify:Order:Result",
		ctx:       ctx,
		ttl:       time.Hour * 1,
	}
}

func (repo *NotifyResultRepo) Set(msg *NotifyResult) error {
	key := fmt.Sprintf("%s:%s", repo.keyPrefix, msg.MsgId)
	data, _ := json.Marshal(msg)
	return repo.rdb.Set(repo.ctx, key, data, repo.ttl).Err()
}

func (repo *NotifyResultRepo) Get(id string) (*NotifyResult, error) {
	key := fmt.Sprintf("%s:%s", repo.keyPrefix, id)
	val, err := repo.rdb.Get(repo.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var u NotifyResult
	err = json.Unmarshal([]byte(val), &u)
	return &u, err
}

func (repo *NotifyResultRepo) Del(key string) (err error) {
	return repo.rdb.Del(repo.ctx, key).Err()
}

func setNotifyResault() {

	id := "test:1234"
	var data = &NotifyResult{
		MsgId:  id,
		Result: "success",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	repo := NewNotifyResultRepo(ctx)
	if err := repo.Set(data); err != nil {
		fmt.Printf("repo set %s err: %s\n", id, err)
		return
	}

	msg, err := repo.Get(id)
	if err != nil {
		fmt.Printf("repo get %s err: %s\n", id, err)
		return
	}
	if msg == nil {
		fmt.Println("not find ", id)
		return
	}
	fmt.Println(msg.MsgId)
}

func playgroundRedis() {
	setNotifyResault()

	//bloomFilter()
}

func bloomFilter() {
	client := cache.RedisClient()
	ctx := context.Background()
	client.BFAdd(ctx, "notify_result", "cg")
	bool, err := client.BFExists(ctx, "notify_result", "cg").Result()
	if err != nil {
		fmt.Printf("create bf error: %s\n", err)
	}
	fmt.Println("bool: ", bool)
}
