package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Log struct {
	Prefix string
	Ctx    []interface{}
}

func New(prefix string, ctx ...interface{}) Logger {
	return &Log{
		Prefix: prefix,
		Ctx:    ctx,
	}
}

func (log *Log) Tracef(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	logrus.Tracef("[%s] %s\n", log.Prefix, message)
}

func (log *Log) Debugf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	logrus.Debugf("[%s] %s\n", log.Prefix, message)

}

func (log *Log) Infof(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	logrus.Infof("[%s] %s\n", log.Prefix, message)
}

func (log *Log) Warnf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	logrus.Warnf("[%s] %s\n", log.Prefix, message)
}

func (log *Log) Errorf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	logrus.Errorf("[%s] %s\n", log.Prefix, message)
}

func (log *Log) Fatalf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	logrus.Fatalf("[%s] %s\n", log.Prefix, message)
}
