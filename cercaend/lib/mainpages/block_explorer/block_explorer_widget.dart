import 'package:flutter/material.dart';
import '/services/blockchain/blockchain_service.dart';

class BlockExplorerWidget extends StatefulWidget {
  const BlockExplorerWidget({super.key});

  @override
  BlockExplorerWidgetState createState() => BlockExplorerWidgetState();
}

class BlockExplorerWidgetState extends State<BlockExplorerWidget> {
  final BlockchainService _blockchainService = BlockchainService();
  List<BlockData> _blocks = [];
  NetworkStatusResponse? _networkStatus;
  NodeStatusResponse? _nodeStatus;
  int _peerCount = 0;
  
  bool _isLoading = true;
  final TextEditingController _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _refreshData();
  }

  Future<void> _refreshData() async {
    if (!mounted) return;
    setState(() => _isLoading = true);
    
    try {
      final futures = await Future.wait([
        _blockchainService.getBlocks(limit: 10),
        _blockchainService.getNetworkStatus(),
        _blockchainService.getNodeStatus(),
        _blockchainService.getPeers()
      ]);
      
      if (mounted) {
        setState(() {
          _blocks = futures[0] as List<BlockData>;
          _networkStatus = futures[1] as NetworkStatusResponse;
          _nodeStatus = futures[2] as NodeStatusResponse;
          _peerCount = (futures[3] as Map<String, dynamic>)['count'] ?? 0;
          _isLoading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _isLoading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error loading network data: $e'), backgroundColor: Colors.red),
        );
      }
    }
  }

  void _search() {
    final query = _searchController.text.trim();
    if (query.isEmpty) return;
    final results = _blocks.where((block) {
      return block.hash.contains(query) || block.index.toString() == query;
    }).toList();

    if (results.isNotEmpty) {
      _showBlockDetail(results.first);
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('No block found in recent blocks.')),
      );
    }
  }

  void _showBlockDetail(BlockData block) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Block #${block.index}', style: const TextStyle(fontSize: 20)),
        content: SingleChildScrollView(
          child: ListBody(
            children: <Widget>[
              _buildDetailRow('Hash:', block.hash),
              _buildDetailRow('Previous Hash:', block.prevHash),
              _buildDetailRow(
                  'Timestamp:',
                  DateTime.fromMillisecondsSinceEpoch(block.timestamp * 1000)
                      .toLocal()
                      .toString()),
              _buildDetailRow('Validator:', block.validator),
              _buildDetailRow('Signature:', block.signature),
              const SizedBox(height: 10),
              Text('Transactions (${block.transactions.length})',
                  style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 18)),
              const SizedBox(height: 5),
              ...block.transactions.map((tx) => Card(
                    margin: const EdgeInsets.symmetric(vertical: 4),
                    child: Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          _buildDetailRow('From:', tx.sender),
                          _buildDetailRow('To:', tx.recipient),
                          _buildDetailRow('Amount:', '${tx.amount} tokens'),
                          _buildDetailRow('Fee:', '${tx.fee} tokens'),
                        ],
                      ),
                    ),
                  )),
            ],
          ),
        ),
        actions: <Widget>[
          TextButton(
            child: const Text('Close'),
            onPressed: () {
              Navigator.of(context).pop();
            },
          ),
        ],
      ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: RichText(
        text: TextSpan(
          style: DefaultTextStyle.of(context).style,
          children: <TextSpan>[
            TextSpan(
                text: '$label ',
                style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Colors.black87)),
            TextSpan(
                text: value.length > 32 ? '${value.substring(0, 32)}...' : value,
                style: const TextStyle(fontFamily: 'monospace', fontSize: 14, color: Colors.black54)),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final isDesktop = MediaQuery.of(context).size.width >= 992;
    
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(
        title: const Text('Explorer & Health', style: TextStyle(fontWeight: FontWeight.w600)),
        backgroundColor: const Color(0xFF1A202C),
        foregroundColor: Colors.white,
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh, color: Colors.white),
            onPressed: _refreshData,
            tooltip: 'Refresh Data',
          )
        ],
      ),
      body: _isLoading 
        ? const Center(child: CircularProgressIndicator()) 
        : SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: isDesktop ? _buildDesktopLayout() : _buildMobileLayout(),
          ),
    );
  }

  Widget _buildDesktopLayout() {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Expanded(
          flex: 4,
          child: Column(
             crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              _buildMetricsPanel(),
              const SizedBox(height: 24),
              _buildNetworkConfigPanel(),
            ],
          ),
        ),
        const SizedBox(width: 24),
        Expanded(
          flex: 6,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              _buildSearchBar(),
              const SizedBox(height: 24),
              _buildBlockList(),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildMobileLayout() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        _buildMetricsPanel(),
        const SizedBox(height: 24),
        _buildNetworkConfigPanel(),
        const SizedBox(height: 24),
        _buildSearchBar(),
        const SizedBox(height: 24),
        _buildBlockList(),
        const SizedBox(height: 100),
      ],
    );
  }

  Widget _buildMetricsPanel() {
    return Container(
      padding: const EdgeInsets.all(28),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF2B6CB0), Color(0xFF2C5282)],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(24),
        boxShadow: [
          BoxShadow(
            color: const Color(0xFF2B6CB0).withValues(alpha: 0.25),
            blurRadius: 24,
            offset: const Offset(0, 10),
          )
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Row(
            children: [
              Icon(Icons.monitor_heart, color: Colors.white, size: 28),
              SizedBox(width: 12),
              Text(
                'Live Network Health',
                style: TextStyle(color: Colors.white, fontSize: 22, fontWeight: FontWeight.bold),
              ),
            ],
          ),
          const SizedBox(height: 24),
          _buildMetricStat('Status', _nodeStatus?.state.toUpperCase() ?? 'UNKNOWN', Icons.fiber_manual_record, color: _nodeStatus?.state == 'running' ? const Color(0xFF68D391) : Colors.amber),
          const SizedBox(height: 20),
          _buildMetricStat('Current Height', '${_networkStatus?.blockHeight ?? 0} Blocks', Icons.layers),
          const SizedBox(height: 20),
          _buildMetricStat('TX Pool Size', '${_networkStatus?.txPoolSize ?? 0} Pending', Icons.swap_horiz_rounded),
          const SizedBox(height: 20),
          _buildMetricStat('Connected Peers', '$_peerCount Nodes', Icons.hub_rounded),
          const SizedBox(height: 20),
          _buildMetricStat('Active Validators', '${_networkStatus?.totalValidators ?? 0} Validators', Icons.verified_user),
        ],
      ),
    );
  }

  Widget _buildNetworkConfigPanel() {
    return Container(
      padding: const EdgeInsets.all(28),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: const Color(0xFFE2E8F0)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.03),
            blurRadius: 20,
            offset: const Offset(0, 8),
          )
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Node Configuration',
            style: TextStyle(color: Color(0xFF2D3748), fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const Divider(height: 32, color: Color(0xFFE2E8F0)),
          _buildInfoRow('Primary Node Uptime', _nodeStatus?.uptime ?? '0s', Icons.timer),
          const SizedBox(height: 16),
          _buildInfoRow('Consensus Engine', 'Pos-DevNet', Icons.extension),
          const SizedBox(height: 16),
          _buildInfoRow('Primary Node Address', '${_nodeStatus?.validatorAddress.substring(0, 12) ?? 'Unknown'}...', Icons.vpn_key),
        ],
      ),
    );
  }

  Widget _buildMetricStat(String label, String value, IconData icon, {Color color = Colors.white}) {
    return Row(
      children: [
        Container(
          padding: const EdgeInsets.all(8),
          decoration: BoxDecoration(
            color: Colors.white.withValues(alpha: 0.1),
            borderRadius: BorderRadius.circular(10),
          ),
          child: Icon(icon, color: color.withValues(alpha: 0.9), size: 20),
        ),
        const SizedBox(width: 16),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(label, style: const TextStyle(color: Colors.white70, fontSize: 13, letterSpacing: 0.5)),
            const SizedBox(height: 2),
            Text(value, style: TextStyle(color: color, fontSize: 18, fontWeight: FontWeight.bold)),
          ],
        )
      ],
    );
  }

  Widget _buildInfoRow(String label, String value, IconData icon) {
    return Row(
      children: [
        Container(
          padding: const EdgeInsets.all(10),
          decoration: BoxDecoration(
            color: const Color(0xFFF7FAFC),
            border: Border.all(color: const Color(0xFFE2E8F0)),
            borderRadius: BorderRadius.circular(10),
          ),
          child: Icon(icon, color: const Color(0xFF718096), size: 18),
        ),
        const SizedBox(width: 16),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(label, style: const TextStyle(color: Color(0xFF718096), fontSize: 13)),
              const SizedBox(height: 2),
              Text(value, style: const TextStyle(color: Color(0xFF2D3748), fontSize: 15, fontWeight: FontWeight.w600)),
            ],
          ),
        )
      ],
    );
  }

  Widget _buildSearchBar() {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: const Color(0xFFE2E8F0)),
      ),
      child: Row(
        children: [
          Expanded(
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: 'Search by block hash or height...',
                hintStyle: const TextStyle(color: Color(0xFFA0AEC0)),
                prefixIcon: const Icon(Icons.search, color: Color(0xFFA0AEC0)),
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(20), borderSide: BorderSide.none),
                filled: true,
                fillColor: Colors.transparent,
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: ElevatedButton(
              onPressed: _search,
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFF2B6CB0),
                padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                elevation: 0,
              ),
              child: const Text('Search', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBlockList() {
    return Container(
      padding: const EdgeInsets.all(28),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: const Color(0xFFE2E8F0)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.02),
            blurRadius: 20,
            offset: const Offset(0, 8),
          )
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                'Latest Blocks',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Color(0xFF2D3748)),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                decoration: BoxDecoration(
                  color: const Color(0xFFEBF8FF),
                  borderRadius: BorderRadius.circular(20),
                  border: Border.all(color: const Color(0xFFBEE3F8))
                ),
                child: Row(
                  children: [
                    Container(
                      width: 8,
                      height: 8,
                      decoration: const BoxDecoration(color: Color(0xFF3182CE), shape: BoxShape.circle),
                    ),
                    const SizedBox(width: 8),
                    const Text('Live Ledger', style: TextStyle(color: Color(0xFF2B6CB0), fontWeight: FontWeight.bold, fontSize: 12)),
                  ],
                ),
              ),
            ],
          ),
          const Divider(height: 32, color: Color(0xFFE2E8F0)),
          if (_blocks.isEmpty) 
            const Padding(
              padding: EdgeInsets.all(32.0),
              child: Center(child: Text('No blocks on network yet. Waiting for Genesis...', style: TextStyle(color: Color(0xFFA0AEC0)))),
            )
          else 
            ListView.separated(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              itemCount: _blocks.length,
              separatorBuilder: (context, index) => const SizedBox(height: 12),
              itemBuilder: (context, index) {
                return _buildBlockItem(_blocks[index]);
              },
            ),
        ],
      ),
    );
  }

  Widget _buildBlockItem(BlockData block) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: const Color(0xFFEDF2F7)),
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          borderRadius: BorderRadius.circular(16),
          onTap: () => _showBlockDetail(block),
          hoverColor: const Color(0xFFF7FAFC),
          child: Padding(
            padding: const EdgeInsets.all(20.0),
            child: Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(14),
                  decoration: BoxDecoration(
                    color: const Color(0xFFEDF2F7),
                    borderRadius: BorderRadius.circular(14),
                  ),
                  child: const Icon(Icons.view_in_ar, color: Color(0xFF4A5568), size: 24),
                ),
                const SizedBox(width: 20),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text('Block #${block.index}', style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Color(0xFF2D3748))),
                      const SizedBox(height: 6),
                      Text(block.hash, 
                          maxLines: 1, 
                          overflow: TextOverflow.ellipsis,
                          style: const TextStyle(fontFamily: 'monospace', color: Color(0xFFA0AEC0), fontSize: 13)),
                    ],
                  ),
                ),
                const SizedBox(width: 16),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                      decoration: BoxDecoration(
                        color: const Color(0xFFF7FAFC),
                        borderRadius: BorderRadius.circular(10),
                      ),
                      child: Text('${block.transactions.length} TXs', style: const TextStyle(fontWeight: FontWeight.w600, color: Color(0xFF4A5568), fontSize: 13)),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      DateTime.fromMillisecondsSinceEpoch(block.timestamp * 1000).toLocal().toString().split('.')[0],
                      style: const TextStyle(color: Color(0xFFA0AEC0), fontSize: 12),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
