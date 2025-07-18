package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var cfg = create_default_config()
var should_quit = false

func bootstrap() {
	load_cfg()
	fmt.Println(cfg.project_scan_limit)
	fmt.Println(cfg.indent.visual.style)

	parse_args()
	start_render()
}

func main() {
	defer dump_all_logs()
	defer exit_raw_mode()

	bootstrap()

	{ // enter raw mode
		old_state, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		defer term.Restore(int(os.Stdin.Fd()), old_state)
	}

	for !should_quit {
		buf := make([]byte, 1)
		_, _ = os.Stdin.Read(buf)
		logf("pressed key: %q (%d)", buf[0], buf[0])

		if buf[0] == 'q' {
			should_quit = true
		}
		render()
	}
}
