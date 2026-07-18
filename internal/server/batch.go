package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/logging"

	"github.com/google/uuid"
)

func (s *Server) postBatch(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read request body in postBatch", logging.Err(err))
		http.Error(w, BadBodyRead, http.StatusBadRequest)
		return
	}

	var b []dto.CreateBatchDTO
	if err := json.Unmarshal(rawBody, &b); err != nil {
		slog.Error("unmarshalling create batch dto", logging.Err(err))
		http.Error(w, FailedToMarshall, http.StatusBadRequest)
	}

	u := r.Context().Value(ContextUserIdKey).(string)
	userId, err := uuid.Parse(u)
	if err != nil || u == "" {
		slog.Error("parsing user id", logging.Err(err), slog.String("user_id", u))
		http.Error(w, FailedToParseUserId, http.StatusBadRequest)
		return
	}

	ms, err := s.svc.CreateBatch(userId, b)
	if err != nil {
		slog.Error("creating batch", logging.Err(err))
		http.Error(w, "error creating batch", http.StatusInternalServerError)
		return
	}

	dtos := make([]dto.BatchDTO, len(ms))
	for i, m := range ms {
		dtos[i] = *new(dto.BatchDTO).FromModel(m)
	}

	jsonResp, err := json.Marshal(dtos)
	if err != nil {
		slog.Error("marshalling batch dtos", logging.Err(err))
		http.Error(w, "Faled to marshal batch dto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(jsonResp); err != nil {
		slog.Error("writing response", logging.Err(err))
		return
	}
}
