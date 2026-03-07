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
  final TextEditingController _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _loadBlocks();
  }

  Future<void> _loadBlocks() async {
    try {
      final blocks = await _blockchainService.getBlocks();
      if (mounted) {
        setState(() {
          _blocks = blocks;
        });
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error loading blocks: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  void _search() {
    final query = _searchController.text.trim();
    if (query.isEmpty) {
      return;
    }
    final results = _blocks.where((block) {
      return block.hash.contains(query) || block.index.toString() == query;
    }).toList();

    if (results.isNotEmpty) {
      _showBlockDetail(results.first);
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('No block found')),
      );
    }
  }

  void _clearSearch() {
    _searchController.clear();
  }

  void _refreshData() {
    _loadBlocks();
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Data refreshed')),
    );
  }

  void _showBlockDetail(BlockData block) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title:
            Text('Block #${block.index}', style: const TextStyle(fontSize: 20)),
        content: SingleChildScrollView(
          child: ListBody(
            children: <Widget>[
              _buildDetailRow('Hash:', block.hash),
              _buildDetailRow('Previous Hash:', block.prevHash),
              _buildDetailRow(
                  'Timestamp:',
                  DateTime.fromMillisecondsSinceEpoch(block.timestamp * 1000)
                      .toString()),
              _buildDetailRow('Validator:', block.validator),
              _buildDetailRow('Signature:', block.signature),
              const SizedBox(height: 10),
              Text('Transactions (${block.transactions.length})',
                  style: const TextStyle(
                      fontWeight: FontWeight.bold, fontSize: 18)),
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
                style:
                    const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
            TextSpan(
                text: value,
                style: const TextStyle(fontFamily: 'monospace', fontSize: 14)),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    // Note: We are inside a Scaffold provided by NavBarPage if used there,
    // but the original code had its own Scaffold. Nested Scaffolds are okay,
    // but if we put this in a tab, we might want to avoid a second Scaffold
    // if the outer one handles the bottom bar.
    // However, the original code uses a CustomScrollView with SliverAppBar.
    // To preserve the exact look, we keep the Scaffold.
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      body: CustomScrollView(
        slivers: [
          SliverAppBar(
            title: const Text('Blockchain Explorer'),
            floating: true,
            snap: true,
            backgroundColor: const Color(0xFF2D3748),
            automaticallyImplyLeading:
                false, // Don't show back button if in tab
            actions: [
              IconButton(
                icon: const Icon(Icons.refresh, color: Colors.white),
                onPressed: _refreshData,
                tooltip: 'Refresh Data',
              )
            ],
          ),
          SliverToBoxAdapter(
            child: Container(
              padding: const EdgeInsets.all(20),
              child: Column(
                children: [
                  _buildSearchBar(),
                  const SizedBox(height: 20),
                  _buildBlockList(),
                ],
              ),
            ),
          ),
          // Add extra padding at bottom so content isn't hidden behind bottom nav bar
          const SliverPadding(padding: EdgeInsets.only(bottom: 100)),
        ],
      ),
    );
  }

  Widget _buildSearchBar() {
    return Container(
      padding: const EdgeInsets.all(25),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.95),
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.08),
            blurRadius: 40,
            offset: const Offset(0, 20),
          )
        ],
      ),
      child: Column(
        children: [
          TextField(
            controller: _searchController,
            decoration: InputDecoration(
              hintText: 'Search by block hash or block number...',
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(8),
                borderSide: const BorderSide(color: Color(0xFFE0E0E0)),
              ),
              focusedBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(8),
                borderSide: const BorderSide(color: Color(0xFF1E3C72)),
              ),
            ),
          ),
          const SizedBox(height: 15),
          Row(
            children: [
              ElevatedButton.icon(
                onPressed: _search,
                icon: const Icon(Icons.search),
                label: const Text('Search'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF1E3C72),
                ),
              ),
              const SizedBox(width: 10),
              TextButton(onPressed: _clearSearch, child: const Text('Clear')),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildBlockList() {
    return Container(
      padding: const EdgeInsets.all(30),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.95),
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.08),
            blurRadius: 40,
            offset: const Offset(0, 20),
          )
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Latest Blocks',
            style: TextStyle(fontSize: 24, fontWeight: FontWeight.w600),
          ),
          const Divider(height: 40),
          _blocks.isEmpty
              ? const Center(child: CircularProgressIndicator())
              : ListView.builder(
                  shrinkWrap: true,
                  physics: const NeverScrollableScrollPhysics(),
                  itemCount: _blocks.length,
                  itemBuilder: (context, index) {
                    final block = _blocks[index];
                    return _buildBlockItem(block);
                  },
                ),
        ],
      ),
    );
  }

  Widget _buildBlockItem(BlockData block) {
    return Card(
      margin: const EdgeInsets.only(bottom: 15),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
        side: const BorderSide(color: Color(0xFF4299E1), width: 2),
      ),
      child: InkWell(
        onTap: () => _showBlockDetail(block),
        child: Padding(
          padding: const EdgeInsets.all(20.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'Block #${block.index}',
                    style: const TextStyle(
                        fontSize: 18, fontWeight: FontWeight.bold),
                  ),
                  Text(
                    '${block.hash.length > 16 ? block.hash.substring(0, 16) : block.hash}...',
                    style: const TextStyle(
                        fontFamily: 'monospace', color: Colors.grey),
                  ),
                ],
              ),
              const SizedBox(height: 10),
              GridView(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
                  maxCrossAxisExtent: 200,
                  mainAxisSpacing: 10,
                  crossAxisSpacing: 10,
                  childAspectRatio: 3,
                ),
                children: [
                  _buildDetailItem(
                      'Timestamp',
                      DateTime.fromMillisecondsSinceEpoch(
                              block.timestamp * 1000)
                          .toIso8601String()),
                  _buildDetailItem(
                      'Transactions', block.transactions.length.toString()),
                  _buildDetailItem('Validator',
                      '${block.validator.length > 16 ? block.validator.substring(0, 16) : block.validator}...'),
                  // _buildDetailItem('Size', '${block.sizeInBytes} bytes'), // Removed sizeInBytes as it's not in BlockData yet
                ],
              )
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildDetailItem(String label, String value) {
    return Container(
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(4),
        border: Border.all(color: const Color(0xFFE0E0E0)),
      ),
      child: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(label,
                style: const TextStyle(color: Colors.grey, fontSize: 12)),
            Text(value,
                style:
                    const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
          ],
        ),
      ),
    );
  }
}

// Removed ExplorerBlock and ExplorerTransaction classes
