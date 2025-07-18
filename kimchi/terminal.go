package main

import "fmt"

func clear_line() {
	fmt.Print("\033[2K")
}
