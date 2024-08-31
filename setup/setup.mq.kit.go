package setup

import (
	"github.com/go-kratos/kratos/v2/log"

	kit "github.com/eden-quan/go-biz-kit"
	"github.com/eden-quan/go-biz-kit/config"
	"github.com/eden-quan/go-biz-kit/config/def"
	"github.com/eden-quan/go-biz-kit/message"
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
