package db

import (
	"testing"

	"searchobject/models"
)

func setupDB(t *testing.T) *DB {
	t.Helper()
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Error abriendo DB: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func crearUsuario(t *testing.T, db *DB) *models.Usuario {
	t.Helper()
	u := &models.Usuario{
		ID:        "user-1",
		Nombre:    "Test",
		Email:     "test@test.com",
		CreatedAt: models.Now(),
		UpdatedAt: models.Now(),
	}
	if err := db.CrearUsuario(u); err != nil {
		t.Fatalf("Error creando usuario: %v", err)
	}
	return u
}

func TestCrearYObtenerUsuario(t *testing.T) {
	db := setupDB(t)
	u := crearUsuario(t, db)

	got, err := db.ObtenerUsuario(u.ID)
	if err != nil {
		t.Fatalf("Error obteniendo usuario: %v", err)
	}
	if got.Nombre != "Test" {
		t.Errorf("Nombre = %q; want Test", got.Nombre)
	}
}

func TestDuplicadoEspacio(t *testing.T) {
	db := setupDB(t)
	_ = crearUsuario(t, db)

	e1 := &models.Espacio{
		ID:         "esp-1",
		UsuarioID:  "user-1",
		Nombre:     "Garage",
		NombreNorm: "garage",
		CreatedAt:  models.Now(),
		UpdatedAt:  models.Now(),
	}
	if err := db.CrearEspacio(e1); err != nil {
		t.Fatalf("Error creando espacio: %v", err)
	}

	e2 := &models.Espacio{
		ID:         "esp-2",
		UsuarioID:  "user-1",
		Nombre:     "garage",
		NombreNorm: "garage",
		CreatedAt:  models.Now(),
		UpdatedAt:  models.Now(),
	}
	if err := db.CrearEspacio(e2); err == nil {
		t.Fatal("Esperaba error de duplicado, pero no lo hubo")
	}
}

func TestJerarquiaEspacios(t *testing.T) {
	db := setupDB(t)
	_ = crearUsuario(t, db)

	padre := "esp-padre"
	db.CrearEspacio(&models.Espacio{
		ID:         padre,
		UsuarioID:  "user-1",
		Nombre:     "Casa",
		NombreNorm: "casa",
		CreatedAt:  models.Now(),
		UpdatedAt:  models.Now(),
	})

	hijo := &models.Espacio{
		ID:         "esp-hijo",
		UsuarioID:  "user-1",
		Nombre:     "Garage",
		NombreNorm: "garage",
		PadreID:    &padre,
		CreatedAt:  models.Now(),
		UpdatedAt:  models.Now(),
	}
	if err := db.CrearEspacio(hijo); err != nil {
		t.Fatalf("Error creando subespacio: %v", err)
	}

	espacios, _ := db.ListarEspacios("user-1")
	var encontrado bool
	for _, e := range espacios {
		if e.PadreID != nil && *e.PadreID == padre {
			encontrado = true
			break
		}
	}
	if !encontrado {
		t.Error("No se encontró el subespacio")
	}
}

func TestBuscarObjetos(t *testing.T) {
	db := setupDB(t)
	_ = crearUsuario(t, db)

	db.CrearEspacio(&models.Espacio{
		ID: "esp", UsuarioID: "user-1", Nombre: "Test", NombreNorm: "test",
		CreatedAt: models.Now(), UpdatedAt: models.Now(),
	})
	db.CrearCaja(&models.Caja{
		ID: "caja", EspacioID: "esp", UsuarioID: "user-1",
		Nombre: "Caja", NombreNorm: "caja",
		CreatedAt: models.Now(), UpdatedAt: models.Now(),
	})
	db.CrearObjeto(&models.Objeto{
		ID: "obj-1", CajaID: "caja", UsuarioID: "user-1",
		Nombre: "Taladro", NombreNorm: "taladro",
		Cantidad: 1, CreatedAt: models.Now(), UpdatedAt: models.Now(),
	})
	db.CrearObjeto(&models.Objeto{
		ID: "obj-2", CajaID: "caja", UsuarioID: "user-1",
		Nombre: "Martillo", NombreNorm: "martillo",
		Cantidad: 1, CreatedAt: models.Now(), UpdatedAt: models.Now(),
	})

	resultados, err := db.BuscarObjetos("user-1", "taladro")
	if err != nil {
		t.Fatalf("Error buscando: %v", err)
	}
	if len(resultados) != 1 {
		t.Errorf("Esperaba 1 resultado, tengo %d", len(resultados))
	}
	if resultados[0].ID != "obj-1" {
		t.Errorf("ID = %q; want obj-1", resultados[0].ID)
	}
}

func TestMovimientos(t *testing.T) {
	db := setupDB(t)
	_ = crearUsuario(t, db)

	db.CrearEspacio(&models.Espacio{
		ID: "esp", UsuarioID: "user-1", Nombre: "X", NombreNorm: "x",
		CreatedAt: models.Now(), UpdatedAt: models.Now(),
	})
	db.CrearCaja(&models.Caja{
		ID: "caja-1", EspacioID: "esp", UsuarioID: "user-1",
		Nombre: "Caja 1", NombreNorm: "caja 1",
		CreatedAt: models.Now(), UpdatedAt: models.Now(),
	})
	db.CrearCaja(&models.Caja{
		ID: "caja-2", EspacioID: "esp", UsuarioID: "user-1",
		Nombre: "Caja 2", NombreNorm: "caja 2",
		CreatedAt: models.Now(), UpdatedAt: models.Now(),
	})
	db.CrearObjeto(&models.Objeto{
		ID: "obj", CajaID: "caja-1", UsuarioID: "user-1",
		Nombre: "Test", NombreNorm: "test", Cantidad: 1,
		CreatedAt: models.Now(), UpdatedAt: models.Now(),
	})

	c2 := "caja-2"
	db.CrearMovimiento(&models.Movimiento{
		ID: "mov-1", ObjetoID: "obj",
		DesdeCajaID: nil, HaciaCajaID: &c2,
		Tipo: models.MovTraslado, Fecha: models.Now(),
		CreatedAt: models.Now(),
	})

	movs, err := db.ListarMovimientos("obj")
	if err != nil {
		t.Fatalf("Error listando movimientos: %v", err)
	}
	if len(movs) != 1 {
		t.Errorf("Esperaba 1 movimiento, tengo %d", len(movs))
	}
}

func TestAlertas(t *testing.T) {
	db := setupDB(t)

	db.CrearAlerta(&models.Alerta{
		ID: "al-1", Tipo: models.AlertaStockBajo,
		EntidadTipo: "objeto", EntidadID: "obj-1",
		Mensaje: "Stock bajo de Tornillos", CreatedAt: models.Now(),
	})
	db.CrearAlerta(&models.Alerta{
		ID: "al-2", Tipo: models.AlertaCajaSaturada,
		EntidadTipo: "caja", EntidadID: "caja-1",
		Mensaje: "Caja llena", CreatedAt: models.Now(),
	})

	alertas, err := db.ListarAlertas(nil)
	if err != nil {
		t.Fatalf("Error listando alertas: %v", err)
	}
	if len(alertas) != 2 {
		t.Errorf("Esperaba 2 alertas, tengo %d", len(alertas))
	}

	db.MarcarAlertaLeida("al-1")
	alertas, _ = db.ListarAlertas(boolPtr(false))
	if len(alertas) != 1 {
		t.Errorf("Esperaba 1 alerta no leída, tengo %d", len(alertas))
	}
}

func boolPtr(b bool) *bool {
	return &b
}
