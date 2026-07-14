import 'dart:convert';
import 'package:flutter/material.dart';
import '../../core/bridge/native_bridge.dart';
import '../../core/models/caja.dart';
import '../objetos/objetos_screen.dart';

class CajasScreen extends StatefulWidget {
  final String espacioId;
  final String espacioNombre;
  final String usuarioId;

  const CajasScreen({
    super.key,
    required this.espacioId,
    required this.espacioNombre,
    required this.usuarioId,
  });

  @override
  State<CajasScreen> createState() => _CajasScreenState();
}

class _CajasScreenState extends State<CajasScreen> {
  final _bridge = NativeBridge.instance;
  List<Caja> _cajas = [];
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
      final json = await _bridge.listarCajas(widget.espacioId);
      final lista = json['cajas'] as List<dynamic>;
      if (!mounted) return;
      setState(() {
        _cajas = lista.map((e) => Caja.fromJson(e)).toList();
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

  Future<void> _crearCaja() async {
    final nombre = await _mostrarDialogo('Nueva caja', 'Nombre de la caja');
    if (nombre == null || nombre.isEmpty) return;

    try {
      final jsonStr = jsonEncode({
        'espacio_id': widget.espacioId,
        'usuario_id': widget.usuarioId,
        'nombre': nombre,
      });
      await _bridge.crearCaja(jsonStr);
      await _cargar();
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(e.toString()), backgroundColor: Colors.red),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.espacioNombre),
        actions: [
          IconButton(icon: const Icon(Icons.refresh), onPressed: _cargar),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _crearCaja,
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

    if (_cajas.isEmpty) {
      return const Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.inbox_outlined, size: 64, color: Colors.grey),
            SizedBox(height: 16),
            Text('No hay cajas todavía', style: TextStyle(color: Colors.grey)),
            SizedBox(height: 8),
            Text('Tocá + para crear una', style: TextStyle(color: Colors.grey)),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: _cargar,
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: _cajas.length,
        itemBuilder: (context, index) {
          final caja = _cajas[index];
          return Card(
            child: ListTile(
              leading: const Icon(Icons.inventory_2_outlined),
              title: Text(caja.nombre),
              subtitle: caja.descripcion.isNotEmpty
                  ? Text(caja.descripcion)
                  : null,
              trailing: const Icon(Icons.chevron_right),
              onTap: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (_) => ObjetosScreen(
                      cajaId: caja.id,
                      cajaNombre: caja.nombre,
                      usuarioId: widget.usuarioId,
                    ),
                  ),
                );
              },
            ),
          );
        },
      ),
    );
  }

  Future<String?> _mostrarDialogo(String title, String hint) =>
      showDialog<String>(
        context: context,
        builder: (ctx) {
          final ctrl = TextEditingController();
          return AlertDialog(
            title: Text(title),
            content: TextField(
              controller: ctrl,
              autofocus: true,
              decoration: InputDecoration(hintText: hint),
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(ctx),
                child: const Text('Cancelar'),
              ),
              FilledButton(
                onPressed: () => Navigator.pop(ctx, ctrl.text),
                child: const Text('Crear'),
              ),
            ],
          );
        },
      );
}
