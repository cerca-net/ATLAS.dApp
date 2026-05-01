import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class CatalogueRecord extends FirestoreRecord {
  CatalogueRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "user_ref" field.
  DocumentReference? _userRef;
  DocumentReference? get userRef => _userRef;
  bool hasUserRef() => _userRef != null;

  // "catalalogue_name" field.
  String? _catalalogueName;
  String get catalalogueName => _catalalogueName ?? '';
  bool hasCatalalogueName() => _catalalogueName != null;

  // "catalogue_choicechips" field.
  List<String>? _catalogueChoicechips;
  List<String> get catalogueChoicechips => _catalogueChoicechips ?? const [];
  bool hasCatalogueChoicechips() => _catalogueChoicechips != null;

  // "catalogue_items" field.
  List<ItemStruct>? _catalogueItems;
  List<ItemStruct> get catalogueItems => _catalogueItems ?? const [];
  bool hasCatalogueItems() => _catalogueItems != null;

  // "total_items" field.
  int? _totalItems;
  int get totalItems => _totalItems ?? 0;
  bool hasTotalItems() => _totalItems != null;

  // "catalogue_buzzwords" field.
  List<String>? _catalogueBuzzwords;
  List<String> get catalogueBuzzwords => _catalogueBuzzwords ?? const [];
  bool hasCatalogueBuzzwords() => _catalogueBuzzwords != null;

  void _initializeFields() {
    _userRef = safeDocRef(snapshotData['user_ref']);
    _catalalogueName = snapshotData['catalalogue_name'] as String?;
    _catalogueChoicechips = getDataList(snapshotData['catalogue_choicechips']);
    _catalogueItems = getStructList(
      snapshotData['catalogue_items'],
      ItemStruct.fromMap,
    );
    _totalItems = castToType<int>(snapshotData['total_items']);
    _catalogueBuzzwords = getDataList(snapshotData['catalogue_buzzwords']);
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('catalogue');

  static Stream<CatalogueRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => CatalogueRecord.fromSnapshot(s));

  static Future<CatalogueRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => CatalogueRecord.fromSnapshot(s));

  static CatalogueRecord fromSnapshot(DocumentSnapshot snapshot) =>
      CatalogueRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static CatalogueRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      CatalogueRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'CatalogueRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is CatalogueRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createCatalogueRecordData({
  DocumentReference? userRef,
  String? catalalogueName,
  int? totalItems,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'user_ref': userRef,
      'catalalogue_name': catalalogueName,
      'total_items': totalItems,
    }.withoutNulls,
  );

  return firestoreData;
}

class CatalogueRecordDocumentEquality implements Equality<CatalogueRecord> {
  const CatalogueRecordDocumentEquality();

  @override
  bool equals(CatalogueRecord? e1, CatalogueRecord? e2) {
    const listEquality = ListEquality();
    return e1?.userRef == e2?.userRef &&
        e1?.catalalogueName == e2?.catalalogueName &&
        listEquality.equals(
            e1?.catalogueChoicechips, e2?.catalogueChoicechips) &&
        listEquality.equals(e1?.catalogueItems, e2?.catalogueItems) &&
        e1?.totalItems == e2?.totalItems &&
        listEquality.equals(e1?.catalogueBuzzwords, e2?.catalogueBuzzwords);
  }

  @override
  int hash(CatalogueRecord? e) => const ListEquality().hash([
        e?.userRef,
        e?.catalalogueName,
        e?.catalogueChoicechips,
        e?.catalogueItems,
        e?.totalItems,
        e?.catalogueBuzzwords
      ]);

  @override
  bool isValidKey(Object? o) => o is CatalogueRecord;
}
