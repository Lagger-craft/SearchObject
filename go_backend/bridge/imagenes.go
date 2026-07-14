package bridge

import (
	"bytes"
	"encoding/json"

	"github.com/google/uuid"

	apperrors "searchobject/errors"
	"searchobject/images"
	"searchobject/models"
)

type imagenDTO struct {
	ID          string   `json:"id"`
	Path        string   `json:"path"`
	ThumbPath   string   `json:"thumb_path"`
	EsPrincipal bool     `json:"es_principal"`
	AreaX       *float64 `json:"area_x,omitempty"`
	AreaY       *float64 `json:"area_y,omitempty"`
	AreaW       *float64 `json:"area_w,omitempty"`
	AreaH       *float64 `json:"area_h,omitempty"`
}

type listarImagenesRes struct {
	Imagenes []imagenDTO `json:"imagenes"`
}

func (app *App) AgregarImagen(objetoID string, imageBytes []byte) (string, error) {
	o, err := app.objetos.Obtener(objetoID)
	if err != nil {
		return "", err
	}
	_ = o

	imagenID := uuid.New().String()
	proc, err := images.Procesar(bytes.NewReader(imageBytes), app.imageDir, objetoID, imagenID, images.DefaultOptions)
	if err != nil {
		return "", err
	}

	img := &models.Imagen{
		ID:        imagenID,
		ObjetoID:  objetoID,
		Path:      proc.Path,
		ThumbPath: proc.ThumbPath,
		CreatedAt: models.Now(),
	}

	imgsExistentes, _ := app.db.ListarImagenes(objetoID)
	if len(imgsExistentes) == 0 {
		img.EsPrincipal = true
	}

	if err := app.db.CrearImagen(img); err != nil {
		images.Eliminar(proc.Path, proc.ThumbPath)
		return "", err
	}

	dto := toImagenDTO(*img)
	data, _ := json.Marshal(dto)
	return string(data), nil
}

func (app *App) AgregarImagenConArea(objetoID string, imageBytes []byte, jsonArea string) (string, error) {
	var area struct {
		X *float64 `json:"x"`
		Y *float64 `json:"y"`
		W *float64 `json:"w"`
		H *float64 `json:"h"`
	}
	if err := json.Unmarshal([]byte(jsonArea), &area); err != nil {
		return "", apperrors.Wrap(err, apperrors.ErrValidation, "área JSON inválida", "bridge.AgregarImagenConArea")
	}

	resultJSON, err := app.AgregarImagen(objetoID, imageBytes)
	if err != nil {
		return "", err
	}

	var img models.Imagen
	json.Unmarshal([]byte(resultJSON), &img)

	img.AreaX = area.X
	img.AreaY = area.Y
	img.AreaW = area.W
	img.AreaH = area.H

	if err := app.db.CrearImagen(&img); err != nil {
		return "", err
	}

	dto := toImagenDTO(img)
	data, _ := json.Marshal(dto)
	return string(data), nil
}

func (app *App) ListarImagenes(objetoID string) (string, error) {
	imgs, err := app.db.ListarImagenes(objetoID)
	if err != nil {
		return "", err
	}

	dtos := make([]imagenDTO, 0, len(imgs))
	for _, i := range imgs {
		dtos = append(dtos, toImagenDTO(i))
	}

	res := listarImagenesRes{Imagenes: dtos}
	data, _ := json.Marshal(res)
	return string(data), nil
}

func (app *App) EliminarImagen(id string) (string, error) {
	// TODO: eliminar archivo físico también
	if err := app.db.EliminarImagen(id); err != nil {
		return "", err
	}
	return `{"ok":true}`, nil
}

func toImagenDTO(i models.Imagen) imagenDTO {
	return imagenDTO{
		ID:          i.ID,
		Path:        i.Path,
		ThumbPath:   i.ThumbPath,
		EsPrincipal: i.EsPrincipal,
		AreaX:       i.AreaX,
		AreaY:       i.AreaY,
		AreaW:       i.AreaW,
		AreaH:       i.AreaH,
	}
}
