import 'dart:convert';
import 'package:http/http.dart' as http;

/// Blockchain service for integrating with ATLAS blockchain
class BlockchainService {
  final String baseUrl;

  BlockchainService({this.baseUrl = 'http://192.168.0.105:8080'}); // LAN testnet

  /// Connect wallet to ATLAS blockchain
  Future<WalletConnectionResponse> connectWallet(
      WalletConnectRequest request) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/flutterflow/connect-wallet'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode(request.toJson()),
      );

      if (response.statusCode == 200) {
        return WalletConnectionResponse.fromJson(jsonDecode(response.body));
      } else {
        throw BlockchainException('Failed to connect wallet: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error connecting wallet: $e');
    }
  }

  /// Get wallet information
  Future<WalletInfoResponse> getWalletInfo(String address) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/flutterflow/wallet-info?address=$address'),
      );

      if (response.statusCode == 200) {
        return WalletInfoResponse.fromJson(jsonDecode(response.body));
      } else {
        throw BlockchainException(
            'Failed to get wallet info: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error getting wallet info: $e');
    }
  }

  /// Send blockchain transaction
  Future<SendTransactionResponse> sendTransaction(
      SendTransactionRequest request) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/submit-transaction'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode(request.toJson()),
      );

      if (response.statusCode == 200 || response.statusCode == 202) {
        return SendTransactionResponse.fromJson(jsonDecode(response.body));
      } else {
        throw BlockchainException(
            'Failed to send transaction: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error sending transaction: $e');
    }
  }

  /// Get transaction history
  Future<TransactionHistoryResponse> getTransactionHistory(
      String address) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/flutterflow/transaction-history?address=$address'),
      );

      if (response.statusCode == 200) {
        return TransactionHistoryResponse.fromJson(jsonDecode(response.body));
      } else {
        throw BlockchainException(
            'Failed to get transaction history: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException(
          'Network error getting transaction history: $e');
    }
  }

  /// Authenticate wallet session
  Future<AuthResponse> authenticate(String sessionToken, String address) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/flutterflow/authenticate'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'sessionToken': sessionToken,
          'address': address,
        }),
      );

      if (response.statusCode == 200) {
        return AuthResponse.fromJson(jsonDecode(response.body));
      } else {
        throw BlockchainException('Failed to authenticate: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error authenticating: $e');
    }
  }

  /// Start the node
  Future<void> startNode() async {
    try {
      final response = await http.post(Uri.parse('$baseUrl/node/start'));
      if (response.statusCode != 200) {
        throw BlockchainException('Failed to start node: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error starting node: $e');
    }
  }

  /// Stop the node
  Future<void> stopNode() async {
    try {
      final response = await http.post(Uri.parse('$baseUrl/node/stop'));
      if (response.statusCode != 200) {
        throw BlockchainException('Failed to stop node: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error stopping node: $e');
    }
  }

  /// Pause the node
  Future<void> pauseNode() async {
    try {
      final response = await http.post(Uri.parse('$baseUrl/node/pause'));
      if (response.statusCode != 200) {
        throw BlockchainException('Failed to pause node: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error pausing node: $e');
    }
  }

  /// Sync the node
  Future<void> syncNode() async {
    try {
      final response = await http.post(Uri.parse('$baseUrl/node/sync'));
      if (response.statusCode != 200) {
        throw BlockchainException('Failed to sync node: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error syncing node: $e');
    }
  }

  /// Get aggregated network status
  Future<NetworkStatusResponse> getNetworkStatus({String? address}) async {
    try {
      final addressParam = (address != null && address.isNotEmpty) ? '?address=$address' : '';
      final response = await http.get(Uri.parse('$baseUrl/status$addressParam'));

      if (response.statusCode == 200) {
        return NetworkStatusResponse.fromJson(jsonDecode(response.body));
      } else {
        throw BlockchainException('Failed to get status from network: ${response.body}');
      }
    } catch (e) {
      print('Network error getting network status: $e');
      return NetworkStatusResponse(
        blockHeight: 0,
        txPoolSize: 0,
        isValidator: false,
        validatorAddress: '',
        stakeAmount: 0,
        rewardsEarned: 0,
        totalValidators: 0,
        walletBalance: 0,
        walletStaked: 0,
        totalBalance: 0,
        mode: 'offline',
      );
    }
  }

  /// Get node status
  Future<NodeStatusResponse> getNodeStatus() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/node/status'));
      if (response.statusCode == 200) {
        return NodeStatusResponse.fromJson(jsonDecode(response.body));
      } else {
        throw BlockchainException(
            'Failed to get node status: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error getting node status: $e');
    }
  }

  /// Get node logs
  Future<NodeLogsResponse> getNodeLogs({int limit = 50}) async {
    try {
      final response =
          await http.get(Uri.parse('$baseUrl/node/logs?limit=$limit'));
      if (response.statusCode == 200) {
        return NodeLogsResponse.fromJson(jsonDecode(response.body));
      } else {
        throw BlockchainException('Failed to get node logs: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error getting node logs: $e');
    }
  }

  /// Get connected P2P peers
  Future<Map<String, dynamic>> getPeers() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/peers'));
      if (response.statusCode == 200) {
        return jsonDecode(response.body) as Map<String, dynamic>;
      } else {
        throw BlockchainException('Failed to get peers: ${response.body}');
      }
    } catch (e) {
      print('Network error getting peers: $e');
      return {'success': false, 'count': 0, 'peers': []};
    }
  }

  void dispose() {
    // No-op
  }

  /// Request faucet tokens
  Future<Map<String, dynamic>> requestFaucet(String address) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/faucet'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'address': address}),
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        throw BlockchainException('Failed to request faucet: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error requesting faucet: $e');
    }
  }

  /// Get current nonce for address to prevent replay attacks
  Future<int> getNonce(String address) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/nonce?address=$address'),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        return data['nonce'] ?? 0;
      } else {
        // If the endpoint doesn't exist yet or errors, default to 0 (or handle as needed)
        print('Warning: Failed to fetch nonce: ${response.body}');
        return 0;
      }
    } catch (e) {
      print('Network error fetching nonce: $e');
      return 0; // Fallback
    }
  }

  /// Register as a validator (Legacy endpoint, now expecting signed transaction)
  Future<Map<String, dynamic>> registerValidator(
      SendTransactionRequest request) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/register-validator'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode(request.toJson()),
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        throw BlockchainException(
            'Failed to register validator: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error registering validator: $e');
    }
  }

  /// Get recent blocks
  Future<List<BlockData>> getBlocks({int limit = 20}) async {
    try {
      final response =
          await http.get(Uri.parse('$baseUrl/blocks?limit=$limit'));
      if (response.statusCode == 200) {
        final List<dynamic> data = jsonDecode(response.body);
        return data.map((json) => BlockData.fromJson(json)).toList();
      } else {
        throw BlockchainException('Failed to get blocks: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error getting blocks: $e');
    }
  }

  // ---------------------------------------------------------------------------
  // SOCIAL LAYER ENDPOINTS
  // ---------------------------------------------------------------------------

  /// Get Social Feed (Global or User-specific)
  Future<List<SocialPost>> getSocialFeed(
      {String? userId, int limit = 20}) async {
    try {
      final queryParams =
          userId != null ? '?user_id=$userId&limit=$limit' : '?limit=$limit';
      final response =
          await http.get(Uri.parse('$baseUrl/social/feed$queryParams'));

      if (response.statusCode == 200) {
        final Map<String, dynamic> jsonResponse = jsonDecode(response.body);
        if (jsonResponse.containsKey('data')) {
          final List<dynamic> data = jsonResponse['data'];
          return data.map((json) => SocialPost.fromJson(json)).toList();
        } else if (jsonResponse is List) {
          // Fallback if API changes to return list directly
          return (jsonResponse as List)
              .map((json) => SocialPost.fromJson(json))
              .toList();
        }
        return [];
      } else {
        // Fallback for empty feed or errors
        print('Warning: Failed to fetch feed: ${response.body}');
        return [];
      }
    } catch (e) {
      print('Network error fetching feed: $e');
      return [];
    }
  }

  /// Create a new Post
  Future<SocialPost> createPost(String author, String content,
      {List<String>? mediaUrls}) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/social/post/create'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'author': author,
          'content': content,
          'media_urls': mediaUrls ?? [],
          'visibility': 'public', // Default for now
          'category': 'general',
        }),
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final Map<String, dynamic> jsonResponse = jsonDecode(response.body);
        if (jsonResponse.containsKey('data')) {
          return SocialPost.fromJson(jsonResponse['data']);
        }
        return SocialPost.fromJson(jsonResponse);
      } else {
        throw BlockchainException('Failed to create post: ${response.body}');
      }
    } catch (e) {
      throw BlockchainException('Network error creating post: $e');
    }
  }

  /// Like a Post (Interact with the "Object")
  Future<bool> likePost(String postId, String userId,
      {String type = 'like'}) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/social/like'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'post_id': postId,
          'user_id': userId,
          'like_type': type,
        }),
      );

      return response.statusCode == 200;
    } catch (e) {
      print('Network error liking post: $e');
      return false;
    }
  }

  /// Tip a Post (Transfer Energy)
  Future<bool> tipPost(String postId, String userId, int amount) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/social/tip'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'post_id': postId,
          'user_id': userId,
          'amount': amount,
        }),
      );

      return response.statusCode == 200;
    } catch (e) {
      print('Network error tipping post: $e');
      return false;
    }
  }

  /// Get User Identity for Social features
  Future<UserIdentity?> getSocialIdentity(String address,
      {String? requester}) async {
    try {
      final requesterParam = requester != null ? '&requester=$requester' : '';
      final response = await http.get(
        Uri.parse('$baseUrl/identity/social?address=$address$requesterParam'),
      );

      if (response.statusCode == 200) {
        final Map<String, dynamic> jsonResponse = jsonDecode(response.body);
        if (jsonResponse['success'] == true &&
            jsonResponse['identity'] != null) {
          return UserIdentity.fromJson(jsonResponse['identity']);
        }
      }
      return null;
    } catch (e) {
      print('Network error getting social identity: $e');
      return null;
    }
  }
}

