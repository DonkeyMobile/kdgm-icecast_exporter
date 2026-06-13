package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestExporterCollect(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		status         int
		wantUp         float64
		wantMetric     string
		wantScrapeErrs float64
	}{
		{
			name:   "multiple sources",
			body:   `{"icestats":{"server_start_iso8601":"2026-01-01T00:00:00+0000","source":[{"listeners":10,"listenurl":"http://localhost:8000/live","server_type":"audio/mpeg","stream_start_iso8601":"2026-06-01T12:00:00+0000"},{"listeners":5,"listenurl":"http://localhost:8000/test","server_type":"audio/ogg","stream_start_iso8601":"2026-06-01T13:00:00+0000"}]}}`,
			status: 200,
			wantUp: 1,
			wantMetric: `
# HELP icecast_listeners The number of currently connected listeners.
# TYPE icecast_listeners gauge
icecast_listeners{listenurl="http://localhost:8000/live",server_type="audio/mpeg"} 10
icecast_listeners{listenurl="http://localhost:8000/test",server_type="audio/ogg"} 5
`,
		},
		{
			name:   "single source (object format)",
			body:   `{"icestats":{"server_start_iso8601":"2026-01-01T00:00:00+0000","source":{"listeners":7,"listenurl":"http://localhost:8000/live","server_type":"audio/mpeg","stream_start_iso8601":"2026-06-01T12:00:00+0000"}}}`,
			status: 200,
			wantUp: 1,
			wantMetric: `
# HELP icecast_listeners The number of currently connected listeners.
# TYPE icecast_listeners gauge
icecast_listeners{listenurl="http://localhost:8000/live",server_type="audio/mpeg"} 7
`,
		},
		{
			name:   "no sources",
			body:   `{"icestats":{"server_start_iso8601":"2026-01-01T00:00:00+0000"}}`,
			status: 200,
			wantUp: 1,
		},
		{
			name:           "server unreachable",
			status:         -1, // signal to close server before request
			wantUp:         0,
			wantScrapeErrs: 1,
		},
		{
			name:           "invalid JSON",
			body:           `not json`,
			status:         200,
			wantUp:         1,
			wantScrapeErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.body))
			}))

			uri := srv.URL
			if tt.status == -1 {
				srv.Close()
			} else {
				defer srv.Close()
			}

			e := NewExporter(uri, 5*time.Second)
			reg := prometheus.NewRegistry()
			reg.MustRegister(e)

			metrics, err := reg.Gather()
			if err != nil {
				t.Fatalf("gather failed: %v", err)
			}

			for _, mf := range metrics {
				if mf.GetName() == "icecast_up" {
					got := mf.GetMetric()[0].GetGauge().GetValue()
					if got != tt.wantUp {
						t.Errorf("icecast_up = %v, want %v", got, tt.wantUp)
					}
				}
			}

			if tt.wantMetric != "" {
				metricName := strings.Fields(strings.Split(tt.wantMetric, "\n")[3])[0]
				metricName = strings.Split(metricName, "{")[0]
				if err := testutil.CollectAndCompare(e, strings.NewReader(tt.wantMetric), metricName); err != nil {
					t.Errorf("metric mismatch: %v", err)
				}
			}

			if tt.wantScrapeErrs > 0 {
				got := testutil.ToFloat64(e.scrapeErrors)
				if got != tt.wantScrapeErrs {
					t.Errorf("scrape_errors = %v, want %v", got, tt.wantScrapeErrs)
				}
			}
		})
	}
}
