package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrFConfigFileDoesNotExist = "config file does not exist: %s"
	ErrFInvalidConfigField     = "invalid config field: %s"
	ErrFUnkownConfigKey        = "unknown config key: %s"
	ErrFInvalidIndentStyle     = "invalid indent style: %s"
	ErrNoBuffers               = "no buffers available"
	ErrNoBufferFocus           = "no buffer selected"
	ErrNoFocusedAreas          = "no focused areas"
	ErrBufferOOB               = "buffer index out of range"
	ErrFUnkownLuaType          = "unknown lua type%s"
)

func Err(msg string) error {
	return errors.New(msg)
}

func ErrF(msg string, args ...any) error {
	return fmt.Errorf(msg, args...)
}

func Throw(err error) {
	ThrowErr(err.Error())
}

func ThrowErr(msg string) {
	fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", msg)
	os.Exit(1)
}

func ThrowErrF(msg string, args ...any) {
	ThrowErr(fmt.Sprintf(msg, args...))
}
