package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type PathKind int

const (
	PathNone PathKind = iota
	PathFile
	PathDir
)

func file_exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func is_dir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func is_file(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

func classify_path(path string) (PathKind, string, error) {
	cs, err := filepath.Abs(path)
	if err != nil {
		return PathNone, "", err
	}

	if !file_exists(cs) {
		return PathNone, cs, nil
	}

	switch {
	case is_dir(cs):
		return PathDir, cs, nil
	case is_file(cs):
		return PathFile, cs, nil
	default:
		return PathNone, cs, fmt.Errorf("unsupported path type: %s", cs)
	}
}
