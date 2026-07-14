package models

type Imagen struct {
	ID          string   `json:"id"`
	ObjetoID    string   `json:"objeto_id"`
	Path        string   `json:"path"`
	ThumbPath   string   `json:"thumb_path"`
	AreaX       *float64 `json:"area_x,omitempty"`
	AreaY       *float64 `json:"area_y,omitempty"`
	AreaW       *float64 `json:"area_w,omitempty"`
	AreaH       *float64 `json:"area_h,omitempty"`
	EsPrincipal bool     `json:"es_principal"`
	CreatedAt   Time     `json:"created_at"`
}
