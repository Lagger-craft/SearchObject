package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"searchobject/bridge"
)

var app *bridge.App
var imageDir string

func main() {
	// Usar directorio persistente en vez de MkdirTemp
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".searchobject")
	os.MkdirAll(dataDir, 0755)

	dbPath := filepath.Join(dataDir, "searchobject.db")
	imageDir = filepath.Join(dataDir, "images")
	os.MkdirAll(imageDir, 0755)

	log.Printf("📁 Datos en: %s", dataDir)

	var err error
	app, err = bridge.New(dbPath, imageDir)
	if err != nil {
		log.Fatalf("❌ Error al iniciar: %v", err)
	}
	defer app.Close()

	mux := http.NewServeMux()

	// Ruta universal: POST /api/{method}
	// El body es JSON según lo que espere cada método del bridge
	mux.HandleFunc("/api/", apiHandler)

	// Servir imágenes: GET /images/{objetoID}/{filename}
	mux.HandleFunc("/images/", imageHandler)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	addr := ":8080"
	log.Printf("🚀 Servidor HTTP escuchando en http://localhost%s", addr)
	log.Printf("   Dashboard: GET /health")
	log.Printf("   API:       POST /api/CrearEspacio (body: JSON)")
	log.Printf("   API:       POST /api/ListarEspacios (body: {\"usuario_id\":\"...\"})")
	log.Fatal(http.ListenAndServe(addr, mux))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"método no permitido, use POST"}`, 405)
		return
	}

	// Extraer método de /api/CrearEspacio → "CrearEspacio"
	method := r.URL.Path[len("/api/"):]
	if method == "" {
		http.Error(w, `{"error":"método requerido en /api/{method}"}`, 400)
		return
	}

	// Leer body
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		body = make(map[string]interface{})
	}

	// Convertir a JSON string
	bodyBytes, _ := json.Marshal(body)

	w.Header().Set("Content-Type", "application/json")

	switch method {
	// Usuarios
	case "CrearUsuario":
		result, err := app.CrearUsuario(stringField(body, "nombre"), stringField(body, "email"))
		respond(w, result, err)
	case "ObtenerUsuario":
		result, err := app.ObtenerUsuario(stringField(body, "id"))
		respond(w, result, err)
	case "ListarUsuarios":
		result, err := app.ListarUsuarios()
		respond(w, result, err)
	case "ActualizarUsuario":
		result, err := app.ActualizarUsuario(string(bodyBytes))
		respond(w, result, err)

	// Espacios
	case "CrearEspacio":
		result, err := app.CrearEspacio(string(bodyBytes))
		respond(w, result, err)
	case "ListarEspacios":
		result, err := app.ListarEspacios(stringField(body, "usuario_id"))
		respond(w, result, err)
	case "ObtenerEspacio":
		result, err := app.ObtenerEspacio(stringField(body, "id"))
		respond(w, result, err)
	case "ActualizarEspacio":
		result, err := app.ActualizarEspacio(string(bodyBytes))
		respond(w, result, err)
	case "EliminarEspacio":
		result, err := app.EliminarEspacio(stringField(body, "id"))
		respond(w, result, err)

	// Cajas
	case "CrearCaja":
		result, err := app.CrearCaja(string(bodyBytes))
		respond(w, result, err)
	case "ListarCajas":
		result, err := app.ListarCajas(stringField(body, "espacio_id"))
		respond(w, result, err)
	case "ObtenerCaja":
		result, err := app.ObtenerCaja(stringField(body, "id"))
		respond(w, result, err)
	case "EliminarCaja":
		result, err := app.EliminarCaja(stringField(body, "id"))
		respond(w, result, err)

	// Objetos
	case "CrearObjeto":
		result, err := app.CrearObjeto(string(bodyBytes))
		respond(w, result, err)
	case "ListarObjetos":
		result, err := app.ListarObjetos(stringField(body, "caja_id"))
		respond(w, result, err)
	case "ObtenerObjeto":
		result, err := app.ObtenerObjeto(stringField(body, "id"))
		respond(w, result, err)
	case "MoverObjeto":
		result, err := app.MoverObjeto(string(bodyBytes))
		respond(w, result, err)
	case "EliminarObjeto":
		result, err := app.EliminarObjeto(stringField(body, "id"))
		respond(w, result, err)
	case "BuscarObjetos":
		result, err := app.BuscarObjetos(stringField(body, "usuario_id"), stringField(body, "termino"))
		respond(w, result, err)

	// Búsqueda completa
	case "Buscar":
		result, err := app.Buscar(stringField(body, "usuario_id"), stringField(body, "termino"))
		respond(w, result, err)

	// Imágenes
	case "AgregarImagen":
		// Aceptar base64 en image_bytes del JSON
		b64 := stringField(body, "image_bytes")
		var imgBytes []byte
		if b64 != "" {
			decoded, err := base64.StdEncoding.DecodeString(b64)
			if err != nil {
				respond(w, "", fmt.Errorf("image_bytes base64 inválido: %w", err))
				return
			}
			imgBytes = decoded
		}
		result, err := app.AgregarImagen(stringField(body, "objeto_id"), imgBytes)
		respond(w, result, err)
	case "ListarImagenes":
		result, err := app.ListarImagenes(stringField(body, "objeto_id"))
		respond(w, result, err)
	case "EliminarImagen":
		result, err := app.EliminarImagen(stringField(body, "id"))
		respond(w, result, err)

	// Alertas
	case "EvaluarAlertas":
		result, err := app.EvaluarAlertas(stringField(body, "usuario_id"))
		respond(w, result, err)
	case "ListarAlertas":
		result, err := app.ListarAlertas(stringField(body, "leidas"))
		respond(w, result, err)
	case "MarcarAlertaLeida":
		result, err := app.MarcarAlertaLeida(stringField(body, "id"))
		respond(w, result, err)
	case "ResolverAlerta":
		result, err := app.ResolverAlerta(stringField(body, "id"))
		respond(w, result, err)

	// Dashboard y reportes
	case "Resumen":
		result, err := app.Resumen(stringField(body, "usuario_id"))
		respond(w, result, err)
	case "Dashboard":
		result, err := app.Dashboard(stringField(body, "usuario_id"))
		respond(w, result, err)
	case "ExportarJSON":
		result, err := app.ExportarJSON(stringField(body, "usuario_id"))
		respond(w, result, err)
	case "ExportarCSV":
		result, err := app.ExportarCSV(stringField(body, "usuario_id"))
		respond(w, result, err)

	default:
		http.Error(w, `{"error":"método '`+method+`' no encontrado"}`, 404)
	}
}

func respond(w http.ResponseWriter, result string, err error) {
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	w.Write([]byte(result))
}

func stringField(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// imageHandler sirve archivos de imagen: GET /images/{objetoID}/{filename}
func imageHandler(w http.ResponseWriter, r *http.Request) {
	// /images/objetoID/imagenID.jpg
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/images/"), "/", 2)
	if len(parts) != 2 {
		http.Error(w, "ruta inválida", 400)
		return
	}

	objetoID := parts[0]
	filename := parts[1]

	// Construir path: imageDir/objetoID/filename
	fullPath := filepath.Join(imageDir, objetoID, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "imagen no encontrada", 404)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	http.ServeFile(w, r, fullPath)
}
