package bridge

import (
	"encoding/json"

	"github.com/google/uuid"

	apperrors "searchobject/errors"
	"searchobject/models"
)

type usuarioDTO struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email,omitempty"`
}

func (app *App) CrearUsuario(nombre, email string) (string, error) {
	if email != "" {
		existing, err := app.db.ObtenerUsuarioPorEmail(email)
		if err == nil {
			data, _ := json.Marshal(usuarioDTO{ID: existing.ID, Nombre: existing.Nombre, Email: existing.Email})
			return string(data), nil
		}
	}

	u := &models.Usuario{
		ID:        uuid.New().String(),
		Nombre:    nombre,
		Email:     email,
		CreatedAt: models.Now(),
		UpdatedAt: models.Now(),
	}

	if err := app.db.CrearUsuario(u); err != nil {
		return "", err
	}

	data, _ := json.Marshal(usuarioDTO{ID: u.ID, Nombre: u.Nombre, Email: u.Email})
	return string(data), nil
}

func (app *App) ActualizarUsuario(jsonReq string) (string, error) {
	var req struct {
		ID     string `json:"id"`
		Nombre string `json:"nombre"`
		Email  string `json:"email"`
	}
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		return "", apperrors.Wrap(err, apperrors.ErrValidation, "JSON inválido", "bridge.ActualizarUsuario")
	}

	u, err := app.db.ObtenerUsuario(req.ID)
	if err != nil {
		return "", err
	}

	if req.Nombre != "" {
		u.Nombre = req.Nombre
	}
	if req.Email != "" {
		u.Email = req.Email
	}

	if err := app.db.ActualizarUsuario(u); err != nil {
		return "", err
	}

	data, _ := json.Marshal(usuarioDTO{ID: u.ID, Nombre: u.Nombre, Email: u.Email})
	return string(data), nil
}

func (app *App) ListarUsuarios() (string, error) {
	usuarios, err := app.db.ListarUsuarios()
	if err != nil {
		return "", err
	}

	var dtos []usuarioDTO
	for _, u := range usuarios {
		dtos = append(dtos, usuarioDTO{ID: u.ID, Nombre: u.Nombre, Email: u.Email})
	}

	data, _ := json.Marshal(dtos)
	return string(data), nil
}

func (app *App) ObtenerUsuario(id string) (string, error) {
	u, err := app.db.ObtenerUsuario(id)
	if err != nil {
		return "", err
	}
	data, _ := json.Marshal(usuarioDTO{ID: u.ID, Nombre: u.Nombre, Email: u.Email})
	return string(data), nil
}
