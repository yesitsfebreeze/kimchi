package main

import (
	"os"
	"path/filepath"
	"strings"
)

func OpenFile(path string) (*Context, bool) {
	content, err := os.ReadFile(path)

	if err != nil {
		LogErr("Failed to read file: ", err)
		return nil, false
	}

	absPath, err := filepath.Abs(path)

	if err != nil {
		LogErr("Failed to get absolute path for: ", path)
	}

	buf := &Buffer{
		Name:     filepath.Base(path),
		Path:     path,
		AbsPath:  absPath,
		Modified: false,
	}

	ctx := &Context{
		Buffer: buf,
		Cursors: []*Cursor{
			&Cursor{
				Position: Vec2{X: 0, Y: 0},
			},
		},
	}

	lines := strings.SplitSeq(string(content), "\n")
	for line := range lines {
		buf.AppendLine([]rune(line))
	}
	buf.Modified = false // explicitly reset

	state.Buffers.Files = append(state.Buffers.Files, buf)

	return ctx, true
}

func SaveBuffer(args ...any) {
	// WithBuffer(func(buf *Buffer) {
	// 	// TODO: implement me
	// })
	// buf := &state.Buffers[state.CurrentBuffer]
	// var lines []string
	// for _, l := range buf.Lines {
	// 	lines = append(lines, string(l.data))
	// }
	// err := os.WriteFile(state.Startfile, []byte(strings.Join(lines, "\n")), 0644)
	// if err != nil {
	// 	LogErr("Failed to save buffer: ", err)
	// 	return
	// }
	// buf.Modified = false
	// Log("Saved file: ", state.Startfile)
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

		DaemonSend(DaemonHighlight, DaemonData{
			Lang: "go",
			Code: ctx.Buffer.ToString(),
		}, func(res string, err error) {
			if err != nil {
				LogErr("Failed to highlight code:", err)
				return
			}
			LogF("Daemon response: %s", res)
		})

	}
}
