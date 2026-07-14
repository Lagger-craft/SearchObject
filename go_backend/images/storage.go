package images

import (
	"os"
	"path/filepath"

	apperrors "searchobject/errors"
)

func Eliminar(paths ...string) error {
	for _, p := range paths {
		if p == "" {
			continue
		}
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			return apperrors.Wrap(err, apperrors.ErrStorage, "no se pudo eliminar la imagen", "images.Eliminar")
		}
	}
	return nil
}

func EliminarDirectorioObjeto(baseDir, objetoID string) error {
	dir := filepath.Join(baseDir, objetoID)
	if err := os.RemoveAll(dir); err != nil {
		return apperrors.Wrap(err, apperrors.ErrStorage, "no se pudo eliminar el directorio del objeto", "images.EliminarDirectorioObjeto")
	}
	return nil
}

func PathExiste(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
