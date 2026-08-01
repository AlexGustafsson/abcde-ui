package abcde

import "os/exec"

func Command(device string) *exec.Cmd {
	return exec.Command("abcde", "-N", "-o", "flac", "-p", "-d", device)
}
