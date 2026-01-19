package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/AlexGustafsson/abcde-ui/internal/abcde"
)

type Server struct {
	mux *http.ServeMux
}

type Data struct {
	IsRipping bool

	Error      string
	ErrorLines int

	Logs     string
	LogLines int

	Info abcde.LogInfo
}

func NewServer(runner *abcde.Runner) *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}

	s.mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		err := Render(w, runner.Running(), runner.Output(), runner.Error())
		if err != nil {
			slog.Error("Failed to render template", slog.Any("error", err))
			w.Write([]byte("<html><body>Failed to render template</body></html>"))
			return
		}
	})

	s.mux.HandleFunc("POST /api/v1/rip", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		fallback := r.FormValue("fallback")
		fmt.Println(fallback)

		err := runner.Start(fallback)

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
