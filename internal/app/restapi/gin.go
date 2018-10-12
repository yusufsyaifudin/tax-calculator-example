package restapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/model"
)

// WrapGin wraps a Handler and turns it into gin compatible handler
// This method should be called with a fresh ctx
func WrapGin(parent context.Context, h Handler) gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		// create span
		ctx, closer := context.WithTimeout(parent, 10*time.Second)
		defer closer()

		// create request and run the handler
		var req = newGinRequest(gCtx)
		resp := h(ctx, req)

		// get the body first
		body, err := resp.Body()
		if err != nil {
			gCtx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "internal server error",
			})
			return
		}

		// then write header
		for k, v := range resp.Header() {
			for _, h := range v {
				gCtx.Writer.Header().Add(k, h)
			}
		}
		gCtx.Writer.Header().Add("Content-Type", resp.ContentType())

		// the last is writing the body
		gCtx.Writer.WriteHeader(resp.StatusCode())
		gCtx.Writer.Write(body)
	}
}

type ginRequest struct {
	gCtx *gin.Context
	user *model.User
}

func newGinRequest(gCtx *gin.Context) Request {
	return &ginRequest{
		gCtx: gCtx,
	}
}

func (r *ginRequest) Bind(out interface{}) error {
	return r.gCtx.Bind(out)
}

func (r *ginRequest) Header() http.Header {
	return r.gCtx.Request.Header
}

func (r *ginRequest) ContentType() string {
	return r.gCtx.ContentType()
}

func (r *ginRequest) RawRequest() *http.Request {
	return r.gCtx.Request
}

func (r *ginRequest) GetParam(key string) string {
	return r.gCtx.Param(key)
}

func (r *ginRequest) User() *model.User {
	return r.user
}

func (r *ginRequest) SetUser(user *model.User) {
	r.user = user
}
