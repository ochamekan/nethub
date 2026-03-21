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
		logger.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// NOTE: На проде я бы использовал validator, но для данного проекта это оверкил
	if strings.TrimSpace(b.IP) == "" {
		logger.Warn("ip field is required")
		http.Error(w, "ip field is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(b.Hostname) == "" {
		logger.Warn("hostname is required")
		http.Error(w, "hostname is required", http.StatusBadRequest)
		return
	}

	d, err := h.repo.CreateDevice(r.Context(), b.IP, b.Hostname, b.IsActive)
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
}

func (h *Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(zap.String("endpoint", "GetDevice"))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		logger.Warn("failed to extract id from searchParams", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	d, err := h.repo.GetDevice(r.Context(), int64(id))
	if err != nil {
		logger.Error("failed to get device", zap.Error(err))
		if errors.Is(repository.ErrNotFound, err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		logger.Error("faield to encode device", zap.Error(err))
		return
	}
}
