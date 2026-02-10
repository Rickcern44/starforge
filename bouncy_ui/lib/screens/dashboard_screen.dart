import 'package:flutter/material.dart';

class DashboardScreen extends StatelessWidget {
  const DashboardScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('League Manager'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () {
              // Handle logout action
              Navigator.of(context).pushNamedAndRemoveUntil(
                '/login',
                (route) =>
                    false, // clear stack so back does not return to dashboard
              );
            },
          ),
        ],
      ),
      body: const Center(
        child: Text('Welcome to your dashboard!'),
      ),
    );
  }
}
