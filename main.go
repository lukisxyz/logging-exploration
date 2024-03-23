package main

import (
	"context"
	"io"
	"log/slog"
	"logging-challenge/greetings"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

var log *slog.Logger

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		oscall := <-ch
		slog.Error("system call", oscall)
		cancel()
	}()

	// start: set up any of your logger configuration here if necessary

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	lf, err := os.OpenFile(
		"logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666,
	)
	if err != nil {
		log := slog.New(slog.NewJSONHandler(os.Stdout, opts))
		log.Error("unable to find log file")
	}
	mw := io.MultiWriter(os.Stdout, lf)
	log = slog.New(&SlogContextHandler{slog.NewJSONHandler(
		mw,
		opts,
	)})
	slog.SetDefault(log)

	// end: set up any of your logger configuration here

	r := chi.NewRouter()
	r.Use(logMiddleware)

	r.Mount("/", greetings.Router())

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to listen and serve http server", err)
		}
	}()
	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Error("failed to shutdown http server gracefully", err)
	}
}
