package respayload

type Login struct {
	AuthenticationToken string `json:"authentication_token" example:"abc"`
	User                *User  `json:"user"`
}

type Register struct {
	AuthenticationToken string `json:"authentication_token" example:"abc"`
	User                *User  `json:"user"`
}

type User struct {
	ID       int64  `json:"id" example:"1"`
	Username string `json:"username" example:"john_doe"`
}
