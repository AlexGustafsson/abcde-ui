package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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

	server := &http.Server{
		Addr:    ":8082",
		Handler: server.NewServer(runner),
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
