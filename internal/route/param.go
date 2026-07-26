package route

import (
	"log/slog"
	"net/http"
	"strconv"
)

func GetIntQueryParam(r *http.Request, key string, def int) int {
	rawValue := r.URL.Query().Get(key)
	val, err := strconv.Atoi(rawValue)
	if err != nil {
		slog.Debug("defaulting quary parameter", slog.String("key", key), slog.String("value", rawValue), slog.Int("default", def))
		return def
	}

	return val
}

func GetStringQueryParam(r *http.Request, key, def string) string {
	rawValue := r.URL.Query().Get(key)
	if rawValue == "" {
		slog.Debug("defaulting query parameter", slog.String("key", key), slog.String("value", rawValue), slog.String("default", def))
		return def
	}
	return rawValue
}

func GetBoolQueryParam(r *http.Request, key string, def bool) bool {
	rawValue := r.URL.Query().Get(key)
	b, err := strconv.ParseBool(rawValue)
	if err != nil {
		slog.Debug("defaulting query parameter", slog.String("key", key), slog.String("value", rawValue), slog.Bool("default", def))
		return def
	}

	return b
}
