import 'dart:math';
import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:hex/hex.dart';
import 'package:elliptic/elliptic.dart' as elliptic;
import 'package:ecdsa/ecdsa.dart' as ecdsa;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'blockchain_service.dart';
import '/auth/firebase_auth/auth_util.dart';
import 'bip39_words.dart';

class WalletService {
  static final WalletService _instance = WalletService._internal();
  factory WalletService() => _instance;
  WalletService._internal();

  final _storage = const FlutterSecureStorage();
  final _blockchainService = BlockchainService();

  static const String _addressKey = 'wallet_address';
  static const String _mnemonicKey = 'wallet_mnemonic';

  /// Create a new wallet locally (Local Identity) and register address with node
  Future<Map<String, String>> createWallet() async {
    try {
      print('WalletService: Starting wallet creation...');

      // 1. Generate Mnemonic
      final random = Random.secure();
      final words = <String>[];
      for (int i = 0; i < 12; i++) {
        final index = random.nextInt(bip39EnglishWords.length);
        words.add(bip39EnglishWords[index]);
      }
      final mnemonic = words.join(' ');
      print('WalletService: Mnemonic generated');

      // 2. Derive Keys locally using pure Dart (Safe for Flutter Web)
      final mnemonicBytes = utf8.encode(mnemonic);
      final seedHash = sha256.convert(mnemonicBytes).bytes;
      final privateKeyBytes =
          BigInt.parse(HEX.encode(seedHash.sublist(0, 32)), radix: 16);

      final curve = elliptic.getP256();
      final privateKey = elliptic.PrivateKey(curve, privateKeyBytes);
      final publicKey = privateKey.publicKey;
      print('WalletService: Deterministic KeyPair generated via pure Dart');

      // 3. Construct ASN.1 DER Encoded Public Key
      final List<int> header = [
        0x30,
        0x59,
        0x30,
        0x13,
        0x06,
        0x07,
        0x2A,
        0x86,
        0x48,
        0xCE,
        0x3D,
        0x02,
        0x01,
        0x06,
        0x08,
        0x2A,
        0x86,
        0x48,
        0xCE,
        0x3D,
        0x03,
        0x01,
        0x07,
        0x03,
        0x42,
        0x00
      ];

      // Pad X and Y to 32 bytes if they are smaller
      final xBytes = _padTo32(publicKey.X.toRadixString(16));
      final yBytes = _padTo32(publicKey.Y.toRadixString(16));

      final List<int> uncompressedPoint = [
        0x04,
        ...HEX.decode(xBytes),
        ...HEX.decode(yBytes)
      ];
      final List<int> derPublicKey = [...header, ...uncompressedPoint];
      print('WalletService: DER public key constructed');

      // 4. Derive Address: SHA256(DER_PubKey) -> Last 20 Bytes -> Hex
      final pubKeyHash = sha256.convert(derPublicKey).bytes;
      final addressBytes = pubKeyHash.sublist(pubKeyHash.length - 20);
      final address = '0x${HEX.encode(addressBytes).toLowerCase()}';
      print('WalletService: Address derived: $address');

      // 5. Connect to Backend
      print('WalletService: Connecting to backend at http://localhost:8080...');
      bool registrationSuccess = false;
      try {
        final response = await _blockchainService.connectWallet(
          WalletConnectRequest(
            action: 'connect',
            address: address,
          ),
        );
        registrationSuccess = response.success;
      } catch (e) {
        print('WalletService: Backend offline or network error: $e');
        // We still proceed so the user can see their local keys
      }

      // 6. Store Securely
      print('WalletService: Storing keys locally...');
      await _storage.write(key: _addressKey, value: address);
      await _storage.write(key: _mnemonicKey, value: mnemonic);
      await _storage.write(key: 'wallet_privkey', value: privateKey.toHex());
      print('WalletService: Local storage complete');

      // 7. Update user document in Firestore
      if (currentUserReference != null) {
        try {
          await currentUserReference!.update({
            'wallet_address': address,
          });
          print('WalletService: Firestore updated with $address');
        } catch (e) {
          print('WalletService: Firestore update failed: $e');
        }
      }

      return {
        'address': address,
        'mnemonic': mnemonic,
        'registered': registrationSuccess.toString(),
      };
    } catch (e, stack) {
      print('WalletService: FATAL ERROR: $e');
      print('Stack Trace: $stack');
      throw Exception('Error creating wallet: $e');
    }
  }

  String _padTo32(String hex) {
    var h = hex;
    if (h.length % 2 != 0) h = '0$h';
    while (h.length < 64) {
      h = '00$h';
    }
    return h;
  }

