package log

import (
	"log"
	"os"
)

type Logger struct {
    ErrorLog *log.Logger
    InfoLog *log.Logger
}

func New() *Logger {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

    return &Logger{
        ErrorLog: errorLog,
        InfoLog: infoLog,
    }
}
