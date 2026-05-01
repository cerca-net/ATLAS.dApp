import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class ThreadRecord extends FirestoreRecord {
  ThreadRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "user_ref" field.
  DocumentReference? _userRef;
  DocumentReference? get userRef => _userRef;
  bool hasUserRef() => _userRef != null;

  // "public_ref" field.
  List<DocumentReference>? _publicRef;
  List<DocumentReference> get publicRef => _publicRef ?? const [];
  bool hasPublicRef() => _publicRef != null;

  // "submission_ref" field.
  DocumentReference? _submissionRef;
  DocumentReference? get submissionRef => _submissionRef;
  bool hasSubmissionRef() => _submissionRef != null;

  // "object_thread" field.
  List<String>? _objectThread;
  List<String> get objectThread => _objectThread ?? const [];
  bool hasObjectThread() => _objectThread != null;

  void _initializeFields() {
    _userRef = safeDocRef(snapshotData['user_ref']);
    _publicRef = getDataList(snapshotData['public_ref']);
    _submissionRef = safeDocRef(snapshotData['submission_ref']);
    _objectThread = getDataList(snapshotData['object_thread']);
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('thread');

  static Stream<ThreadRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => ThreadRecord.fromSnapshot(s));

  static Future<ThreadRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => ThreadRecord.fromSnapshot(s));

  static ThreadRecord fromSnapshot(DocumentSnapshot snapshot) => ThreadRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static ThreadRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      ThreadRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'ThreadRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is ThreadRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createThreadRecordData({
  DocumentReference? userRef,
  DocumentReference? submissionRef,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'user_ref': userRef,
      'submission_ref': submissionRef,
    }.withoutNulls,
  );

  return firestoreData;
}

class ThreadRecordDocumentEquality implements Equality<ThreadRecord> {
  const ThreadRecordDocumentEquality();

  @override
  bool equals(ThreadRecord? e1, ThreadRecord? e2) {
    const listEquality = ListEquality();
    return e1?.userRef == e2?.userRef &&
        listEquality.equals(e1?.publicRef, e2?.publicRef) &&
        e1?.submissionRef == e2?.submissionRef &&
        listEquality.equals(e1?.objectThread, e2?.objectThread);
  }

  @override
  int hash(ThreadRecord? e) => const ListEquality()
      .hash([e?.userRef, e?.publicRef, e?.submissionRef, e?.objectThread]);

  @override
  bool isValidKey(Object? o) => o is ThreadRecord;
}
