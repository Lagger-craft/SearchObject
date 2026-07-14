import 'dart:convert';
import 'package:flutter/material.dart';
import '../../core/bridge/native_bridge.dart';
import '../../core/models/espacio.dart';
import '../../core/models/caja.dart';
import '../../core/models/objeto.dart';
import 'detalle_objeto_screen.dart';

class ObjetosScreen extends StatefulWidget {
  final String cajaId;
  final String cajaNombre;
  final String usuarioId;

  const ObjetosScreen({
    super.key,
    required this.cajaId,
    required this.cajaNombre,
    required this.usuarioId,
  });

  @override
  State<ObjetosScreen> createState() => _ObjetosScreenState();
}

class _ObjetosScreenState extends State<ObjetosScreen> {
  final _bridge = NativeBridge.instance;
  List<Objeto> _objetos = [];
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _cargar();
  }

  Future<void> _cargar() async {
    setState(() => _loading = true);
    try {
      final json = await _bridge.listarObjetos(widget.cajaId);
      final lista = json['objetos'] as List<dynamic>;
      if (!mounted) return;
      setState(() {
        _objetos = lista.map((e) => Objeto.fromJson(e)).toList();
        _loading = false;
      });
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  Future<void> _crearObjeto() async {
    final result = await _mostrarDialogoCrear();
    if (result == null) return;

    try {
      final jsonStr = jsonEncode({
        'caja_id': widget.cajaId,
        'usuario_id': widget.usuarioId,
        'nombre': result['nombre'],
        'descripcion': result['descripcion'],
        'cantidad': result['cantidad'],
      });
      await _bridge.crearObjeto(jsonStr);
      await _cargar();
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(e.toString()), backgroundColor: Colors.red),
      );
    }
  }

  void _abrirDetalle(Objeto obj) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (_) => DetalleObjetoScreen(objeto: obj),
      ),
    ).then((_) => _cargar()); // Recargar al volver por si se agregaron fotos
  }

  Future<void> _moverObjeto(Objeto obj) async {
    try {
      final cajaDestino = await _seleccionarCajaDestino();
      if (cajaDestino == null || cajaDestino == widget.cajaId) return;

      final jsonStr = jsonEncode({
        'id': obj.id,
        'caja_id': cajaDestino,
        'nota': '',
      });
      await _bridge.moverObjeto(jsonStr);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Objeto movido'),
          backgroundColor: Colors.green,
        ),
      );
      await _cargar();
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(e.toString()), backgroundColor: Colors.red),
      );
    }
  }

  Future<String?> _seleccionarCajaDestino() async {
    final espacio = await _mostrarSelectorEspacios();
    if (espacio == null) return null;
    return _mostrarSelectorCajas(espacio.id, espacio.nombre);
  }

  Future<Espacio?> _mostrarSelectorEspacios() async {
    final json = await _bridge.listarEspacios(widget.usuarioId);
    final lista = json['espacios'] as List<dynamic>;
    final espacios = lista.map((e) => Espacio.fromJson(e)).toList();

    if (!mounted || espacios.isEmpty) return null;

    return showDialog<Espacio>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Mover a espacio...'),
        content: SizedBox(
          width: double.maxFinite,
          child: ListView.builder(
            shrinkWrap: true,
            itemCount: espacios.length,
            itemBuilder: (_, i) {
              final esp = espacios[i];
              return ListTile(
                leading: const Icon(Icons.folder),
                title: Text(esp.nombre),
                trailing: const Icon(Icons.chevron_right),
                onTap: () => Navigator.pop(ctx, esp),
              );
            },
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cancelar'),
          ),
        ],
      ),
    );
  }

  Future<String?> _mostrarSelectorCajas(
      String espacioId, String espacioNombre) async {
    final json = await _bridge.listarCajas(espacioId);
    final lista = json['cajas'] as List<dynamic>;
    final cajas = lista.map((e) => Caja.fromJson(e)).toList();

    if (!mounted || cajas.isEmpty) return null;

    return showDialog<String>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('Mover a caja en "$espacioNombre"...'),
        content: SizedBox(
          width: double.maxFinite,
          child: ListView.builder(
            shrinkWrap: true,
            itemCount: cajas.length,
            itemBuilder: (_, i) {
              final caja = cajas[i];
              return ListTile(
                leading: const Icon(Icons.inventory_2_outlined),
                title: Text(caja.nombre),
                subtitle: caja.descripcion.isNotEmpty
                    ? Text(caja.descripcion)
                    : null,
                onTap: () => Navigator.pop(ctx, caja.id),
              );
            },
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cancelar'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.cajaNombre),
        actions: [
          IconButton(icon: const Icon(Icons.refresh), onPressed: _cargar),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _crearObjeto,
        child: const Icon(Icons.add),
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (_loading) return const Center(child: CircularProgressIndicator());
    if (_error != null) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, size: 48, color: Colors.red),
            const SizedBox(height: 16),
            Text(_error!, textAlign: TextAlign.center),
            const SizedBox(height: 16),
            ElevatedButton(onPressed: _cargar, child: const Text('Reintentar')),
          ],
        ),
      );
    }

    if (_objetos.isEmpty) {
      return const Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.inventory_outlined, size: 64, color: Colors.grey),
            SizedBox(height: 16),
            Text('No hay objetos todavía', style: TextStyle(color: Colors.grey)),
            SizedBox(height: 8),
            Text('Tocá + para crear uno', style: TextStyle(color: Colors.grey)),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: _cargar,
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: _objetos.length,
        itemBuilder: (context, index) {
          final obj = _objetos[index];
          return Card(
            child: ListTile(
              leading: Icon(
                obj.esInsumo ? Icons.construction : Icons.inventory_2,
              ),
              title: Text(obj.nombre),
              subtitle: Row(
                children: [
                  if (obj.descripcion.isNotEmpty)
                    Expanded(
                      child: Text(
                        obj.descripcion,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                  const SizedBox(width: 8),
                  Text('x${obj.cantidad}'),
                ],
              ),
              trailing: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  if (obj.tieneFotos)
                    const Icon(Icons.photo, size: 20, color: Colors.blue),
                  IconButton(
                    icon: const Icon(Icons.camera_alt, size: 20),
                    tooltip: 'Ver fotos',
                    onPressed: () => _abrirDetalle(obj),
                  ),
                  IconButton(
                    icon: const Icon(Icons.open_with, size: 20),
                    tooltip: 'Mover',
                    onPressed: () => _moverObjeto(obj),
                  ),
                ],
              ),
              onTap: () => _abrirDetalle(obj),
            ),
          );
        },
      ),
    );
  }

  Future<Map<String, dynamic>?> _mostrarDialogoCrear() {
    final nombreCtrl = TextEditingController();
    final descCtrl = TextEditingController();
    final cantidadCtrl = TextEditingController(text: '1');

    return showDialog<Map<String, dynamic>>(
      context: context,
      builder: (ctx) {
        return AlertDialog(
          title: const Text('Nuevo objeto'),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                TextField(
                  controller: nombreCtrl,
                  autofocus: true,
                  decoration: const InputDecoration(hintText: 'Nombre'),
                ),
                const SizedBox(height: 8),
                TextField(
                  controller: descCtrl,
                  decoration: const InputDecoration(hintText: 'Descripción'),
                ),
                const SizedBox(height: 8),
                TextField(
                  controller: cantidadCtrl,
                  keyboardType: TextInputType.number,
                  decoration: const InputDecoration(hintText: 'Cantidad'),
                ),
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(ctx),
              child: const Text('Cancelar'),
            ),
            FilledButton(
              onPressed: () {
                final nombre = nombreCtrl.text.trim();
                if (nombre.isEmpty) return;
                Navigator.pop(ctx, {
                  'nombre': nombre,
                  'descripcion': descCtrl.text.trim(),
                  'cantidad': int.tryParse(cantidadCtrl.text) ?? 1,
                });
              },
              child: const Text('Crear'),
            ),
          ],
        );
      },
    );
  }
}
