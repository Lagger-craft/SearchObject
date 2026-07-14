package db

import (
	"database/sql"

	apperrors "searchobject/errors"
	"searchobject/models"

	"github.com/google/uuid"
)

func (db *DB) CrearEspacio(e *models.Espacio) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = models.Now()
		e.UpdatedAt = models.Now()
	}
	_, err := db.Exec(`INSERT INTO espacios (id,usuario_id,nombre,nombre_norm,descripcion,padre_id,created_at,updated_at)
		VALUES (?,?,?,?,?,?,?,?)`, e.ID, e.UsuarioID, e.Nombre, e.NombreNorm, nullStr(e.Descripcion), e.PadreID, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		if isUnique(err) {
			return apperrors.Wrap(err, apperrors.ErrDuplicate, "ya existe un espacio con ese nombre", "db.CrearEspacio")
		}
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo crear el espacio", "db.CrearEspacio")
	}
	return nil
}

func (db *DB) ListarEspacios(usuarioID string) ([]models.Espacio, error) {
	rows, err := db.Query(`SELECT id,usuario_id,nombre,nombre_norm,descripcion,padre_id,created_at,updated_at
		FROM espacios WHERE usuario_id=? ORDER BY nombre_norm`, usuarioID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudieron listar espacios", "db.ListarEspacios")
	}
	defer rows.Close()

	return scanEspacios(rows)
}

func (db *DB) ObtenerEspacio(id string) (*models.Espacio, error) {
	row := db.QueryRow(`SELECT id,usuario_id,nombre,nombre_norm,descripcion,padre_id,created_at,updated_at
		FROM espacios WHERE id=?`, id)
	e, err := scanEspacio(row)
	if err == sql.ErrNoRows {
		return nil, apperrors.New(apperrors.ErrNotFound, "espacio no encontrado", "db.ObtenerEspacio")
	}
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo obtener el espacio", "db.ObtenerEspacio")
	}
	return e, nil
}

func (db *DB) ActualizarEspacio(e *models.Espacio) error {
	r, err := db.Exec(`UPDATE espacios SET nombre=?,nombre_norm=?,descripcion=?,padre_id=?,updated_at=?
		WHERE id=?`, e.Nombre, e.NombreNorm, nullStr(e.Descripcion), e.PadreID, e.UpdatedAt, e.ID)
	if err != nil {
		if isUnique(err) {
			return apperrors.Wrap(err, apperrors.ErrDuplicate, "ya existe otro espacio con ese nombre", "db.ActualizarEspacio")
		}
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo actualizar el espacio", "db.ActualizarEspacio")
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return apperrors.New(apperrors.ErrNotFound, "espacio no encontrado", "db.ActualizarEspacio")
	}
	return nil
}

func (db *DB) EliminarEspacio(id string) error {
	r, err := db.Exec(`DELETE FROM espacios WHERE id=?`, id)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo eliminar el espacio", "db.EliminarEspacio")
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return apperrors.New(apperrors.ErrNotFound, "espacio no encontrado", "db.EliminarEspacio")
	}
	return nil
}

func scanEspacios(rows *sql.Rows) ([]models.Espacio, error) {
	var es []models.Espacio
	for rows.Next() {
		var e models.Espacio
		var desc, padre sql.NullString
		if err := rows.Scan(&e.ID, &e.UsuarioID, &e.Nombre, &e.NombreNorm, &desc, &padre, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		e.Descripcion = desc.String
		e.PadreID = nullFromSQL(padre)
		es = append(es, e)
	}
	return es, rows.Err()
}

func scanEspacio(row *sql.Row) (*models.Espacio, error) {
	var e models.Espacio
	var desc, padre sql.NullString
	if err := row.Scan(&e.ID, &e.UsuarioID, &e.Nombre, &e.NombreNorm, &desc, &padre, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return nil, err
	}
	e.Descripcion = desc.String
	e.PadreID = nullFromSQL(padre)
	return &e, nil
}
