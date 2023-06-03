package logger

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	DefaultFormat   = "$now $traceCode $counter $level $funcName[$fileName:$lineNumber] $value"
	LogErrorLevel   = "ERROR"
	LogInfoLevel    = "INFO"
	LogWarningLevel = "WARNING"
)

type Logger struct {
	consolePrint            bool
	outputLogFile           string
	createLogFileIfNotExist bool
	traceCode               string
	counter                 int
	format                  string
}

func New(consolePrint bool, outputLogFile string, createLogFileIfNotExist bool, traceCode string, format string) *Logger {
	return &Logger{
		consolePrint:            consolePrint,
		outputLogFile:           outputLogFile,
		createLogFileIfNotExist: createLogFileIfNotExist,
		traceCode:               traceCode,
		counter:                 0,
		format:                  format,
	}
}

func (lgr *Logger) LogError(s string) {
	lgr.Log(s, LogErrorLevel, 1)
}

func (lgr *Logger) LogInfo(s string) {
	lgr.Log(s, LogInfoLevel, 1)
}

func (lgr *Logger) LogWarning(s string) {
	lgr.Log(s, LogWarningLevel, 1)
}

func (lgr *Logger) LogInnerError(s string, skip int) {
	lgr.Log(s, LogErrorLevel, skip+1)
}

func (lgr *Logger) LogInnerInfo(s string, skip int) {
	lgr.Log(s, LogInfoLevel, skip+1)
}

func (lgr *Logger) LogInnerWarning(s string, skip int) {
	lgr.Log(s, LogWarningLevel, skip+1)
}

func (lgr *Logger) Log(s string, level string, skip int) {
	pc, filename, line, _ := runtime.Caller(skip + 1)
	vars := map[string]string{
		"now":        time.Now().UTC().Format(time.RFC3339),
		"traceCode":  lgr.traceCode,
		"counter":    fmt.Sprintf("%04d", lgr.counter),
		"level":      level,
		"funcName":   runtime.FuncForPC(pc).Name(),
		"fileName":   filename,
		"lineNumber": string(rune(line)),
		"value":      s,
	}
	l := lgr.replaceVariables(lgr.format, vars) + "\n"
	if lgr.consolePrint {
		log.Print(l)
	}
	if lgr.outputLogFile != "" {
		lgr.writeToFile(l, lgr.outputLogFile)
	}
}

func (lgr *Logger) writeToFile(l string, file string) {
	var f *os.File
	if lgr.createLogFileIfNotExist {
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

func (lgr *Logger) ResetCounter() {
	lgr.counter = 0
}

func (lgr *Logger) replaceVariables(str string, variables map[string]string) string {
	var builder strings.Builder
	start := 0

	for {
		placeholderStart := strings.IndexByte(str[start:], '$')
		if placeholderStart == -1 {
			builder.WriteString(str[start:])
			break
		}

		builder.WriteString(str[start : start+placeholderStart])

		placeholderEnd := strings.IndexByte(str[start+placeholderStart:], ' ')
		if placeholderEnd == -1 {
			placeholderEnd = len(str)
		} else {
			placeholderEnd += start + placeholderStart
		}

		placeholder := str[start+placeholderStart : placeholderEnd]

		if replacement, exists := variables[placeholder]; exists {
			builder.WriteString(replacement)
		} else {
			builder.WriteString(placeholder)
		}

		start = placeholderEnd
	}

	return builder.String()
}

func autoGenerateTraceCode(prefix string) string {
	const digits = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	code := make([]byte, 6)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}

	return prefix + "." + string(code)
}
