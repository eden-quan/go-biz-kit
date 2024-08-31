package message

import (
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/eden-quan/go-biz-kit/config"
	"github.com/eden-quan/go-biz-kit/config/def"
)

type factoryImpl struct {
	logger   log.Logger             // 日志记录器
	conf     *def.Configuration     // 配置中心
	local    *config.LocalConfigure // 本地配置
	conn     *amqp.Connection       // amqp 的连接
	exitChan chan bool              // 退出信号
	lock     sync.Mutex             // 获取及创建连接的同步锁
}

type connProxy struct {
	f *factoryImpl
}

func (c *connProxy) Get() *amqp.Connection {
	return c.f.GetConn()
}

// NewQueueFactory 创建一个消息队列的管理器，用于创建消费者及生产者,
// logger 提供日志记录能力， conf 为配置中心
func NewQueueFactory(logger log.Logger, conf *def.Configuration, local *config.LocalConfigure) (QueueFactory, error) {
	f := &factoryImpl{
		logger:   log.With(logger, "mq", "rabbitmq"),
		exitChan: make(chan bool),
		conf:     conf,
		local:    local,
	}

	err := f.reconnect()
	if err == nil {
		err = f.watchConn()
	}
	return f, err
}

func (f *factoryImpl) reconnect() error {
	f.lock.Lock()
	defer f.lock.Unlock()

	conn, err := amqp.DialConfig(f.conf.MessageQueue.GetAddresses(), amqp.Config{
		Vhost:     f.conf.MessageQueue.GetVhost(),
		Heartbeat: f.conf.MessageQueue.GetHeartbeat().AsDuration(),
		//Properties: amqp.Table{},
	})

	if err == nil {
		f.conn = conn
	}
	return err
}

func (f *factoryImpl) watchConn() error {
	if f.conn == nil || f.conn.IsClosed() {
		return amqp.ErrClosed
	}

	receiver := make(chan *amqp.Error)
	f.conn.NotifyClose(receiver)
	go func() {
		for {
			select {
			case err, ok := <-receiver:
				if err == nil && ok {
					continue
				}

				conErr := f.reconnect()

				if conErr != nil {
					time.Sleep(time.Second * 10)
				}

				receiver = make(chan *amqp.Error)
				f.conn.NotifyClose(receiver)
			case <-time.After(time.Second * 30):
				log.NewHelper(f.logger).Debug("rabbit connection monitor alive...")
			}
		}
	}()

	return nil
}
func (f *factoryImpl) Producer(topic *TopicConfig) (Producer, error) {
	return NewRabbitProducer(f.local, f.logger, &connProxy{f: f}, topic)
}

func (f *factoryImpl) Consumer(topic *TopicConfig) (Consumer, error) {
	return NewRabbitConsumer(f.local, f.logger, &connProxy{f: f}, topic)
}

// GetConn 获取连接，通过 lock 确保连接出错时堵塞该操作
func (f *factoryImpl) GetConn() *amqp.Connection {
	f.lock.Lock()
	defer f.lock.Unlock()

	return f.conn
}
