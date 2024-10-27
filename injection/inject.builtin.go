package injection

import (
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/fx"
)

type injectDep struct{}

func injectEntrance(_ *injectDep) {
	//log.NewHelper(logger).Info("Injected service finished...")
}

func collectServiceProvider(logger log.Logger, httpRegister []fmt.Stringer, grpcRegister []fmt.Stringer) *injectDep {
	helper := log.NewHelper(logger)
	helper.Info("\tinject http: ", httpRegister)
	helper.Info("\tinject grpc: ", grpcRegister)
	return &injectDep{}
}

type graphOption struct {
	Output string
}

func (g *graphOption) Provide() fx.Option {
	return fx.Invoke(func(dot fx.DotGraph) {
		f, err := os.OpenFile(g.Output, os.O_WRONLY|os.O_CREATE, 0644)

		if err == nil {
			_, err = f.WriteString(string(dot))
		}

		defer func() {
			_ = f.Close()
		}()

		if err != nil {
			fx.Error(err)
		}
	})
}

// WithGraph 在启动时将依赖信息输出到 output 指定的文件中
func WithGraph(output string) Option {
	return &graphOption{Output: output}
}

// invokeOption 用于在完成依赖注入后触发指定的函数调用，触发的调用需要依赖
// 已注入的组件
type invokeOption struct {
	invoker interface{}
}

func (i *invokeOption) Provide() fx.Option {
	return fx.Invoke(i.invoker)
}

// WithInvoke 用于在完成依赖注入后触发指定的函数调用，触发的调用需要依赖
// 已注入的组件, invoker 可以是任意参数的函数
func WithInvoke(invoker interface{}) Option {
	return &invokeOption{invoker: invoker}
}

type replaceOption struct {
	replacer interface{}
}

func (r *replaceOption) Provide() fx.Option {
	return fx.Decorate(r.replacer)
}

// WithReplace 允许用户替换已注入的依赖，如已注入 Logger，可通过 Replace 来为 Logger 加入新的前缀
// replacer 参数为任意函数，他可以接受已注入的依赖，并返回新的依赖，新的依赖会根据类型替换已有的依赖
func WithReplace(replacer interface{}) Option {
	return &replaceOption{replacer: replacer}
}
