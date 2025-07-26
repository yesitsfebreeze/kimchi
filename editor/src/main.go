package main

import (
	"fmt"
	"os"
	"runtime"
)

const DEBUG = true

func Shutdown() {
	CloseLogFile()
	ResetScreen()
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "\nPANIC: %v\n", r)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, true)
		fmt.Fprintf(os.Stderr, "Stack trace:\n%s\n", buf[:stackSize])
	}
}

func main() {
	defer Shutdown()

	InitState()
	InitScreen()
	Init()
	StartRender()
}

func Init() {
	InitStatusBar()
	InitPanes()
	InitPrompt()
	OpenInputFile()
}

func Update() {
	UpdatePanes()
	UpdateStatusBar()
	UpdatePrompt()
}
