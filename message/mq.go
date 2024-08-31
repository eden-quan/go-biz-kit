package message

import (
	"context"

	"github.com/eden/go-biz-kit/config/def"
)

/*
    -- 全局配置
	全局配置 -> 配置 RocketMQ 的 NamedServer 信息
	Topic 配置 -> 配置 Topic 的名称 / 类型 / 标签，可通过配置中心，按照 Topic 名称进行查找，得到指定 Topic 的配置信息
                 类型配置：定时时间 / 等
    默认客户端配置 (是否需要?)

    -- 个性化配置
	客户端 配置 -> 每个 Topic 一个配置, 消费者与生产者共用，配置客户端的信息，每个客户端需要消费 Topic 时，都需要在自身的服务目录下配置客户端信息
				   该配置信息包括客户端 ID / 消费者组 (可选随机模式，用于实现广播) / 标签信息 (可选)
*/

type TopicConfig struct {
	Exchange *def.ExchangeConfig
	Queue    *def.QueueConfig
	Consume  *def.ConsumeConfig
}

type PushOption interface {
	Do(opt *defaultOpt)
}

type PullOption interface {
	Do(opt *defaultOpt)
}

// Producer 为已与指定 Topic 绑定的生产者, Push 会根据 Topic 的配置来触发不同的行为，
// 如发送的信息根据配置可能是广播/延迟信息等
type Producer interface {
	// Push 往指定的 Topic 推送 data 消息, 所有推送的消息都会通过 Json.Marshal 序列化为 JSON 后进行传输
	Push(context context.Context, data interface{}, option ...PushOption) error
}

type Consumer interface {
	// Puller 返回一个用户获取数据的 Channel, 在有数据可消费时该 Channel 会返回一条类型为 Message 的消息
	// 使用者自行根据消息进行f处理 (反序列化等)
	Puller(option ...PullOption) <-chan *Message

	// Pull 返回一条接收到的数据，在没有消息可消费时该函数会 Block,
	// 与 Producer 对应，拉取的数据是序列化后的 Json，因此消费者需要通过 Json.Unmarshal
	// 将数据反序列化为所需的目标类型
	Pull(context context.Context, option ...PullOption) (*Message, error)
}

// QueueFactory 提供 Producer 及 Consumer 的创建能力，
// 当前可根据 Topic 获取对应的消费者及生产者，其具体的行为由配置信息来决定，后续有复杂需求时再开放灵活度更高的消息队列接口
type QueueFactory interface {
	// Producer 创建一个绑定了 topic 的生产者
	Producer(topic *TopicConfig) (Producer, error)
	// Consumer 创建一个绑定了 topic 的消费者
	Consumer(topic *TopicConfig) (Consumer, error)
}

// NewProducerConfig 创建一个新的生产者配置
func NewProducerConfig(ex *def.ExchangeConfig, qe *def.QueueConfig) *TopicConfig {
	return &TopicConfig{
		Exchange: ex,
		Queue:    qe,
	}
}

// NewConsumerConfig 创建一个新的消费者配置
func NewConsumerConfig(qe *def.QueueConfig, con *def.ConsumeConfig) *TopicConfig {
	return &TopicConfig{
		Queue:   qe,
		Consume: con,
	}
}
