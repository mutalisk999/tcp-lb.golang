package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	l        *log.Logger
	logLevel int
}

var (
	DEBUG = 10
	INFO  = 20
	WARN  = 30
	ERROR = 40

	logSetLevel = 10

	Debug *Logger
	Info  *Logger
	Warn  *Logger
	Error *Logger
)

func InitLog(infoFile, errorFile string, setLevel int) {
	logSetLevel = setLevel
	log.Println("logSetLevel:", setLevel)

	log.Println("infoFile:", infoFile)
	log.Println("errorFile:", errorFile)

	infoFd, err := os.OpenFile(infoFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open info log file:", err)
	}

	errorFd, err := os.OpenFile(errorFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Debug = &Logger{log.New(io.MultiWriter(os.Stdout, infoFd), "[D] ", log.Ldate|log.Lmicroseconds|log.Lshortfile), DEBUG}
	Info = &Logger{log.New(io.MultiWriter(os.Stdout, infoFd), "[I] ", log.Ldate|log.Lmicroseconds|log.Lshortfile), INFO}
	Warn = &Logger{log.New(io.MultiWriter(os.Stderr, infoFd, errorFd), "[W] ", log.Ldate|log.Lmicroseconds|log.Lshortfile), WARN}
	Error = &Logger{log.New(io.MultiWriter(os.Stderr, infoFd, errorFd), "[E] ", log.Ldate|log.Lmicroseconds|log.Lshortfile), ERROR}
}

func (l *Logger) Print(v ...interface{}) {
	if logSetLevel <= l.logLevel {
		_ = l.l.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Println(v ...interface{}) {
	if logSetLevel <= l.logLevel {
		_ = l.l.Output(2, fmt.Sprintln(v...))
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	if logSetLevel <= l.logLevel {
		_ = l.l.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	_ = l.l.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func (l *Logger) Fatalln(v ...interface{}) {
	_ = l.l.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	_ = l.l.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	_ = l.l.Output(2, s)
	panic(s)
}

func (l *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	_ = l.l.Output(2, s)
	panic(s)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	_ = l.l.Output(2, s)
	panic(s)
}
