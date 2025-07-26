package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

var tempLogLines []string

const LOG_BUFFER_NAME = "KITLOG"
const MAX_LOG_LINES = 10

func InitLogBuffer() *Buffer {
	buf := &Buffer{
		Name:     LOG_BUFFER_NAME,
		MaxLines: MAX_LOG_LINES,
	}

	// If buffer just created, flush temp lines
	if len(tempLogLines) > 0 {
		for _, l := range tempLogLines {
			state.Buffers.Log.AppendLine([]rune(l))
		}
	}

	return buf
}

func CloseLogFile() {
	if state.LogFile != nil {
		state.LogFile.Close()
	}
}

func InitLogFile() {
	if !DEBUG {
		return
	}
	exePath, err := os.Executable()
	if err != nil {
		LogErr("Could not determine executable path:", err)
		return
	}
	exeDir := filepath.Dir(exePath)
	logPath := filepath.Join(exeDir, "kitsune.log")

	state.LogFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0644)
	if err != nil {
		LogErr("Could not open log file:", err)
		return
	}

	log.SetOutput(state.LogFile)

	if len(tempLogLines) > MAX_LOG_LINES {
		tempLogLines = tempLogLines[len(tempLogLines)-MAX_LOG_LINES:]
	}

	for _, line := range tempLogLines {
		log.Println(line)
	}
	tempLogLines = nil
}

func appendLogLine(line string) {
	// Buffer may not be initialized yet
	if state.Buffers.Log == nil {
		tempLogLines = append(tempLogLines, line)
		if len(tempLogLines) > MAX_LOG_LINES {
			tempLogLines = tempLogLines[1:]
		}
	} else {
		state.Buffers.Log.AppendLine([]rune(line))
		state.Buffers.Log.Modified = true

		// Truncate log buffer if it exceeds max
		if len(state.Buffers.Log.Lines) > MAX_LOG_LINES {
			excess := len(state.Buffers.Log.Lines) - MAX_LOG_LINES
			state.Buffers.Log.Lines = state.Buffers.Log.Lines[excess:]
		}
	}

	if !DEBUG {
		return
	}

	state.LogFile.Truncate(0)
	state.LogFile.Seek(0, 0)

	// Redump current log buffer
	if state.Buffers.Log != nil {
		for _, line := range state.Buffers.Log.Lines {
			str := string(line.ToRunes())
			log.Println(str)
		}
	}
}

func CreateLogMsg(prefix string, args ...any) string {
	parts := make([]string, len(args))
	for i, arg := range args {
		parts[i] = fmt.Sprint(arg)
	}
	msg := strings.Join(parts, " ")

	if prefix != "" {
		msg = fmt.Sprintf("[%s] %s", prefix, msg)
	}

	if DEBUG {
		if _, file, line, ok := runtime.Caller(2); ok {
			short := filepath.Base(file)
			msg = fmt.Sprintf("(%s:%d) %s", short, line, msg)
		}
	}

	return msg
}

func Log(args ...any) {
	msg := CreateLogMsg("", args...)
	appendLogLine(msg)
}

func LogF(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	msg = CreateLogMsg("log", msg)
	appendLogLine(msg)
}

func LogErr(args ...any) {
	msg := CreateLogMsg("error", args...)
	appendLogLine(msg)
}

func LogErrF(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	msg = CreateLogMsg("error", msg)
	appendLogLine(msg)
}

func DumpLog() {
	spew.Dump(state.Buffers.Log.Lines)
}
