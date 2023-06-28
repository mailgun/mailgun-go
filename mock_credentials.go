package mailgun

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (ms *mockServer) addCredentialsRoutes(r chi.Router) {
	r.Get("/domains/{domain}/credentials", ms.listCredentials)
	r.Put("/domains/{domain}/credentials/{login}", ms.updateCredential)
	r.Delete("/domains/{domain}/credentials/{login}", ms.deleteCredential)
	r.Post("/domains/{domain}/credentials", ms.createCredential)

	ms.credentials = append(ms.credentials, Credential{
		CreatedAt: RFC2822Time(time.Now()),
		Login:     "alice",
		Password:  "alices_password",
	})

	ms.credentials = append(ms.credentials, Credential{
		CreatedAt: RFC2822Time(time.Now()),
		Login:     "bob",
		Password:  "bobs_password",
	})
}

func (ms *mockServer) listCredentials(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	limit := stringToInt(r.FormValue("limit"))
	if limit == 0 {
		limit = 100
	}

	skip := stringToInt(r.FormValue("skip"))

	var results []Credential

	if skip > 0 {
		if len(ms.credentials[skip:]) < limit {
			results = ms.credentials[skip:]
		} else {
			results = ms.credentials[skip : skip+limit-1]
		}
	} else {
		if len(ms.credentials) < limit {
			results = ms.credentials
		} else {
			results = ms.credentials[:limit-1]
		}
	}

	toJSON(w, credentialsListResponse{
		Items:      results,
		TotalCount: len(results),
	})
}

func (ms *mockServer) createCredential(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	if len(login) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\": \"login is required\"}"))
		return
	}

	if len(password) < 5 || len(password) > 32 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\": \"password should be in range 5-32 characters\"}"))
		return
	}

	for _, existingCredential := range ms.credentials {
		if existingCredential.Login == login {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"message\": \"Created 1 credentials pair(s)\"}"))
			return
		}
	}

	ms.credentials = append(ms.credentials, Credential{Login: login, Password: password, CreatedAt: RFC2822Time(time.Now())})

	toJSON(w, map[string]interface{}{
		"message": "Credentials created",
	})
}

func (ms *mockServer) deleteCredential(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	login := chi.URLParam(r, "login")
	domain := chi.URLParam(r, "domain")

	for i, credential := range ms.credentials {
		if credential.Login == login || credential.Login == login+"@"+domain {
			continue
		}

		ms.credentials = append(ms.credentials[:i], ms.credentials[i+1:]...)
		toJSON(w, map[string]interface{}{
			"message": "Credentials have been deleted",
			"spec":    login,
		})
		return
	}

	w.WriteHeader(http.StatusNotFound)
	toJSON(w, map[string]string{"message": "Credentials not found"})
}

func (ms *mockServer) updateCredential(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	domain := chi.URLParam(r, "domain")
	login := chi.URLParam(r, "login")
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	password := r.FormValue("password")
	if len(password) < 5 || len(password) > 32 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\": \"password should be in range 5-32 characters\"}"))
		return
	}

	for i, credential := range ms.credentials {
		if credential.Login == login || credential.Login == login+"@"+domain {
			ms.credentials[i].Password = password
			toJSON(w, map[string]interface{}{
				"message": "Password changed",
			})
			return
		}

	}

	w.WriteHeader(http.StatusNotFound)
	toJSON(w, map[string]string{"message": "Credentials not found"})
}
