package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

type Cell struct {
	Position Vec2
	Rune     rune
	Style    Style
}

type DirtyCoord = [2]int

type Screen struct {
	Obj          tcell.Screen
	Size         Vec2
	OldXTerm     *term.State
	OldTerm      string
	ModifiedTerm bool
	Cells        [][]Cell
	Dirty        map[DirtyCoord]bool // List of dirty cells to flush
}

func ActivateXTerm() {
	if oldState, err := term.GetState(int(os.Stdin.Fd())); err == nil {
		state.Screen.OldXTerm = oldState
	}

	state.Screen.OldTerm = os.Getenv("TERM")
	os.Setenv("TERM", "xterm-256color")
	state.Screen.ModifiedTerm = true
}

func CreateScreen() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	if err := screen.Init(); err != nil {
		return err
	}

	state.Screen.Obj = screen

	return nil
}

func InitScreen() error {
	state.Screen = &Screen{
		OldTerm:      os.Getenv("TERM"),
		ModifiedTerm: false,
	}

	if state.Config.XTerm {
		ActivateXTerm()
	}

	err := CreateScreen()
	if err != nil {
		return err
	}

	if state.Config.CopyPaste {
		state.Screen.Obj.EnablePaste()
	}

	// if state.Screen.ModifiedTerm {
	// 	os.Setenv("TERM", state.Screen.OldTerm)
	// }

	if state.Config.MouseEnabled {
		state.Screen.Obj.EnableMouse()
	}

	state.Screen.UpdateSize()
	state.Screen.Dirty = make(map[DirtyCoord]bool)

	return nil
}

func (s *Screen) UpdateSize() {
	state.Screen.Obj.Sync()
	w, h := state.Screen.Obj.Size()

	s.Size = Vec2{X: w, Y: h}

	s.Cells = make([][]Cell, h)
	for y := range h {
		s.Cells[y] = make([]Cell, w)
		for x := range w {
			s.Cells[y][x] = Cell{
				Rune:     ' ',
				Position: Vec2{X: x, Y: y},
				Style:    state.Theme.Text,
			}
		}
	}

	// for _, area := range state.Areas {
	// 	// area.UpdateViewSize()
	// }

}

func (c *Cell) Clear() {
	if c.Rune == ' ' && c.Style == state.Theme.Text {
		return // no need to clear empty cells
	}
	c.Rune = ' '
	c.Style = state.Theme.Text
	c.Dirty()
}

func (c *Cell) Dirty() {
	coord := DirtyCoord{c.Position.X, c.Position.Y}
	state.Screen.Dirty[coord] = true
}

func (s *Screen) Clear() {
	for y := range s.Cells {
		for x := range s.Cells[y] {
			s.Cells[y][x].Clear()
		}
	}
}

func (s *Screen) ClearRegion(x, y, w, h int) {
	for cy := y; cy < y+h; cy++ {
		for cx := x; cx < x+w; cx++ {
			s.Cells[y][x].Clear()
		}
	}
}

func (s *Screen) ClearLine(y int) {
	if y < 0 || y >= len(s.Cells) {
		return // silently ignore out-of-bounds
	}

	for x := range s.Cells[y] {
		s.Cells[y][x].Clear()
	}
}

func (s *Screen) CellXY(x, y int, r rune, style Style) {
	if y < 0 || y >= state.Screen.Size.Y || x < 0 || x >= state.Screen.Size.X {
		return // silently ignore out-of-bounds
	}

	cell := &s.Cells[y][x]

	if cell.Rune == r && cell.Style == style {
		return
	}

	cell.Rune = r
	cell.Style = style
	cell.Position = Vec2{X: x, Y: y}
	cell.Rune = r
	cell.Style = style
	cell.Dirty()
}

func (s *Screen) Cell(coord Vec2, r rune, style Style) {
	s.CellXY(coord.X, coord.Y, r, style)
}

func (c *Cell) SetStyle(style Style) {
	c.Style = style
}

func (c *Cell) Flush() {
	state.Screen.Obj.SetContent(
		c.Position.X,
		c.Position.Y,
		c.Rune,
		nil,
		c.Style.TCellStyle(),
	)
}

func (s *Screen) Poll() tcell.Event {
	if s.Obj == nil {
		LogErr("Screen object is nil, cannot poll")
		return nil
	}
	return s.Obj.PollEvent()
}

func (s *Screen) Show() {
	if s.Obj == nil {
		LogErr("Screen object is nil, cannot show")
		return
	}
	s.Obj.Show()
	for k := range s.Dirty {
		delete(s.Dirty, k)
	}
}

func (s *Screen) Flush() {
	// Log(len(s.Dirty), "dirty cells of", s.Size.X*s.Size.Y)

	for coord := range s.Dirty {
		cell := &s.Cells[coord[1]][coord[0]]
		cell.Flush()
	}

	s.Show()
}

func (s *Screen) FlushAll() {
	// dirtyCount := len(s.Dirty)

	// if dirtyCount != 0 {
	// 	// Log(len(s.Dirty), "dirty cells of", s.Size.X*s.Size.Y)
	// }

	for y := range s.Size.Y {
		for x := range s.Size.X {
			cell := &s.Cells[y][x]
			cell.Flush()
		}
	}

	s.Show()
}

func ResetScreen() {
	if state.Screen == nil {
		return // nothing to reset
	}

	if state.Screen.Obj == nil {
		return // nothing to reset
	}

	state.Screen.Reset()
}

func (s *Screen) Reset() {
	if s.Obj == nil {
		return // nothing to reset
	}
	term.Restore(int(os.Stdin.Fd()), s.OldXTerm)
	fmt.Fprint(os.Stderr, "\033[?25h\033[0m\n")
	os.Setenv("TERM", s.OldTerm)

	s.Obj.DisableMouse()
	s.Obj.Fini()
}
