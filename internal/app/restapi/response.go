package restapi

import (
	"encoding/json"
	"net/http"
)

// Response represents an api response
type Response interface {
	StatusCode() int
	Body() ([]byte, error)
	Header() http.Header
	ContentType() string
}

type jsonResponse struct {
	data       interface{}
	header     http.Header
	statusCode int
}

// newJSONResponse creates a json response
func newJSONResponse(code int, data interface{}) Response {
	return &jsonResponse{
		statusCode: code,
		data:       data,
		header:     http.Header{},
	}
}

func (r *jsonResponse) StatusCode() int {
	return r.statusCode
}

func (r *jsonResponse) Body() ([]byte, error) {
	b, err := json.Marshal(r.data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *jsonResponse) Header() http.Header {
	return r.header
}

func (r *jsonResponse) ContentType() string {
	return ContentTypeJSON
}

type dummyResponse struct {
	statusCode  int
	err         error
	body        []byte
	header      http.Header
	contentType string
}

func (r *dummyResponse) StatusCode() int {
	return r.statusCode
}

func (r *dummyResponse) Body() ([]byte, error) {
	return r.body, r.err
}

func (r *dummyResponse) Header() http.Header {
	return r.header
}

func (r *dummyResponse) ContentType() string {
	return r.contentType
}