/// Request/Response models
class WalletConnectRequest {
  final String action; // 'create', 'import', 'connect'
  final String? privateKey;
  final String? address;
  final String? sessionId;

  WalletConnectRequest({
    required this.action,
    this.privateKey,
    this.address,
    this.sessionId,
  });

  Map<String, dynamic> toJson() => {
        'action': action,
        if (privateKey != null) 'privateKey': privateKey,
        if (address != null) 'address': address,
        if (sessionId != null) 'sessionId': sessionId,
      };
}

class WalletConnectionResponse {
  final bool success;
  final String message;
  final WalletData data;

  WalletConnectionResponse({
    required this.success,
    required this.message,
    required this.data,
  });

  factory WalletConnectionResponse.fromJson(Map<String, dynamic> json) {
    return WalletConnectionResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
      data: WalletData.fromJson(json['data'] ?? {}),
    );
  }
}

class WalletData {
  final String address;
  final String sessionToken;
  final double balance;
  final bool isValidator;
  final String? mnemonic;

  WalletData({
    required this.address,
    required this.sessionToken,
    required this.balance,
    required this.isValidator,
    this.mnemonic,
  });

  factory WalletData.fromJson(Map<String, dynamic> json) {
    return WalletData(
      address: json['address'] ?? '',
      sessionToken: json['sessionToken'] ?? '',
      balance: (json['balance'] ?? 0).toDouble(),
      isValidator: json['isValidator'] ?? false,
      mnemonic: json['mnemonic'],
    );
  }
}

