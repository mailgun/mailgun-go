package mailgun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/mailgun/errors"
)

var validURL = regexp.MustCompile(`/v[2-5].*`)

type httpRequest struct {
	URL                string
	Parameters         map[string][]string
	Headers            map[string]string
	BasicAuthUser      string
	BasicAuthPassword  string
	Client             *http.Client
	capturedCurlOutput string
}

type httpResponse struct {
	Code int
	Data []byte
}

type payload interface {
	getPayloadBuffer() (*bytes.Buffer, error)
	getContentType() string
	getValues() []keyValuePair
}

type keyValuePair struct {
	key   string
	value string
}

type keyNameRC struct {
	key   string
	name  string
	value io.ReadCloser
}

type keyNameBuff struct {
	key   string
	name  string
	value []byte
}

type formDataPayload struct {
	contentType string
	Values      []keyValuePair
	Files       []keyValuePair
	ReadClosers []keyNameRC
	Buffers     []keyNameBuff
}

type urlEncodedPayload struct {
	Values []keyValuePair
}

type jsonEncodedPayload struct {
	payload interface{}
}

func newHTTPRequest(url string) *httpRequest {
	return &httpRequest{URL: url, Client: http.DefaultClient}
}

func (r *httpRequest) addParameter(name, value string) {
	if r.Parameters == nil {
		r.Parameters = make(map[string][]string)
	}
	r.Parameters[name] = append(r.Parameters[name], value)
}

func (r *httpRequest) setClient(c *http.Client) {
	r.Client = c
}

func (r *httpRequest) setBasicAuth(user, password string) {
	r.BasicAuthUser = user
	r.BasicAuthPassword = password
}

func newJSONEncodedPayload(payload interface{}) *jsonEncodedPayload {
	return &jsonEncodedPayload{payload: payload}
}

func (j *jsonEncodedPayload) getPayloadBuffer() (*bytes.Buffer, error) {
	b, err := json.Marshal(j.payload)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}

func (j *jsonEncodedPayload) getContentType() string {
	return "application/json"
}

func (j *jsonEncodedPayload) getValues() []keyValuePair {
	return nil
}

func newUrlEncodedPayload() *urlEncodedPayload {
	return &urlEncodedPayload{}
}

func (f *urlEncodedPayload) addValue(key, value string) {
	f.Values = append(f.Values, keyValuePair{key: key, value: value})
}

func (f *urlEncodedPayload) getPayloadBuffer() (*bytes.Buffer, error) {
	data := url.Values{}
	for _, keyVal := range f.Values {
		data.Add(keyVal.key, keyVal.value)
	}
	return bytes.NewBufferString(data.Encode()), nil
}

func (f *urlEncodedPayload) getContentType() string {
	return "application/x-www-form-urlencoded"
}

func (f *urlEncodedPayload) getValues() []keyValuePair {
	return f.Values
}

func (r *httpResponse) parseFromJSON(v interface{}) error {
	return json.Unmarshal(r.Data, v)
}

func newFormDataPayload() *formDataPayload {
	return &formDataPayload{}
}

func (f *formDataPayload) getValues() []keyValuePair {
	return f.Values
}

func (f *formDataPayload) addValue(key, value string) {
	f.Values = append(f.Values, keyValuePair{key: key, value: value})
}

func (f *formDataPayload) addFile(key, file string) {
	f.Files = append(f.Files, keyValuePair{key: key, value: file})
}

func (f *formDataPayload) addBuffer(key, file string, buff []byte) {
	f.Buffers = append(f.Buffers, keyNameBuff{key: key, name: file, value: buff})
}

func (f *formDataPayload) addReadCloser(key, name string, rc io.ReadCloser) {
	f.ReadClosers = append(f.ReadClosers, keyNameRC{key: key, name: name, value: rc})
}

