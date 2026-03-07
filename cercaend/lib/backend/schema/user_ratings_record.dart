import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class UserRatingsRecord extends FirestoreRecord {
  UserRatingsRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "rated_user" field.
  DocumentReference? _ratedUser;
  DocumentReference? get ratedUser => _ratedUser;
  bool hasRatedUser() => _ratedUser != null;

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

  void _initializeFields() {
    _ratedUser = snapshotData['rated_user'] as DocumentReference?;
    _value = castToType<int>(snapshotData['value']);
    _comment = snapshotData['comment'] as String?;
    _date = snapshotData['date'] as DateTime?;
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('user_ratings');

  static Stream<UserRatingsRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => UserRatingsRecord.fromSnapshot(s));

  static Future<UserRatingsRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => UserRatingsRecord.fromSnapshot(s));

  static UserRatingsRecord fromSnapshot(DocumentSnapshot snapshot) =>
      UserRatingsRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static UserRatingsRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      UserRatingsRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'UserRatingsRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is UserRatingsRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createUserRatingsRecordData({
  DocumentReference? ratedUser,
  int? value,
  String? comment,
  DateTime? date,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'rated_user': ratedUser,
      'value': value,
      'comment': comment,
      'date': date,
    }.withoutNulls,
  );

  return firestoreData;
}

class UserRatingsRecordDocumentEquality implements Equality<UserRatingsRecord> {
  const UserRatingsRecordDocumentEquality();

  @override
  bool equals(UserRatingsRecord? e1, UserRatingsRecord? e2) {
    return e1?.ratedUser == e2?.ratedUser &&
        e1?.value == e2?.value &&
        e1?.comment == e2?.comment &&
        e1?.date == e2?.date;
  }

  @override
  int hash(UserRatingsRecord? e) =>
      const ListEquality().hash([e?.ratedUser, e?.value, e?.comment, e?.date]);

  @override
  bool isValidKey(Object? o) => o is UserRatingsRecord;
}
