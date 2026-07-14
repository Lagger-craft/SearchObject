package services

import (
	"testing"

	"searchobject/db"
	"searchobject/models"
)

type testEnv struct {
	db       *db.DB
	userID   string
	espacios *EspacioService
	objetos  *ObjetoService
	alertas  *AlertaEngine
	reportes *ReporteService
}

func setup(t *testing.T) *testEnv {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("Error abriendo DB: %v", err)
	}
	t.Cleanup(func() { database.Close() })

	u := &models.Usuario{
		ID:        "user-test",
		Nombre:    "Test",
		CreatedAt: models.Now(),
		UpdatedAt: models.Now(),
	}
	if err := database.CrearUsuario(u); err != nil {
		t.Fatalf("Error creando usuario: %v", err)
	}

	return &testEnv{
		db:       database,
		userID:   u.ID,
		espacios: NewEspacioService(database),
		objetos:  NewObjetoService(database),
		alertas:  NewAlertaEngine(database),
		reportes: NewReporteService(database),
	}
}

func TestCrearEspacioConNormalizacion(t *testing.T) {
	env := setup(t)

	e1, err := env.espacios.Crear(env.userID, "Garage", "", nil)
	if err != nil {
		t.Fatalf("Error creando espacio: %v", err)
	}
	if e1.Nombre != "Garage" {
		t.Errorf("Nombre debería preservarse = %q", e1.Nombre)
	}

	// Duplicado con tilde y mayúsculas
	_, err = env.espacios.Crear(env.userID, "GáraGe  ", "", nil)
	if err == nil {
		t.Fatal("Debería fallar por duplicado (GáraGe normaliza a garage)")
	}
}

func TestCrearObjetoConNormalizacion(t *testing.T) {
	env := setup(t)

	esp, _ := env.espacios.Crear(env.userID, "Test", "", nil)

	c := &models.Caja{
		EspacioID:  esp.ID,
		UsuarioID:  env.userID,
		Nombre:     "Caja Test",
		NombreNorm: "caja test",
		CreatedAt:  models.Now(),
		UpdatedAt:  models.Now(),
	}
	env.db.CrearCaja(c)

	o1, err := env.objetos.Crear(CrearObjetoReq{
		CajaID:    c.ID,
		UsuarioID: env.userID,
		Nombre:    "Taladro Percutor",
		Cantidad:  1,
	})
	if err != nil {
		t.Fatalf("Error creando objeto: %v", err)
	}
	if o1.Nombre != "Taladro Percutor" {
		t.Errorf("Nombre = %q; want Taladro Percutor", o1.Nombre)
	}

	// Buscar con tilde
	resultados, err := env.objetos.Buscar(env.userID, "táladro")
	if err != nil {
		t.Fatalf("Error buscando: %v", err)
	}
	if len(resultados) != 1 {
		t.Errorf("Esperaba 1 resultado, tengo %d", len(resultados))
	}
}

func TestMoverObjeto(t *testing.T) {
	env := setup(t)

	esp, _ := env.espacios.Crear(env.userID, "Test", "", nil)

	c1 := &models.Caja{EspacioID: esp.ID, UsuarioID: env.userID, Nombre: "Caja 1", NombreNorm: "caja 1"}
	c2 := &models.Caja{EspacioID: esp.ID, UsuarioID: env.userID, Nombre: "Caja 2", NombreNorm: "caja 2"}
	env.db.CrearCaja(c1)
	env.db.CrearCaja(c2)

	o, _ := env.objetos.Crear(CrearObjetoReq{
		CajaID: c1.ID, UsuarioID: env.userID, Nombre: "Test", Cantidad: 1,
	})

	movido, err := env.objetos.Mover(o.ID, c2.ID, "Traslado de prueba")
	if err != nil {
		t.Fatalf("Error moviendo: %v", err)
	}
	if movido.CajaID != c2.ID {
		t.Errorf("CajaID = %q; want %q", movido.CajaID, c2.ID)
	}

	// Mover al mismo lugar debe fallar
	_, err = env.objetos.Mover(o.ID, c2.ID, "")
	if err == nil {
		t.Fatal("Mover al mismo lugar debería fallar")
	}
}

func TestAlertasStockBajo(t *testing.T) {
	env := setup(t)

	esp, _ := env.espacios.Crear(env.userID, "Test", "", nil)
	c := &models.Caja{EspacioID: esp.ID, UsuarioID: env.userID, Nombre: "Caja", NombreNorm: "caja"}
	env.db.CrearCaja(c)

	// Insumo con cantidad baja
	env.objetos.Crear(CrearObjetoReq{
		CajaID: c.ID, UsuarioID: env.userID,
		Nombre: "Tornillos 5mm", Cantidad: 3, EsInsumo: true,
	})

	// Objeto normal (no insumo)
	env.objetos.Crear(CrearObjetoReq{
		CajaID: c.ID, UsuarioID: env.userID,
		Nombre: "Martillo", Cantidad: 1, EsInsumo: false,
	})

	if err := env.alertas.Evaluar(env.userID); err != nil {
		t.Fatalf("Error evaluando alertas: %v", err)
	}

	alertas, _ := env.db.ListarAlertas(nil)
	if len(alertas) == 0 {
		t.Fatal("Debería haber al menos 1 alerta de stock bajo")
	}
	if alertas[0].Tipo != models.AlertaStockBajo {
		t.Errorf("Tipo = %q; want stock_bajo", alertas[0].Tipo)
	}
}

func TestReporteJSON(t *testing.T) {
	env := setup(t)

	esp, _ := env.espacios.Crear(env.userID, "Test", "", nil)
	c := &models.Caja{EspacioID: esp.ID, UsuarioID: env.userID, Nombre: "Caja", NombreNorm: "caja"}
	env.db.CrearCaja(c)
	env.objetos.Crear(CrearObjetoReq{
		CajaID: c.ID, UsuarioID: env.userID,
		Nombre: "Objeto 1", Cantidad: 1,
	})

	json, err := env.reportes.ExportarJSON(env.userID)
	if err != nil {
		t.Fatalf("Error exportando JSON: %v", err)
	}
	if len(json) == 0 {
		t.Fatal("Reporte JSON vacío")
	}
}

func TestDashboard(t *testing.T) {
	env := setup(t)

	inv := NewInventarioService(env.db)
	stats, err := inv.Dashboard(env.userID)
	if err != nil {
		t.Fatalf("Error en dashboard: %v", err)
	}
	if stats.TotalEspacios < 0 {
		t.Error("TotalEspacios negativo")
	}
}
