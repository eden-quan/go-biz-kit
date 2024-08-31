package inject

import (
	"github.com/eden-quan/go-biz-kit/injection"
	"github.com/eden-quan/go-biz-kit/message"
)

func Inject() {
	injection.Inject(message.NewQueueFactory)
}

// InjectIns 使用实例化的方式注入消息队列
func InjectIns(inj *injection.Injector) {
	inj.Inject(message.NewQueueFactory)
}
