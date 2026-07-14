package models

type Tag struct {
	ID         string `json:"id"`
	Nombre     string `json:"nombre"`
	NombreNorm string `json:"-"`
	Color      string `json:"color,omitempty"`
	CreatedAt  Time   `json:"created_at"`
}

type ObjetoTag struct {
	ObjetoID string `json:"objeto_id"`
	TagID    string `json:"tag_id"`
}
