import 'dart:convert';
import 'package:http/http.dart' as http;
import 'blockchain_service.dart';

/// Represents the energy physics state of any object (post or item)
/// bridged from Firebase via its document ID.
class ObjectEnergyState {
  final String objectId;
  final int tipBalance;
  final double influenceScore;
  final String status; // "active", "fossilized"
  final int upvotes;
  final int downvotes;
  final String objectType; // "post" or "item"

  ObjectEnergyState({
    required this.objectId,
    required this.tipBalance,
    required this.influenceScore,
    required this.status,
    required this.upvotes,
    required this.downvotes,
    required this.objectType,
  });

  factory ObjectEnergyState.fromJson(Map<String, dynamic> json) {
    return ObjectEnergyState(
      objectId: json['object_id'] ?? '',
      tipBalance: (json['tip_balance'] ?? 0).toInt(),
      influenceScore: (json['influence_score'] ?? 0.0).toDouble(),
      status: json['status'] ?? 'active',
      upvotes: (json['upvotes'] ?? 0).toInt(),
      downvotes: (json['downvotes'] ?? 0).toInt(),
      objectType: json['object_type'] ?? 'post',
    );
  }

  bool get isFossilized => status == 'fossilized';
  bool get isActive => status == 'active';

  /// Energy level as a 0.0 to 1.0 ratio (capped at 1.0 for display)
  /// Based on fossilization threshold of 50
  double get energyRatio {
    const maxDisplay = 200.0; // Display scale max
    return (tipBalance / maxDisplay).clamp(0.0, 1.0);
  }
}

/// Service that bridges Firebase objects with blockchain energy state.
/// Every SubmissionRecord (post or item) can have its energy queried
/// and modified through this service.
class EnergyService {
  static final EnergyService _instance = EnergyService._internal();

  factory EnergyService() => _instance;

  EnergyService._internal();

  String get _baseUrl => BlockchainService().baseUrl;

  /// Get energy state for a Firebase submission object.
  /// If the object doesn't exist in the blockchain yet, it gets auto-registered
  /// with default energy (100 TCOIN grace period).
  Future<ObjectEnergyState?> getObjectEnergy(
    String firebaseDocId, {
    String objectType = 'post',
  }) async {
    try {
      final response = await http.get(
        Uri.parse(
            '$_baseUrl/social/object/energy?object_id=$firebaseDocId&object_type=$objectType'),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        if (data['success'] == true && data['data'] != null) {
          return ObjectEnergyState.fromJson(data['data']);
        }
      }
      return null;
    } catch (e) {
      print('EnergyService: Error fetching energy for $firebaseDocId: $e');
      return null;
    }
  }

  /// Send energy (TCOIN) to an object. Works for both posts and items.
  /// Handles revival of fossilized objects automatically.
  Future<ObjectEnergyState?> energizeObject(
    String firebaseDocId,
    String walletAddress,
    int amount,
  ) async {
    try {
      final response = await http.post(
        Uri.parse('$_baseUrl/social/object/energize'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'object_id': firebaseDocId,
          'user_id': walletAddress,
          'amount': amount,
        }),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        if (data['success'] == true && data['data'] != null) {
          return ObjectEnergyState.fromJson(data['data']);
        }
      }
      return null;
    } catch (e) {
      print('EnergyService: Error energizing $firebaseDocId: $e');
      return null;
    }
  }
}
