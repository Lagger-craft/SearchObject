package models

type AlertaTipo string

const (
	AlertaFaltante        AlertaTipo = "faltante"
	AlertaStockBajo       AlertaTipo = "stock_bajo"
	AlertaPrestamoVencido AlertaTipo = "prestamo_vencido"
	AlertaObjetoPerdido   AlertaTipo = "objeto_perdido"
	AlertaCajaSaturada    AlertaTipo = "caja_saturada"
)

type Alerta struct {
	ID          string     `json:"id"`
	Tipo        AlertaTipo `json:"tipo"`
	EntidadTipo string     `json:"entidad_tipo"`
	EntidadID   string     `json:"entidad_id"`
	Mensaje     string     `json:"mensaje"`
	Leida       bool       `json:"leida"`
	CreatedAt   Time       `json:"created_at"`
	ResueltaAt  *Time      `json:"resuelta_at,omitempty"`
}
