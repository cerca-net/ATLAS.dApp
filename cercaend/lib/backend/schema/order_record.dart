import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';
import '/backend/schema/enums/enums.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class OrderRecord extends FirestoreRecord {
  OrderRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "date" field.
  DateTime? _date;
  DateTime? get date => _date;
  bool hasDate() => _date != null;

  // "ref_value" field.
  double? _refValue;
  double get refValue => _refValue ?? 0.0;
  bool hasRefValue() => _refValue != null;

  // "total_ref_value" field.
  double? _totalRefValue;
  double get totalRefValue => _totalRefValue ?? 0.0;
  bool hasTotalRefValue() => _totalRefValue != null;

  // "order_completed" field.
  bool? _orderCompleted;
  bool get orderCompleted => _orderCompleted ?? false;
  bool hasOrderCompleted() => _orderCompleted != null;

  // "payment_method" field.
  String? _paymentMethod;
  String get paymentMethod => _paymentMethod ?? '';
  bool hasPaymentMethod() => _paymentMethod != null;

  // "order_comlpletion_date" field.
  DateTime? _orderComlpletionDate;
  DateTime? get orderComlpletionDate => _orderComlpletionDate;
  bool hasOrderComlpletionDate() => _orderComlpletionDate != null;

  // "user_ref" field.
  DocumentReference? _userRef;
  DocumentReference? get userRef => _userRef;
  bool hasUserRef() => _userRef != null;

  // "items" field.
  List<BagItemStruct>? _items;
  List<BagItemStruct> get items => _items ?? const [];
  bool hasItems() => _items != null;

  // "root_item" field.
  ItemStruct? _rootItem;
  ItemStruct get rootItem => _rootItem ?? ItemStruct();
  bool hasRootItem() => _rootItem != null;

  // "ordersCompleted" field.
  List<DocumentReference>? _ordersCompleted;
  List<DocumentReference> get ordersCompleted => _ordersCompleted ?? const [];
  bool hasOrdersCompleted() => _ordersCompleted != null;

  // "items_inorder" field.
  List<DocumentReference>? _itemsInorder;
  List<DocumentReference> get itemsInorder => _itemsInorder ?? const [];
  bool hasItemsInorder() => _itemsInorder != null;

  // "publicuser_ref" field.
  DocumentReference? _publicuserRef;
  DocumentReference? get publicuserRef => _publicuserRef;
  bool hasPublicuserRef() => _publicuserRef != null;

  // "order_stats" field.
  OrderStatuses? _orderStats;
  OrderStatuses? get orderStats => _orderStats;
  bool hasOrderStats() => _orderStats != null;

  // "order_method" field.
  DocumentReference? _orderMethod;
  DocumentReference? get orderMethod => _orderMethod;
  bool hasOrderMethod() => _orderMethod != null;

  // "wallet_method" field.
  DocumentReference? _walletMethod;
  DocumentReference? get walletMethod => _walletMethod;
  bool hasWalletMethod() => _walletMethod != null;

  // "username" field.
  String? _username;
  String get username => _username ?? '';
  bool hasUsername() => _username != null;

  // "publicusername" field.
  String? _publicusername;
  String get publicusername => _publicusername ?? '';
  bool hasPublicusername() => _publicusername != null;

  // "order_image" field.
  String? _orderImage;
  String get orderImage => _orderImage ?? '';
  bool hasOrderImage() => _orderImage != null;

  // "is_orderimage_uploaded" field.
  bool? _isOrderimageUploaded;
  bool get isOrderimageUploaded => _isOrderimageUploaded ?? false;
  bool hasIsOrderimageUploaded() => _isOrderimageUploaded != null;

  // "method_order" field.
  DocumentReference? _methodOrder;
  DocumentReference? get methodOrder => _methodOrder;
  bool hasMethodOrder() => _methodOrder != null;

  // "method_wallet" field.
  DocumentReference? _methodWallet;
  DocumentReference? get methodWallet => _methodWallet;
  bool hasMethodWallet() => _methodWallet != null;

  // "is_order_accepted" field.
  bool? _isOrderAccepted;
  bool get isOrderAccepted => _isOrderAccepted ?? false;
  bool hasIsOrderAccepted() => _isOrderAccepted != null;

  // "is_orderimage_accepted" field.
  bool? _isOrderimageAccepted;
  bool get isOrderimageAccepted => _isOrderimageAccepted ?? false;
  bool hasIsOrderimageAccepted() => _isOrderimageAccepted != null;

  // "key_1" field.
  String? _key1;
  String get key1 => _key1 ?? '';
  bool hasKey1() => _key1 != null;

  // "key_2" field.
  String? _key2;
  String get key2 => _key2 ?? '';
  bool hasKey2() => _key2 != null;

  // "rev_user" field.
  int? _revUser;
  int get revUser => _revUser ?? 0;
  bool hasRevUser() => _revUser != null;

  // "rev_pubuser" field.
  int? _revPubuser;
  int get revPubuser => _revPubuser ?? 0;
  bool hasRevPubuser() => _revPubuser != null;

  // "ref_list_orders" field.
  List<double>? _refListOrders;
  List<double> get refListOrders => _refListOrders ?? const [];
  bool hasRefListOrders() => _refListOrders != null;

  // "revby_maker" field.
  bool? _revbyMaker;
  bool get revbyMaker => _revbyMaker ?? false;
  bool hasRevbyMaker() => _revbyMaker != null;

  // "revby_taker" field.
  bool? _revbyTaker;
  bool get revbyTaker => _revbyTaker ?? false;
  bool hasRevbyTaker() => _revbyTaker != null;

  // "order_finished" field.
  bool? _orderFinished;
  bool get orderFinished => _orderFinished ?? false;
  bool hasOrderFinished() => _orderFinished != null;

  // "order_closed" field.
  bool? _orderClosed;
  bool get orderClosed => _orderClosed ?? false;
  bool hasOrderClosed() => _orderClosed != null;

  // "rating" field.
  List<int>? _rating;
  List<int> get rating => _rating ?? const [];
  bool hasRating() => _rating != null;

  void _initializeFields() {
    _date = snapshotData['date'] is DateTime ? snapshotData['date'] as DateTime : null;
    _refValue = castToType<double>(snapshotData['ref_value']);
    _totalRefValue = castToType<double>(snapshotData['total_ref_value']);
    _orderCompleted = snapshotData['order_completed'] is bool ? snapshotData['order_completed'] as bool : null;
    _paymentMethod = snapshotData['payment_method'] as String?;
    _orderComlpletionDate = snapshotData['order_comlpletion_date'] is DateTime ? snapshotData['order_comlpletion_date'] as DateTime : null;
    _userRef = safeDocRef(snapshotData['user_ref']);
    _items = getStructList(
      snapshotData['items'],
      BagItemStruct.fromMap,
    );
    _rootItem = snapshotData['root_item'] is ItemStruct
        ? snapshotData['root_item']
        : ItemStruct.maybeFromMap(snapshotData['root_item']);
    _ordersCompleted = getDataList(snapshotData['ordersCompleted']);
    _itemsInorder = getDataList(snapshotData['items_inorder']);
    _publicuserRef = safeDocRef(snapshotData['publicuser_ref']);
    _orderStats = snapshotData['order_stats'] is OrderStatuses
        ? snapshotData['order_stats']
        : deserializeEnum<OrderStatuses>(snapshotData['order_stats']);
    _orderMethod = safeDocRef(snapshotData['order_method']);
    _walletMethod = safeDocRef(snapshotData['wallet_method']);
    _username = snapshotData['username'] as String?;
    _publicusername = snapshotData['publicusername'] as String?;
    _orderImage = snapshotData['order_image'] as String?;
    _isOrderimageUploaded = snapshotData['is_orderimage_uploaded'] is bool ? snapshotData['is_orderimage_uploaded'] as bool : null;
    _methodOrder = safeDocRef(snapshotData['method_order']);
    _methodWallet = safeDocRef(snapshotData['method_wallet']);
    _isOrderAccepted = snapshotData['is_order_accepted'] is bool ? snapshotData['is_order_accepted'] as bool : null;
    _isOrderimageAccepted = snapshotData['is_orderimage_accepted'] is bool ? snapshotData['is_orderimage_accepted'] as bool : null;
    _key1 = snapshotData['key_1'] as String?;
    _key2 = snapshotData['key_2'] as String?;
    _revUser = castToType<int>(snapshotData['rev_user']);
    _revPubuser = castToType<int>(snapshotData['rev_pubuser']);
    _refListOrders = getDataList(snapshotData['ref_list_orders']);
    _revbyMaker = snapshotData['revby_maker'] is bool ? snapshotData['revby_maker'] as bool : null;
    _revbyTaker = snapshotData['revby_taker'] is bool ? snapshotData['revby_taker'] as bool : null;
    _orderFinished = snapshotData['order_finished'] is bool ? snapshotData['order_finished'] as bool : null;
    _orderClosed = snapshotData['order_closed'] is bool ? snapshotData['order_closed'] as bool : null;
    _rating = getDataList(snapshotData['rating']);
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('order');

  static Stream<OrderRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => OrderRecord.fromSnapshot(s));

  static Future<OrderRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => OrderRecord.fromSnapshot(s));

  static OrderRecord fromSnapshot(DocumentSnapshot snapshot) => OrderRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static OrderRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      OrderRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'OrderRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is OrderRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createOrderRecordData({
  DateTime? date,
  double? refValue,
  double? totalRefValue,
  bool? orderCompleted,
  String? paymentMethod,
  DateTime? orderComlpletionDate,
  DocumentReference? userRef,
  ItemStruct? rootItem,
  DocumentReference? publicuserRef,
  OrderStatuses? orderStats,
  DocumentReference? orderMethod,
  DocumentReference? walletMethod,
  String? username,
  String? publicusername,
  String? orderImage,
  bool? isOrderimageUploaded,
  DocumentReference? methodOrder,
  DocumentReference? methodWallet,
  bool? isOrderAccepted,
  bool? isOrderimageAccepted,
  String? key1,
  String? key2,
  int? revUser,
  int? revPubuser,
  bool? revbyMaker,
  bool? revbyTaker,
  bool? orderFinished,
  bool? orderClosed,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'date': date,
      'ref_value': refValue,
      'total_ref_value': totalRefValue,
      'order_completed': orderCompleted,
      'payment_method': paymentMethod,
      'order_comlpletion_date': orderComlpletionDate,
      'user_ref': userRef,
      'root_item': ItemStruct().toMap(),
      'publicuser_ref': publicuserRef,
      'order_stats': orderStats,
      'order_method': orderMethod,
      'wallet_method': walletMethod,
      'username': username,
      'publicusername': publicusername,
      'order_image': orderImage,
      'is_orderimage_uploaded': isOrderimageUploaded,
      'method_order': methodOrder,
      'method_wallet': methodWallet,
      'is_order_accepted': isOrderAccepted,
      'is_orderimage_accepted': isOrderimageAccepted,
      'key_1': key1,
      'key_2': key2,
      'rev_user': revUser,
      'rev_pubuser': revPubuser,
      'revby_maker': revbyMaker,
      'revby_taker': revbyTaker,
      'order_finished': orderFinished,
      'order_closed': orderClosed,
    }.withoutNulls,
  );

  // Handle nested data for "root_item" field.
  addItemStructData(firestoreData, rootItem, 'root_item');

  return firestoreData;
}

