/**
 * 日志二次包装工具
 * @author duhaifeng
 * @date   2021/04/15
 */
package common

import (
	log "github.com/duhaifeng/loglet"
)

var Log = newLogWrapper()

func InitLog() {
	Log = newLogWrapper()
}

func newLogWrapper() *LogWrapper {
	logger := new(LogWrapper)
	originLogger := log.NewLogger()
	if Configs != nil {
		originLogger.SetLogLevel(Configs.System.LogLevel)
		logConfigs := make(map[string]string)
		logConfigs["log_level"] = Configs.Log.LogLevel
		logConfigs["writers"] = Configs.Log.LogOutput
		logConfigs["log_file"] = Configs.Log.LogFilePath
		logConfigs["file_number"] = Configs.Log.LogFileNumber
		logConfigs["max_size"] = Configs.Log.LogFileSize
		originLogger.Init(logConfigs)
	}
	originLogger.SetLogPositionOffset(1) //避免日志点都打在当前文件上，对日志打印点增加一个偏移量
	logger.SetOriginalLogger(originLogger)
	return logger
}

type LogWrapper struct {
	logger *log.Logger
}

func (this *LogWrapper) SetOriginalLogger(logger *log.Logger) {
	this.logger = logger
}

func (this *LogWrapper) GetOriginalLogger() *log.Logger {
	return this.logger
}

func (this *LogWrapper) Debug(content string, contentArgs ...interface{}) {
	this.logger.Debug(RequestIdHolder.GetRoutineReqId()+" "+content, contentArgs...)
}

func (this *LogWrapper) Info(content string, contentArgs ...interface{}) {
	this.logger.Info(RequestIdHolder.GetRoutineReqId()+" "+content, contentArgs...)
}

func (this *LogWrapper) Warn(content string, contentArgs ...interface{}) {
	this.logger.Warn(RequestIdHolder.GetRoutineReqId()+" "+content, contentArgs...)
}

func (this *LogWrapper) Error(content string, contentArgs ...interface{}) {
	this.logger.Error(RequestIdHolder.GetRoutineReqId()+" "+content, contentArgs...)
}
