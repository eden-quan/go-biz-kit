package message

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/eden-quan/go-biz-kit/config"
)

type consumer struct {
	conn     *connProxy
	channel  *amqp.Channel
	logger   log.Logger
	topic    *TopicConfig
	conf     *config.LocalConfigure
	queue    *amqp.Queue
	delivery <-chan amqp.Delivery
	msgChan  chan *Message
	lock     sync.Mutex
}

func NewRabbitConsumer(conf *config.LocalConfigure, logger log.Logger, conn *connProxy, topic *TopicConfig) (Consumer, error) {
	c := consumer{
		conn:    conn,
		channel: nil,
		logger:  logger,
		topic:   topic,
		conf:    conf,
		msgChan: make(chan *Message),
	}

	err := c.reconnect()
	if err == nil {
		err = c.watchChan()
	}
	return &c, err
}

func (c *consumer) reconnect() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	topic := c.topic
	conn := c.conn.Get()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		topic.Queue.GetName(), topic.Queue.GetDurable(), topic.Queue.GetAutoDelete(),
		false, false, nil,
	)

	if err != nil {
		return fmt.Errorf("[MQ] Declare Queue %s with error %s", topic.Queue.GetName(), err)
	}

	delivery, err := ch.Consume(q.Name, topic.Consume.GetName(), topic.Consume.GetAutoAck(),
		false, false, topic.Consume.GetNoWait(), nil)

	if err == nil {
		qosErr := ch.Qos(int(topic.Consume.QosPrefetchCount), int(topic.Consume.QosPrefetchSize), false)
		if qosErr != nil {
			err = fmt.Errorf("[MQ] Setting Qos for topic(%s, %s) failed with error %s", topic.Queue.GetName(), topic.Exchange.GetName(), qosErr)
		}
	}

	c.channel = ch
	c.queue = &q
	c.delivery = delivery

	if err == nil {
		go c.startConsumer()
	}

	return err
}

func (c *consumer) startConsumer() {
	curChan := c.channel
	for {
		select {
		case msg := <-c.delivery:
			if msg.Body == nil && msg.Exchange == "" && curChan.IsClosed() {
				return
			}
			c.msgChan <- &Message{
				MessageId:       msg.MessageId,
				Queue:           c.queue.Name,
				Exchange:        msg.Exchange,
				RoutingKey:      msg.RoutingKey,
				Priority:        int(msg.Priority),
				Payload:         msg.Body,
				Timestamp:       msg.Timestamp,
				App:             msg.AppId,
				ContentType:     msg.ContentType,
				ContentEncoding: msg.ContentEncoding,
			}
		}
	}
}

// watchChan 用于监听 Channel 的关闭状态，实现自动重连的能力
// TODO: 该函数可以抽取为通用函数同时供 Producer 及 Consumer 使用
func (c *consumer) watchChan() error {
	if c.channel == nil || c.channel.IsClosed() {
		return amqp.ErrClosed
	}

	receiver := make(chan *amqp.Error)
	c.channel.NotifyClose(receiver)
	go func() {
		for {
			select {
			case err, ok := <-receiver:
				if err == nil && ok {
					continue
				}

				conErr := c.reconnect() // create new consumer
				if conErr != nil {
					time.Sleep(10 * time.Second)
				}

				receiver = make(chan *amqp.Error)
				c.channel.NotifyClose(receiver)

			case <-time.After(time.Second * 30):
				log.NewHelper(c.logger).Debugf(
					"consumer (queue(%s):consumer(%s) channel monitor alive...",
					c.topic.Queue.GetName(), c.topic.Consume.GetName(),
				)
			}
		}
	}()

	return nil
}

func (c *consumer) Pull(context context.Context, _ ...PullOption) (*Message, error) {
	select {
	case msg := <-c.msgChan:
		return msg, nil
	case <-context.Done():
		return nil, context.Err()
	}
}

func (c *consumer) Puller(_ ...PullOption) <-chan *Message {
	if c.channel.IsClosed() {
		ch := make(chan *Message)
		close(ch)
		return ch
	}

	return c.msgChan
}
