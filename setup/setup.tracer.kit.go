package setup

import (
	"github.com/go-kratos/kratos/v2/log"

	"gitlab.lainuoniao.cn/eden-quan/go-biz-kit/config"
	"gitlab.lainuoniao.cn/eden-quan/go-biz-kit/config/def"
	"gitlab.lainuoniao.cn/eden-quan/go-biz-kit/tracing"
)

// NewTracing 创建链路跟踪
func NewTracing(conf *def.Configuration, logger log.Logger, local *config.LocalConfigure) (*tracing.TracerInitializer, error) {
	return tracing.InitTracing(conf, logger, local)
}
