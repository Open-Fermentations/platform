package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/logging"
)

func (s *Server) postDevices(w http.ResponseWriter, r *http.Request) {
	var b []dto.CreateDeviceDTO
	if err := readBody(r, &b); err != nil {
		slog.Error("reading create device dto body", logging.Err(err))
		http.Error(w, FailedReadingBody, http.StatusBadRequest)
		return
	}

	userId := getUserId(r)

	ms, err := s.svc.CreateDevices(userId, b)
	if err != nil {
		slog.Error("creating devices", logging.Err(err))
		http.Error(w, "error creating devices", http.StatusInternalServerError)
		return
	}

	dtos := make([]dto.DeviceDTO, len(ms))
	for i, m := range ms {
		dtos[i] = *new(dto.DeviceDTO).FromModel(&m)
	}

	jsonResp, err := json.Marshal(dtos)
	if err != nil {
		slog.Error("marshalling created device dtos", logging.Err(err))
		http.Error(w, "Failed to marshal created device dtos", http.StatusInternalServerError)
		return
	}

	setContentTypeJson(w)
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(jsonResp); err != nil {
		slog.Error("writing response", logging.Err(err))
		return
	}
}
