package config

import (
	"fmt"
	"time"

	"github.com/go-kratos/kratos/contrib/config/etcd/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/encoding"
	etcdclient "go.etcd.io/etcd/client/v3"
)

type KratosConfWrapper struct {
	conf config.Config
}

func (k *KratosConfWrapper) Load(object interface{}) error {
	panic("implement me")
}

func (k *KratosConfWrapper) AddPrefix(prefix string, priority int, ignoreEmpty bool) {
	panic("implement me")
}

func (k *KratosConfWrapper) Start() error {
	return nil
}

func (k *KratosConfWrapper) LoadAndStart(object interface{}) error {
	panic("implement me")
}

func (k *KratosConfWrapper) LoadWithPath(object interface{}, path string) error {
	val := k.conf.Value(path)
	return val.Scan(object)
}

func NewConfig(configure *LocalConfigure) (config.Config, ConfigureWatcherRepo, error) {
	timeOut, err := time.ParseDuration(configure.ConfigCenter.Timeout)
	if err != nil {
		return nil, nil, err
	}
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints:   configure.ConfigCenter.Endpoints,
		Username:    configure.ConfigCenter.Username,
		Password:    configure.ConfigCenter.Password,
		DialTimeout: timeOut,
	})
	if err != nil {
		panic(fmt.Sprint("connect to etcd failed with error", err))
	}
	etcdSource, err := etcd.New(client, etcd.WithPath("/"), etcd.WithPrefix(true))
	if err != nil {
		panic(fmt.Sprint("new etcd source failed with error", err))
	}
	conf := config.New(
		config.WithSource(etcdSource),
		config.WithDecoder(confDecoder),
	)
	if err := conf.Load(); err != nil {
		panic(fmt.Sprint("load config failed with error", err))
	}
	return conf, &KratosConfWrapper{conf: conf}, nil
}

func confDecoder(src *config.KeyValue, target map[string]interface{}) error {
	var codec encoding.Codec
	if src.Format == "" {
		codec = encoding.GetCodec("json")
	} else {
		codec = encoding.GetCodec(src.Format)
	}
	if codec == nil {
		return fmt.Errorf("unsupported key: %s format: %s", src.Key, src.Format)
	}
	if src.Format != "" {
		return codec.Unmarshal(src.Value, &target)
	}
	container := make(map[string]interface{})
	err := codec.Unmarshal(src.Value, &container)
	if err != nil {
		return err
	}
	target[src.Key] = container
	return nil
}
