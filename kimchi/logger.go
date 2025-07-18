package main

import (
	"fmt"
	"strings"
)

const MAX_LOG_LINES = 256

var log_buffer = &Buffer{
	name:    "[Log]",
	content: [][]rune{},
}

func __create_log_msg(prefix string, args ...interface{}) string {
	msg := fmt.Sprint(append([]interface{}{""}, args...)...) // adds space at start
	msg = strings.TrimSpace(msg)

	return fmt.Sprintf("[%s] %s", prefix, msg)
}

func log(args ...interface{}) {
	msg := __create_log_msg("log", args...)
	append_log_line(msg)
}

func logf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	append_log_line(__create_log_msg("log", msg))

	fmt.Println(msg)
}

func log_err(args ...interface{}) {
	msg := __create_log_msg("error", args...)
	append_log_line(msg)
}

func logf_err(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	append_log_line(__create_log_msg("error", msg))
}

func append_log_line(line string) {
	runes := []rune(line)

	log_buffer.content = append(log_buffer.content, runes)

	if len(log_buffer.content) > MAX_LOG_LINES {
		log_buffer.content = log_buffer.content[len(log_buffer.content)-MAX_LOG_LINES:]
	}

	log_buffer.modified = true
}

func dump_all_logs() {
	fmt.Println("\n--- LOG DUMP ---")
	for _, line := range log_buffer.content {
		fmt.Println(string(line))
	}
}

func render_log() {
	start := len(log_buffer.content) - 5
	if start < 0 {
		start = 0
	}

	for _, msg := range log_buffer.content[start:] {
		move_cursor(0, 0)
		clear_line()
		fmt.Print(string(msg))
	}
}
