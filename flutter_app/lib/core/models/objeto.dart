class Objeto {
  final String id;
  final String cajaId;
  final String nombre;
  final String descripcion;
  final int cantidad;
  final bool esInsumo;
  final double? valor;
  final bool tieneFotos;

  Objeto({
    required this.id,
    required this.cajaId,
    required this.nombre,
    this.descripcion = '',
    this.cantidad = 1,
    this.esInsumo = false,
    this.valor,
    this.tieneFotos = false,
  });

  factory Objeto.fromJson(Map<String, dynamic> json) => Objeto(
        id: json['id'] as String,
        cajaId: json['caja_id'] as String,
        nombre: json['nombre'] as String,
        descripcion: json['descripcion'] as String? ?? '',
        cantidad: json['cantidad'] as int? ?? 1,
        esInsumo: json['es_insumo'] as bool? ?? false,
        valor: json['valor'] as double?,
        tieneFotos: json['tiene_fotos'] as bool? ?? false,
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'nombre': nombre,
      };
}
