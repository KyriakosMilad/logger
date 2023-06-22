package logger

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	DefaultFormat   = "${now} ${traceCode} ${counter} ${level} ${funcName}[${fileName}:${lineNumber}] ${value}"
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

func New(consolePrint bool, outputLogsDir string, createLogFileIfNotExist bool, traceCode string, format string) *Logger {
	if outputLogsDir[len(outputLogsDir)-1] != '/' {
		outputLogsDir += "/"
	}
	logsDir := outputLogsDir + time.Now().UTC().Format("20060102") + ".log"
	return &Logger{
		consolePrint:            consolePrint,
		outputLogFile:           logsDir,
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
		"lineNumber": fmt.Sprintf("%d", line),
		"value":      s,
	}
	l := lgr.replaceVariables(lgr.format, vars) + "\n"
	if lgr.consolePrint {
		log.Print(l)
	}
	if lgr.outputLogFile != "" {
		lgr.writeToFile(l, lgr.outputLogFile)
	}
	lgr.counter += 1
}

func (lgr *Logger) writeToFile(l string, file string) {
	flag := os.O_APPEND | os.O_WRONLY
	if lgr.createLogFileIfNotExist {
		flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}
	var err error
	f, err := os.OpenFile(file, flag, 0644)
	if err != nil {
		panic(fmt.Sprintf("error writing to log file %s: %s", file, err))
	}

	// close the log file after writing
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(fmt.Sprintf("error closing log file %s: %s", file, err))
		}
	}(f)

	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(l)
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

func (lgr *Logger) replaceVariables(input string, variables map[string]string) string {
	var output strings.Builder
	variableStart := -1
	insideVariable := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		if insideVariable {
			if char == '}' {
				insideVariable = false
				variableName := input[variableStart+2 : i]
				if value, ok := variables[variableName]; ok {
					output.WriteString(value)
				} else {
					output.WriteString("${" + variableName + "}")
				}
			}
		} else {
			if char == '$' && i+1 < len(input) && input[i+1] == '{' {
				insideVariable = true
				variableStart = i
				i++ // Skip '{' character
			} else {
				output.WriteByte(char)
			}
		}
	}

	return output.String()
}

func AutoGenerateTraceCode(prefix string, length uint8) string {
	const digits = "abcdefghijklmnopqrstuvwxyz0123456789"

	code := make([]byte, length)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}

	if prefix != "" {
		return prefix + "." + string(code)
	}
	return string(code)
}
