---

## 📱 Captura de Imágenes

El sistema se encarga de tomar imágenes del usuario. La app permitirá indicarle el área en cuestión pinchando encima de la imagen. Según donde se marque, se seleccionará un área delimitada en donde estarán guardados los objetos.

---

## 🏷️ Sistema de Tags

El usuario podrá agregarle un tag para realizar un orden a las ubicaciones, además de poder mostrarlas de manera más ordenada desde la app.

---

## 🔍 Detección de Objetos

De momento, el usuario tendrá que definirlo. La app no se ha indicado que tenga que adivinar o saber qué objeto en concreto se está indicando.

---

## ❌ Elementos Descartados

| Nivel | Elemento | Estado |
|-------|----------|--------|
| NIVEL 3 — CAJAS | Opción QR | Descartado |
| NIVEL 4 - OBJETOS | "Foto rápido" | Descartado (pierde sentido) |
| NIVEL 9 — REPORTES | IA | Descartado |

> **Nota:** La opción "manual y por voz" en NIVEL 4 se refiere a: anotación manual o nota de voz.

---

## ⏳ Pendiente (sin descartar)

- **NIVEL 5 — BÚSQUEDA VISUAL:** Agregar "zoom" del móvil — pendiente hasta entender cómo se llevaría a cabo.

---

## 🔎 Sistema de Búsqueda

La app tendrá una opción de búsqueda mencionada al inicio:

```markdown
5. Mejor UX

La búsqueda podría ser:

Buscar → Taladro

Resultado:

📍 Garage  
📸 Mostrar imagen  
🔴 Resaltado aquí
```

---

## 🤖 Inteligencia Artificial

La IA mencionada en [[TU NUEVA IDEA.docx]] **NO se implementará**. Se descarta nuevamente.

> *Nota:* Antiguamente con Arduino ya se pudo implementar un escaneo de objetos, por lo que se podrá llevar a cabo sin IA.

---

## 📊 Análisis del Organigrama

Cambiando a un archivo más completo: [[ORGANIGRAMA ESTRUCTURAL DEL SISTEMA.docx]]

---

### NIVEL 8 — ALERTAS

> Las alertas son resultado del análisis.

```
INVENTARIO  
│  
▼  
ANALÍTICA  
│  
├── faltantes          → Avisará todo (herramientas y/o insumos)
├── stock bajo         → Solo insumos (Tornillos, tuercas, etc.)
├── préstamos vencidos  
├── objetos perdidos  
└── cajas saturadas
```

---

### NIVEL 9 — REPORTES

> La IA mencionada se **descarta** nuevamente.

---

## 🗂️ Estructura General del Sistema

```
USUARIOS  
   │  
   ▼  
ESPACIOS  
   │  
   ▼  
CAJAS  
   │  
   ▼  
OBJETOS  
   │  
   ▼  
ESCANEOS  
   │  
   ▼  
INVENTARIO  
   │  
   ▼  
MOVIMIENTOS  
   │  
   ▼  
HISTORIAL  
   │  
   ▼  
BÚSQUEDA  
   │  
   ▼  
ALERTAS  
   │  
   ▼  
REPORTES
```

---

## 📈 Dashboard

```
DASHBOARD  
│  
├── total objetos  
├── valor estimado  
├── uso frecuente  
├── objetos perdidos  
├── movimientos  
└── estadísticas
```

---

## 🏠 Ejemplo de Estructura de Guardado

```
CASA  
└── ESPACIO  
    └── SUBESPACIO  
        └── CAJA  
            └── OBJETO  
                └── MOVIMIENTO  
                    └── HISTORIAL
```

---