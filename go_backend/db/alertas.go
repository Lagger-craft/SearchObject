package db

import (
	"database/sql"

	apperrors "searchobject/errors"
	"searchobject/models"

	"github.com/google/uuid"
)

func (db *DB) CrearAlerta(a *models.Alerta) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = models.Now()
	}
	_, err := db.Exec(`INSERT INTO alertas (id,tipo,entidad_tipo,entidad_id,mensaje,leida,created_at,resuelta_at)
		VALUES (?,?,?,?,?,?,?,?)`, a.ID, a.Tipo, a.EntidadTipo, a.EntidadID, a.Mensaje, boolInt(a.Leida), a.CreatedAt, a.ResueltaAt)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo crear la alerta", "db.CrearAlerta")
	}
	return nil
}

func (db *DB) ListarAlertas(leidas *bool) ([]models.Alerta, error) {
	var rows *sql.Rows
	var err error
	if leidas != nil {
		rows, err = db.Query(`SELECT id,tipo,entidad_tipo,entidad_id,mensaje,leida,created_at,resuelta_at
			FROM alertas WHERE leida=? ORDER BY created_at DESC`, boolInt(*leidas))
	} else {
		rows, err = db.Query(`SELECT id,tipo,entidad_tipo,entidad_id,mensaje,leida,created_at,resuelta_at
			FROM alertas ORDER BY created_at DESC`)
	}
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudieron listar alertas", "db.ListarAlertas")
	}
	defer rows.Close()

	var as []models.Alerta
	for rows.Next() {
		var a models.Alerta
		var resuelta sql.NullString
		if err := rows.Scan(&a.ID, &a.Tipo, &a.EntidadTipo, &a.EntidadID, &a.Mensaje, &a.Leida, &a.CreatedAt, &resuelta); err != nil {
			return nil, err
		}
		if resuelta.Valid {
			t := parseTime(resuelta.String)
			a.ResueltaAt = &t
		}
		as = append(as, a)
	}
	return as, rows.Err()
}

func (db *DB) MarcarAlertaLeida(id string) error {
	r, err := db.Exec(`UPDATE alertas SET leida=1 WHERE id=?`, id)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo marcar la alerta como leída", "db.MarcarAlertaLeida")
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return apperrors.New(apperrors.ErrNotFound, "alerta no encontrada", "db.MarcarAlertaLeida")
	}
	return nil
}

func (db *DB) ResolverAlerta(id string) error {
	r, err := db.Exec(`UPDATE alertas SET leida=1, resuelta_at=datetime('now') WHERE id=?`, id)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo resolver la alerta", "db.ResolverAlerta")
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return apperrors.New(apperrors.ErrNotFound, "alerta no encontrada", "db.ResolverAlerta")
	}
	return nil
}

func (db *DB) AlertasNoLeidasCount() (int, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM alertas WHERE leida=0`).Scan(&count)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo contar alertas", "db.AlertasNoLeidasCount")
	}
	return count, nil
}