class WalletInfoResponse {
  final bool success;
  final WalletInfoData data;

  WalletInfoResponse({
    required this.success,
    required this.data,
  });

  factory WalletInfoResponse.fromJson(Map<String, dynamic> json) {
    return WalletInfoResponse(
      success: json['success'] ?? false,
      data: WalletInfoData.fromJson(json['data'] ?? {}),
    );
  }
}

class WalletInfoData {
  final String address;
  final double balance;
  final bool isValidator;
  final List<TransactionData> recentTransactions;
  final int nonce;

  WalletInfoData({
    required this.address,
    required this.balance,
    required this.isValidator,
    required this.recentTransactions,
    required this.nonce,
  });

  factory WalletInfoData.fromJson(Map<String, dynamic> json) {
    return WalletInfoData(
      address: json['address'] ?? '',
      balance: (json['balance'] ?? 0).toDouble(),
      isValidator: json['isValidator'] ?? false,
      recentTransactions: (json['recentTransactions'] as List<dynamic>?)
              ?.map((tx) => TransactionData.fromJson(tx))
              .toList() ??
          [],
      nonce: json['nonce'] ?? 0,
    );
  }
}

class SendTransactionRequest {
  final String sender;
  final String senderPublicKey; // Added
  final String recipient;
  final int amount;
  final int fee;
  final int timestamp; // Added
  final int nonce; // Added
  final String? data;
  final String signature;
  final String type; // Added type field

