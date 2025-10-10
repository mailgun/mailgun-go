package mocks

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addIPWarmupsRoutes(r chi.Router) {
	r.Get("/ip_warmups", ms.listIPWarmups)
}

func (ms *Server) listIPWarmups(w http.ResponseWriter, r *http.Request) {

	var items []mtypes.IPWarmup
	if r.FormValue("page") == "" {
		items = []mtypes.IPWarmup{
			{
				IP:               "172.0.0.1",
				SentWithinStage:  "100",
				Throttle:         100,
				StageNumber:      1,
				StageStartVolume: 100,
				StageStartTime:   "2021-01-01T00:00:00Z",
				StageVolumeLimit: 1000,
			},
			{
				IP:               "172.0.0.2",
				SentWithinStage:  "100",
				Throttle:         78,
				StageNumber:      2,
				StageStartVolume: 1000,
				StageStartTime:   "2021-01-01T02:00:00Z",
				StageVolumeLimit: 10000,
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
