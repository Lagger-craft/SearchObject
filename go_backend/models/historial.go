package models

type Accion string

const (
	AccionCreado      Accion = "creado"
	AccionActualizado Accion = "actualizado"
	AccionEliminado   Accion = "eliminado"
	AccionMovido      Accion = "movido"
)

type Historial struct {
	ID          string `json:"id"`
	EntidadTipo string `json:"entidad_tipo"`
	EntidadID   string `json:"entidad_id"`
	Accion      Accion `json:"accion"`
	Detalle     string `json:"detalle,omitempty"`
	UsuarioID   string `json:"usuario_id"`
	CreatedAt   Time   `json:"created_at"`
}
