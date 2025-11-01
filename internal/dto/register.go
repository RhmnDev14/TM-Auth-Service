package dto

type Register struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Login struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
