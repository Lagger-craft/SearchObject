import 'package:flutter/material.dart';
import '../../core/bridge/native_bridge.dart';
import '../../core/models/search_result.dart';
import '../objetos/objetos_screen.dart';

class SearchScreen extends StatefulWidget {
  final String usuarioId;

  const SearchScreen({super.key, required this.usuarioId});

  @override
  State<SearchScreen> createState() => _SearchScreenState();
}

class _SearchScreenState extends State<SearchScreen> {
  final _bridge = NativeBridge.instance;
  final _controller = TextEditingController();
  final _focusNode = FocusNode();
  List<SearchResult> _results = [];
  bool _loading = false;
  bool _searched = false;

  @override
  void dispose() {
    _controller.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  Future<void> _buscar(String termino) async {
    final t = termino.trim();
    if (t.isEmpty) return;

    setState(() {
      _loading = true;
      _searched = true;
    });

    try {
      final json = await _bridge.buscar(widget.usuarioId, t);
      final lista = json['resultados'] as List<dynamic>;
      if (!mounted) return;
      setState(() {
        _results = lista.map((e) => SearchResult.fromJson(e)).toList();
        _loading = false;
      });
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _results = [];
        _loading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: TextField(
          controller: _controller,
          focusNode: _focusNode,
          autofocus: true,
          decoration: const InputDecoration(
            hintText: 'Buscar objetos...',
            border: InputBorder.none,
          ),
          onSubmitted: _buscar,
          textInputAction: TextInputAction.search,
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: () => _buscar(_controller.text),
          ),
        ],
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (_loading) return const Center(child: CircularProgressIndicator());

    if (!_searched) {
      return const Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.search, size: 64, color: Colors.grey),
            SizedBox(height: 16),
            Text('Buscá por nombre de objeto',
                style: TextStyle(color: Colors.grey)),
          ],
        ),
      );
    }

    if (_results.isEmpty) {
      return const Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.search_off, size: 64, color: Colors.grey),
            SizedBox(height: 16),
            Text('No se encontraron resultados',
                style: TextStyle(color: Colors.grey)),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: _results.length,
      itemBuilder: (context, index) {
        final r = _results[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 8),
          child: ListTile(
            leading: Icon(
              r.objetoEsInsumo ? Icons.construction : Icons.inventory_2,
            ),
            title: Text(r.objetoNombre),
            subtitle: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (r.objetoDescripcion.isNotEmpty)
                  Text(r.objetoDescripcion, maxLines: 1, overflow: TextOverflow.ellipsis),
                const SizedBox(height: 4),
                Row(
                  children: [
                    Icon(Icons.folder_outlined, size: 14, color: Colors.grey[600]),
                    const SizedBox(width: 4),
                    Text(r.espacioNombre, style: TextStyle(fontSize: 12, color: Colors.grey[600])),
                    const SizedBox(width: 8),
                    Icon(Icons.inventory_2_outlined, size: 14, color: Colors.grey[600]),
                    const SizedBox(width: 4),
                    Text(r.cajaNombre, style: TextStyle(fontSize: 12, color: Colors.grey[600])),
                    const SizedBox(width: 8),
                    Text('x${r.objetoCantidad}',
                        style: TextStyle(fontSize: 12, color: Colors.grey[600])),
                  ],
                ),
              ],
            ),
            trailing: r.tieneFotos ? const Icon(Icons.photo, size: 20) : null,
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (_) => ObjetosScreen(
                    cajaId: r.cajaId,
                    cajaNombre: r.cajaNombre,
                    usuarioId: widget.usuarioId,
                  ),
                ),
              );
            },
          ),
        );
      },
    );
  }
}
