import 'dart:async';

import 'package:collection/collection.dart';

import '/backend/schema/util/firestore_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class UsersRecord extends FirestoreRecord {
  UsersRecord._(
    super.reference,
    super.data,
  ) {
    _initializeFields();
  }

  // "bio" field.
  String? _bio;
  String get bio => _bio ?? '';
  bool hasBio() => _bio != null;

  // "email" field.
  String? _email;
  String get email => _email ?? '';
  bool hasEmail() => _email != null;

  // "location" field.
  String? _location;
  String get location => _location ?? '';
  bool hasLocation() => _location != null;

  // "name" field.
  String? _name;
  String get name => _name ?? '';
  bool hasName() => _name != null;

  // "surname" field.
  String? _surname;
  String get surname => _surname ?? '';
  bool hasSurname() => _surname != null;

  // "username" field.
  String? _username;
  String get username => _username ?? '';
  bool hasUsername() => _username != null;

  // "display_name" field.
  String? _displayName;
  String get displayName => _displayName ?? '';
  bool hasDisplayName() => _displayName != null;

  // "photo_url" field.
  String? _photoUrl;
  String get photoUrl => _photoUrl ?? '';
  bool hasPhotoUrl() => _photoUrl != null;

  // "uid" field.
  String? _uid;
  String get uid => _uid ?? '';
  bool hasUid() => _uid != null;

  // "created_time" field.
  DateTime? _createdTime;
  DateTime? get createdTime => _createdTime;
  bool hasCreatedTime() => _createdTime != null;

  // "phone_number" field.
  String? _phoneNumber;
  String get phoneNumber => _phoneNumber ?? '';
  bool hasPhoneNumber() => _phoneNumber != null;

  // "last_active_time" field.
  DateTime? _lastActiveTime;
  DateTime? get lastActiveTime => _lastActiveTime;
  bool hasLastActiveTime() => _lastActiveTime != null;

  // "role" field.
  String? _role;
  String get role => _role ?? '';
  bool hasRole() => _role != null;

  // "title" field.
  String? _title;
  String get title => _title ?? '';
  bool hasTitle() => _title != null;

  // "banner" field.
  String? _banner;
  String get banner => _banner ?? '';
  bool hasBanner() => _banner != null;

  // "profile_picture" field.
  String? _profilePicture;
  String get profilePicture => _profilePicture ?? '';
  bool hasProfilePicture() => _profilePicture != null;

  // "usernames" field.
  String? _usernames;
  String get usernames => _usernames ?? '';
  bool hasUsernames() => _usernames != null;

  // "pinned_users" field.
  List<DocumentReference>? _pinnedUsers;
  List<DocumentReference> get pinnedUsers => _pinnedUsers ?? const [];
  bool hasPinnedUsers() => _pinnedUsers != null;

  // "user_pins" field.
  List<DocumentReference>? _userPins;
  List<DocumentReference> get userPins => _userPins ?? const [];
  bool hasUserPins() => _userPins != null;

  // "bag_ref" field.
  DocumentReference? _bagRef;
  DocumentReference? get bagRef => _bagRef;
  bool hasBagRef() => _bagRef != null;

  // "user_occupations" field.
  List<String>? _userOccupations;
  List<String> get userOccupations => _userOccupations ?? const [];
  bool hasUserOccupations() => _userOccupations != null;

  // "user_interests" field.
  List<String>? _userInterests;
  List<String> get userInterests => _userInterests ?? const [];
  bool hasUserInterests() => _userInterests != null;

  // "pinned_objects" field.
  List<DocumentReference>? _pinnedObjects;
  List<DocumentReference> get pinnedObjects => _pinnedObjects ?? const [];
  bool hasPinnedObjects() => _pinnedObjects != null;

  // "user_objects" field.
  List<DocumentReference>? _userObjects;
  List<DocumentReference> get userObjects => _userObjects ?? const [];
  bool hasUserObjects() => _userObjects != null;

  // "wallet_methods_user" field.
  List<DocumentReference>? _walletMethodsUser;
  List<DocumentReference> get walletMethodsUser =>
      _walletMethodsUser ?? const [];
  bool hasWalletMethodsUser() => _walletMethodsUser != null;

  // "wallet_address" field.
  String? _walletAddress;
  String get walletAddress => _walletAddress ?? '';
  bool hasWalletAddress() => _walletAddress != null;

  // "user_bag_objects" field.
  List<DocumentReference>? _userBagObjects;
  List<DocumentReference> get userBagObjects => _userBagObjects ?? const [];
  bool hasUserBagObjects() => _userBagObjects != null;

  // "order_ref" field.
  DocumentReference? _orderRef;
  DocumentReference? get orderRef => _orderRef;
  bool hasOrderRef() => _orderRef != null;

  // "order_methods_user" field.
  List<DocumentReference>? _orderMethodsUser;
  List<DocumentReference> get orderMethodsUser => _orderMethodsUser ?? const [];
  bool hasOrderMethodsUser() => _orderMethodsUser != null;

  // "user_places" field.
  List<String>? _userPlaces;
  List<String> get userPlaces => _userPlaces ?? const [];
  bool hasUserPlaces() => _userPlaces != null;

  // "user_type" field.
  String? _userType;
  String get userType => _userType ?? '';
  bool hasUserType() => _userType != null;

  // "user_items" field.
  List<DocumentReference>? _userItems;
  List<DocumentReference> get userItems => _userItems ?? const [];
  bool hasUserItems() => _userItems != null;

  // "user_posts" field.
  List<DocumentReference>? _userPosts;
  List<DocumentReference> get userPosts => _userPosts ?? const [];
  bool hasUserPosts() => _userPosts != null;

  // "user_transactions" field.
  List<DocumentReference>? _userTransactions;
  List<DocumentReference> get userTransactions => _userTransactions ?? const [];
  bool hasUserTransactions() => _userTransactions != null;

  // "shortDescription" field.
  String? _shortDescription;
  String get shortDescription => _shortDescription ?? '';
  bool hasShortDescription() => _shortDescription != null;

  // "analytics_ref" field.
  DocumentReference? _analyticsRef;
  DocumentReference? get analyticsRef => _analyticsRef;
  bool hasAnalyticsRef() => _analyticsRef != null;

  // "user_verified" field.
  bool? _userVerified;
  bool get userVerified => _userVerified ?? false;
  bool hasUserVerified() => _userVerified != null;

  // "user_verified_pending" field.
  bool? _userVerifiedPending;
  bool get userVerifiedPending => _userVerifiedPending ?? false;
  bool hasUserVerifiedPending() => _userVerifiedPending != null;

  // "user_is_admin" field.
  bool? _userIsAdmin;
  bool get userIsAdmin => _userIsAdmin ?? false;
  bool hasUserIsAdmin() => _userIsAdmin != null;

  // "u_ratings" field.
  List<DocumentReference>? _uRatings;
  List<DocumentReference> get uRatings => _uRatings ?? const [];
  bool hasURatings() => _uRatings != null;

  // "avg_u_rating" field.
  double? _avgURating;
  double get avgURating => _avgURating ?? 0.0;
  bool hasAvgURating() => _avgURating != null;

  // "ratings" field.
  int? _ratings;
  int get ratings => _ratings ?? 0;
  bool hasRatings() => _ratings != null;

  // "user_ratings" field.
  List<DocumentReference>? _userRatings;
  List<DocumentReference> get userRatings => _userRatings ?? const [];
  bool hasUserRatings() => _userRatings != null;

  // "userCredits" field.
  DocumentReference? _userCredits;
  DocumentReference? get userCredits => _userCredits;
  bool hasUserCredits() => _userCredits != null;

  // "following_users" field.
  List<DocumentReference>? _followingUsers;
  List<DocumentReference> get followingUsers => _followingUsers ?? const [];
  bool hasFollowingUsers() => _followingUsers != null;

  // "users_follwoing_me" field.
  List<DocumentReference>? _usersFollwoingMe;
  List<DocumentReference> get usersFollwoingMe => _usersFollwoingMe ?? const [];
  bool hasUsersFollwoingMe() => _usersFollwoingMe != null;

  // "verification_id_url" field.
  String? _verificationIdUrl;
  String get verificationIdUrl => _verificationIdUrl ?? '';
  bool hasVerificationIdUrl() => _verificationIdUrl != null;

  // "verification_selfie_url" field.
  String? _verificationSelfieUrl;
  String get verificationSelfieUrl => _verificationSelfieUrl ?? '';
  bool hasVerificationSelfieUrl() => _verificationSelfieUrl != null;

  // "verification_location_url" field.
  String? _verificationLocationUrl;
  String get verificationLocationUrl => _verificationLocationUrl ?? '';
  bool hasVerificationLocationUrl() => _verificationLocationUrl != null;

  // "verification_municipality" field.
  String? _verificationMunicipality;
  String get verificationMunicipality => _verificationMunicipality ?? '';
  bool hasVerificationMunicipality() => _verificationMunicipality != null;

  // "verification_submitted_at" field.
  DateTime? _verificationSubmittedAt;
  DateTime? get verificationSubmittedAt => _verificationSubmittedAt;
  bool hasVerificationSubmittedAt() => _verificationSubmittedAt != null;

  void _initializeFields() {
    _bio = snapshotData['bio'] as String?;
    _email = snapshotData['email'] as String?;
    _location = snapshotData['location'] as String?;
    _name = snapshotData['name'] as String?;
    _surname = snapshotData['surname'] as String?;
    _username = snapshotData['username'] as String?;
    _displayName = snapshotData['display_name'] as String?;
    _photoUrl = snapshotData['photo_url'] as String?;
    _uid = snapshotData['uid'] as String?;
    _createdTime = snapshotData['created_time'] as DateTime?;
    _phoneNumber = snapshotData['phone_number'] as String?;
    _lastActiveTime = snapshotData['last_active_time'] as DateTime?;
    _role = snapshotData['role'] as String?;
    _title = snapshotData['title'] as String?;
    _banner = snapshotData['banner'] as String?;
    _profilePicture = snapshotData['profile_picture'] as String?;
    _usernames = snapshotData['usernames'] as String?;
    _pinnedUsers = getDataList(snapshotData['pinned_users']);
    _userPins = getDataList(snapshotData['user_pins']);
    _bagRef = snapshotData['bag_ref'] as DocumentReference?;
    _userOccupations = getDataList(snapshotData['user_occupations']);
    _userInterests = getDataList(snapshotData['user_interests']);
    _pinnedObjects = getDataList(snapshotData['pinned_objects']);
    _userObjects = getDataList(snapshotData['user_objects']);
    _walletMethodsUser = getDataList(snapshotData['wallet_methods_user']);
    _walletAddress = snapshotData['wallet_address'] as String?;
    _userBagObjects = getDataList(snapshotData['user_bag_objects']);
    _orderRef = snapshotData['order_ref'] as DocumentReference?;
    _orderMethodsUser = getDataList(snapshotData['order_methods_user']);
    _userPlaces = getDataList(snapshotData['user_places']);
    _userType = snapshotData['user_type'] as String?;
    _userItems = getDataList(snapshotData['user_items']);
    _userPosts = getDataList(snapshotData['user_posts']);
    _userTransactions = getDataList(snapshotData['user_transactions']);
    _shortDescription = snapshotData['shortDescription'] as String?;
    _analyticsRef = snapshotData['analytics_ref'] as DocumentReference?;
    _userVerified = snapshotData['user_verified'] as bool?;
    _userVerifiedPending = snapshotData['user_verified_pending'] as bool?;
    _userIsAdmin = snapshotData['user_is_admin'] as bool?;
    _uRatings = getDataList(snapshotData['u_ratings']);
    _avgURating = castToType<double>(snapshotData['avg_u_rating']);
    _ratings = castToType<int>(snapshotData['ratings']);
    _userRatings = getDataList(snapshotData['user_ratings']);
    _userCredits = snapshotData['userCredits'] as DocumentReference?;
    _followingUsers = getDataList(snapshotData['following_users']);
    _usersFollwoingMe = getDataList(snapshotData['users_follwoing_me']);
    _verificationIdUrl = snapshotData['verification_id_url'] as String?;
    _verificationSelfieUrl = snapshotData['verification_selfie_url'] as String?;
    _verificationLocationUrl =
        snapshotData['verification_location_url'] as String?;
    _verificationMunicipality =
        snapshotData['verification_municipality'] as String?;
    _verificationSubmittedAt =
        snapshotData['verification_submitted_at'] as DateTime?;
  }

  static CollectionReference get collection =>
      FirebaseFirestore.instance.collection('users');

  static Stream<UsersRecord> getDocument(DocumentReference ref) =>
      ref.snapshots().map((s) => UsersRecord.fromSnapshot(s));

  static Future<UsersRecord> getDocumentOnce(DocumentReference ref) =>
      ref.get().then((s) => UsersRecord.fromSnapshot(s));

  static UsersRecord fromSnapshot(DocumentSnapshot snapshot) => UsersRecord._(
        snapshot.reference,
        mapFromFirestore(snapshot.data() as Map<String, dynamic>),
      );

  static UsersRecord getDocumentFromData(
    Map<String, dynamic> data,
    DocumentReference reference,
  ) =>
      UsersRecord._(reference, mapFromFirestore(data));

  @override
  String toString() =>
      'UsersRecord(reference: ${reference.path}, data: $snapshotData)';

  @override
  int get hashCode => reference.path.hashCode;

  @override
  bool operator ==(other) =>
      other is UsersRecord &&
      reference.path.hashCode == other.reference.path.hashCode;
}

