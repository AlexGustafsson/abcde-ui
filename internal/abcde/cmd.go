package abcde

import "os/exec"

func Command() *exec.Cmd {
	return exec.Command("abcde", "-n", "-N", "-o", "flac", "-p")
}
