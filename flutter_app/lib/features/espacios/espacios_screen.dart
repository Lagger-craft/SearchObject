import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../core/bridge/native_bridge.dart';
import '../../core/models/espacio.dart';
import '../cajas/cajas_screen.dart';
import '../search/search_screen.dart';

class EspaciosScreen extends StatefulWidget {
  const EspaciosScreen({super.key, required this.userId, required this.userName});

  final String userId;
  final String userName;

  @override
  State<EspaciosScreen> createState() => _EspaciosScreenState();
}

class _EspaciosScreenState extends State<EspaciosScreen> {
  final _bridge = NativeBridge.instance;
  late String _userId;
  late String _userName;
  List<Espacio> _espacios = [];
  bool _loading = true;
  String? _error;
  int _currentTab = 0;

  @override
  void initState() {
    super.initState();
    _userId = widget.userId;
    _userName = widget.userName;
    _cargarEspacios();
  }

  Future<void> _editarNombre() async {
    final ctrl = TextEditingController(text: _userName);
    final nuevo = await showDialog<String>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Tu nombre'),
        content: TextField(
          controller: ctrl,
          autofocus: true,
          textCapitalization: TextCapitalization.words,
          decoration: const InputDecoration(
            hintText: 'Cómo querés que te digamos',
            prefixIcon: Icon(Icons.person),
          ),
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Cancelar')),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, ctrl.text.trim()),
            child: const Text('Guardar'),
          ),
        ],
      ),
    );

    ctrl.dispose();
    if (nuevo == null || nuevo.isEmpty || nuevo == _userName) return;

    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('userName', nuevo);
    if (!mounted) return;
    setState(() => _userName = nuevo);
  }

  Future<void> _cargarEspacios() async {
    final json = await _bridge.listarEspacios(_userId!);
    final lista = json['espacios'] as List<dynamic>;
    if (!mounted) return;
    setState(() {
      _espacios = lista.map((e) => Espacio.fromJson(e)).toList();
      _loading = false;
    });
  }

  Future<void> _crearEspacio() async {
    final nombre = await _mostrarDialogo();
    if (nombre == null || nombre.isEmpty) return;

    try {
      final jsonStr = jsonEncode({
        'usuario_id': _userId!,
        'nombre': nombre,
      });
      await _bridge.crearEspacio(jsonStr);
      await _cargarEspacios();
    } catch (e) {
      if (!mounted) return;
      _mostrarError(e.toString());
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(_currentTab == 0 ? 'SearchObject' : 'Buscar'),
        actions: [
          if (_currentTab == 0)
            IconButton(
              icon: const Icon(Icons.refresh),
              onPressed: _cargarEspacios,
            ),
          IconButton(
            icon: const Icon(Icons.account_circle),
            tooltip: _userName,
            onPressed: _editarNombre,
          ),
        ],
      ),
      floatingActionButton: _currentTab == 0
          ? FloatingActionButton(
              onPressed: _crearEspacio,
              child: const Icon(Icons.add),
            )
          : null,
      body: _currentTab == 0 ? _buildEspacios() : _buildSearch(),
      bottomNavigationBar: NavigationBar(
        selectedIndex: _currentTab,
        onDestinationSelected: (i) => setState(() => _currentTab = i),
        destinations: const [
          NavigationDestination(
            icon: Icon(Icons.home_outlined),
            selectedIcon: Icon(Icons.home),
            label: 'Espacios',
          ),
          NavigationDestination(
            icon: Icon(Icons.search_outlined),
            selectedIcon: Icon(Icons.search),
            label: 'Buscar',
          ),
        ],
      ),
    );
  }

  Widget _buildEspacios() {
    if (_loading) {
      return const Center(child: CircularProgressIndicator());
    }
    if (_error != null) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, size: 48, color: Colors.red),
            const SizedBox(height: 16),
            Text(_error!, textAlign: TextAlign.center),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: _cargarEspacios,
              child: const Text('Reintentar'),
            ),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: _cargarEspacios,
      child: _espacios.isEmpty
          ? const Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.folder_outlined, size: 64, color: Colors.grey),
                  SizedBox(height: 16),
                  Text('No hay espacios todavía',
                      style: TextStyle(color: Colors.grey)),
                  SizedBox(height: 8),
                  Text('Tocá + para crear uno',
                      style: TextStyle(color: Colors.grey)),
                ],
              ),
            )
          : ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: _espacios.length,
              itemBuilder: (context, index) {
                final esp = _espacios[index];
                return Card(
                  child: ListTile(
                    leading: const Icon(Icons.folder),
                    title: Text(esp.nombre),
                    subtitle: esp.descripcion.isNotEmpty
                        ? Text(esp.descripcion)
                        : null,
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (_) => CajasScreen(
                            espacioId: esp.id,
                            espacioNombre: esp.nombre,
                            usuarioId: _userId!,
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

  Widget _buildSearch() {
    if (_userId == null) {
      return const Center(child: CircularProgressIndicator());
    }
    return SearchScreen(usuarioId: _userId!);
  }

  Future<String?> _mostrarDialogo() => showDialog<String>(
        context: context,
        builder: (ctx) {
          final ctrl = TextEditingController();
          return AlertDialog(
            title: const Text('Nuevo espacio'),
            content: TextField(
              controller: ctrl,
              autofocus: true,
              decoration: const InputDecoration(
                hintText: 'Nombre del espacio',
              ),
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

  void _mostrarError(String msg) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(msg), backgroundColor: Colors.red),
    );
  }
}
