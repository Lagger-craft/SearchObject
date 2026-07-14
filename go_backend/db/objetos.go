package db

import (
	"database/sql"

	apperrors "searchobject/errors"
	"searchobject/models"

	"github.com/google/uuid"
)

func (db *DB) CrearObjeto(o *models.Objeto) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	if o.CreatedAt.IsZero() {
		o.CreatedAt = models.Now()
		o.UpdatedAt = models.Now()
	}
	_, err := db.Exec(`INSERT INTO objetos (id,caja_id,usuario_id,nombre,nombre_norm,descripcion,cantidad,es_insumo,valor_estimado,created_at,updated_at)
		VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		o.ID, o.CajaID, o.UsuarioID, o.Nombre, o.NombreNorm, nullStr(o.Descripcion),
		o.Cantidad, boolInt(o.EsInsumo), o.ValorEstimado, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		if isUnique(err) {
			return apperrors.Wrap(err, apperrors.ErrDuplicate, "ya existe un objeto con ese nombre en esta caja", "db.CrearObjeto")
		}
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo crear el objeto", "db.CrearObjeto")
	}
	return nil
}

func (db *DB) ListarObjetos(cajaID string) ([]models.Objeto, error) {
	rows, err := db.Query(`SELECT id,caja_id,usuario_id,nombre,nombre_norm,descripcion,cantidad,es_insumo,valor_estimado,created_at,updated_at
		FROM objetos WHERE caja_id=? ORDER BY nombre_norm`, cajaID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudieron listar objetos", "db.ListarObjetos")
	}
	defer rows.Close()
	return scanObjetos(rows)
}

func (db *DB) ObtenerObjeto(id string) (*models.Objeto, error) {
	row := db.QueryRow(`SELECT id,caja_id,usuario_id,nombre,nombre_norm,descripcion,cantidad,es_insumo,valor_estimado,created_at,updated_at
		FROM objetos WHERE id=?`, id)
	return scanObjeto(row)
}

func (db *DB) ActualizarObjeto(o *models.Objeto) error {
	r, err := db.Exec(`UPDATE objetos SET caja_id=?,nombre=?,nombre_norm=?,descripcion=?,cantidad=?,es_insumo=?,valor_estimado=?,updated_at=?
		WHERE id=?`, o.CajaID, o.Nombre, o.NombreNorm, nullStr(o.Descripcion), o.Cantidad, boolInt(o.EsInsumo), o.ValorEstimado, o.UpdatedAt, o.ID)
	if err != nil {
		if isUnique(err) {
			return apperrors.Wrap(err, apperrors.ErrDuplicate, "ya existe otro objeto con ese nombre", "db.ActualizarObjeto")
		}
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo actualizar el objeto", "db.ActualizarObjeto")
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return apperrors.New(apperrors.ErrNotFound, "objeto no encontrado", "db.ActualizarObjeto")
	}
	return nil
}

func (db *DB) EliminarObjeto(id string) error {
	r, err := db.Exec(`DELETE FROM objetos WHERE id=?`, id)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo eliminar el objeto", "db.EliminarObjeto")
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return apperrors.New(apperrors.ErrNotFound, "objeto no encontrado", "db.EliminarObjeto")
	}
	return nil
}

func (db *DB) BuscarObjetos(usuarioID, termino string) ([]models.Objeto, error) {
	rows, err := db.Query(`SELECT o.id,o.caja_id,o.usuario_id,o.nombre,o.nombre_norm,o.descripcion,o.cantidad,o.es_insumo,o.valor_estimado,o.created_at,o.updated_at
		FROM objetos o
		JOIN cajas c ON c.id = o.caja_id
		JOIN espacios e ON e.id = c.espacio_id
		WHERE o.usuario_id=? AND (o.nombre_norm LIKE ? OR o.descripcion LIKE ?)
		ORDER BY o.nombre_norm`, usuarioID, "%"+termino+"%", "%"+termino+"%")
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "error en búsqueda", "db.BuscarObjetos")
	}
	defer rows.Close()
	return scanObjetos(rows)
}

func scanObjetos(rows *sql.Rows) ([]models.Objeto, error) {
	var os []models.Objeto
	for rows.Next() {
		var o models.Objeto
		var desc sql.NullString
		if err := rows.Scan(&o.ID, &o.CajaID, &o.UsuarioID, &o.Nombre, &o.NombreNorm, &desc, &o.Cantidad, &o.EsInsumo, &o.ValorEstimado, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		o.Descripcion = desc.String
		os = append(os, o)
	}
	return os, rows.Err()
}

func scanObjeto(row *sql.Row) (*models.Objeto, error) {
	var o models.Objeto
	var desc sql.NullString
	if err := row.Scan(&o.ID, &o.CajaID, &o.UsuarioID, &o.Nombre, &o.NombreNorm, &desc, &o.Cantidad, &o.EsInsumo, &o.ValorEstimado, &o.CreatedAt, &o.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.New(apperrors.ErrNotFound, "objeto no encontrado", "db.ObtenerObjeto")
		}
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo obtener el objeto", "db.ObtenerObjeto")
	}
	o.Descripcion = desc.String
	return &o, nil
}
