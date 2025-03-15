package mocks

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

const metricsEndpoint = "analytics/metrics"

func (ms *Server) addAnalyticsRoutes(r chi.Router) {
	r.Post("/v1/"+metricsEndpoint, ms.listMetrics)
}

func (ms *Server) listMetrics(w http.ResponseWriter, _ *http.Request) {
	start, _ := mtypes.NewRFC2822Time("Tue, 24 Sep 2024 00:00:00 +0000")
	end, _ := mtypes.NewRFC2822Time("Tue, 24 Oct 2024 00:00:00 +0000")

	resp := mtypes.MetricsResponse{
		Start:      start,
		End:        end,
		Resolution: "day",
		Duration:   "30d",
		Dimensions: []string{"time"},
		Items: []mtypes.MetricsItem{
			{
				Dimensions: []mtypes.MetricsDimension{{
					Dimension:    "time",
					Value:        "Tue, 24 Sep 2024 00:00:00 +0000",
					DisplayValue: "Tue, 24 Sep 2024 00:00:00 +0000",
				}},
				Metrics: mtypes.Metrics{
					SentCount:      ptr(uint64(4)),
					DeliveredCount: ptr(uint64(3)),
					OpenedCount:    ptr(uint64(2)),
					FailedCount:    ptr(uint64(1)),
				},
			},
		},
		Pagination: mtypes.MetricsPagination{
			Sort:  "",
			Skip:  0,
			Limit: 10,
			Total: 1,
		},
	}

	toJSON(w, resp)
}
