import 'dart:io';
import 'package:flutter/material.dart';
import '../../core/bridge/native_bridge.dart';
import '../../core/models/imagen.dart';
import '../../core/models/objeto.dart';
import '../../core/services/camera_service.dart';

class DetalleObjetoScreen extends StatefulWidget {
  final Objeto objeto;

  const DetalleObjetoScreen({super.key, required this.objeto});

  @override
  State<DetalleObjetoScreen> createState() => _DetalleObjetoScreenState();
}

class _DetalleObjetoScreenState extends State<DetalleObjetoScreen> {
  final _bridge = NativeBridge.instance;
  final _camera = CameraService.instance;
  List<Imagen> _imagenes = [];
  bool _loading = true;
  bool _subiendo = false;

  @override
  void initState() {
    super.initState();
    _cargarImagenes();
  }

  Future<void> _cargarImagenes() async {
    setState(() => _loading = true);
    try {
      final json = await _bridge.listarImagenes(widget.objeto.id);
      final lista = json['imagenes'] as List<dynamic>;
      if (!mounted) return;
      setState(() {
        _imagenes = lista.map((e) => Imagen.fromJson(e)).toList();
        _loading = false;
      });
    } catch (e) {
      if (!mounted) return;
      setState(() => _loading = false);
    }
  }

  Future<void> _tomarFoto() async {
    final bytes = await _camera.takePhoto();
    if (bytes == null) return;
    await _subirImagen(bytes);
  }

  Future<void> _seleccionarGaleria() async {
    final bytes = await _camera.pickFromGallery();
    if (bytes == null) return;
    await _subirImagen(bytes);
  }

  Future<void> _subirImagen(List<int> bytes) async {
    setState(() => _subiendo = true);
    try {
      await _bridge.agregarImagen(widget.objeto.id, bytes);
      await _cargarImagenes();
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Foto agregada'),
          backgroundColor: Colors.green,
        ),
      );
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
      );
    } finally {
      if (mounted) setState(() => _subiendo = false);
    }
  }

  Future<void> _eliminarImagen(Imagen img) async {
    final confirm = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Eliminar foto'),
        content: const Text('¿Seguro que querés eliminar esta foto?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Cancelar'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            style: FilledButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('Eliminar'),
          ),
        ],
      ),
    );

    if (confirm != true) return;

    try {
      await _bridge.eliminarImagen(img.id);
      await _cargarImagenes();
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
      );
    }
  }

  void _mostrarFotoCompleta(Imagen img) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (_) => _FullImageView(imagen: img, objetoId: widget.objeto.id),
      ),
    );
  }

  void _mostrarOpcionesFoto() {
    showModalBottomSheet(
      context: context,
      builder: (ctx) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 8),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              ListTile(
                leading: const Icon(Icons.camera_alt),
                title: const Text('Tomar foto'),
                onTap: () {
                  Navigator.pop(ctx);
                  _tomarFoto();
                },
              ),
              ListTile(
                leading: const Icon(Icons.photo_library),
                title: const Text('Elegir de galería'),
                onTap: () {
                  Navigator.pop(ctx);
                  _seleccionarGaleria();
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.objeto.nombre),
        actions: [
          IconButton(icon: const Icon(Icons.refresh), onPressed: _cargarImagenes),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _subiendo ? null : _mostrarOpcionesFoto,
        child: _subiendo
            ? const SizedBox(
                width: 24,
                height: 24,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  color: Colors.white,
                ),
              )
            : const Icon(Icons.camera_alt),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Icon(
                          widget.objeto.esInsumo
                              ? Icons.construction
                              : Icons.inventory_2,
                          size: 32,
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                widget.objeto.nombre,
                                style: Theme.of(context).textTheme.titleLarge,
                              ),
                              if (widget.objeto.descripcion.isNotEmpty)
                                Text(
                                  widget.objeto.descripcion,
                                  style: Theme.of(context).textTheme.bodyMedium,
                                ),
                            ],
                          ),
                        ),
                      ],
                    ),
                    const Divider(height: 24),
                    Row(
                      children: [
                        _InfoChip(
                          icon: Icons.numbers,
                          label: 'x${widget.objeto.cantidad}',
                        ),
                        const SizedBox(width: 8),
                        if (widget.objeto.tieneFotos)
                          _InfoChip(
                            icon: Icons.photo,
                            label: '${_imagenes.length} fotos',
                          ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),
            Text(
              'Fotos',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            _buildGaleria(),
          ],
        ),
      ),
    );
  }

  Widget _buildGaleria() {
    if (_loading) {
      return const Center(
        child: Padding(
          padding: EdgeInsets.all(32),
          child: CircularProgressIndicator(),
        ),
      );
    }

    if (_imagenes.isEmpty) {
      return Card(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Center(
            child: Column(
              children: [
                Icon(Icons.camera_alt, size: 48, color: Colors.grey[400]),
                const SizedBox(height: 12),
                Text(
                  'Sin fotos todavía',
                  style: TextStyle(color: Colors.grey[600]),
                ),
                const SizedBox(height: 4),
                Text(
                  'Tocá la cámara para agregar una',
                  style: TextStyle(color: Colors.grey[500], fontSize: 12),
                ),
              ],
            ),
          ),
        ),
      );
    }

    return GridView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 3,
        crossAxisSpacing: 8,
        mainAxisSpacing: 8,
      ),
      itemCount: _imagenes.length,
      itemBuilder: (context, index) {
        final img = _imagenes[index];
        return _Thumbnail(
          imagen: img,
          objetoId: widget.objeto.id,
          onTap: () => _mostrarFotoCompleta(img),
          onLongPress: () => _eliminarImagen(img),
        );
      },
    );
  }
}