  /// Import an existing wallet from mnemonic
  Future<Map<String, String>> importWallet(String mnemonic) async {
    try {
      print('WalletService: Starting wallet import...');

      // Derive Keys locally using pure Dart (Safe for Flutter Web)
      final mnemonicBytes = utf8.encode(mnemonic.trim());
      final seedHash = sha256.convert(mnemonicBytes).bytes;
      final privateKeyBytes =
          BigInt.parse(HEX.encode(seedHash.sublist(0, 32)), radix: 16);

      final curve = elliptic.getP256();
      final privateKey = elliptic.PrivateKey(curve, privateKeyBytes);
      final publicKey = privateKey.publicKey;

      // Construct ASN.1 DER Encoded Public Key
      final List<int> header = [
        0x30,
        0x59,
        0x30,
        0x13,
        0x06,
        0x07,
        0x2A,
        0x86,
        0x48,
        0xCE,
        0x3D,
        0x02,
        0x01,
        0x06,
        0x08,
        0x2A,
        0x86,
        0x48,
        0xCE,
        0x3D,
        0x03,
        0x01,
        0x07,
        0x03,
        0x42,
        0x00
      ];

      final xBytes = _padTo32(publicKey.X.toRadixString(16));
      final yBytes = _padTo32(publicKey.Y.toRadixString(16));
      final List<int> uncompressedPoint = [
        0x04,
        ...HEX.decode(xBytes),
        ...HEX.decode(yBytes)
      ];
      final List<int> derPublicKey = [...header, ...uncompressedPoint];

      // Derive Address
      final pubKeyHash = sha256.convert(derPublicKey).bytes;
      final addressBytes = pubKeyHash.sublist(pubKeyHash.length - 20);
      final address = '0x${HEX.encode(addressBytes).toLowerCase()}';

      print('WalletService: Address imported: $address');

      bool registrationSuccess = false;
      try {
        final response = await _blockchainService.connectWallet(
          WalletConnectRequest(
            action: 'connect',
            address: address,
          ),
        );
        registrationSuccess = response.success;
      } catch (e) {
        print('WalletService: Backend offline or network error: $e');
      }

      await _storage.write(key: _addressKey, value: address);
      await _storage.write(key: _mnemonicKey, value: mnemonic.trim());
      await _storage.write(key: 'wallet_privkey', value: privateKey.toHex());

      if (currentUserReference != null) {
        await currentUserReference!.update({
          'wallet_address': address,
        });
      }

      return {
        'address': address,
        'mnemonic': mnemonic.trim(),
        'registered': registrationSuccess.toString(),
      };
    } catch (e) {
      throw Exception('Error importing wallet: $e');
    }
  }

  Future<String?> getAddress() async => await _storage.read(key: _addressKey);
  Future<String?> getMnemonic() async => await _storage.read(key: _mnemonicKey);

  Future<String?> getPublicKey() async {
    final privateKeyHex = await _storage.read(key: 'wallet_privkey');
    if (privateKeyHex == null) return null;

    final curve = elliptic.getP256();
    final privateKey = elliptic.PrivateKey.fromHex(curve, privateKeyHex);
    final publicKey = privateKey.publicKey;

    // Construct ASN.1 DER Encoded Public Key (SPKI) to match HTML's exportKey('spki')
    // This allows the backend to parse it correctly.
    final List<int> header = [
      0x30,
      0x59,
      0x30,
      0x13,
      0x06,
      0x07,
      0x2A,
      0x86,
      0x48,
      0xCE,
      0x3D,
      0x02,
      0x01,
      0x06,
      0x08,
      0x2A,
      0x86,
      0x48,
      0xCE,
      0x3D,
      0x03,
      0x01,
      0x07,
      0x03,
      0x42,
      0x00
    ];

    final xBytes = _padTo32(publicKey.X.toRadixString(16));
    final yBytes = _padTo32(publicKey.Y.toRadixString(16));

    final List<int> uncompressedPoint = [
      0x04,
      ...HEX.decode(xBytes),
      ...HEX.decode(yBytes)
    ];
    final List<int> derPublicKey = [...header, ...uncompressedPoint];

    return HEX.encode(derPublicKey);
  }

  /// Returns the stored wallet details if they exist
  Future<Map<String, String>?> getStoredWallet() async {
    final address = await getAddress();
    final mnemonic = await getMnemonic();
    if (address != null && mnemonic != null) {
      return {
        'address': address,
        'mnemonic': mnemonic,
      };
    }
    return null;
  }

