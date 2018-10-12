package reqpayload

type (
	Register struct {
		Username string `json:"username" form:"username" validate:"required" example:"john_doe"`
		Password string `json:"password" form:"password" validate:"required" example:"secret"`
	}

	Login struct {
		Username string `json:"username" form:"username" validate:"required" example:"john_doe"`
		Password string `json:"password" form:"password" validate:"required" example:"secret"`
	}
)
