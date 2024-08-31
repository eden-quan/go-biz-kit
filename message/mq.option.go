package message

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	headerpkg "github.com/eden-quan/go-kratos-pkg/header"
	uuidpkg "github.com/eden-quan/go-kratos-pkg/uuid"
)

type defaultOpt struct {
	exchangeName string
	routingKey   string
	mandatory    bool
	immediately  bool
	msg          amqp.Publishing
}

func (p *producer) defaultOption() *defaultOpt {
	return &defaultOpt{
		exchangeName: p.topic.Exchange.GetName(),
		routingKey:   p.topic.Queue.GetRoutingKey(),
		mandatory:    false,
		immediately:  false,
		msg: amqp.Publishing{
			Headers:         nil,
			ContentType:     headerpkg.ContentTypeJSON,
			ContentEncoding: "gzip",
			DeliveryMode:    2, // persistent
			Priority:        0,
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       uuidpkg.New(),
			Timestamp:       time.Now(),
			Type:            "",
			UserId:          "",
			AppId:           p.conf.APP.Name,
			Body:            nil,
		},
	}
}

type PriorityPushOption struct {
	Priority uint8
}

// WithPriority 为 Push 操作的消息增加优先级信息, 优先级可以是 0-9 的数值，数值越大优先级越高
func WithPriority(priority uint8) PriorityPushOption {
	return PriorityPushOption{Priority: priority}
}

func (p *PriorityPushOption) Do(opt *defaultOpt) {
	opt.msg.Priority = p.Priority
}

type RoutingKeyPushOption struct {
	RoutingKey string
}

// WithRoutingKey 为消息发送增加路由信息，允许消息根据路由提交到指定的消费者
func WithRoutingKey(key string) RoutingKeyPushOption {
	return RoutingKeyPushOption{RoutingKey: key}
}

func (p *RoutingKeyPushOption) Do(opt *defaultOpt) {
	opt.routingKey = p.RoutingKey
}

type ExchangePushOption struct {
	ExchangeName string
}

// WithExchange 使用指定的交换机替换 Producer 配置中的交换机
func WithExchange(exchange string) ExchangePushOption {
	return ExchangePushOption{ExchangeName: exchange}
}

func (p *ExchangePushOption) Do(opt *defaultOpt) {
	opt.exchangeName = p.ExchangeName
}

type MsgIdPushOption struct {
	MessageId string
}

func WithMessageId(msgId string) MsgIdPushOption {
	return MsgIdPushOption{msgId}
}

func (p *MsgIdPushOption) Do(opt *defaultOpt) {
	opt.msg.MessageId = p.MessageId
}
