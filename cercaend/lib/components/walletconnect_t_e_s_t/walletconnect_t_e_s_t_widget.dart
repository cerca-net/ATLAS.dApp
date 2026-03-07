import '/flutter_flow/flutter_flow_theme.dart';
import '/flutter_flow/flutter_flow_util.dart';
import '/flutter_flow/flutter_flow_widgets.dart';
import '/services/blockchain/wallet_service.dart';
import '/services/blockchain/blockchain_service.dart';
import 'package:barcode_widget/barcode_widget.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import 'walletconnect_t_e_s_t_model.dart';
export 'walletconnect_t_e_s_t_model.dart';

class WalletconnectTESTWidget extends StatefulWidget {
  const WalletconnectTESTWidget({super.key});

  @override
  State<WalletconnectTESTWidget> createState() =>
      _WalletconnectTESTWidgetState();
}

class _WalletconnectTESTWidgetState extends State<WalletconnectTESTWidget> {
  late WalletconnectTESTModel _model;
  final WalletService _walletService = WalletService();
  final BlockchainService _blockchainService = BlockchainService();

  String? _mnemonic;
  bool _isConnecting = false;
  bool _hasSavedKeys = false;

  @override
  void setState(VoidCallback callback) {
    super.setState(callback);
    _model.onUpdate();
  }

  @override
  void initState() {
    super.initState();
    _model = createModel(context, () => WalletconnectTESTModel());
    _checkExistingWallet();
  }

  Future<void> _checkExistingWallet() async {
    final storedWallet = await _walletService.getStoredWallet();
    if (storedWallet != null) {
      final appState = FFAppState();
      appState.update(() {
        appState.isWalletConnected = true;
        appState.walletAddress = storedWallet['address'] ?? '';
      });
      _fetchBalance();
    }
  }

  Future<void> _fetchBalance() async {
    final appState = FFAppState();
    if (appState.walletAddress.isNotEmpty) {
      try {
        final info =
            await _blockchainService.getWalletInfo(appState.walletAddress);
        appState.update(() {
          appState.walletBalance = info.data.balance;
        });
      } catch (e) {
        debugPrint('Error fetching balance: $e');
      }
    }
  }