func (f *formDataPayload) getPayloadBuffer() (*bytes.Buffer, error) {
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	defer writer.Close()

	for _, keyVal := range f.Values {
		if tmp, err := writer.CreateFormField(keyVal.key); err == nil {
			tmp.Write([]byte(keyVal.value))
		} else {
			return nil, err
		}
	}

	for _, file := range f.Files {
		if tmp, err := writer.CreateFormFile(file.key, path.Base(file.value)); err == nil {
			if fp, err := os.Open(file.value); err == nil {
				defer fp.Close()
				io.Copy(tmp, fp)
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	for _, file := range f.ReadClosers {
		if tmp, err := writer.CreateFormFile(file.key, file.name); err == nil {
			defer file.value.Close()
			io.Copy(tmp, file.value)
		} else {
			return nil, err
		}
	}

	for _, buff := range f.Buffers {
		if tmp, err := writer.CreateFormFile(buff.key, buff.name); err == nil {
			r := bytes.NewReader(buff.value)
			io.Copy(tmp, r)
		} else {
			return nil, err
		}
	}

	f.contentType = writer.FormDataContentType()

	return data, nil
}

func (f *formDataPayload) getContentType() string {
	if f.contentType == "" {
		f.getPayloadBuffer()
	}
	return f.contentType
}

func (r *httpRequest) addHeader(name, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[name] = value
}

func (r *httpRequest) makeGetRequest(ctx context.Context) (*httpResponse, error) {
	return r.makeRequest(ctx, "GET", nil)
}

func (r *httpRequest) makePostRequest(ctx context.Context, payload payload) (*httpResponse, error) {
	return r.makeRequest(ctx, "POST", payload)
}

func (r *httpRequest) makePutRequest(ctx context.Context, payload payload) (*httpResponse, error) {
	return r.makeRequest(ctx, "PUT", payload)
}

func (r *httpRequest) makeDeleteRequest(ctx context.Context) (*httpResponse, error) {
	return r.makeRequest(ctx, "DELETE", nil)
}

func (r *httpRequest) NewRequest(ctx context.Context, method string, payload payload) (*http.Request, error) {
	url, err := r.generateUrlWithParameters()
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if payload != nil {
		if body, err = payload.getPayloadBuffer(); err != nil {
			return nil, err
		}
	} else {
		body = nil
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	if payload != nil && payload.getContentType() != "" {
		req.Header.Add("Content-Type", payload.getContentType())
	}

	if r.BasicAuthUser != "" && r.BasicAuthPassword != "" {
		req.SetBasicAuth(r.BasicAuthUser, r.BasicAuthPassword)
	}

	for header, value := range r.Headers {
		// Special case, override the Host header
		if header == "Host" {
			req.Host = value
			continue
		}
		req.Header.Add(header, value)
	}
	return req, nil
}

func (r *httpRequest) makeRequest(ctx context.Context, method string, payload payload) (*httpResponse, error) {
	req, err := r.NewRequest(ctx, method, payload)
	if err != nil {
		return nil, err
	}

	if Debug {
		if CaptureCurlOutput {
			r.capturedCurlOutput = r.curlString(req, payload)
		} else {
			fmt.Println(r.curlString(req, payload))
		}
	}

	response := httpResponse{}

	resp, err := r.Client.Do(req)
	if resp != nil {
		response.Code = resp.StatusCode
	}
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			if urlErr.Err == io.EOF {
				return nil, errors.Wrap(err, "remote server prematurely closed connection")
			}
		}
		return nil, errors.Wrap(err, "while making http request")
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "while reading response body")
	}

	response.Data = responseBody
	return &response, nil
}

func (r *httpRequest) generateUrlWithParameters() (string, error) {
	url, err := url.Parse(r.URL)
	if err != nil {
		return "", err
	}

	if !validURL.MatchString(url.Path) {
		return "", errors.New(`BaseAPI must end with a /v2, /v3 or /v4; setBaseAPI("https://host/v3")`)
	}

	q := url.Query()
	if r.Parameters != nil && len(r.Parameters) > 0 {
		for name, values := range r.Parameters {
			for _, value := range values {
				q.Add(name, value)
			}
		}
	}
	url.RawQuery = q.Encode()

	return url.String(), nil
}

func (r *httpRequest) curlString(req *http.Request, p payload) string {

	parts := []string{"curl", "-i", "-X", req.Method, req.URL.String()}
	for key, value := range req.Header {
		if key == "Authorization" {
			parts = append(parts, fmt.Sprintf("-H \"%s: %s\"", key, "<redacted>"))
		} else {
			parts = append(parts, fmt.Sprintf("-H \"%s: %s\"", key, value[0]))
		}
	}

	// Special case for Host
	if req.Host != "" {
		parts = append(parts, fmt.Sprintf("-H \"Host: %s\"", req.Host))
	}

	// parts = append(parts, fmt.Sprintf(" --user '%s:%s'", r.BasicAuthUser, r.BasicAuthPassword))

	if p != nil {
		if p.getContentType() == "application/json" {
			b, err := p.getPayloadBuffer()
			if err != nil {
				return "Unable to get payload buffer: " + err.Error()
			}
			parts = append(parts, fmt.Sprintf("--data '%s'", b.String()))
		} else {
			for _, param := range p.getValues() {
				parts = append(parts, fmt.Sprintf(" -F %s='%s'", param.key, param.value))
			}
		}
	}
	return strings.Join(parts, " ")
}
