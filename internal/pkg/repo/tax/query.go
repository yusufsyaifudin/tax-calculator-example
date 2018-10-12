package tax

import (
	"context"

	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/conn"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/model"
)

// Create will insert new tax related to the specific user id.
func Create(parent context.Context, userID int64, name string, code int, price int64) (Tax *model.Tax, err error) {
	Tax = &model.Tax{}
	err = conn.GetDBConnection().Writer().Query(parent, Tax, sqlInsertTax, userID, name, code, price)
	return
}

// GetTaxesByUserID get taxes by user ID.
func GetTaxesByUserID(parent context.Context, userID int64) (Taxes []*model.Tax, err error) {
	Taxes = []*model.Tax{}
	err = conn.GetDBConnection().Reader().Query(parent, &Taxes, sqlGetTaxesByUserId, userID)
	return
}
