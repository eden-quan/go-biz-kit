package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	kit "github.com/eden-quan/go-biz-kit"
	errorutil "github.com/eden-quan/go-biz-kit/error"
)

const DBContextKey = "eden.db.ctx.key"

// dbContextValue 在开启全局事务时，会保存在 Context 中，并使用 level 记录当前使用者的深度，
// 在每个使用者启动事务时 +1, 在每个使用者提交事务时 -1，最终提交事务时该值应被变更为 0，否则说明使用过程中存在事务使用错误
type dbContextValue struct {
	level int32
	tx    *transImpl
}

// DBImpl 为 kit.Database 的实现，该类型提供了数据库访问的入口
type DBImpl struct {
	db *sqlx.DB
}

// NewDB 创建一个新的 DBImpl 实例
func NewDB(db *sqlx.DB) *DBImpl {
	return &DBImpl{db: db}
}

// Begin 构建一个新的 ctx，后续使用该 Ctx 的数据库操作都将使用同一个事务
func Begin(ctx context.Context) context.Context {
	_, ok := getTransValue(ctx)
	if ok {
		return ctx
	}

	return context.WithValue(ctx, DBContextKey, &dbContextValue{
		level: 0,
		tx:    nil,
	})
}

func getTransValue(ctx context.Context) (*dbContextValue, bool) {
	v, exists := ctx.Value(DBContextKey).(*dbContextValue)
	return v, exists
}

// Commit 提交上下文 ctx 中关联的全局事务，并对事务使用的正确性进行检查，如果上下文中不包含全局事务，
// 则不进行任何操作
func Commit(ctx context.Context, err error) error {
	if v, ok := getTransValue(ctx); ok && v.tx != nil {
		return v.tx.Commit(err)
	}

	return err
}

func IsNoRowErr(err error) bool {
	if err == nil {
		return false
	}

	for err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true
		}

		err = errors.Unwrap(err)
	}

	return false
}

func (db *DBImpl) IsNoRowErr(err error) bool {
	return IsNoRowErr(err)
}

// Get 获取当前数据库的实例
func (db *DBImpl) Get() *sqlx.DB {
	return db.db
}

func (db *DBImpl) WithTx(ctx context.Context, f kit.WithTxFunc) (err error) {
	ctx = Begin(ctx)
	tx, ctx, txErr := db.GetTx(ctx)

	if txErr != nil {
		err = errorutil.DBTxError.FromError(txErr)
		return
	}

	err = f(ctx, tx)
	err = tx.Commit(err)

	return
}

// GetTx 返回事务及相关的上下文，后续事务的提交需要使用本次返回的上下文，否则将会导致提交失败
// 上下文主要的作用是在上下文中维护事务信息，为使用同一上下文的代码共享同一个事务, 确保一致性
func (db *DBImpl) GetTx(ctx context.Context) (kit.Transaction, context.Context, error) {

	var value *dbContextValue
	var ok bool
	var trans *sqlx.Tx

	// 如果 ok 为 false, 说明未开启全局事务，仍使用独立的事务处理功能
	value, ok = getTransValue(ctx)
	if ok { // 开启了全局事务
		if value.tx == nil {
			trans, err := db.newSQLTx(ctx)
			value.tx = newTransImpl(ctx, trans, err)
		}

		value.level += 1
		return value.tx, ctx, value.tx.commitErr
	}

	// 否则创建新的事务，并关联到对应的上下文
	trans, err := db.newSQLTx(ctx)
	tx := newTransImpl(ctx, trans, err)
	return tx, ctx, tx.commitErr
}

// newSQLTx 创建一个数据库事务，默认级别为可重复度
func (db *DBImpl) newSQLTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := db.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	err = errorutil.DBGetError.FromError(err)
	return tx, err
}

// In 将 query 中的 IN 条件根据 args 的个数，扩展到为多个 ? 的形式
func (db *DBImpl) In(query string, args ...interface{}) (string, []interface{}, error) {
	//db.Get().NamedExecContext()
	return In(query, args...)
}

/*
NamedIn 将 query 中
的变量绑定条件 :var_name 的形式以 ? 的形式扩展，并将 slice 类型以 ? 的形式展开到 IN 条件中

```

	query := `SELECT id FROM table WHERE name IN (:name) AND is_deleted = :is_deleted`
	arg := { "name": ["a", "b", "c"], "is_deleted": 0 }
	query, args, err := database.NamedIn(query, arg)
	// query => SELECT id FROM table WHERE name IN (?, ?, ?) AND is_deleted = ?
	// args = ["a", "b", "c", 0]

```
*/
func (db *DBImpl) NamedIn(query string, arg interface{}) (string, []interface{}, error) {
	return NamedIn(query, arg)
}

// NamedIn 将 query 中的变量绑定条件 :var_name 的形式以 ? 的形式扩展，并将 slice 类型以 ? 的形式展开到 IN 条件中
func NamedIn(query string, arg interface{}) (string, []interface{}, error) {
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return query, nil, err
	}

	query, args, err = sqlx.In(query, args...)
	return query, args, err
}

// In 将 query 中的 IN 条件根据 args 的个数，扩展到为多个 ? 的形式
func In(query string, args ...interface{}) (string, []interface{}, error) {
	return sqlx.In(query, args...)
}
