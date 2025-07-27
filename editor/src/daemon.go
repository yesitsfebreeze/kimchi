package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

const DAEMON_NAME = "kitsuned"
const DAEMON_MAX_RETRIES = 32
const DAEMON_RETRY_INTERVAL = 25 * time.Millisecond

type DaemonCmd int

const (
	DaemonConnect DaemonCmd = iota
	DaemonDisconnect
	DaemonInstallLanguage
	DaemonHighlight
)

var DaemonCmdEnum = NewEnumMap(map[string]DaemonCmd{
	"Connect":         DaemonConnect,
	"Disconnect":      DaemonDisconnect,
	"InstallLanguage": DaemonInstallLanguage,
	"Highlight":       DaemonHighlight,
})

type DaemonData struct {
	User        string `json:"user,omitempty"`
	Lang        string `json:"lang,omitempty"`
	Code        string `json:"code,omitempty"`
	Path        string `json:"path,omitempty"`
	RegionStart int    `json:"region_start,omitempty"`
	RegionEnd   int    `json:"region_end,omitempty"`
}

var DaemonQueue []DaemonMessage
var DaemonQueueLock = sync.Mutex{}

type DaemonMessage struct {
	Cmd      DaemonCmd
	Data     DaemonData
	Callback func(string, error)
}

func EnsureDaemonConnection() {
	var path string = "/tmp/" + DAEMON_NAME + ".sock"

	connect := func() bool {
		if _, err := os.Stat(path); err != nil {
			return false
		}
		conn, err := net.Dial("unix", path)
		if err != nil {
			LogErr("Could not connect to daemon:", err)
			return false
		}
		state.Daemon = &conn
		Log("Connected to daemon.")
		FlushDaemonQueue()
		return true
	}

	connected := connect()

	go func() {
		for range DAEMON_MAX_RETRIES {
			if connected {
				return
			}
			connected = connect()
			if connected {
				return
			}
			time.Sleep(DAEMON_RETRY_INTERVAL)
		}
		if !connected {
			LogErr("Failed to connect to daemon after multiple attempts.")
			os.Exit(1)
		}
	}()
}

func GetDaemonBinPath() (string, error) {
	exe_path, err := os.Executable()
	if err != nil {
		return "", err
	}
	daemon_bin_path := filepath.Join(filepath.Dir(exe_path), DAEMON_NAME)
	if !FileExists(daemon_bin_path) {
		return "", fmt.Errorf("daemon binary not found at %s", daemon_bin_path)
	}
	return daemon_bin_path, nil
}

func DaemonStart() error {
	daemon_bin, err := GetDaemonBinPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(daemon_bin)

	// TODO impl forward to log
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	LogF("Daemon started with PID %d", cmd.Process.Pid)
	EnsureDaemonConnection()

	return nil
}

func ConnectToDaemon() {
	DaemonSend(DaemonConnect, DaemonData{User: GetCurrentUser()}, func(res string, err error) {
		if err != nil {
			LogErr("Failed to connect to daemon:", err)
			return
		}
		Log("Connected to daemon.")

		for _, lang := range state.Config.DefaultLanguages {
			DaemonSend(DaemonInstallLanguage, DaemonData{Lang: lang}, func(res string, err error) {
				if err != nil {
					LogErr("Failed to install language:", err)
					return
				}
				Log("Installed language:", lang)
			})
		}
	})
}

func InitDaemon() {
	EnsureDaemonConnection()
	if state.Daemon != nil {
		Log("Daemon is already running.")
		ConnectToDaemon()
		return
	}
	DaemonStart()
	ConnectToDaemon()
}

func DisconnectFromDaemon() {
	DaemonSend(DaemonDisconnect, DaemonData{User: GetCurrentUser()}, func(res string, err error) {
		if err != nil {
			LogErr("Failed to disconnect from daemon:", err)
			return
		}
		Log("Disconnected from daemon.")
	})
}

func DaemonSend(cmd DaemonCmd, data DaemonData, f func(string, error)) {
	if state.Daemon == nil {
		DaemonQueueLock.Lock()
		defer DaemonQueueLock.Unlock()
		DaemonQueue = append(DaemonQueue, DaemonMessage{Cmd: cmd, Data: data, Callback: f})
		Log("Daemon not ready. Queued command:", DaemonCmdEnum.String(cmd))
		return
	}

	job := DaemonCmdEnum.String(cmd)
	wrapped := map[string]any{
		job: data,
	}
	send_data, err := json.Marshal(wrapped)

	if err != nil {
		f("", err)
		return
	}

	conn := *state.Daemon
	if _, err := conn.Write(send_data); err != nil {
		f("", err)
		return
	}

	// write newline because of `BufReader::read_line()`
	if _, err := conn.Write([]byte("\n")); err != nil {
		f("", err)
		return
	}

	if res, err := DaemonReadMessage(); err == nil {
		f(res, nil)
	}
}

func DaemonReadMessage() (string, error) {
	conn := *state.Daemon
	// Read the 4-byte big-endian length prefix
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return "", fmt.Errorf("failed to read length: %w", err)
	}

	length := binary.BigEndian.Uint32(lenBuf)
	if length > 10*1024*1024 {
		return "", fmt.Errorf("message too large: %d bytes", length) // max 10MB
	}

	// Read the message body
	msgBuf := make([]byte, length)
	if _, err := io.ReadFull(conn, msgBuf); err != nil {
		return "", fmt.Errorf("failed to read message body: %w", err)
	}

	return string(msgBuf), nil
}

func FlushDaemonQueue() {
	DaemonQueueLock.Lock()
	defer DaemonQueueLock.Unlock()

	if state.Daemon == nil {
		LogErr("Cannot flush daemon queue — still not connected.")
		return
	}

	for _, msg := range DaemonQueue {
		LogF("Flushing queued command: %s", DaemonCmdEnum.String(msg.Cmd))
		DaemonSend(msg.Cmd, msg.Data, msg.Callback)
	}

	DaemonQueue = nil // clear
}
