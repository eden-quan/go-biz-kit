package inject

import (
	"github.com/eden-quan/go-biz-kit/config"
	"github.com/eden-quan/go-biz-kit/injection"
)

/*
Inject 为项目注入 配置库
 1. NewConfigWithFiles 注入本地文件配置，本地文件配置包含应用配置及配置中心的地址，
    其他配置库都依赖于该基础组件，可通过 *config.LocalConfigure 获取
 2. NewConfigWatcher 注入配置监听器, 他依赖 LocalConfigure 提供的配置中心地址连接到配置中心，
    并提供了监听配置的能力，使用者可通过 *ConfigureWatcherRepo 获取
 3. NewConfiguration 注入基础配置, 他依赖于 ConfigureWatcherRepo 为他提供监听能力，
    并加载基础中间件配置，其中包括了数据库/日志等基础配置信息，使用者 *def.Configuration 获取基础配置信息
*/
func Inject() {
	InjectIns(injection.GlobalInjector())
}

// InjectIns 使用实例化的方式注册配置依赖项
func InjectIns(inj *injection.Injector) {
	inj.InjectMany(
		config.NewConfigWithFiles,
		config.NewConfigWatcher,
		//def.NewConfiguration,
	)
}
