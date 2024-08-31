package inject

import (
	"github.com/eden/go-biz-kit/injection"
	"github.com/eden/go-biz-kit/setup"
)

/*
Inject 提供了默认的 setup 依赖注入，主要包括了各个中间件的注入，需要使用以下的中间件，需要先注入 config 模块，通过 config 模块
为中间件提供配置中心的配置信息，通过该函数可以得到以下注入信息

 1. Logger 注入，提供了基础的日志功能，可通过 log.Logger 得到
 2. LoggerManager 注入，提供了日志库管理能力，可通过 *LogManager 得到
 3. Redis 注入，提供缓存的访问及管理能力，可通过 kit.Redis 得到
 4. MongoDB 注入，提供了 MongoDB 数据库的访问能力，可通过 kit.MongoDB 得到
 5. MySQL 注入，提供了 MySQL 数据库的访问能力，可通过 kit.MySQL 得到
 6. Messaging 注入，提供了基于 RabbitMQ 的消息队列能力, 可通过 kt.MessageQueue 得到
 7. Tracing 注入，提供了全局的链路跟踪能力，所有通过依赖注入的客户端都能够自动得到链路跟踪的能力
*/
func Inject() {
	InjectIns(injection.GlobalInjector())
}

/*
InjectIns 使用创建实例的方式提供了默认的 setup 依赖注入，主要包括了各个中间件的注入，需要使用以下的中间件，需要先注入 config 模块，通过 config 模块
为中间件提供配置中心的配置信息，通过该函数可以得到以下注入信息

 1. Logger 注入，提供了基础的日志功能，可通过 log.Logger 得到
 2. LoggerManager 注入，提供了日志库管理能力，可通过 *LogManager 得到
 3. Redis 注入，提供缓存的访问及管理能力，可通过 kit.Redis 得到
 4. MongoDB 注入，提供了 MongoDB 数据库的访问能力，可通过 kit.MongoDB 得到
 5. MySQL 注入，提供了 MySQL 数据库的访问能力，可通过 kit.MySQL 得到
 6. Messaging 注入，提供了基于 RabbitMQ 的消息队列能力, 可通过 kt.MessageQueue 得到
 7. Tracing 注入，提供了全局的链路跟踪能力，所有通过依赖注入的客户端都能够自动得到链路跟踪的能力
*/
func InjectIns(inj *injection.Injector) {
	// inj.InjectHTTPMiddleware
	inj.InjectMany(
		setup.NewLogger,
		setup.NewLoggerManager,
		setup.NewRedis,
		setup.NewMongoDB,
		setup.NewMySQLDatabase,
		setup.NewSQLDatabase,
		setup.NewMessageQueue,
		setup.NewTracing,
		setup.NewHTTPClientFactory,
		setup.NewGRPCClientFactory,
	)

	inj.InjectWithParam(setup.NewSQLActionManager, []string{"", `group:"sql_action_register"`}, nil)

	inj.Invoke(
		injection.WithInvoke(
			setup.NewTracing,
		),
	)

	inj.Invoke(
		injection.WithInvoke(
			setup.NewProfile,
		),
	)
}
