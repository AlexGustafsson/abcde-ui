package abcde

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	ErrAlreadyRunning = errors.New("already running")
)

type Runner struct {
	// Dir specifies the working directory of the command.
	Dir string

	mutex  sync.Mutex
	cmd    *exec.Cmd
	output strings.Builder
	err    error
}

func (r *Runner) Start() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.cmd != nil {
		return ErrAlreadyRunning
	}

	r.output.Reset()

	cmd := Command()
	cmd.Dir = r.Dir
	cmd.Stdout = &r.output
	cmd.Stderr = &r.output

	slog.Info("Starting abcde")
	r.err = nil
	err := cmd.Start()
	if err != nil {
		r.err = err
		return err
	}

	go func() {
		err := cmd.Wait()
		if err == nil {
			slog.Info("Ran abcde successfully")
		} else {
			slog.Error("Failed to run abcde", slog.Any("error", err))
			// Fallthrough
		}

		r.mutex.Lock()
		defer r.mutex.Unlock()

		r.cmd = nil
		r.err = err
	}()

	r.cmd = cmd
	return nil
}

func (r *Runner) Running() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.cmd != nil
}

func (r *Runner) Error() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.err
}

func (r *Runner) Output() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.output.String()
}

func (r *Runner) Shutdown() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.cmd == nil {
		return nil
	}

	if err := r.cmd.Process.Signal(os.Interrupt); err != nil {
		return err
	}

	return r.cmd.Wait()
}

func (r *Runner) Kill() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.cmd == nil {
		return nil
	}

	if err := r.cmd.Process.Kill(); err != nil {
		return err
	}

	return r.cmd.Wait()
}