  SendTransactionRequest({
    required this.sender,
    required this.senderPublicKey,
    required this.recipient,
    required this.amount,
    required this.fee,
    required this.timestamp,
    required this.nonce,
    this.data,
    required this.signature,
    this.type = 'regular', // Default to regular
  });

  Map<String, dynamic> toJson() => {
        'type': type,
        'Sender': sender,
        'SenderPublicKey': senderPublicKey,
        'Recipient': recipient,
        'Amount': amount,
        'Fee': fee,
        'Timestamp': timestamp,
        'Nonce': nonce,
        if (data != null) 'Data': data,
        'Signature': signature,
      };
}

class SendTransactionResponse {
  final bool success;
  final String message;
  final TransactionResultData data;

  SendTransactionResponse({
    required this.success,
    required this.message,
    required this.data,
  });

  factory SendTransactionResponse.fromJson(Map<String, dynamic> json) {
    return SendTransactionResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
      data: TransactionResultData.fromJson(json['data'] ?? {}),
    );
  }
}

class TransactionResultData {
  final String transactionHash;
  final String from;
  final String to;
  final int amount;
  final String status;

  TransactionResultData({
    required this.transactionHash,
    required this.from,
    required this.to,
    required this.amount,
    required this.status,
  });

  factory TransactionResultData.fromJson(Map<String, dynamic> json) {
    return TransactionResultData(
      transactionHash: json['transactionHash'] ?? '',
      from: json['from'] ?? '',
      to: json['to'] ?? '',
      amount: json['amount'] ?? 0,
      status: json['status'] ?? '',
    );
  }
}

class TransactionHistoryResponse {
  final bool success;
  final TransactionHistoryData data;

  TransactionHistoryResponse({
    required this.success,
    required this.data,
  });

  factory TransactionHistoryResponse.fromJson(Map<String, dynamic> json) {
    return TransactionHistoryResponse(
      success: json['success'] ?? false,
      data: TransactionHistoryData.fromJson(json['data'] ?? {}),
    );
  }
}

class TransactionHistoryData {
  final String address;
  final List<TransactionData> transactions;
  final int totalCount;

  TransactionHistoryData({
    required this.address,
    required this.transactions,
    required this.totalCount,
  });

  factory TransactionHistoryData.fromJson(Map<String, dynamic> json) {
    return TransactionHistoryData(
      address: json['address'] ?? '',
      transactions: (json['transactions'] as List<dynamic>?)
              ?.map((tx) => TransactionData.fromJson(tx))
              .toList() ??
          [],
      totalCount: json['totalCount'] ?? 0,
    );
  }
}

class TransactionData {
  final String hash;
  final int blockIndex;
  final String sender;
  final String recipient;
  final int amount;
  final int fee;
  final int timestamp;
  final String type;

