package main

import (
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"path/filepath"
)

const DAEMON_SOCKET_PATH = "/tmp/kitsuned.sock"

type DaemonCmd int

const (
	DaemonConnect DaemonCmd = iota
	DaemonDisconnect
	DaemonShutdown
	DaemonHighlight
)

var DaemonCmdEnum = NewEnumMap(map[string]DaemonCmd{
	"Connect":    DaemonConnect,
	"Disconnect": DaemonDisconnect,
	"Highlight":  DaemonHighlight,
})

type DaemonData struct {
	User        string `json:"user,omitempty"`
	Lang        string `json:"lang,omitempty"`
	Code        string `json:"code,omitempty"`
	Path        string `json:"path,omitempty"`
	RegionStart int    `json:"region_start,omitempty"`
	RegionEnd   int    `json:"region_end,omitempty"`
}

func DaemonRunning() bool {
	conn, err := net.Dial("unix", DAEMON_SOCKET_PATH)
	if err != nil {
		return false
	}
	conn.Close()

	return true
}

func DaemonStart() error {
	exe_path, err := os.Executable()
	if err != nil {
		return err
	}
	daemon_path := filepath.Join(filepath.Dir(exe_path), "kitsuned")

	LogF("Starting daemon at %s", daemon_path)
	if !FileExists(daemon_path) {
		return os.ErrNotExist
	}

	cmd := exec.Command(daemon_path)

	// TODO impl forward to log
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	LogF("Daemon started with PID %d", cmd.Process.Pid)
	LogF("Daemon is running: %v", DaemonRunning())

	return nil
}

func InitDaemon() {
	if !DaemonRunning() {
		if err := DaemonStart(); err != nil {
			LogF("Failed to start daemon: %v", err)
			os.Exit(1)
		}
	} else {
		Log("Daemon is already running.")
	}
	DaemonSend(DaemonConnect, DaemonData{User: GetCurrentUser()})
}

func DaemonStop() {
	if !DaemonRunning() {
		Log("No daemon is running.")
		return
	}
	DaemonSend(DaemonDisconnect, DaemonData{User: GetCurrentUser()})
}

func DaemonSend(cmd DaemonCmd, data DaemonData) (string, error) {
	conn, err := net.Dial("unix", DAEMON_SOCKET_PATH)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	job := DaemonCmdEnum.String(cmd)
	wrapped := map[string]interface{}{
		job: data,
	}
	send_data, err := json.Marshal(wrapped)

	if err != nil {
		return "", err
	}

	if _, err := conn.Write(send_data); err != nil {
		return "", err
	}

	// write newline because of `BufReader::read_line()`
	if _, err := conn.Write([]byte("\n")); err != nil {
		return "", err
	}

	// read response
	buf := make([]byte, 4096) // increase if you expect large output
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}
