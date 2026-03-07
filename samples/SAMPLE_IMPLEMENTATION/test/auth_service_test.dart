import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:firebase_auth/firebase_auth.dart' as firebase;
import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:wallet_integration/services/auth_service.dart';
import 'package:wallet_integration/services/api_client.dart';
import 'package:wallet_integration/models/user.dart' as app;
import 'package:wallet_integration/models/wallet.dart';

// Create mock classes for Firebase Auth and Firestore
class MockFirebaseAuth extends Mock implements firebase.FirebaseAuth {}
class MockFirebaseUser extends Mock implements firebase.User {}
class MockUserCredential extends Mock implements firebase.UserCredential {}
class MockFirestore extends Mock implements FirebaseFirestore {}
class MockCollectionReference extends Mock implements CollectionReference {}
class MockDocumentReference extends Mock implements DocumentReference {}

@GenerateMocks([ApiClient])
import 'auth_service_test.mocks.dart';

void main() {
  group('AuthService Tests', () {
    late AuthService authService;
    late MockFirebaseAuth mockFirebaseAuth;
    late MockFirestore mockFirestore;
    late MockApiClient mockApiClient;
    late MockUserCredential mockUserCredential;
    late MockFirebaseUser mockFirebaseUser;
    
    setUp(() {
      mockFirebaseAuth = MockFirebaseAuth();
      mockFirestore = MockFirestore();
      mockApiClient = MockApiClient();
      mockUserCredential = MockUserCredential();
      mockFirebaseUser = MockFirebaseUser();
      
      authService = AuthService(
        mockFirebaseAuth,
        mockFirestore,
        mockApiClient,
      );
    });
    
    test('register creates user in Firebase and wallet', () async {
      // Arrange
      final email = 'test@example.com';
      final password = 'Password123!';
      final displayName = 'Test User';
      
      // Mock Firebase Auth
      when(mockFirebaseAuth.createUserWithEmailAndPassword(
        email: email,
        password: password,
      )).thenAnswer((_) async => mockUserCredential);
      
      when(mockUserCredential.user).thenReturn(mockFirebaseUser);
      when(mockFirebaseUser.uid).thenReturn('user123');
      when(mockFirebaseUser.email).thenReturn(email);
      
      // Mock update profile
      when(mockFirebaseUser.updateDisplayName(displayName))
          .thenAnswer((_) async {});
      
      // Mock Firestore
      final mockCollection = MockCollectionReference();
      final mockDocument = MockDocumentReference();
      
      when(mockFirestore.collection('users'))
          .thenReturn(mockCollection as CollectionReference<Map<String, dynamic>>);
      when(mockCollection.doc('user123'))
          .thenReturn(mockDocument as DocumentReference<Map<String, dynamic>>);
      when(mockDocument.set(any))
          .thenAnswer((_) async {});
      
      // Mock API Client
      final walletResponse = WalletResponse(
        address: '0xABC123',
        privateKey: '0x123',
        sessionToken: 'token123',
      );
      
      when(mockApiClient.createWallet())
          .thenAnswer((_) async => walletResponse);
      
      // Act
      final result = await authService.register(email, password, displayName);
      
      // Assert
      expect(result, isA<app.User>());
      expect(result.uid, equals('user123'));
      expect(result.email, equals(email));
      expect(result.displayName, equals(displayName));
      expect(result.walletAddress, equals('0xABC123'));
      
      // Verify interactions
      verify(mockFirebaseAuth.createUserWithEmailAndPassword(
        email: email,
        password: password,
      )).called(1);
      verify(mockFirebaseUser.updateDisplayName(displayName)).called(1);
      verify(mockApiClient.createWallet()).called(1);
      verify(mockDocument.set(any)).called(1);
    });
    
    test('login authenticates user and retrieves profile', () async {
      // Arrange
      final email = 'test@example.com';
      final password = 'Password123!';
      
      // Mock Firebase Auth
      when(mockFirebaseAuth.signInWithEmailAndPassword(
        email: email,
        password: password,
      )).thenAnswer((_) async => mockUserCredential);
      
      when(mockUserCredential.user).thenReturn(mockFirebaseUser);
      when(mockFirebaseUser.uid).thenReturn('user123');
      when(mockFirebaseUser.email).thenReturn(email);
      when(mockFirebaseUser.displayName).thenReturn('Test User');
      
      // Mock Firestore
      final mockCollection = MockCollectionReference();
      final mockDocument = MockDocumentReference();
      final mockSnapshot = MockDocumentSnapshot();
      
      when(mockFirestore.collection('users'))
          .thenReturn(mockCollection as CollectionReference<Map<String, dynamic>>);
      when(mockCollection.doc('user123'))
          .thenReturn(mockDocument as DocumentReference<Map<String, dynamic>>);
      when(mockDocument.get())
          .thenAnswer((_) async => mockSnapshot as DocumentSnapshot<Map<String, dynamic>>);
      when(mockSnapshot.data()).thenReturn({
        'uid': 'user123',
        'email': email,
        'displayName': 'Test User',
        'walletAddress': '0xABC123',
      });
      
      // Act
      final result = await authService.login(email, password);
      
      // Assert
      expect(result, isA<app.User>());
      expect(result.uid, equals('user123'));
      expect(result.email, equals(email));
      expect(result.displayName, equals('Test User'));
      expect(result.walletAddress, equals('0xABC123'));
      
      // Verify interactions
      verify(mockFirebaseAuth.signInWithEmailAndPassword(
        email: email,
        password: password,
      )).called(1);
      verify(mockDocument.get()).called(1);
    });
    
    test('getCurrentUser returns null when not logged in', () async {
      // Arrange
      when(mockFirebaseAuth.currentUser).thenReturn(null);
      
      // Act
      final result = await authService.getCurrentUser();
      
      // Assert
      expect(result, isNull);
    });
    
    test('logout signs out user', () async {
      // Arrange
      when(mockFirebaseAuth.signOut()).thenAnswer((_) async {});
      
      // Act
      await authService.logout();
      
      // Assert
      verify(mockFirebaseAuth.signOut()).called(1);
    });
  });
}

// Mock DocumentSnapshot for Firestore
class MockDocumentSnapshot extends Mock implements DocumentSnapshot {}