  Future<void> _connectWallet() async {
    setState(() => _isConnecting = true);
    try {
      final walletInfo = await _walletService.createWallet();
      if (mounted) {
        setState(() {
          _mnemonic = walletInfo['mnemonic'];
        });
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isConnecting = false);
      }
    }
  }

  void _onSavedKeys() {
    final appState = FFAppState();
    appState.update(() {
      appState.isWalletConnected = true;
    });
    setState(() {
      _hasSavedKeys = true;
      _mnemonic = null;
    });
    _fetchBalance();
  }

  @override
  void dispose() {
    _model.maybeDispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    context.watch<FFAppState>();
    final appState = FFAppState();
    final isWalletConnected = appState.isWalletConnected;

    return Align(
      alignment: const AlignmentDirectional(0.0, 0.0),
      child: Padding(
        padding: const EdgeInsets.all(4.0),
        child: Container(
          width: 320.0,
          height: 380.0,
          decoration: BoxDecoration(
            boxShadow: const [
              BoxShadow(
                blurRadius: 6.0,
                color: Color(0x4B1A1F24),
                offset: Offset(0.0, 2.0),
              )
            ],
            gradient: LinearGradient(
              colors: [
                FlutterFlowTheme.of(context).secondary,
                FlutterFlowTheme.of(context).secondaryText
              ],
              stops: const [0.0, 1.0],
              begin: const AlignmentDirectional(0.94, -1.0),
              end: const AlignmentDirectional(-0.94, 1.0),
            ),
            borderRadius: BorderRadius.circular(12.0),
          ),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              mainAxisSize: MainAxisSize.max,
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                // Top Header: Balance or Title
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text(
                      'ATLAS WALLET',
                      style: TextStyle(
                        fontFamily: 'Arial',
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                        fontSize: 14,
                      ),
                    ),
                    if (isWalletConnected)
                      Text(
                        '${appState.walletBalance.toStringAsFixed(2)} TCOIN',
                        style: const TextStyle(
                          fontFamily: 'Arial',
                          color: Colors.white,
                          fontWeight: FontWeight.bold,
                          fontSize: 16,
                        ),
                      ),
                  ],
                ),

                // Center Content: QR Code, Mnemonic, or Placeholder
                Expanded(
                  child: Center(
                    child: Builder(
                      builder: (context) {
                        if (_isConnecting) {
                          return const CircularProgressIndicator(
                              color: Colors.white);
                        }
                        if (_mnemonic != null) {
                          return Column(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              const Text(
                                'Save your 12-key phrase:',
                                style: TextStyle(
                                  fontFamily: 'Arial',
                                  color: Colors.white70,
                                  fontSize: 12,
                                ),
                              ),
                              const SizedBox(height: 8),
                              Container(
                                padding: const EdgeInsets.all(8),
                                decoration: BoxDecoration(
                                  color: Colors.black26,
                                  borderRadius: BorderRadius.circular(8),
                                ),
                                child: Text(
                                  _mnemonic!,
                                  textAlign: TextAlign.center,
                                  style: const TextStyle(
                                    fontFamily:
                                        'Courier New', // Safe system font for code
                                    color: Colors.white,
                                    fontSize: 14,
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                              ),
                            ],
                          );
                        }
                        if (isWalletConnected) {
                          return BarcodeWidget(
                            data: appState.walletAddress,
                            barcode: Barcode.qrCode(),
                            width: 160.0,
                            height: 160.0,
                            color: Colors.white,
                            backgroundColor: Colors.transparent,
                          );
                        }
                        return const Icon(
                          Icons.account_balance_wallet_rounded,
                          color: Colors.white54,
                          size: 80,
                        );
                      },
                    ),
                  ),
                ),

                // Bottom Actions
                if (_mnemonic != null)
                  FFButtonWidget(
                    onPressed: _onSavedKeys,
                    text: 'I have saved my keys',
                    options: FFButtonOptions(
                      width: double.infinity,
                      height: 40.0,
                      color: FlutterFlowTheme.of(context).primary,
                      textStyle: const TextStyle(
                        fontFamily: 'Arial',
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                      ),
                      borderRadius: BorderRadius.circular(8.0),
                    ),
                  )
                else if (!isWalletConnected)
                  FFButtonWidget(
                    onPressed: _connectWallet,
                    text: 'Connect Wallet',
                    icon: const Icon(Icons.account_balance_wallet, size: 15.0),
                    options: FFButtonOptions(
                      width: double.infinity,
                      height: 40.0,
                      color: FlutterFlowTheme.of(context).primary,
                      textStyle: const TextStyle(
                        fontFamily: 'Arial',
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                      ),
                      borderRadius: BorderRadius.circular(8.0),
                    ),
                  )
                else
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      _buildActionBtn(
                          context, Icons.send_rounded, 'Send', _onSendTap),
                      _buildActionBtn(context, Icons.call_received_rounded,
                          'Receive', _onReceiveTap),
                      _buildActionBtn(context, Icons.history_rounded, 'History',
                          _onHistoryTap),
                      _buildActionBtn(context, Icons.water_drop_rounded,
                          'Faucet', _onFaucetTap),
                    ],
                  ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildActionBtn(
      BuildContext context, IconData icon, String label, VoidCallback onTap) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        IconButton(
          icon: Icon(icon, color: Colors.white, size: 24),
          onPressed: onTap,
        ),
        Text(
          label,
          style: const TextStyle(
            fontFamily: 'Arial',
            color: Colors.white,
            fontSize: 10,
          ),
        ),
      ],
    );
  }

  void _onReceiveTap() {
    final address = FFAppState().walletAddress;
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Receive Tokens'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            BarcodeWidget(
              data: address,
              barcode: Barcode.qrCode(),
              width: 200,
              height: 200,
            ),
            const SizedBox(height: 16),
            SelectableText(
              address,
              style: const TextStyle(fontSize: 12),
              textAlign: TextAlign.center,
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () {
              Clipboard.setData(ClipboardData(text: address));
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('Address copied!')),
              );
            },
            child: const Text('Copy'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Close'),
          ),
        ],
      ),
    );
  }

  void _onHistoryTap() {
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            const Text('Transaction History',
                style: TextStyle(fontWeight: FontWeight.bold, fontSize: 18)),
            const SizedBox(height: 10),
            Expanded(
              child: FutureBuilder<TransactionHistoryResponse>(
                future: _blockchainService
                    .getTransactionHistory(FFAppState().walletAddress),
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return const Center(child: CircularProgressIndicator());
                  }
                  if (snapshot.hasError) {
                    return Center(child: Text('Error: ${snapshot.error}'));
                  }
                  final txs = snapshot.data?.data.transactions ?? [];
                  if (txs.isEmpty) {
                    return const Center(child: Text('No transactions found.'));
                  }
                  return ListView.builder(
                    itemCount: txs.length,
                    itemBuilder: (context, index) {
                      final tx = txs[index];
                      final isReceived = tx.recipient.toLowerCase() ==
                          FFAppState().walletAddress.toLowerCase();
                      return ListTile(
                        leading: Icon(
                          isReceived
                              ? Icons.arrow_downward
                              : Icons.arrow_upward,
                          color: isReceived ? Colors.green : Colors.red,
                        ),
                        title: Text('${tx.amount} Tokens'),
                        subtitle: Text(isReceived
                            ? 'From: ${tx.sender}'
                            : 'To: ${tx.recipient}'),
                        trailing: Text(DateTime.fromMillisecondsSinceEpoch(
                                tx.timestamp * 1000)
                            .toString()
                            .split(' ')[0]),
                      );
                    },
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _onFaucetTap() async {
    final address = FFAppState().walletAddress;
    if (address.isEmpty) return;

    setState(() => _isConnecting = true); // Show loading
    try {
      await _blockchainService.requestFaucet(address);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Faucet tokens received!')),
        );
        _fetchBalance();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
              content: Text('Faucet failed: $e'), backgroundColor: Colors.red),
        );
      }
    } finally {
      if (mounted) setState(() => _isConnecting = false);
    }
  }

  void _onSendTap() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        height: MediaQuery.of(context).size.height * 0.85,
        decoration: const BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
        ),
        padding: const EdgeInsets.all(20),
        child: SendTransactionForm(
          onSend: (to, amount, msg) => _processSend(to, amount, msg),
        ),
      ),
    );
  }

  Future<void> _processSend(String to, double amount, String msg) async {
    Navigator.pop(context); // Close modal
    setState(() => _isConnecting = true); // Use loading state
    try {
      final sender = await _walletService.getAddress();
      final senderPubKey = await _walletService.getPublicKey();
      if (sender == null || senderPubKey == null) {
        throw Exception('Wallet not loaded');
      }

      final nonce = await _blockchainService.getNonce(sender);
      final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;

      // Construct TX to sign matches 'SendTransactionRequest' structure
      final txData = {
        'Sender': sender,
        'Recipient': to,
        'Amount': amount.toInt(),
        'Fee': 0,
        'Timestamp': timestamp,
        'Nonce': nonce,
        'Data': msg,
      };

      final signature = await _walletService.signTransaction(txData);

      final request = SendTransactionRequest(
        sender: sender,
        senderPublicKey: senderPubKey,
        recipient: to,
        amount: amount.toInt(),
        fee: 0,
        timestamp: timestamp,
        nonce: nonce,
        data: msg,
        signature: signature,
      );

      await _blockchainService.sendTransaction(request);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Transaction Sent Successfully!')),
        );
        _fetchBalance(); // Refresh balance
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
              content: Text('Failed to send: $e'), backgroundColor: Colors.red),
        );
      }
    } finally {
      if (mounted) setState(() => _isConnecting = false);
    }
  }
}

