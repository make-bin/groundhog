// @AI_GENERATED
package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// Daemon manages a background process via a PID file.
type Daemon struct {
	pidFile string
	logFile string
}

// New creates a Daemon with the given PID file and log file paths.
func New(pidFile, logFile string) *Daemon {
	return &Daemon{
		pidFile: pidFile,
		logFile: logFile,
	}
}

// Start forks the current executable as a background process.
// args are the command-line arguments to pass to the child process
// (should NOT include the --daemon flag to avoid infinite forking).
func (d *Daemon) Start(executable string, args []string) error {
	// Open (or create) the log file for stdout/stderr redirection.
	logF, err := os.OpenFile(d.logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("daemon: open log file %s: %w", d.logFile, err)
	}

	cmd := exec.Command(executable, args...)
	cmd.Stdout = logF
	cmd.Stderr = logF
	cmd.Stdin = nil

	// Detach from the current process group.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		logF.Close()
		return fmt.Errorf("daemon: start process: %w", err)
	}

	// Write PID file.
	if err := os.WriteFile(d.pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644); err != nil {
		// Best-effort: process is running but we couldn't write the PID file.
		return fmt.Errorf("daemon: write pid file %s: %w", d.pidFile, err)
	}

	// Detach from the child — we don't wait for it.
	if err := cmd.Process.Release(); err != nil {
		return fmt.Errorf("daemon: release process: %w", err)
	}

	fmt.Printf("daemon started (pid %d), logging to %s\n", cmd.Process.Pid, d.logFile)
	return nil
}

// Stop reads the PID file and sends SIGTERM to the process.
func (d *Daemon) Stop() error {
	pid, err := d.readPID()
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("daemon: find process %d: %w", pid, err)
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("daemon: send SIGTERM to %d: %w", pid, err)
	}

	// Remove PID file after signalling.
	_ = os.Remove(d.pidFile)

	fmt.Printf("daemon stopped (pid %d)\n", pid)
	return nil
}

// Status returns (running bool, pid int, err).
func (d *Daemon) Status() (bool, int, error) {
	pid, err := d.readPID()
	if err != nil {
		return false, 0, nil // no PID file → not running
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return false, pid, nil
	}

	// Signal 0 checks if the process exists without sending a real signal.
	if err := proc.Signal(syscall.Signal(0)); err != nil {
		return false, pid, nil
	}

	return true, pid, nil
}

// readPID reads and parses the PID file.
func (d *Daemon) readPID() (int, error) {
	data, err := os.ReadFile(d.pidFile)
	if err != nil {
		return 0, fmt.Errorf("daemon: read pid file %s: %w", d.pidFile, err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, fmt.Errorf("daemon: invalid pid in %s: %w", d.pidFile, err)
	}

	return pid, nil
}

// @AI_GENERATED: end
