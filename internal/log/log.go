package log

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// func init() {
// 	log.SetFormatter(&logrus.TextFormatter{
// 		TimestampFormat: "2006-01-02 15:04:05", //时间格式化
// 	})
// }

// Config configuration of log
type Config struct {
	Target       string
	Level        string
	PrintCaller  bool
	Formatter    string
	DisableColor bool
}

// TpLogger logger model
type TpLogger struct {
	*logrus.Logger
	prefix string
}

// SetConfig set the config for logger
func SetConfig(config *Config) {

	switch config.Target {
	case "stdout", "":
		log.SetOutput(os.Stdout)
	default:
		file, err := os.OpenFile(config.Target, os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			panic("Failed to log to file")
		}
	}

	switch config.Level {
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	case "INFO", "":
		log.SetLevel(logrus.InfoLevel)
	case "WARN":
		log.SetLevel(logrus.WarnLevel)
	case "ERROR":
		log.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		log.SetLevel(logrus.FatalLevel)
	default:
		panic("The level config error")
	}

	switch config.Formatter {
	case "json":
		formatter := &logrus.JSONFormatter{
			TimestampFormat:  "2006-01-02 15:04:05",
			DisableTimestamp: false,
		}
		log.SetFormatter(formatter)
	case "text", "":
		formatter := &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		}
		formatter.DisableColors = config.DisableColor

		log.SetFormatter(formatter)
	default:
		panic("The formatter config error")
	}

	log.SetReportCaller(config.PrintCaller)

}

// NewLogger create a logger
func NewLogger(prefixes ...string) Logger {
	logger := &TpLogger{Logger: log}

	for _, p := range prefixes {
		logger.AddLogPrefix(p)
	}

	return logger
}

// AddLogPrefix add prefix message to the logger
func (lg *TpLogger) AddLogPrefix(prefix string) {
	lg.prefix += "[" + prefix + "]"
}

func (lg *TpLogger) format(fmtstr string, args ...interface{}) string {
	return fmt.Sprintf(lg.prefix+fmtstr, args...)
}

// Debug log debug level message
func (lg *TpLogger) Debug(arg0 string, args ...interface{}) {
	lg.Logger.Debug(lg.format(arg0, args...))
}

// Info log info level message
func (lg *TpLogger) Info(arg0 string, args ...interface{}) {
	lg.Logger.Info(lg.format(arg0, args...))
}

// Warn log warn level message
func (lg *TpLogger) Warn(arg0 string, args ...interface{}) {
	lg.Logger.Warn(lg.format(arg0, args...))
}

// Error log error level message and return a error
func (lg *TpLogger) Error(arg0 string, args ...interface{}) error {
	err := (lg.format(arg0, args...))
	lg.Logger.Error(err)
	return errors.New(err)
}

// Fatal log fatal level message and exit program
func (lg *TpLogger) Fatal(arg0 string, args ...interface{}) {
	lg.Logger.Fatal(lg.format(arg0, args...))
}

// Panic log panic level message and raise a panic
func (lg *TpLogger) Panic(arg0 string, args ...interface{}) {
	err := (lg.format(arg0, args...))
	lg.Logger.Panic(err)
	panic(err)
}
