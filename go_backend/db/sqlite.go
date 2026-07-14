package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	apperrors "searchobject/errors"
)

type DB struct {
	*sql.DB
}

func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo abrir la base de datos", "db.Open")
	}

	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS usuarios (
			id TEXT PRIMARY KEY,
			nombre TEXT NOT NULL,
			email TEXT,
			avatar TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS espacios (
			id TEXT PRIMARY KEY,
			usuario_id TEXT NOT NULL REFERENCES usuarios(id),
			nombre TEXT NOT NULL,
			nombre_norm TEXT NOT NULL,
			descripcion TEXT,
			padre_id TEXT REFERENCES espacios(id),
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			UNIQUE(usuario_id, nombre_norm)
		)`,
		`CREATE TABLE IF NOT EXISTS cajas (
			id TEXT PRIMARY KEY,
			espacio_id TEXT NOT NULL REFERENCES espacios(id),
			usuario_id TEXT NOT NULL REFERENCES usuarios(id),
			nombre TEXT NOT NULL,
			nombre_norm TEXT NOT NULL,
			descripcion TEXT,
			capacidad_max INTEGER,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			UNIQUE(espacio_id, nombre_norm)
		)`,
		`CREATE TABLE IF NOT EXISTS objetos (
			id TEXT PRIMARY KEY,
			caja_id TEXT NOT NULL REFERENCES cajas(id),
			usuario_id TEXT NOT NULL REFERENCES usuarios(id),
			nombre TEXT NOT NULL,
			nombre_norm TEXT NOT NULL,
			descripcion TEXT,
			cantidad INTEGER NOT NULL DEFAULT 1,
			es_insumo INTEGER NOT NULL DEFAULT 0,
			valor_estimado REAL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			UNIQUE(caja_id, nombre_norm)
		)`,
		`CREATE TABLE IF NOT EXISTS imagenes (
			id TEXT PRIMARY KEY,
			objeto_id TEXT NOT NULL REFERENCES objetos(id) ON DELETE CASCADE,
			path TEXT NOT NULL,
			thumb_path TEXT NOT NULL,
			area_x REAL,
			area_y REAL,
			area_w REAL,
			area_h REAL,
			es_principal INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tags (
			id TEXT PRIMARY KEY,
			nombre TEXT NOT NULL,
			nombre_norm TEXT NOT NULL UNIQUE,
			color TEXT,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS objeto_tags (
			objeto_id TEXT NOT NULL REFERENCES objetos(id) ON DELETE CASCADE,
			tag_id TEXT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
			PRIMARY KEY (objeto_id, tag_id)
		)`,
		`CREATE TABLE IF NOT EXISTS movimientos (
			id TEXT PRIMARY KEY,
			objeto_id TEXT NOT NULL REFERENCES objetos(id),
			desde_caja_id TEXT REFERENCES cajas(id),
			hacia_caja_id TEXT REFERENCES cajas(id),
			tipo TEXT NOT NULL CHECK(tipo IN ('entrada','salida','traslado','prestamo')),
			nota TEXT,
			fecha TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS historial (
			id TEXT PRIMARY KEY,
			entidad_tipo TEXT NOT NULL,
			entidad_id TEXT NOT NULL,
			accion TEXT NOT NULL CHECK(accion IN ('creado','actualizado','eliminado','movido')),
			detalle TEXT,
			usuario_id TEXT NOT NULL REFERENCES usuarios(id),
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS alertas (
			id TEXT PRIMARY KEY,
			tipo TEXT NOT NULL CHECK(tipo IN ('faltante','stock_bajo','prestamo_vencido','objeto_perdido','caja_saturada')),
			entidad_tipo TEXT NOT NULL,
			entidad_id TEXT NOT NULL,
			mensaje TEXT NOT NULL,
			leida INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			resuelta_at TEXT
		)`,
	}

	for i, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("migracion %d: %w", i+1, err)
		}
	}

	return nil
}
