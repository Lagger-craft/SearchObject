import 'dart:convert';
import 'package:flutter/services.dart';
import 'package:http/http.dart' as http;

/// Bridge unificado: MethodChannel (gomobile/prod) o HTTP (dev).
/// Detecta automáticamente el modo y envía los params en el formato correcto.
class NativeBridge {
  static const _channel = MethodChannel('searchobject/bridge');
  static final NativeBridge instance = NativeBridge._();
  NativeBridge._();

  /// true = intenta MethodChannel primero (gomobile en el celular)
  /// false = siempre usa HTTP (dev en la PC)
  bool useMethodChannel = true; // Default: prod mode

  String _httpBaseUrl = 'http://10.0.2.2:8080';

  void setHttpBaseUrl(String url) => _httpBaseUrl = url;

  // ─── Core call ───

  Future<Map<String, dynamic>> call(
    String method,
    Map<String, dynamic> args,
  ) async {
    if (useMethodChannel) {
      try {
        return await _methodChannelCall(method, args);
      } catch (_) {
        // MethodChannel no disponible (dev mode en PC), fallback HTTP
        return _httpCall(method, args);
      }
    }
    return _httpCall(method, args);
  }

  /// HTTP: usa snake_case keys (como el server Go espera)
  Future<Map<String, dynamic>> _httpCall(
    String method,
    Map<String, dynamic> args,
  ) async {
    final uri = Uri.parse('$_httpBaseUrl/api/$method');
    final response = await http.post(
      uri,
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode(args),
    );
    if (response.statusCode != 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      throw Exception(body['error'] ?? 'Error ${response.statusCode}');
    }
    return jsonDecode(response.body) as Map<String, dynamic>;
  }

  /// MethodChannel: usa camelCase keys (como Kotlin espera)
  Future<Map<String, dynamic>> _methodChannelCall(
    String method,
    Map<String, dynamic> args,
  ) async {
    final result = await _channel.invokeMethod(method, args);
    if (result == null) return {};
    if (result is String) {
      return jsonDecode(result) as Map<String, dynamic>;
    }
    return Map<String, dynamic>.from(result as Map);
  }

  // ─── Auth ───

  Future<Map<String, dynamic>> crearUsuario(String nombre, String email) =>
      call('crearUsuario', {'nombre': nombre, 'email': email});

  Future<Map<String, dynamic>> obtenerUsuario(String id) =>
      call('obtenerUsuario', {'id': id});

  Future<Map<String, dynamic>> listarUsuarios() =>
      call('listarUsuarios', {});

  Future<Map<String, dynamic>> actualizarUsuario(String id, String nombre) =>
      call('actualizarUsuario', {'id': id, 'nombre': nombre});

  // ─── Espacios ───

  Future<Map<String, dynamic>> crearEspacio(String json) =>
      call('crearEspacio', {'json': json});

  Future<Map<String, dynamic>> listarEspacios(String usuarioId) =>
      call('listarEspacios', {'usuarioId': usuarioId});

  Future<Map<String, dynamic>> obtenerEspacio(String id) =>
      call('obtenerEspacio', {'id': id});

  Future<Map<String, dynamic>> actualizarEspacio(String json) =>
      call('actualizarEspacio', {'json': json});

  Future<Map<String, dynamic>> eliminarEspacio(String id) =>
      call('eliminarEspacio', {'id': id});

  // ─── Cajas ───

  Future<Map<String, dynamic>> crearCaja(String json) =>
      call('crearCaja', {'json': json});

  Future<Map<String, dynamic>> listarCajas(String espacioId) =>
      call('listarCajas', {'espacioId': espacioId});

  Future<Map<String, dynamic>> obtenerCaja(String id) =>
      call('obtenerCaja', {'id': id});

  Future<Map<String, dynamic>> eliminarCaja(String id) =>
      call('eliminarCaja', {'id': id});

  // ─── Objetos ───

  Future<Map<String, dynamic>> crearObjeto(String json) =>
      call('crearObjeto', {'json': json});

  Future<Map<String, dynamic>> listarObjetos(String cajaId) =>
      call('listarObjetos', {'cajaId': cajaId});

  Future<Map<String, dynamic>> obtenerObjeto(String id) =>
      call('obtenerObjeto', {'id': id});

  Future<Map<String, dynamic>> moverObjeto(String json) =>
      call('moverObjeto', {'json': json});

  Future<Map<String, dynamic>> eliminarObjeto(String id) =>
      call('eliminarObjeto', {'id': id});

  // ─── Búsqueda ───

  Future<Map<String, dynamic>> buscarObjetos(
          String usuarioId, String termino) =>
      call('buscarObjetos', {'usuarioId': usuarioId, 'termino': termino});

  Future<Map<String, dynamic>> buscar(String usuarioId, String termino) =>
      call('buscar', {'usuarioId': usuarioId, 'termino': termino});

