package telemetry_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/adrianpk/rida/internal/telemetry"
)

func TestAreaFromQuery(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		want    telemetry.Query
		wantErr bool
	}{
		{
			name: "valid area",
			params: map[string]string{
				"minLat": "51.0",
				"minLng": "17.0",
				"maxLat": "52.0",
				"maxLng": "18.0",
				"status": "free",
			},
			want: telemetry.Query{
				Area: telemetry.Area{
					MinLat: 51.0,
					MinLng: 17.0,
					MaxLat: 52.0,
					MaxLng: 18.0,
				},
				Status: "free",
			},
			wantErr: false,
		},
		{
			name: "missing param",
			params: map[string]string{
				"minLat": "51.0",
				"minLng": "17.0",
				"maxLat": "52.0",
				// maxLng missing
				"status": "free",
			},
			wantErr: true,
		},
		{
			name: "invalid float",
			params: map[string]string{
				"minLat": "not-a-float",
				"minLng": "17.0",
				"maxLat": "52.0",
				"maxLng": "18.0",
				"status": "free",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := url.Values{}
			for k, v := range tt.params {
				u.Set(k, v)
			}
			r := &http.Request{URL: &url.URL{RawQuery: u.Encode()}}
			qry, err := telemetry.NewQuery(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && qry != tt.want {
				t.Errorf("expected query: %+v, got: %+v", tt.want, qry)
			}
		})
	}
}
