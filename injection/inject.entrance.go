package injection

type Entrance struct {
	injs []*Injector
}

func NewEntrance(injs ...*Injector) *Entrance {
	if injs == nil {
		injs = make([]*Injector, 0)
	}
	return &Entrance{
		injs: injs,
	}
}

// Replace 将 replacer 应用到关联的所有 Injector 中,
func (e *Entrance) Replace(replacer interface{}) {
	panic("implement me")
}

// DoIt 启动已注入的所有 Injector 并对
func (e *Entrance) DoIt(opts ...Option) {
	panic("implement me")
}

// MixUp 将多个 injs 合并成一个统一的依赖注入入口 Entrance，通过该入口启动的 Injector 将由入口进行管理
// 通过该入口注入的依赖将会被注入到所有的 Injector 中
// 该入口主要是替换那些需要以单例方式提供的组件， 因为当前每个 Injector 中都会持有所需依赖的一个实例，
// 通过该入口注入的实例将可以在多个 Injector 中共享，如可用共享的 Kratos APP 实例来管理所有的服务
func MixUp(injs []*Injector) *Entrance {
	ent := NewEntrance(injs...)

	// TODO: MixUp 操作列表
	//       1. GRPC 的 Server 在注册时会依赖 grpc.ServiceRegistrar，因此要将该依赖换成 MixUp 后的单例
	//       2. HTTP 的 Server 在注册时同 GRPC 处理，但他依赖的事 http.Server
	//       3. 最终的 Server 都会注册到统一的 APP 中，并由 APP 进行管理，因此需要将 APP 也替换为一个单例的 Provider
	//       4. 如何保证在所有的依赖都注册完成后才启动 APP -> 为 APP 提供一个需要大部分基础组件的依赖，确保其他组件完成后才触发 APP 的注入
	//       5. Server 会在有自定义中间件时通过 Option 注入中间件，这些中间件不应该跨越 Injector, 因此需要在构建 Server 中间件时，为中间件增加一层封装
	// 			这层封装可以通过 Injector 所在的应用名称来进行过滤

	return ent
}
