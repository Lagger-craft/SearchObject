package bridge

import (
	"encoding/json"

	"github.com/google/uuid"

	apperrors "searchobject/errors"
	"searchobject/models"
	"searchobject/normalize"
)

type crearCajaReq struct {
	EspacioID    string `json:"espacio_id"`
	UsuarioID    string `json:"usuario_id"`
	Nombre       string `json:"nombre"`
	Descripcion  string `json:"descripcion"`
	CapacidadMax *int   `json:"capacidad_max"`
}

type listarCajasRes struct {
	Cajas []cajaDTO `json:"cajas"`
}

func (app *App) CrearCaja(jsonReq string) (string, error) {
	var req crearCajaReq
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		return "", apperrors.Wrap(err, apperrors.ErrValidation, "JSON inválido", "bridge.CrearCaja")
	}

	if req.Nombre == "" {
		return "", apperrors.New(apperrors.ErrValidation, "el nombre de la caja no puede estar vacío", "bridge.CrearCaja")
	}

	norm := normalize.Nombre(req.Nombre)

	caja := &models.Caja{
		ID:           uuid.New().String(),
		EspacioID:    req.EspacioID,
		UsuarioID:    req.UsuarioID,
		Nombre:       req.Nombre,
		NombreNorm:   norm,
		Descripcion:  req.Descripcion,
		CapacidadMax: req.CapacidadMax,
		CreatedAt:    models.Now(),
		UpdatedAt:    models.Now(),
	}

	if err := app.db.CrearCaja(caja); err != nil {
		return "", err
	}

	data, _ := json.Marshal(toCajaDTO(*caja))
	return string(data), nil
}

func (app *App) ListarCajas(espacioID string) (string, error) {
	cajas, err := app.db.ListarCajas(espacioID)
	if err != nil {
		return "", err
	}

	dtos := make([]cajaDTO, 0, len(cajas))
	for _, c := range cajas {
		dtos = append(dtos, toCajaDTO(c))
	}

	res := listarCajasRes{Cajas: dtos}
	data, _ := json.Marshal(res)
	return string(data), nil
}

func (app *App) ObtenerCaja(id string) (string, error) {
	c, err := app.db.ObtenerCaja(id)
	if err != nil {
		return "", err
	}
	data, _ := json.Marshal(toCajaDTO(*c))
	return string(data), nil
}

func (app *App) EliminarCaja(id string) (string, error) {
	if err := app.db.EliminarCaja(id); err != nil {
		return "", err
	}
	return `{"ok":true}`, nil
}
