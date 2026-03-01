package abcde

import "os/exec"

func Command() *exec.Cmd {
	return exec.Command("abcde", "-N", "-o", "flac", "-p")
}
