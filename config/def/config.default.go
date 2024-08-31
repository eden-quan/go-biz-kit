package def

import (
	"github.com/eden/go-biz-kit/config"
)

type Configuration struct {
	Server       *Server  `conf_path:"/basic/config"`               // 服务的地址配置，包括监听的地址，端口等
	Registry     Registry `conf_path:"/basic/online"`               // 其他在线服务
	Profile      Profile  `conf_path:"/basic/profile/config"`       // 性能分析配置
	Log          Log      `conf_path:"/middleware/log/config"`      // 服务的日志配置
	Redis        Redis    `conf_path:"/middleware/redis/config"`    // Redis 配置
	Mongo        Mongo    `conf_path:"/middleware/mongodb/config"`  // MongoDB 配置
	Database     Database `conf_path:"/middleware/database/config"` // MySQL 配置
	MessageQueue RabbitMQ `conf_path:"/middleware/rabbitmq/config"` // RabbitMQ 配置
	Tracing      Tracing  `conf_path:"/middleware/tracing/config"`  // 链路跟踪配置
}

// NewConfiguration 创建一个新的配置实例，该实例支持热更新等能力, 为了保证全局统一，该实例为单例模式
// 注意：一般情况下业务系统会在基础配置上再添加自定义的配置
func NewConfiguration(config config.ConfigureWatcherRepo) (*Configuration, error) {
	conf := &Configuration{}
	err := config.LoadAndStart(&conf)

	return conf, err
}
