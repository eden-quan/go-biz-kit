package setup

import (
	"io"
	stdlog "log"
	"strings"
	"sync"

	logpkg "github.com/eden/go-kratos-pkg/log"
	writerpkg "github.com/eden/go-kratos-pkg/writer"

	"github.com/eden/go-biz-kit/config/def"
)

// getLoggerFileWriter 文件日志写手柄
func (m *LoggerManager) getLoggerFileWriter() (io.Writer, error) {

	var err error
	m.loggerFileWriterMutex.Do(func() {
		m.loggerFileWriter, err = m.loadingLoggerFileWriter()
	})

	if err != nil {
		m.loggerFileWriterMutex = sync.Once{}
	}
	return m.loggerFileWriter, err
}

// loadingLoggerFileWriter 启动日志文件写手柄
func (m *LoggerManager) loadingLoggerFileWriter() (io.Writer, error) {

	fileConf := m.conf.Log.GetFile()
	if !fileConf.GetEnable() {
		stdlog.Println("|*** 加载：日志工具：虚拟的文件写手柄")
		return writerpkg.NewDummyWriter()
	}

	rotateConfig := &writerpkg.ConfigRotate{
		Dir:            fileConf.GetDir(),
		Filename:       fileConf.GetFilename(),
		RotateTime:     fileConf.GetRotateTime().AsDuration(),
		RotateSize:     fileConf.GetRotateSize(),
		StorageCounter: uint(fileConf.GetStorageCounter()),
		StorageAge:     fileConf.GetStorageAge().AsDuration(),
	}

	fileName := []string{rotateConfig.Filename}
	replaceHandler := strings.NewReplacer(
		" ", "-",
		"/", "--",
	)
	if m.localConf.APP.Env != "" {
		fileName = append(fileName, replaceHandler.Replace(m.localConf.APP.Env))
	}
	if m.localConf.APP.Version != "" {
		fileName = append(fileName, replaceHandler.Replace(m.localConf.APP.Version))
	}

	rotateConfig.Filename = strings.Join(fileName, "_")

	return writerpkg.NewRotateFile(rotateConfig)
}

// getLoggerGraylogWriter graylog日志写手柄
func (m *LoggerManager) getLoggerGraylogWriter() (io.Writer, error) {
	var err error
	m.loggerGraylogWriterMutex.Do(func() {
		m.loggerGraylogWriter, err = m.loadingLoggerGraylogWriter()
	})
	if err != nil {
		m.loggerGraylogWriterMutex = sync.Once{}
	}
	return m.loggerGraylogWriter, err
}

// loadingLoggerGraylogWriter graylog日志文件写手柄
func (m *LoggerManager) loadingLoggerGraylogWriter() (io.Writer, error) {
	graylogLoggerConfig := m.conf.Log.GetGraylog()
	if !graylogLoggerConfig.GetEnable() {
		return writerpkg.NewDummyWriter()
	}

	graylogConfig := m.genGraylogConfig(graylogLoggerConfig)
	return logpkg.NewGraylogWriter(graylogConfig)
}

// genGraylogConfig ...
func (m *LoggerManager) genGraylogConfig(graylogLoggerConfig *def.Log_Graylog) *logpkg.GraylogConfig {
	graylogConfig := &logpkg.GraylogConfig{
		Facility:      graylogLoggerConfig.Facility,
		Proto:         graylogLoggerConfig.Proto,
		Addr:          graylogLoggerConfig.Addr,
		AsyncPoolSize: int(graylogLoggerConfig.AsyncPoolSize),
	}
	return graylogConfig
}
