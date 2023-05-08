package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	ConsolePrint            = false
	OutputLogFile           = ""
	CreateLogFileIfNotExist = true
	TraceCode               = ""
	Counter                 = 0
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
	now := time.Now().UTC().Format(time.RFC3339)
	l := fmt.Sprintf("%s %s %04d %s %s[%s:%d] %v\n", now, TraceCode, Counter, level, runtime.FuncForPC(pc).Name(), filename, line, s)
	if ConsolePrint {
		log.Printf(l)
	}
	if OutputLogFile != "" {
		writeToFile(l, OutputLogFile)
	}
}

func writeToFile(l string, file string) {
	var f *os.File
	if CreateLogFileIfNotExist {
		err := os.MkdirAll(filepath.Dir(file), os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("error creating log file path %s: %s", file, err))
		}
		f, err = os.Create(file)
		if err != nil {
			panic(fmt.Sprintf("error creating log file %s: %s", file, err))
		}
	} else {
		var err error
		f, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
	}

	// close the log file after writing
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(fmt.Sprintf("error closing log file %s: %s", file, err))
		}
	}(f)

	writer := bufio.NewWriter(f)
	_, err := writer.WriteString(l)
	if err != nil {
		panic(fmt.Sprintf("error writing to log file %s: %s", file, err))
	}
	// make sure all data is written
	err = writer.Flush()
	if err != nil {
		panic(fmt.Sprintf("error flushing log file %s: %s", file, err))
	}
}

func ResetCounter() {
	Counter = 0
}
