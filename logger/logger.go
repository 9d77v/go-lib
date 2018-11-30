package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

//some log levels
const (
	LevelDEBUG = "[DEBUG]"
	LevelINFO  = "[INFO]"
	LevelERROR = "[ERROR]"
	LevelWARN  = "[WARN]"
)

//Logger custom logger
type Logger struct {
	*log.Logger
}

//NewLog is the default log for the app
//when the mode is prod,only save logs to files
//when the mode is dev,output logs to both file and console
func NewLog(logPath string, mode string) *Logger {
	logger := new(Logger)
	os.Mkdir(path.Dir(logPath), os.ModePerm)
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	writers := make([]io.Writer, 0)
	writers = append(writers, f)
	if mode != "prod" {
		writers = append(writers, os.Stdout)
		fileAndStdoutWriter := io.MultiWriter(writers...)
		logger.Logger = log.New(fileAndStdoutWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
		return logger
	}
	logger.Logger = log.New(f, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
	return logger
}

//output log msgs
func (l *Logger) output(level string, args ...interface{}) {
	l.Logger.SetPrefix(level)
	l.Output(3, fmt.Sprintln(args...))
}

//DEBUG print DEBUG msg
func (l *Logger) DEBUG(args ...interface{}) {
	l.output(LevelDEBUG, args...)
}

//INFO print INFO msg
func (l *Logger) INFO(args ...interface{}) {
	l.output(LevelINFO, args...)
}

//WARN print INFO msg
func (l *Logger) WARN(args ...interface{}) {
	l.output(LevelWARN, args...)
}

//ERROR print INFO msg
func (l *Logger) ERROR(args ...interface{}) {
	l.output(LevelERROR, args...)
}

//outputf log msgs
func (l *Logger) outputf(level, format string, args ...interface{}) {
	l.Logger.SetPrefix(level)
	l.Output(3, fmt.Sprintf(format, args...))
}

//DEBUGf print DEBUG msg
func (l *Logger) DEBUGf(format string, args ...interface{}) {
	l.outputf(LevelDEBUG, format, args...)
}

//INFOf print INFO msg
func (l *Logger) INFOf(format string, args ...interface{}) {
	l.outputf(LevelINFO, format, args...)
}

//WARNf print INFO msg
func (l *Logger) WARNf(format string, args ...interface{}) {
	l.outputf(LevelWARN, format, args...)
}

//ERRORf print INFO msg
func (l *Logger) ERRORf(format string, args ...interface{}) {
	l.outputf(LevelERROR, fmt.Sprintf(format, args))
}
