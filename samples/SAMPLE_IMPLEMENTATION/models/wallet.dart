class Wallet {
  final String address;
  final double balance;
  final int nonce;
  final List<Transaction> recentTransactions;
  
  Wallet({
    required this.address,
    required this.balance,
    required this.nonce,
    required this.recentTransactions,
  });
  
  factory Wallet.fromJson(Map<String, dynamic> json) {
    return Wallet(
      address: json['address'] ?? '',
      balance: (json['balance'] ?? 0.0).toDouble(),
      nonce: json['nonce'] ?? 0,
      recentTransactions: (json['recentTransactions'] as List<dynamic>?)
          ?.map((tx) => Transaction.fromJson(tx))
          .toList() ?? [],
    );
  }
  
  Map<String, dynamic> toJson() {
    return {
      'address': address,
      'balance': balance,
      'nonce': nonce,
      'recentTransactions': recentTransactions.map((tx) => tx.toJson()).toList(),
    };
  }
}

class WalletResponse {
  final String address;
  final String privateKey;
  final String sessionToken;
  
  WalletResponse({
    required this.address,
    required this.privateKey,
    required this.sessionToken,
  });
  
  factory WalletResponse.fromJson(Map<String, dynamic> json) {
    return WalletResponse(
      address: json['address'] ?? '',
      privateKey: json['privateKey'] ?? '',
      sessionToken: json['sessionToken'] ?? '',
    );
  }
  
  Map<String, dynamic> toJson() {
    return {
      'address': address,
      'privateKey': privateKey,
      'sessionToken': sessionToken,
    };
  }
}

class Transaction {
  final String hash;
  final String from;
  final String to;
  final double amount;
  final int timestamp;
  final String status;
  
  Transaction({
    required this.hash,
    required this.from,
    required this.to,
    required this.amount,
    required this.timestamp,
    required this.status,
  });
  
  factory Transaction.fromJson(Map<String, dynamic> json) {
    return Transaction(
      hash: json['hash'] ?? '',
      from: json['from'] ?? '',
      to: json['to'] ?? '',
      amount: (json['amount'] ?? 0.0).toDouble(),
      timestamp: json['timestamp'] ?? 0,
      status: json['status'] ?? 'pending',
    );
  }
  
  Map<String, dynamic> toJson() {
    return {
      'hash': hash,
      'from': from,
      'to': to,
      'amount': amount,
      'timestamp': timestamp,
      'status': status,
    };
  }
}