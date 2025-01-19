package model

type LoginRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type AutheticationResponse struct {
	AccessToken string `json:"access_token"`
}
