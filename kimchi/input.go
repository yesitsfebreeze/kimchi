package main

import "strings"

type Mods struct {
	Ctrl  bool
	Alt   bool
	Shift bool
	Key   rune
}

func ctrl_mod(char rune) byte {
	return byte(char) & 0x1F
}

func parse_input(s string) byte {
	if strings.HasPrefix(s, "ctrl+") && len(s) == 6 {
		c := s[len(s)-1]
		return c - 'a' + 1
	}
	return 0 // or panic/error
}

func input_loop() {
	// // Draw current_buffer and main_cursor
	// fmt.Printf("\033[H")  // move to top left
	// fmt.Printf("\033[2K") // clear line
	// fmt.Print("> ")
	// for i, ch := range current_buffer {
	// 	if i == main_cursor {
	// 		fmt.Printf("\033[7m%c\033[0m", ch) // inverse main_cursor
	// 	} else {
	// 		fmt.Printf("%c", ch)
	// 	}
	// }
	// if main_cursor == len(current_buffer) {
	// 	fmt.Print("\033[7m \033[0m") // show empty main_cursor at end
	// }
	// fmt.Printf("\033[%d;1H", render_height) // move to bottom line

	// // Read a key
	// b := make([]byte, 3)
	// n, _ := os.Stdin.Read(b)
	// if n == 1 {
	// 	switch b[0] {
	// 	case 3:
	// 		return // Ctrl-C
	// 	case 127:
	// 		// Backspace
	// 		if main_cursor > 0 {
	// 			current_buffer = append(current_buffer[:main_cursor-1], current_buffer[main_cursor:]...)
	// 			main_cursor--
	// 		}
	// 	case 13:
	// 		// Enter
	// 		fmt.Printf("\n[Entered]: %s\n", string(current_buffer))
	// 		current_buffer = []rune{}
	// 		main_cursor = 0
	// 	case 27:
	// 		// ESC — ignore
	// 	default:
	// 		// Insert printable character
	// 		current_buffer = append(current_buffer[:main_cursor], append([]rune{rune(b[0])}, current_buffer[main_cursor:]...)...)
	// 		main_cursor++
	// 	}
	// } else if n == 3 && b[0] == 27 && b[1] == 91 {
	// 	switch b[2] {
	// 	case 67:
	// 		// → Right arrow
	// 		if main_cursor < len(current_buffer) {
	// 			main_cursor++
	// 		}
	// 	case 68:
	// 		// ← Left arrow
	// 		if main_cursor > 0 {
	// 			main_cursor--
	// 		}
	// 	}
	// }
}
