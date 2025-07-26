package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type Buffers struct {
	Log    *Buffer
	Prompt *Buffer
	Files  []*Buffer
}

type Times struct {
	BootTime *Timer
}

type State struct {
	Quit             bool
	Delta            int64
	Args             Args
	LogFile          *os.File
	InputFile        string
	Project          string
	ProjectFile      string
	ConfigFile       string
	DefaultConfig    Config
	Config           Config
	Theme            Theme
	CoreBinds        Binds
	Binds            Binds
	Screen           *Screen
	StatusBar        *StatusBar
	Panes            *Panes
	Areas            []*Area
	NamedAreas       map[string]*Area
	FocusedArea      *Area
	ZIndexedAreas    map[int][]*Area
	Buffers          Buffers
	DisplayLogBuffer bool
	Prompt           *Prompt
	Times            Times
}

var state *State

// order matters
func InitState() {
	state = &State{
		Buffers: Buffers{
			Log:   InitLogBuffer(),
			Files: []*Buffer{},
		},
	}

	ParseArgs()
	if state.Args.Version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	InitLogFile()

	state.CoreBinds = CoreBinds()
	state.Binds = DefaultBinds()
	state.DefaultConfig = DefaultConfig()
	state.Config = DefaultConfig()
	ConfigInit()

	state.NamedAreas = make(map[string]*Area)

	state.Buffers = Buffers{
		Log:   state.Buffers.Log,
		Files: []*Buffer{},
	}

	state.Times = Times{
		BootTime: NewTimer(),
	}
	state.Times.BootTime.Start()

	if state.Args.Log {
		state.DisplayLogBuffer = true
	}

	state.Times.BootTime.Log("Boot completed in: %s")

	if state.Args.Dump {
		DumpState()
		os.Exit(0)
	}

	if DEBUG {
		GetActionProgress()
	}

}

func DumpState() {
	state.Buffers.Log = &Buffer{}
	spew.Dump(state)
}