Map<String, dynamic> createUsersRecordData({
  String? bio,
  String? email,
  String? location,
  String? name,
  String? surname,
  String? username,
  String? displayName,
  String? photoUrl,
  String? uid,
  DateTime? createdTime,
  String? phoneNumber,
  DateTime? lastActiveTime,
  String? role,
  String? title,
  String? banner,
  String? profilePicture,
  String? usernames,
  DocumentReference? bagRef,
  DocumentReference? orderRef,
  String? userType,
  String? shortDescription,
  DocumentReference? analyticsRef,
  bool? userVerified,
  bool? userVerifiedPending,
  bool? userIsAdmin,
  double? avgURating,
  int? ratings,
  DocumentReference? userCredits,
  String? walletAddress,
  String? verificationIdUrl,
  String? verificationSelfieUrl,
  String? verificationLocationUrl,
  String? verificationMunicipality,
  DateTime? verificationSubmittedAt,
}) {
  final firestoreData = mapToFirestore(
    <String, dynamic>{
      'bio': bio,
      'email': email,
      'location': location,
      'name': name,
      'surname': surname,
      'username': username,
      'display_name': displayName,
      'photo_url': photoUrl,
      'uid': uid,
      'created_time': createdTime,
      'phone_number': phoneNumber,
      'last_active_time': lastActiveTime,
      'role': role,
      'title': title,
      'banner': banner,
      'profile_picture': profilePicture,
      'usernames': usernames,
      'bag_ref': bagRef,
      'order_ref': orderRef,
      'user_type': userType,
      'shortDescription': shortDescription,
      'analytics_ref': analyticsRef,
      'user_verified': userVerified,
      'user_verified_pending': userVerifiedPending,
      'user_is_admin': userIsAdmin,
      'avg_u_rating': avgURating,
      'ratings': ratings,
      'userCredits': userCredits,
      'wallet_address': walletAddress,
      'verification_id_url': verificationIdUrl,
      'verification_selfie_url': verificationSelfieUrl,
      'verification_location_url': verificationLocationUrl,
      'verification_municipality': verificationMunicipality,
      'verification_submitted_at': verificationSubmittedAt,
    }.withoutNulls,
  );

  return firestoreData;
}

