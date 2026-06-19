import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import '../../services/blockchain/blockchain_service.dart';
import '../../app_state.dart';

class NodeDashboardWidget extends StatefulWidget {
  const NodeDashboardWidget({super.key});

  static const String routeName = 'NodeDashboard';
  static const String routePath = '/nodeDashboard';

  @override
  NodeDashboardWidgetState createState() => NodeDashboardWidgetState();
}

class NodeDashboardWidgetState extends State<NodeDashboardWidget> {
  final BlockchainService _blockchainService = BlockchainService();
  
  bool _isLoading = true;
  NodeStatusResponse? _nodeStatus;
  NetworkStatusResponse? _networkStatus;
  List<NodeLogEntry> _logs = [];
  
  Timer? _pollingTimer;
  final ScrollController _logScrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    _refreshData();
    if (FFAppState().isLocalNodeMode) {
      _pollLogs();
    }
    _pollingTimer = Timer.periodic(const Duration(seconds: 3), (_) {
      if (FFAppState().isLocalNodeMode) {
        _pollLogs();
      }
      _refreshData(isBackground: true);
    });
  }

  @override
  void dispose() {
    _pollingTimer?.cancel();
    _logScrollController.dispose();
    super.dispose();
  }

  Future<void> _refreshData({bool isBackground = false}) async {
    if (!mounted) return;
    if (!isBackground) setState(() => _isLoading = true);

    try {
      final futures = await Future.wait([
        _blockchainService.getNodeStatus(),
        _blockchainService.getNetworkStatus(),
      ]);

      if (mounted) {
        setState(() {
          _nodeStatus = futures[0] as NodeStatusResponse;
          _networkStatus = futures[1] as NetworkStatusResponse;
          _isLoading = false;
        });
      }
    } catch (e) {
      if (mounted && !isBackground) {
        setState(() => _isLoading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error loading dashboard data: $e'), backgroundColor: Colors.red),
        );
      }
    }
  }

  Future<void> _pollLogs() async {
    if (!FFAppState().isLocalNodeMode) return;
    try {
      final logsReq = await _blockchainService.getNodeLogs(limit: 50);
      if (mounted) {
        setState(() {
          _logs = logsReq.logs;
        });
        if (_logScrollController.hasClients) {
          _logScrollController.animateTo(
            _logScrollController.position.maxScrollExtent,
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeOut,
          );
        }
      }
    } catch (e) {
      // Ignore background log poll errors
    }
  }

  Future<void> _executeNodeCommand(String action) async {
    try {
      if (action == 'start') { await _blockchainService.startNode(); }
      else if (action == 'stop') { await _blockchainService.stopNode(); }
      else if (action == 'pause') { await _blockchainService.pauseNode(); }
      else if (action == 'sync') { await _blockchainService.syncNode(); }
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Command $action executed successfully'), backgroundColor: Colors.green),
      );
      _refreshData();
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error executing $action: $e'), backgroundColor: Colors.red),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF1E1E2C), // Dark theme
      appBar: AppBar(
        title: const Text('Treasury Node Dashboard', style: TextStyle(fontWeight: FontWeight.bold, color: Colors.white)),
        backgroundColor: const Color(0xFF13131A),
        iconTheme: const IconThemeData(color: Colors.white),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh, color: Colors.white),
            onPressed: () {
               _refreshData();
               _pollLogs();
            },
          )
        ],
      ),
      body: _isLoading && _nodeStatus == null
          ? const Center(child: CircularProgressIndicator())
          : _buildLayout(),
    );
  }

  Widget _buildLayout() {
    final isDesktop = MediaQuery.of(context).size.width >= 992;
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: isDesktop ? Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(flex: 2, child: _buildLeftPanel()),
          const SizedBox(width: 24),
          Expanded(flex: 3, child: _buildRightPanel()),
        ],
      ) : Column(
        children: [
          _buildLeftPanel(),
          const SizedBox(height: 24),
          _buildRightPanel(),
        ],
      ),
    );
  }

  Widget _buildLeftPanel() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        _buildControlsCard(),
        const SizedBox(height: 24),
        _buildMetricsCard(),
      ],
    );
  }

  Widget _buildRightPanel() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        _buildTerminalCard(),
        const SizedBox(height: 24),
        _buildValidatorCard(),
      ],
    );
  }

  Widget _buildControlsCard() {
    final isLocal = FFAppState().isLocalNodeMode;
    return Card(
      color: const Color(0xFF27293D),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('Node Controls', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
            const SizedBox(height: 16),
            if (kIsWeb && isLocal) ...[
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.orangeAccent.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(color: Colors.orangeAccent.withValues(alpha: 0.2)),
                ),
                child: const Row(
                  children: [
                    Icon(Icons.info_outline, color: Colors.orangeAccent, size: 20),
                    SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        'Running in Web Browser: Web apps cannot spawn or directly control background node processes. Please download the ATLAS Desktop app or run the node daemon manually on your machine.',
                        style: TextStyle(color: Colors.white70, fontSize: 12),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 16),
            ],
            if (!isLocal) ...[
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.redAccent.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(color: Colors.redAccent.withValues(alpha: 0.2)),
                ),
                child: const Row(
                  children: [
                    Icon(Icons.warning_amber_rounded, color: Colors.redAccent, size: 20),
                    SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        'Local node controls are disabled when connected to a remote seed node.',
                        style: TextStyle(color: Colors.white70, fontSize: 12),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 16),
            ],
            Wrap(
              spacing: 12,
              runSpacing: 12,
              children: [
                _buildButton(Icons.play_arrow, 'Start', Colors.green, isLocal ? () => _executeNodeCommand('start') : null),
                _buildButton(Icons.pause, 'Pause', Colors.orange, isLocal ? () => _executeNodeCommand('pause') : null),
                _buildButton(Icons.stop, 'Stop', Colors.red, isLocal ? () => _executeNodeCommand('stop') : null),
                _buildButton(Icons.sync, 'Sync', Colors.blue, isLocal ? () => _executeNodeCommand('sync') : null),
              ],
            )
          ],
        ),
      ),
    );
  }

  Widget _buildButton(IconData icon, String label, Color color, VoidCallback? onTap) {
    return ElevatedButton.icon(
      icon: Icon(icon, color: onTap == null ? Colors.white30 : Colors.white, size: 18),
      label: Text(label, style: TextStyle(color: onTap == null ? Colors.white30 : Colors.white)),
      style: ElevatedButton.styleFrom(
        backgroundColor: onTap == null ? const Color(0xFF3A3B4C) : color,
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      ),
      onPressed: onTap,
    );
  }

  Widget _buildMetricsCard() {
    return Card(
      color: const Color(0xFF27293D),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('Network & System Metrics', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
            const SizedBox(height: 16),
            _buildMetricRow('Node State', _nodeStatus?.state.toUpperCase() ?? 'UNKNOWN', isHighlight: true),
            _buildMetricRow('Uptime', _nodeStatus?.uptime ?? '-'),
            const Divider(color: Colors.white24, height: 24),
            _buildMetricRow('Block Height', '\${_networkStatus?.blockHeight ?? 0}'),
            _buildMetricRow('TX Pool Size', '\${_networkStatus?.txPoolSize ?? 0}'),
          ],
        ),
      ),
    );
  }

  Widget _buildMetricRow(String label, String value, {bool isHighlight = false}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(color: Colors.white70, fontSize: 15)),
          Text(value, style: TextStyle(
            color: isHighlight ? Colors.greenAccent : Colors.white, 
            fontSize: 15, 
            fontWeight: isHighlight ? FontWeight.bold : FontWeight.normal
          )),
        ],
      ),
    );
  }

  Widget _buildTerminalCard() {
    final isLocal = FFAppState().isLocalNodeMode;
    return Card(
      color: const Color(0xFF13131A), // Darker for terminal
      shape: RoundedRectangleBorder(side: const BorderSide(color: Colors.white12), borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Row(
              children: [
                Icon(Icons.terminal, color: Colors.greenAccent),
                SizedBox(width: 8),
                Text('Terminal Window', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
              ],
            ),
            const SizedBox(height: 16),
            Container(
              height: 300,
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Colors.black,
                borderRadius: BorderRadius.circular(8),
              ),
              child: isLocal ? ListView.builder(
                controller: _logScrollController,
                itemCount: _logs.length,
                itemBuilder: (context, index) {
                  final log = _logs[index];
                  Color logColor = Colors.white;
                  if (log.level.toLowerCase().contains('error') || log.level.toLowerCase().contains('err')) { logColor = Colors.redAccent; }
                  else if (log.level.toLowerCase().contains('warn')) { logColor = Colors.orangeAccent; }
                  else if (log.level.toLowerCase().contains('info')) { logColor = Colors.greenAccent; }

                  return Padding(
                    padding: const EdgeInsets.symmetric(vertical: 2.0),
                    child: Text('[${log.level.toUpperCase()}] ${log.timestamp}: ${log.message}', style: TextStyle(color: logColor, fontFamily: 'monospace', fontSize: 12)),
                  );
                },
              ) : const Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(Icons.lock_outline, color: Colors.white30, size: 48),
                    SizedBox(height: 12),
                    Text(
                      'Terminal logs are disabled in remote node mode.',
                      style: TextStyle(color: Colors.white30, fontSize: 14),
                      textAlign: TextAlign.center,
                    ),
                  ],
                ),
              ),
            )
          ],
        ),
      ),
    );
  }

  Widget _buildValidatorCard() {
    final isVal = _networkStatus?.isValidator ?? false;
    return Card(
      color: const Color(0xFF27293D),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('Validator Status', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
            const SizedBox(height: 16),
            Row(
              children: [
                Container(
                  width: 12, height: 12,
                  decoration: BoxDecoration(
                    color: isVal ? Colors.greenAccent : Colors.redAccent,
                    shape: BoxShape.circle
                  ),
                ),
                const SizedBox(width: 12),
                Text(isVal ? 'Active Validator' : 'Not a Validator', style: const TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.bold)),
              ],
            ),
            const SizedBox(height: 16),
            _buildMetricRow('Validator Address', _networkStatus?.validatorAddress.isNotEmpty == true ? '${_networkStatus!.validatorAddress.substring(0, 16)}...' : 'Not registered'),
            _buildMetricRow('Total Validators', '\${_networkStatus?.totalValidators ?? 0}'),
            _buildMetricRow('Staked Amount', '\${_networkStatus?.stakeAmount ?? 0} TCOIN'),
            _buildMetricRow('Rewards Earned', '\${_networkStatus?.rewardsEarned ?? 0} TCOIN'),
          ],
        ),
      ),
    );
  }
}
