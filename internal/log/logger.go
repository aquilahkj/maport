package log

// Logger log interface
type Logger interface {
	AddLogPrefix(string)
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{}) error
	Fatal(string, ...interface{})
	Panic(string, ...interface{})
}