  TransactionData({
    required this.hash,
    required this.blockIndex,
    required this.sender,
    required this.recipient,
    required this.amount,
    required this.fee,
    required this.timestamp,
    required this.type,
  });

  factory TransactionData.fromJson(Map<String, dynamic> json) {
    return TransactionData(
      hash: json['Hash'] ?? json['hash'] ?? '',
      blockIndex: json['BlockIndex'] ?? json['blockIndex'] ?? 0,
      sender: json['Sender'] ?? json['sender'] ?? '',
      recipient: json['Recipient'] ?? json['recipient'] ?? '',
      amount: (json['Amount'] ?? json['amount'] ?? 0).toInt(),
      fee: (json['Fee'] ?? json['fee'] ?? 0).toInt(),
      timestamp: json['Timestamp'] ?? json['timestamp'] ?? 0,
      type: json['Type'] ?? json['type'] ?? '',
      // data: json['Data'] ?? json['data'] ?? '', // Add if needed in model
    );
  }
}

class AuthResponse {
  final bool success;
  final String message;
  final AuthData data;

  AuthResponse({
    required this.success,
    required this.message,
    required this.data,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    return AuthResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
      data: AuthData.fromJson(json['data'] ?? {}),
    );
  }
}

class AuthData {
  final String address;
  final double balance;
  final bool isValidator;

  AuthData({
    required this.address,
    required this.balance,
    required this.isValidator,
  });

  factory AuthData.fromJson(Map<String, dynamic> json) {
    return AuthData(
      address: json['address'] ?? '',
      balance: (json['balance'] ?? 0).toDouble(),
      isValidator: json['isValidator'] ?? false,
    );
  }
}

class DisconnectResponse {
  final bool success;
  final String message;
  final DisconnectData data;

  DisconnectResponse({
    required this.success,
    required this.message,
    required this.data,
  });

  factory DisconnectResponse.fromJson(Map<String, dynamic> json) {
    return DisconnectResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
      data: DisconnectData.fromJson(json['data'] ?? {}),
    );
  }
}

class DisconnectData {
  final String address;
  final int disconnectedAt;

  DisconnectData({
    required this.address,
    required this.disconnectedAt,
  });

  factory DisconnectData.fromJson(Map<String, dynamic> json) {
    return DisconnectData(
      address: json['address'] ?? '',
      disconnectedAt: json['disconnectedAt'] ?? 0,
    );
  }
}

class NetworkStatusResponse {
  final int blockHeight;
  final int txPoolSize;
  final bool isValidator;
  final String validatorAddress;
  final int stakeAmount;
  final int rewardsEarned;
  final int totalValidators;
  final double walletBalance; // Changed to double to match other models
  final int walletStaked;
  final double totalBalance; // Changed to double
  final String mode;

  NetworkStatusResponse({
    required this.blockHeight,
    required this.txPoolSize,
    required this.isValidator,
    required this.validatorAddress,
    required this.stakeAmount,
    required this.rewardsEarned,
    required this.totalValidators,
    required this.walletBalance,
    required this.walletStaked,
    required this.totalBalance,
    required this.mode,
  });

  factory NetworkStatusResponse.fromJson(Map<String, dynamic> json) {
    return NetworkStatusResponse(
      blockHeight: json['blockHeight'] ?? 0,
      txPoolSize: json['txPoolSize'] ?? 0,
      isValidator: json['isValidator'] ?? false,
      validatorAddress: json['validatorAddress'] ?? '',
      stakeAmount: json['stakeAmount'] ?? 0,
      rewardsEarned: json['rewardsEarned'] ?? 0,
      totalValidators: json['totalValidators'] ?? 0,
      walletBalance: (json['walletBalance'] ?? 0).toDouble(),
      walletStaked: json['walletStaked'] ?? 0,
      totalBalance: (json['totalBalance'] ?? 0).toDouble(),
      mode: json['mode'] ?? 'observer',
    );
  }
}

