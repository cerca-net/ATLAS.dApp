import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class AnalyticsRecord extends FirestoreRecord {
  AnalyticsRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "user_ref" field.
  DocumentReference? _userRef;
  DocumentReference? get userRef => _userRef;
  bool hasUserRef() => _userRef != null;

  // "user_objects" field.
  List<DocumentReference>? _userObjects;
  List<DocumentReference> get userObjects => _userObjects ?? const [];
  bool hasUserObjects() => _userObjects != null;

  // "user_orders" field.
  List<DocumentReference>? _userOrders;
  List<DocumentReference> get userOrders => _userOrders ?? const [];
  bool hasUserOrders() => _userOrders != null;

  // "object_views_accum" field.
  int? _objectViewsAccum;
  int get objectViewsAccum => _objectViewsAccum ?? 0;
  bool hasObjectViewsAccum() => _objectViewsAccum != null;

  // "user_object_count" field.
  int? _userObjectCount;
  int get userObjectCount => _userObjectCount ?? 0;
  bool hasUserObjectCount() => _userObjectCount != null;

  // "object_reach" field.
  int? _objectReach;
  int get objectReach => _objectReach ?? 0;
  bool hasObjectReach() => _objectReach != null;

  // "user_impressions" field.
  int? _userImpressions;
  int get userImpressions => _userImpressions ?? 0;
  bool hasUserImpressions() => _userImpressions != null;

  // "object_voterate" field.
  double? _objectVoterate;
  double get objectVoterate => _objectVoterate ?? 0.0;
  bool hasObjectVoterate() => _objectVoterate != null;

  // "object_pin" field.
  double? _objectPin;
  double get objectPin => _objectPin ?? 0.0;
  bool hasObjectPin() => _objectPin != null;

  // "object_share" field.
  double? _objectShare;
  double get objectShare => _objectShare ?? 0.0;
  bool hasObjectShare() => _objectShare != null;

  // "top_performer" field.
  String? _topPerformer;
  String get topPerformer => _topPerformer ?? '';
  bool hasTopPerformer() => _topPerformer != null;

  // "order_avg_ref" field.
  double? _orderAvgRef;
  double get orderAvgRef => _orderAvgRef ?? 0.0;
  bool hasOrderAvgRef() => _orderAvgRef != null;

  // "order_avg_time" field.
  int? _orderAvgTime;
  int get orderAvgTime => _orderAvgTime ?? 0;
  bool hasOrderAvgTime() => _orderAvgTime != null;

  // "order_avg_review" field.
  double? _orderAvgReview;
  double get orderAvgReview => _orderAvgReview ?? 0.0;
  bool hasOrderAvgReview() => _orderAvgReview != null;

  // "bag_order_ratio" field.
  double? _bagOrderRatio;
  double get bagOrderRatio => _bagOrderRatio ?? 0.0;
  bool hasBagOrderRatio() => _bagOrderRatio != null;

  // "order_accum_ref" field.
  double? _orderAccumRef;
  double get orderAccumRef => _orderAccumRef ?? 0.0;
  bool hasOrderAccumRef() => _orderAccumRef != null;

  // "order_accum" field.
  int? _orderAccum;
  int get orderAccum => _orderAccum ?? 0;
  bool hasOrderAccum() => _orderAccum != null;

  // "user_pin" field.
  int? _userPin;
  int get userPin => _userPin ?? 0;
  bool hasUserPin() => _userPin != null;

  // "market_index" field.
  int? _marketIndex;
  int get marketIndex => _marketIndex ?? 0;
  bool hasMarketIndex() => _marketIndex != null;

  // "performance" field.
  int? _performance;
  int get performance => _performance ?? 0;
  bool hasPerformance() => _performance != null;

  // "user_hash" field.
  double? _userHash;
  double get userHash => _userHash ?? 0.0;
  bool hasUserHash() => _userHash != null;

  // "objectAvgReference" field.
  List<double>? _objectAvgReference;
  List<double> get objectAvgReference => _objectAvgReference ?? const [];
  bool hasObjectAvgReference() => _objectAvgReference != null;

  // "orderAvgReference" field.
  List<double>? _orderAvgReference;
  List<double> get orderAvgReference => _orderAvgReference ?? const [];
  bool hasOrderAvgReference() => _orderAvgReference != null;

  // "order_pubrev_list" field.
  List<int>? _orderPubrevList;
  List<int> get orderPubrevList => _orderPubrevList ?? const [];
  bool hasOrderPubrevList() => _orderPubrevList != null;

  // "order_user_rev_list" field.
  List<int>? _orderUserRevList;
  List<int> get orderUserRevList => _orderUserRevList ?? const [];
  bool hasOrderUserRevList() => _orderUserRevList != null;

  // "pubusers_ref" field.
  List<DocumentReference>? _pubusersRef;
  List<DocumentReference> get pubusersRef => _pubusersRef ?? const [];
  bool hasPubusersRef() => _pubusersRef != null;

  // "orders_as_maker" field.
  List<DocumentReference>? _ordersAsMaker;
  List<DocumentReference> get ordersAsMaker => _ordersAsMaker ?? const [];
  bool hasOrdersAsMaker() => _ordersAsMaker != null;

  // "date_created" field.
  DateTime? _dateCreated;
  DateTime? get dateCreated => _dateCreated;
  bool hasDateCreated() => _dateCreated != null;

  // "ratings" field.
  List<int>? _ratings;
  List<int> get ratings => _ratings ?? const [];
  bool hasRatings() => _ratings != null;

  void _initializeFields() {
    _userRef = snapshotData['user_ref'] as DocumentReference?;
    _userObjects = getDataList(snapshotData['user_objects']);
    _userOrders = getDataList(snapshotData['user_orders']);
    _objectViewsAccum = castToType<int>(snapshotData['object_views_accum']);
    _userObjectCount = castToType<int>(snapshotData['user_object_count']);
    _objectReach = castToType<int>(snapshotData['object_reach']);
    _userImpressions = castToType<int>(snapshotData['user_impressions']);
    _objectVoterate = castToType<double>(snapshotData['object_voterate']);
    _objectPin = castToType<double>(snapshotData['object_pin']);
    _objectShare = castToType<double>(snapshotData['object_share']);
    _topPerformer = snapshotData['top_performer'] as String?;
    _orderAvgRef = castToType<double>(snapshotData['order_avg_ref']);
    _orderAvgTime = castToType<int>(snapshotData['order_avg_time']);
    _orderAvgReview = castToType<double>(snapshotData['order_avg_review']);
    _bagOrderRatio = castToType<double>(snapshotData['bag_order_ratio']);
    _orderAccumRef = castToType<double>(snapshotData['order_accum_ref']);
    _orderAccum = castToType<int>(snapshotData['order_accum']);
    _userPin = castToType<int>(snapshotData['user_pin']);
    _marketIndex = castToType<int>(snapshotData['market_index']);
    _performance = castToType<int>(snapshotData['performance']);
    _userHash = castToType<double>(snapshotData['user_hash']);
    _objectAvgReference = getDataList(snapshotData['objectAvgReference']);
    _orderAvgReference = getDataList(snapshotData['orderAvgReference']);
    _orderPubrevList = getDataList(snapshotData['order_pubrev_list']);
    _orderUserRevList = getDataList(snapshotData['order_user_rev_list']);
    _pubusersRef = getDataList(snapshotData['pubusers_ref']);
    _ordersAsMaker = getDataList(snapshotData['orders_as_maker']);
    _dateCreated = snapshotData['date_created'] as DateTime?;
    _ratings = getDataList(snapshotData['ratings']);
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('analytics');

  static Stream<AnalyticsRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => AnalyticsRecord.fromSnapshot(s));

  static Future<AnalyticsRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => AnalyticsRecord.fromSnapshot(s));

  static AnalyticsRecord fromSnapshot(DocumentSnapshot snapshot) =>
      AnalyticsRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static AnalyticsRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      AnalyticsRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'AnalyticsRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is AnalyticsRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createAnalyticsRecordData({
  DocumentReference? userRef,
  int? objectViewsAccum,
  int? userObjectCount,
  int? objectReach,
  int? userImpressions,
  double? objectVoterate,
  double? objectPin,
  double? objectShare,
  String? topPerformer,
  double? orderAvgRef,
  int? orderAvgTime,
  double? orderAvgReview,
  double? bagOrderRatio,
  double? orderAccumRef,
  int? orderAccum,
  int? userPin,
  int? marketIndex,
  int? performance,
  double? userHash,
  DateTime? dateCreated,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'user_ref': userRef,
      'object_views_accum': objectViewsAccum,
      'user_object_count': userObjectCount,
      'object_reach': objectReach,
      'user_impressions': userImpressions,
      'object_voterate': objectVoterate,
      'object_pin': objectPin,
      'object_share': objectShare,
      'top_performer': topPerformer,
      'order_avg_ref': orderAvgRef,
      'order_avg_time': orderAvgTime,
      'order_avg_review': orderAvgReview,
      'bag_order_ratio': bagOrderRatio,
      'order_accum_ref': orderAccumRef,
      'order_accum': orderAccum,
      'user_pin': userPin,
      'market_index': marketIndex,
      'performance': performance,
      'user_hash': userHash,
      'date_created': dateCreated,
    }.withoutNulls,
  );

  return firestoreData;
}

