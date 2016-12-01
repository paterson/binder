package logger

import (
	"io"
	"os"
)

type Logger struct {
	file *os.File
}

func (logger Logger) Write(p []byte) (n int, err error) {
	logger.file.WriteString(string(p))
	return 0, nil
}

func New(logfile string) io.Writer {
	os.Remove(logfile)
	os.Create(logfile)
	f, _ := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY, 0600)
	return &Logger{file: f}
}
