package logging

import (
	"log"
	"os"
	"strconv"
)

type Log struct {
	Verbosity int
}

const (
	Fatal int = iota
	Error
	Warning
	Info
	Debug
)

func NewLogger() *Log {
	verbosity := Info
	v, isSet := os.LookupEnv("VERBOSITY")
	if isSet {
		verbosity, _ = strconv.Atoi(v)
	}

	return &Log{
		Verbosity: verbosity,
	}
}

func println(level string, msg string) {
	log.Println(level, msg)
}

func (l *Log) Debug(s string) {
	if l.Verbosity >= Debug {
		println("DEBUG", s)
	}
}

func (l *Log) Info(s string) {
	if l.Verbosity >= Info {
		println("INFO", s)
	}
}

func (l *Log) Warning(s string) {
	if l.Verbosity >= Warning {
		println("WARNING", s)
	}
}
func (l *Log) Error(s error) {
	if l.Verbosity >= Error {
		println("ERROR", s.Error())
	}
}

func (l *Log) Fatal(s error) {
	log.Panicln("FATAL " + s.Error())
}
