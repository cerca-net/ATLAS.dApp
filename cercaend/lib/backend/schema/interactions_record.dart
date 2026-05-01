import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';
import '/backend/schema/enums/enums.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class InteractionsRecord extends FirestoreRecord {
  InteractionsRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "user_ref" field.
  DocumentReference? _userRef;
  DocumentReference? get userRef => _userRef;
  bool hasUserRef() => _userRef != null;

  // "type" field.
  InteractionTypes? _type;
  InteractionTypes? get type => _type;
  bool hasType() => _type != null;

  // "value" field.
  double? _value;
  double get value => _value ?? 0.0;
  bool hasValue() => _value != null;

  void _initializeFields() {
    _userRef = safeDocRef(snapshotData['user_ref']);
    _type = snapshotData['type'] is InteractionTypes
        ? snapshotData['type']
        : deserializeEnum<InteractionTypes>(snapshotData['type']);
    _value = castToType<double>(snapshotData['value']);
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('interactions');

  static Stream<InteractionsRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => InteractionsRecord.fromSnapshot(s));

  static Future<InteractionsRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => InteractionsRecord.fromSnapshot(s));

  static InteractionsRecord fromSnapshot(DocumentSnapshot snapshot) =>
      InteractionsRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static InteractionsRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      InteractionsRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'InteractionsRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is InteractionsRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createInteractionsRecordData({
  DocumentReference? userRef,
  InteractionTypes? type,
  double? value,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'user_ref': userRef,
      'type': type,
      'value': value,
    }.withoutNulls,
  );

  return firestoreData;
}

class InteractionsRecordDocumentEquality
    implements Equality<InteractionsRecord> {
  const InteractionsRecordDocumentEquality();

  @override
  bool equals(InteractionsRecord? e1, InteractionsRecord? e2) {
    return e1?.userRef == e2?.userRef &&
        e1?.type == e2?.type &&
        e1?.value == e2?.value;
  }

  @override
  int hash(InteractionsRecord? e) =>
      const ListEquality().hash([e?.userRef, e?.type, e?.value]);

  @override
  bool isValidKey(Object? o) => o is InteractionsRecord;
}
