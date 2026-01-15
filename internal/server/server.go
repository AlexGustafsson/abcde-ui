package server

import (
	"log/slog"
	"net/http"
	"net/url"
	"text/template"

	"github.com/AlexGustafsson/abcde-ui/internal/abcde"
)

type Server struct {
	mux *http.ServeMux
}

type Data struct {
	IsRipping bool
	Error     string
	Logs      string
}

func NewServer(runner *abcde.Runner) *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}

	s.mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFS(templates, "templates/index.html.gotmpl")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		errString := ""
		if err := runner.Error(); err != nil {
			errString = err.Error()
		}
		data := &Data{
			IsRipping: runner.Running(),
			Error:     errString,
			Logs:      runner.Output(),
		}

		w.Header().Set("Content-Type", "text/html")
		if err := t.Execute(w, data); err != nil {
			slog.Error("Failed to render template", slog.Any("error", err))
			// Fallthrough
		}
	})

	s.mux.HandleFunc("POST /api/v1/rip", func(w http.ResponseWriter, r *http.Request) {
		err := runner.Start()

		query := make(url.Values)
		switch err {
		case abcde.ErrAlreadyRunning:
			query.Set("error", "running")
		case nil:
			// Do nothing
		default:
			slog.Error("Failed to start abcde", slog.Any("error", err))
			query.Set("error", "internal")
		}

		w.Header().Set("Location", "/#"+query.Encode())
		w.WriteHeader(http.StatusFound)
	})

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
