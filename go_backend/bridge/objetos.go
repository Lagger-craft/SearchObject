package bridge

import (
	"encoding/json"

	apperrors "searchobject/errors"
	"searchobject/models"
	"searchobject/normalize"
	"searchobject/services"
)

type crearObjetoReq struct {
	CajaID      string   `json:"caja_id"`
	UsuarioID   string   `json:"usuario_id"`
	Nombre      string   `json:"nombre"`
	Descripcion string   `json:"descripcion"`
	Cantidad    int      `json:"cantidad"`
	EsInsumo    bool     `json:"es_insumo"`
	Valor       *float64 `json:"valor"`
}

type listarObjetosRes struct {
	Objetos []objetoDTO `json:"objetos"`
}

type objetoDTO struct {
	ID          string   `json:"id"`
	CajaID      string   `json:"caja_id"`
	Nombre      string   `json:"nombre"`
	Descripcion string   `json:"descripcion"`
	Cantidad    int      `json:"cantidad"`
	EsInsumo    bool     `json:"es_insumo"`
	Valor       *float64 `json:"valor,omitempty"`
	TieneFotos  bool     `json:"tiene_fotos"`
}

func (app *App) CrearObjeto(jsonReq string) (string, error) {
	var req crearObjetoReq
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		return "", apperrors.Wrap(err, apperrors.ErrValidation, "JSON inválido", "bridge.CrearObjeto")
	}

	o, err := app.objetos.Crear(services.CrearObjetoReq{
		CajaID:        req.CajaID,
		UsuarioID:     req.UsuarioID,
		Nombre:        req.Nombre,
		Descripcion:   req.Descripcion,
		Cantidad:      req.Cantidad,
		EsInsumo:      req.EsInsumo,
		ValorEstimado: req.Valor,
	})
	if err != nil {
		return "", err
	}

	dto := toObjetoDTO(*o, false)
	data, _ := json.Marshal(dto)
	return string(data), nil
}

func (app *App) ListarObjetos(cajaID string) (string, error) {
	objetos, err := app.objetos.Listar(cajaID)
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

func (app *App) ObtenerObjeto(id string) (string, error) {
	o, err := app.objetos.Obtener(id)
	if err != nil {
		return "", err
	}
	imgs, _ := app.db.ListarImagenes(o.ID)
	dto := toObjetoDTO(*o, len(imgs) > 0)
	data, _ := json.Marshal(dto)
	return string(data), nil
}

func (app *App) MoverObjeto(jsonReq string) (string, error) {
	var req struct {
		ID     string `json:"id"`
		CajaID string `json:"caja_id"`
		Nota   string `json:"nota"`
	}
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		return "", apperrors.Wrap(err, apperrors.ErrValidation, "JSON inválido", "bridge.MoverObjeto")
	}

	o, err := app.objetos.Mover(req.ID, req.CajaID, req.Nota)
	if err != nil {
		return "", err
	}
	dto := toObjetoDTO(*o, false)
	data, _ := json.Marshal(dto)
	return string(data), nil
}

func (app *App) EliminarObjeto(id string) (string, error) {
	if err := app.objetos.Eliminar(id); err != nil {
		return "", err
	}
	return `{"ok":true}`, nil
}

func (app *App) BuscarObjetos(usuarioID, termino string) (string, error) {
	norm := normalize.Busqueda(termino)
	if norm == "" {
		return `{"objetos":[]}`, nil
	}

	objetos, err := app.objetos.Buscar(usuarioID, norm)
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

func toObjetoDTO(o models.Objeto, tieneFotos bool) objetoDTO {
	return objetoDTO{
		ID:          o.ID,
		CajaID:      o.CajaID,
		Nombre:      o.Nombre,
		Descripcion: o.Descripcion,
		Cantidad:    o.Cantidad,
		EsInsumo:    o.EsInsumo,
		Valor:       o.ValorEstimado,
		TieneFotos:  tieneFotos,
	}
}
