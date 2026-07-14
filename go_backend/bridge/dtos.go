package bridge

import "searchobject/models"

type cajaDTO struct {
	ID           string  `json:"id"`
	EspacioID    string  `json:"espacio_id"`
	Nombre       string  `json:"nombre"`
	Descripcion  string  `json:"descripcion"`
	CapacidadMax *int    `json:"capacidad_max,omitempty"`
}

func toCajaDTO(c models.Caja) cajaDTO {
	return cajaDTO{
		ID:           c.ID,
		EspacioID:    c.EspacioID,
		Nombre:       c.Nombre,
		Descripcion:  c.Descripcion,
		CapacidadMax: c.CapacidadMax,
	}
}
