package services

import (
	"time"

	"searchobject/db"
	"searchobject/models"
)

type InventarioService struct {
	db *db.DB
}

func NewInventarioService(database *db.DB) *InventarioService {
	return &InventarioService{db: database}
}

type ResumenInventario struct {
	TotalObjetos  int     `json:"total_objetos"`
	ValorEstimado float64 `json:"valor_estimado"`
	Faltantes     int     `json:"faltantes"`
	StockBajo     int     `json:"stock_bajo"`
	Prestamos     int     `json:"prestamos"`
	CajasSaturadas int    `json:"cajas_saturadas"`
}

func (s *InventarioService) Resumen(usuarioID string) (*ResumenInventario, error) {
	r := &ResumenInventario{}

	alertas, err := s.db.ListarAlertas(nil)
	if err != nil {
		return nil, err
	}

	for _, a := range alertas {
		if a.ResueltaAt != nil {
			continue
		}
		switch a.Tipo {
		case models.AlertaFaltante:
			r.Faltantes++
		case models.AlertaStockBajo:
			r.StockBajo++
		case models.AlertaPrestamoVencido:
			r.Prestamos++
		case models.AlertaCajaSaturada:
			r.CajasSaturadas++
		}
	}

	return r, nil
}

func (s *InventarioService) VerificarFaltantes(usuarioID string) ([]models.Alerta, error) {
	espacios, err := s.db.ListarEspacios(usuarioID)
	if err != nil {
		return nil, err
	}

	var alertas []models.Alerta

	for _, esp := range espacios {
		cajas, err := s.db.ListarCajas(esp.ID)
		if err != nil {
			continue
		}
		for _, c := range cajas {
			objetos, err := s.db.ListarObjetos(c.ID)
			if err != nil {
				continue
			}
			for _, o := range objetos {
				movs, err := s.db.ListarMovimientos(o.ID)
				if err != nil || len(movs) == 0 {
					continue
				}
				ultimoMov := movs[0]
				if ultimoMov.Tipo == models.MovSalida {
					a := models.Alerta{
						Tipo:        models.AlertaFaltante,
						EntidadTipo: "objeto",
						EntidadID:   o.ID,
						Mensaje:     "El objeto " + o.Nombre + " fue retirado y no se registró su ubicación actual",
					}
					alertas = append(alertas, a)
				}
			}
		}
	}

	return alertas, nil
}

type ItemStockBajo struct {
	Objeto    models.Objeto `json:"objeto"`
	Ubicacion string        `json:"ubicacion"`
}

func (s *InventarioService) StockBajo(usuarioID string, limite int) ([]ItemStockBajo, error) {
	if limite <= 0 {
		limite = 5
	}

	espacios, err := s.db.ListarEspacios(usuarioID)
	if err != nil {
		return nil, err
	}

	var items []ItemStockBajo

	for _, esp := range espacios {
		cajas, err := s.db.ListarCajas(esp.ID)
		if err != nil {
			continue
		}
		for _, c := range cajas {
			objetos, err := s.db.ListarObjetos(c.ID)
			if err != nil {
				continue
			}
			for _, o := range objetos {
				if o.EsInsumo && o.Cantidad <= limite {
					items = append(items, ItemStockBajo{
						Objeto:    o,
						Ubicacion: esp.Nombre + " > " + c.Nombre,
					})
				}
			}
		}
	}

	return items, nil
}

type StatsDashboard struct {
	TotalObjetos    int     `json:"total_objetos"`
	TotalEspacios   int     `json:"total_espacios"`
	TotalCajas      int     `json:"total_cajas"`
	ValorEstimado   float64 `json:"valor_estimado"`
	MovimientosHoy  int     `json:"movimientos_hoy"`
	AlertasNoLeidas int     `json:"alertas_no_leidas"`
}

func (s *InventarioService) Dashboard(usuarioID string) (*StatsDashboard, error) {
	stats := &StatsDashboard{}

	espacios, err := s.db.ListarEspacios(usuarioID)
	if err != nil {
		return nil, err
	}
	stats.TotalEspacios = len(espacios)

	for _, esp := range espacios {
		cajas, err := s.db.ListarCajas(esp.ID)
		if err != nil {
			continue
		}
		stats.TotalCajas += len(cajas)

		for _, c := range cajas {
			objetos, err := s.db.ListarObjetos(c.ID)
			if err != nil {
				continue
			}
			stats.TotalObjetos += len(objetos)
			for _, o := range objetos {
				if o.ValorEstimado != nil {
					stats.ValorEstimado += *o.ValorEstimado
				}
			}
		}
	}

	movs, err := s.db.MovimientosRecientes(usuarioID, 100)
	if err == nil {
		hoy := time.Now().Truncate(24 * time.Hour)
		for _, m := range movs {
			if m.Fecha.After(hoy) {
				stats.MovimientosHoy++
			}
		}
	}

	count, err := s.db.AlertasNoLeidasCount()
	if err == nil {
		stats.AlertasNoLeidas = count
	}

	return stats, nil
}
