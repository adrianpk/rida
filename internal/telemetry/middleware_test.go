package telemetry

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	validKeys := []string{"demo-api-key"}
	auth := AuthMiddleware(validKeys)

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok")) // ignore error for test
	})

	handler := auth(dummyHandler)

	tests := []struct {
		name       string
		apiKey     string
		clientID   string
		wantStatus int
	}{
		{"valid key", "demo-api-key", "test-client", http.StatusOK},
		{"invalid key", "bad-key", "test-client", http.StatusUnauthorized},
		{"missing key", "", "test-client", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}

			if tt.clientID != "" {
				req.Header.Set("X-Client-ID", tt.clientID)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
			}
		})
	}
}
