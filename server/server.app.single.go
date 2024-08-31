package servers

import (
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/eden/go-biz-kit/config"
)

var _app *kratos.App = nil
var _once sync.Once

func NewSingleApp(localConf *config.LocalConfigure, gs *grpc.Server, hs *http.Server, logger log.Logger) (*kratos.App, error) {
	_once.Do(func() {
		var err error
		_app, err = NewApp(localConf, gs, hs, logger)
		if err != nil {
			panic(fmt.Sprintf("create application failed with error %s", err))
		}
	})

	return _app, nil
}
