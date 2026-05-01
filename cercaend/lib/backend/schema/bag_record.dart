import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class BagRecord extends FirestoreRecord {
  BagRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "user_ref" field.
  DocumentReference? _userRef;
  DocumentReference? get userRef => _userRef;
  bool hasUserRef() => _userRef != null;

  // "items" field.
  List<BagItemStruct>? _items;
  List<BagItemStruct> get items => _items ?? const [];
  bool hasItems() => _items != null;

  // "in_bag_items" field.
  List<DocumentReference>? _inBagItems;
  List<DocumentReference> get inBagItems => _inBagItems ?? const [];
  bool hasInBagItems() => _inBagItems != null;

  // "public_user_ref" field.
  DocumentReference? _publicUserRef;
  DocumentReference? get publicUserRef => _publicUserRef;
  bool hasPublicUserRef() => _publicUserRef != null;

  void _initializeFields() {
    _userRef = safeDocRef(snapshotData['user_ref']);
    _items = getStructList(
      snapshotData['items'],
      BagItemStruct.fromMap,
    );
    _inBagItems = getDataList(snapshotData['in_bag_items']);
    _publicUserRef = safeDocRef(snapshotData['public_user_ref']);
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('bag');

  static Stream<BagRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => BagRecord.fromSnapshot(s));

  static Future<BagRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => BagRecord.fromSnapshot(s));

  static BagRecord fromSnapshot(DocumentSnapshot snapshot) => BagRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static BagRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      BagRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'BagRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is BagRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createBagRecordData({
  DocumentReference? userRef,
  DocumentReference? publicUserRef,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'user_ref': userRef,
      'public_user_ref': publicUserRef,
    }.withoutNulls,
  );

  return firestoreData;
}

class BagRecordDocumentEquality implements Equality<BagRecord> {
  const BagRecordDocumentEquality();

  @override
  bool equals(BagRecord? e1, BagRecord? e2) {
    const listEquality = ListEquality();
    return e1?.userRef == e2?.userRef &&
        listEquality.equals(e1?.items, e2?.items) &&
        listEquality.equals(e1?.inBagItems, e2?.inBagItems) &&
        e1?.publicUserRef == e2?.publicUserRef;
  }

  @override
  int hash(BagRecord? e) => const ListEquality()
      .hash([e?.userRef, e?.items, e?.inBagItems, e?.publicUserRef]);

  @override
  bool isValidKey(Object? o) => o is BagRecord;
}
