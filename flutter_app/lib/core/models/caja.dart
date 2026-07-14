class Caja {
  final String id;
  final String espacioId;
  final String nombre;
  final String descripcion;
  final int? capacidadMax;

  Caja({
    required this.id,
    required this.espacioId,
    required this.nombre,
    this.descripcion = '',
    this.capacidadMax,
  });

  factory Caja.fromJson(Map<String, dynamic> json) => Caja(
        id: json['id'] as String,
        espacioId: json['espacio_id'] as String,
        nombre: json['nombre'] as String,
        descripcion: json['descripcion'] as String? ?? '',
        capacidadMax: json['capacidad_max'] as int?,
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'nombre': nombre,
      };
}
