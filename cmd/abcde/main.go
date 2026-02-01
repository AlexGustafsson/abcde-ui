package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/AlexGustafsson/abcde-ui/internal/abcde"
	"github.com/AlexGustafsson/abcde-ui/internal/server"
)

func main() {
	runner := &abcde.Runner{
		Dir:               ".",
		GrapevineEndpoint: os.Getenv("ABCDE_UI_GRAPEVINE_ENDPOINT"),
		GrapevineTopic:    os.Getenv("ABCDE_UI_GRAPEVINE_TOPIC"),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /livez", healthcheck)
	mux.HandleFunc("GET /readyz", healthcheck)

	mux.Handle("/", server.NewServer(runner))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		signals := make(chan os.Signal, 2)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

		caught := 0
		for range signals {
			caught++
			if caught == 1 {
				slog.Info("Caught signal, exiting gracefully")
				runner.Shutdown()
				server.Shutdown(context.Background())
			} else {
				slog.Info("Caught signal, exiting now")
				os.Exit(1)
			}
		}
	}()

	err := server.ListenAndServe()
	if err != http.ErrServerClosed && err != nil {
		slog.Error("Failed to serve", slog.Any("error", err))
		os.Exit(1)
	}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	errors := make([]string, 0)

	{
		sr, _ := filepath.Glob("/dev/sr*")
		cdrom, _ := filepath.Glob("/dev/cdrom")
		if len(sr)+len(cdrom) == 0 {
			errors = append(errors, "No disk reader available")
		}
	}

	{
		_, err := exec.LookPath("abcde")
		if err != nil {
			errors = append(errors, "Can't find abcde")
		}
	}

	if len(errors) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(strings.Join(errors, "\n")))
	}
}
