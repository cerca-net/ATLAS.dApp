import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class TransactionsRecord extends FirestoreRecord {
  TransactionsRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "user_madeby" field.
  DocumentReference? _userMadeby;
  DocumentReference? get userMadeby => _userMadeby;
  bool hasUserMadeby() => _userMadeby != null;

  // "user_madeto" field.
  DocumentReference? _userMadeto;
  DocumentReference? get userMadeto => _userMadeto;
  bool hasUserMadeto() => _userMadeto != null;

  // "items" field.
  List<DocumentReference>? _items;
  List<DocumentReference> get items => _items ?? const [];
  bool hasItems() => _items != null;

  // "id_transaction" field.
  String? _idTransaction;
  String get idTransaction => _idTransaction ?? '';
  bool hasIdTransaction() => _idTransaction != null;

  // "date" field.
  DateTime? _date;
  DateTime? get date => _date;
  bool hasDate() => _date != null;

  // "total_ref_value_transaction" field.
  double? _totalRefValueTransaction;
  double get totalRefValueTransaction => _totalRefValueTransaction ?? 0.0;
  bool hasTotalRefValueTransaction() => _totalRefValueTransaction != null;

  // "order_ref" field.
  DocumentReference? _orderRef;
  DocumentReference? get orderRef => _orderRef;
  bool hasOrderRef() => _orderRef != null;

  // "w_method" field.
  DocumentReference? _wMethod;
  DocumentReference? get wMethod => _wMethod;
  bool hasWMethod() => _wMethod != null;

  // "o_method" field.
  DocumentReference? _oMethod;
  DocumentReference? get oMethod => _oMethod;
  bool hasOMethod() => _oMethod != null;

  // "UserIn" field.
  String? _userIn;
  String get userIn => _userIn ?? '';
  bool hasUserIn() => _userIn != null;

  // "userOut" field.
  String? _userOut;
  String get userOut => _userOut ?? '';
  bool hasUserOut() => _userOut != null;

  void _initializeFields() {
    _userMadeby = snapshotData['user_madeby'] as DocumentReference?;
    _userMadeto = snapshotData['user_madeto'] as DocumentReference?;
    _items = getDataList(snapshotData['items']);
    _idTransaction = snapshotData['id_transaction'] as String?;
    _date = snapshotData['date'] as DateTime?;
    _totalRefValueTransaction =
        castToType<double>(snapshotData['total_ref_value_transaction']);
    _orderRef = snapshotData['order_ref'] as DocumentReference?;
    _wMethod = snapshotData['w_method'] as DocumentReference?;
    _oMethod = snapshotData['o_method'] as DocumentReference?;
    _userIn = snapshotData['UserIn'] as String?;
    _userOut = snapshotData['userOut'] as String?;
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('transactions');

  static Stream<TransactionsRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => TransactionsRecord.fromSnapshot(s));

  static Future<TransactionsRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => TransactionsRecord.fromSnapshot(s));

  static TransactionsRecord fromSnapshot(DocumentSnapshot snapshot) =>
      TransactionsRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static TransactionsRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      TransactionsRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'TransactionsRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is TransactionsRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createTransactionsRecordData({
  DocumentReference? userMadeby,
  DocumentReference? userMadeto,
  String? idTransaction,
  DateTime? date,
  double? totalRefValueTransaction,
  DocumentReference? orderRef,
  DocumentReference? wMethod,
  DocumentReference? oMethod,
  String? userIn,
  String? userOut,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'user_madeby': userMadeby,
      'user_madeto': userMadeto,
      'id_transaction': idTransaction,
      'date': date,
      'total_ref_value_transaction': totalRefValueTransaction,
      'order_ref': orderRef,
      'w_method': wMethod,
      'o_method': oMethod,
      'UserIn': userIn,
      'userOut': userOut,
    }.withoutNulls,
  );

  return firestoreData;
}

class TransactionsRecordDocumentEquality
    implements Equality<TransactionsRecord> {
  const TransactionsRecordDocumentEquality();

  @override
  bool equals(TransactionsRecord? e1, TransactionsRecord? e2) {
    const listEquality = ListEquality();
    return e1?.userMadeby == e2?.userMadeby &&
        e1?.userMadeto == e2?.userMadeto &&
        listEquality.equals(e1?.items, e2?.items) &&
        e1?.idTransaction == e2?.idTransaction &&
        e1?.date == e2?.date &&
        e1?.totalRefValueTransaction == e2?.totalRefValueTransaction &&
        e1?.orderRef == e2?.orderRef &&
        e1?.wMethod == e2?.wMethod &&
        e1?.oMethod == e2?.oMethod &&
        e1?.userIn == e2?.userIn &&
        e1?.userOut == e2?.userOut;
  }

  @override
  int hash(TransactionsRecord? e) => const ListEquality().hash([
        e?.userMadeby,
        e?.userMadeto,
        e?.items,
        e?.idTransaction,
        e?.date,
        e?.totalRefValueTransaction,
        e?.orderRef,
        e?.wMethod,
        e?.oMethod,
        e?.userIn,
        e?.userOut
      ]);

  @override
  bool isValidKey(Object? o) => o is TransactionsRecord;
}
