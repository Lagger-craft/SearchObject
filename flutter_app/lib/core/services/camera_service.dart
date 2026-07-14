import 'dart:io';
import 'package:image_picker/image_picker.dart';

/// Servicio para capturar imágenes desde la cámara nativa o galería.
/// Provee los bytes listos para enviar al backend Go via bridge.
class CameraService {
  CameraService._();
  static final instance = CameraService._();

  final _picker = ImagePicker();

  /// Toma una foto con la cámara nativa del teléfono.
  /// Devuelve los bytes de la imagen o null si el usuario canceló.
  Future<List<int>?> takePhoto() async {
    final file = await _picker.pickImage(
      source: ImageSource.camera,
      maxWidth: 1920,
      maxHeight: 1920,
      imageQuality: 85,
    );
    if (file == null) return null;
    return File(file.path).readAsBytes();
  }

  /// Abre la galería para seleccionar una imagen.
  /// Devuelve los bytes o null si el usuario canceló.
  Future<List<int>?> pickFromGallery() async {
    final file = await _picker.pickImage(
      source: ImageSource.gallery,
      maxWidth: 1920,
      maxHeight: 1920,
      imageQuality: 85,
    );
    if (file == null) return null;
    return File(file.path).readAsBytes();
  }

  /// Muestra un bottom sheet preguntando cámara o galería.
  /// Devuelve los bytes o null.
  Future<List<int>?> captureOrPick({
    required Future<List<int>?> Function(String source) onCapture,
  }) async {
    // Se deja la lógica de UI en la pantalla que llama,
    // acá solo tenemos las dos opciones puras.
    return null;
  }
}
