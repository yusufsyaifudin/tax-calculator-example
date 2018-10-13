package model

import "time"

// Tax represent data structure on database in table taxes.
type Tax struct {
	ID        int64
	UserID    int64
	Name      string
	TaxCode   TaxCode
	Price     int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetTaxCodeString is a helper to return the name of tax category in string (instead using integer code that we save in db).
func (t *Tax) GetTaxCodeString() string {
	switch t.TaxCode {
	case TaxCodeFood:
		return "Food & Beverage"
	case TaxCodeTobacco:
		return "Tobacco"
	case TaxCodeEntertainment:
		return "Entertainment"
	default:
		return "unknown"
	}
}

// IsRefundable returns whether this tax type is refundable or not.
func (t *Tax) IsRefundable() bool {
	switch t.TaxCode {
	case TaxCodeFood:
		return true
	case TaxCodeTobacco:
		return false
	case TaxCodeEntertainment:
		return false
	default:
		return false
	}
}

// GetTaxValue calculates the tax value based on system specification.
func (t *Tax) GetTaxValue() float64 {
	switch t.TaxCode {
	case TaxCodeFood:
		// 10% of Price
		return (float64(10) / float64(100)) * float64(t.Price)
	case TaxCodeTobacco:
		// 10 + (2% of Price )
		return float64(10) + float64((float64(2)/float64(100))*float64(t.Price))
	case TaxCodeEntertainment:
		// Price >= 100: 1% of ( Price - 100)
		if t.Price >= 100 {
			return float64(1) / float64(100) * (float64(t.Price) - float64(100))
		}

		// 0 < Price < 100: tax-free
		return 0
	default:
		return 0
	}
}

// GetAmount will returns the total amount of this tax item.
func (t *Tax) GetAmount() float64 {
	return float64(t.Price) + t.GetTaxValue()
}
