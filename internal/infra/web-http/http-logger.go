package webhttp

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

type HttpLogger struct {
	logger *slog.Logger
}

func NewHttpLogger() HttpLogger {
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)

	return HttpLogger{
		logger: slog.New(jsonHandler),
	}
}

func (h *HttpLogger) Log(c echo.Context, err error) {
	var body map[string]int
	_ = c.Bind(&body)
	h.logger.LogAttrs(context.Background(), slog.LevelError, "Unexpected Error",
		slog.String("request_method", c.Request().Method),
		slog.String("request_endpoint", c.Request().RequestURI),
		slog.Any("request_header", c.Request().Header),
		slog.Any("request_body", body),
		slog.String("error_message", err.Error()),
	)
}
