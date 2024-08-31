package clientutil

import (
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"

	"github.com/eden/go-biz-kit/config"
	"github.com/eden/go-biz-kit/config/def"
)

func RegisterHTTPClient(nameAndType []string, repo config.ConfigureWatcherRepo, logger log.Logger) (*http.Client, error) {
	if len(nameAndType) != 2 {
		return nil, errors.New("generated code for provide service name with wrong format, we need array like ['service_name', 'http']")
	}

	serviceName := nameAndType[0]
	serviceType := nameAndType[1]

	server := &def.Server{}
	if serviceName[0] == '/' {
		serviceName = serviceName[1:]
	}
	err := repo.LoadWithPath(server, fmt.Sprintf("/registry/%s/config", serviceName))
	err = repo.Start()

	if err != nil {
		return nil, err
	}

	if serviceType == "http" {
		if server.GetHttp() == nil {
			return nil, fmt.Errorf("can not find http registry info for %s, please check the config center (etcd)", serviceName)
		}
		return NewHttpClientConn(server, logger)
	}

	return nil, fmt.Errorf("registry got wrong type %s for service %s", serviceType, serviceName)
}

func RegisterGRPCClient(nameAndType []string, repo config.ConfigureWatcherRepo, logger log.Logger) (grpc.ClientConnInterface, error) {
	if len(nameAndType) != 2 {
		return nil, errors.New("generated code for provide service name with wrong format, we need array like ['service_name', 'grpc']")
	}

	serviceName := nameAndType[0]
	serviceType := nameAndType[1]

	server := &def.Server{}
	if serviceName[0] == '/' {
		serviceName = serviceName[1:]
	}
	err := repo.LoadWithPath(server, fmt.Sprintf("/registry/%s/config", serviceName))
	err = repo.Start()

	//err := repo.LoadAndStart(server)
	if err != nil {
		return nil, err
	}

	if serviceType == "grpc" {
		if server.GetGrpc() == nil {
			return nil, fmt.Errorf("can not find grpc registry info for %s, please check the config center (etcd)", serviceName)
		}
		return NewGrpcClientConn(server, logger)
	}

	return nil, fmt.Errorf("registry got wrong type %s for service %s", serviceType, serviceName)
}

type RegisterGRPCClientFactoryType = func(conf *def.Server) (grpc.ClientConnInterface, error)
type RegisterHTTPClientFactoryType = func(conf *def.Server) (*http.Client, error)
