package logger

import (
	"os"
)

func New(logfile string) *os.File {
	os.Remove(logfile)
	f, _ := os.Create(logfile)
	return f
}

var AuthServiceLogger = New("./log/authservice.log")
var DirectoryServiceLogger = New("./log/directoryservice.log")
var FileServerLogger = New("./log/fileserver.log")
