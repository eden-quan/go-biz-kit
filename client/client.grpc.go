package clientutil

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/config/def"
	"google.golang.org/grpc"
)

// GrpcClientConn 是对 grpc.ClientConn 的简易封装，用于后续实现热更新等能力
// TODO: 当前使用网关来实现该功能，后续根据需要再支持多节点负载均衡等能力
type GrpcClientConn struct {
	*grpc.ClientConn                      // 基于 grpc ClientConn 进行简易封装
	server           *def.Server          // server 为当前客户端的服务发现信息
	options          []kgrpc.ClientOption // options 为当前客户端的配置项
	logger           log.Logger           // 客户端的日志记录器
	helper           *log.Helper
	isOk             bool // 该连接是否可正常使用
}

func NewGrpcClientConn(server *def.Server, logger log.Logger, options ...kgrpc.ClientOption) (*GrpcClientConn, error) {
	var opts []kgrpc.ClientOption

	opts = append(opts, kgrpc.WithEndpoint(server.GetGrpc().GetAddress()))

	// 默认超时 5 分钟
	timeOut := time.Minute * 5
	if server.GetGrpc() != nil && server.GetGrpc().GetTimeout() != nil {
		timeOut = server.GetGrpc().GetTimeout().AsDuration()
	}
	opts = append(opts, kgrpc.WithTimeout(timeOut))
	opts = append(opts, kgrpc.WithMiddleware(DefaultClientMiddlewares(logger)...))
	opts = append(opts, options...)

	client := &GrpcClientConn{
		ClientConn: nil,
		server:     server,
		options:    opts,
		logger:     logger,
		isOk:       false,
		helper:     log.NewHelper(logger, log.WithMessageKey(fmt.Sprintf("GRPCClient %s", server.Grpc.Name))),
	}

	err := client.connect()
	return client, err
}

func (c *GrpcClientConn) connect() error {
	conn, err := kgrpc.DialInsecure(context.Background(), c.options...)
	c.ClientConn = conn
	c.isOk = err == nil
	return err
}

func (c *GrpcClientConn) Invoke(ctx context.Context, method string, args, reply any, opts ...grpc.CallOption) error {

	tryTimes := 1
	var lastError error

	for tryTimes <= 3 {

		if c.ClientConn == nil {
			_ = c.connect()
		}

		if c.ClientConn.GetState() == connectivity.Shutdown || c.ClientConn.GetState() == connectivity.TransientFailure {
			_ = c.connect()
		}
		err := c.ClientConn.Invoke(ctx, method, args, reply, opts...)
		if err == nil {
			return nil
		}

		lastError = err
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.DeadlineExceeded:
				c.helper.Errorf("grpc client calling %s with deadline exceeded", method)
				return err
			case codes.Canceled:
				c.helper.Errorf("grpc client calling %s with cancelled", method)
				return err
			case codes.Internal:
				c.helper.Errorf("grpc client calling %s with internal error", method)
				return err
			case codes.InvalidArgument:
				c.helper.Errorf("grpc client calling %s with invalid argument", method)
				return err
			case codes.Unimplemented:
				c.helper.Errorf("grpc client calling %s is unimplemented", method)
				return err
			case codes.Unavailable:
				c.helper.Errorf("grpc client calling %s with unavailable - retry %d times", method, tryTimes)
				break
			default:
				c.helper.Errorf("grpc client calling %s with unhandled status %d and message %s", method, st.Code(), st.Message())
				return err
			}

			tryTimes += 1
		} else {
			return err
		}
	}

	return lastError
}

func (c *GrpcClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return c.ClientConn.NewStream(ctx, desc, method, opts...)
}
