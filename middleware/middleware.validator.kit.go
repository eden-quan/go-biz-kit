package middlewareutil

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"strconv"

	"github.com/go-kratos/kratos/v2/middleware"
)

// validator ...
type validator interface {
	Validate() error
	ValidateAll() error
}

type validateError interface {
	ErrorName() string
	Code() int64
	HttpCode() int64
	Reason() string
	Cause() error
}

// Validator 是参数验证中间件，
func Validator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			err = checkValidator(req)
			if err != nil {
				return nil, err
			}

			return handler(ctx, req)
		}
	}
}

func checkValidator(req interface{}) (err error) {
	v, ok := req.(validator)
	if !ok {
		return
	}

	err = v.Validate()
	if err == nil {
		return
	}

	detail, ok := err.(validateError)
	if !ok {
		return
	}

	err = errors.New(
		int(detail.HttpCode()), detail.ErrorName(), detail.Reason(),
	).WithMetadata(map[string]string{
		"BizCode":        strconv.Itoa(int(detail.Code())),
		"DefaultMessage": detail.ErrorName(),
	})

	return
}
