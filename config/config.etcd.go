package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/contrib/config/etcd/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
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

const (
	ENV_DEV    = "dev"
	ENV_TEST   = "test"
	ENV_HOTFIX = "hotfix"
	ENV_PRE    = "pre"
	ENV_SIT    = "sit"
	ENV_UAT    = "uat"
	ENV_PROD   = "prod"
)

var envMap = map[string]string{
	ENV_DEV:    "/dev/",
	ENV_TEST:   "/test/",
	ENV_HOTFIX: "/hotfix/",
	ENV_PRE:    "/pre/",
	ENV_SIT:    "/sit/",
	ENV_UAT:    "/uat/",
	ENV_PROD:   "/prod/",
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
	var path etcd.Option
	if configure.APP.Env == ENV_PROD {
		path = etcd.WithPath("/")
	} else {
		path = etcd.WithPath("/" + configure.APP.Env)
	}
	etcdSource, err := etcd.New(client, path, etcd.WithPrefix(true))
	if err != nil {
		panic(fmt.Sprint("new etcd source failed with error", err))
	}
	var source config.Option
	if configure.ConfigCenter.LocalFile != "" {
		source = config.WithSource(
			etcdSource,
			file.NewSource(configure.ConfigCenter.LocalFile),
		)
	} else {
		source = config.WithSource(etcdSource)
	}
	conf := config.New(source, config.WithDecoder(confDecoder))
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
	var key = src.Key
	for _, v := range envMap {
		if strings.HasPrefix(key, v) {
			key = strings.TrimPrefix(key, v[:len(v)-1])
			break
		}
	}
	target[key] = container
	return nil
}
