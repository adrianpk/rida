package telemetry_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adrianpk/rida/internal/telemetry"
	"github.com/google/uuid"
)

func TestGetScooterHandler(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name       string
		id         string
		svc        *mockService
		wantStatus int
		wantBody   string
	}{
		{
			name: "happy path",
			id:   id.String(),
			svc: &mockService{
				GetScooterFunc: happyGetScooter(id),
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"id":"` + id.String() + `","status":"free","lat":0,"lng":0,"updatedAt":"`,
		},
		{
			name: "not found",
			id:   id.String(),
			svc: &mockService{
				GetScooterFunc: notFoundGetScooter,
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "not found",
		},
		{
			name: "invalid id",
			id:   "not-a-uuid",
			svc: &mockService{
				GetScooterFunc: alwaysNilGetScooter,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid scooter id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := telemetry.NewHandler(tt.svc)
			r := httptest.NewRequest(http.MethodGet, "/scooters/"+tt.id, nil)
			r.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()
			h.GetScooter(w, r)

			resp := w.Result()
			body := w.Body.String()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			if tt.wantStatus == http.StatusOK {
				if len(body) == 0 || body[:len(tt.wantBody)] != tt.wantBody {
					t.Errorf("expected body to start with %q, got %q", tt.wantBody, body)
				}
			} else {
				if !bytes.Contains([]byte(body), []byte(tt.wantBody)) {
					t.Errorf("expected body to contain %q, got %q", tt.wantBody, body)
				}
			}
		})
	}
}

func TestUpdateScooterHandler(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name       string
		id         string
		body       telemetry.Scooter
		svc        *mockService
		wantStatus int
		wantBody   string
	}{
		{
			name: "happy path",
			id:   id.String(),
			body: telemetry.Scooter{ID: uuid.New(), Status: telemetry.StatusOccupied},
			svc: &mockService{
				UpdateScooterFunc: happyUpdateScooter(id),
			},
			wantStatus: http.StatusNoContent,
			wantBody:   "",
		},
		{
			name: "invalid id",
			id:   "not-a-uuid",
			body: telemetry.Scooter{},
			svc: &mockService{
				UpdateScooterFunc: alwaysNilUpdateScooter,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid scooter id",
		},
		{
			name: "unmarshal error",
			id:   id.String(),
			body: telemetry.Scooter{}, // will send invalid JSON
			svc: &mockService{
				UpdateScooterFunc: alwaysNilUpdateScooter,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "unmarshalable request body",
		},
		{
			name: "service error",
			id:   id.String(),
			body: telemetry.Scooter{ID: id, Status: telemetry.StatusOccupied},
			svc: &mockService{
				UpdateScooterFunc: failUpdateScooter,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := telemetry.NewHandler(tt.svc)
			var reqBody []byte
			if tt.name == "unmarshal error" {
				reqBody = []byte("not-json")
			} else {
				var err error
				reqBody, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("marshal error: %v", err)
				}
			}

			r := httptest.NewRequest(http.MethodPut, "/scooters/"+tt.id, bytes.NewReader(reqBody))
			r.SetPathValue("id", tt.id)

			w := httptest.NewRecorder()
			h.UpdateScooter(w, r)

			resp := w.Result()
			body := w.Body.String()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			if tt.wantBody != "" && !bytes.Contains([]byte(body), []byte(tt.wantBody)) {
				t.Errorf("expected body to contain %q, got %q", tt.wantBody, body)
			}
		})
	}
}

func TestFindScootersHandler(t *testing.T) {
	id := uuid.New()
	qry := telemetry.Query{
		Area: telemetry.Area{
			MinLat: 51,
			MinLng: 17,
			MaxLat: 52,
			MaxLng: 18,
		},
		Status: telemetry.StatusFree,
	}
	tests := []struct {
		name       string
		params     string
		svc        *mockService
		wantStatus int
		wantBody   string
	}{
		{
			name:   "happy path",
			params: "?minLat=51&minLng=17&maxLat=52&maxLng=18&status=free",
			svc: &mockService{
				FindScootersFunc: happyFindScooters(qry, []telemetry.Scooter{{ID: id, Status: telemetry.StatusFree}}),
			},
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":"` + id.String() + `","status":"free","lat":0,"lng":0,"updatedAt":"`,
		},
		{
			name:   "invalid area param",
			params: "?minLat=not-a-float&minLng=17&maxLat=52&maxLng=18&status=free",
			svc: &mockService{
				FindScootersFunc: alwaysNilFindScooters,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid area parameters",
		},
		{
			name:   "service error",
			params: "?minLat=51&minLng=17&maxLat=52&maxLng=18&status=free",
			svc: &mockService{
				FindScootersFunc: failFindScooters,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := telemetry.NewHandler(tt.svc)
			r := httptest.NewRequest(http.MethodGet, "/scooters"+tt.params, nil)
			w := httptest.NewRecorder()

			h.FindScooters(w, r)

			resp := w.Result()
			body := w.Body.String()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			if tt.wantStatus == http.StatusOK {
				if len(body) == 0 || body[:len(tt.wantBody)] != tt.wantBody {
					t.Errorf("expected body to start with %q, got %q", tt.wantBody, body)
				}
			} else {
				if !bytes.Contains([]byte(body), []byte(tt.wantBody)) {
					t.Errorf("expected body to contain %q, got %q", tt.wantBody, body)
				}
			}
		})
	}
}

func TestReportEventHandler(t *testing.T) {
	validEvent := telemetry.Event{
		ID:        uuid.New(),
		ScooterID: uuid.New(),
		Type:      telemetry.EventTripStart,
		Timestamp: time.Now(),
		Lat:       51.1,
		Lng:       17.0,
	}

	tests := []struct {
		name       string
		body       interface{}
		mockSvc    *mockService
		wantStatus int
		wantBody   string
	}{
		{
			name: "happy path",
			body: validEvent,
			mockSvc: &mockService{
				ReportEventFunc: func(ctx context.Context, e telemetry.Event) error {
					if e.ID != validEvent.ID ||
						e.ScooterID != validEvent.ScooterID ||
						e.Type != validEvent.Type ||
						e.Lat != validEvent.Lat ||
						e.Lng != validEvent.Lng {
						return errors.New("event mismatch")
					}
					// Do NOT compare e.Timestamp
					return nil
				},
			},
			wantStatus: http.StatusCreated,
			wantBody:   "",
		},
		{
			name: "invalid payload",
			body: "not-json",
			mockSvc: &mockService{
				ReportEventFunc: func(ctx context.Context, e telemetry.Event) error { return nil },
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid event payload",
		},
		{
			name: "service error",
			body: validEvent,
			mockSvc: &mockService{
				ReportEventFunc: func(ctx context.Context, e telemetry.Event) error { return errors.New("fail") },
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := telemetry.NewHandler(tt.mockSvc)
			var reqBody []byte
			if s, ok := tt.body.(string); ok {
				reqBody = []byte(s)
			} else {
				var err error
				reqBody, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("marshal error: %v", err)
				}
			}

			r := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			h.ReportEvent(w, r)

			resp := w.Result()
			body := w.Body.String()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
			if tt.wantBody != "" && !bytes.Contains([]byte(body), []byte(tt.wantBody)) {
				t.Errorf("expected body to contain %q, got %q", tt.wantBody, body)
			}
		})
	}
}

type mockService struct {
	GetScooterFunc    func(ctx context.Context, id uuid.UUID) (telemetry.Scooter, error)
	UpdateScooterFunc func(ctx context.Context, s telemetry.Scooter) error
	FindScootersFunc  func(ctx context.Context, qry telemetry.Query) ([]telemetry.Scooter, error)
	ReportEventFunc   func(ctx context.Context, e telemetry.Event) error
}

func (m *mockService) GetScooter(ctx context.Context, id uuid.UUID) (telemetry.Scooter, error) {
	return m.GetScooterFunc(ctx, id)
}

func (m *mockService) UpdateScooter(ctx context.Context, s telemetry.Scooter) error {
	return m.UpdateScooterFunc(ctx, s)
}

func (m *mockService) FindScooters(ctx context.Context, qry telemetry.Query) ([]telemetry.Scooter, error) {
	return m.FindScootersFunc(ctx, qry)
}

func (m *mockService) ReportEvent(ctx context.Context, e telemetry.Event) error {
	if m.ReportEventFunc != nil {
		return m.ReportEventFunc(ctx, e)
	}
	return nil
}

func happyGetScooter(expectedID uuid.UUID) func(context.Context, uuid.UUID) (telemetry.Scooter, error) {
	return func(ctx context.Context, gotID uuid.UUID) (telemetry.Scooter, error) {
		if gotID != expectedID {
			return telemetry.Scooter{}, errors.New("wrong id")
		}
		return telemetry.Scooter{ID: expectedID, Status: telemetry.StatusFree}, nil
	}
}

func notFoundGetScooter(context.Context, uuid.UUID) (telemetry.Scooter, error) {
	return telemetry.Scooter{}, errors.New("not found")
}

func alwaysNilGetScooter(context.Context, uuid.UUID) (telemetry.Scooter, error) {
	return telemetry.Scooter{}, nil
}

func happyUpdateScooter(expectedID uuid.UUID) func(context.Context, telemetry.Scooter) error {
	return func(ctx context.Context, s telemetry.Scooter) error {
		if s.ID != expectedID {
			return errors.New("id not overwritten")
		}
		return nil
	}
}

func alwaysNilUpdateScooter(context.Context, telemetry.Scooter) error {
	return nil
}

func failUpdateScooter(context.Context, telemetry.Scooter) error {
	return errors.New("fail")
}

func happyFindScooters(expectedQuery telemetry.Query, result []telemetry.Scooter) func(context.Context, telemetry.Query) ([]telemetry.Scooter, error) {
	return func(ctx context.Context, gotQuery telemetry.Query) ([]telemetry.Scooter, error) {
		if gotQuery != expectedQuery {
			return nil, errors.New("wrong params")
		}
		return result, nil
	}
}

func alwaysNilFindScooters(context.Context, telemetry.Query) ([]telemetry.Scooter, error) {
	return nil, nil
}

func failFindScooters(context.Context, telemetry.Query) ([]telemetry.Scooter, error) {
	return nil, errors.New("fail")
}
