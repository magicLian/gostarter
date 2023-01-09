package log

type Logger interface {
	Debugf(msg string, ctx ...interface{})
	Infof(msg string, ctx ...interface{})
	Warnf(msg string, ctx ...interface{})
	Errorf(msg string, ctx ...interface{})
	Fatalf(msg string, ctx ...interface{})
}
