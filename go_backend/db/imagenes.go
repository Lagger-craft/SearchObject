package db

import (
	"database/sql"

	apperrors "searchobject/errors"
	"searchobject/models"

	"github.com/google/uuid"
)

func (db *DB) CrearImagen(i *models.Imagen) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	_, err := db.Exec(`INSERT INTO imagenes (id,objeto_id,path,thumb_path,area_x,area_y,area_w,area_h,es_principal,created_at)
		VALUES (?,?,?,?,?,?,?,?,?,?)`,
		i.ID, i.ObjetoID, i.Path, i.ThumbPath, i.AreaX, i.AreaY, i.AreaW, i.AreaH, boolInt(i.EsPrincipal), i.CreatedAt)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo guardar la imagen", "db.CrearImagen")
	}
	return nil
}

func (db *DB) ListarImagenes(objetoID string) ([]models.Imagen, error) {
	rows, err := db.Query(`SELECT id,objeto_id,path,thumb_path,area_x,area_y,area_w,area_h,es_principal,created_at
		FROM imagenes WHERE objeto_id=? ORDER BY es_principal DESC, created_at`, objetoID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudieron listar imágenes", "db.ListarImagenes")
	}
	defer rows.Close()

	var is []models.Imagen
	for rows.Next() {
		var im models.Imagen
		if err := rows.Scan(&im.ID, &im.ObjetoID, &im.Path, &im.ThumbPath, &im.AreaX, &im.AreaY, &im.AreaW, &im.AreaH, &im.EsPrincipal, &im.CreatedAt); err != nil {
			return nil, err
		}
		is = append(is, im)
	}
	return is, rows.Err()
}

func (db *DB) ObtenerImagenPrincipal(objetoID string) (*models.Imagen, error) {
	row := db.QueryRow(`SELECT id,objeto_id,path,thumb_path,area_x,area_y,area_w,area_h,es_principal,created_at
		FROM imagenes WHERE objeto_id=? AND es_principal=1 LIMIT 1`, objetoID)
	return scanImagen(row)
}

func (db *DB) EliminarImagen(id string) error {
	r, err := db.Exec(`DELETE FROM imagenes WHERE id=?`, id)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo eliminar la imagen", "db.EliminarImagen")
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return apperrors.New(apperrors.ErrNotFound, "imagen no encontrada", "db.EliminarImagen")
	}
	return nil
}

func scanImagen(row *sql.Row) (*models.Imagen, error) {
	var i models.Imagen
	if err := row.Scan(&i.ID, &i.ObjetoID, &i.Path, &i.ThumbPath, &i.AreaX, &i.AreaY, &i.AreaW, &i.AreaH, &i.EsPrincipal, &i.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.New(apperrors.ErrNotFound, "imagen no encontrada", "db.ObtenerImagenPrincipal")
		}
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "error al leer imagen", "db.ObtenerImagenPrincipal")
	}
	return &i, nil
}

func (db *DB) EliminarImagenesPorObjeto(objetoID string) error {
	_, err := db.Exec(`DELETE FROM imagenes WHERE objeto_id=?`, objetoID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudieron eliminar imágenes del objeto", "db.EliminarImagenesPorObjeto")
	}
	return nil
}
