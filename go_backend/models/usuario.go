package models

type Usuario struct {
	ID        string `json:"id"`
	Nombre    string `json:"nombre"`
	Email     string `json:"email,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	CreatedAt Time   `json:"created_at"`
	UpdatedAt Time   `json:"updated_at"`
}
