package mailgun

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (ms *mockServer) addStatsRoutes(r chi.Router) {
	r.Get("/{domain}/stats", ms.getStats)
	r.Get("/{domain}/stats/total", ms.getStatsTotal)

	ms.stats = append(ms.stats, Stats{
		Time:         "",
		Accepted:     Accepted{Total: 100, Incoming: 30, Outgoing: 70},
		Delivered:    Delivered{Total: 100, Http: 100, Smtp: 0},
		Failed:       Failed{Temporary: Temporary{Espblock: 0}, Permanent: Permanent{Total: 0}},
		Stored:       Total{Total: 5},
		Opened:       Total{Total: 40},
		Clicked:      Total{Total: 50},
		Unsubscribed: Total{Total: 5},
		Complained:   Total{Total: 2},
	})

}

func (ms *mockServer) getStats(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("{\"message\": \"this route was deprecated 2 years ago and is no longer available\"}"))
}

func (ms *mockServer) getStatsTotal(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	eventName := r.Form.Get("event")
	if len(eventName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("please provide at least one event: Missing mandatory parameter: event"))
		return
	}

	resolution := r.Form.Get("resolution")
	if len(resolution) > 0 && resolution != "day" && resolution != "hour" && resolution != "month" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("'resolution' parameter is invalid, valid values: 'month', 'day', 'hour'"))
		return
	}

	end, _ := NewRFC2822Time(time.Now().String())
	start, _ := NewRFC2822Time(time.Now().Add(-10 * time.Minute).String())
	toJSON(w, &statsTotalResponse{
		Stats: ms.stats,
		End:   end.String(),
		Start: start.String(),
	})
}
