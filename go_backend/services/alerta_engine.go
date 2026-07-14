package services

import (
	"searchobject/db"
	"searchobject/models"
)

type AlertaEngine struct {
	db *db.DB
}

func NewAlertaEngine(database *db.DB) *AlertaEngine {
	return &AlertaEngine{db: database}
}

func (e *AlertaEngine) Evaluar(usuarioID string) error {
	if err := e.verificarStockBajo(usuarioID); err != nil {
		return err
	}
	if err := e.verificarCajasSaturadas(usuarioID); err != nil {
		return err
	}
	return nil
}

func (e *AlertaEngine) verificarStockBajo(usuarioID string) error {
	espacios, err := e.db.ListarEspacios(usuarioID)
	if err != nil {
		return err
	}

	umbral := 5

	for _, esp := range espacios {
		cajas, err := e.db.ListarCajas(esp.ID)
		if err != nil {
			continue
		}
		for _, c := range cajas {
			objetos, err := e.db.ListarObjetos(c.ID)
			if err != nil {
				continue
			}
			for _, o := range objetos {
				if o.EsInsumo && o.Cantidad <= umbral {
					a := &models.Alerta{
						Tipo:        models.AlertaStockBajo,
						EntidadTipo: "objeto",
						EntidadID:   o.ID,
						Mensaje:     "Stock bajo de " + o.Nombre + " (quedan " + itoa(o.Cantidad) + ")",
					}
					e.db.CrearAlerta(a)
				}
			}
		}
	}
	return nil
}

func (e *AlertaEngine) verificarCajasSaturadas(usuarioID string) error {
	espacios, err := e.db.ListarEspacios(usuarioID)
	if err != nil {
		return err
	}

	for _, esp := range espacios {
		cajas, err := e.db.ListarCajas(esp.ID)
		if err != nil {
			continue
		}
		for _, c := range cajas {
			if c.CapacidadMax == nil {
				continue
			}
			objetos, err := e.db.ListarObjetos(c.ID)
			if err != nil {
				continue
			}
			if len(objetos) >= *c.CapacidadMax {
				a := &models.Alerta{
					Tipo:        models.AlertaCajaSaturada,
					EntidadTipo: "caja",
					EntidadID:   c.ID,
					Mensaje:     "La caja " + c.Nombre + " está llena (" + itoa(len(objetos)) + "/" + itoa(*c.CapacidadMax) + ")",
				}
				e.db.CrearAlerta(a)
			}
		}
	}
	return nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
