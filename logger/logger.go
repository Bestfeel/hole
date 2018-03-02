package logger

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"io"
)

type Logger struct {
	Log              *logrus.Logger
	LogPath          string
	LogName          string
	LogTimeSplit     string
	WithMaxAge       time.Duration
	WithRotationTime time.Duration
}

func New() *Logger {

	return &Logger{
		Log: &logrus.Logger{
			Out:       os.Stderr,
			Formatter: &logrus.TextFormatter{ForceColors: true, DisableColors: false, DisableTimestamp: false, FullTimestamp: true, DisableSorting: true},
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.InfoLevel,
		},
		LogPath:          "/tmp/hole.log",
		LogTimeSplit:     "%Y%m%d",
		WithMaxAge:       7 * 24 * time.Hour,
		WithRotationTime: 24 * time.Hour,
	}
}

var std = New()

func AddHook() {

	writer, err := rotatelogs.New(
		std.LogPath+"."+std.LogTimeSplit,
		rotatelogs.WithLinkName(std.LogPath),              // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(std.WithMaxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(std.WithRotationTime), // 日志切割时间间隔
	)

	if err != nil {
		logrus.Error(errors.WithStack(err))
	}

	lfHook := NewHook(WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{DisableTimestamp: false})

	std.Log.AddHook(lfHook)
}

func SetOutput(out io.Writer) {
	std.Log.Out = out
}
func SetFormatter(formatter logrus.Formatter) {
	std.Log.Formatter = formatter
}

func SetLogTimeSplit(timeSplit string) {
	std.LogTimeSplit = timeSplit
}

func SetLogPath(path string) {
	std.LogPath = path
}

func SetMaxAge(time time.Duration) {

	std.WithMaxAge = time
}

func SetRotationTime(time time.Duration) {

	std.WithRotationTime = time
}

func SetLevel(level string) {

	switch level {
	case "debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stderr)
	case "info", "INFO":
		setNull()
		logrus.SetLevel(logrus.InfoLevel)
	case "warn", "WARN":
		setNull()
		logrus.SetLevel(logrus.WarnLevel)
	case "error", "ERROR":
		setNull()
		logrus.SetLevel(logrus.ErrorLevel)
	case "panic", "PANIC":
		setNull()
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal", "FATAL":
		setNull()
		logrus.SetLevel(logrus.FatalLevel)
	default:
		setNull()
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	std.Log.Out = writer
}

func ErrorWithFields(msg interface{}, fields logrus.Fields) {

	std.Log.WithFields(fields).Error(msg)
}

func Error(msg interface{}) {

	std.Log.Error(msg)
}

func Errorf(format string, msg ...interface{}) {

	std.Log.Errorf(format, msg)
}

func InfoWithFields(msg interface{}, fields logrus.Fields) {

	std.Log.WithFields(fields).Info(msg)
}

func Info(msg interface{}) {

	std.Log.Info(msg)
}

func Infof(format string, msg ...interface{}) {

	std.Log.Infof(format, msg)
}

func DebugWithFields(msg interface{}, fields logrus.Fields) {

	std.Log.WithFields(fields).Debug(msg)
}

func Debug(data interface{}) {
	std.Log.Debug(data)
}

func Debugf(format string, msg ...interface{}) {

	std.Log.Debugf(format, msg)
}

func WarnWithFields(msg interface{}, fields logrus.Fields) {

	std.Log.WithFields(fields).Warn(msg)
}

func Warn(msg interface{}) {
	std.Log.Warn(msg)
}

func Warnf(format string, msg ...interface{}) {

	std.Log.Warnf(format, msg)
}

func Fatal(msg interface{}) {
	std.Log.Fatal(msg)
}

func FatalWithFields(msg interface{}, fields logrus.Fields) {

	std.Log.WithFields(fields).Fatal(msg)
}

func Fatalf(format string, msg ...interface{}) {

	std.Log.Fatalf(format, msg)
}

func PanicWithFields(msg interface{}, fields logrus.Fields) {

	std.Log.WithFields(fields).Panic(msg)
}

func Panic(msg interface{}) {

	std.Log.Panic(msg)
}

func Panicf(format string, msg ...interface{}) {

	std.Log.Panicf(format, msg)
}
