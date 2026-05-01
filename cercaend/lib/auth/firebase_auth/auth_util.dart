import 'dart:async';
import 'package:supabase_flutter/supabase_flutter.dart';
import 'package:flutter/material.dart';
import '../../backend/supabase/supabase.dart';
import '../../backend/supabase/supabase_shim.dart';
import '../../backend/backend.dart';
import 'package:stream_transform/stream_transform.dart';
import '../supabase_auth/supabase_auth_manager.dart';
import '../../flutter_flow/flutter_flow_util.dart';

export '../supabase_auth/supabase_auth_manager.dart';
export '../supabase_auth/supabase_user_provider.dart';

final _authManager = SupabaseAuthManager();
SupabaseAuthManager get authManager => _authManager;

String get currentUserEmail => currentUser?.email ?? '';
String get currentUserUid => currentUser?.uid ?? '';
String get currentUserDisplayName => currentUser?.displayName ?? '';
String get currentUserPhoto => currentUser?.photoUrl ?? '';
String get currentPhoneNumber => currentUser?.phoneNumber ?? '';
String get currentJwtToken => _currentJwtToken ?? '';
bool get currentUserEmailVerified => currentUser?.emailVerified ?? false;

String? _currentJwtToken;
final jwtTokenStream = Supabase.instance.client.auth.onAuthStateChange
    .map((authState) async => _currentJwtToken = authState.session?.accessToken)
    .asBroadcastStream();

// Keeping currentUserReference signature but returning null for now as we transition off Firestore 
// (or mapping it to an empty mock object if needed).
// To prevent breaking UI that relies on this, we'll remove DocumentReference dependencies step-by-step.


DocumentReference? get currentUserReference {
  final uid = currentUserUid.isNotEmpty ? currentUserUid : 'mock_user_123';
  return DocumentReference('users', uid);
}

// The UI heavily relies on `UsersRecord`, we mock a stream wrapper until Phase 2
// For now, emit null so widgets fall back, or you can implement a Supabase wrapper here
// that returns a dummy `UsersRecord`.
UsersRecord? currentUserDocument;
final authenticatedUserStream = Supabase.instance.client.auth.onAuthStateChange
    .map<String>((authState) {
      debugPrint('Auth state changed! User ID: ${authState.session?.user.id}');
      return authState.session?.user.id ?? '';
    })
    .switchMap(
      (uid) {
        if (uid.isEmpty) return Stream.value(null);
        debugPrint('Querying users table for $uid...');
        return SupaFlow.client.from('users').stream(primaryKey: ['id']).eq('id', uid).map((rows) {
          debugPrint('Users table returned rows: ${rows.length}');
          final data = rows.isNotEmpty ? rows.first : <String, dynamic>{'id': uid, 'wallet_address': ''};
          return UsersRecord.fromSnapshot(DocumentSnapshot(uid, data, DocumentReference('users', uid)));
        }).handleError((e) {
          debugPrint('Error streaming users table: $e');
        });
      },
    )
    .map((userData) {
  debugPrint('Broadcasting currentUserDocument: $userData');
  currentUserDocument = userData as UsersRecord?;
  return currentUserDocument;
}).asBroadcastStream();

class AuthUserStreamWidget extends StatelessWidget {
  const AuthUserStreamWidget({super.key, required this.builder});

  final WidgetBuilder builder;

  @override
  Widget build(BuildContext context) => StreamBuilder(
        stream: authenticatedUserStream,
        builder: (context, _) => builder(context),
      );
}
