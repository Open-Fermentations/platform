package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/logging"
	"open-fermentations/internal/route"

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

	w.Header().Set("Content-Type", route.ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(jsonResp); err != nil {
		slog.Error("writing response", logging.Err(err))
		return
	}
}

func (s *Server) deleteBatch(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue(IDKey)
	id, err := uuid.Parse(idString)
	if err != nil {
		slog.Error("parsing batch id", slog.String(IDKey, idString), logging.Err(err))
		http.Error(w, FailedToParsePathId, http.StatusBadRequest)
		return
	}

	if err = s.svc.DeleteBatch(id); err != nil {
		slog.Error("deleting batch", logging.Err(err))
		http.Error(w, "Failed to successfully delete batch", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) getBatches(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	limit := getIntQueryParam(r, "limit", 50)
	offset := getIntQueryParam(r, "offset", 0)

	batches, total, err := s.svc.GetBatches(fmt.Sprintf("%%%v%%", name), limit, offset)
	if err != nil {
		slog.Error("getting batches from service", logging.Err(err), slog.Group("query",
			slog.Int("limit", limit),
			slog.Int("offset", offset),
			slog.String("name", name)))
		http.Error(w, "Failed to get batches", http.StatusInternalServerError)
	}

	batchDtos := make([]dto.BatchDTO, len(batches))
	for i, b := range batches {
		batchDtos[i] = *new(dto.BatchDTO).FromModel(b)
	}

	page := dto.PageDTO[dto.BatchDTO]{
		Limit:  limit,
		Offset: offset,
		Total:  total,
		Data:   batchDtos,
	}

	pageJson, err := json.Marshal(page)
	if err != nil {
		slog.Error("marshalling batch page to json", logging.Err(err))
		http.Error(w, FailedToMarshall, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", route.ContentTypeJSON)
	if _, err := w.Write(pageJson); err != nil {
		slog.Error("writing batch page", logging.Err(err))
		return
	}
}

func (s *Server) getBatchById(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue(IDKey)
	id, err := uuid.Parse(rawId)
	if err != nil {
		slog.Error("parsing batch id from path", logging.Err(err), slog.String("id", rawId))
		http.Error(w, FailedToParsePathId, http.StatusBadRequest)
		return
	}

	batch, err := s.svc.GetBatchById(id)
	if err != nil {
		slog.Error("fetching batch by id", logging.Err(err), slog.String("id", id.String()))
		http.Error(w, "failed to fetch batch by id", http.StatusInternalServerError)
		return
	}

	if batch == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	batchDto := new(dto.BatchDTO).FromModel(*batch)

	batchJson, err := json.Marshal(batchDto)
	if err != nil {
		slog.Error("marshalling batch to json", []any{logging.Err(err), batchDto.Slog()}...)
		http.Error(w, FailedToMarshall, http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(batchJson); err != nil {
		slog.Error("failed to write batch", []any{logging.Err(err), batchDto.Slog()}...)
		return
	}
}

func (s *Server) putBatchById(w http.ResponseWriter, r *http.Request) {
	rawBatchId := r.PathValue(IDKey)
	batchId, err := uuid.Parse(rawBatchId)
	if err != nil {
		slog.Error("parsing batch id from path", logging.Err(err), slog.String("id", rawBatchId))
		http.Error(w, FailedToParsePathId, http.StatusBadRequest)
		return
	}

	var d dto.UpdateBatchDTO
	if err := readBody(r, &d); err != nil {
		slog.Error("reading body", logging.Err(err))
		http.Error(w, BadBodyRead, http.StatusBadRequest)
		return
	}

	batch, err := s.svc.UpdateBatch(batchId, d.Name)
	if err != nil {
		slog.Error("updating batch", []any{logging.Err(err), d.Slog()}...)
		http.Error(w, "failed to successfully update batch", http.StatusInternalServerError)
		return
	}

	batchDto := new(dto.BatchDTO).FromModel(*batch)

	jsonResp, err := json.Marshal(batchDto)
	if err != nil {
		slog.Error("marshalling batch dto", []any{logging.Err(err), batchDto.Slog()}...)
		http.Error(w, FailedToMarshall, http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(jsonResp); err != nil {
		slog.Error("writing response", logging.Err(err))
		return
	}
}
