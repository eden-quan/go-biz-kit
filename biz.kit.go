package biz_kit

import (
	"context"
	"github.com/redis/go-redis/v9"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/eden-quan/go-biz-kit/message"
)

type MongoDB interface {
	Get() *mongo.Database
}

type WithTxFunc = func(ctx context.Context, tx Transaction) error
type WithTxGlobalFunc func(ctx context.Context) error

// MySQL 为简易的 DB 封装层，为后续的优化提供切入点
type MySQL interface {
	Get() *sqlx.DB
	GetTx(ctx context.Context) (Transaction, context.Context, error)
	WithTx(ctx context.Context, f WithTxFunc) error
}

// Database 定义了满足 SQL 规范的的接口，后续用于替换各种数据库实现
type Database interface {
	Get() *sqlx.DB
	GetTx(ctx context.Context) (Transaction, context.Context, error)
	WithTx(ctx context.Context, f WithTxFunc) error
}

// Transaction 定义了满足 SQL 规范的事务接口
type Transaction interface {
	Get() *sqlx.Tx
	// Commit 根据参数决定提交事务或回滚, 存在错误并进行回滚时，如果回滚也发生错误，则会将原始错误和回滚的错误进行组合, 否则会返回业务的错误
	Commit(error) error
}

// Redis 为对 Redis 客户端的简易包装，为后续的优化提供切入点
type Redis interface {
	Get() redis.UniversalClient
}

// MessageQueue 为对消息队列的封装, 抽象具体的操作后允许用户不关注具体使用的消息队列中间件
type MessageQueue interface {
	Get() message.QueueFactory
}
