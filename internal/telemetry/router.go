package telemetry

import "net/http"

func NewRouter(handler *Handler, apiKeys ...string) *http.ServeMux {
	mux := http.NewServeMux()

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /api/v1/scooters", handler.FindScooters)
	apiMux.HandleFunc("POST /api/v1/events", handler.ReportEvent)

	mux.Handle("/api/v1/", AuthMiddleware(apiKeys)(apiMux))
	mux.HandleFunc("GET /healthz", HealthzHandler)

	return mux
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}
