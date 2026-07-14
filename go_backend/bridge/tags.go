package bridge

import (
	"encoding/json"

	apperrors "searchobject/errors"
	"searchobject/models"
	"searchobject/normalize"
)

type crearTagReq struct {
	Nombre string `json:"nombre"`
	Color  string `json:"color,omitempty"`
}

type tagDTO struct {
	ID        string `json:"id"`
	Nombre    string `json:"nombre"`
	Color     string `json:"color,omitempty"`
	CreatedAt models.Time `json:"created_at"`
}

type listarTagsRes struct {
	Tags []tagDTO `json:"tags"`
}

func (app *App) CrearTag(jsonReq string) (string, error) {
	var req crearTagReq
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		return "", apperrors.Wrap(err, apperrors.ErrValidation, "JSON inválido", "bridge.CrearTag")
	}

	if req.Nombre == "" {
		return "", apperrors.New(apperrors.ErrValidation, "el nombre es requerido", "bridge.CrearTag")
	}

	t := &models.Tag{
		Nombre:     req.Nombre,
		NombreNorm: normalize.Nombre(req.Nombre),
		Color:      req.Color,
	}

	if err := app.db.CrearTag(t); err != nil {
		return "", err
	}

	dto := tagDTO{ID: t.ID, Nombre: t.Nombre, Color: t.Color, CreatedAt: t.CreatedAt}
	data, _ := json.Marshal(dto)
	return string(data), nil
}

func (app *App) ListarTags() (string, error) {
	tags, err := app.db.ListarTags()
	if err != nil {
		return "", err
	}

	dtos := make([]tagDTO, 0, len(tags))
	for _, t := range tags {
		dtos = append(dtos, tagDTO{ID: t.ID, Nombre: t.Nombre, Color: t.Color, CreatedAt: t.CreatedAt})
	}

	res := listarTagsRes{Tags: dtos}
	data, _ := json.Marshal(res)
	return string(data), nil
}

func (app *App) AgregarTagAObjeto(jsonReq string) (string, error) {
	var req struct {
		ObjetoID string `json:"objeto_id"`
		TagID    string `json:"tag_id"`
	}
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		return "", apperrors.Wrap(err, apperrors.ErrValidation, "JSON inválido", "bridge.AgregarTagAObjeto")
	}

	if err := app.db.AgregarTagAObjeto(req.ObjetoID, req.TagID); err != nil {
		return "", err
	}

	return `{"ok":true}`, nil
}

func (app *App) QuitarTagDeObjeto(jsonReq string) (string, error) {
	var req struct {
		ObjetoID string `json:"objeto_id"`
		TagID    string `json:"tag_id"`
	}
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		return "", apperrors.Wrap(err, apperrors.ErrValidation, "JSON inválido", "bridge.QuitarTagDeObjeto")
	}

	if err := app.db.QuitarTagDeObjeto(req.ObjetoID, req.TagID); err != nil {
		return "", err
	}

	return `{"ok":true}`, nil
}

func (app *App) TagsDeObjeto(objetoID string) (string, error) {
	tags, err := app.db.TagsDeObjeto(objetoID)
	if err != nil {
		return "", err
	}

	dtos := make([]tagDTO, 0, len(tags))
	for _, t := range tags {
		dtos = append(dtos, tagDTO{ID: t.ID, Nombre: t.Nombre, Color: t.Color, CreatedAt: t.CreatedAt})
	}

	res := listarTagsRes{Tags: dtos}
	data, _ := json.Marshal(res)
	return string(data), nil
}

func (app *App) BuscarPorTag(usuarioID, tagID string) (string, error) {
	objetos, err := app.db.ObjetosPorTag(usuarioID, tagID)
	if err != nil {
		return "", err
	}

	dtos := make([]objetoDTO, 0, len(objetos))
	for _, o := range objetos {
		imgs, _ := app.db.ListarImagenes(o.ID)
		dtos = append(dtos, toObjetoDTO(o, len(imgs) > 0))
	}

	res := listarObjetosRes{Objetos: dtos}
	data, _ := json.Marshal(res)
	return string(data), nil
}
