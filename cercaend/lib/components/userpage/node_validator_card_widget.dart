import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../mainpages/userpage/userpage_model.dart';

class NodeValidatorCardWidget extends StatelessWidget {
  final UserpageModel model;
  final VoidCallback onStartNode;
  final VoidCallback onPauseNode;
  final VoidCallback onStopNode;
  final VoidCallback onSyncNode;
  final VoidCallback onRegisterValidator;

  const NodeValidatorCardWidget({
    super.key,
    required this.model,
    required this.onStartNode,
    required this.onPauseNode,
    required this.onStopNode,
    required this.onSyncNode,
    required this.onRegisterValidator,
  });

  Widget _buildControlBtn(
      IconData icon, String label, Color color, VoidCallback onTap) {
    return InkWell(
      onTap: onTap,
      child: Column(
        children: [
          Container(
            padding: const EdgeInsets.all(10),
            decoration: BoxDecoration(
              color: color.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: color.withOpacity(0.3)),
            ),
            child: Icon(icon, color: color, size: 24),
          ),
          const SizedBox(height: 6),
          Text(label, style: const TextStyle(color: Colors.white70, fontSize: 10)),
        ],
      ),
    );
  }

  Widget _buildLogEntry(String text, String level) {
    Color color = Colors.white70;
    if (level == 'ERROR') color = Colors.redAccent;
    if (level == 'WARN') color = Colors.amber;
    if (level == 'SUCCESS') color = Colors.greenAccent;

    return Padding(
      padding: const EdgeInsets.only(bottom: 4),
      child: Text(
        text,
        style: TextStyle(
          color: color,
          fontSize: 10,
          fontFamily: 'monospace',
        ),
      ),
    );
  }

  Widget _buildNetworkMetric(String label, String value, IconData icon) {
    return Container(
      padding: const EdgeInsets.all(10),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: Colors.white.withOpacity(0.05)),
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(icon, color: Colors.blueAccent, size: 20),
          const SizedBox(height: 6),
          Text(value,
              style: GoogleFonts.outfit(
                  color: Colors.white,
                  fontSize: 14,
                  fontWeight: FontWeight.bold)),
          const SizedBox(height: 2),
          Text(label,
              style: const TextStyle(color: Colors.white54, fontSize: 10)),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        // Validator & Node Card - Credit Card Style
        Container(
          width: double.infinity,
          padding: const EdgeInsets.all(20.0),
          decoration: BoxDecoration(
            gradient: const LinearGradient(
              colors: [
                Color(0xFF2D3748),
                Color(0xFF1A202C),
              ],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
            borderRadius: BorderRadius.circular(16.0),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.3),
                blurRadius: 12,
                offset: const Offset(0, 6),
              )
            ],
            border: Border.all(color: Colors.white.withOpacity(0.1)),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Header
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.all(8),
                        decoration: BoxDecoration(
                          color: (model.nodeStatus?.state == 'running')
                              ? Colors.green.withOpacity(0.2)
                              : (model.nodeStatus?.state == 'paused')
                                  ? Colors.amber.withOpacity(0.2)
                                  : Colors.red.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Icon(
                          Icons.dns_rounded,
                          color: (model.nodeStatus?.state == 'running')
                              ? Colors.greenAccent
                              : (model.nodeStatus?.state == 'paused'
                                  ? Colors.amber
                                  : Colors.redAccent),
                          size: 20,
                        ),
                      ),
                      const SizedBox(width: 12),
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            (model.networkStatus?.isValidator == true)
                                ? 'VALIDATOR NODE'
                                : 'OBSERVER MODE',
                            style: GoogleFonts.outfit(
                              color: Colors.white,
                              fontSize: 14,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          Text(
                            (model.nodeStatus?.state == 'running')
                                ? 'Node Running & Earning'
                                : (model.nodeStatus?.state == 'paused'
                                    ? 'Node Paused'
                                    : 'Node Offline'),
                            style: GoogleFonts.outfit(
                                color: Colors.white54, fontSize: 11),
                          ),
                        ],
                      ),
                    ],
                  ),
                  Container(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10, vertical: 4),
                    decoration: BoxDecoration(
                      color: (model.networkStatus?.isValidator == true)
                          ? Colors.green.withOpacity(0.2)
                          : Colors.grey.withOpacity(0.2),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Row(
                      children: [
                        Container(
                          width: 8,
                          height: 8,
                          decoration: BoxDecoration(
                            color: (model.nodeStatus?.state == 'running')
                                ? Colors.greenAccent
                                : (model.nodeStatus?.state == 'paused'
                                    ? Colors.amber
                                    : Colors.grey),
                            shape: BoxShape.circle,
                          ),
                        ),
                        const SizedBox(width: 6),
                        Text(
                          (model.networkStatus?.isValidator == true)
                              ? 'ONLINE'
                              : 'OFFLINE',
                          style: const TextStyle(
                              color: Colors.white70,
                              fontSize: 10,
                              fontWeight: FontWeight.w600),
                        ),
                      ],
                    ),
                  ),
                ],
              ),

              const SizedBox(height: 20),

              // Balance Cards Row
              Row(
                children: [
                  Expanded(
                    child: Container(
                      padding: const EdgeInsets.all(14),
                      decoration: BoxDecoration(
                        gradient: const LinearGradient(
                          colors: [Color(0xFFF1C40F), Color(0xFFF39C12)],
                          begin: Alignment.topLeft,
                          end: Alignment.bottomRight,
                        ),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Row(
                            children: [
                              Icon(Icons.lock, color: Colors.black54, size: 16),
                              SizedBox(width: 4),
                              Text('STAKED',
                                  style: TextStyle(
                                      color: Colors.black54,
                                      fontSize: 10,
                                      fontWeight: FontWeight.w600)),
                            ],
                          ),
                          const SizedBox(height: 6),
                          Text(
                            '${model.networkStatus?.stakeAmount ?? 0}',
                            style: GoogleFonts.outfit(
                              color: Colors.black87,
                              fontSize: 22,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const Text('TCOIN',
                              style: TextStyle(
                                  color: Colors.black54, fontSize: 10)),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Container(
                      padding: const EdgeInsets.all(14),
                      decoration: BoxDecoration(
                        color: Colors.white.withOpacity(0.1),
                        borderRadius: BorderRadius.circular(12),
                        border:
                            Border.all(color: Colors.white.withOpacity(0.1)),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Row(
                            children: [
                              Icon(Icons.emoji_events,
                                  color: Colors.amber, size: 16),
                              SizedBox(width: 4),
                              Text('REWARDS',
                                  style: TextStyle(
                                      color: Colors.white54,
                                      fontSize: 10,
                                      fontWeight: FontWeight.w600)),
                            ],
                          ),
                          const SizedBox(height: 6),
                          Text(
                            '${model.networkStatus?.rewardsEarned ?? 0}',
                            style: GoogleFonts.outfit(
                              color: Colors.white,
                              fontSize: 22,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const Text('TCOIN',
                              style: TextStyle(
                                  color: Colors.white38, fontSize: 10)),
                        ],
                      ),
                    ),
                  ),
                ],
              ),

              const SizedBox(height: 16),

              // Staking Form (only if not validator)
              if (model.networkStatus?.isValidator != true) ...[
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.05),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Become a Validator',
                          style: TextStyle(
                              color: Colors.white,
                              fontWeight: FontWeight.w600,
                              fontSize: 13)),
                      const SizedBox(height: 4),
                      const Text(
                          'Stake TCOIN to earn rewards for securing the network.',
                          style: TextStyle(
                              color: Colors.white54, fontSize: 11)),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: TextFormField(
                              controller: model.stakeAmountController,
                              style: const TextStyle(color: Colors.white),
                              keyboardType: TextInputType.number,
                              decoration: InputDecoration(
                                hintText: 'Amount (min 1000)',
                                hintStyle: const TextStyle(
                                    color: Colors.white38, fontSize: 12),
                                filled: true,
                                fillColor: Colors.black26,
                                contentPadding: const EdgeInsets.symmetric(
                                    horizontal: 12, vertical: 10),
                                border: OutlineInputBorder(
                                    borderRadius: BorderRadius.circular(8),
                                    borderSide: BorderSide.none),
                              ),
                            ),
                          ),
                          const SizedBox(width: 12),
                          ElevatedButton(
                            onPressed: model.isWalletLoading
                                ? null
                                : onRegisterValidator,
                            style: ElevatedButton.styleFrom(
                              backgroundColor: const Color(0xFFF83B46),
                              foregroundColor: Colors.white,
                              padding: const EdgeInsets.symmetric(
                                  horizontal: 20, vertical: 12),
                            ),
                            child: model.isWalletLoading
                                ? const SizedBox(
                                    width: 16,
                                    height: 16,
                                    child: CircularProgressIndicator(
                                        strokeWidth: 2, color: Colors.white))
                                : const Text('STAKE',
                                    style: TextStyle(
                                        fontWeight: FontWeight.w600)),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ],

              // Node Controls (only if validator)
              if (model.networkStatus?.isValidator == true) ...[
                const SizedBox(height: 4),
                const Text('Node Controls',
                    style: TextStyle(color: Colors.white54, fontSize: 11)),
                const SizedBox(height: 8),
                Container(
                  padding: const EdgeInsets.symmetric(
                      vertical: 10, horizontal: 8),
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.05),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      _buildControlBtn(
                          Icons.play_arrow, 'Start', Colors.green, onStartNode),
                      _buildControlBtn(
                          Icons.pause, 'Pause', Colors.amber, onPauseNode),
                      _buildControlBtn(
                          Icons.stop, 'Stop', Colors.red, onStopNode),
                      _buildControlBtn(
                          Icons.sync, 'Sync', Colors.blue, onSyncNode),
                    ],
                  ),
                ),

                const SizedBox(height: 16),

                // Node Logs Section
                const Text('Node Logs',
                    style: TextStyle(color: Colors.white54, fontSize: 11)),
                const SizedBox(height: 8),
                Container(
                  height: 100,
                  width: double.infinity,
                  padding: const EdgeInsets.all(10),
                  decoration: BoxDecoration(
                    color: const Color(0xFF0D1117),
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(color: Colors.white.withOpacity(0.05)),
                  ),
                  child: SingleChildScrollView(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: model.nodeLogs.isEmpty
                          ? [
                              const Text('Waiting for node activity...',
                                  style: TextStyle(
                                      color: Colors.white24, fontSize: 10))
                            ]
                          : model.nodeLogs
                              .map((log) => _buildLogEntry(
                                  '[${log.timestamp.split(' ').last}] ${log.message}',
                                  log.level))
                              .toList(),
                    ),
                  ),
                ),
              ],
            ],
          ),
        ),

        const SizedBox(height: 16),

        // Network Stats Grid
        if (model.networkStatus != null)
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.05),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('Network Overview',
                    style: TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.w600,
                        fontSize: 13)),
                const SizedBox(height: 12),
                GridView.count(
                  shrinkWrap: true,
                  crossAxisCount: 3,
                  crossAxisSpacing: 8,
                  mainAxisSpacing: 8,
                  physics: const NeverScrollableScrollPhysics(),
                  childAspectRatio: 1.3,
                  children: [
                    _buildNetworkMetric('Block Height',
                        '${model.networkStatus!.blockHeight}', Icons.layers),
                    _buildNetworkMetric('TX Pool',
                        '${model.networkStatus!.txPoolSize}', Icons.receipt_long),
                    _buildNetworkMetric('Validators',
                        '${model.networkStatus!.totalValidators}', Icons.people),
                  ],
                ),
              ],
            ),
          ),
      ],
    );
  }
}
