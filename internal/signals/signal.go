// +build !windows,!darwin

// Package signals provides utilities for signal handling and process signaling.
package signals

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

var signalMap = map[string]os.Signal{"SIGABRT": syscall.SIGABRT,
	"SIGALRM":   syscall.SIGALRM,
	"SIGBUS":    syscall.SIGBUS,
	"SIGCHLD":   syscall.SIGCHLD,
	"SIGCLD":    syscall.SIGCLD,
	"SIGCONT":   syscall.SIGCONT,
	"SIGFPE":    syscall.SIGFPE,
	"SIGHUP":    syscall.SIGHUP,
	"SIGILL":    syscall.SIGILL,
	"SIGINT":    syscall.SIGINT,
	"SIGIO":     syscall.SIGIO,
	"SIGIOT":    syscall.SIGIOT,
	"SIGKILL":   syscall.SIGKILL,
	"SIGPIPE":   syscall.SIGPIPE,
	"SIGPOLL":   syscall.SIGPOLL,
	"SIGPROF":   syscall.SIGPROF,
	"SIGPWR":    syscall.SIGPWR,
	"SIGQUIT":   syscall.SIGQUIT,
	"SIGSEGV":   syscall.SIGSEGV,
	"SIGSTKFLT": syscall.SIGSTKFLT,
	"SIGSTOP":   syscall.SIGSTOP,
	"SIGSYS":    syscall.SIGSYS,
	"SIGTERM":   syscall.SIGTERM,
	"SIGTRAP":   syscall.SIGTRAP,
	"SIGTSTP":   syscall.SIGTSTP,
	"SIGTTIN":   syscall.SIGTTIN,
	"SIGTTOU":   syscall.SIGTTOU,
	"SIGUNUSED": syscall.SIGUNUSED,
	"SIGURG":    syscall.SIGURG,
	"SIGUSR1":   syscall.SIGUSR1,
	"SIGUSR2":   syscall.SIGUSR2,
	"SIGVTALRM": syscall.SIGVTALRM,
	"SIGWINCH":  syscall.SIGWINCH,
	"SIGXCPU":   syscall.SIGXCPU,
	"SIGXFSZ":   syscall.SIGXFSZ}

// ToSignal returns OS dependent signal name for given signal name (or syscall.SIGTERM if garbage given).
func ToSignal(signalName string) (os.Signal, error) {
	if !strings.HasPrefix(signalName, "SIG") {
		signalName = "SIG" + signalName
	}
	if sig, ok := signalMap[signalName]; ok {
		return sig, nil
	}
	return syscall.SIGTERM, nil
}

// Kill sends signal to the process.
//
// Args:.
//    process - the process which the signal should be sent to
//    sig - the signal will be sent
//    sigChildren - true if the signal needs to be sent to the children also
//
func Kill(process *os.Process, sig os.Signal, sigChildren bool) error {
	localSig, ok := sig.(syscall.Signal)
	if !ok {
		return fmt.Errorf("signal type assertion failed: expected syscall.Signal, got %T", sig)
	}
	pid := process.Pid
	if sigChildren {
		pid = -pid
	}
	if err := syscall.Kill(pid, localSig); err != nil {
		return fmt.Errorf("failed to send signal to process %d: %w", pid, err)
	}
	return nil
}
