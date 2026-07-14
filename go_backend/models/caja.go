package models

type Caja struct {
	ID           string  `json:"id"`
	EspacioID    string  `json:"espacio_id"`
	UsuarioID    string  `json:"usuario_id"`
	Nombre       string  `json:"nombre"`
	NombreNorm   string  `json:"-"`
	Descripcion  string  `json:"descripcion,omitempty"`
	CapacidadMax *int    `json:"capacidad_max,omitempty"`
	CreatedAt    Time    `json:"created_at"`
	UpdatedAt    Time    `json:"updated_at"`
}
