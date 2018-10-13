package restapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/model"
)

// Request represents an api request
type Request interface {
	ContentType() string
	Bind(out interface{}) error
	RawRequest() *http.Request
	GetParam(string) string

	/* helper methods */
	User() *model.User
	SetUser(*model.User)
}

// DummyRequest is for testing purpose. So instead using gin context, it will use http.Request.
type DummyRequest struct {
	encodedBody []byte
	req         *http.Request

	user *model.User
}

// NewDummyRequest creates a new dummy request. This implements the Request interface.
func NewDummyRequest() Request {
	req, _ := http.NewRequest("GET", "/", nil)
	return &DummyRequest{
		req: req,
	}
}

// ContentType will return the content type of the request.
func (r *DummyRequest) ContentType() string {
	return r.req.Header.Get("content-type")
}

// Bind will bind the request parameter (query, post form or raw body) into out variable, which can be struct.
func (r *DummyRequest) Bind(out interface{}) error {
	err := binding.Default(r.req.Method, r.ContentType()).Bind(r.req, out)
	r.setReqBody()
	return err
}

// RawRequest will returns the raw request payload, so we gan get the header or etc here.
func (r *DummyRequest) RawRequest() *http.Request {
	return r.req
}

// GetParam will get url parameter, for example on URL /api/:user_name, we can get user_name value by calling
// GetParam("user_name").
func (r *DummyRequest) GetParam(key string) string {
	return r.req.PostFormValue(key)
}

// User get current user of this request.
func (r *DummyRequest) User() *model.User {
	return r.user
}

// SetUser sets the current user based on authentication token. This usually set in middleware.
func (r *DummyRequest) SetUser(user *model.User) {
	r.user = user
}

// setReqBody set the request body so it can be read using
func (r *DummyRequest) setReqBody() {
	r.req.Body = ioutil.NopCloser(bytes.NewBuffer(r.encodedBody))
}

// AddPOSTParam will parse the encoded body in method post into current body.
func (r *DummyRequest) AddPOSTParam(key, val string) *DummyRequest {
	if len(r.encodedBody) > 0 {
		r.encodedBody = []byte(fmt.Sprintf("%s&%s=%s", string(r.encodedBody), key, val))
	} else {
		r.encodedBody = []byte(fmt.Sprintf("%s=%s", key, val))
	}

	r.req.Method = "POST"
	r.req.Header.Set("content-type", ContentTypePostForm)
	r.setReqBody()
	return r
}

// SetMethod set the HTTP method on this request.
func (r *DummyRequest) SetMethod(method string) *DummyRequest {
	r.req.Method = method
	return r
}

// SetContentType set content type of this request.
func (r *DummyRequest) SetContentType(contentType string) *DummyRequest {
	r.req.Header.Set("content-type", contentType)
	return r
}

// AddHeader add header on this request.
func (r *DummyRequest) AddHeader(key, val string) *DummyRequest {
	r.req.Header.Add(key, val)
	return r
}

// SetJSONBody set the raw body json request.
func (r *DummyRequest) SetJSONBody(p interface{}) *DummyRequest {
	c, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	r.encodedBody = c
	r.req.Method = "POST"
	r.req.Header.Set("content-type", ContentTypeJSON)
	r.setReqBody()
	return r
}
