package respayload

// Login is the response model when user success to login.
// I separate this model with Register model even it looks similar to ensure that
// if we need to add or remove some property here, it don't affecting register response.
type Login struct {
	AuthenticationToken string `json:"authentication_token" example:"abc"`
	User                *User  `json:"user"`
}

// Register is the response model when user success to register.
// I separate this model with Login model even it looks similar to ensure that
// if we need to add or remove some property here, it don't affecting login response.
type Register struct {
	AuthenticationToken string `json:"authentication_token" example:"abc"`
	User                *User  `json:"user"`
}

// User is the entity model which the user entity should look in http response.
// This to make every User object is consistent in every response, since I know that it is painful in the client side,
// if we return inconsistent structure for the same object.
// I mean, in endpoint A, User model only return id, which in endpoint B, User model return id and username. It's hard to parse.
// So, this struct is intented to avoid that inconsistency.
type User struct {
	ID       int64  `json:"id" example:"1"`
	Username string `json:"username" example:"john_doe"`
}
