package bridge

import (
	"os"
	"path/filepath"

	"searchobject/db"
	"searchobject/services"
)

type App struct {
	db          *db.DB
	espacios    *services.EspacioService
	objetos     *services.ObjetoService
	busqueda    *services.BusquedaService
	inventario  *services.InventarioService
	alertas     *services.AlertaEngine
	reportes    *services.ReporteService
	imageDir    string
}

func New(dbPath, imageDir string) (*App, error) {
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return nil, err
	}

	database, err := db.Open(dbPath)
	if err != nil {
		return nil, err
	}

	return &App{
		db:         database,
		espacios:   services.NewEspacioService(database),
		objetos:    services.NewObjetoService(database),
		busqueda:   services.NewBusquedaService(database),
		inventario: services.NewInventarioService(database),
		alertas:    services.NewAlertaEngine(database),
		reportes:   services.NewReporteService(database),
		imageDir:   imageDir,
	}, nil
}

func (app *App) Close() error {
	if app.db != nil {
		return app.db.Close()
	}
	return nil
}

func (app *App) ImageDir() string {
	return app.imageDir
}

func (app *App) PathParaObjeto(objetoID string) string {
	return filepath.Join(app.imageDir, objetoID)
}
