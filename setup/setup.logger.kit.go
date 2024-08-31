package setup

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel/trace"

	contextpkg "github.com/eden-quan/go-kratos-pkg/context"
	headerpkg "github.com/eden-quan/go-kratos-pkg/header"
	ippkg "github.com/eden-quan/go-kratos-pkg/ip"
	logpkg "github.com/eden-quan/go-kratos-pkg/log"

	config2 "github.com/eden-quan/go-biz-kit/config"
	config "github.com/eden-quan/go-biz-kit/config/def"
)

// LoggerPrefixField with logger fields.
type LoggerPrefixField struct {
	AppName    string `json:"name"`
	AppVersion string `json:"version"`
	AppEnv     string `json:"env"`
	Hostname   string `json:"hostname"`
	ServerIP   string `json:"serverIP"`
}

// Prefix 日志前缀
func (s *LoggerPrefixField) Prefix() []interface{} {
	return []interface{}{
		"service", s.AppName,
		"app.hostname", s.Hostname,
		"app.env", s.AppEnv,
		"app.version", s.AppVersion,
	}
}

type LoggerManager struct {
	logger           log.Logger
	loggerMutex      sync.Once
	loggerHelper     log.Logger
	helperMutex      sync.Once
	loggerMiddleware log.Logger
	middlewareMutex  sync.Once

	loggerPrefixFieldMutex sync.Once
	loggerPrefixField      *LoggerPrefixField

	loggerFileWriter      io.Writer
	loggerFileWriterMutex sync.Once

	loggerGraylogWriter      io.Writer
	loggerGraylogWriterMutex sync.Once

	conf      *config.Configuration
	localConf *config2.LocalConfigure
}

// NewLogger 根据配置信息创建日志实例
// TODO: 增加 Fluent 支持，分离日志系统与代码的耦合
func NewLogger(manager *LoggerManager) (log.Logger, error) {
	return manager.Logger()
}

// NewLoggerManager 提供一个日志管理器，可以用他创建日志实例以及进行一些初始化操作
func NewLoggerManager(local *config2.LocalConfigure, configuration *config.Configuration) *LoggerManager {
	return &LoggerManager{
		conf:      configuration,
		localConf: local,
	}
}

// Logger 日志处理示例
func (m *LoggerManager) Logger() (log.Logger, error) {
	var err error
	m.loggerMutex.Do(func() {
		m.logger, err = m.loadingLogger()
	})
	if err != nil {
		m.loggerMutex = sync.Once{}
	}
	return m.logger, err
}

// LoggerHelper 日志处理示例
func (m *LoggerManager) LoggerHelper() (log.Logger, error) {
	var err error
	m.helperMutex.Do(func() {
		m.loggerHelper, err = m.loadingLoggerHelper()
	})
	if err != nil {
		m.helperMutex = sync.Once{}
	}
	return m.loggerHelper, err
}

// LoggerMiddleware 中间件的日志处理示例
func (m *LoggerManager) LoggerMiddleware() (log.Logger, error) {
	var err error
	m.middlewareMutex.Do(func() {
		m.loggerMiddleware, err = m.loadingLoggerMiddleware()
	})
	if err != nil {
		m.middlewareMutex = sync.Once{}
	}
	return m.loggerMiddleware, err
}

// loadingLogger 初始化日志输出实例
func (m *LoggerManager) loadingLogger() (logger log.Logger, err error) {
	skip := logpkg.CallerSkipForLogger
	//return s.loadingLoggerWithCallerSkip(skip)
	logger, err = m.loadingLoggerWithCallerSkip(skip)
	if err != nil {
		return logger, err
	}
	logger = m.withLoggerPrefix(logger)
	return logger, err
}

// loadingLoggerHelper 初始化日志工具输出实例
func (m *LoggerManager) loadingLoggerHelper() (logger log.Logger, err error) {
	skip := logpkg.CallerSkipForHelper
	//return s.loadingLoggerWithCallerSkip(skip)
	logger, err = m.loadingLoggerWithCallerSkip(skip)
	if err != nil {
		return logger, err
	}
	logger = m.withLoggerPrefix(logger)
	return logger, err
}

// loadingLoggerMiddleware 初始化中间价的日志输出实例
func (m *LoggerManager) loadingLoggerMiddleware() (logger log.Logger, err error) {
	skip := logpkg.CallerSkipForMiddleware
	//return s.loadingLoggerWithCallerSkip(skip)
	logger, err = m.loadingLoggerWithCallerSkip(skip)
	if err != nil {
		return logger, err
	}
	logger = m.withLoggerPrefix(logger)
	return logger, err
}

