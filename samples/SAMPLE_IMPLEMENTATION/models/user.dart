import 'package:cloud_firestore/cloud_firestore.dart';

class User {
  final String uid;
  final String email;
  final String? displayName;
  final String? walletAddress;
  
  User({
    required this.uid,
    required this.email,
    this.displayName,
    this.walletAddress,
  });
  
  // Create User from Firestore data
  factory User.fromMap(Map<String, dynamic> map) {
    return User(
      uid: map['uid'] ?? '',
      email: map['email'] ?? '',
      displayName: map['displayName'],
      walletAddress: map['walletAddress'],
    );
  }
  
  // Convert User to Map for Firestore
  Map<String, dynamic> toMap() {
    return {
      'uid': uid,
      'email': email,
      'displayName': displayName,
      'walletAddress': walletAddress,
      'createdAt': FieldValue.serverTimestamp(),
    };
  }
  
  // Create a copy of User with updated fields
  User copyWith({
    String? uid,
    String? email,
    String? displayName,
    String? walletAddress,
  }) {
    return User(
      uid: uid ?? this.uid,
      email: email ?? this.email,
      displayName: displayName ?? this.displayName,
      walletAddress: walletAddress ?? this.walletAddress,
    );
  }
  
  @override
  String toString() {
    return 'User(uid: $uid, email: $email, displayName: $displayName, walletAddress: $walletAddress)';
  }
}