class SendTransactionForm extends StatefulWidget {
  final Function(String, double, String) onSend;
  const SendTransactionForm({super.key, required this.onSend});

  @override
  State<SendTransactionForm> createState() => _SendTransactionFormState();
}

class _SendTransactionFormState extends State<SendTransactionForm> {
  final _toController = TextEditingController();
  final _amountController = TextEditingController();
  final _msgController = TextEditingController();
  final _formKey = GlobalKey<FormState>();

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: Column(
        children: [
          const Text('Send Transaction',
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
          const SizedBox(height: 20),
          TextFormField(
            controller: _toController,
            decoration: const InputDecoration(
                labelText: 'Recipient Address', border: OutlineInputBorder()),
            validator: (v) => v!.isEmpty ? 'Required' : null,
          ),
          const SizedBox(height: 10),
          TextFormField(
            controller: _amountController,
            decoration: const InputDecoration(
                labelText: 'Amount', border: OutlineInputBorder()),
            keyboardType: TextInputType.number,
            validator: (v) => v!.isEmpty || double.tryParse(v) == null
                ? 'Invalid amount'
                : null,
          ),
          const SizedBox(height: 10),
          TextFormField(
            controller: _msgController,
            decoration: const InputDecoration(
                labelText: 'Message (Optional)', border: OutlineInputBorder()),
          ),
          const SizedBox(height: 20),
          ElevatedButton(
            onPressed: () {
              if (_formKey.currentState!.validate()) {
                widget.onSend(_toController.text,
                    double.parse(_amountController.text), _msgController.text);
              }
            },
            style: ElevatedButton.styleFrom(
              minimumSize: const Size(double.infinity, 50),
              backgroundColor: Colors.blueAccent,
              foregroundColor: Colors.white,
            ),
            child: const Text('Confirm Send'),
          ),
        ],
      ),
    );
  }
}
