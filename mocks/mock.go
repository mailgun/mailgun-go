package mocks

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
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/events"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

// Server is a Mailgun API mock suitable for testing
type Server struct {
	srv *httptest.Server

	domainIPS        []string
	domainList       []DomainContainer
	exportList       []mtypes.Export
	mailingList      []MailingListContainer
	routeList        []mtypes.Route
	events           []events.Event
	templates        []mtypes.Template
	templateVersions map[string][]mtypes.TemplateVersion
	unsubscribes     []mtypes.Unsubscribe
	complaints       []mtypes.Complaint
	bounces          []mtypes.Bounce
	credentials      []mtypes.Credential
	tags             []mtypes.Tag
	subaccountList   []mtypes.Subaccount
	webhooks         mtypes.WebHooksListResponse
	mutex            sync.Mutex
}

func (ms *Server) DomainIPS() []string {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.domainIPS
}

func (ms *Server) DomainList() []DomainContainer {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.domainList
}

func (ms *Server) ExportList() []mtypes.Export {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.exportList
}

func (ms *Server) MailingList() []MailingListContainer {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.mailingList
}

func (ms *Server) RouteList() []mtypes.Route {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.routeList
}

func (ms *Server) Events() []events.Event {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.events
}

func (ms *Server) Webhooks() mtypes.WebHooksListResponse {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.webhooks
}

func (ms *Server) Templates() []mtypes.Template {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.templates
}

func (ms *Server) Unsubscribes() []mtypes.Unsubscribe {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.unsubscribes
}

func (ms *Server) SubaccountList() []mtypes.Subaccount {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	return ms.subaccountList
}

// NewServer creates a new instance of the mailgun API mock server
func NewServer() *Server {
	ms := Server{}

	// Add all our handlers
	r := chi.NewRouter()

	r.Route("/v3", func(r chi.Router) {
		ms.addIPRoutes(r)
		ms.addExportRoutes(r)
		ms.addMailingListRoutes(r)
		ms.addEventRoutes(r)
		ms.addMessagesRoutes(r)
		ms.addRoutes(r)
		ms.addWebhookRoutes(r)
		ms.addTemplateRoutes(r)
		ms.addTemplateVersionRoutes(r)
		ms.addUnsubscribesRoutes(r)
		ms.addComplaintsRoutes(r)
		ms.addBouncesRoutes(r)
		ms.addCredentialsRoutes(r)
		ms.addTagsRoutes(r)
	})
	r.Route("/v5", func(r chi.Router) {
		ms.addSubaccountRoutes(r)
	})
	ms.addDomainRoutes(r) // mix of v3 and v4
	ms.addValidationRoutes(r)
	ms.addAnalyticsRoutes(r)

	// Start the server
	ms.srv = httptest.NewServer(r)
	return &ms
}

// Stop the server
func (ms *Server) Stop() {
	ms.srv.Close()
}

func (ms *Server) URL() string {
	return ms.srv.URL
}

func toJSON(w http.ResponseWriter, obj any) {
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

func stringToMap(v string) map[string]any {
	if v == "" {
		return nil
	}

	result := make(map[string]any)
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

type okResp struct {
	ID      string `json:"id,omitempty"`
	Message string `json:"message"`
}

func ptr[T any](v T) *T {
	return &v
}
