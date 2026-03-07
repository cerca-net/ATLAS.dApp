import 'package:get_it/get_it.dart';
import 'package:firebase_auth/firebase_auth.dart';
import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:web3dart/web3dart.dart';
import 'package:http/http.dart' as http;

import '../services/api_client.dart';
import '../services/auth_service.dart';
import '../services/secure_storage.dart';
import '../services/wallet_service.dart';

final GetIt serviceLocator = GetIt.instance;

Future<void> setupServiceLocator() async {
  // External services
  serviceLocator.registerSingleton<FirebaseAuth>(FirebaseAuth.instance);
  serviceLocator.registerSingleton<FirebaseFirestore>(FirebaseFirestore.instance);
  
  // Create HTTP client
  final httpClient = http.Client();
  serviceLocator.registerSingleton<http.Client>(httpClient);
  
  // Create Web3 client
  final web3Client = Web3Client('https://rpc.atlas.bc.com', httpClient);
  serviceLocator.registerSingleton<Web3Client>(web3Client);
  
  // App services
  serviceLocator.registerSingleton<SecureStorage>(SecureStorage());
  
  serviceLocator.registerSingleton<ApiClient>(
    ApiClient(
      baseUrl: 'https://api.atlas.bc.com',
      secureStorage: serviceLocator<SecureStorage>(),
    ),
  );
  
  serviceLocator.registerSingleton<AuthService>(
    AuthService(
      serviceLocator<FirebaseAuth>(),
      serviceLocator<FirebaseFirestore>(),
      serviceLocator<ApiClient>(),
    ),
  );
  
  serviceLocator.registerSingleton<WalletService>(
    WalletService(
      serviceLocator<ApiClient>(),
      serviceLocator<SecureStorage>(),
      serviceLocator<Web3Client>(),
    ),
  );
}