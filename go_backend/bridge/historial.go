package bridge

import (
	"encoding/json"

	apperrors "searchobject/errors"
	"searchobject/models"
)

type historialDTO struct {
	ID          string      `json:"id"`
	EntidadTipo string      `json:"entidad_tipo"`
	EntidadID   string      `json:"entidad_id"`
	Accion      models.Accion `json:"accion"`
	Detalle     string      `json:"detalle,omitempty"`
	UsuarioID   string      `json:"usuario_id"`
	CreatedAt   models.Time `json:"created_at"`
}

type listarHistorialRes struct {
	Historial []historialDTO `json:"historial"`
}

func (app *App) RegistrarHistorial(jsonReq string) (string, error) {
	var req struct {
		EntidadTipo string `json:"entidad_tipo"`
		EntidadID   string `json:"entidad_id"`
		Accion      string `json:"accion"`
		Detalle     string `json:"detalle,omitempty"`
		UsuarioID   string `json:"usuario_id"`
	}
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		return "", apperrors.Wrap(err, apperrors.ErrValidation, "JSON inválido", "bridge.RegistrarHistorial")
	}

	h := &models.Historial{
		EntidadTipo: req.EntidadTipo,
		EntidadID:   req.EntidadID,
		Accion:      models.Accion(req.Accion),
		Detalle:     req.Detalle,
		UsuarioID:   req.UsuarioID,
	}

	if err := app.db.CrearHistorial(h); err != nil {
		return "", err
	}

	dto := historialDTO{
		ID:          h.ID,
		EntidadTipo: h.EntidadTipo,
		EntidadID:   h.EntidadID,
		Accion:      h.Accion,
		Detalle:     h.Detalle,
		UsuarioID:   h.UsuarioID,
		CreatedAt:   h.CreatedAt,
	}
	data, _ := json.Marshal(dto)
	return string(data), nil
}

func (app *App) ListarHistorial(usuarioID string) (string, error) {
	historial, err := app.db.ListarHistorial(usuarioID, 50)
	if err != nil {
		return "", err
	}

	dtos := make([]historialDTO, 0, len(historial))
	for _, h := range historial {
		dtos = append(dtos, historialDTO{
			ID:          h.ID,
			EntidadTipo: h.EntidadTipo,
			EntidadID:   h.EntidadID,
			Accion:      h.Accion,
			Detalle:     h.Detalle,
			UsuarioID:   h.UsuarioID,
			CreatedAt:   h.CreatedAt,
		})
	}

	res := listarHistorialRes{Historial: dtos}
	data, _ := json.Marshal(res)
	return string(data), nil
}

func (app *App) HistorialPorEntidad(usuarioID, entidadTipo, entidadID string) (string, error) {
	historial, err := app.db.HistorialPorEntidad(usuarioID, entidadTipo, entidadID)
	if err != nil {
		return "", err
	}

	dtos := make([]historialDTO, 0, len(historial))
	for _, h := range historial {
		dtos = append(dtos, historialDTO{
			ID:          h.ID,
			EntidadTipo: h.EntidadTipo,
			EntidadID:   h.EntidadID,
			Accion:      h.Accion,
			Detalle:     h.Detalle,
			UsuarioID:   h.UsuarioID,
			CreatedAt:   h.CreatedAt,
		})
	}

	res := listarHistorialRes{Historial: dtos}
	data, _ := json.Marshal(res)
	return string(data), nil
}
