package clientutil

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/eden-quan/go-biz-kit/config/def"
)

func NewHttpClientConn(server *def.Server, logger log.Logger) (*http.Client, error) {
	opts := make([]http.ClientOption, 0)

	opts = append(opts, http.WithEndpoint(server.GetHttp().GetAddress()))

	// 默认超时 5 分钟
	timeOut := time.Minute * 5
	if server.GetHttp() != nil && server.GetHttp().GetTimeout() != nil {
		timeOut = server.GetHttp().GetTimeout().AsDuration()
	}
	opts = append(opts, http.WithTimeout(timeOut))

	opts = append(opts, http.WithMiddleware(DefaultClientMiddlewares(logger)...))

	client, err := http.NewClient(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	return client, err
}
