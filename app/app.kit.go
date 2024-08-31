package apputil

import (
	"strings"

	apppkg "github.com/eden/go-kratos-pkg/app"

	common "github.com/eden/go-biz-kit/common/def"
)

const (
	_appIDSep      = ":"
	_configPathSep = "/"
)

// ID 程序ID
// 例：go-srv-services/DEVELOP/main/v1.0.0/user-service
func ID(appConfig *common.App) string {
	return appIdentifier(appConfig, _appIDSep)
}

// ConfigPath 配置路径；用于配置中心，如：consul、etcd、...
// @result = app.BelongTo + "/" + app.RuntimeEnv + "/" + app.Branch + "/" + app.Version + "/" + app.Name
// 例：go-srv-services/DEVELOP/main/v1.0.0/user-service
func ConfigPath(appConfig *common.App) string {
	return appIdentifier(appConfig, _configPathSep)
}

// appIdentifier app 唯一标准
// @result = app.BelongTo + "/" + app.RuntimeEnv + "/" + app.Branch + "/" + app.Version + "/" + app.Name
// 例：go-srv-services/DEVELOP/main/v1.0.0/user-service
func appIdentifier(appConfig *common.App, sep string) string {
	var ss = make([]string, 0, 5)
	ss = append(ss, ParseEnv(appConfig.Env).String())
	if appConfig.Env != "" {
		branchString := strings.Replace(appConfig.Env, " ", ":", -1)
		ss = append(ss, branchString)
	}
	if appConfig.Version != "" {
		ss = append(ss, appConfig.Version)
	}
	if appConfig.Name != "" {
		ss = append(ss, appConfig.Name)
	}
	return strings.Join(ss, sep)
}

// SetRuntimeEnv ...
func SetRuntimeEnv(appEnv common.RuntimeEnvEnum_RuntimeEnv) {
	rv := apppkg.RuntimeEnvProduction
	switch appEnv {
	case common.RuntimeEnvEnum_LOCAL:
		rv = apppkg.RuntimeEnvLocal
	case common.RuntimeEnvEnum_DEVELOP:
		rv = apppkg.RuntimeEnvDevelop
	case common.RuntimeEnvEnum_TESTING:
		rv = apppkg.RuntimeEnvTesting
	case common.RuntimeEnvEnum_PREVIEW:
		rv = apppkg.RuntimeEnvPreview
	case common.RuntimeEnvEnum_PRODUCTION:
		rv = apppkg.RuntimeEnvProduction
	}
	apppkg.SetRuntimeEnv(rv)
}

// IsLocalMode ...
func IsLocalMode() bool {
	return apppkg.GetRuntimeEnv() == apppkg.RuntimeEnvLocal
}

// IsDebugMode ...
func IsDebugMode() bool {
	return apppkg.IsDebugMode()
}

// ParseEnv ...
func ParseEnv(appEnv string) (envEnum common.RuntimeEnvEnum_RuntimeEnv) {
	envInt32, ok := common.RuntimeEnvEnum_RuntimeEnv_value[strings.ToUpper(appEnv)]
	if ok {
		envEnum = common.RuntimeEnvEnum_RuntimeEnv(envInt32)
	}
	if envEnum == common.RuntimeEnvEnum_UNKNOWN {
		envEnum = common.RuntimeEnvEnum_PRODUCTION
		return envEnum
	}
	return envEnum
}
