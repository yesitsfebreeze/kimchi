package main

import (
	"fmt"
	"os"
	"runtime"
)

const DEBUG = true

func Shutdown() {
	DisconnectFromDaemon()
	ResetScreen()
	CloseLogFile()
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "\nPANIC: %v\n", r)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, true)
		fmt.Fprintf(os.Stderr, "Stack trace:\n%s\n", buf[:stackSize])
	}
}

// check if the deamon is running
// if not start it

func main() {
	defer Shutdown()
	InitState()
	InitScreen()
	InitDaemon()
	InitStatusBar()
	InitPanes()
	InitPrompt()
	OpenInputFile()
	state.Times.BootTime.Log("Boot completed in %s")
	StartRender()
}

func Update() {
	UpdatePanes()
	UpdateStatusBar()
	UpdatePrompt()
}
