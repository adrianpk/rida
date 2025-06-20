package telemetry

import (
	"context"
	"log"
	"net/http"
)

type contextKey string

const clientIDKey contextKey = "clientID"

func AuthMiddleware(validAPIKeys []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			clientID := r.Header.Get("X-Client-ID")

			valid := false
			for _, k := range validAPIKeys {
				if apiKey == k {
					valid = true
					break
				}
			}

			if !valid {
				log.Printf("auth: invalid API key: %q, client: %q, path: %s", mask(apiKey), clientID, r.URL.Path)
				http.Error(w, "invalid API key", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), clientIDKey, clientID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func mask(s string) string {
	if len(s) > 4 {
		return s[:2] + "..." + s[len(s)-2:]
	}

	return "[redacted]"
}

func ClientID(ctx context.Context) (string, bool) {
	clientID, ok := ctx.Value(clientIDKey).(string)
	return clientID, ok
}
