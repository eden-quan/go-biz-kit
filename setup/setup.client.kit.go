package setup

import (
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"

	"github.com/eden/go-biz-kit/client"
	"github.com/eden/go-biz-kit/config/def"
)

func NewGRPCClientFactory(logger log.Logger) clientutil.RegisterGRPCClientFactoryType {
	return func(conf *def.Server) (grpc.ClientConnInterface, error) {
		if conf.Grpc == nil {
			return nil, errors.New("wanna create grpc client but grpc config is empty")
		}

		return clientutil.NewGrpcClientConn(conf, logger)
	}
}

func NewHTTPClientFactory(logger log.Logger) clientutil.RegisterHTTPClientFactoryType {
	return func(conf *def.Server) (*http.Client, error) {
		if conf.Http == nil {
			return nil, errors.New("wanna create http client but grpc config is empty")
		}

		return clientutil.NewHttpClientConn(conf, logger)
	}
}
