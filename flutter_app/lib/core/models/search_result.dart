class SearchResult {
  final String objetoId;
  final String objetoNombre;
  final String objetoDescripcion;
  final int objetoCantidad;
  final bool objetoEsInsumo;
  final bool tieneFotos;
  final String espacioId;
  final String espacioNombre;
  final String cajaId;
  final String cajaNombre;

  SearchResult({
    required this.objetoId,
    required this.objetoNombre,
    this.objetoDescripcion = '',
    this.objetoCantidad = 1,
    this.objetoEsInsumo = false,
    this.tieneFotos = false,
    required this.espacioId,
    required this.espacioNombre,
    required this.cajaId,
    required this.cajaNombre,
  });

  factory SearchResult.fromJson(Map<String, dynamic> json) {
    final obj = json['objeto'] as Map<String, dynamic>;
    final esp = json['espacio'] as Map<String, dynamic>;
    final caja = json['caja'] as Map<String, dynamic>;
    return SearchResult(
      objetoId: obj['id'] as String,
      objetoNombre: obj['nombre'] as String,
      objetoDescripcion: obj['descripcion'] as String? ?? '',
      objetoCantidad: obj['cantidad'] as int? ?? 1,
      objetoEsInsumo: obj['es_insumo'] as bool? ?? false,
      tieneFotos: obj['tiene_fotos'] as bool? ?? false,
      espacioId: esp['id'] as String,
      espacioNombre: esp['nombre'] as String,
      cajaId: caja['id'] as String,
      cajaNombre: caja['nombre'] as String,
    );
  }
}
