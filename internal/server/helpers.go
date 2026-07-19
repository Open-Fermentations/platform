package server

import (
	"log/slog"
	"net/http"
	"strconv"
)

func getIntQueryParam(r *http.Request, key string, def int) int {
	rawValue := r.URL.Query().Get(key)
	val, err := strconv.Atoi(rawValue)
	if err != nil {
		slog.Debug("defaulting quary parameter", slog.String("key", key), slog.String("value", rawValue))
		return def
	}

	return val
}
