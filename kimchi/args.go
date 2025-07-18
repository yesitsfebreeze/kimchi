package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	config_path  string
	show_version bool
	dump_config  bool
)

func arg_is_flag(s string) bool {
	return strings.HasPrefix(s, "-")
}

func init_args() {
	flag.StringVar(&config_path, "config", "", "Optional path to .kimchi.lua")
	flag.StringVar(&config_path, "c", "", "Optional path to .kimchi.lua")

	flag.BoolVar(&show_version, "version", false, "Print version and exit")
	flag.BoolVar(&show_version, "v", false, "Print version and exit (shorthand)")

	flag.BoolVar(&dump_config, "dump-config", false, "Print default config in Lua format")
	flag.BoolVar(&dump_config, "dc", false, "Print default config in Lua format")
}

func parse_arg_locations(loc string) {
	pathtype, path, err := classify_path(loc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	if pathtype == PathKind(PathFile) {
		cfg.startfile = path
		cfg.workspace = filepath.Dir(path)
	}

	if pathtype == PathKind(PathDir) {
		cfg.workspace = path
		cfg.startfile = ""
	}

	fmt.Println(cfg.project_scan_limit)
	// find_closest_kimchi_lua()
}

func parse_args() {
	init_args()

	flag.Parse()
	args := flag.Args()

	// get the first non flag as workspace/file path
	for _, arg := range args {
		if !arg_is_flag(arg) {
			parse_arg_locations(arg)

			fmt.Println(cfg.startfile)
			fmt.Println(cfg.workspace)
			os.Exit(1)
			break
		}
	}

	if len(args) > 0 && args[0] == "dump-config" || dump_config {
		fmt.Printf("dump config")
		os.Exit(1)
	}

	if len(args) > 0 && args[0] == "version" || show_version {
		fmt.Printf("kim v%s\n", VERSION)
		os.Exit(1)
	}

	// if len(args) == 0 {
	// 	fmt.Fprintln(os.Stderr, "Usage: kim [file or folder] [--config=...]")
	// 	os.Exit(1)
	// }

}
