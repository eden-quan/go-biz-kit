package setup

import (
	"github.com/go-kratos/kratos/v2/log"

	kit "gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git"
	"gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/config"
	"gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/config/def"
	"gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/message"
)

type messageQueueImpl struct {
	factory message.QueueFactory
}

func (m *messageQueueImpl) Get() message.QueueFactory {
	return m.factory
}

func NewMessageQueue(logger log.Logger, conf *def.Configuration, local *config.LocalConfigure) (kit.MessageQueue, error) {
	factory, err := message.NewQueueFactory(logger, conf, local)
	return &messageQueueImpl{factory: factory}, err
}
