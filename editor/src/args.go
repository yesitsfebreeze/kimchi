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

func OpenInputFile() {
	if !IsFile(state.Args.Path) {
		return
	}

	if ctx, ok := OpenFile(state.Args.Path); ok {
		if state.Panes.One == nil || state.Panes.Two == nil {
			InitPanes()
		}

		// state.Panes.Layout = PaneLayoutVertical

		// TODO: this should be the state
		state.Panes.One.Area.SetContext(ctx)
		state.Panes.One.Visible = true
		state.Panes.One.Area.Focus()

		state.Panes.Two.Area.SetContext(ctx)
		state.Panes.Two.Visible = true
		state.Panes.Two.Area.Focus()

		result, err := DaemonSend(DaemonHighlight, DaemonData{
			Lang: "go",
			Code: ctx.Buffer.ToString(),
		})

		if err == nil {
			Log("OpenInputFile: %s", result)
		}

	}
}