// loadingLoggerWithCallerSkip 初始化日志输出实例
func (m *LoggerManager) loadingLoggerWithCallerSkip(skip int) (logger log.Logger, err error) {
	// loggers
	var loggers []log.Logger

	// DummyLogger
	stdLogger, err := logpkg.NewDummyLogger()
	if err != nil {
		return logger, err
	}

	conf := &m.conf.Log

	// 日志 输出到控制台
	consoleConf := conf.GetConsole()
	if consoleConf.GetEnable() {
		stdLoggerConfig := &logpkg.ConfigStd{
			Level:          logpkg.ParseLevel(consoleConf.GetLevel()),
			CallerSkip:     skip,
			UseJSONEncoder: consoleConf.GetUseJsonEncoder(),
		}
		stdLoggerImpl, err := logpkg.NewStdLogger(stdLoggerConfig)
		if err != nil {
			return logger, err
		}
		stdLogger = stdLoggerImpl

		// 覆盖 stdLogger
		loggers = append(loggers, stdLogger)
	}

	// 日志 输出到文件
	loggerConfigForFile := conf.GetFile()
	if loggerConfigForFile.GetEnable() {
		// file logger
		fileLoggerConfig := &logpkg.ConfigFile{
			Level:      logpkg.ParseLevel(loggerConfigForFile.GetLevel()),
			CallerSkip: skip,

			Dir:      loggerConfigForFile.Dir,
			Filename: loggerConfigForFile.Filename,

			RotateTime: loggerConfigForFile.RotateTime.AsDuration(),
			RotateSize: loggerConfigForFile.RotateSize,

			StorageCounter: uint(loggerConfigForFile.StorageCounter),
			StorageAge:     loggerConfigForFile.StorageAge.AsDuration(),
		}
		writer, err := m.getLoggerFileWriter()
		if err != nil {
			return logger, err
		}
		fileLogger, err := logpkg.NewFileLogger(
			fileLoggerConfig,
			logpkg.WithWriter(writer),
		)

		loggers = append(loggers, fileLogger)
	}

	// 日志 输出到Graylog
	loggerConfigForGraylog := conf.GetGraylog()
	if loggerConfigForGraylog.GetEnable() {
		writer, err := m.getLoggerGraylogWriter()
		if err != nil {
			return logger, err
		}
		graylogLoggerConfig := &logpkg.ConfigGraylog{
			Level:         logpkg.ParseLevel(loggerConfigForGraylog.GetLevel()),
			CallerSkip:    skip,
			GraylogConfig: *m.genGraylogConfig(loggerConfigForGraylog),
		}
		graylogLogger, err := logpkg.NewGraylogLogger(
			graylogLoggerConfig,
			logpkg.WithWriter(writer),
		)
		if err != nil {
			return logger, err
		}
		loggers = append(loggers, graylogLogger)
	}

	// 日志工具
	if len(loggers) == 0 {
		return logger, err
	}
	return logpkg.NewMultiLogger(loggers...), err
}

// assemblyLoggerPrefixField 组装日志前缀
func (m *LoggerManager) assemblyLoggerPrefixField() *LoggerPrefixField {
	fields := &LoggerPrefixField{
		AppName:    m.localConf.APP.Name,
		AppVersion: m.localConf.APP.Version,
		AppEnv:     m.localConf.APP.Env,
		ServerIP:   ippkg.LocalIP(),
	}

	fields.Hostname, _ = os.Hostname()
	return fields
}

// LoggerPrefixField .
func (m *LoggerManager) LoggerPrefixField() *LoggerPrefixField {
	m.loggerPrefixFieldMutex.Do(func() {
		m.loggerPrefixField = m.assemblyLoggerPrefixField()
	})
	return m.loggerPrefixField
}

// withLoggerPrefix ...
func (m *LoggerManager) withLoggerPrefix(logger log.Logger) log.Logger {
	var kvs = m.LoggerPrefixField().Prefix()
	kvs = append(kvs, "tracer.trace_id", m.withLoggerTracerTraceId())
	kvs = append(kvs, "tracer.span_id", m.withLoggerTracerSpanId())
	return log.With(logger, kvs...)
}

// withLoggerTracerTraceId returns a tracing trace id valuer.
func (m *LoggerManager) withLoggerTracerTraceId() log.Valuer {
	return func(ctx context.Context) interface{} {

		span := trace.SpanContextFromContext(ctx)
		if span.HasTraceID() {
			return span.TraceID().String()
		}

		var tr transport.Transporter = nil
		tr, ok := contextpkg.MatchHTTPServerContext(ctx)
		if !ok {
			tr, ok = contextpkg.MatchGRPCServerContext(ctx)
		}

		if ok && tr != nil {
			return tr.RequestHeader().Get(headerpkg.RequestID)
		}

		return ""
	}
}

// withLoggerTracerSpanId returns tracing span id valuer.
func (m *LoggerManager) withLoggerTracerSpanId() log.Valuer {
	return func(ctx context.Context) interface{} {
		span := trace.SpanContextFromContext(ctx)
		if span.HasSpanID() {
			return span.SpanID().String()
		}

		return ""
	}
}
