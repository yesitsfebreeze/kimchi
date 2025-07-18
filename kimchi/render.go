package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var render_width int
var render_height int

var old_state *term.State

func enter_raw_mode() {
	var err error
	old_state, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
}

func exit_raw_mode() {
	term.Restore(int(os.Stdin.Fd()), old_state)
	fmt.Print("\033[?25h") // make sure to show cursor again
}

func start_render() {
	enter_raw_mode()
	clear_screen()
}

func render() {
	clear_screen()
	// render_buffer() // your editor content
	render_log() // log overlay at bottom
	// move_cursor(cursor_x, cursor_y)
}

func clear_screen() {
	render_width, render_height, _ = term.GetSize(int(os.Stdout.Fd()))
	fmt.Print("\033[2J\033[H")
	fmt.Print("\033[?25l")
	fmt.Printf("\033[%d;%dH", 0+1, 0+1)
}
