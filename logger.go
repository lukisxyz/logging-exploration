package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey string

const (
	slogFields ctxKey = "slog_fields"
)

type SlogContextHandler struct {
	slog.Handler
}

func (s SlogContextHandler) Handle(c context.Context, r slog.Record) error {
	if attrs, ok := c.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}
	return s.Handler.Handle(c, r)
}

func AppendCtx(c context.Context, attrs []slog.Attr) context.Context {
	if c == nil {
		c = context.Background()
	}

	if v, ok := c.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attrs...)
		return context.WithValue(c, slogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attrs...)
	return context.WithValue(c, slogFields, v)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := AppendCtx(
			r.Context(),
			[]slog.Attr{
				slog.String("request_id", uuid.New().String()),
				slog.Group(
					"request",
					slog.String("url", r.URL.String()),
					slog.String("method", r.Method),
				),
			},
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
