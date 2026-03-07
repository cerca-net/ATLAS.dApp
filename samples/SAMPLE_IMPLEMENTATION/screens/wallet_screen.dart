import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';
import '../services/auth_service.dart';
import '../services/wallet_service.dart';
import '../models/wallet.dart';
import '../models/transaction.dart';

class WalletScreen extends StatefulWidget {
  @override
  _WalletScreenState createState() => _WalletScreenState();
}

class _WalletScreenState extends State<WalletScreen> {
  final _walletService = GetIt.instance<WalletService>();
  final _authService = GetIt.instance<AuthService>();
  
  late Future<Wallet> _walletFuture;
  bool _isLoading = false;
  
  @override
  void initState() {
    super.initState();
    _loadWalletData();
  }
  
  Future<void> _loadWalletData() async {
    final user = await _authService.getCurrentUserWithWallet();
    if (user != null && user.walletAddress != null) {
      setState(() {
        _walletFuture = _walletService.getWalletInfo(user.walletAddress!);
      });
    }
  }
  
  Future<void> _refreshWalletData() async {
    setState(() {
      _loadWalletData();
    });
  }
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Wallet'),
        actions: [
          IconButton(
            icon: Icon(Icons.refresh),
            onPressed: _refreshWalletData,
          ),
        ],
      ),
      body: FutureBuilder<Wallet>(
        future: _walletFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text('Error: ${snapshot.error}'),
                  SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: _refreshWalletData,
                    child: Text('Retry'),
                  ),
                ],
              ),
            );
          } else if (snapshot.hasData) {
            final wallet = snapshot.data!;
            return RefreshIndicator(
              onRefresh: _refreshWalletData,
              child: SingleChildScrollView(
                physics: AlwaysScrollableScrollPhysics(),
                padding: EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _buildWalletCard(wallet),
                    SizedBox(height: 24),
                    _buildActionButtons(wallet),
                    SizedBox(height: 24),
                    Text(
                      'Recent Transactions',
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                    SizedBox(height: 8),
                    _buildTransactionList(wallet),
                  ],
                ),
              ),
            );
          } else {
            return Center(child: Text('No wallet data available'));
          }
        },
      ),
    );
  }
  
  Widget _buildWalletCard(Wallet wallet) {
    return Card(
      elevation: 4,
      child: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Wallet Address',
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
            SizedBox(height: 4),
            Text(
              wallet.address,
              style: TextStyle(fontSize: 12),
            ),
            Divider(),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'Balance',
                  style: TextStyle(fontWeight: FontWeight.bold),
                ),
                Text(
                  '${wallet.balance} TOKENS',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
  
  Widget _buildActionButtons(Wallet wallet) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        ElevatedButton.icon(
          icon: Icon(Icons.send),
          label: Text('Send'),
          onPressed: () => _showSendDialog(context),
        ),
        ElevatedButton.icon(
          icon: Icon(Icons.water_drop),
          label: Text('Request Tokens'),
          onPressed: () => _requestTestTokens(wallet.address),
        ),
      ],
    );
  }
  
  Widget _buildTransactionList(Wallet wallet) {
    if (wallet.recentTransactions.isEmpty) {
      return Padding(
        padding: EdgeInsets.symmetric(vertical: 16),
        child: Center(child: Text('No transactions yet')),
      );
    }
    
    return ListView.builder(
      shrinkWrap: true,
      physics: NeverScrollableScrollPhysics(),
      itemCount: wallet.recentTransactions.length,
      itemBuilder: (context, index) {
        final tx = wallet.recentTransactions[index];
        return Card(
          margin: EdgeInsets.symmetric(vertical: 4),
          child: ListTile(
            title: Text(
              tx.hash.length > 10 ? '${tx.hash.substring(0, 10)}...' : tx.hash,
            ),
            subtitle: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('${tx.amount} TOKENS'),
                Text(
                  'Date: ${tx.timestamp.toString().substring(0, 16)}',
                  style: TextStyle(fontSize: 12),
                ),
              ],
            ),
            trailing: _buildStatusChip(tx.status),
            leading: Icon(
              tx.from == wallet.address ? Icons.arrow_upward : Icons.arrow_downward,
              color: tx.from == wallet.address ? Colors.red : Colors.green,
            ),
          ),
        );
      },
    );
  }
  
  Widget _buildStatusChip(String status) {
    Color color;
    switch (status.toLowerCase()) {
      case 'completed':
        color = Colors.green;
        break;
      case 'pending':
        color = Colors.orange;
        break;
      case 'failed':
        color = Colors.red;
        break;
      default:
        color = Colors.grey;
    }
    
    return Chip(
      label: Text(
        status,
        style: TextStyle(color: Colors.white, fontSize: 12),
      ),
      backgroundColor: color,
      padding: EdgeInsets.zero,
      labelPadding: EdgeInsets.symmetric(horizontal: 8),
    );
  }
  
  void _showSendDialog(BuildContext context) {
    final addressController = TextEditingController();
    final amountController = TextEditingController();
    
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Send Tokens'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: addressController,
              decoration: InputDecoration(
                labelText: 'Recipient Address',
                hintText: '0x...',
              ),
            ),
            SizedBox(height: 16),
            TextField(
              controller: amountController,
              decoration: InputDecoration(
                labelText: 'Amount',
                hintText: '0.0',
              ),
              keyboardType: TextInputType.numberWithOptions(decimal: true),
            ),
          ],
        ),
        actions: [
          TextButton(
            child: Text('Cancel'),
            onPressed: () => Navigator.of(context).pop(),
          ),
          ElevatedButton(
            child: Text('Send'),
            onPressed: () {
              Navigator.of(context).pop();
              _sendTransaction(
                addressController.text,
                double.tryParse(amountController.text) ?? 0,
              );
            },
          ),
        ],
      ),
    );
  }
  
  Future<void> _requestTestTokens(String address) async {
    setState(() {
      _isLoading = true;
    });
    
    try {
      final result = await _walletService.requestTestTokens(address);
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(result)),
      );
      // Refresh wallet data
      _refreshWalletData();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: $e')),
      );
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }
  
  Future<void> _sendTransaction(String to, double amount) async {
    if (to.isEmpty || amount <= 0) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Invalid recipient or amount')),
      );
      return;
    }
    
    setState(() {
      _isLoading = true;
    });
    
    try {
      final txHash = await _walletService.sendTransaction(to, amount);
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Transaction sent: $txHash')),
      );
      // Refresh wallet data
      _refreshWalletData();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: $e')),
      );
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }
}