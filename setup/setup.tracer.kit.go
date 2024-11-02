package setup

import (
	"github.com/go-kratos/kratos/v2/log"

	"gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/config"
	"gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/config/def"
	"gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/tracing"
)

// NewTracing 创建链路跟踪
func NewTracing(conf *def.Configuration, logger log.Logger, local *config.LocalConfigure) (*tracing.TracerInitializer, error) {
	return tracing.InitTracing(conf, logger, local)
}
