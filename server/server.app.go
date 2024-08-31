package servers

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.uber.org/fx"

	"github.com/eden/go-biz-kit/config"
)

// NewApp 通过配置信息提供 Kratos 的 APP 示例，以及对应的 Server (http/grpc), 供后续的实现
func NewApp(localConf *config.LocalConfigure, gs *grpc.Server, hs *http.Server, logger log.Logger) (*kratos.App, error) {

	servers := make([]transport.Server, 0)

	if gs != nil {
		servers = append(servers, gs)
	}
	if hs != nil {
		servers = append(servers, hs)
	}

	appOptions := []kratos.Option{
		kratos.ID(localConf.APP.Name),
		kratos.Name(localConf.APP.Name),
		kratos.Logger(logger),
		kratos.Server(servers...),
	}

	app := kratos.New(appOptions...)
	return app, nil
}

func StartKratosApp(lifecycle fx.Lifecycle, app *kratos.App, logger log.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := app.Run(); err != nil {
					err := fmt.Errorf("app.Run %+v\n", err)
					panic(err)
				}

			}()

			return nil
		},
		OnStop: func(_ context.Context) error {
			return app.Stop()
		},
	})
}
