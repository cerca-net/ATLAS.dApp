import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class CreditsRecord extends FirestoreRecord {
  CreditsRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "userRef" field.
  DocumentReference? _userRef;
  DocumentReference? get userRef => _userRef;
  bool hasUserRef() => _userRef != null;

  // "f_x" field.
  double? _fX;
  double get fX => _fX ?? 0.0;
  bool hasFX() => _fX != null;

  // "generation_i" field.
  double? _generationI;
  double get generationI => _generationI ?? 0.0;
  bool hasGenerationI() => _generationI != null;

  // "participation_i" field.
  double? _participationI;
  double get participationI => _participationI ?? 0.0;
  bool hasParticipationI() => _participationI != null;

  // "transition_i" field.
  double? _transitionI;
  double get transitionI => _transitionI ?? 0.0;
  bool hasTransitionI() => _transitionI != null;

  void _initializeFields() {
    _userRef = snapshotData['userRef'] as DocumentReference?;
    _fX = castToType<double>(snapshotData['f_x']);
    _generationI = castToType<double>(snapshotData['generation_i']);
    _participationI = castToType<double>(snapshotData['participation_i']);
    _transitionI = castToType<double>(snapshotData['transition_i']);
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('credits');

  static Stream<CreditsRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => CreditsRecord.fromSnapshot(s));

  static Future<CreditsRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => CreditsRecord.fromSnapshot(s));

  static CreditsRecord fromSnapshot(DocumentSnapshot snapshot) =>
      CreditsRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static CreditsRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      CreditsRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'CreditsRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is CreditsRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createCreditsRecordData({
  DocumentReference? userRef,
  double? fX,
  double? generationI,
  double? participationI,
  double? transitionI,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'userRef': userRef,
      'f_x': fX,
      'generation_i': generationI,
      'participation_i': participationI,
      'transition_i': transitionI,
    }.withoutNulls,
  );

  return firestoreData;
}

class CreditsRecordDocumentEquality implements Equality<CreditsRecord> {
  const CreditsRecordDocumentEquality();

  @override
  bool equals(CreditsRecord? e1, CreditsRecord? e2) {
    return e1?.userRef == e2?.userRef &&
        e1?.fX == e2?.fX &&
        e1?.generationI == e2?.generationI &&
        e1?.participationI == e2?.participationI &&
        e1?.transitionI == e2?.transitionI;
  }

  @override
  int hash(CreditsRecord? e) => const ListEquality().hash(
      [e?.userRef, e?.fX, e?.generationI, e?.participationI, e?.transitionI]);

  @override
  bool isValidKey(Object? o) => o is CreditsRecord;
}
