package config

import (
	"flag"
	"fmt"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	pkgerrors "github.com/pkg/errors"
)

type APPConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Env     string `json:"env"`
}

type LocalConfigure struct {
	APP APPConfig `json:"app"` // APP 基本配置
	// ConfigureCenter 配置中心的地址
	ConfigCenter struct {
		// 本地配置文件，如果配置了则不使用配置中心，而是使用本地的配置文件
		LocalFile string   `json:"local_file"`
		Endpoints []string `json:"endpoints"`
		Username  string   `json:"username"`
		Password  string   `json:"password"`
		Timeout   string   `json:"timeout"` // example 10s
	} `json:"config_center"`
}

var (
	configFilepath string // 配置文件 所在的目录
)

func init() {
	flag.StringVar(&configFilepath, "conf", "./configs", "加载本地配置文件获取配置中心地址, example: xxx --conf ./configs")
}

// NewConfigWithFiles 初始化配置手柄 ,
// WARN: kratos 的实现会产生错误日志，暂时忽略
func NewConfigWithFiles() (conf *LocalConfigure, err error) {
	// parses the command-line flags
	if !flag.Parsed() {
		flag.Parse()
	}

	var opts []config.Option
	opts = append(opts, config.WithSource(file.NewSource(configFilepath)))
	handler := config.New(opts...)

	if err = handler.Load(); err != nil {
		err = pkgerrors.WithStack(err)
		return
	}

	conf = &LocalConfigure{}
	if err = handler.Scan(conf); err != nil {
		err = pkgerrors.WithStack(err)
		return
	}

	err = handler.Close()
	if err != nil {
		panic(fmt.Sprint("process local config file with error ", err))
	}

	return conf, err
}
