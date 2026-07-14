package com.example.search_object

import android.content.Context
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel
import bridge.bridge.Bridge

class MainActivity : FlutterActivity() {
    private val CHANNEL = "searchobject/bridge"
    private var app: bridge.bridge.App? = null

    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)

        val dbPath = getDatabasePath("searchobject.db").absolutePath
        val imageDir = filesDir.absolutePath + "/images"

        app = Bridge.new_(dbPath, imageDir)

        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, CHANNEL)
            .setMethodCallHandler { call, result ->
                try {
                    val response = when (call.method) {
                        "crearUsuario" -> app!!.crearUsuario(
                            call.argument<String>("nombre")!!,
                            call.argument<String>("email")!!
                        )
                        "obtenerUsuario" -> app!!.obtenerUsuario(
                            call.argument<String>("id")!!
                        )
                        "crearEspacio" -> app!!.crearEspacio(
                            call.argument<String>("json")!!
                        )
                        "listarEspacios" -> app!!.listarEspacios(
                            call.argument<String>("usuarioId")!!
                        )
                        "obtenerEspacio" -> app!!.obtenerEspacio(
                            call.argument<String>("id")!!
                        )
                        "actualizarEspacio" -> app!!.actualizarEspacio(
                            call.argument<String>("json")!!
                        )
                        "eliminarEspacio" -> app!!.eliminarEspacio(
                            call.argument<String>("id")!!
                        )
                        "crearCaja" -> app!!.crearCaja(
                            call.argument<String>("json")!!
                        )
                        "listarCajas" -> app!!.listarCajas(
                            call.argument<String>("espacioId")!!
                        )
                        "obtenerCaja" -> app!!.obtenerCaja(
                            call.argument<String>("id")!!
                        )
                        "eliminarCaja" -> app!!.eliminarCaja(
                            call.argument<String>("id")!!
                        )
                        "crearObjeto" -> app!!.crearObjeto(
                            call.argument<String>("json")!!
                        )
                        "listarObjetos" -> app!!.listarObjetos(
                            call.argument<String>("cajaId")!!
                        )
                        "obtenerObjeto" -> app!!.obtenerObjeto(
                            call.argument<String>("id")!!
                        )
                        "moverObjeto" -> app!!.moverObjeto(
                            call.argument<String>("json")!!
                        )
                        "eliminarObjeto" -> app!!.eliminarObjeto(
                            call.argument<String>("id")!!
                        )
                        "buscarObjetos" -> app!!.buscarObjetos(
                            call.argument<String>("usuarioId")!!,
                            call.argument<String>("termino")!!
                        )
                        "buscar" -> app!!.buscar(
                            call.argument<String>("usuarioId")!!,
                            call.argument<String>("termino")!!
                        )
                        "resumen" -> app!!.resumen(
                            call.argument<String>("usuarioId")!!
                        )
                        "dashboard" -> app!!.dashboard(
                            call.argument<String>("usuarioId")!!
                        )
                        "evaluarAlertas" -> app!!.evaluarAlertas(
                            call.argument<String>("usuarioId")!!
                        )
                        "listarAlertas" -> app!!.listarAlertas(
                            call.argument<String>("jsonLeidas")!!
                        )
                        "marcarAlertaLeida" -> app!!.marcarAlertaLeida(
                            call.argument<String>("id")!!
                        )
                        "resolverAlerta" -> app!!.resolverAlerta(
                            call.argument<String>("id")!!
                        )
                        "exportarJSON" -> app!!.exportarJSON(
                            call.argument<String>("usuarioId")!!
                        )
                        "exportarCSV" -> app!!.exportarCSV(
                            call.argument<String>("usuarioId")!!
                        )
                        "agregarImagen" -> app!!.agregarImagen(
                            call.argument<String>("objetoId")!!,
                            call.argument<ByteArray>("imageBytes")!!
                        )
                        "agregarImagenConArea" -> app!!.agregarImagenConArea(
                            call.argument<String>("objetoId")!!,
                            call.argument<ByteArray>("imageBytes")!!,
                            call.argument<String>("jsonArea")!!
                        )
                        "listarImagenes" -> app!!.listarImagenes(
                            call.argument<String>("objetoId")!!
                        )
                        "eliminarImagen" -> app!!.eliminarImagen(
                            call.argument<String>("id")!!
                        )
                        "imageDir" -> app!!.imageDir()
                        "pathParaObjeto" -> app!!.pathParaObjeto(
                            call.argument<String>("objetoId")!!
                        )
                        "close" -> {
                            app!!.close()
                            null
                        }
                        else -> {
                            result.notImplemented()
                            return@setMethodCallHandler
                        }
                    }
                    result.success(response)
                } catch (e: Exception) {
                    result.error("BRIDGE_ERROR", e.message, null)
                }
            }
    }

    override fun onDestroy() {
        app?.close()
        super.onDestroy()
    }
}