class UsersRecordDocumentEquality implements Equality<UsersRecord> {
  const UsersRecordDocumentEquality();

  @override
  bool equals(UsersRecord? e1, UsersRecord? e2) {
    const listEquality = ListEquality();
    return e1?.bio == e2?.bio &&
        e1?.email == e2?.email &&
        e1?.location == e2?.location &&
        e1?.name == e2?.name &&
        e1?.surname == e2?.surname &&
        e1?.username == e2?.username &&
        e1?.displayName == e2?.displayName &&
        e1?.photoUrl == e2?.photoUrl &&
        e1?.uid == e2?.uid &&
        e1?.createdTime == e2?.createdTime &&
        e1?.phoneNumber == e2?.phoneNumber &&
        e1?.lastActiveTime == e2?.lastActiveTime &&
        e1?.role == e2?.role &&
        e1?.title == e2?.title &&
        e1?.banner == e2?.banner &&
        e1?.profilePicture == e2?.profilePicture &&
        e1?.usernames == e2?.usernames &&
        listEquality.equals(e1?.pinnedUsers, e2?.pinnedUsers) &&
        listEquality.equals(e1?.userPins, e2?.userPins) &&
        e1?.bagRef == e2?.bagRef &&
        listEquality.equals(e1?.userOccupations, e2?.userOccupations) &&
        listEquality.equals(e1?.userInterests, e2?.userInterests) &&
        listEquality.equals(e1?.pinnedObjects, e2?.pinnedObjects) &&
        listEquality.equals(e1?.userObjects, e2?.userObjects) &&
        listEquality.equals(e1?.walletMethodsUser, e2?.walletMethodsUser) &&
        listEquality.equals(e1?.userBagObjects, e2?.userBagObjects) &&
        e1?.orderRef == e2?.orderRef &&
        listEquality.equals(e1?.orderMethodsUser, e2?.orderMethodsUser) &&
        listEquality.equals(e1?.userPlaces, e2?.userPlaces) &&
        e1?.userType == e2?.userType &&
        listEquality.equals(e1?.userItems, e2?.userItems) &&
        listEquality.equals(e1?.userPosts, e2?.userPosts) &&
        listEquality.equals(e1?.userTransactions, e2?.userTransactions) &&
        e1?.shortDescription == e2?.shortDescription &&
        e1?.analyticsRef == e2?.analyticsRef &&
        e1?.userVerified == e2?.userVerified &&
        e1?.userVerifiedPending == e2?.userVerifiedPending &&
        e1?.userIsAdmin == e2?.userIsAdmin &&
        listEquality.equals(e1?.uRatings, e2?.uRatings) &&
        e1?.avgURating == e2?.avgURating &&
        e1?.ratings == e2?.ratings &&
        listEquality.equals(e1?.userRatings, e2?.userRatings) &&
        e1?.userCredits == e2?.userCredits &&
        e1?.walletAddress == e2?.walletAddress &&
        listEquality.equals(e1?.followingUsers, e2?.followingUsers) &&
        listEquality.equals(e1?.usersFollwoingMe, e2?.usersFollwoingMe) &&
        e1?.verificationIdUrl == e2?.verificationIdUrl &&
        e1?.verificationSelfieUrl == e2?.verificationSelfieUrl &&
        e1?.verificationLocationUrl == e2?.verificationLocationUrl &&
        e1?.verificationMunicipality == e2?.verificationMunicipality &&
        e1?.verificationSubmittedAt == e2?.verificationSubmittedAt;
  }

