package models

type Espacio struct {
	ID          string  `json:"id"`
	UsuarioID   string  `json:"usuario_id"`
	Nombre      string  `json:"nombre"`
	NombreNorm  string  `json:"-"`
	Descripcion string  `json:"descripcion,omitempty"`
	PadreID     *string `json:"padre_id,omitempty"`
	CreatedAt   Time    `json:"created_at"`
	UpdatedAt   Time    `json:"updated_at"`
}
