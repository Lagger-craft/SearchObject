package bridge

import (
	"encoding/json"
)

type itemStockBajoDTO struct {
	ObjetoID  string `json:"objeto_id"`
	Nombre    string `json:"nombre"`
	Cantidad  int    `json:"cantidad"`
	Ubicacion string `json:"ubicacion"`
}

type listarStockBajoRes struct {
	Items []itemStockBajoDTO `json:"items"`
}

func (app *App) StockBajo(usuarioID string, limite int) (string, error) {
	items, err := app.inventario.StockBajo(usuarioID, limite)
	if err != nil {
		return "", err
	}

	dtos := make([]itemStockBajoDTO, 0, len(items))
	for _, item := range items {
		dtos = append(dtos, itemStockBajoDTO{
			ObjetoID:  item.Objeto.ID,
			Nombre:    item.Objeto.Nombre,
			Cantidad:  item.Objeto.Cantidad,
			Ubicacion: item.Ubicacion,
		})
	}

	res := listarStockBajoRes{Items: dtos}
	data, _ := json.Marshal(res)
	return string(data), nil
}