// ─── Widgets internos ───

class _InfoChip extends StatelessWidget {
  final IconData icon;
  final String label;

  const _InfoChip({required this.icon, required this.label});

  @override
  Widget build(BuildContext context) {
    return Chip(
      avatar: Icon(icon, size: 16),
      label: Text(label, style: const TextStyle(fontSize: 12)),
      padding: const EdgeInsets.symmetric(horizontal: 4),
      materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
      visualDensity: VisualDensity.compact,
    );
  }
}

/// Thumbnail que carga imagen desde HTTP o filesystem local.
class _Thumbnail extends StatefulWidget {
  final Imagen imagen;
  final String objetoId;
  final VoidCallback onTap;
  final VoidCallback onLongPress;

  const _Thumbnail({
    required this.imagen,
    required this.objetoId,
    required this.onTap,
    required this.onLongPress,
  });

  @override
  State<_Thumbnail> createState() => _ThumbnailState();
}

class _ThumbnailState extends State<_Thumbnail> {
  String? _imageUrl;
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _loadImage();
  }

  Future<void> _loadImage() async {
    final bridge = NativeBridge.instance;
    // Extract imagenId from thumbPath (e.g., "/tmp/.../images/objetoId/imagenId_thumb.jpg")
    final thumbPath = widget.imagen.thumbPath;
    String? imagenId;
    if (thumbPath.isNotEmpty) {
      final parts = thumbPath.split('/');
      final filename = parts.last; // imagenId_thumb.jpg
      imagenId = filename.replaceAll('_thumb.jpg', '').replaceAll('.jpg', '');
    }

    if (imagenId != null) {
      final url = await bridge.imageUrl(widget.objetoId, imagenId);
      if (mounted) setState(() { _imageUrl = url; _loading = false; });
    } else {
      if (mounted) setState(() { _loading = false; });
    }
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: widget.onTap,
      onLongPress: widget.onLongPress,
      child: Stack(
        fit: StackFit.expand,
        children: [
          Container(
            decoration: BoxDecoration(
              color: Colors.grey[300],
              borderRadius: BorderRadius.circular(8),
            ),
            clipBehavior: Clip.antiAlias,
            child: _buildImage(),
          ),
          if (widget.imagen.esPrincipal)
            Positioned(
              top: 4,
              left: 4,
              child: Container(
                padding: const EdgeInsets.all(2),
                decoration: BoxDecoration(
                  color: Colors.blue,
                  borderRadius: BorderRadius.circular(4),
                ),
                child: const Icon(Icons.star, size: 12, color: Colors.white),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildImage() {
    if (_loading) {
      return const Center(child: SizedBox(
        width: 20, height: 20,
        child: CircularProgressIndicator(strokeWidth: 2),
      ));
    }

    final url = _imageUrl ?? '';
    if (url.isEmpty) {
      return const Center(child: Icon(Icons.photo, size: 32, color: Colors.grey));
    }

    // HTTP URL → Image.network
    if (url.startsWith('http')) {
      return Image.network(
        url,
        fit: BoxFit.cover,
        errorBuilder: (_, __, ___) => const Center(
          child: Icon(Icons.broken_image, size: 32, color: Colors.grey),
        ),
      );
    }

    // Filesystem path → Image.file
    final file = File(url);
    if (file.existsSync()) {
      return Image.file(file, fit: BoxFit.cover);
    }

    return const Center(child: Icon(Icons.photo, size: 32, color: Colors.grey));
  }
}

/// Vista completa: carga imagen desde HTTP o filesystem.
class _FullImageView extends StatefulWidget {
  final Imagen imagen;
  final String objetoId;

  const _FullImageView({required this.imagen, required this.objetoId});

  @override
  State<_FullImageView> createState() => _FullImageViewState();
}

class _FullImageViewState extends State<_FullImageView> {
  String? _imageUrl;
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _loadImage();
  }

  Future<void> _loadImage() async {
    final bridge = NativeBridge.instance;
    final path = widget.imagen.path;
    String? imagenId;
    if (path.isNotEmpty) {
      final parts = path.split('/');
      final filename = parts.last;
      imagenId = filename.replaceAll('.jpg', '');
    }

    if (imagenId != null) {
      final url = await bridge.imageUrl(widget.objetoId, imagenId);
      if (mounted) setState(() { _imageUrl = url; _loading = false; });
    } else {
      if (mounted) setState(() { _loading = false; });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.black,
        iconTheme: const IconThemeData(color: Colors.white),
      ),
      body: Center(child: _buildBody()),
    );
  }

  Widget _buildBody() {
    if (_loading) {
      return const CircularProgressIndicator(color: Colors.white);
    }

    final url = _imageUrl ?? '';
    if (url.isEmpty) {
      return _errorPlaceholder();
    }

    Widget image;
    if (url.startsWith('http')) {
      image = Image.network(
        url,
        fit: BoxFit.contain,
        errorBuilder: (_, __, ___) => _errorPlaceholder(),
      );
    } else {
      final file = File(url);
      if (file.existsSync()) {
        image = Image.file(file, fit: BoxFit.contain);
      } else {
        return _errorPlaceholder();
      }
    }

    return InteractiveViewer(
      minScale: 0.5,
      maxScale: 4.0,
      child: image,
    );
  }

  Widget _errorPlaceholder() {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        const Icon(Icons.broken_image, size: 64, color: Colors.grey),
        const SizedBox(height: 16),
        Text(
          widget.imagen.path.split('/').last,
          style: const TextStyle(color: Colors.white70),
        ),
      ],
    );
  }
}
