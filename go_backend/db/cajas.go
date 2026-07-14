package db

import (
	"database/sql"

	apperrors "searchobject/errors"
	"searchobject/models"

	"github.com/google/uuid"
)

func (db *DB) CrearCaja(c *models.Caja) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = models.Now()
		c.UpdatedAt = models.Now()
	}
	_, err := db.Exec(`INSERT INTO cajas (id,espacio_id,usuario_id,nombre,nombre_norm,descripcion,capacidad_max,created_at,updated_at)
		VALUES (?,?,?,?,?,?,?,?,?)`, c.ID, c.EspacioID, c.UsuarioID, c.Nombre, c.NombreNorm, nullStr(c.Descripcion), c.CapacidadMax, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		if isUnique(err) {
			return apperrors.Wrap(err, apperrors.ErrDuplicate, "ya existe una caja con ese nombre en este espacio", "db.CrearCaja")
		}
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo crear la caja", "db.CrearCaja")
	}
	return nil
}

func (db *DB) ListarCajas(espacioID string) ([]models.Caja, error) {
	rows, err := db.Query(`SELECT id,espacio_id,usuario_id,nombre,nombre_norm,descripcion,capacidad_max,created_at,updated_at
		FROM cajas WHERE espacio_id=? ORDER BY nombre_norm`, espacioID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudieron listar cajas", "db.ListarCajas")
	}
	defer rows.Close()

	var cs []models.Caja
	for rows.Next() {
		var c models.Caja
		var desc, cap sql.NullString
		if err := rows.Scan(&c.ID, &c.EspacioID, &c.UsuarioID, &c.Nombre, &c.NombreNorm, &desc, &cap, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.Descripcion = desc.String
		if cap.Valid {
			v := parseInt(cap.String)
			c.CapacidadMax = &v
		}
		cs = append(cs, c)
	}
	return cs, rows.Err()
}

func (db *DB) ObtenerCaja(id string) (*models.Caja, error) {
	row := db.QueryRow(`SELECT id,espacio_id,usuario_id,nombre,nombre_norm,descripcion,capacidad_max,created_at,updated_at
		FROM cajas WHERE id=?`, id)
	var c models.Caja
	var desc, cap sql.NullString
	if err := row.Scan(&c.ID, &c.EspacioID, &c.UsuarioID, &c.Nombre, &c.NombreNorm, &desc, &cap, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.New(apperrors.ErrNotFound, "caja no encontrada", "db.ObtenerCaja")
		}
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo obtener la caja", "db.ObtenerCaja")
	}
	c.Descripcion = desc.String
	if cap.Valid {
		v := parseInt(cap.String)
		c.CapacidadMax = &v
	}
	return &c, nil
}

func (db *DB) ActualizarCaja(c *models.Caja) error {
	r, err := db.Exec(`UPDATE cajas SET nombre=?,nombre_norm=?,descripcion=?,capacidad_max=?,updated_at=?
		WHERE id=?`, c.Nombre, c.NombreNorm, nullStr(c.Descripcion), c.CapacidadMax, c.UpdatedAt, c.ID)
	if err != nil {
		if isUnique(err) {
			return apperrors.Wrap(err, apperrors.ErrDuplicate, "ya existe otra caja con ese nombre", "db.ActualizarCaja")
		}
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo actualizar la caja", "db.ActualizarCaja")
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return apperrors.New(apperrors.ErrNotFound, "caja no encontrada", "db.ActualizarCaja")
	}
	return nil
}

func (db *DB) EliminarCaja(id string) error {
	r, err := db.Exec(`DELETE FROM cajas WHERE id=?`, id)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo eliminar la caja", "db.EliminarCaja")
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return apperrors.New(apperrors.ErrNotFound, "caja no encontrada", "db.EliminarCaja")
	}
	return nil
}
