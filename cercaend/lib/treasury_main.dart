import 'package:flutter/material.dart';
import 'mainpages/node_dashboard/node_dashboard_widget.dart';

void main() {
  runApp(const TreasuryDashboardApp());
}

class TreasuryDashboardApp extends StatelessWidget {
  const TreasuryDashboardApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'ATLAS Treasury Node Dashboard',
      theme: ThemeData.dark(),
      home: const NodeDashboardWidget(),
      debugShowCheckedModeBanner: false,
    );
  }
}