/// Node status response from /node/status
class NodeStatusResponse {
  final String state; // running, paused, stopped, syncing
  final bool isValidator;
  final String validatorAddress;
  final int stakeAmount;
  final int blockHeight;
  final int txPoolSize;
  final int totalValidators;
  final String uptime;

  NodeStatusResponse({
    required this.state,
    required this.isValidator,
    required this.validatorAddress,
    required this.stakeAmount,
    required this.blockHeight,
    required this.txPoolSize,
    required this.totalValidators,
    required this.uptime,
  });

  factory NodeStatusResponse.fromJson(Map<String, dynamic> json) {
    return NodeStatusResponse(
      state: json['state'] ?? 'stopped',
      isValidator: json['isValidator'] ?? false,
      validatorAddress: json['validatorAddress'] ?? '',
      stakeAmount: json['stakeAmount'] ?? 0,
      blockHeight: json['blockHeight'] ?? 0,
      txPoolSize: json['txPoolSize'] ?? 0,
      totalValidators: json['totalValidators'] ?? 0,
      uptime: json['uptime'] ?? '',
    );
  }
}

/// Single log entry from node
class NodeLogEntry {
  final String timestamp;
  final String level; // info, success, warning, error
  final String message;

  NodeLogEntry({
    required this.timestamp,
    required this.level,
    required this.message,
  });

  factory NodeLogEntry.fromJson(Map<String, dynamic> json) {
    return NodeLogEntry(
      timestamp: json['timestamp'] ?? '',
      level: json['level'] ?? 'info',
      message: json['message'] ?? '',
    );
  }
}

/// Node logs response from /node/logs
class NodeLogsResponse {
  final List<NodeLogEntry> logs;
  final int total;
  final String state;

  NodeLogsResponse({
    required this.logs,
    required this.total,
    required this.state,
  });

  factory NodeLogsResponse.fromJson(Map<String, dynamic> json) {
    return NodeLogsResponse(
      logs: (json['logs'] as List<dynamic>?)
              ?.map((e) => NodeLogEntry.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      total: json['total'] ?? 0,
      state: json['state'] ?? 'stopped',
    );
  }
}

class BlockData {
  final int index;
  final String hash;
  final String prevHash;
  final int timestamp;
  final List<TransactionData> transactions;
  final String validator;
  final String signature;

  BlockData({
    required this.index,
    required this.hash,
    required this.prevHash,
    required this.timestamp,
    required this.transactions,
    required this.validator,
    required this.signature,
  });

  factory BlockData.fromJson(Map<String, dynamic> json) {
    return BlockData(
      index: json['Index'] ?? json['index'] ?? 0,
      hash: json['Hash'] ?? json['hash'] ?? '',
      prevHash: json['PrevHash'] ?? json['prevHash'] ?? '',
      timestamp: json['Timestamp'] ?? json['timestamp'] ?? 0,
      transactions:
          ((json['Transactions'] ?? json['transactions']) as List<dynamic>?)
                  ?.map((tx) => TransactionData.fromJson(tx))
                  .toList()
                  .cast<TransactionData>() ??
              [],
      validator: json['Validator'] ?? json['validator'] ?? '',
      signature: json['Signature'] ?? json['signature'] ?? '',
    );
  }
}

/// Custom exception for blockchain errors
class BlockchainException implements Exception {
  final String message;
  BlockchainException(this.message);

  @override
  String toString() => 'BlockchainException: $message';
}

// ---------------------------------------------------------------------------
// SOCIAL MODELS
// ---------------------------------------------------------------------------

class SocialPost {
  final String id;
  final String author;
  final String content;
  final List<String> mediaUrls;
  final int likes;
  final int comments;
  final int shares;
  final int views;
  final int timestamp;
  final Map<String, dynamic> metadata;

  // Physics Properties (Smart Contract State)
  final int tipBalance;
  final double influenceScore;
  final int upvotes;
  final int downvotes;

  SocialPost({
    required this.id,
    required this.author,
    required this.content,
    required this.mediaUrls,
    required this.likes,
    required this.comments,
    required this.shares,
    required this.views,
    required this.timestamp,
    required this.metadata,
    this.tipBalance = 0,
    this.influenceScore = 0.0,
    this.upvotes = 0,
    this.downvotes = 0,
  });

