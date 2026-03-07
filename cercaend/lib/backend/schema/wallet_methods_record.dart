import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class WalletMethodsRecord extends FirestoreRecord {
  WalletMethodsRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "method_poster" field.
  DocumentReference? _methodPoster;
  DocumentReference? get methodPoster => _methodPoster;
  bool hasMethodPoster() => _methodPoster != null;

  // "method_name" field.
  String? _methodName;
  String get methodName => _methodName ?? '';
  bool hasMethodName() => _methodName != null;

  // "method_type" field.
  String? _methodType;
  String get methodType => _methodType ?? '';
  bool hasMethodType() => _methodType != null;

  // "method_id" field.
  String? _methodId;
  String get methodId => _methodId ?? '';
  bool hasMethodId() => _methodId != null;

  // "method_account" field.
  String? _methodAccount;
  String get methodAccount => _methodAccount ?? '';
  bool hasMethodAccount() => _methodAccount != null;

  // "method_thread" field.
  List<String>? _methodThread;
  List<String> get methodThread => _methodThread ?? const [];
  bool hasMethodThread() => _methodThread != null;

  void _initializeFields() {
    _methodPoster = snapshotData['method_poster'] as DocumentReference?;
    _methodName = snapshotData['method_name'] as String?;
    _methodType = snapshotData['method_type'] as String?;
    _methodId = snapshotData['method_id'] as String?;
    _methodAccount = snapshotData['method_account'] as String?;
    _methodThread = getDataList(snapshotData['method_thread']);
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('wallet_methods');

  static Stream<WalletMethodsRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => WalletMethodsRecord.fromSnapshot(s));

  static Future<WalletMethodsRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => WalletMethodsRecord.fromSnapshot(s));

  static WalletMethodsRecord fromSnapshot(DocumentSnapshot snapshot) =>
      WalletMethodsRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static WalletMethodsRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      WalletMethodsRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'WalletMethodsRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is WalletMethodsRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createWalletMethodsRecordData({
  DocumentReference? methodPoster,
  String? methodName,
  String? methodType,
  String? methodId,
  String? methodAccount,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'method_poster': methodPoster,
      'method_name': methodName,
      'method_type': methodType,
      'method_id': methodId,
      'method_account': methodAccount,
    }.withoutNulls,
  );

  return firestoreData;
}

class WalletMethodsRecordDocumentEquality
    implements Equality<WalletMethodsRecord> {
  const WalletMethodsRecordDocumentEquality();

  @override
  bool equals(WalletMethodsRecord? e1, WalletMethodsRecord? e2) {
    const listEquality = ListEquality();
    return e1?.methodPoster == e2?.methodPoster &&
        e1?.methodName == e2?.methodName &&
        e1?.methodType == e2?.methodType &&
        e1?.methodId == e2?.methodId &&
        e1?.methodAccount == e2?.methodAccount &&
        listEquality.equals(e1?.methodThread, e2?.methodThread);
  }

  @override
  int hash(WalletMethodsRecord? e) => const ListEquality().hash([
        e?.methodPoster,
        e?.methodName,
        e?.methodType,
        e?.methodId,
        e?.methodAccount,
        e?.methodThread
      ]);

  @override
  bool isValidKey(Object? o) => o is WalletMethodsRecord;
}
