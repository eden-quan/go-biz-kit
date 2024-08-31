package config

import "fmt"

// pathPriorityHistory 负责根据优先级管理 key 及其对应的 value，实现高优先级配置覆盖低优先级配置的能力，
// 以此提供在基于基础配置来减轻个业务系统配置复杂度的前提下，为各个服务提供个性化的能力
type pathPriorityHistory struct {
	Path    string
	History map[int][]byte // History 以优先级作为 Key， 配置值作为 Value
}

func newPathHistory(path string, priority int, value []byte) *pathPriorityHistory {
	his := &pathPriorityHistory{
		Path:    path,
		History: make(map[int][]byte),
	}

	his.setValue(priority, value)

	return his
}

func (p *pathPriorityHistory) setValue(priority int, value []byte) {
	p.History[priority] = value
}

func (p *pathPriorityHistory) getValue() ([]byte, error) {
	maxPriority := -1
	var value []byte = nil

	for p, v := range p.History {
		if p > maxPriority {
			maxPriority = p
			value = v
		}
	}

	if maxPriority == -1 {
		return nil, fmt.Errorf("value of %s doesn't exists", p.Path)
	}

	return value, nil
}

func (p *pathPriorityHistory) removePriority(priority int) {
	_, exists := p.History[priority]
	if exists {
		delete(p.History, priority)
	}
}
