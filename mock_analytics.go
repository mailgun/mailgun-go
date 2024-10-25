package mailgun

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (ms *mockServer) addAnalyticsRoutes(r chi.Router) {
	r.Post("/v1/"+metricsEndpoint, ms.listMetrics)
}

func (ms *mockServer) listMetrics(w http.ResponseWriter, _ *http.Request) {
	start, _ := NewRFC2822Time("Tue, 24 Sep 2024 00:00:00 +0000")
	end, _ := NewRFC2822Time("Tue, 24 Oct 2024 00:00:00 +0000")

	resp := MetricsResponse{
		Start:      start,
		End:        end,
		Resolution: "day",
		Duration:   "30d",
		Dimensions: []string{"time"},
		Items: []MetricsItem{
			{
				Dimensions: []MetricsDimension{{
					Dimension:    "time",
					Value:        "Tue, 24 Sep 2024 00:00:00 +0000",
					DisplayValue: "Tue, 24 Sep 2024 00:00:00 +0000",
				}},
				Metrics: Metrics{
					SentCount:      ptr(uint64(4)),
					DeliveredCount: ptr(uint64(3)),
					OpenedCount:    ptr(uint64(2)),
					FailedCount:    ptr(uint64(1)),
				},
			},
		},
		Pagination: MetricsPagination{
			Sort:  "",
			Skip:  0,
			Limit: 10,
			Total: 1,
		},
	}

	toJSON(w, resp)
}
