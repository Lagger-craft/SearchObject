package db

import (
	"database/sql"

	apperrors "searchobject/errors"
	"searchobject/models"

	"github.com/google/uuid"
)

func (db *DB) CrearTag(t *models.Tag) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	_, err := db.Exec(`INSERT INTO tags (id,nombre,nombre_norm,color,created_at) VALUES (?,?,?,?,?)`,
		t.ID, t.Nombre, t.NombreNorm, nullStr(t.Color), t.CreatedAt)
	if err != nil {
		if isUnique(err) {
			return apperrors.Wrap(err, apperrors.ErrDuplicate, "ya existe un tag con ese nombre", "db.CrearTag")
		}
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo crear el tag", "db.CrearTag")
	}
	return nil
}

func (db *DB) ListarTags() ([]models.Tag, error) {
	rows, err := db.Query(`SELECT id,nombre,nombre_norm,color,created_at FROM tags ORDER BY nombre_norm`)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudieron listar tags", "db.ListarTags")
	}
	defer rows.Close()

	var ts []models.Tag
	for rows.Next() {
		var t models.Tag
		var color sql.NullString
		if err := rows.Scan(&t.ID, &t.Nombre, &t.NombreNorm, &color, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.Color = color.String
		ts = append(ts, t)
	}
	return ts, rows.Err()
}

func (db *DB) AgregarTagAObjeto(objetoID, tagID string) error {
	_, err := db.Exec(`INSERT OR IGNORE INTO objeto_tags (objeto_id,tag_id) VALUES (?,?)`, objetoID, tagID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo agregar el tag al objeto", "db.AgregarTagAObjeto")
	}
	return nil
}

func (db *DB) QuitarTagDeObjeto(objetoID, tagID string) error {
	_, err := db.Exec(`DELETE FROM objeto_tags WHERE objeto_id=? AND tag_id=?`, objetoID, tagID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo quitar el tag del objeto", "db.QuitarTagDeObjeto")
	}
	return nil
}

func (db *DB) TagsDeObjeto(objetoID string) ([]models.Tag, error) {
	rows, err := db.Query(`SELECT t.id,t.nombre,t.nombre_norm,t.color,t.created_at
		FROM tags t JOIN objeto_tags ot ON ot.tag_id=t.id WHERE ot.objeto_id=?`, objetoID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudieron obtener tags del objeto", "db.TagsDeObjeto")
	}
	defer rows.Close()

	var ts []models.Tag
	for rows.Next() {
		var t models.Tag
		var color sql.NullString
		if err := rows.Scan(&t.ID, &t.Nombre, &t.NombreNorm, &color, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.Color = color.String
		ts = append(ts, t)
	}
	return ts, rows.Err()
}

func (db *DB) ObjetosPorTag(usuarioID, tagID string) ([]models.Objeto, error) {
	rows, err := db.Query(`SELECT o.id,o.usuario_id,o.caja_id,o.nombre,o.descripcion,o.cantidad,o.es_insumo,o.valor_estimado,o.created_at,o.updated_at
		FROM objetos o
		JOIN objeto_tags ot ON ot.objeto_id=o.id
		WHERE ot.tag_id=? AND o.usuario_id=?`, tagID, usuarioID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudieron obtener objetos por tag", "db.ObjetosPorTag")
	}
	defer rows.Close()

	var os []models.Objeto
	for rows.Next() {
		var o models.Objeto
		var valor sql.NullFloat64
		if err := rows.Scan(&o.ID, &o.UsuarioID, &o.CajaID, &o.Nombre, &o.Descripcion, &o.Cantidad, &o.EsInsumo, &valor, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		if valor.Valid {
			o.ValorEstimado = &valor.Float64
		}
		os = append(os, o)
	}
	return os, rows.Err()
}
