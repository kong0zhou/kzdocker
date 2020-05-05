package log

import (
	"fmt"

	"go.uber.org/zap"
)

// Debug 正常调试输出，只会输出到终端
func Debug(msg string, fields ...zap.Field) {
	fields = append([]zap.Field{zap.Int64(`goid`, getGroutineID())}, fields...)
	log.Debug(msg, fields...)
}

// Debugf 格式输出
func Debugf(format string, val ...interface{}) {
	log.Debug(fmt.Sprintf(format, val...), zap.Int64(`goid`, getGroutineID()))
}

// Info 正常输出，会输出到文件，和其他
func Info(msg string, fields ...zap.Field) {
	fields = append([]zap.Field{zap.Int64(`goid`, getGroutineID())}, fields...)
	log.Info(msg, fields...)
}

// Infof 格式输出
func Infof(format string, val ...interface{}) {
	log.Info(fmt.Sprintf(format, val...), zap.Int64(`goid`, getGroutineID()))
}

// Warn 警告输出
func Warn(msg string, fields ...zap.Field) {
	fields = append([]zap.Field{zap.Int64(`goid`, getGroutineID())}, fields...)
	log.Warn(msg, fields...)
}

// Warnf 格式输出
func Warnf(format string, val ...interface{}) {
	log.Warn(fmt.Sprintf(format, val...), zap.Int64(`goid`, getGroutineID()))
}

// Error 错误输出
func Error(msg string, fields ...zap.Field) {
	fields = append([]zap.Field{zap.Int64(`goid`, getGroutineID())}, fields...)
	log.Error(msg, fields...)
}

// Errorf 格式输出
func Errorf(format string, val ...interface{}) {
	log.Error(fmt.Sprintf(format, val...), zap.Int64(`goid`, getGroutineID()))
}

// DPanic 开发状态使用，如果logger处于开发模式，它就会死机
func DPanic(msg string, fields ...zap.Field) {
	fields = append([]zap.Field{zap.Int64(`goid`, getGroutineID())}, fields...)
	log.DPanic(msg, fields...)
}

// DPanicf 格式输出
func DPanicf(format string, val ...interface{}) {
	log.DPanic(fmt.Sprintf(format, val...), zap.Int64(`goid`, getGroutineID()))
}

// Panic 打印日志，然后panic
func Panic(msg string, fields ...zap.Field) {
	fields = append([]zap.Field{zap.Int64(`goid`, getGroutineID())}, fields...)
	log.Panic(msg, fields...)
}

// Panicf 格式输出
func Panicf(format string, val ...interface{}) {
	log.Panic(fmt.Sprintf(format, val...), zap.Int64(`goid`, getGroutineID()))
}

// Fatal 打印日志，然后os.Exit(1)
func Fatal(msg string, fields ...zap.Field) {
	fields = append([]zap.Field{zap.Int64(`goid`, getGroutineID())}, fields...)
	log.Fatal(msg, fields...)
}

// Fatalf 格式输出
func Fatalf(format string, val ...interface{}) {
	log.Fatal(fmt.Sprintf(format, val...), zap.Int64(`goid`, getGroutineID()))
}

// Sync calls the underlying Core's Sync method, flushing any buffered log entries. Applications should take care to call Sync before exiting.
var Sync = log.Sync
