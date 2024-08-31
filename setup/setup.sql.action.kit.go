package setup

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
)

type ActionManager struct {
	actionMap map[string]ExecutableSQLAction
}

type ExecutableSQLAction interface {
	ExecuteQuery(ctx context.Context, args interface{}) (interface{}, error)
	String() string
}

func NewSQLActionManager(logger log.Logger, stringers []fmt.Stringer) *ActionManager {
	helper := log.NewHelper(logger, log.WithMessageKey("SQLAction"))
	actions := make(map[string]ExecutableSQLAction)

	for _, a := range stringers {
		if action, ok := a.(ExecutableSQLAction); ok {
			if _, exists := actions[action.String()]; exists {
				helper.Warnf("SQLAction duplicate with key %s", action.String())
			}

			actions[action.String()] = action
		}
	}

	return &ActionManager{
		actionMap: actions,
	}
}

// Find 根据 Key 查找具有相同 Key 的 SQLAction, 当 key 不存在时返回 nil
func (a *ActionManager) Find(key string) ExecutableSQLAction {
	action := a.actionMap[key]
	return action
}
