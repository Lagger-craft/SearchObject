import 'package:flutter/material.dart';
import 'features/auth/auth_screen.dart';

class SearchObjectApp extends StatelessWidget {
  const SearchObjectApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'SearchObject',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.indigo),
        useMaterial3: true,
      ),
      home: const AuthScreen(),
    );
  }
}