class OrderRecordDocumentEquality implements Equality<OrderRecord> {
  const OrderRecordDocumentEquality();

  @override
  bool equals(OrderRecord? e1, OrderRecord? e2) {
    const listEquality = ListEquality();
    return e1?.date == e2?.date &&
        e1?.refValue == e2?.refValue &&
        e1?.totalRefValue == e2?.totalRefValue &&
        e1?.orderCompleted == e2?.orderCompleted &&
        e1?.paymentMethod == e2?.paymentMethod &&
        e1?.orderComlpletionDate == e2?.orderComlpletionDate &&
        e1?.userRef == e2?.userRef &&
        listEquality.equals(e1?.items, e2?.items) &&
        e1?.rootItem == e2?.rootItem &&
        listEquality.equals(e1?.ordersCompleted, e2?.ordersCompleted) &&
        listEquality.equals(e1?.itemsInorder, e2?.itemsInorder) &&
        e1?.publicuserRef == e2?.publicuserRef &&
        e1?.orderStats == e2?.orderStats &&
        e1?.orderMethod == e2?.orderMethod &&
        e1?.walletMethod == e2?.walletMethod &&
        e1?.username == e2?.username &&
        e1?.publicusername == e2?.publicusername &&
        e1?.orderImage == e2?.orderImage &&
        e1?.isOrderimageUploaded == e2?.isOrderimageUploaded &&
        e1?.methodOrder == e2?.methodOrder &&
        e1?.methodWallet == e2?.methodWallet &&
        e1?.isOrderAccepted == e2?.isOrderAccepted &&
        e1?.isOrderimageAccepted == e2?.isOrderimageAccepted &&
        e1?.key1 == e2?.key1 &&
        e1?.key2 == e2?.key2 &&
        e1?.revUser == e2?.revUser &&
        e1?.revPubuser == e2?.revPubuser &&
        listEquality.equals(e1?.refListOrders, e2?.refListOrders) &&
        e1?.revbyMaker == e2?.revbyMaker &&
        e1?.revbyTaker == e2?.revbyTaker &&
        e1?.orderFinished == e2?.orderFinished &&
        e1?.orderClosed == e2?.orderClosed &&
        listEquality.equals(e1?.rating, e2?.rating);
  }

  @override
  int hash(OrderRecord? e) => const ListEquality().hash([
        e?.date,
        e?.refValue,
        e?.totalRefValue,
        e?.orderCompleted,
        e?.paymentMethod,
        e?.orderComlpletionDate,
        e?.userRef,
        e?.items,
        e?.rootItem,
        e?.ordersCompleted,
        e?.itemsInorder,
        e?.publicuserRef,
        e?.orderStats,
        e?.orderMethod,
        e?.walletMethod,
        e?.username,
        e?.publicusername,
        e?.orderImage,
        e?.isOrderimageUploaded,
        e?.methodOrder,
        e?.methodWallet,
        e?.isOrderAccepted,
        e?.isOrderimageAccepted,
        e?.key1,
        e?.key2,
        e?.revUser,
        e?.revPubuser,
        e?.refListOrders,
        e?.revbyMaker,
        e?.revbyTaker,
        e?.orderFinished,
        e?.orderClosed,
        e?.rating
      ]);

  @override
  bool isValidKey(Object? o) => o is OrderRecord;
}
