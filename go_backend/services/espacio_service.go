package services

import (
	"searchobject/db"
	apperrors "searchobject/errors"
	"searchobject/models"
	"searchobject/normalize"
)

type EspacioService struct {
	db *db.DB
}

func NewEspacioService(database *db.DB) *EspacioService {
	return &EspacioService{db: database}
}

func (s *EspacioService) Crear(usuarioID, nombre, descripcion string, padreID *string) (*models.Espacio, error) {
	if nombre == "" {
		return nil, apperrors.New(apperrors.ErrValidation, "el nombre del espacio no puede estar vacío", "EspacioService.Crear")
	}

	norm := normalize.Nombre(nombre)

	exist, err := s.buscarPorNombreNorm(usuarioID, norm)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, apperrors.New(apperrors.ErrDuplicate,
			"ya tenés un espacio llamado '"+exist.Nombre+"'", "EspacioService.Crear")
	}

	e := &models.Espacio{
		UsuarioID:   usuarioID,
		Nombre:      nombre,
		NombreNorm:  norm,
		Descripcion: descripcion,
		PadreID:     padreID,
		CreatedAt:   models.Now(),
		UpdatedAt:   models.Now(),
	}

	if err := s.db.CrearEspacio(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *EspacioService) Listar(usuarioID string) ([]models.Espacio, error) {
	return s.db.ListarEspacios(usuarioID)
}

func (s *EspacioService) Obtener(id string) (*models.Espacio, error) {
	return s.db.ObtenerEspacio(id)
}

func (s *EspacioService) Actualizar(id, nombre, descripcion string, padreID *string) (*models.Espacio, error) {
	e, err := s.db.ObtenerEspacio(id)
	if err != nil {
		return nil, err
	}

	if nombre != "" {
		norm := normalize.Nombre(nombre)
		exist, err := s.buscarPorNombreNorm(e.UsuarioID, norm)
		if err != nil {
			return nil, err
		}
		if exist != nil && exist.ID != id {
			return nil, apperrors.New(apperrors.ErrDuplicate,
				"ya tenés un espacio llamado '"+exist.Nombre+"'", "EspacioService.Actualizar")
		}
		e.Nombre = nombre
		e.NombreNorm = norm
	}

	if descripcion != "" {
		e.Descripcion = descripcion
	}
	if padreID != nil {
		e.PadreID = padreID
	}
	e.UpdatedAt = models.Now()

	if err := s.db.ActualizarEspacio(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *EspacioService) Eliminar(id string) error {
	return s.db.EliminarEspacio(id)
}

func (s *EspacioService) buscarPorNombreNorm(usuarioID, norm string) (*models.Espacio, error) {
	espacios, err := s.db.ListarEspacios(usuarioID)
	if err != nil {
		return nil, err
	}
	for _, e := range espacios {
		if e.NombreNorm == norm {
			return &e, nil
		}
	}
	return nil, nil
}
