package config

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd "go.etcd.io/etcd/client/v3"
)

const (
	EventTypePut    = 0
	EventTypeDelete = 1
)

type ChangeEvent struct {
	EventType int
	Priority  int
	Prefix    string
	Key       string
	FullKey   string
	Value     []byte
}

type ChangeCallback func(event ChangeEvent) error

type managerInstance struct {
	Path             []string       // Path 为需要监听的地址
	Client           *etcd.Client   // Client 为用来建立监听的客户端
	Prefix           string         // Prefix 为当前实例在监听时需要为所有 path 统一添加的前缀
	Priority         int            // Priority 是该实例的优先级
	Callback         ChangeCallback // Callback 为发生数据变化时的通知通道
	IgnoreEmptyCheck bool           // IgnoreEmptyCheck 设置是否忽略未配置的选项, 该配置为 true 时遇到缺失的配置项不抛出错误
}

func newConfigManagerInstance(client *etcd.Client, prefix string, priority int, callback ChangeCallback) *managerInstance {
	e := managerInstance{
		Path:             make([]string, 0),
		Client:           client,
		Priority:         priority,
		Prefix:           prefix,
		Callback:         callback,
		IgnoreEmptyCheck: false,
	}

	return &e
}

func (e *managerInstance) IgnoreEmpty() {
	e.IgnoreEmptyCheck = true
}

// addPath 传递需要监听的 path，path 为完整的路径，如 /middleware/database/mysql
func (e *managerInstance) addPath(path string) error {

	exists := slices.Index(e.Path, path) != -1
	if !exists {
		//return errors.New(fmt.Sprintf("path %s already added to config center", path))
		// Path 当前为调试用，实际功能并无用到
		e.Path = append(e.Path, path)
		//return nil
	}

	if len(e.Prefix) != 0 {
		path = e.Prefix + path
	}

	err := e.getValue(path)
	if !exists {
		e.watchChan(path)
	}
	return err
}

// watchChan 监听 path 的配置变化，并在数据发生变化时通过 traceInfo 分析的结果更新对应的数据
// 传递的 path 是不包含 "/" 起始符的路径, prefix 也为不包含 / 起始符的路径
func (e *managerInstance) watchChan(path string) {

	prevRevision := int64(-1)
	c := e.Client.Watch(context.Background(), path, etcd.WithPrefix())

	go func() {
		for {
			msg := <-c
			if msg.Err() != nil {
				fmt.Printf("Watching %s failed with error %s\n", path, msg.Err())
			}

			if msg.CompactRevision < prevRevision {
				continue
			}

			prevRevision = msg.CompactRevision

			fmt.Printf("IsProgressNotify %t", msg.IsProgressNotify())
			for _, event := range msg.Events {
				_ = e.Callback(e.newEvent(event.Kv, event.Type))
			}
		}
	}()
}

func (e *managerInstance) getValue(path string) error {
	resp, err := e.Client.Get(context.Background(), path, etcd.WithPrefix())
	if err != nil {
		return fmt.Errorf("reading config %s from etcd with error %s", path, err)
	}

	if resp.Count == 0 && !e.IgnoreEmptyCheck {
		return fmt.Errorf("path %s is empty, please configure before use", path)
	}

	for _, kv := range resp.Kvs {
		event := e.newEvent(kv, etcd.EventTypePut)
		err = e.Callback(event)

		if err != nil {
			break
		}
	}

	return err
}

func (e *managerInstance) newEvent(kv *mvccpb.KeyValue, event mvccpb.Event_EventType) ChangeEvent {
	eventType := EventTypeDelete
	if event == etcd.EventTypePut {
		eventType = EventTypePut
	}

	return ChangeEvent{
		EventType: eventType,
		Priority:  e.Priority,
		Prefix:    e.Prefix,
		FullKey:   string(kv.Key),
		Key:       strings.Replace(string(kv.Key), e.Prefix, "", 1),
		Value:     kv.Value,
	}
}
