package dto

type RegisterResponse struct {
	ID      uint64 `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Message string `json:"message"`
}
