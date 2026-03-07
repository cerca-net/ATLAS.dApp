import 'package:firebase_auth/firebase_auth.dart' as firebase;
import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:wallet_integration/models/user.dart';
import 'api_client.dart';

class AuthException implements Exception {
  final String message;
  final dynamic cause;
  
  AuthException(this.message, [this.cause]);
  
  @override
  String toString() => 'AuthException: $message${cause != null ? ' ($cause)' : ''}';
}

class AuthService {
  final firebase.FirebaseAuth _firebaseAuth;
  final FirebaseFirestore _firestore;
  final ApiClient _apiClient;
  
  AuthService(
    this._firebaseAuth,
    this._firestore,
    this._apiClient,
  );
  
  // Register a new user
  Future<User> register(String email, String password, String displayName) async {
    try {
      // Create user in Firebase
      final userCredential = await _firebaseAuth.createUserWithEmailAndPassword(
        email: email,
        password: password,
      );
      
      final firebaseUser = userCredential.user;
      if (firebaseUser == null) {
        throw AuthException('Failed to create user');
      }
      
      // Update display name
      await firebaseUser.updateDisplayName(displayName);
      
      // Create wallet
      final walletResponse = await _apiClient.createWallet();
      
      // Create user document in Firestore
      final user = User(
        uid: firebaseUser.uid,
        email: email,
        displayName: displayName,
        walletAddress: walletResponse.address,
      );
      
      await _firestore.collection('users').doc(user.uid).set(user.toMap());
      
      return user;
    } on firebase.FirebaseAuthException catch (e) {
      throw AuthException(_getFirebaseErrorMessage(e.code), e);
    } on ApiException catch (e) {
      throw AuthException('Failed to create wallet: ${e.message}', e);
    } catch (e) {
      throw AuthException('Registration failed', e);
    }
  }
  
  // Login existing user
  Future<User> login(String email, String password) async {
    try {
      // Sign in with Firebase
      final userCredential = await _firebaseAuth.signInWithEmailAndPassword(
        email: email,
        password: password,
      );
      
      final firebaseUser = userCredential.user;
      if (firebaseUser == null) {
        throw AuthException('Failed to login');
      }
      
      // Get user data from Firestore
      final docSnapshot = await _firestore.collection('users').doc(firebaseUser.uid).get();
      
      if (!docSnapshot.exists) {
        throw AuthException('User data not found');
      }
      
      return User.fromMap(docSnapshot.data()!);
    } on firebase.FirebaseAuthException catch (e) {
      throw AuthException(_getFirebaseErrorMessage(e.code), e);
    } catch (e) {
      throw AuthException('Login failed', e);
    }
  }
  
  // Get current user
  Future<User?> getCurrentUser() async {
    final firebaseUser = _firebaseAuth.currentUser;
    if (firebaseUser == null) {
      return null;
    }
    
    try {
      final docSnapshot = await _firestore.collection('users').doc(firebaseUser.uid).get();
      
      if (!docSnapshot.exists) {
        return null;
      }
      
      return User.fromMap(docSnapshot.data()!);
    } catch (e) {
      throw AuthException('Failed to get current user', e);
    }
  }
  
  // Logout
  Future<void> logout() async {
    await _firebaseAuth.signOut();
  }
  
  // Helper method to get user-friendly error messages
  String _getFirebaseErrorMessage(String code) {
    switch (code) {
      case 'user-not-found':
        return 'No user found with this email';
      case 'wrong-password':
        return 'Incorrect password';
      case 'email-already-in-use':
        return 'Email is already in use';
      case 'weak-password':
        return 'Password is too weak';
      case 'invalid-email':
        return 'Invalid email address';
      default:
        return 'Authentication failed';
    }
  }
}