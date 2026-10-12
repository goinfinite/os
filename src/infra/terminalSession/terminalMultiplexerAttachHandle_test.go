package terminalSessionInfra

import (
	"os/exec"
	"sync"
	"syscall"
	"testing"

	"github.com/creack/pty"
)

func exitedByKillSignal(command *exec.Cmd) bool {
	processState := command.ProcessState
	if processState == nil {
		return false
	}

	waitStatus, isWaitStatus := processState.Sys().(syscall.WaitStatus)
	if !isWaitStatus {
		return false
	}

	return waitStatus.Signaled() && waitStatus.Signal() == syscall.SIGKILL
}

func TestAttachHandleClose(t *testing.T) {
	newAttachHandleForClose := func(t *testing.T, command *exec.Cmd) *terminalMultiplexerAttachHandle {
		ptyFile, ptyErr := pty.Start(command)
		if ptyErr != nil {
			t.Fatalf("StartAttachPtyFailed: %v", ptyErr)
		}

		return &terminalMultiplexerAttachHandle{
			attachCmd: command,
			ptyFile:   ptyFile,
		}
	}

	t.Run("StopsTheProcessOnce", func(t *testing.T) {
		handle := newAttachHandleForClose(t, exec.Command("sleep", "30"))

		firstCloseErr := handle.Close()
		if firstCloseErr != nil {
			t.Fatalf("FirstCloseShouldSucceed: %v", firstCloseErr)
		}

		secondCloseErr := handle.Close()
		if secondCloseErr != nil {
			t.Errorf("SecondCloseShouldStayIdempotent: %v", secondCloseErr)
		}

		if handle.attachCmd.ProcessState == nil {
			t.Fatal("CloseShouldReapTheChildProcess")
		}

		if !exitedByKillSignal(handle.attachCmd) {
			t.Error("CloseShouldStopTheChildWithAKillSignal")
		}
	})

	t.Run("StopsTheProcessOnceUnderConcurrentClosers", func(t *testing.T) {
		handle := newAttachHandleForClose(t, exec.Command("sleep", "30"))

		concurrentClosers := 10
		var closersWaitGroup sync.WaitGroup
		for range concurrentClosers {
			closersWaitGroup.Go(func() {
				if closeErr := handle.Close(); closeErr != nil {
					t.Errorf("ConcurrentCloseShouldSucceed: %v", closeErr)
				}
			})
		}
		closersWaitGroup.Wait()

		if handle.attachCmd.ProcessState == nil {
			t.Fatal("ConcurrentCloseShouldReapTheChildProcess")
		}

		if !exitedByKillSignal(handle.attachCmd) {
			t.Error("ConcurrentCloseShouldStopTheChildWithAKillSignal")
		}
	})

	t.Run("DoesNotReportAKillSignalForAnIndependentExit", func(t *testing.T) {
		handle := newAttachHandleForClose(t, exec.Command("true"))

		waitErr := handle.Wait()
		if waitErr != nil {
			t.Fatalf("WaitShouldSucceed: %v", waitErr)
		}

		closeErr := handle.Close()
		if closeErr != nil {
			t.Fatalf("CloseShouldSucceed: %v", closeErr)
		}

		if exitedByKillSignal(handle.attachCmd) {
			t.Error("OwnExitShouldNotCarryAKillSignal")
		}
	})
}
