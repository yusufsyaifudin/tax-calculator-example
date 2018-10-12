package respayload

type Tax struct {
	Name       string `json:"name" example:"Big Mac"`
	TaxCode    int    `json:"tax_code" exaple:"1"`
	Type       string `json:"type" example:"Food and Beverage"`
	Price      int64  `json:"price" example:"1000"`
	Tax        string `json:"tax" example:"100"`
	Amount     string `json:"amount" example:"1100"`
	Refundable bool   `json:"refundable" example:"false"`
}

type TaxesForCurrentUser struct {
	PriceSubTotal int64  `json:"price_sub_total" example:"2150"`
	TaxSubTotal   string `json:"tax_sub_total" example:"120.5"`
	GrandTotal    string `json:"grand_total" example:"2270.5"`
	Taxes         []Tax  `json:"taxes"`
}
