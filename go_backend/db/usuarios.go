package db

import (
	"database/sql"

	apperrors "searchobject/errors"
	"searchobject/models"

	"github.com/google/uuid"
)

func (db *DB) CrearUsuario(u *models.Usuario) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	_, err := db.Exec(`INSERT INTO usuarios (id,nombre,email,avatar,created_at,updated_at)
		VALUES (?,?,?,?,?,?)`, u.ID, u.Nombre, nullStr(u.Email), nullStr(u.Avatar), u.CreatedAt, u.UpdatedAt)
	if err != nil {
		if isUnique(err) {
			return apperrors.Wrap(err, apperrors.ErrDuplicate, "ya existe un usuario con ese email", "db.CrearUsuario")
		}
		return apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo crear el usuario", "db.CrearUsuario")
	}
	return nil
}

func (db *DB) ObtenerUsuario(id string) (*models.Usuario, error) {
	row := db.QueryRow(`SELECT id,nombre,email,avatar,created_at,updated_at FROM usuarios WHERE id=?`, id)
	var u models.Usuario
	var email, avatar sql.NullString
	if err := row.Scan(&u.ID, &u.Nombre, &email, &avatar, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.New(apperrors.ErrNotFound, "usuario no encontrado", "db.ObtenerUsuario")
		}
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo obtener el usuario", "db.ObtenerUsuario")
	}
	u.Email = email.String
	u.Avatar = avatar.String
	return &u, nil
}

func (db *DB) ObtenerUsuarioPorEmail(email string) (*models.Usuario, error) {
	row := db.QueryRow(`SELECT id,nombre,email,avatar,created_at,updated_at FROM usuarios WHERE email=? ORDER BY rowid ASC LIMIT 1`, email)
	var u models.Usuario
	var mail, avatar sql.NullString
	if err := row.Scan(&u.ID, &u.Nombre, &mail, &avatar, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.New(apperrors.ErrNotFound, "usuario no encontrado", "db.ObtenerUsuarioPorEmail")
		}
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudo obtener el usuario por email", "db.ObtenerUsuarioPorEmail")
	}
	u.Email = mail.String
	u.Avatar = avatar.String
	return &u, nil
}

func (db *DB) ListarUsuarios() ([]models.Usuario, error) {
	rows, err := db.Query(`SELECT id,nombre,email,avatar,created_at,updated_at FROM usuarios ORDER BY nombre`)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrDatabase, "no se pudieron listar usuarios", "db.ListarUsuarios")
	}
	defer rows.Close()

	var us []models.Usuario
	for rows.Next() {
		var u models.Usuario
		var email, avatar sql.NullString
		if err := rows.Scan(&u.ID, &u.Nombre, &email, &avatar, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		u.Email = email.String
		u.Avatar = avatar.String
		us = append(us, u)
	}
	return us, rows.Err()
}
