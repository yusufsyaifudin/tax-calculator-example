package reqpayload

type (
	// Register is a payload required when user register on this system.
	Register struct {
		Username string `json:"username" form:"username" validate:"required" example:"john_doe"`
		Password string `json:"password" form:"password" validate:"required" example:"secret"`
	}

	// Login is a payload required when user login to this system.
	Login struct {
		Username string `json:"username" form:"username" validate:"required" example:"john_doe"`
		Password string `json:"password" form:"password" validate:"required" example:"secret"`
	}
)
