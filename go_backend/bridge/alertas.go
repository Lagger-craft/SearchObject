package bridge

import (
	"encoding/json"
)

func (app *App) Resumen(usuarioID string) (string, error) {
	resumen, err := app.inventario.Resumen(usuarioID)
	if err != nil {
		return "", err
	}
	data, _ := json.Marshal(resumen)
	return string(data), nil
}

func (app *App) Dashboard(usuarioID string) (string, error) {
	stats, err := app.inventario.Dashboard(usuarioID)
	if err != nil {
		return "", err
	}
	data, _ := json.Marshal(stats)
	return string(data), nil
}

func (app *App) EvaluarAlertas(usuarioID string) (string, error) {
	if err := app.alertas.Evaluar(usuarioID); err != nil {
		return "", err
	}
	return `{"ok":true}`, nil
}

func (app *App) ListarAlertas(jsonLeidas string) (string, error) {
	var leidas *bool
	if jsonLeidas != "" {
		v := jsonLeidas == "true"
		leidas = &v
	}

	alertas, err := app.db.ListarAlertas(leidas)
	if err != nil {
		return "", err
	}

	data, _ := json.Marshal(alertas)
	return string(data), nil
}

func (app *App) MarcarAlertaLeida(id string) (string, error) {
	if err := app.db.MarcarAlertaLeida(id); err != nil {
		return "", err
	}
	return `{"ok":true}`, nil
}

func (app *App) ResolverAlerta(id string) (string, error) {
	if err := app.db.ResolverAlerta(id); err != nil {
		return "", err
	}
	return `{"ok":true}`, nil
}

func (app *App) ExportarJSON(usuarioID string) (string, error) {
	return app.reportes.ExportarJSON(usuarioID)
}

func (app *App) ExportarCSV(usuarioID string) (string, error) {
	return app.reportes.ExportarCSV(usuarioID)
}