  @override
  int hash(UsersRecord? e) => const ListEquality().hash([
        e?.bio,
        e?.email,
        e?.location,
        e?.name,
        e?.surname,
        e?.username,
        e?.displayName,
        e?.photoUrl,
        e?.uid,
        e?.createdTime,
        e?.phoneNumber,
        e?.lastActiveTime,
        e?.role,
        e?.title,
        e?.banner,
        e?.profilePicture,
        e?.usernames,
        e?.pinnedUsers,
        e?.userPins,
        e?.bagRef,
        e?.userOccupations,
        e?.userInterests,
        e?.pinnedObjects,
        e?.userObjects,
        e?.walletMethodsUser,
        e?.userBagObjects,
        e?.orderRef,
        e?.orderMethodsUser,
        e?.userPlaces,
        e?.userType,
        e?.userItems,
        e?.userPosts,
        e?.userTransactions,
        e?.shortDescription,
        e?.analyticsRef,
        e?.userVerified,
        e?.userVerifiedPending,
        e?.userIsAdmin,
        e?.uRatings,
        e?.avgURating,
        e?.ratings,
        e?.userRatings,
        e?.userCredits,
        e?.walletAddress,
        e?.followingUsers,
        e?.usersFollwoingMe,
        e?.verificationIdUrl,
        e?.verificationSelfieUrl,
        e?.verificationLocationUrl,
        e?.verificationMunicipality,
        e?.verificationSubmittedAt,
      ]);

  @override
  bool isValidKey(Object? o) => o is UsersRecord;
}
