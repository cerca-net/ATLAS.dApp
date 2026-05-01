import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class RatingsRecord extends FirestoreRecord {
  RatingsRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "user_ref" field.
  DocumentReference? _userRef;
  DocumentReference? get userRef => _userRef;
  bool hasUserRef() => _userRef != null;

  // "value" field.
  int? _value;
  int get value => _value ?? 0;
  bool hasValue() => _value != null;

  // "comment" field.
  String? _comment;
  String get comment => _comment ?? '';
  bool hasComment() => _comment != null;

  // "date" field.
  DateTime? _date;
  DateTime? get date => _date;
  bool hasDate() => _date != null;

  DocumentReference get parentReference => reference.parent.parent!;

  void _initializeFields() {
    _userRef = safeDocRef(snapshotData['user_ref']);
    _value = castToType<int>(snapshotData['value']);
    _comment = snapshotData['comment'] as String?;
    _date = snapshotData['date'] is DateTime ? snapshotData['date'] as DateTime : null;
  }

  static Query<Map<String, dynamic>> collection([DocumentReference? parent]) =>
      parent != null
          ? parent.collection('ratings')
          : FirebaseFirestore.instance.collectionGroup('ratings');

  static DocumentReference createDoc(DocumentReference parent, {String? id}) =>
      parent.collection('ratings').doc(id);

  static Stream<RatingsRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => RatingsRecord.fromSnapshot(s));

  static Future<RatingsRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => RatingsRecord.fromSnapshot(s));

  static RatingsRecord fromSnapshot(DocumentSnapshot snapshot) =>
      RatingsRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static RatingsRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      RatingsRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'RatingsRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is RatingsRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createRatingsRecordData({
  DocumentReference? userRef,
  int? value,
  String? comment,
  DateTime? date,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'user_ref': userRef,
      'value': value,
      'comment': comment,
      'date': date,
    }.withoutNulls,
  );

  return firestoreData;
}

class RatingsRecordDocumentEquality implements Equality<RatingsRecord> {
  const RatingsRecordDocumentEquality();

  @override
  bool equals(RatingsRecord? e1, RatingsRecord? e2) {
    return e1?.userRef == e2?.userRef &&
        e1?.value == e2?.value &&
        e1?.comment == e2?.comment &&
        e1?.date == e2?.date;
  }

  @override
  int hash(RatingsRecord? e) =>
      const ListEquality().hash([e?.userRef, e?.value, e?.comment, e?.date]);

  @override
  bool isValidKey(Object? o) => o is RatingsRecord;
}
