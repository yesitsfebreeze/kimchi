package main

import (
	"fmt"
	"strings"
)

const StatusBarName = ":KitsuneStatusBar:"

type StatusBar struct {
	Pane   *Pane
	Area   *Area
	Buffer *Buffer
	Stats  *Stats
}

type Stats struct {
	FullPath string
	BaseName string
	Encoding string
	// git branch
	// git status
	// git commit
	// time
	// encoding
	// filetype
	// file size
	// file permissions (maybe)
	// file encoding
	// file modified time
	// Line:Col
	// Percentage of file scrolled
	// Line count
}

func InitStatusBar() {
	bar := &StatusBar{
		Buffer: &Buffer{
			Name: StatusBarName,
		},
		Area: CreateArea(StatusBarName).
			IgnoreClamp(true).
			SetZIndex(MID_ZINDEX).
			SetStyle(state.Theme.StatusBar),
	}

	bar.Area.SetContext(&Context{
		Buffer: bar.Buffer,
	})
	bar.SetSizeAndPosition()

	state.StatusBar = bar
}

func (b *StatusBar) SetSizeAndPosition() {
	var y int = 0
	height := 1
	if state.Config.Statusbar.UseSeparator {
		height = 2
	}

	if state.Config.Statusbar.Position == StatusbarTop {
		y = 0
	} else {
		y = state.Screen.Size.Y - height
	}
	if b.Area == nil {
		LogErr("StatusBar Area is nil")
		return
	}
	if b.Buffer == nil {
		LogErr("StatusBar Buffer is nil")
		return
	}

	b.Buffer.MaxLines = height
	b.Area.SetSize(Vec2{X: state.Screen.Size.X, Y: height})
	b.Area.SetPosition(Vec2{X: 0, Y: y})
}

func (b *StatusBar) CollectStats() {
	if b.Stats == nil {
		b.Stats = &Stats{}
	}

	// if state.FocusedArea != nil {
	// 	if state.FocusedArea.Buffer != nil {
	// 		buf := state.FocusedArea.Buffer
	// 		b.Stats.FullPath = buf.AbsPath
	// 		b.Stats.BaseName = filepath.Base(buf.AbsPath)
	// 	}
	// }
}

func (b *StatusBar) RenderContent() {

	s := b.Stats

	// TODO: restructure this correctly
	// to use the config layout
	// also check height, currently 2 but doesnt render properly

	if b.Area == nil || b.Buffer == nil {
		LogErr("StatusBar Area or Buffer is nil")
		return
	}

	width := b.Area.Size.X
	if width <= 0 {
		return
	}

	var (
		left   = fmt.Sprintf("%s | %s", "master*", s.BaseName)
		center = ""
		right  = fmt.Sprintf("%s | %s", s.Encoding, GetFormattedTime())
	)

	leftLen := len([]rune(left))
	rightLen := len([]rune(right))
	centerLen := len([]rune(center))

	paddingTotal := width - (leftLen + centerLen + rightLen)
	if paddingTotal < 0 {
		// Fallback: truncate or bail out
		b.Buffer.SetLine(0, []rune(left))
		return
	}

	paddingLeft := paddingTotal / 2
	paddingRight := paddingTotal - paddingLeft

	text := left
	text += strings.Repeat(" ", paddingLeft)
	text += center
	text += strings.Repeat(" ", paddingRight)
	text += right

	if state.Config.Statusbar.UseSeparator {
		if state.Config.Statusbar.Position == StatusbarTop {
			b.Buffer.SetLine(0, []rune(text))
			b.Buffer.SetLine(1, []rune(strings.Repeat(state.Config.Statusbar.Separator, width)))
		} else {
			b.Buffer.SetLine(0, []rune(strings.Repeat(state.Config.Statusbar.Separator, width)))
			b.Buffer.SetLine(1, []rune(text))
		}
	} else {
		b.Buffer.SetLine(0, []rune(text))
	}

}

func UpdateStatusBar() {
	sb := state.StatusBar
	if sb != nil {
		sb.CollectStats()
		sb.SetSizeAndPosition()
		sb.RenderContent()
	}
}
