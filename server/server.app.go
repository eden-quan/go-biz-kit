package servers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"

	"gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/config"
)

// NewApp 通过配置信息提供 Kratos 的 APP 示例，以及对应的 Server (http/grpc), 供后续的实现
func NewApp(localConf *config.LocalConfigure, gs *grpc.Server, hs *http.Server, logger log.Logger) (*kratos.App, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   localConf.ConfigCenter.Endpoints,
		Username:    localConf.ConfigCenter.Username,
		Password:    localConf.ConfigCenter.Password,
		DialTimeout: time.Second,
	})
	if err != nil {
		panic(err)
	}
	reg := etcd.New(
		cli,
		etcd.Namespace(fmt.Sprintf("/microservices/%s", localConf.APP.Env)),
		etcd.RegisterTTL(time.Second*15),
	)
	servers := make([]transport.Server, 0)

	if gs != nil {
		servers = append(servers, gs)
	}
	if hs != nil {
		servers = append(servers, hs)
	}
	id, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	appOptions := []kratos.Option{
		kratos.ID(id),
		kratos.Name(localConf.APP.Name),
		kratos.Logger(logger),
		kratos.Server(servers...),
		kratos.Registrar(reg),
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