  // ─── Imágenes ───

  /// Sube imagen: MethodChannel recibe ByteArray directo,
  /// HTTP recibe base64 en image_bytes.
  Future<Map<String, dynamic>> agregarImagen(
      String objetoId, List<int> imageBytes) async {
    if (useMethodChannel) {
      try {
        return await _methodChannelCall('agregarImagen', {
          'objetoId': objetoId,
          'imageBytes': imageBytes,
        });
      } catch (_) {
        // Fallback HTTP
      }
    }
    return _httpCall('agregarImagen', {
      'objeto_id': objetoId,
      'image_bytes': base64Encode(imageBytes),
    });
  }

  Future<Map<String, dynamic>> agregarImagenConArea(
      String objetoId, List<int> imageBytes, String jsonArea) async {
    if (useMethodChannel) {
      try {
        return await _methodChannelCall('agregarImagenConArea', {
          'objetoId': objetoId,
          'imageBytes': imageBytes,
          'jsonArea': jsonArea,
        });
      } catch (_) {
        // Fallback HTTP
      }
    }
    return _httpCall('agregarImagen', {
      'objeto_id': objetoId,
      'image_bytes': base64Encode(imageBytes),
      'json_area': jsonArea,
    });
  }

  Future<Map<String, dynamic>> listarImagenes(String objetoId) =>
      call('listarImagenes', {'objetoId': objetoId});

  Future<Map<String, dynamic>> eliminarImagen(String id) =>
      call('eliminarImagen', {'id': id});

  /// Obtiene la URL/Path para cargar una imagen.
  /// En MethodChannel: devuelve el path local del archivo.
  /// En HTTP: devuelve la URL del server.
  Future<String> imageUrl(String objetoId, String imagenId) async {
    if (useMethodChannel) {
      try {
        final dir = await _methodChannelCall('pathParaObjeto', {
          'objetoId': objetoId,
        });
        return '${dir['path']}/$imagenId.jpg';
      } catch (_) {}
    }
    return '$_httpBaseUrl/images/$objetoId/$imagenId.jpg';
  }

  // ─── Dashboard ───

  Future<Map<String, dynamic>> resumen(String usuarioId) =>
      call('resumen', {'usuarioId': usuarioId});

  Future<Map<String, dynamic>> dashboard(String usuarioId) =>
      call('dashboard', {'usuarioId': usuarioId});

  // ─── Alertas ───

  Future<Map<String, dynamic>> evaluarAlertas(String usuarioId) =>
      call('evaluarAlertas', {'usuarioId': usuarioId});

  Future<Map<String, dynamic>> listarAlertas(String jsonLeidas) =>
      call('listarAlertas', {'jsonLeidas': jsonLeidas});

  Future<Map<String, dynamic>> marcarAlertaLeida(String id) =>
      call('marcarAlertaLeida', {'id': id});

  Future<Map<String, dynamic>> resolverAlerta(String id) =>
      call('resolverAlerta', {'id': id});

  // ─── Reportes ───

  Future<Map<String, dynamic>> exportarJSON(String usuarioId) =>
      call('exportarJSON', {'usuarioId': usuarioId});

  Future<Map<String, dynamic>> exportarCSV(String usuarioId) =>
      call('exportarCSV', {'usuarioId': usuarioId});

  // ─── Tags ───

  Future<Map<String, dynamic>> crearTag(String json) =>
      call('crearTag', {'json': json});

  Future<Map<String, dynamic>> listarTags() =>
      call('listarTags', {});

  Future<Map<String, dynamic>> agregarTagAObjeto(String json) =>
      call('agregarTagAObjeto', {'json': json});

  Future<Map<String, dynamic>> quitarTagDeObjeto(String json) =>
      call('quitarTagDeObjeto', {'json': json});

  Future<Map<String, dynamic>> tagsDeObjeto(String objetoId) =>
      call('tagsDeObjeto', {'objetoId': objetoId});

  Future<Map<String, dynamic>> buscarPorTag(String usuarioId, String tagId) =>
      call('buscarPorTag', {'usuarioId': usuarioId, 'tagId': tagId});

  // ─── Historial ───

  Future<Map<String, dynamic>> registrarHistorial(String json) =>
      call('registrarHistorial', {'json': json});

  Future<Map<String, dynamic>> listarHistorial(String usuarioId) =>
      call('listarHistorial', {'usuarioId': usuarioId});

  Future<Map<String, dynamic>> historialPorEntidad(
          String usuarioId, String entidadTipo, String entidadId) =>
      call('historialPorEntidad', {
        'usuarioId': usuarioId,
        'entidadTipo': entidadTipo,
        'entidadId': entidadId,
      });

  // ─── Inventario ───

  Future<Map<String, dynamic>> stockBajo(String usuarioId, int limite) =>
      call('stockBajo', {'usuarioId': usuarioId, 'limite': limite});
}
