package restapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/repo/tax"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/reqpayload"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/respayload"
	"github.com/yusufsyaifudin/tax-calculator-example/pkg/validator"
)

// Create new tax
// @Summary Add tax record to your account
// @Description Add tax record to your account
// @ID create-tax
//
// @Param Authentication-Token header string true "Authentication-Token your-token"
// @Param tax body reqpayload.CreateNewTax true "tax info"
// @Accept  json
// @Produce  json
// @Success 200 {object} respayload.Tax
// @Failure 400 {object} respayload.Error
// @Failure 422 {object} respayload.Error
// @Router /tax [post]
func createNewTax(parent context.Context, req Request) Response {
	form := &reqpayload.CreateNewTax{}
	err := req.Bind(form)
	if err != nil {
		return newJSONResponse(http.StatusUnprocessableEntity, respayload.Error{
			HttpStatusCode: http.StatusUnprocessableEntity,
			ErrorCode:      respayload.ErrorBindingBodyRequest,
			Message:        fmt.Sprintf("error while binding request body %s", err.Error()),
		})
	}

	// trim spaces
	form.Name = strings.TrimSpace(form.Name)

	if errs := validator.Validate(form); errs != nil {
		return newJSONResponse(http.StatusBadRequest, respayload.Error{
			HttpStatusCode: http.StatusBadRequest,
			ErrorCode:      respayload.ErrorGeneralValidationError,
			Message:        errs.String(),
		})
	}

	// try inserting new tax to DB
	Tax, err := tax.Create(parent, req.User().ID, form.Name, form.TaxCode, form.Price)
	if err != nil {
		return newJSONResponse(http.StatusUnprocessableEntity, respayload.Error{
			HttpStatusCode: http.StatusUnprocessableEntity,
			ErrorCode:      respayload.ErrorCodeTaxCantBeCreated,
			Message:        fmt.Sprintf("db error when insert %s", err.Error()),
		})
	}

	return newJSONResponse(http.StatusOK, respayload.Tax{
		Name:       Tax.Name,
		TaxCode:    int(Tax.TaxCode),
		Type:       Tax.GetTaxCodeString(),
		Price:      Tax.Price,
		Tax:        fmt.Sprintf("%2f", Tax.GetTaxValue()),
		Amount:     fmt.Sprintf("%2f", Tax.GetAmount()),
		Refundable: Tax.IsRefundable(),
	})
}

// TODO: sorry for long inline description, swag doesn't support multi-line description yet. https://github.com/swaggo/swag/issues/191
// Get taxes related to current user
// @Summary Get taxes related to current user
// @Description Get taxes related to current user. Tax calculation is based on following calculation rule: 1. Food and Beverage: 10% of Price, for example if the price is 1000 then the tax is 100, hence the amount is 1100. 2. Tobacco: 10 + (2% of Price), for example if the price is 1000 then the tax is 10 + (2% * 1000) = 10 + 20 = 30, hence the amount is 1030. 3. Entertainment: if the price is equal or more than 100 is 1% of (Price - 100), otherwise is free. For instance, if the price is 150, then the tax is 1% * (150-100) = 1% * 50 = 0.5, hence the final amount is 150.5.
//
// @ID get-taxes
//
// @Param Authentication-Token header string true "Authentication-Token your-token"
// @Accept  json
// @Produce  json
// @Success 200 {object} respayload.TaxesForCurrentUser
// @Failure 400 {object} respayload.Error
// @Failure 422 {object} respayload.Error
// @Router /tax [get]
func getTaxes(parent context.Context, req Request) Response {
	Taxes, err := tax.GetTaxesByUserID(parent, req.User().ID)
	if err != nil {
		return newJSONResponse(http.StatusUnprocessableEntity, respayload.Error{
			HttpStatusCode: http.StatusUnprocessableEntity,
			ErrorCode:      respayload.ErrorCodeTaxDBError,
			Message:        fmt.Sprintf("db error when insert %s", err.Error()),
		})
	}

	priceSubTotal := int64(0)
	taxSubTotal := float64(0)
	grandTotal := float64(0)

	var taxesResponse []respayload.Tax
	for _, Tax := range Taxes {
		priceSubTotal += Tax.Price
		taxSubTotal += float64(Tax.GetTaxValue())
		grandTotal += Tax.GetAmount()

		taxesResponse = append(taxesResponse, respayload.Tax{
			Name:       Tax.Name,
			TaxCode:    int(Tax.TaxCode),
			Type:       Tax.GetTaxCodeString(),
			Price:      Tax.Price,
			Tax:        fmt.Sprintf("%2f", Tax.GetTaxValue()),
			Amount:     fmt.Sprintf("%2f", Tax.GetAmount()),
			Refundable: Tax.IsRefundable(),
		})
	}

	return newJSONResponse(http.StatusOK, respayload.TaxesForCurrentUser{
		PriceSubTotal: priceSubTotal,
		TaxSubTotal:   fmt.Sprintf("%2f", taxSubTotal),
		GrandTotal:    fmt.Sprintf("%2f", grandTotal),
		Taxes:         taxesResponse,
	})
}
