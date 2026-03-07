import 'package:web3dart/web3dart.dart';
import 'package:wallet_integration/models/wallet.dart';
import 'package:wallet_integration/models/transaction.dart';
import 'api_client.dart';
import 'secure_storage.dart';

class WalletException implements Exception {
  final String message;
  final dynamic cause;
  
  WalletException(this.message, [this.cause]);
  
  @override
  String toString() => 'WalletException: $message${cause != null ? ' ($cause)' : ''}';
}

class WalletService {
  final ApiClient _apiClient;
  final SecureStorage _secureStorage;
  final Web3Client _web3Client;
  
  WalletService(
    this._apiClient,
    this._secureStorage,
    this._web3Client,
  );
  
  // Get wallet information
  Future<Wallet> getWalletInfo(String address) async {
    try {
      return await _apiClient.getWalletInfo(address);
    } on ApiException catch (e) {
      throw WalletException('Failed to get wallet info: ${e.message}', e);
    } catch (e) {
      throw WalletException('Failed to get wallet info', e);
    }
  }
  
  // Request test tokens
  Future<String> requestTestTokens(String address) async {
    try {
      return await _apiClient.requestTestTokens(address);
    } on ApiException catch (e) {
      throw WalletException('Failed to request test tokens: ${e.message}', e);
    } catch (e) {
      throw WalletException('Failed to request test tokens', e);
    }
  }
  
  // Send transaction
  Future<String> sendTransaction(String to, double amount) async {
    try {
      // Get private key from secure storage
      final privateKey = await _secureStorage.getPrivateKey();
      if (privateKey == null) {
        throw WalletException('Private key not found');
      }
      
      // Create credentials from private key
      final credentials = EthPrivateKey.fromHex(privateKey);
      
      // Get sender address
      final address = await credentials.extractAddress();
      
      // Convert amount to Wei
      final amountInWei = EtherAmount.fromUnitAndValue(
        EtherUnit.ether,
        BigInt.from(amount),
      );
      
      // Create transaction
      final transaction = Transaction(
        hash: '',
        from: address.hex,
        to: to,
        amount: amount,
        timestamp: DateTime.now().millisecondsSinceEpoch,
        status: 'pending',
      );
      
      // Sign transaction
      final signedTx = await _web3Client.signTransaction(
        credentials,
        Transaction(
          from: address.hex,
          to: to,
          value: amountInWei,
        ),
        chainId: 1337, // Replace with actual chain ID
      );
      
      // Send transaction
      final txHash = await _apiClient.sendTransaction(
        signedTx.toString(),
        to,
      );
      
      return txHash;
    } on ApiException catch (e) {
      throw WalletException('Failed to send transaction: ${e.message}', e);
    } catch (e) {
      throw WalletException('Failed to send transaction', e);
    }
  }
  
  // Import wallet
  Future<WalletResponse> importWallet(String privateKey) async {
    try {
      final response = await _apiClient.importWallet(privateKey);
      
      // Store private key and session token
      await _secureStorage.storePrivateKey(response.privateKey);
      await _secureStorage.storeSessionToken(response.sessionToken);
      
      return response;
    } on ApiException catch (e) {
      throw WalletException('Failed to import wallet: ${e.message}', e);
    } catch (e) {
      throw WalletException('Failed to import wallet', e);
    }
  }
}