package mocks

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addIPWarmupsRoutes(r chi.Router) {
	r.Get("/ip_warmups", ms.listIPWarmups)
	r.Get("/ip_warmups/{ip}", ms.getIPWarmup)
	r.Post("/ip_warmups/{ip}", ms.createIPWarmup)
	r.Delete("/ip_warmups/{ip}", ms.deleteIPWarmup)
}

func (ms *Server) listIPWarmups(w http.ResponseWriter, r *http.Request) {

	var items []mtypes.IPWarmup
	if r.FormValue("page") == "" {
		items = []mtypes.IPWarmup{
			{
				IP:               "1.0.0.1",
				SentWithinStage:  "0%",
				Throttle:         78,
				StageNumber:      3,
				StageStartVolume: 14000,
				StageStartTime:   time.Now(),
				StageVolumeLimit: 4000,
			},
			{
				IP:               "1.0.0.2",
				SentWithinStage:  "25%",
				Throttle:         90,
				StageNumber:      4,
				StageStartVolume: 10000,
				StageStartTime:   time.Now(),
				StageVolumeLimit: 8000,
			},
		}
	}

	toJSON(w, mtypes.ListIPWarmupsResponse{
		Items: items,
		Paging: mtypes.Paging{
			First: getPageURL(r, url.Values{
				"page": []string{"first"},
			}),
			Last: getPageURL(r, url.Values{
				"page": []string{"last"},
			}),
			Next: getPageURL(r, url.Values{
				"page": []string{"next"},
			}),
			Previous: getPageURL(r, url.Values{
				"page": []string{"prev"},
			}),
		},
	})
}

func (ms *Server) getIPWarmup(w http.ResponseWriter, r *http.Request) {
	toJSON(w, mtypes.IPWarmupDetailsResponse{
		Details: mtypes.IPWarmupDetails{
			SentWithinStage:   "20%",
			Throttle:          78,
			StageNumber:       2,
			StageStartVolume:  10000,
			StageStartTime:    time.Now(),
			StageVolumeLimit:  4000,
			StageStartedAt:    time.Now(),
			HourStartedAt:     time.Now(),
			PlanStartedAt:     time.Now(),
			PlanLastUpdatedAt: time.Now(),
			TotalStages:       15,
			StageHistory: []mtypes.IPWarmupStageHistory{
				{
					FirstUpdatedAt: time.Now(),
					CompletedAt:    time.Now(),
					Limit:          1000,
				},
				{
					FirstUpdatedAt: time.Now(),
					CompletedAt:    time.Now(),
					Limit:          2000,
				},
			},
		},
	})

}

func (ms *Server) createIPWarmup(w http.ResponseWriter, r *http.Request) {
	toJSON(w, okResp{Message: "success"})
}

func (ms *Server) deleteIPWarmup(w http.ResponseWriter, r *http.Request) {
	toJSON(w, okResp{Message: "success"})
}
