package bridge

import (
	"encoding/json"

	"github.com/google/uuid"

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
