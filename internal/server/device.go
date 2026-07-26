package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/logging"

	"github.com/google/uuid"
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

func (s *Server) searchDevices(w http.ResponseWriter, r *http.Request) {
	search := new(dto.SearchDTO).FromRequest(r)

	devicesPage, err := s.svc.SearchDevices(*search)
	if err != nil {
		slog.Error("searching devices", []any{logging.Err(err), search.Slog()}...)
		http.Error(w, "Failed to search devices", http.StatusInternalServerError)
		return
	}

	deviceDtos := make([]dto.DeviceDTO, len(devicesPage.Data))
	for i, b := range devicesPage.Data {
		deviceDtos[i] = *new(dto.DeviceDTO).FromModel(&b)
	}

	page := dto.PageDTO[dto.DeviceDTO]{
		Limit:  search.Limit,
		Offset: search.Offset,
		Total:  devicesPage.Total,
		Data:   deviceDtos,
	}

	pageJson, err := json.Marshal(page)
	if err != nil {
		slog.Error("marshalling device page to json", logging.Err(err))
		http.Error(w, FailedMarshalling, http.StatusInternalServerError)
		return
	}

	setContentTypeJson(w)
	if _, err := w.Write(pageJson); err != nil {
		slog.Error("writing device page", logging.Err(err))
		return
	}
}

func (s *Server) getDeviceById(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue(IDKey)
	id, err := uuid.Parse(rawId)
	if err != nil {
		slog.Error("parsing device id from path", logging.Err(err), slog.String("id", rawId))
		http.Error(w, FailedParsingPathId, http.StatusBadRequest)
		return
	}

	device, err := s.svc.GetDeviceById(id)
	if err != nil {
		slog.Error("fetching device by id", logging.Err(err), slog.String("id", id.String()))
		http.Error(w, "failed to fetch device by id", http.StatusInternalServerError)
		return
	}

	if device == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	deviceDto := new(dto.DeviceDTO).FromModel(device)

	deviceJson, err := json.Marshal(deviceDto)
	if err != nil {
		slog.Error("marshalling device to json", []any{logging.Err(err), deviceDto.Slog()}...)
		http.Error(w, FailedMarshalling, http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(deviceJson); err != nil {
		slog.Error("failed to write device", []any{logging.Err(err), deviceDto.Slog()}...)
		return
	}
}

func (s *Server) deleteDeviceById(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue(IDKey)
	id, err := uuid.Parse(rawId)
	if err != nil {
		slog.Error("parsing device id", slog.String(IDKey, rawId), logging.Err(err))
		http.Error(w, FailedParsingPathId, http.StatusBadRequest)
		return
	}

	if err = s.svc.DeleteDevice(id); err != nil {
		slog.Error("deleting device", logging.Err(err))
		http.Error(w, "Failed to successfully delete device", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) updateDevice(w http.ResponseWriter, r *http.Request) {
	deviceDto := &dto.UpdateDeviceDTO{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&deviceDto); err != nil {
		slog.Error("unable to decode body on device update", []any{logging.Err(err), deviceDto.Slog()}...)
		http.Error(w, FailedReadingBody, http.StatusBadRequest)
		return
	}

	err := s.validate.Struct(deviceDto)
	if err != nil {
		slog.Error("faled to validate update device dto", logging.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.Context().Value(IDKey).(uuid.UUID)

	device, err := s.svc.UpdateDevice(id, deviceDto.ToUpdateDeviceParams(id))
	if err != nil {
		slog.Error("failed to successfully update device", logging.Err(err))
		http.Error(w, "Failed to successfully update device", http.StatusInternalServerError)
		return
	}

	returnDto := new(dto.DeviceDTO).FromModel(device)

	deviceJson, err := json.Marshal(returnDto)
	if err != nil {
		slog.Error("failed to marshal device after update", []any{logging.Err(err), returnDto.Slog()}...)
		http.Error(w, FailedMarshalling, http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(deviceJson); err != nil {
		slog.Error("failed to write", logging.Err(err))
		return
	}
}
