package mailgun

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/go-chi/chi"
)

type MockServer struct {
	srv *httptest.Server

	domainIPS   []string
	domainList  []domainContainer
	exportList  []Export
	mailingList []mailingListContainer
}

func NewMockServer() MockServer {
	ms := MockServer{}

	// Add all our handlers
	r := chi.NewRouter()

	r.Route("/v3", func(r chi.Router) {
		ms.addIPRoutes(r)
		ms.addExportRoutes(r)
		ms.addDomainRoutes(r)
		ms.addMailingListRoutes(r)
	})

	// Start the server
	ms.srv = httptest.NewServer(r)
	return ms
}

func (ms *MockServer) Stop() {
	ms.srv.Close()
}

func (ms *MockServer) URL() string {
	return ms.srv.URL + "/v3"
}

func toJSON(w http.ResponseWriter, obj interface{}) {
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}

func stringToBool(b string) bool {
	result, err := strconv.ParseBool(b)
	if err != nil {
		panic(err)
	}
	return result
}

func stringToInt(b string) int {
	if b == "" {
		return 0
	}

	result, err := strconv.ParseInt(b, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(result)
}

// Given the page direction, pivot value and limit, calculate the offsets for the slice
func pageOffsets(pivotIdx []string, pivotDir, pivotVal string, limit int) (int, int) {
	switch pivotDir {
	case "first":
		if limit < len(pivotIdx) {
			return 0, limit
		}
		return 0, len(pivotIdx)
	case "last":
		if limit < len(pivotIdx) {
			return len(pivotIdx) - limit, len(pivotIdx)
		}
		return 0, len(pivotIdx)
	case "next":
		for i, item := range pivotIdx {
			if item == pivotVal {
				offset := i + 1 + limit
				if offset > len(pivotIdx) {
					offset = len(pivotIdx)
				}
				return i + 1, offset
			}
		}
		return 0, 0
	case "prev":
		for i, item := range pivotIdx {
			if item == pivotVal {
				offset := i + limit
				if offset > len(pivotIdx) {
					offset = len(pivotIdx)
				}
				return i, offset
			}
		}
		return 0, 0
	}

	if limit > len(pivotIdx) {
		return 0, len(pivotIdx)
	}
	return 0, limit
}
