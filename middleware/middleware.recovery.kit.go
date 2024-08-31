package middlewareutil

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware/recovery"

	errorv1 "github.com/eden-quan/go-biz-kit/common/def"
)

var _ = recovery.ErrUnknownRequest

// RecoveryHandler ...
func RecoveryHandler() recovery.HandlerFunc {
	return func(ctx context.Context, req, err interface{}) error {
		message := fmt.Sprintf("%v", err)
		e := errorv1.ErrorInternalServer(message)
		return e
	}
}
