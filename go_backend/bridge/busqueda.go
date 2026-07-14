package bridge

import (
	"encoding/json"

	"searchobject/normalize"
)

type resultadoBusquedaDTO struct {
	Objeto  objetoDTO  `json:"objeto"`
	Espacio espacioDTO `json:"espacio"`
	Caja    cajaDTO    `json:"caja"`
	Imagen  *imagenDTO `json:"imagen,omitempty"`
}

type resultadosBusquedaRes struct {
	Resultados []resultadoBusquedaDTO `json:"resultados"`
}

func (app *App) Buscar(usuarioID, termino string) (string, error) {
	norm := normalize.Busqueda(termino)
	if norm == "" {
		return `{"resultados":[]}`, nil
	}

	resultados, err := app.busqueda.Buscar(usuarioID, norm)
	if err != nil {
		return "", err
	}

	dtos := make([]resultadoBusquedaDTO, 0, len(resultados))
	for _, r := range resultados {
		dto := resultadoBusquedaDTO{
			Objeto:  toObjetoDTO(r.Objeto, r.Imagen != nil),
			Espacio: toEspacioDTO(r.Espacio, false),
			Caja:    toCajaDTO(r.Caja),
		}
		if r.Imagen != nil {
			img := toImagenDTO(*r.Imagen)
			dto.Imagen = &img
		}
		dtos = append(dtos, dto)
	}

	res := resultadosBusquedaRes{Resultados: dtos}
	data, _ := json.Marshal(res)
	return string(data), nil
}
