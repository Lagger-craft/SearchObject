package models

type Objeto struct {
	ID            string   `json:"id"`
	CajaID        string   `json:"caja_id"`
	UsuarioID     string   `json:"usuario_id"`
	Nombre        string   `json:"nombre"`
	NombreNorm    string   `json:"-"`
	Descripcion   string   `json:"descripcion,omitempty"`
	Cantidad      int      `json:"cantidad"`
	EsInsumo      bool     `json:"es_insumo"`
	ValorEstimado *float64 `json:"valor_estimado,omitempty"`
	CreatedAt     Time     `json:"created_at"`
	UpdatedAt     Time     `json:"updated_at"`
}
