class Imagen {
  final String id;
  final String path;
  final String thumbPath;
  final bool esPrincipal;
  final double? areaX;
  final double? areaY;
  final double? areaW;
  final double? areaH;

  Imagen({
    required this.id,
    required this.path,
    required this.thumbPath,
    this.esPrincipal = false,
    this.areaX,
    this.areaY,
    this.areaW,
    this.areaH,
  });

  factory Imagen.fromJson(Map<String, dynamic> json) => Imagen(
        id: json['id'] as String,
        path: json['path'] as String,
        thumbPath: json['thumb_path'] as String,
        esPrincipal: json['es_principal'] as bool? ?? false,
        areaX: (json['area_x'] as num?)?.toDouble(),
        areaY: (json['area_y'] as num?)?.toDouble(),
        areaW: (json['area_w'] as num?)?.toDouble(),
        areaH: (json['area_h'] as num?)?.toDouble(),
      );
}
