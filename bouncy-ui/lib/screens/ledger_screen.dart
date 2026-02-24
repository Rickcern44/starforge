import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class LedgerScreen extends StatefulWidget {
  const LedgerScreen({super.key});

  @override
  State<LedgerScreen> createState() => _LedgerScreenState();
}

class _LedgerScreenState extends State<LedgerScreen> {
  final List<Map<String, dynamic>> items = [
    {
      'id': '1',
      'type': 'charge',
      'title': 'Game vs Red Dragons',
      'date': 'Oct 10, 2023',
      'amount': 15.0,
      'status': 'unpaid'
    },
    {
      'id': '2',
      'type': 'charge',
      'title': 'Game vs Blue Hawks',
      'date': 'Oct 17, 2023',
      'amount': 15.0,
      'status': 'unpaid'
    },
    {
      'id': '3',
      'type': 'payment',
      'title': 'Bulk Payment',
      'date': 'Oct 05, 2023',
      'amount': 45.0,
      'status': 'completed'
    },
    {
      'id': '4',
      'type': 'charge',
      'title': 'Court Rental Fee',
      'date': 'Oct 01, 2023',
      'amount': 10.0,
      'status': 'paid'
    },
  ];

  Set<String> selectedChargeIds = {};

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final unpaidCharges = items
        .where((i) => i['type'] == 'charge' && i['status'] == 'unpaid')
        .toList();
    final double selectedTotal = items
        .where((i) => selectedChargeIds.contains(i['id']))
        .fold(0, (sum, item) => sum + item['amount']);

    return Scaffold(
      appBar: AppBar(
        title: Text('League Ledger',
            style: GoogleFonts.poppins(fontWeight: FontWeight.w600)),
      ),
      body: Column(
        children: [
          _buildBalanceHeader(colorScheme),
          Expanded(
            child: ListView.builder(
              itemCount: items.length,
              itemBuilder: (context, index) {
                final item = items[index];
                final isCharge = item['type'] == 'charge';
                final isUnpaid = item['status'] == 'unpaid';
                final isSelected = selectedChargeIds.contains(item['id']);

                return ListTile(
                  leading: CircleAvatar(
                    backgroundColor: isCharge
                        ? (isUnpaid
                            ? Colors.orange.withOpacity(0.1)
                            : Colors.red.withOpacity(0.1))
                        : Colors.green.withOpacity(0.1),
                    child: Icon(
                      isCharge ? Icons.receipt_long : Icons.payment,
                      color: isCharge
                          ? (isUnpaid ? Colors.orange : Colors.red)
                          : Colors.green,
                    ),
                  ),
                  title: Text(item['title'],
                      style: const TextStyle(fontWeight: FontWeight.bold)),
                  subtitle: Text(item['date']),
                  trailing: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        '${isCharge ? "-" : "+"}\$${item['amount'].toStringAsFixed(2)}',
                        style: TextStyle(
                          fontWeight: FontWeight.bold,
                          color: isCharge ? Colors.red : Colors.green,
                        ),
                      ),
                      if (isUnpaid)
                        Checkbox(
                          value: isSelected,
                          onChanged: (val) {
                            setState(() {
                              if (val == true) {
                                selectedChargeIds.add(item['id']);
                              } else {
                                selectedChargeIds.remove(item['id']);
                              }
                            });
                          },
                        ),
                    ],
                  ),
                );
              },
            ),
          ),
          if (selectedChargeIds.isNotEmpty)
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: colorScheme.surface,
                boxShadow: const [
                  BoxShadow(
                      color: Colors.black12,
                      blurRadius: 10,
                      offset: Offset(0, -2))
                ],
              ),
              child: SafeArea(
                child: Row(
                  children: [
                    Expanded(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text('Selected Total',
                              style: TextStyle(
                                  color: colorScheme.onSurfaceVariant)),
                          Text('\$${selectedTotal.toStringAsFixed(2)}',
                              style: const TextStyle(
                                  fontSize: 20, fontWeight: FontWeight.bold)),
                        ],
                      ),
                    ),
                    ElevatedButton(
                      onPressed: () {
                        // Logic to "pay" selected charges
                        ScaffoldMessenger.of(context).showSnackBar(SnackBar(
                            content: Text(
                                'Payment of \$${selectedTotal.toStringAsFixed(2)} processed!')));
                        setState(() {
                          selectedChargeIds.clear();
                        });
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: colorScheme.primary,
                        foregroundColor: colorScheme.onPrimary,
                        padding: const EdgeInsets.symmetric(
                            horizontal: 24, vertical: 12),
                      ),
                      child: const Text('Make Payment'),
                    ),
                  ],
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildBalanceHeader(ColorScheme colorScheme) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(24),
      color: colorScheme.primaryContainer.withOpacity(0.3),
      child: Column(
        children: [
          Text('Outstanding Balance',
              style: TextStyle(color: colorScheme.onSurfaceVariant)),
          const SizedBox(height: 4),
          Text('\$30.00',
              style: GoogleFonts.poppins(
                  fontSize: 32,
                  fontWeight: FontWeight.bold,
                  color: Colors.red[700])),
        ],
      ),
    );
  }
}
