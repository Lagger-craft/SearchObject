package db

import (
	"database/sql"

	apperrors "searchobject/errors"
	"searchobject/models"

	"github.com/google/uuid"
)

func (db *DB) CrearMovimiento(m *models.Movimiento) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = models.Now()
	}
	if m.Fecha.IsZero() {
		m.Fecha = models.Now()
	}
	_, err := db.Exec(`INSERT INTO movimientos (id,objeto_id,desde_caja_id,hacia_caja_id,tipo,nota,fecha,created_at)
		VALUES (?,?,?,?,?,?,?,?)`, m.ID, m.ObjetoID, m.DesdeCajaID, m.HaciaCajaID, m.Tipo, nullStr(m.Nota), m.Fecha, m.CreatedAt)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo registrar el movimiento", "db.CrearMovimiento")
	}
	return nil
}

func (db *DB) ListarMovimientos(objetoID string) ([]models.Movimiento, error) {
	rows, err := db.Query(`SELECT id,objeto_id,desde_caja_id,hacia_caja_id,tipo,nota,fecha,created_at
		FROM movimientos WHERE objeto_id=? ORDER BY fecha DESC`, objetoID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudieron listar movimientos", "db.ListarMovimientos")
	}
	defer rows.Close()

	var ms []models.Movimiento
	for rows.Next() {
		var m models.Movimiento
		var desde, hacia, nota sql.NullString
		if err := rows.Scan(&m.ID, &m.ObjetoID, &desde, &hacia, &m.Tipo, &nota, &m.Fecha, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.DesdeCajaID = nullFromSQL(desde)
		m.HaciaCajaID = nullFromSQL(hacia)
		m.Nota = nota.String
		ms = append(ms, m)
	}
	return ms, rows.Err()
}

func (db *DB) MovimientosRecientes(usuarioID string, limite int) ([]models.Movimiento, error) {
	rows, err := db.Query(`SELECT m.id,m.objeto_id,m.desde_caja_id,m.hacia_caja_id,m.tipo,m.nota,m.fecha,m.created_at
		FROM movimientos m
		JOIN objetos o ON o.id=m.objeto_id
		WHERE o.usuario_id=?
		ORDER BY m.created_at DESC LIMIT ?`, usuarioID, limite)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudieron obtener movimientos recientes", "db.MovimientosRecientes")
	}
	defer rows.Close()

	var ms []models.Movimiento
	for rows.Next() {
		var m models.Movimiento
		var desde, hacia, nota sql.NullString
		if err := rows.Scan(&m.ID, &m.ObjetoID, &desde, &hacia, &m.Tipo, &nota, &m.Fecha, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.DesdeCajaID = nullFromSQL(desde)
		m.HaciaCajaID = nullFromSQL(hacia)
		m.Nota = nota.String
		ms = append(ms, m)
	}
	return ms, rows.Err()
}

func (db *DB) HistorialMovimiento(movimientoID string) (*models.Historial, error) {
	row := db.QueryRow(`SELECT id,entidad_tipo,entidad_id,accion,detalle,usuario_id,created_at
		FROM historial WHERE entidad_id=? AND entidad_tipo='movimiento' LIMIT 1`, movimientoID)
	var h models.Historial
	var detalle sql.NullString
	if err := row.Scan(&h.ID, &h.EntidadTipo, &h.EntidadID, &h.Accion, &detalle, &h.UsuarioID, &h.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.New(apperrors.ErrNotFound, "historial no encontrado", "db.HistorialMovimiento")
		}
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo obtener el historial", "db.HistorialMovimiento")
	}
	h.Detalle = detalle.String
	return &h, nil
}
