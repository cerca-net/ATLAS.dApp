import 'package:dio/dio.dart';
import 'package:wallet_integration/models/wallet.dart';
import 'package:wallet_integration/models/transaction.dart';
import 'secure_storage.dart';

class ApiException implements Exception {
  final String message;
  final int? statusCode;
  
  ApiException(this.message, {this.statusCode});
  
  @override
  String toString() => 'ApiException: $message (Status code: $statusCode)';
}

class AuthInterceptor extends Interceptor {
  final SecureStorage secureStorage;
  
  AuthInterceptor(this.secureStorage);
  
  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final token = await secureStorage.getSessionToken();
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    return handler.next(options);
  }
}

class ApiClient {
  final Dio _dio;
  final String baseUrl;
  final SecureStorage secureStorage;
  
  ApiClient({
    required this.baseUrl,
    required this.secureStorage,
  }) : _dio = Dio(BaseOptions(
          baseUrl: baseUrl,
          connectTimeout: const Duration(seconds: 10),
          receiveTimeout: const Duration(seconds: 10),
          headers: {'Content-Type': 'application/json'},
        )) {
    _dio.interceptors.add(AuthInterceptor(secureStorage));
  }
  
  // Create a new wallet
  Future<WalletResponse> createWallet() async {
    try {
      final response = await _dio.post('/wallet/create');
      return WalletResponse.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }
  
  // Import an existing wallet
  Future<WalletResponse> importWallet(String privateKey) async {
    try {
      final response = await _dio.post(
        '/wallet/import',
        data: {'privateKey': privateKey},
      );
      return WalletResponse.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }
  
  // Get wallet information
  Future<Wallet> getWalletInfo(String address) async {
    try {
      final response = await _dio.get('/wallet/$address');
      return Wallet.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }
  
  // Request test tokens
  Future<String> requestTestTokens(String address) async {
    try {
      final response = await _dio.post(
        '/faucet/request',
        data: {'address': address},
      );
      return response.data['message'];
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }
  
  // Send transaction
  Future<String> sendTransaction(String signedTx, String to) async {
    try {
      final response = await _dio.post(
        '/transaction/send',
        data: {
          'signedTransaction': signedTx,
          'to': to,
        },
      );
      return response.data['transactionHash'];
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }
  
  // Handle API errors
  ApiException _handleError(DioException e) {
    String message = 'An error occurred';
    int? statusCode = e.response?.statusCode;
    
    if (e.response != null) {
      if (e.response!.data is Map) {
        message = e.response!.data['message'] ?? message;
      }
    } else if (e.type == DioExceptionType.connectionTimeout) {
      message = 'Connection timeout';
    } else if (e.type == DioExceptionType.receiveTimeout) {
      message = 'Receive timeout';
    } else if (e.type == DioExceptionType.sendTimeout) {
      message = 'Send timeout';
    } else if (e.type == DioExceptionType.connectionError) {
      message = 'Connection error';
    }
    
    return ApiException(message, statusCode: statusCode);
  }
}