  factory SocialPost.fromJson(Map<String, dynamic> json) {
    // Handle specific timestamp format from Go (RFC3339 string) to int if needed
    // Assuming backend might send string ISO date or unix int.
    int ts = 0;
    if (json['created_at'] is String) {
      try {
        ts = DateTime.parse(json['created_at']).millisecondsSinceEpoch ~/ 1000;
      } catch (e) {
        ts = 0;
      }
    } else if (json['created_at'] is int) {
      ts = json['created_at'];
    }

    return SocialPost(
      id: json['id'] ?? '',
      author: json['author'] ?? '',
      content: json['content'] ?? '',
      mediaUrls: (json['media_urls'] as List<dynamic>?)?.cast<String>() ?? [],
      likes: json['likes'] ?? 0,
      comments: json['comments'] ?? 0,
      shares: json['shares'] ?? 0,
      views: json['views'] ?? 0,
      timestamp: ts,
      metadata: json['metadata'] ?? {},
      tipBalance: json['tip_balance'] ?? 0,
      influenceScore: (json['influence_score'] ?? 0).toDouble(),
      upvotes: json['upvotes'] ?? 0,
      downvotes: json['downvotes'] ?? 0,
    );
  }
}

class UserIdentity {
  final String address;
  final String username;
  final UserProfile? profile;
  final ReputationScore? reputation;
  final ActivityMetrics? activity;

  UserIdentity({
    required this.address,
    required this.username,
    this.profile,
    this.reputation,
    this.activity,
  });

  factory UserIdentity.fromJson(Map<String, dynamic> json) {
    return UserIdentity(
      address: json['address'] ?? '',
      username: json['username'] ?? '',
      profile: json['profile'] != null
          ? UserProfile.fromJson(json['profile'])
          : null,
      reputation: json['reputation'] != null
          ? ReputationScore.fromJson(json['reputation'])
          : null,
      activity: json['activity'] != null
          ? ActivityMetrics.fromJson(json['activity'])
          : null,
    );
  }
}

class UserProfile {
  final String displayName;
  final String bio;
  final String avatar;
  final bool isPublic;

  UserProfile({
    required this.displayName,
    required this.bio,
    required this.avatar,
    required this.isPublic,
  });

  factory UserProfile.fromJson(Map<String, dynamic> json) {
    return UserProfile(
      displayName: json['display_name'] ?? '',
      bio: json['bio'] ?? '',
      avatar: json['avatar'] ?? '',
      isPublic: json['is_public'] ?? true,
    );
  }
}

class ReputationScore {
  final double overall;
  final double social;
  final double commerce;
  final double governance;

  ReputationScore({
    required this.overall,
    required this.social,
    required this.commerce,
    required this.governance,
  });

  factory ReputationScore.fromJson(Map<String, dynamic> json) {
    return ReputationScore(
      overall: (json['overall'] ?? 0).toDouble(),
      social: (json['social'] ?? 0).toDouble(),
      commerce: (json['commerce'] ?? 0).toDouble(),
      governance: (json['governance'] ?? 0).toDouble(),
    );
  }
}

class ActivityMetrics {
  final int postsCreated;
  final int commentsMade;
  final int likesGiven;
  final int likesReceived;
  final int totalTokensEarned;

  ActivityMetrics({
    required this.postsCreated,
    required this.commentsMade,
    required this.likesGiven,
    required this.likesReceived,
    required this.totalTokensEarned,
  });

  factory ActivityMetrics.fromJson(Map<String, dynamic> json) {
    return ActivityMetrics(
      postsCreated: json['posts_created'] ?? 0,
      commentsMade: json['comments_made'] ?? 0,
      likesGiven: json['likes_given'] ?? 0,
      likesReceived: json['likes_received'] ?? 0,
      totalTokensEarned: json['total_tokens_earned'] ?? 0,
    );
  }
}
