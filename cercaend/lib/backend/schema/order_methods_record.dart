import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class OrderMethodsRecord extends FirestoreRecord {
  OrderMethodsRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "method_poster" field.
  DocumentReference? _methodPoster;
  DocumentReference? get methodPoster => _methodPoster;
  bool hasMethodPoster() => _methodPoster != null;

  // "method_type" field.
  String? _methodType;
  String get methodType => _methodType ?? '';
  bool hasMethodType() => _methodType != null;

  // "method_tag" field.
  String? _methodTag;
  String get methodTag => _methodTag ?? '';
  bool hasMethodTag() => _methodTag != null;

  // "method_name" field.
  String? _methodName;
  String get methodName => _methodName ?? '';
  bool hasMethodName() => _methodName != null;

  // "method_thread" field.
  List<String>? _methodThread;
  List<String> get methodThread => _methodThread ?? const [];
  bool hasMethodThread() => _methodThread != null;

  void _initializeFields() {
    _methodPoster = safeDocRef(snapshotData['method_poster']);
    _methodType = snapshotData['method_type'] as String?;
    _methodTag = snapshotData['method_tag'] as String?;
    _methodName = snapshotData['method_name'] as String?;
    _methodThread = getDataList(snapshotData['method_thread']);
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('order_methods');

  static Stream<OrderMethodsRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => OrderMethodsRecord.fromSnapshot(s));

  static Future<OrderMethodsRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => OrderMethodsRecord.fromSnapshot(s));

  static OrderMethodsRecord fromSnapshot(DocumentSnapshot snapshot) =>
      OrderMethodsRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static OrderMethodsRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      OrderMethodsRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'OrderMethodsRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is OrderMethodsRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createOrderMethodsRecordData({
  DocumentReference? methodPoster,
  String? methodType,
  String? methodTag,
  String? methodName,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'method_poster': methodPoster,
      'method_type': methodType,
      'method_tag': methodTag,
      'method_name': methodName,
    }.withoutNulls,
  );

  return firestoreData;
}

class OrderMethodsRecordDocumentEquality
    implements Equality<OrderMethodsRecord> {
  const OrderMethodsRecordDocumentEquality();

  @override
  bool equals(OrderMethodsRecord? e1, OrderMethodsRecord? e2) {
    const listEquality = ListEquality();
    return e1?.methodPoster == e2?.methodPoster &&
        e1?.methodType == e2?.methodType &&
        e1?.methodTag == e2?.methodTag &&
        e1?.methodName == e2?.methodName &&
        listEquality.equals(e1?.methodThread, e2?.methodThread);
  }

  @override
  int hash(OrderMethodsRecord? e) => const ListEquality().hash([
        e?.methodPoster,
        e?.methodType,
        e?.methodTag,
        e?.methodName,
        e?.methodThread
      ]);

  @override
  bool isValidKey(Object? o) => o is OrderMethodsRecord;
}
