package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"searchobject/bridge"
)

func main() {
	tmpDir, _ := os.MkdirTemp("", "searchobject-test")
	defer os.RemoveAll(tmpDir)

	app, err := bridge.New(filepath.Join(tmpDir, "test.db"), filepath.Join(tmpDir, "images"))
	if err != nil {
		log.Fatalf("❌ Error al iniciar: %v", err)
	}
	defer app.Close()

	fmt.Println("🔧 SearchObject — Test End-to-End")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// ─── 1. USUARIO ───
	fmt.Print("📌 Creando usuario... ")
	uj, err := app.CrearUsuario("Lager", "lager@test.com")
	must(err)
	userID := jval(uj, "id")
	fmt.Println("✅", userID)

	// ─── 2. ESPACIOS ───
	fmt.Print("📌 Creando espacios... ")
	_, err = app.CrearEspacio(jobj("usuario_id", userID, "nombre", "Casa"))
	must(err)
	garage, err := app.CrearEspacio(jobj("usuario_id", userID, "nombre", "GARAGE"))
	must(err)

	_, err = app.CrearEspacio(jobj("usuario_id", userID, "nombre", "Casa"))
	if err == nil { log.Fatal("❌ duplicado Casa debería fallar") }
	fmt.Print("🚫Casa ")

	_, err = app.CrearEspacio(jobj("usuario_id", userID, "nombre", "  CáSa  "))
	if err == nil { log.Fatal("❌ duplicado CáSa debería fallar") }
	fmt.Print("🚫CáSa ")

	_, err = app.CrearEspacio(jobj("usuario_id", userID, "nombre", "Estante A", "padre_id", jval(garage, "id")))
	must(err)

	lj, _ := app.ListarEspacios(userID)
	espacios := jarr(lj, "espacios")
	fmt.Printf(" → %d espacios ✅\n", len(espacios))

	// ─── 3. CAJAS ───
	fmt.Print("📌 Cajas... ")
	gID := jval(garage, "id")
	tornillos, err := app.CrearCaja(jobj("espacio_id", gID, "usuario_id", userID, "nombre", "Caja Tornillos"))
	must(err)

	_, err = app.CrearCaja(jobj("espacio_id", gID, "usuario_id", userID, "nombre", "caja tornillos"))
	if err == nil { log.Fatal("❌ duplicado") }
	fmt.Print("🚫dup ")

	herr, err := app.CrearCaja(jobj("espacio_id", gID, "usuario_id", userID, "nombre", "Caja Herramientas"))
	must(err)
	fmt.Println("✅")

	// ─── 4. OBJETOS ───
	fmt.Print("📌 Objetos... ")
	tID := jval(tornillos, "id")
	_, err = app.CrearObjeto(jobj("caja_id", tID, "usuario_id", userID, "nombre", "Taladro Percutor", "cantidad", 1, "valor", 150.0))
	must(err)
	_, err = app.CrearObjeto(jobj("caja_id", tID, "usuario_id", userID, "nombre", "taladro percutor"))
	if err == nil { log.Fatal("❌ duplicado") }
	fmt.Print("🚫dup ")

	_, err = app.CrearObjeto(jobj("caja_id", tID, "usuario_id", userID, "nombre", "Tornillos 5mm", "cantidad", 3, "es_insumo", true))
	must(err)
	_, err = app.CrearObjeto(jobj("caja_id", tID, "usuario_id", userID, "nombre", "Mártillo", "cantidad", 1))
	must(err)
	fmt.Println("✅")

	// ─── 5. BÚSQUEDA (acentos!) ───
	fmt.Print("🔍 Buscando... ")
	s1, _ := app.Buscar(userID, "taladro")
	r1 := jarr(s1, "resultados")

	s2, _ := app.Buscar(userID, "martillo")
	r2 := jarr(s2, "resultados")

	s3, _ := app.Buscar(userID, "Mártillo")
	r3 := jarr(s3, "resultados")
	fmt.Printf("'taladro':%d 'martillo':%d 'Mártillo':%d ✅\n", len(r1), len(r2), len(r3))

	// ─── 6. MOVER ───
	fmt.Print("📦 Mover objeto... ")
	obs, _ := app.ListarObjetos(tID)
	obsArr := jarr(obs, "objetos")
	talID := obsArr[0].(map[string]interface{})["id"].(string)

	herrID := jval(herr, "id")

	_, err = app.MoverObjeto(jobj("id", talID, "caja_id", herrID, "nota", "Movido a herramientas"))
	must(err)
	fmt.Println("✅")

	// ─── 7. ALERTAS ───
	fmt.Print("🔔 Alertas... ")
	_, err = app.EvaluarAlertas(userID)
	must(err)
	aj, _ := app.ListarAlertas("")
	var alertas []interface{}
	json.Unmarshal([]byte(aj), &alertas)
	fmt.Printf("%d alertas ✅\n", len(alertas))

	// ─── 8. REPORTES ───
	fmt.Print("📊 Reportes... ")
	rj, err := app.ExportarJSON(userID)
	must(err)
	var res map[string]interface{}
	json.Unmarshal([]byte(rj), &res)
	total := res["resumen"].(map[string]interface{})["total_objetos"]
	fmt.Printf("JSON:%v ", total)

	csv, err := app.ExportarCSV(userID)
	must(err)
	fmt.Printf("CSV:%dbytes ✅\n", len(csv))

	// ─── 9. DASHBOARD ───
	fmt.Print("📈 Dashboard... ")
	dj, err := app.Dashboard(userID)
	must(err)
	var dash map[string]interface{}
	json.Unmarshal([]byte(dj), &dash)
	objs := dash["total_objetos"]
	fmt.Printf("✅ %v objetos\n", objs)

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎉 TODO OK")
}

// Helpers JSON rápidos
func jval(s, key string) string {
	var m map[string]interface{}
	json.Unmarshal([]byte(s), &m)
	v, _ := m[key].(string)
	return v
}

func jarr(s, key string) []interface{} {
	var m map[string]interface{}
	json.Unmarshal([]byte(s), &m)
	a, _ := m[key].([]interface{})
	return a
}

func jobj(kv ...interface{}) string {
	m := make(map[string]interface{})
	for i := 0; i < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func must(err error) {
	if err != nil {
		log.Fatalf("❌ %v", err)
	}
}
