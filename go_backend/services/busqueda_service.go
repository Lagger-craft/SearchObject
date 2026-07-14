package services

import (
	"searchobject/db"
	"searchobject/models"
)

type BusquedaService struct {
	db *db.DB
}

func NewBusquedaService(database *db.DB) *BusquedaService {
	return &BusquedaService{db: database}
}

type ResultadoBusqueda struct {
	Objeto    models.Objeto    `json:"objeto"`
	Caja      models.Caja      `json:"caja"`
	Espacio   models.Espacio   `json:"espacio"`
	Imagen    *models.Imagen   `json:"imagen,omitempty"`
}

func (s *BusquedaService) Buscar(usuarioID, termino string) ([]ResultadoBusqueda, error) {
	objetos, err := s.db.BuscarObjetos(usuarioID, termino)
	if err != nil {
		return nil, err
	}

	resultados := make([]ResultadoBusqueda, 0, len(objetos))
	for _, o := range objetos {
		r := ResultadoBusqueda{Objeto: o}

		caja, err := s.db.ObtenerCaja(o.CajaID)
		if err == nil {
			r.Caja = *caja

			espacio, err := s.db.ObtenerEspacio(caja.EspacioID)
			if err == nil {
				r.Espacio = *espacio
			}
		}

		img, _ := s.db.ObtenerImagenPrincipal(o.ID)
		r.Imagen = img

		resultados = append(resultados, r)
	}

	return resultados, nil
}
