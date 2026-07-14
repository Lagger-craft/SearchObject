package images

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"

	apperrors "searchobject/errors"
)

type Procesado struct {
	Path      string
	ThumbPath string
}

type Options struct {
	MaxDimension int
	Quality      int
	ThumbSize    int
}

var DefaultOptions = Options{
	MaxDimension: 1920,
	Quality:      82,
	ThumbSize:    256,
}

func Procesar(lectura io.Reader, baseDir, objetoID, imagenID string, opts Options) (*Procesado, error) {
	data, err := io.ReadAll(lectura)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrImageProcess, "no se pudo leer la imagen", "images.Procesar")
	}

	img, formato, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrImageProcess, "formato de imagen no soportado", "images.Procesar")
	}
	_ = formato

	imgDir := filepath.Join(baseDir, objetoID)
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrStorage, "no se pudo crear el directorio de imágenes", "images.Procesar")
	}

	fullPath := filepath.Join(imgDir, imagenID+".jpg")
	thumbPath := filepath.Join(imgDir, imagenID+"_thumb.jpg")

	redim := redimensionar(img, opts.MaxDimension)
	if err := guardarJPEG(redim, fullPath, opts.Quality); err != nil {
		return nil, err
	}

	thumb := redimensionar(img, opts.ThumbSize)
	if err := guardarJPEG(thumb, thumbPath, 70); err != nil {
		return nil, err
	}

	return &Procesado{Path: fullPath, ThumbPath: thumbPath}, nil
}

func redimensionar(img image.Image, maxDim int) image.Image {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w <= maxDim && h <= maxDim {
		return img
	}

	nuevoW, nuevoH := w, h
	if w > h {
		nuevoW = maxDim
		nuevoH = h * maxDim / w
	} else {
		nuevoH = maxDim
		nuevoW = w * maxDim / h
	}

	dst := image.NewRGBA(image.Rect(0, 0, nuevoW, nuevoH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

func guardarJPEG(img image.Image, path string, quality int) error {
	f, err := os.Create(path)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrStorage, "no se pudo crear el archivo de imagen", "images.guardarJPEG")
	}
	defer f.Close()

	return jpeg.Encode(f, img, &jpeg.Options{Quality: quality})
}
