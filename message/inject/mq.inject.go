package inject

import (
	"gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/injection"
	"gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/message"
)

func Inject() {
	injection.Inject(message.NewQueueFactory)
}

// InjectIns 使用实例化的方式注入消息队列
func InjectIns(inj *injection.Injector) {
	inj.Inject(message.NewQueueFactory)
}
