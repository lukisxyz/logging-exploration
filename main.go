package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		oscall := <-ch
		log.Warn().Msgf("system call:%+v", oscall)
		cancel()
	}()

	r := mux.NewRouter()
	r.HandleFunc("/", handler)

	// start: set up any of your logger configuration here if necessary
	// set log level to debug
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// prepare file for save log
	lf, err := os.OpenFile(
		"logs/app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open log file")
	}

	multiwriters := zerolog.MultiLevelWriter(os.Stdout, lf)
	log.Logger = zerolog.New(multiwriters).With().Timestamp().Logger()
	// end: set up any of your logger configuration here

	r.Use(logMiddleware)

	listenerAddr := ":8080"
	server := &http.Server{
		Addr:    listenerAddr,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to listen and serve http server")
		}
	}()
	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server gracefully")
	}
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := log.Logger.With().
			Str("request_id", uuid.New().String()).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("query", r.URL.RawQuery).
			Logger()
		ctx := log.WithContext(r.Context())
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	ctx := r.Context()
	log := log.Ctx(ctx).With().Str("func", "handler").Logger()
	log.Debug().Msg("processing request")
	res, err := greeting(ctx, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(res))
}

func greeting(ctx context.Context, name string) (string, error) {
	log := log.Ctx(ctx).With().Str("func", "greeting").Logger()
	log.Debug().Msg("processing greeting")
	if len(name) < 5 {
		return fmt.Sprintf("Hello %s! Your name is to short\n", name), nil
	}
	return fmt.Sprintf("Hello %s!", name), nil
}
