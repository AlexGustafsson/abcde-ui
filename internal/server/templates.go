package server

import (
	"embed"
	"io"
	"strings"
	"text/template"

	"github.com/AlexGustafsson/abcde-ui/internal/abcde"
)

//go:embed templates
var templates embed.FS

func Render(w io.Writer, running bool, runnerLogs string, runnerErr error) error {
	t, err := template.ParseFS(templates, "templates/index.html.gotmpl")
	if err != nil {
		return err
	}

	errString := ""
	errLines := 0
	if runnerErr != nil {
		errString = runnerErr.Error()
		errLines = strings.Count(errString, "\n")
	}
	data := &Data{
		IsRipping: running,

		Error:      errString,
		ErrorLines: errLines,

		Logs:     runnerLogs,
		LogLines: strings.Count(runnerLogs, "\n"),

		Info: abcde.ParseLogInfo(strings.NewReader(runnerLogs)),
	}

	return t.Execute(w, data)
}
