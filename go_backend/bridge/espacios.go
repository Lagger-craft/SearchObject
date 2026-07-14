package bridge

import (
	"encoding/json"

	apperrors "searchobject/errors"
	"searchobject/models"
)

type crearEspacioReq struct {
	UsuarioID   string  `json:"usuario_id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	PadreID     *string `json:"padre_id"`
}

type listarEspaciosRes struct {
	Espacios []espacioDTO `json:"espacios"`
}

type espacioDTO struct {
	ID          string  `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	PadreID     *string `json:"padre_id,omitempty"`
	TieneHijos  bool    `json:"tiene_hijos"`
}

func (app *App) CrearEspacio(jsonReq string) (string, error) {
	var req crearEspacioReq
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		return "", apperrors.Wrap(err, apperrors.ErrValidation, "JSON inválido", "bridge.CrearEspacio")
	}

	e, err := app.espacios.Crear(req.UsuarioID, req.Nombre, req.Descripcion, req.PadreID)
	if err != nil {
		return "", err
	}

	dto := toEspacioDTO(*e, false)
	data, _ := json.Marshal(dto)
	return string(data), nil
}

func (app *App) ListarEspacios(usuarioID string) (string, error) {
	espacios, err := app.espacios.Listar(usuarioID)
	if err != nil {
		return "", err
	}

	dtos := make([]espacioDTO, 0, len(espacios))
	for _, e := range espacios {
		hijos, _ := app.db.ListarEspacios(usuarioID)
		h := hijosPorPadre(hijos, e.ID)
		dtos = append(dtos, toEspacioDTO(e, len(h) > 0))
	}

	res := listarEspaciosRes{Espacios: dtos}
	data, _ := json.Marshal(res)
	return string(data), nil
}

func (app *App) ObtenerEspacio(id string) (string, error) {
	e, err := app.espacios.Obtener(id)
	if err != nil {
		return "", err
	}
	dto := toEspacioDTO(*e, false)
	data, _ := json.Marshal(dto)
	return string(data), nil
}

func (app *App) ActualizarEspacio(jsonReq string) (string, error) {
	var req struct {
		ID          string  `json:"id"`
		Nombre      string  `json:"nombre"`
		Descripcion string  `json:"descripcion"`
		PadreID     *string `json:"padre_id"`
	}
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		return "", apperrors.Wrap(err, apperrors.ErrValidation, "JSON inválido", "bridge.ActualizarEspacio")
	}

	e, err := app.espacios.Actualizar(req.ID, req.Nombre, req.Descripcion, req.PadreID)
	if err != nil {
		return "", err
	}
	dto := toEspacioDTO(*e, false)
	data, _ := json.Marshal(dto)
	return string(data), nil
}

func (app *App) EliminarEspacio(id string) (string, error) {
	if err := app.espacios.Eliminar(id); err != nil {
		return "", err
	}
	return `{"ok":true}`, nil
}

func toEspacioDTO(e models.Espacio, tieneHijos bool) espacioDTO {
	return espacioDTO{
		ID:          e.ID,
		Nombre:      e.Nombre,
		Descripcion: e.Descripcion,
		PadreID:     e.PadreID,
		TieneHijos:  tieneHijos,
	}
}

func hijosPorPadre(espacios []models.Espacio, padreID string) []models.Espacio {
	var hijos []models.Espacio
	for _, e := range espacios {
		if e.PadreID != nil && *e.PadreID == padreID {
			hijos = append(hijos, e)
		}
	}
	return hijos
}
