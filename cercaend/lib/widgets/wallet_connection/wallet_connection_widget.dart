import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '/services/blockchain/wallet_service.dart';
import '/app_state.dart';
import '/auth/auth_util.dart';

class WalletConnectionWidget extends StatefulWidget {
  const WalletConnectionWidget({super.key});

  @override
  State<WalletConnectionWidget> createState() => _WalletConnectionWidgetState();
}

class _WalletConnectionWidgetState extends State<WalletConnectionWidget> {
  final WalletService _walletService = WalletService();

  bool _isLoading = false;
  bool _isConnected = false;
  String _walletAddress = '';
  String _mnemonic = '';
  String _errorMessage = '';
  bool _needsImport = false;
  final TextEditingController _mnemonicController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _checkExistingWallet();
  }

  Future<void> _checkExistingWallet() async {
    try {
      final storedWallet = await _walletService.getStoredWallet();
      if (storedWallet != null) {
        if (mounted) {
          setState(() {
            _walletAddress = storedWallet['address'] ?? '';
            _isConnected = true;
          });

          final appState = Provider.of<FFAppState>(context, listen: false);
          appState.isWalletConnected = true;
          appState.walletAddress = _walletAddress;
        }
      } else {
        if (currentUserDocument != null &&
            currentUserDocument!.walletAddress.isNotEmpty) {
          if (mounted) {
            setState(() {
              _needsImport = true;
            });
          }
        }
      }
    } catch (e) {
      debugPrint('Error checking existing wallet: $e');
    }
  }

  Future<void> _connectWallet() async {
    if (!mounted) return;

    setState(() {
      _isLoading = true;
      _errorMessage = '';
    });

    try {
      final walletInfo = await _walletService.createWallet();

      if (mounted) {
        setState(() {
          _walletAddress = walletInfo['address']!;
          _mnemonic = walletInfo['mnemonic']!;
          _isConnected = true;
          _isLoading = false;
        });

        final appState = Provider.of<FFAppState>(context, listen: false);
        appState.isWalletConnected = true;
        appState.walletAddress = _walletAddress;
        appState.sessionToken = '';
        appState.walletBalance = 0.0;

        try {
          if (currentUserReference != null) {
            await currentUserReference!.update({'wallet_address': _walletAddress});
          }
        } catch (e) {
          debugPrint('Failed to save wallet address to database: $e');
        }

        _showMnemonicDialog(_mnemonic);
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _errorMessage = 'Failed to connect wallet: ${e.toString()}';
          _isLoading = false;
        });
      }
    }
  }

  Future<void> _importWallet() async {
    if (_mnemonicController.text.trim().isEmpty) {
      setState(() => _errorMessage = 'Please enter mnemonic phrase');
      return;
    }

    if (!mounted) return;

    setState(() {
      _isLoading = true;
      _errorMessage = '';
    });

    try {
      final walletInfo =
          await _walletService.importWallet(_mnemonicController.text.trim());

      if (mounted) {
        setState(() {
          _walletAddress = walletInfo['address']!;
          _isConnected = true;
          _needsImport = false;
          _isLoading = false;
        });

        final appState = Provider.of<FFAppState>(context, listen: false);
        appState.isWalletConnected = true;
        appState.walletAddress = _walletAddress;
        appState.sessionToken = '';
        appState.walletBalance = 0.0;

        try {
          if (currentUserReference != null) {
            await currentUserReference!.update({'wallet_address': _walletAddress});
          }
        } catch (e) {
          debugPrint('Failed to save imported wallet address to database: $e');
        }
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _errorMessage = 'Failed to import wallet: ${e.toString()}';
          _isLoading = false;
        });
      }
    }
  }

  Future<void> _disconnectWallet() async {
    if (!mounted) return;

    setState(() {
      _isLoading = true;
    });

    try {
      await _walletService.logout();

      if (mounted) {
        setState(() {
          _isConnected = false;
          _walletAddress = '';
          _mnemonic = '';
          _isLoading = false;
        });

        final appState = Provider.of<FFAppState>(context, listen: false);
        appState.isWalletConnected = false;
        appState.walletAddress = '';
        appState.sessionToken = '';
        appState.walletBalance = 0.0;
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _errorMessage = 'Failed to disconnect wallet: ${e.toString()}';
          _isLoading = false;
        });
      }
    }
  }

  void _showMnemonicDialog(String mnemonic) {
    if (!mounted) return;

    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => AlertDialog(
        title: const Text('Wallet Created Successfully!'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Save this mnemonic phrase securely. You will need it to recover your wallet:',
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 16),
            SelectableText(
              mnemonic,
              style: const TextStyle(
                fontFamily: 'monospace',
                fontSize: 16,
                letterSpacing: 2,
              ),
            ),
            const SizedBox(height: 16),
            const Text(
              '⚠️ WARNING: Never share this mnemonic with anyone!',
              style: TextStyle(color: Colors.red, fontWeight: FontWeight.bold),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () {
              if (mounted) {
                Navigator.of(context).pop();
              }
            },
            child: const Text('I Understand'),
          ),
        ],
      ),
    );
  }

  void _showWalletInfo() {
    if (!mounted) return;

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Wallet Information'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('Address:'),
            SelectableText(
              _walletAddress,
              style: const TextStyle(
                fontFamily: 'monospace',
                fontSize: 12,
              ),
            ),
            const SizedBox(height: 16),
            const Text('Connection Status:'),
            Text(
              _isConnected ? 'Connected' : 'Disconnected',
              style: TextStyle(
                color: _isConnected ? Colors.green : Colors.red,
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () {
              if (mounted) {
                Navigator.of(context).pop();
              }
            },
            child: const Text('Close'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 4,
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Row(
              children: [
                Icon(
                  Icons.account_balance_wallet,
                  color: Colors.blue,
                  size: 24,
                ),
                SizedBox(width: 8),
                Text(
                  'Wallet Connection',
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            if (_errorMessage.isNotEmpty)
              Padding(
                padding: const EdgeInsets.only(bottom: 16.0),
                child: Text(
                  _errorMessage,
                  style: const TextStyle(
                    color: Colors.red,
                    fontSize: 14,
                  ),
                ),
              ),
            if (_isLoading)
              const Padding(
                padding: EdgeInsets.all(16.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    CircularProgressIndicator(),
                    SizedBox(width: 16),
                    Text('Connecting...'),
                  ],
                ),
              )
            else if (_isConnected)
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Row(
                    children: [
                      Icon(
                        Icons.check_circle,
                        color: Colors.green,
                        size: 20,
                      ),
                      SizedBox(width: 8),
                      Text(
                        'Wallet Connected',
                        style: TextStyle(
                          color: Colors.green,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Text(
                    'Address: ${_walletAddress.substring(0, 10)}...',
                    style: const TextStyle(fontSize: 14),
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        child: ElevatedButton.icon(
                          icon: const Icon(Icons.info),
                          label: const Text('View Details'),
                          onPressed: _showWalletInfo,
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.blue,
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: ElevatedButton.icon(
                          icon: const Icon(Icons.logout),
                          label: const Text('Disconnect'),
                          onPressed: _disconnectWallet,
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.red,
                          ),
                        ),
                      ),
                    ],
                  ),
                ],
              )
            else if (_needsImport)
              Column(
                children: [
                  const Text(
                    'A wallet is linked to your profile, but it\'s missing from this device. Please import it using your 12-word mnemonic phrase.',
                    style: TextStyle(fontSize: 14),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 16),
                  TextField(
                    controller: _mnemonicController,
                    decoration: const InputDecoration(
                      labelText: 'Enter Mnemonic Phrase',
                      border: OutlineInputBorder(),
                    ),
                    maxLines: 3,
                  ),
                  const SizedBox(height: 16),
                  ElevatedButton.icon(
                    icon: const Icon(Icons.download),
                    label: const Text('Import Wallet'),
                    onPressed: _importWallet,
                    style: ElevatedButton.styleFrom(
                      minimumSize: const Size(double.infinity, 44),
                      backgroundColor: Colors.orange,
                    ),
                  ),
                ],
              )
            else
              Column(
                children: [
                  const Text(
                    'Connect your wallet to access blockchain features',
                    style: TextStyle(fontSize: 14),
                  ),
                  const SizedBox(height: 16),
                  ElevatedButton.icon(
                    icon: const Icon(Icons.add),
                    label: const Text('Connect Wallet'),
                    onPressed: _connectWallet,
                    style: ElevatedButton.styleFrom(
                      minimumSize: const Size(double.infinity, 44),
                    ),
                  ),
                ],
              ),
          ],
        ),
      ),
    );
  }
}
