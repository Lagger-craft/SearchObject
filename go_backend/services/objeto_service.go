package services

import (
	"searchobject/db"
	apperrors "searchobject/errors"
	"searchobject/models"
	"searchobject/normalize"
)

type ObjetoService struct {
	db *db.DB
}

func NewObjetoService(database *db.DB) *ObjetoService {
	return &ObjetoService{db: database}
}

type CrearObjetoReq struct {
	CajaID      string
	UsuarioID   string
	Nombre      string
	Descripcion string
	Cantidad    int
	EsInsumo    bool
	ValorEstimado *float64
}

func (s *ObjetoService) Crear(req CrearObjetoReq) (*models.Objeto, error) {
	if req.Nombre == "" {
		return nil, apperrors.New(apperrors.ErrValidation, "el nombre del objeto no puede estar vacío", "ObjetoService.Crear")
	}
	if req.Cantidad < 1 {
		req.Cantidad = 1
	}

	norm := normalize.Nombre(req.Nombre)

	exist, err := s.buscarPorNombreNorm(req.CajaID, norm)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, apperrors.New(apperrors.ErrDuplicate,
			"ya existe un objeto llamado '"+exist.Nombre+"' en esta caja", "ObjetoService.Crear")
	}

	o := &models.Objeto{
		CajaID:       req.CajaID,
		UsuarioID:    req.UsuarioID,
		Nombre:       req.Nombre,
		NombreNorm:   norm,
		Descripcion:  req.Descripcion,
		Cantidad:     req.Cantidad,
		EsInsumo:     req.EsInsumo,
		ValorEstimado: req.ValorEstimado,
		CreatedAt:    models.Now(),
		UpdatedAt:    models.Now(),
	}

	if err := s.db.CrearObjeto(o); err != nil {
		return nil, err
	}

	s.registrarMovimientoEntrada(o)

	return o, nil
}

func (s *ObjetoService) Mover(objetoID, haciaCajaID, nota string) (*models.Objeto, error) {
	o, err := s.db.ObtenerObjeto(objetoID)
	if err != nil {
		return nil, err
	}

	if o.CajaID == haciaCajaID {
		return nil, apperrors.New(apperrors.ErrValidation,
			"el objeto ya está en esa caja", "ObjetoService.Mover")
	}

	desdeID := o.CajaID
	o.CajaID = haciaCajaID
	o.UpdatedAt = models.Now()

	if err := s.db.ActualizarObjeto(o); err != nil {
		return nil, err
	}

	m := &models.Movimiento{
		ObjetoID:    objetoID,
		DesdeCajaID: &desdeID,
		HaciaCajaID: &haciaCajaID,
		Tipo:        models.MovTraslado,
		Nota:        nota,
		Fecha:       models.Now(),
		CreatedAt:   models.Now(),
	}
	s.db.CrearMovimiento(m)

	return o, nil
}

func (s *ObjetoService) Listar(cajaID string) ([]models.Objeto, error) {
	return s.db.ListarObjetos(cajaID)
}

func (s *ObjetoService) Obtener(id string) (*models.Objeto, error) {
	return s.db.ObtenerObjeto(id)
}

func (s *ObjetoService) Eliminar(id string) error {
	o, err := s.db.ObtenerObjeto(id)
	if err != nil {
		return err
	}
	if err := s.db.EliminarImagenesPorObjeto(id); err != nil {
		return err
	}
	return s.db.EliminarObjeto(o.ID)
}

func (s *ObjetoService) Buscar(usuarioID, termino string) ([]models.Objeto, error) {
	norm := normalize.Busqueda(termino)
	if norm == "" {
		return nil, apperrors.New(apperrors.ErrValidation, "ingresá un término de búsqueda", "ObjetoService.Buscar")
	}
	return s.db.BuscarObjetos(usuarioID, norm)
}

func (s *ObjetoService) buscarPorNombreNorm(cajaID, norm string) (*models.Objeto, error) {
	objetos, err := s.db.ListarObjetos(cajaID)
	if err != nil {
		return nil, err
	}
	for _, o := range objetos {
		if o.NombreNorm == norm {
			return &o, nil
		}
	}
	return nil, nil
}

func (s *ObjetoService) registrarMovimientoEntrada(o *models.Objeto) {
	m := &models.Movimiento{
		ObjetoID:   o.ID,
		HaciaCajaID: &o.CajaID,
		Tipo:       models.MovEntrada,
		Fecha:      o.CreatedAt,
		CreatedAt:  o.CreatedAt,
	}
	s.db.CrearMovimiento(m)
}
