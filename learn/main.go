package main

import (
	"context"
	"time"
)

type User struct {
	Name string
	Age  int
}

func getUsers() []User {
	users := []User{
		{"A", 1},
		{"B", 2},
	}
	return users
}

type IGetUsers interface {
	GetUsers() []User
}

type options struct {
	connectTimeout time.Duration
	readTimeot     time.Duration
	writeTimeout   time.Duration
	logError       func(ctx context.Context, err error)
}

type Cluster struct {
	opts options
}

// 优雅设置配置参数

type Option func(c *options)

// 设置参数
func WithConnectTimeout(d time.Duration) Option {
	return func(c *options) {
		c.connectTimeout = d
	}
}
func WithReadTimeout(d time.Duration) Option {
	return func(c *options) {
		c.readTimeot = d
	}
}
func NewCluster(opts ...Option) *Cluster {
	cluterOpts := &options{}
	for _, opt := range opts {
		opt(cluterOpts)
	}

	return &Cluster{
		opts: *cluterOpts,
	}
}

// 组合继承
// 基类
type MsgModel struct {
	msgId   int
	msgType int
}

func (msg *MsgModel) SetId(id int) {
	msg.msgId = id
}
func (msg *MsgModel) SetType(msgType int) {
	msg.msgType = msgType
}

// 子类
type GroupMsgModel struct {
	MsgModel // 组合继承
	msgId    int
}

func (group *GroupMsgModel) GetId() int {
	return group.msgId
}
func main() {

	//单例模式
	/* group := &GroupMsgModel{}
	group.SetId(100)
	group.SetType(1)
	fmt.Println("group.msgId =", group.msgId, "\tgroup.MsgModel.msgId =", group.MsgModel.msgId)
	fmt.Println("group.msgType =", group.msgType, "\tgroup.MsgModel.msgType =", group.MsgModel.msgType) */
	/* opts := []Option{
		WithConnectTimeout(10 * time.Second),
		WithReadTimeout(5 * time.Second),
	}

	cluster := NewCluster(opts...)

	println(cluster.opts.connectTimeout.String())
	println(cluster.opts.readTimeot.String()) */
	/* s1 := []int{1, 2, 3}
	s2 := s1
	s2[0] = 100
	fmt.Println(s1[0])
	production := "GoLand"
	getUsers()
	var i interface{} = production
	fmt.Println(i)
	fmt.Println(production)
	fmt.Println(&production) */
}
