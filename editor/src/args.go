package main

import (
	arg "github.com/alexflint/go-arg"
)

type Args struct {
	Path    string `arg:"positional"`
	Config  string `arg:"-c" help:"optional location for config file"`
	Version bool   `arg:"-v" help:"prints version and exits"`
	Dump    bool   `arg:"-d" help:"dumps the current configuration and exists"`
	Log     bool   `arg:"-l" help:"immediately shows the debug log"`
}

func ParseArgs() {
	arg.MustParse(&state.Args)
	// TODO config file handling
}
