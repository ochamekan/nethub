package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ochamekan/nethub/backend/internal/dto"
	"github.com/ochamekan/nethub/backend/internal/handlers"
	"github.com/ochamekan/nethub/backend/internal/repository"
	"github.com/ochamekan/nethub/backend/pkg/models"
	"go.uber.org/zap"
)

type mockRepo struct {
	createDeviceFn func(ctx context.Context, ip, hostname, location string, isActive bool) (models.Device, error)
	getDevicesFn   func(ctx context.Context, data dto.GetDevicesRequest) ([]models.Device, error)
	getDeviceFn    func(ctx context.Context, id int64) (models.Device, error)
	updateDeviceFn func(ctx context.Context, id int64, data dto.UpdateDeviceRequest) (models.Device, error)
	deleteDeviceFn func(ctx context.Context, id int64) (models.Device, error)
}

func (m *mockRepo) CreateDevice(ctx context.Context, ip, hostname, location string, isActive bool) (models.Device, error) {
	return m.createDeviceFn(ctx, ip, hostname, location, isActive)
}
func (m *mockRepo) GetDevices(ctx context.Context, data dto.GetDevicesRequest) ([]models.Device, error) {
	return m.getDevicesFn(ctx, data)
}
func (m *mockRepo) GetDevice(ctx context.Context, id int64) (models.Device, error) {
	return m.getDeviceFn(ctx, id)
}
func (m *mockRepo) UpdateDevice(ctx context.Context, id int64, data dto.UpdateDeviceRequest) (models.Device, error) {
	return m.updateDeviceFn(ctx, id, data)
}
func (m *mockRepo) DeleteDevice(ctx context.Context, id int64) (models.Device, error) {
	return m.deleteDeviceFn(ctx, id)
}

func newTestHandler(repo repository.DeviceRepository) *handlers.Handler {
	return handlers.New(repo, zap.NewNop())
}

type CreateDeviceTest struct {
	name           string
	requestBody    any
	repoResult     models.Device
	repoErr        error
	wantStatusCode int
}

func TestCreateDevice(t *testing.T) {
	tests := []CreateDeviceTest{
		{
			name: "success - valid device",
			requestBody: dto.CreateDeviceRequest{
				IP:       "192.168.1.1",
				Hostname: "router-01",
				Location: "Server Room A",
				IsActive: true,
			},
			repoResult:     models.Device{ID: 1, IP: "192.168.1.1", Hostname: "router-01", Location: "Server Room A", IsActive: true},
			repoErr:        nil,
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "success - IsActive defaults to false",
			requestBody: dto.CreateDeviceRequest{
				IP:       "10.10.10.10",
				Hostname: "ap-lobby",
				Location: "Lobby",
			},
			repoResult:     models.Device{ID: 2, IP: "10.10.10.10", Hostname: "ap-lobby", Location: "Lobby", IsActive: false},
			repoErr:        nil,
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "fail - missing IP",
			requestBody: dto.CreateDeviceRequest{
				IP:       "",
				Hostname: "router-01",
				Location: "Server Room A",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fail - whitespace-only IP",
			requestBody: dto.CreateDeviceRequest{
				IP:       "   ",
				Hostname: "router-01",
				Location: "Server Room A",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fail - missing location",
			requestBody: dto.CreateDeviceRequest{
				IP:       "10.0.0.1",
				Hostname: "switch-02",
				Location: "",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fail - missing hostname",
			requestBody: dto.CreateDeviceRequest{
				IP:       "10.0.0.2",
				Hostname: "",
				Location: "DC West",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fail - invalid JSON body",
			requestBody:    "this is not json",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fail - empty body",
			requestBody:    nil,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fail - repo returns internal error",
			requestBody: dto.CreateDeviceRequest{
				IP:       "172.16.0.1",
				Hostname: "server-03",
				Location: "Rack 7",
			},
			repoResult:     models.Device{},
			repoErr:        errors.New("connection refused"),
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "fail - repo returns duplicate error",
			requestBody: dto.CreateDeviceRequest{
				IP:       "192.168.1.1",
				Hostname: "router-duplicate",
				Location: "Server Room A",
				IsActive: true,
			},
			repoResult:     models.Device{},
			repoErr:        errors.New("duplicate key value violates unique constraint"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bodyBytes []byte
			if tc.requestBody != nil {
				switch v := tc.requestBody.(type) {
				case string:
					bodyBytes = []byte(v)
				default:
					b, err := json.Marshal(v)
					if err != nil {
						t.Fatalf("test setup: marshal failed: %v", err)
					}
					bodyBytes = b
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/devices", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			mock := &mockRepo{
				createDeviceFn: func(_ context.Context, _, _, _ string, _ bool) (models.Device, error) {
					return tc.repoResult, tc.repoErr
				},
			}

			h := newTestHandler(mock)
			mux := http.NewServeMux()
			mux.HandleFunc("POST /devices", h.CreateDevice)
			mux.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("status: got %d, want %d", rec.Code, tc.wantStatusCode)
			}

			if tc.wantStatusCode == http.StatusCreated {
				var got models.Device
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Errorf("expected valid JSON body on 201, got: %v", err)
				}
				if got.IP != tc.repoResult.IP || got.Hostname != tc.repoResult.Hostname || got.Location != tc.repoResult.Location || got.IsActive != tc.repoResult.IsActive {
					t.Errorf("response body mismatch")
				}
			}
		})
	}
}

type GetDeviceTest struct {
	name           string
	id             string
	repoResult     models.Device
	repoErr        error
	wantStatusCode int
}

func TestGetDevice(t *testing.T) {
	tests := []GetDeviceTest{
		{
			name:           "success - valid id",
			id:             "1",
			repoResult:     models.Device{ID: 1, IP: "192.168.1.1", Hostname: "router-01", Location: "Server Room A", IsActive: true},
			repoErr:        nil,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fail - valid id",
			id:             "5",
			repoResult:     models.Device{},
			repoErr:        repository.ErrNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "fail - not valid non-numeric id",
			id:             "abc",
			repoResult:     models.Device{},
			repoErr:        nil,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fail - valid id but server connection refused",
			id:             "3",
			repoResult:     models.Device{},
			repoErr:        errors.New("connection  refused"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/devices/"+tc.id, nil)

			rec := httptest.NewRecorder()
			mock := &mockRepo{
				getDeviceFn: func(_ context.Context, _ int64) (models.Device, error) {
					return tc.repoResult, tc.repoErr
				},
			}
			h := newTestHandler(mock)
			mux := http.NewServeMux()
			mux.HandleFunc("GET /devices/{id}", h.GetDevice)
			mux.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("status: got %d, want %d", rec.Code, tc.wantStatusCode)
			}
			if tc.wantStatusCode == http.StatusOK {
				var got models.Device
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Errorf("expected valid JSON body on 200, got: %v", err)
				}
				if got.ID != tc.repoResult.ID {
					t.Errorf("device ID: got %d, want %d", got.ID, tc.repoResult.ID)
				}
			}
		})
	}
}
