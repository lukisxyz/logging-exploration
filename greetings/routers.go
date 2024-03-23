package greetings

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	r := chi.NewMux()

	r.Get("/", handler)

	return r
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.URL.Query().Get("name")
	slog.InfoContext(
		ctx,
		"processing request",
		slog.String("name", name),
	)
	res, err := greeting(ctx, name)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed request",
			slog.Int("http_status", http.StatusInternalServerError),
			slog.String("error_message", err.Error()),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(res))
}

func greeting(ctx context.Context, name string) (string, error) {
	slog.InfoContext(
		ctx,
		"processing greeting",
	)
	if len(name) < 5 {
		return fmt.Sprintf("Hello %s! Your name is too short\n", name), nil
	}
	if len(name) > 255 {
		msg := fmt.Sprintf("Hello %s! Your name is more than 255 character\n", name)
		return msg, errors.New(msg)
	}
	return fmt.Sprintf("Hi %s", name), nil
}
