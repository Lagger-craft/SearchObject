package normalize

import "testing"

func TestNombre(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Garage", "garage"},
		{"  GARAGE  ", "garage"},
		{"Garage      ", "garage"},
		{"Mártillo", "martillo"},
		{"MARTILLO", "martillo"},
		{"Martillo", "martillo"},
		{"MáRtIlO", "martilo"},
		{"año", "año"},
		{"AÑO", "año"},
		{"Mañana", "mañana"},
		{"pingüino", "pinguino"},
		{"Pingüino", "pinguino"},
		{"estación", "estacion"},
		{"teléfono", "telefono"},
		{"álbum", "album"},
		{"", ""},
		{"   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Nombre(tt.input)
			if got != tt.expected {
				t.Errorf("Nombre(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestTags(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Herramientas", "herramientas"},
		{"Jardín", "jardin"},
		{"Máquina de Cortar", "maquina_de_cortar"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Tags(tt.input)
			if got != tt.expected {
				t.Errorf("Tags(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIguales(t *testing.T) {
	if !Iguales("Martillo", "Mártillo") {
		t.Error("Martillo != Mártillo")
	}
	if !Iguales("GARAGE", "garage") {
		t.Error("GARAGE != garage")
	}
	if Iguales("Martillo", "Martero") {
		t.Error("Martillo == Martero (deberían ser distintas)")
	}
}