class AnalyticsRecordDocumentEquality implements Equality<AnalyticsRecord> {
  const AnalyticsRecordDocumentEquality();

  @override
  bool equals(AnalyticsRecord? e1, AnalyticsRecord? e2) {
    const listEquality = ListEquality();
    return e1?.userRef == e2?.userRef &&
        listEquality.equals(e1?.userObjects, e2?.userObjects) &&
        listEquality.equals(e1?.userOrders, e2?.userOrders) &&
        e1?.objectViewsAccum == e2?.objectViewsAccum &&
        e1?.userObjectCount == e2?.userObjectCount &&
        e1?.objectReach == e2?.objectReach &&
        e1?.userImpressions == e2?.userImpressions &&
        e1?.objectVoterate == e2?.objectVoterate &&
        e1?.objectPin == e2?.objectPin &&
        e1?.objectShare == e2?.objectShare &&
        e1?.topPerformer == e2?.topPerformer &&
        e1?.orderAvgRef == e2?.orderAvgRef &&
        e1?.orderAvgTime == e2?.orderAvgTime &&
        e1?.orderAvgReview == e2?.orderAvgReview &&
        e1?.bagOrderRatio == e2?.bagOrderRatio &&
        e1?.orderAccumRef == e2?.orderAccumRef &&
        e1?.orderAccum == e2?.orderAccum &&
        e1?.userPin == e2?.userPin &&
        e1?.marketIndex == e2?.marketIndex &&
        e1?.performance == e2?.performance &&
        e1?.userHash == e2?.userHash &&
        listEquality.equals(e1?.objectAvgReference, e2?.objectAvgReference) &&
        listEquality.equals(e1?.orderAvgReference, e2?.orderAvgReference) &&
        listEquality.equals(e1?.orderPubrevList, e2?.orderPubrevList) &&
        listEquality.equals(e1?.orderUserRevList, e2?.orderUserRevList) &&
        listEquality.equals(e1?.pubusersRef, e2?.pubusersRef) &&
        listEquality.equals(e1?.ordersAsMaker, e2?.ordersAsMaker) &&
        e1?.dateCreated == e2?.dateCreated &&
        listEquality.equals(e1?.ratings, e2?.ratings);
  }

  @override
  int hash(AnalyticsRecord? e) => const ListEquality().hash([
        e?.userRef,
        e?.userObjects,
        e?.userOrders,
        e?.objectViewsAccum,
        e?.userObjectCount,
        e?.objectReach,
        e?.userImpressions,
        e?.objectVoterate,
        e?.objectPin,
        e?.objectShare,
        e?.topPerformer,
        e?.orderAvgRef,
        e?.orderAvgTime,
        e?.orderAvgReview,
        e?.bagOrderRatio,
        e?.orderAccumRef,
        e?.orderAccum,
        e?.userPin,
        e?.marketIndex,
        e?.performance,
        e?.userHash,
        e?.objectAvgReference,
        e?.orderAvgReference,
        e?.orderPubrevList,
        e?.orderUserRevList,
        e?.pubusersRef,
        e?.ordersAsMaker,
        e?.dateCreated,
        e?.ratings
      ]);

  @override
  bool isValidKey(Object? o) => o is AnalyticsRecord;
}
