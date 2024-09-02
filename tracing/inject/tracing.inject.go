package inject

import (
	"gitlab.lainuoniao.cn/eden-quan/go-biz-kit/injection"
	"gitlab.lainuoniao.cn/eden-quan/go-biz-kit/tracing"
)

func Inject() {
	injection.Invoke(
		injection.WithInvoke(
			tracing.InitTracing,
		))
}

// InjectIns 使用创建实例的方式注入链路跟踪
func InjectIns(inj *injection.Injector) {
	inj.Invoke(
		injection.WithInvoke(
			tracing.InitTracing,
		))
}
