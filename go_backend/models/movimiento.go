package models

type MovimientoTipo string

const (
	MovEntrada  MovimientoTipo = "entrada"
	MovSalida   MovimientoTipo = "salida"
	MovTraslado MovimientoTipo = "traslado"
	MovPrestamo MovimientoTipo = "prestamo"
)

type Movimiento struct {
	ID          string         `json:"id"`
	ObjetoID    string         `json:"objeto_id"`
	DesdeCajaID *string        `json:"desde_caja_id,omitempty"`
	HaciaCajaID *string        `json:"hacia_caja_id,omitempty"`
	Tipo        MovimientoTipo `json:"tipo"`
	Nota        string         `json:"nota,omitempty"`
	Fecha       Time           `json:"fecha"`
	CreatedAt   Time           `json:"created_at"`
}
