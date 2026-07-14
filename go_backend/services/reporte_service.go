package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"searchobject/db"
	"searchobject/models"
)

type ReporteService struct {
	db *db.DB
}

func NewReporteService(database *db.DB) *ReporteService {
	return &ReporteService{db: database}
}

type ReporteInventario struct {
	GeneradoEn time.Time            `json:"generado_en"`
	Espacios   []ReporteEspacio     `json:"espacios"`
	Resumen    ResumenReporte       `json:"resumen"`
}

type ReporteEspacio struct {
	Espacio models.Espacio `json:"espacio"`
	Cajas   []ReporteCaja  `json:"cajas"`
}

type ReporteCaja struct {
	Caja    models.Caja     `json:"caja"`
	Objetos []models.Objeto `json:"objetos"`
}

type ResumenReporte struct {
	TotalEspacios int     `json:"total_espacios"`
	TotalCajas    int     `json:"total_cajas"`
	TotalObjetos  int     `json:"total_objetos"`
	ValorTotal    float64 `json:"valor_total"`
}

func (s *ReporteService) Generar(usuarioID string) (*ReporteInventario, error) {
	reporte := &ReporteInventario{
		GeneradoEn: time.Now(),
	}

	espacios, err := s.db.ListarEspacios(usuarioID)
	if err != nil {
		return nil, fmt.Errorf("listar espacios: %w", err)
	}

	for _, esp := range espacios {
		re := ReporteEspacio{Espacio: esp}

		cajas, err := s.db.ListarCajas(esp.ID)
		if err != nil {
			continue
		}

		for _, c := range cajas {
			rc := ReporteCaja{Caja: c}
			objetos, err := s.db.ListarObjetos(c.ID)
			if err != nil {
				continue
			}
			rc.Objetos = objetos
			re.Cajas = append(re.Cajas, rc)
		}

		reporte.Espacios = append(reporte.Espacios, re)
	}

	s.calcularResumen(reporte)
	return reporte, nil
}

func (s *ReporteService) calcularResumen(r *ReporteInventario) {
	r.Resumen.TotalEspacios = len(r.Espacios)
	for _, e := range r.Espacios {
		r.Resumen.TotalCajas += len(e.Cajas)
		for _, c := range e.Cajas {
			r.Resumen.TotalObjetos += len(c.Objetos)
			for _, o := range c.Objetos {
				if o.ValorEstimado != nil {
					r.Resumen.ValorTotal += *o.ValorEstimado
				}
			}
		}
	}
}

func (s *ReporteService) ExportarJSON(usuarioID string) (string, error) {
	reporte, err := s.Generar(usuarioID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(reporte, "", "  ")
	if err != nil {
		return "", fmt.Errorf("serializar JSON: %w", err)
	}
	return string(data), nil
}

func (s *ReporteService) ExportarCSV(usuarioID string) (string, error) {
	reporte, err := s.Generar(usuarioID)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	buf.WriteString("Espacio,Caja,Objeto,Cantidad,Tipo,Valor\n")

	for _, e := range reporte.Espacios {
		for _, c := range e.Cajas {
			for _, o := range c.Objetos {
				tipo := "Herramienta"
				if o.EsInsumo {
					tipo = "Insumo"
				}
				valor := ""
				if o.ValorEstimado != nil {
					valor = fmt.Sprintf("%.2f", *o.ValorEstimado)
				}
				buf.WriteString(fmt.Sprintf("%s,%s,%s,%d,%s,%s\n",
					escapeCSV(e.Espacio.Nombre),
					escapeCSV(c.Caja.Nombre),
					escapeCSV(o.Nombre),
					o.Cantidad, tipo, valor))
			}
		}
	}

	return buf.String(), nil
}

func escapeCSV(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		s = strings.ReplaceAll(s, "\"", "\"\"")
		s = "\"" + s + "\""
	}
	return s
}
