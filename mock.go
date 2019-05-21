package mailgun

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

// A mailgun api mock suitable for testing
type MockServer struct {
	srv *httptest.Server

	domainIPS   []string
	domainList  []domainContainer
	exportList  []Export
	mailingList []mailingListContainer
	routeList   []Route
	events      []Event
	webhooks    WebHooksListResponse
}

// Create a new instance of the mailgun API mock server
func NewMockServer() MockServer {
	ms := MockServer{}

	// Add all our handlers
	r := chi.NewRouter()

	r.Route("/v3", func(r chi.Router) {
		ms.addIPRoutes(r)
		ms.addExportRoutes(r)
		ms.addDomainRoutes(r)
		ms.addMailingListRoutes(r)
		ms.addEventRoutes(r)
		ms.addMessagesRoutes(r)
		ms.addValidationRoutes(r)
		ms.addRoutes(r)
		ms.addWebhookRoutes(r)
	})

	// Start the server
	ms.srv = httptest.NewServer(r)
	return ms
}

// Stop the server
func (ms *MockServer) Stop() {
	ms.srv.Close()
}

// URL returns the URL used to connect to the mock server
func (ms *MockServer) URL() string {
	return ms.srv.URL + "/v3"
}

func toJSON(w http.ResponseWriter, obj interface{}) {
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}

func stringToBool(v string) bool {
	lower := strings.ToLower(v)
	if lower == "yes" || lower == "no" {
		return lower == "yes"
	}

	if v == "" {
		return false
	}

	result, err := strconv.ParseBool(v)
	if err != nil {
		panic(err)
	}
	return result
}

func stringToInt(v string) int {
	if v == "" {
		return 0
	}

	result, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(result)
}

func stringToMap(v string) map[string]interface{} {
	if v == "" {
		return nil
	}

	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(v), &result)
	if err != nil {
		panic(err)
	}
	return result
}

func parseAddress(v string) string {
	if v == "" {
		return ""
	}
	e, err := mail.ParseAddress(v)
	if err != nil {
		panic(err)
	}
	return e.Address
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
				if i == 0 {
					return 0, 0
				}

				offset := i - limit
				if offset < 0 {
					offset = 0
				}
				return offset, i
			}
		}
		return 0, 0
	}

	if limit > len(pivotIdx) {
		return 0, len(pivotIdx)
	}
	return 0, limit
}

func getPageURL(r *http.Request, params url.Values) string {
	if r.FormValue("limit") != "" {
		params.Add("limit", r.FormValue("limit"))
	}
	return "http://" + r.Host + r.URL.EscapedPath() + "?" + params.Encode()
}

// randomString generates a string of given length, but random content.
// All content will be within the ASCII graphic character set.
// (Implementation from Even Shaw's contribution on
// http://stackoverflow.com/questions/12771930/what-is-the-fastest-way-to-generate-a-long-random-string-in-go).
func randomString(n int, prefix string) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return prefix + string(bytes)
}

func randomEmail(prefix, domain string) string {
	return strings.ToLower(fmt.Sprintf("%s@%s", randomString(20, prefix), domain))
}
