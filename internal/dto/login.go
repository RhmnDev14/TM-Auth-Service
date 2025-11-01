package dto

type UserData struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	User  UserData `json:"user"`
}
