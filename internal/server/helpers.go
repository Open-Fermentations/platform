package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
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

func getUserId(r *http.Request) uuid.UUID {
	id := r.Context().Value(ContextUserIdKey).(uuid.UUID)

	return id
}

func readBody(r *http.Request, d interface{}) error {
	rawBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rawBody, d); err != nil {
		return err
	}
	return nil
}
