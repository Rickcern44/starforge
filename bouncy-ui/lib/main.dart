import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:logging/logging.dart';
import 'package:ui_test/screens/home_screen.dart';
import 'package:ui_test/screens/login_screen.dart';
import 'package:ui_test/screens/register_screen.dart';
import 'package:ui_test/screens/create_event_screen.dart';
import 'package:ui_test/screens/ledger_screen.dart';

void main() {
  // Configure logging
  Logger.root.level = Level.ALL; // Log everything
  Logger.root.onRecord.listen((record) {
    if (kDebugMode) {
      print('${record.level.name}: ${record.time}: ${record.loggerName}: ${record.message}');
    }
  });

  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'League Manager',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.indigo),
        useMaterial3: true,
      ),
      initialRoute: "/",
      routes: {
        "/auth/login": (context) => const LoginScreen(),
        "/auth/register": (context) => const RegisterScreen(),
        "/": (context) => const HomeScreen(),
        "/events/create": (context) => const CreateEventScreen(),
        "/ledger": (context) => const LedgerScreen(),
      },
    );
  }
}
