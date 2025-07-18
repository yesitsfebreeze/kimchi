package main

import "fmt"

const MAIN_CURSOR int = 0

type Cursor struct {
	x int
	y int
}

type Cursors struct {
	list []Cursor
}

func move_cursor(x, y int) {
	fmt.Printf("\033[%d;%dH", y+1, x+1) // ANSI is row;col (1-based)
}

// func set_cursor_pos(x,y int, index int) {
// 	fmt.Printf("\033[%d;%dH", x, y)
// }
