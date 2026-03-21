package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/ochamekan/nethub/backend/internal/dto"
	"github.com/ochamekan/nethub/backend/internal/repository"
	"go.uber.org/zap"
)

type Handler struct {
	repo   repository.DeviceRepository
	logger *zap.Logger
}

func New(repo repository.DeviceRepository, logger *zap.Logger) *Handler {
	return &Handler{repo: repo, logger: logger.With(zap.String("component", "handler"))}
}

func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(zap.String("endpoint", "CreateDevice"))
	var b dto.CreateDeviceRequest
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		logger.Warn("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// NOTE: На проде я бы использовал validator, но для данного проекта это оверкил
	if strings.TrimSpace(b.IP) == "" {
		logger.Warn("ip field is required")
		http.Error(w, "ip field is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(b.Location) == "" {
		logger.Warn("location field is required")
		http.Error(w, "location field is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(b.Hostname) == "" {
		logger.Warn("hostname is required")
		http.Error(w, "hostname is required", http.StatusBadRequest)
		return
	}

	logger.Info("creating device")
	d, err := h.repo.CreateDevice(r.Context(), b.IP, b.Hostname, b.Location, b.IsActive)
	if err != nil {
		logger.Error("failed to create new device", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		logger.Error("failed to encode json", zap.Error(err))
		return
	}
	logger.Info("device successfully created")
}

func (h *Handler) GetDevices(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(zap.String("endpoint", "GetDevices"))
	var d dto.GetDevicesRequest
	if s := r.URL.Query().Get("search"); s != "" {
		d.Search = &s
	}
	if a := r.URL.Query().Get("is_active"); a != "" {
		d.IsActive = &a
	}

	logger.Info("getting device list")
	devices, err := h.repo.GetDevices(r.Context(), d)
	if err != nil {
		logger.Error("failed to get devices list", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(devices)
	if err != nil {
		logger.Error("faield to encode devices", zap.Error(err))
		return
	}
	logger.Info("device list successfully retrieved")
}

func (h *Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(zap.String("endpoint", "GetDevice"))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		logger.Warn("failed to extract id from searchParams", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info("getting device by id")
	d, err := h.repo.GetDevice(r.Context(), int64(id))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Warn("device not found")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logger.Error("failed to get device", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		logger.Error("faield to encode device", zap.Error(err))
		return
	}
	logger.Info("device successfully retrieved")

}

func (h *Handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(zap.String("endpoint", "UpdateDevice"))
	var b dto.UpdateDeviceRequest
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		logger.Warn("failed to extract id from search params", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		logger.Warn("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info("updating device by id")
	d, err := h.repo.UpdateDevice(r.Context(), int64(id), b)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Warn("device not found")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logger.Error("failed to update device", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		logger.Error("faield to encode device", zap.Error(err))
		return
	}
	logger.Info("device successfully updated")
}

func (h *Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(zap.String("endpoint", "DeleteDevice"))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		logger.Warn("failed to extract id from searchParams", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info("deleting device by id")
	d, err := h.repo.DeleteDevice(r.Context(), int64(id))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Warn("device not found")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logger.Error("failed to delete device", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		logger.Error("faield to encode device", zap.Error(err))
		return
	}
	logger.Info("device successfully deleted")
}
