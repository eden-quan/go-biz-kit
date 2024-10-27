package message

import (
	"encoding/json"
	"time"
)

type Message struct {
	MessageId       string    // 消息的标识符，生产者在发送时设置，如果无设置，则会生成随机字符串
	Queue           string    // 消息来自哪个 Queue
	Exchange        string    // 消息通过哪个交换机发送到当前的 Queue
	RoutingKey      string    // 消息通过什么 RoutingKey 对应的规则路由到当前 Queue
	Priority        int       // 消息的优先级，由生产者设置，范围为 0-9，数值约大优先级越高
	Payload         []byte    // 消息内容
	Timestamp       time.Time // 消息发送时的时间戳
	App             string    // 消息来自于哪个应用
	ContentType     string    // MIME content type, 默认情况下都为 application/json
	ContentEncoding string    // MIME content encoding, 默认情况下都为 utf8
}

// UnmarshalPayload 将 Body 中的数据反序列化到 obj 对应的对象中
func (m *Message) UnmarshalPayload(obj interface{}) error {
	return json.Unmarshal(m.Payload, obj)
}

func (m *Message) Ack() {
	// TODO: 暂时默认所有消息都为自动提交的
}
