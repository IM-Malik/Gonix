package logger

import (
	"log"
	"os"
	"io"
)

type Logger struct {
	Info *log.Logger
	Warn *log.Logger
	Error *log.Logger
}

