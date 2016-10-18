package logger

import (
	"fmt"
	"io"
	"log/syslog"
	"os"
	"time"
)

type Logger struct {
	Syslog bool
	Debug  bool

	writer LogWriter
}

type LogWriter interface {
	Debug(message string) error
	Info(message string) error
	Emerg(message string) error
}

type StdWriter struct {
	Out io.Writer
}

const (
	StdColorDebug = 33
	StdColorInfo  = 32
	StdColorEmerg = 31
)

func (writer *StdWriter) Debug(message string) error {
	return writer.output(message, StdColorDebug)
}

func (writer *StdWriter) Info(message string) error {
	return writer.output(message, StdColorInfo)
}

func (writer *StdWriter) Emerg(message string) error {
	return writer.output(message, StdColorEmerg)
}

func (writer *StdWriter) output(message string, color int) error {
	fmt.Fprintf(writer.Out, "%v \033[%dm%s\033[39m\n", time.Now(), color, message)
	return nil

}

var Log *Logger = &Logger{}

func (logger *Logger) Writer() LogWriter {
	if logger.writer == nil {
		if logger.Syslog {
			syslogWriter, err := syslog.New(syslog.LOG_DAEMON, "edwig")
			if err != nil {
				panic("Can't write to syslog")
			}
			logger.writer = syslogWriter
		} else {
			logger.writer = &StdWriter{Out: os.Stderr}
		}
	}
	return logger.writer
}

func (logger *Logger) Debugf(format string, values ...interface{}) {
	if logger.Debug {
		logger.Writer().Debug(fmt.Sprintf(format, values...))
	}
}

func (logger *Logger) Printf(format string, values ...interface{}) {
	logger.Writer().Info(fmt.Sprintf(format, values...))
}

func (logger *Logger) Panicf(format string, values ...interface{}) {
	message := fmt.Sprintf(format, values...)
	logger.Writer().Emerg(message)
	panic(message)
}
