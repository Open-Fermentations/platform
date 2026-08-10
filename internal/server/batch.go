package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/logging"
	"open-fermentations/internal/route"
	"open-fermentations/internal/service"

	"github.com/google/uuid"
)

func (s *Server) postBatch(w http.ResponseWriter, r *http.Request) {
	var b []dto.CreateBatchDTO
	if err := readBody(r, &b); err != nil {
		slog.Error("reading create batch dto body", logging.Err(err))
		http.Error(w, FailedReadingBody, http.StatusBadRequest)
		return
	}

	ms, err := s.svc.CreateBatches(r.Context(), b)
	if err != nil {
		slog.Error("creating batch", logging.Err(err))
		http.Error(w, "error creating batch", http.StatusInternalServerError)
		return
	}

	dtos := make([]dto.BatchDTO, len(ms))
	for i, m := range ms {
		dtos[i] = *new(dto.BatchDTO).FromModel(&m)
	}

	jsonResp, err := json.Marshal(dtos)
	if err != nil {
		slog.Error("marshalling created batch dtos", logging.Err(err))
		http.Error(w, "Faled to marshal created batch dtos", http.StatusInternalServerError)
		return
	}

	setContentTypeJson(w)
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
		http.Error(w, FailedParsingPathId, http.StatusBadRequest)
		return
	}

	if err = s.svc.DeleteBatch(r.Context(), id); err != nil {
		slog.Error("deleting batch", logging.Err(err))
		http.Error(w, "Failed to successfully delete batch", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) searchBatches(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	limit := route.GetIntQueryParam(r, "limit", 50)
	offset := route.GetIntQueryParam(r, "offset", 0)

	batches, total, err := s.svc.SearchBatches(r.Context(), fmt.Sprintf("%%%v%%", name), limit, offset)
	if err != nil {
		slog.Error("getting batches from service", logging.Err(err), slog.Group("query",
			slog.Int("limit", limit),
			slog.Int("offset", offset),
			slog.String("name", name)))
		http.Error(w, "Failed to search batches", http.StatusInternalServerError)
		return
	}

	batchDtos := make([]dto.BatchDTO, len(batches))
	for i, b := range batches {
		batchDtos[i] = *new(dto.BatchDTO).FromModel(&b)
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
		http.Error(w, FailedMarshalling, http.StatusInternalServerError)
		return
	}

	setContentTypeJson(w)
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
		http.Error(w, FailedParsingPathId, http.StatusBadRequest)
		return
	}

	batch, err := s.svc.GetBatchById(r.Context(), id)
	if err != nil {
		slog.Error("fetching batch by id", logging.Err(err), slog.String("id", id.String()))
		http.Error(w, "failed to fetch batch by id", http.StatusInternalServerError)
		return
	}

	if batch == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	batchDto := new(dto.BatchDTO).FromModel(batch)

	batchJson, err := json.Marshal(batchDto)
	if err != nil {
		slog.Error("marshalling batch to json", []any{logging.Err(err), batchDto.Slog()}...)
		http.Error(w, FailedMarshalling, http.StatusInternalServerError)
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
		http.Error(w, FailedParsingPathId, http.StatusBadRequest)
		return
	}

	var d dto.UpdateBatchDTO
	if err := readBody(r, &d); err != nil {
		slog.Error("reading body", logging.Err(err))
		http.Error(w, FailedReadingBody, http.StatusBadRequest)
		return
	}

	batch, err := s.svc.UpdateBatch(r.Context(), batchId, d.Name)
	if err != nil {
		slog.Error("updating batch", []any{logging.Err(err), d.Slog()}...)
		http.Error(w, "failed to successfully update batch", http.StatusInternalServerError)
		return
	}

	batchDto := new(dto.BatchDTO).FromModel(batch)

	jsonResp, err := json.Marshal(batchDto)
	if err != nil {
		slog.Error("marshalling batch dto", []any{logging.Err(err), batchDto.Slog()}...)
		http.Error(w, FailedMarshalling, http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(jsonResp); err != nil {
		slog.Error("writing response", logging.Err(err))
		return
	}
}

func (s *Server) addDeviceToBatch(w http.ResponseWriter, r *http.Request) {
	batchIdRaw := r.PathValue(IDKey)
	deviceIdRaw := r.PathValue(IDKey2)

	batchId, err := uuid.Parse(batchIdRaw)
	if err != nil {
		slog.Error("add device to batch parsing batch id", slog.String("batchId", batchIdRaw), logging.Err(err))
		http.Error(w, FailedParsingPathId, http.StatusBadRequest)
		return
	}

	deviceId, err := uuid.Parse(deviceIdRaw)
	if err != nil {
		slog.Error("add device to batch parsing device id", slog.String("deviceId", deviceIdRaw), logging.Err(err))
		http.Error(w, FailedParsingPathId, http.StatusBadRequest)
		return
	}

	batchDevice, err := s.svc.AddDeviceToBatch(r.Context(), batchId, deviceId)
	if err != nil {
		notFound := &service.ErrNotFound{}
		if errors.As(err, notFound) {
			slog.Error("could not find batch or device", logging.Err(err),
				slog.String("batchId", batchId.String()), slog.String("deviceId", deviceId.String()))
			http.Error(w, fmt.Sprintf("%v not found", notFound.Name), http.StatusNotFound)
			return
		}
		slog.Error("adding device to batch", logging.Err(err),
			slog.String("batchId", batchId.String()), slog.String("deviceId", deviceId.String()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	batchDeviceDto := new(dto.BatchDeviceDTO).FromModel(batchDevice)

	jsonResp, err := json.Marshal(batchDeviceDto)
	if err != nil {
		slog.Error("marshalling batch device dto", logging.Err(err))
		http.Error(w, FailedMarshalling, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(jsonResp); err != nil {
		slog.Error("writing response", logging.Err(err))
		return
	}
}
