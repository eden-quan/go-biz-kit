package config

type ConfigureWatcherRepo interface {
	// Load 解析配置对象 object 中的 `conf_path` 标签及 json 标签，获取配置对象配置中心路径的映射
	Load(object interface{}) error
	// AddPrefix 为当前的 Manager 对象添加一个前缀并指定优先级，指定前缀中后与配置路径之间会进行合并, 优先级高的配置会覆盖优先级低的配置,
	// 当 ignoreEmpty 为 true 时，不检查缺失的配置项
	AddPrefix(prefix string, priority int, ignoreEmpty bool)
	// Start 启动 ConfigureRepo 对配置中心的监听，配置中心中所有的配置变更会实时更新到之前通过 Load 绑定的对象上
	Start() error

	// LoadAndStart 提供了组合调用 Load 跟 Start 的能力，该接口为用户减少需要处理的错误状态
	LoadAndStart(object interface{}) error

	// LoadWithPath 为 object 建立 path 的监听，并在 path 发生变化时为其提供热更新能力, 该接口一般用于运行时需要动态监听配置的情形
	LoadWithPath(object interface{}, path string) error
}
