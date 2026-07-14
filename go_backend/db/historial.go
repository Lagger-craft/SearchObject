package db

import (
	"database/sql"

	apperrors "searchobject/errors"
	"searchobject/models"

	"github.com/google/uuid"
)

func (db *DB) CrearHistorial(h *models.Historial) error {
	if h.ID == "" {
		h.ID = uuid.New().String()
	}
	if h.CreatedAt.IsZero() {
		h.CreatedAt = models.Now()
	}
	_, err := db.Exec(`INSERT INTO historial (id,entidad_tipo,entidad_id,accion,detalle,usuario_id,created_at)
		VALUES (?,?,?,?,?,?,?)`, h.ID, h.EntidadTipo, h.EntidadID, h.Accion, nullStr(h.Detalle), h.UsuarioID, h.CreatedAt)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo crear el historial", "db.CrearHistorial")
	}
	return nil
}

func (db *DB) ListarHistorial(usuarioID string, limite int) ([]models.Historial, error) {
	if limite <= 0 {
		limite = 50
	}
	rows, err := db.Query(`SELECT id,entidad_tipo,entidad_id,accion,detalle,usuario_id,created_at
		FROM historial WHERE usuario_id=? ORDER BY created_at DESC LIMIT ?`, usuarioID, limite)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo listar el historial", "db.ListarHistorial")
	}
	defer rows.Close()

	var hs []models.Historial
	for rows.Next() {
		var h models.Historial
		var detalle sql.NullString
		if err := rows.Scan(&h.ID, &h.EntidadTipo, &h.EntidadID, &h.Accion, &detalle, &h.UsuarioID, &h.CreatedAt); err != nil {
			return nil, err
		}
		h.Detalle = detalle.String
		hs = append(hs, h)
	}
	return hs, rows.Err()
}

func (db *DB) HistorialPorEntidad(usuarioID, entidadTipo, entidadID string) ([]models.Historial, error) {
	rows, err := db.Query(`SELECT id,entidad_tipo,entidad_id,accion,detalle,usuario_id,created_at
		FROM historial WHERE usuario_id=? AND entidad_tipo=? AND entidad_id=? ORDER BY created_at DESC`,
		usuarioID, entidadTipo, entidadID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo obtener historial de la entidad", "db.HistorialPorEntidad")
	}
	defer rows.Close()

	var hs []models.Historial
	for rows.Next() {
		var h models.Historial
		var detalle sql.NullString
		if err := rows.Scan(&h.ID, &h.EntidadTipo, &h.EntidadID, &h.Accion, &detalle, &h.UsuarioID, &h.CreatedAt); err != nil {
			return nil, err
		}
		h.Detalle = detalle.String
		hs = append(hs, h)
	}
	return hs, rows.Err()
}
