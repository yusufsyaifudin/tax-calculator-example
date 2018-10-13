package model

// TaxCode is a abstraction type of the Tax code, so we can use constant instead pass an integer in every get and set.
// This also leverage code readability.
type TaxCode int

const (
	// TaxCodeFood is a tax code for Food and Beverage type.
	TaxCodeFood TaxCode = 1

	// TaxCodeTobacco is a tax code for Tobacco type.
	TaxCodeTobacco TaxCode = 2

	// TaxCodeEntertainment is a tax code for Entertainment type.
	TaxCodeEntertainment TaxCode = 3
)
