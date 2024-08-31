package message

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/eden-quan/go-biz-kit/config"
)

// Producer 为已与指定 Topic 绑定的生产者, Push 会根据 Topic 的配置来触发不同的行为，
// 如发送的信息根据配置可能是广播/延迟信息等
type producer struct {
	conn    *connProxy
	channel *amqp.Channel
	queue   *amqp.Queue
	topic   *TopicConfig
	lock    sync.Mutex
	logger  log.Logger
	conf    *config.LocalConfigure
}

func NewRabbitProducer(conf *config.LocalConfigure, logger log.Logger, conn *connProxy, topic *TopicConfig) (Producer, error) {
	p := &producer{
		conn:   conn,
		topic:  topic,
		logger: logger,
		conf:   conf,
	}

	err := p.reconnect()
	if err == nil {
		err = p.watchChan()
	}
	return p, err
}

func (p *producer) watchChan() error {
	if p.conn.Get() == nil || p.conn.Get().IsClosed() {
		return amqp.ErrClosed
	}

	receiver := make(chan *amqp.Error)
	p.channel.NotifyClose(receiver)
	go func() {
		for {
			select {
			case err, ok := <-receiver:
				if err == nil && ok {
					continue
				}

				conErr := p.reconnect()
				if conErr != nil {
					time.Sleep(10 * time.Second)
				}

				receiver = make(chan *amqp.Error)
				p.channel.NotifyClose(receiver)
			case <-time.After(time.Second * 30):
				log.NewHelper(p.logger).Debugf(
					"producer (exchange(%s):queue(%s) channel monitor alive...",
					p.topic.Exchange.GetName(), p.topic.Queue.GetName())
			}
		}
	}()

	return nil
}

func (p *producer) reconnect() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	topic := p.topic
	conn := p.conn.Get()

	c, err := conn.Channel()
	if err == nil {
		err = c.ExchangeDeclare(
			p.topic.Exchange.GetName(), topic.Exchange.GetKind(), topic.Exchange.GetDurable(),
			p.topic.Exchange.GetAutoDelete(), false, false, nil)
	}

	if err != nil {
		return err
	}

	queue, err := c.QueueDeclare(
		topic.Queue.GetName(), topic.Queue.GetDurable(), topic.Queue.GetAutoDelete(),
		false, false, nil,
	)

	if err == nil {
		err = c.QueueBind(queue.Name, topic.Queue.GetRoutingKey(), topic.Exchange.GetName(), false, nil)
	}

	p.channel = c
	p.queue = &queue

	return err
}

func (p *producer) Push(context context.Context, data interface{}, option ...PushOption) (err error) {

	body := make([]byte, 0)

	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return
		}
	}

	opt := p.defaultOption()
	opt.msg.Body = body

	for _, o := range option {
		o.Do(opt)
	}

	err = p.channel.PublishWithContext(context, opt.exchangeName, opt.routingKey, opt.mandatory, opt.immediately, opt.msg)
	return
}