  Future<void> logout() async {
    await _storage.delete(key: _addressKey);
    await _storage.delete(key: _mnemonicKey);
    await _storage.delete(key: 'wallet_privkey');
  }

  Future<String> signTransaction(Map<String, dynamic> tx) async {
    try {
      final privateKeyHex = await _storage.read(key: 'wallet_privkey');
      if (privateKeyHex == null) throw Exception('No private key found');

      // Reconstruct PrivateKey object
      final curve = elliptic.getP256();
      final privateKey = elliptic.PrivateKey.fromHex(curve, privateKeyHex);

      // Construct data string for signing: Sender + Recipient + Amount + Fee + Timestamp + Nonce + Data
      // IMPORTANT: This MUST match the Go backend's expectation exactly.
      final String data =
          '${tx['Sender']}${tx['Recipient']}${tx['Amount']}${tx['Fee']}${tx['Timestamp']}${tx['Nonce']}${tx['Data'] ?? ''}';

      // Hash the data
      final bytes = utf8.encode(data);
      // Note: The elliptic package's `signature` method expects the hash of the message, not the message itself,
      // OR it handles hashing internally depending on usage.
      // Standard ECDSA signs the hash. The 'crypto' package sha256 returns a Digest.
      final hash = sha256.convert(bytes).bytes;

      // Sign
      final signature = ecdsa.signature(privateKey, hash);

      // Return as Hex string (R + S concatenated) - Backend expects 64-byte raw signature (r||s)
      return _padTo32(signature.R.toRadixString(16)) +
          _padTo32(signature.S.toRadixString(16));
    } catch (e) {
      throw Exception('Failed to sign transaction: $e');
    }
  }

  Future<Map<String, dynamic>> sendTransaction({
    required String recipient,
    required double amount,
    required String sender,
    String? data,
    String type = 'regular',
  }) async {
    try {
      // 1. Get current nonce for sender to prevent replay
      final nonce = await _blockchainService.getNonce(sender);
      final publicKey = await getPublicKey();
      if (publicKey == null) throw Exception('Local keys not found');

      final timestamp = DateTime.now().millisecondsSinceEpoch;
      const fee = 1;
      final amountInt = (amount * 10).toInt(); // Convert to integer units

      // 2. Prepare for signing
      final txMap = {
        'Sender': sender,
        'Recipient': recipient,
        'Amount': amountInt,
        'Fee': fee,
        'Timestamp': timestamp,
        'Nonce': nonce,
        'Data': data ?? '',
        'Type': type,
      };

      // 3. Sign
      final signature = await signTransaction(txMap);

      // 4. Construct Request
      final request = SendTransactionRequest(
        sender: sender,
        senderPublicKey: publicKey,
        recipient: recipient,
        amount: amountInt,
        fee: fee,
        timestamp: timestamp,
        nonce: nonce,
        signature: signature,
        data: data,
        type: type,
      );

      // 5. Submit to Backend
      final response = await _blockchainService.sendTransaction(request);

      return {
        'success': response.success,
        'message': response.message,
        'hash': response.data.transactionHash,
      };
    } catch (e) {
      rethrow;
    }
  }

  /// Register as a validator (Stake tokens)
  Future<Map<String, dynamic>> registerAsValidator({
    required int stake,
    required Map<String, dynamic> kyc,
  }) async {
    try {
      final address = await getAddress();
      if (address == null) throw Exception('Wallet not found');

      final publicKey = await getPublicKey();
      if (publicKey == null) throw Exception('Local keys not found');

      // 1. Get Nonce
      final nonce = await _blockchainService.getNonce(address);

      final timestamp = DateTime.now().millisecondsSinceEpoch;
      const fee = 1; // Standard fee

      final kycJson = jsonEncode(kyc);

      // 2. Prepare Tx for Signing
      // Use self-address as recipient for staking transaction compliance
      final txMap = {
        'Sender': address,
        'Recipient': address,
        'Amount': stake,
        'Fee': fee,
        'Timestamp': timestamp,
        'Nonce': nonce,
        'Data': kycJson,
      };

      // 3. Sign
      final signature = await signTransaction(txMap);

      // 4. Construct Request
      final request = SendTransactionRequest(
        sender: address,
        senderPublicKey: publicKey,
        recipient: address,
        amount: stake,
        fee: fee,
        timestamp: timestamp,
        nonce: nonce,
        data: kycJson,
        signature: signature,
        type: 'stake',
      );

      // 5. Submit
      return await _blockchainService.registerValidator(request);
    } catch (e) {
      rethrow;
    }
  }
}
