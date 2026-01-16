package abcde

import (
	"bufio"
	"io"
	"strings"
)

type LogInfo struct {
	Album       string
	TotalTracks int
	Ripping     string
	Ripped      []string
	Finished    bool
}

func ParseLogInfo(r io.Reader) LogInfo {
	info := LogInfo{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Grabbing entire CD - tracks:  ") {
			info.TotalTracks = strings.Count(line[30:], " ") + 1
		} else if strings.HasPrefix(line, "Selected: ") {
			info.Album = line[14 : len(line)-1]
		} else if strings.HasPrefix(line, "Grabbing track ") {
			info.Ripping = line[19 : len(line)-3]
		} else if line == "Done." {
			info.Ripped = append(info.Ripped, info.Ripping)
			info.Ripping = ""
		} else if line == "Finished." {
			info.Finished = true
		}
	}

	return info
}
