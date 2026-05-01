import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class SubmissionRecord extends FirestoreRecord {
  SubmissionRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "image" field.
  String? _image;
  String get image => _image ?? '';
  bool hasImage() => _image != null;

  // "date" field.
  DateTime? _date;
  DateTime? get date => _date;
  bool hasDate() => _date != null;

  // "header" field.
  String? _header;
  String get header => _header ?? '';
  bool hasHeader() => _header != null;

  // "poster" field.
  DocumentReference? _poster;
  DocumentReference? get poster => _poster;
  bool hasPoster() => _poster != null;

  // "imagesextra" field.
  List<String>? _imagesextra;
  List<String> get imagesextra => _imagesextra ?? const [];
  bool hasImagesextra() => _imagesextra != null;

  // "type1choice" field.
  String? _type1choice;
  String get type1choice => _type1choice ?? '';
  bool hasType1choice() => _type1choice != null;

  // "type2choice" field.
  String? _type2choice;
  String get type2choice => _type2choice ?? '';
  bool hasType2choice() => _type2choice != null;

  // "video" field.
  String? _video;
  String get video => _video ?? '';
  bool hasVideo() => _video != null;

  // "audio" field.
  String? _audio;
  String get audio => _audio ?? '';
  bool hasAudio() => _audio != null;

  // "body" field.
  String? _body;
  String get body => _body ?? '';
  bool hasBody() => _body != null;

  // "upvote" field.
  List<DocumentReference>? _upvote;
  List<DocumentReference> get upvote => _upvote ?? const [];
  bool hasUpvote() => _upvote != null;

  // "downvote" field.
  List<DocumentReference>? _downvote;
  List<DocumentReference> get downvote => _downvote ?? const [];
  bool hasDownvote() => _downvote != null;

  // "refvalue" field.
  String? _refvalue;
  String get refvalue => _refvalue ?? '';
  bool hasRefvalue() => _refvalue != null;

  // "type0choice" field.
  String? _type0choice;
  String get type0choice => _type0choice ?? '';
  bool hasType0choice() => _type0choice != null;

  // "submitted_date" field.
  DateTime? _submittedDate;
  DateTime? get submittedDate => _submittedDate;
  bool hasSubmittedDate() => _submittedDate != null;

  // "type_order" field.
  String? _typeOrder;
  String get typeOrder => _typeOrder ?? '';
  bool hasTypeOrder() => _typeOrder != null;

  // "type_activity" field.
  String? _typeActivity;
  String get typeActivity => _typeActivity ?? '';
  bool hasTypeActivity() => _typeActivity != null;

  // "type_ticket" field.
  String? _typeTicket;
  String get typeTicket => _typeTicket ?? '';
  bool hasTypeTicket() => _typeTicket != null;

  // "submitter" field.
  String? _submitter;
  String get submitter => _submitter ?? '';
  bool hasSubmitter() => _submitter != null;

  // "profilePictureObjectPost" field.
  String? _profilePictureObjectPost;
  String get profilePictureObjectPost => _profilePictureObjectPost ?? '';
  bool hasProfilePictureObjectPost() => _profilePictureObjectPost != null;

  // "ups" field.
  int? _ups;
  int get ups => _ups ?? 0;
  bool hasUps() => _ups != null;

  // "pins" field.
  List<DocumentReference>? _pins;
  List<DocumentReference> get pins => _pins ?? const [];
  bool hasPins() => _pins != null;

  // "object_is_edited" field.
  bool? _objectIsEdited;
  bool get objectIsEdited => _objectIsEdited ?? false;
  bool hasObjectIsEdited() => _objectIsEdited != null;

  // "shares" field.
  List<DocumentReference>? _shares;
  List<DocumentReference> get shares => _shares ?? const [];
  bool hasShares() => _shares != null;

  // "type0object" field.
  bool? _type0object;
  bool get type0object => _type0object ?? false;
  bool hasType0object() => _type0object != null;

  // "type_object" field.
  String? _typeObject;
  String get typeObject => _typeObject ?? '';
  bool hasTypeObject() => _typeObject != null;

  // "thread" field.
  List<String>? _thread;
  List<String> get thread => _thread ?? const [];
  bool hasThread() => _thread != null;

  // "object_thread" field.
  ThreadStruct? _objectThread;
  ThreadStruct get objectThread => _objectThread ?? ThreadStruct();
  bool hasObjectThread() => _objectThread != null;

  // "object_in_bag" field.
  List<DocumentReference>? _objectInBag;
  List<DocumentReference> get objectInBag => _objectInBag ?? const [];
  bool hasObjectInBag() => _objectInBag != null;

  // "object_refvalue" field.
  double? _objectRefvalue;
  double get objectRefvalue => _objectRefvalue ?? 0.0;
  bool hasObjectRefvalue() => _objectRefvalue != null;

  // "analytics_ref" field.
  DocumentReference? _analyticsRef;
  DocumentReference? get analyticsRef => _analyticsRef;
  bool hasAnalyticsRef() => _analyticsRef != null;

  // "object_views" field.
  int? _objectViews;
  int get objectViews => _objectViews ?? 0;
  bool hasObjectViews() => _objectViews != null;

  // "object_viewers" field.
  List<DocumentReference>? _objectViewers;
  List<DocumentReference> get objectViewers => _objectViewers ?? const [];
  bool hasObjectViewers() => _objectViewers != null;

  // "objectViewsS" field.
  List<DocumentReference>? _objectViewsS;
  List<DocumentReference> get objectViewsS => _objectViewsS ?? const [];
  bool hasObjectViewsS() => _objectViewsS != null;

  void _initializeFields() {
    _image = snapshotData['image'] as String?;
    _date = snapshotData['date'] is DateTime ? snapshotData['date'] as DateTime : null;
    _header = snapshotData['header'] as String?;
    _poster = safeDocRef(snapshotData['poster']);
    _imagesextra = getDataList(snapshotData['imagesextra']);
    _type1choice = snapshotData['type1choice'] as String?;
    _type2choice = snapshotData['type2choice'] as String?;
    _video = snapshotData['video'] as String?;
    _audio = snapshotData['audio'] as String?;
    _body = snapshotData['body'] as String?;
    _upvote = getDataList(snapshotData['upvote']);
    _downvote = getDataList(snapshotData['downvote']);
    _refvalue = snapshotData['refvalue'] as String?;
    _type0choice = snapshotData['type0choice'] as String?;
    _submittedDate = snapshotData['submitted_date'] is DateTime ? snapshotData['submitted_date'] as DateTime : null;
    _typeOrder = snapshotData['type_order'] as String?;
    _typeActivity = snapshotData['type_activity'] as String?;
    _typeTicket = snapshotData['type_ticket'] as String?;
    _submitter = snapshotData['submitter'] as String?;
    _profilePictureObjectPost =
        snapshotData['profilePictureObjectPost'] as String?;
    _ups = castToType<int>(snapshotData['ups']);
    _pins = getDataList(snapshotData['pins']);
    _objectIsEdited = snapshotData['object_is_edited'] is bool ? snapshotData['object_is_edited'] as bool : null;
    _shares = getDataList(snapshotData['shares']);
    _type0object = snapshotData['type0object'] is bool ? snapshotData['type0object'] as bool : null;
    _typeObject = snapshotData['type_object'] as String?;
    _thread = getDataList(snapshotData['thread']);
    _objectThread = snapshotData['object_thread'] is ThreadStruct
        ? snapshotData['object_thread']
        : ThreadStruct.maybeFromMap(snapshotData['object_thread']);
    _objectInBag = getDataList(snapshotData['object_in_bag']);
    _objectRefvalue = castToType<double>(snapshotData['object_refvalue']);
    _analyticsRef = safeDocRef(snapshotData['analytics_ref']);
    _objectViews = castToType<int>(snapshotData['object_views']);
    _objectViewers = getDataList(snapshotData['object_viewers']);
    _objectViewsS = getDataList(snapshotData['objectViewsS']);
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('submission');

  static Stream<SubmissionRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => SubmissionRecord.fromSnapshot(s));

  static Future<SubmissionRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => SubmissionRecord.fromSnapshot(s));

  static SubmissionRecord fromSnapshot(DocumentSnapshot snapshot) =>
      SubmissionRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static SubmissionRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      SubmissionRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'SubmissionRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is SubmissionRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createSubmissionRecordData({
  String? image,
  DateTime? date,
  String? header,
  DocumentReference? poster,
  String? type1choice,
  String? type2choice,
  String? video,
  String? audio,
  String? body,
  String? refvalue,
  String? type0choice,
  DateTime? submittedDate,
  String? typeOrder,
  String? typeActivity,
  String? typeTicket,
  String? submitter,
  String? profilePictureObjectPost,
  int? ups,
  bool? objectIsEdited,
  bool? type0object,
  String? typeObject,
  ThreadStruct? objectThread,
  double? objectRefvalue,
  DocumentReference? analyticsRef,
  int? objectViews,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'image': image,
      'date': date,
      'header': header,
      'poster': poster,
      'type1choice': type1choice,
      'type2choice': type2choice,
      'video': video,
      'audio': audio,
      'body': body,
      'refvalue': refvalue,
      'type0choice': type0choice,
      'submitted_date': submittedDate,
      'type_order': typeOrder,
      'type_activity': typeActivity,
      'type_ticket': typeTicket,
      'submitter': submitter,
      'profilePictureObjectPost': profilePictureObjectPost,
      'ups': ups,
      'object_is_edited': objectIsEdited,
      'type0object': type0object,
      'type_object': typeObject,
      'object_thread': ThreadStruct().toMap(),
      'object_refvalue': objectRefvalue,
      'analytics_ref': analyticsRef,
      'object_views': objectViews,
    }.withoutNulls,
  );

  // Handle nested data for "object_thread" field.
  addThreadStructData(firestoreData, objectThread, 'object_thread');

  return firestoreData;
}

