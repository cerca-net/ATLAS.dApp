import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class SecureStorage {
  final _storage = const FlutterSecureStorage();
  
  // Key constants
  static const String _privateKeyKey = 'private_key';
  static const String _sessionTokenKey = 'session_token';
  
  // Store private key securely
  Future<void> storePrivateKey(String privateKey) async {
    await _storage.write(key: _privateKeyKey, value: privateKey);
  }
  
  // Retrieve private key
  Future<String?> getPrivateKey() async {
    return await _storage.read(key: _privateKeyKey);
  }
  
  // Store session token
  Future<void> storeSessionToken(String token) async {
    await _storage.write(key: _sessionTokenKey, value: token);
  }
  
  // Retrieve session token
  Future<String?> getSessionToken() async {
    return await _storage.read(key: _sessionTokenKey);
  }
  
  // Clear all stored data (for logout)
  Future<void> clearStorage() async {
    await _storage.deleteAll();
  }
}