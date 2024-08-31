package config

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
)

type Manager struct {
	Client      *etcd.Client
	Instances   []*managerInstance
	TraceInfo   traceInfo
	WatchPath   []string
	PathHistory map[string]*pathPriorityHistory
}

// NewConfigWatcher 创建配置中心监听器, 他依赖于本地配置提供的配置中心地址
func NewConfigWatcher(configure *LocalConfigure) (ConfigureWatcherRepo, error) {
	// TODO: 将 ETCD  Client 抽取出来作为独立的依赖项
	timeOut, err := time.ParseDuration(configure.ConfigCenter.Timeout)
	if err != nil {
		return nil, err
	}

	client, err := etcd.New(etcd.Config{
		Endpoints:   configure.ConfigCenter.Endpoints,
		Username:    configure.ConfigCenter.Username,
		Password:    configure.ConfigCenter.Password,
		DialTimeout: timeOut,
	})

	if err != nil {
		panic(fmt.Sprint("connect to etcd failed with error", err))
		return nil, err
	}

	e := Manager{
		Instances:   make([]*managerInstance, 0),
		Client:      client,
		PathHistory: make(map[string]*pathPriorityHistory),
	}

	e.AddPrefix("", 0, false) // 添加默认监听器

	return &e, nil
}

// changeCallback 处理配置变动，维护 path 对应的各个优先级配置，当高优先级被删除时，替换回低优先级配置
func (c *Manager) changeCallback(event ChangeEvent) error {
	history, exists := c.PathHistory[event.Key]
	if !exists {
		history = newPathHistory(event.Key, event.Priority, event.Value)
	}

	if event.EventType == EventTypeDelete {
		history.removePriority(event.Priority)
	} else if event.EventType == EventTypePut {
		history.setValue(event.Priority, event.Value)
	}

	value, err := history.getValue()
	if err == nil {
		err = c.TraceInfo.setValue(event.Key, value)
	}

	// if something happen, keep previews value
	if err == nil {
		c.PathHistory[event.Key] = history
	}

	return err

}

func (c *Manager) AddPrefix(prefix string, priority int, ignoreEmpty bool) {
	if len(c.Instances) == 2 { // TODO: 暂时只支持两层优先级
		return
	}

	ins := newConfigManagerInstance(c.Client, prefix, priority, c.changeCallback)
	if ignoreEmpty {
		ins.IgnoreEmpty()
	}
	c.Instances = append(c.Instances, ins)
}

func (c *Manager) Load(object interface{}) error {
	fakeRoot := newTraceObject(nil)
	err := c.TraceInfo.analyseObj(reflect.ValueOf(object), fakeRoot)

	tree := newMergeTree(1)
	for _, o := range c.TraceInfo.objects {
		tree.addNodes(o.pathArray)
	}

	pathNeedWatch := tree.lookCommonPath(1)
	if err != nil {
		return err
	}

	for _, pathArray := range pathNeedWatch {
		path := "/" + strings.Join(pathArray, "/")
		if slices.Index(c.WatchPath, path) == -1 {
			c.WatchPath = append(c.WatchPath, path)
		}
	}

	return nil
}

func (c *Manager) LoadWithPath(object interface{}, path string) error {
	annoObj := struct {
		Obj interface{}
	}{
		Obj: object,
	}
	t := reflect.TypeOf(annoObj)
	v := reflect.ValueOf(&annoObj)
	f := t.Field(0)
	vv := v.Elem().Field(0)

	reflect.ValueOf(&annoObj.Obj)

	traceObj := &traceObject{
		field:       &f,
		value:       &vv,
		path:        path,
		pathArray:   strings.Split(path, "/"),
		objectType:  TagNamePath,
		valueFields: nil,
	}

	c.TraceInfo.addPathField(traceObj)
	c.WatchPath = append(c.WatchPath, traceObj.path)

	return nil
}

func (c *Manager) Start() error {

	if c.TraceInfo.service != nil {
		// 统一 service 前缀为 /service/{service_name}
		prefix := "/service/" + c.TraceInfo.service.path
		c.AddPrefix(prefix, 1, true)
	}

	var errs []error
	for _, path := range c.WatchPath {
		for _, ins := range c.Instances {
			err := ins.addPath(path)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

func (c *Manager) LoadAndStart(object interface{}) error {
	err := c.Load(object)
	if err == nil {
		err = c.Start()
	}

	return err
}
