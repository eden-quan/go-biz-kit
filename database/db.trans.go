package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"

	bizkit "github.com/eden/go-biz-kit"
)

// transImpl 实现了 kit.Transaction 的接口，提供了数据库事务的创建/提交等操作
type transImpl struct {
	alreadyCommit bool
	commitErr     error
	trans         *sqlx.Tx
	ctx           context.Context // 以后还可以使用该 ctx 中包含的信息来记录分布式事务
}

func newTransImpl(ctx context.Context, trans *sqlx.Tx, err error) *transImpl {
	return &transImpl{
		trans:     trans,
		ctx:       ctx,
		commitErr: err,
	}
}

func baseType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = reflectx.Deref(t)
	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}
	return t, nil
}

func (tx *transImpl) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	//
	//rows, err := tx.trans.QueryxContext(ctx, query, args...)
	//if err != nil {
	//	return err
	//}
	//
	//columns, err := rows.Columns()
	//if err != nil {
	//	return err
	//}
	//
	//value := reflect.ValueOf(dest)
	//direct := reflect.Indirect(value)
	//
	//slice, err := baseType(value.Type(), reflect.Slice)
	//if err != nil {
	//	return err
	//}
	//
	////isPtr := slice.Elem().Kind() == reflect.Ptr
	//base := reflectx.Deref(slice.Elem())
	//
	//fields := tx.trans.Mapper.TraversalsByName(base, columns)
	//
	//values := make([]interface{}, len(columns))
	//
	//for rows.Next() {
	//	vp := reflect.New(base)
	//	v := reflect.Indirect(vp)
	//}

	return nil
}

func (tx *transImpl) Get() *sqlx.Tx {
	return tx.trans
}

func (tx *transImpl) Commit(err error) error {
	if tx.alreadyCommit || tx.commitErr != nil {
		return tx.commitErr
	}

	// when error occur, we process the transaction immediately
	if err != nil {
		commitErr := tx.trans.Rollback()
		tx.commitErr = errors.Join(err, commitErr)
		tx.alreadyCommit = true
		return tx.commitErr
	}

	// if currently is fine, check the counter
	if v, ok := getTransValue(tx.ctx); ok {
		v.level -= 1

		if v.level > 0 {
			return nil
		}
	}

	tx.commitErr = tx.Get().Commit()
	tx.alreadyCommit = true
	return errors.Join(err, tx.commitErr)
}

func WithTx(ctx context.Context, f bizkit.WithTxGlobalFunc) error {
	ctx = Begin(ctx)
	v, _ := getTransValue(ctx)
	v.level += 1

	err := f(ctx)
	return Commit(ctx, err)
}
