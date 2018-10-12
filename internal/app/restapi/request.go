package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

type DummyRequest struct {
	encodedBody []byte
	req         *http.Request

	user *model.User
}

// NewDummyRequest creates a new dummy request
func NewDummyRequest() *DummyRequest {
	req, _ := http.NewRequest("GET", "/", nil)
	return &DummyRequest{
		req: req,
	}
}

func (r *DummyRequest) ContentType() string {
	return r.req.Header.Get("content-type")
}

func (r *DummyRequest) Bind(out interface{}) error {
	err := binding.Default(r.req.Method, r.ContentType()).Bind(r.req, out)
	r.setReqBody()
	return err
}

func (r *DummyRequest) RawRequest() *http.Request {
	return r.req
}

func (r *DummyRequest) GetParam(key string) string {
	return r.req.PostFormValue(key)
}

func (r *DummyRequest) User() *model.User {
	return r.user
}
func (r *DummyRequest) SetUser(user *model.User) {
	r.user = user
}

func (r *DummyRequest) setReqBody() {
	r.req.Body = ioutil.NopCloser(bytes.NewBuffer(r.encodedBody))
}

/* helper method to construct a request */
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

func (r *DummyRequest) SetMethod(method string) *DummyRequest {
	r.req.Method = method
	return r
}

func (r *DummyRequest) SetContentType(contentType string) *DummyRequest {
	r.req.Header.Set("content-type", contentType)
	return r
}

func (r *DummyRequest) AddHeader(key, val string) *DummyRequest {
	r.req.Header.Add(key, val)
	return r
}

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
