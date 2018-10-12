package restapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/auth"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/repo/user"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/reqpayload"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/respayload"
	"github.com/yusufsyaifudin/tax-calculator-example/pkg/validator"
)

// this should not be here, we must set it from env variable
const secretKey = "my-secret"

// Register
// @Summary Register new account
// @Description Register new account
// @ID user-register
// @Param user body reqpayload.Register true "user info"
// @Accept  json
// @Produce  json
// @Success 200 {object} respayload.Register
// @Failure 400 {object} respayload.Error
// @Failure 422 {object} respayload.Error
// @Router /register [post]
func register(parent context.Context, req Request) Response {
	form := &reqpayload.Register{}
	err := req.Bind(form)
	if err != nil {
		return newJSONResponse(http.StatusUnprocessableEntity, respayload.Error{
			HttpStatusCode: http.StatusUnprocessableEntity,
			ErrorCode:      respayload.ErrorBindingBodyRequest,
			Message:        fmt.Sprintf("error while binding request body %s", err.Error()),
		})
	}

	// trim head and tail spaces
	form.Username = strings.TrimSpace(form.Username)
	form.Password = strings.TrimSpace(form.Password)

	if errs := validator.Validate(form); errs != nil {
		return newJSONResponse(http.StatusBadRequest, respayload.Error{
			HttpStatusCode: http.StatusBadRequest,
			ErrorCode:      respayload.ErrorGeneralValidationError,
			Message:        errs.String(),
		})
	}

	// generate password using bcrypt
	password, err := auth.HashPassword(form.Password)
	if err != nil {
		return newJSONResponse(http.StatusUnprocessableEntity, respayload.Error{
			HttpStatusCode: http.StatusUnprocessableEntity,
			ErrorCode:      respayload.ErrorCodeUserCantBeCreated,
			Message:        fmt.Sprintf("fail when hashing your password %s", err.Error()),
		})
	}

	User, err := user.Create(parent, form.Username, password)
	if err != nil {
		return newJSONResponse(http.StatusUnprocessableEntity, respayload.Error{
			HttpStatusCode: http.StatusUnprocessableEntity,
			ErrorCode:      respayload.ErrorCodeUserCantBeCreated,
			Message:        fmt.Sprintf("db error when insert %s", err.Error()),
		})
	}

	authToken, err := auth.GenerateJWTToken(parent, secretKey, User.ID)
	if err != nil {
		return newJSONResponse(http.StatusUnprocessableEntity, respayload.Error{
			HttpStatusCode: http.StatusUnprocessableEntity,
			ErrorCode:      respayload.ErrorCodeUserCantBeCreated,
			Message:        fmt.Sprintf("db error when generating auth token %s", err.Error()),
		})
	}

	return newJSONResponse(http.StatusOK, respayload.Register{
		AuthenticationToken: authToken,
		User: &respayload.User{
			ID:       User.ID,
			Username: User.Username,
		},
	})
}

// Login
// @Summary Login account
// @Description Login using username and password
// @ID user-login
// @Param user body reqpayload.Login true "user info"
// @Accept  json
// @Produce  json
// @Success 200 {object} respayload.Login
// @Failure 400 {object} respayload.Error
// @Failure 401 {object} respayload.Error
// @Failure 404 {object} respayload.Error
// @Failure 422 {object} respayload.Error
// @Router /login [post]
func login(parent context.Context, req Request) Response {
	form := &reqpayload.Login{}
	err := req.Bind(form)
	if err != nil {
		return newJSONResponse(http.StatusUnprocessableEntity, respayload.Error{
			HttpStatusCode: http.StatusUnprocessableEntity,
			ErrorCode:      respayload.ErrorBindingBodyRequest,
			Message:        fmt.Sprintf("error while binding request body %s", err.Error()),
		})
	}

	// trim head and tail spaces
	form.Username = strings.TrimSpace(form.Username)
	form.Password = strings.TrimSpace(form.Password)

	if errs := validator.Validate(form); errs != nil {
		return newJSONResponse(http.StatusBadRequest, respayload.Error{
			HttpStatusCode: http.StatusBadRequest,
			ErrorCode:      respayload.ErrorGeneralValidationError,
			Message:        errs.String(),
		})
	}

	User, err := user.FindByUsername(parent, form.Username)
	if User == nil || User.ID == 0 {
		return newJSONResponse(http.StatusNotFound, respayload.Error{
			HttpStatusCode: http.StatusNotFound,
			ErrorCode:      respayload.ErrorCodeUserCantBeFound,
			Message:        "cannot find user",
		})
	}

	if err != nil {
		return newJSONResponse(http.StatusUnprocessableEntity, respayload.Error{
			HttpStatusCode: http.StatusUnprocessableEntity,
			ErrorCode:      respayload.ErrorCodeUserCantBeFound,
			Message:        fmt.Sprintf("db error when find user %s", err.Error()),
		})
	}

	ok := auth.CheckPasswordHash(form.Password, User.Password)
	if !ok {
		return newJSONResponse(http.StatusUnauthorized, respayload.Error{
			HttpStatusCode: http.StatusUnauthorized,
			ErrorCode:      respayload.ErrorCodeUserWrongPassword,
			Message:        "your password is mismatch",
		})
	}

	authToken, err := auth.GenerateJWTToken(parent, secretKey, User.ID)
	if err != nil {
		return newJSONResponse(http.StatusUnprocessableEntity, respayload.Error{
			HttpStatusCode: http.StatusUnprocessableEntity,
			ErrorCode:      respayload.ErrorCodeUserCantBeFound,
			Message:        fmt.Sprintf("db error when generating auth token %s", err.Error()),
		})
	}

	return newJSONResponse(http.StatusOK, respayload.Login{
		AuthenticationToken: authToken,
		User: &respayload.User{
			ID:       User.ID,
			Username: User.Username,
		},
	})
}