class SubmissionRecordDocumentEquality implements Equality<SubmissionRecord> {
  const SubmissionRecordDocumentEquality();

  @override
  bool equals(SubmissionRecord? e1, SubmissionRecord? e2) {
    const listEquality = ListEquality();
    return e1?.image == e2?.image &&
        e1?.date == e2?.date &&
        e1?.header == e2?.header &&
        e1?.poster == e2?.poster &&
        listEquality.equals(e1?.imagesextra, e2?.imagesextra) &&
        e1?.type1choice == e2?.type1choice &&
        e1?.type2choice == e2?.type2choice &&
        e1?.video == e2?.video &&
        e1?.audio == e2?.audio &&
        e1?.body == e2?.body &&
        listEquality.equals(e1?.upvote, e2?.upvote) &&
        listEquality.equals(e1?.downvote, e2?.downvote) &&
        e1?.refvalue == e2?.refvalue &&
        e1?.type0choice == e2?.type0choice &&
        e1?.submittedDate == e2?.submittedDate &&
        e1?.typeOrder == e2?.typeOrder &&
        e1?.typeActivity == e2?.typeActivity &&
        e1?.typeTicket == e2?.typeTicket &&
        e1?.submitter == e2?.submitter &&
        e1?.profilePictureObjectPost == e2?.profilePictureObjectPost &&
        e1?.ups == e2?.ups &&
        listEquality.equals(e1?.pins, e2?.pins) &&
        e1?.objectIsEdited == e2?.objectIsEdited &&
        listEquality.equals(e1?.shares, e2?.shares) &&
        e1?.type0object == e2?.type0object &&
        e1?.typeObject == e2?.typeObject &&
        listEquality.equals(e1?.thread, e2?.thread) &&
        e1?.objectThread == e2?.objectThread &&
        listEquality.equals(e1?.objectInBag, e2?.objectInBag) &&
        e1?.objectRefvalue == e2?.objectRefvalue &&
        e1?.analyticsRef == e2?.analyticsRef &&
        e1?.objectViews == e2?.objectViews &&
        listEquality.equals(e1?.objectViewers, e2?.objectViewers) &&
        listEquality.equals(e1?.objectViewsS, e2?.objectViewsS);
  }

  @override
  int hash(SubmissionRecord? e) => const ListEquality().hash([
        e?.image,
        e?.date,
        e?.header,
        e?.poster,
        e?.imagesextra,
        e?.type1choice,
        e?.type2choice,
        e?.video,
        e?.audio,
        e?.body,
        e?.upvote,
        e?.downvote,
        e?.refvalue,
        e?.type0choice,
        e?.submittedDate,
        e?.typeOrder,
        e?.typeActivity,
        e?.typeTicket,
        e?.submitter,
        e?.profilePictureObjectPost,
        e?.ups,
        e?.pins,
        e?.objectIsEdited,
        e?.shares,
        e?.type0object,
        e?.typeObject,
        e?.thread,
        e?.objectThread,
        e?.objectInBag,
        e?.objectRefvalue,
        e?.analyticsRef,
        e?.objectViews,
        e?.objectViewers,
        e?.objectViewsS
      ]);

  @override
  bool isValidKey(Object? o) => o is SubmissionRecord;
}
