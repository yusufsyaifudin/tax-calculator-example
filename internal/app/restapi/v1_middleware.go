package restapi

import (
	"context"
	"net/http"
	"strings"

	"fmt"

	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/auth"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/repo/user"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/respayload"
)

func middlewareAuthTokenCheck(next Handler) Handler {
	return func(parent context.Context, req Request) Response {
		accessToken := req.RawRequest().Header.Get("Authentication-Token")
		accessToken = strings.TrimSpace(accessToken)

		if accessToken == "" {
			return newJSONResponse(http.StatusUnauthorized, respayload.Error{
				HttpStatusCode: http.StatusUnauthorized,
				ErrorCode:      respayload.ErrorCodeUserWrongAuthToken,
				Message:        "missing auth token",
			})
		}

		userID, err := auth.ValidateJWTToken(parent, secretKey, accessToken)
		if err != nil {
			return newJSONResponse(http.StatusUnauthorized, respayload.Error{
				HttpStatusCode: http.StatusUnauthorized,
				ErrorCode:      respayload.ErrorCodeUserWrongAuthToken,
				Message:        fmt.Sprintf("wrong auth token, %s", err.Error()),
			})
		}

		User, err := user.FindByID(parent, userID)
		if User == nil || User.ID == 0 {
			return newJSONResponse(http.StatusUnauthorized, respayload.Error{
				HttpStatusCode: http.StatusUnauthorized,
				ErrorCode:      respayload.ErrorCodeUserCantBeFound,
				Message:        "related user is not exist in db anymore",
			})
		}

		if err != nil {
			return newJSONResponse(http.StatusUnauthorized, respayload.Error{
				HttpStatusCode: http.StatusUnauthorized,
				ErrorCode:      respayload.ErrorCodeUserCantBeFound,
				Message:        fmt.Sprintf("db error when get user from token, %s", err.Error()),
			})
		}

		// set user to context, so it can be get from handler
		req.SetUser(User)

		// run the wrapped handler
		return next(parent, req)
	}
}
