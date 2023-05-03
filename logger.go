package logger

import (
	"log"
	"runtime"
)

const (
	LogErrorLevel   = "ERROR"
	LogInfoLevel    = "INFO"
	LogWarningLevel = "WARNING"
)

func LogError(s string) {
	Log(s, LogErrorLevel, 1)
}

func LogInfo(s string) {
	Log(s, LogInfoLevel, 1)
}

func LogWarning(s string) {
	Log(s, LogWarningLevel, 1)
}

func LogInnerError(s string, skip int) {
	Log(s, LogErrorLevel, skip+1)
}

func LogInnerInfo(s string, skip int) {
	Log(s, LogInfoLevel, skip+1)
}

func LogInnerWarning(s string, skip int) {
	Log(s, LogWarningLevel, skip+1)
}

func Log(s string, level string, skip int) {
	pc, filename, line, _ := runtime.Caller(skip + 1)
	log.Printf("[%s] %s[%s:%d] %v", level, runtime.FuncForPC(pc).Name(), filename, line, s)
}
