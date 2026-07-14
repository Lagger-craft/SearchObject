class Espacio {
  final String id;
  final String nombre;
  final String descripcion;
  final String? padreId;
  final bool tieneHijos;

  Espacio({
    required this.id,
    required this.nombre,
    this.descripcion = '',
    this.padreId,
    this.tieneHijos = false,
  });

  factory Espacio.fromJson(Map<String, dynamic> json) => Espacio(
        id: json['id'] as String,
        nombre: json['nombre'] as String,
        descripcion: json['descripcion'] as String? ?? '',
        padreId: json['padre_id'] as String?,
        tieneHijos: json['tiene_hijos'] as bool? ?? false,
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'nombre': nombre,
      };
}
