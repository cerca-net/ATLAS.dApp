// ignore_for_file: unnecessary_getters_setters

import '/backend/supabase/supabase_shim.dart';

import '/backend/schema/util/firestore_util.dart';

import '/flutter_flow/flutter_flow_util.dart';

class BagItemStruct extends FFFirebaseStruct {
  BagItemStruct({
    DocumentReference? itemRef,
    String? name,
    String? image,
    double? refValue,
    FirestoreUtilData firestoreUtilData = const FirestoreUtilData(),
  })  : _itemRef = itemRef,
        _name = name,
        _image = image,
        _refValue = refValue,
        super(firestoreUtilData);

  // "item_ref" field.
  DocumentReference? _itemRef;
  DocumentReference? get itemRef => _itemRef;
  set itemRef(DocumentReference? val) => _itemRef = val;

  bool hasItemRef() => _itemRef != null;

  // "name" field.
  String? _name;
  String get name => _name ?? '';
  set name(String? val) => _name = val;

  bool hasName() => _name != null;

  // "image" field.
  String? _image;
  String get image => _image ?? '';
  set image(String? val) => _image = val;

  bool hasImage() => _image != null;

  // "ref_value" field.
  double? _refValue;
  double get refValue => _refValue ?? 0.0;
  set refValue(double? val) => _refValue = val;

  void incrementRefValue(double amount) => refValue = refValue + amount;

  bool hasRefValue() => _refValue != null;

  static BagItemStruct fromMap(Map<String, dynamic> data) => BagItemStruct(
        itemRef: safeDocRef(data['item_ref']),
        name: data['name'] as String?,
        image: data['image'] as String?,
        refValue: castToType<double>(data['ref_value']),
      );

  static BagItemStruct? maybeFromMap(dynamic data) =>
      data is Map ? BagItemStruct.fromMap(data.cast<String, dynamic>()) : null;

  Map<String, dynamic> toMap() => {
        'item_ref': _itemRef,
        'name': _name,
        'image': _image,
        'ref_value': _refValue,
      }.withoutNulls;

  @override
  Map<String, dynamic> toSerializableMap() => {
        'item_ref': serializeParam(
          _itemRef,
          ParamType.DocumentReference,
        ),
        'name': serializeParam(
          _name,
          ParamType.String,
        ),
        'image': serializeParam(
          _image,
          ParamType.String,
        ),
        'ref_value': serializeParam(
          _refValue,
          ParamType.double,
        ),
      }.withoutNulls;

  static BagItemStruct fromSerializableMap(Map<String, dynamic> data) =>
      BagItemStruct(
        itemRef: deserializeParam(
          data['item_ref'],
          ParamType.DocumentReference,
          false,
          collectionNamePath: ['submission'],
        ),
        name: deserializeParam(
          data['name'],
          ParamType.String,
          false,
        ),
        image: deserializeParam(
          data['image'],
          ParamType.String,
          false,
        ),
        refValue: deserializeParam(
          data['ref_value'],
          ParamType.double,
          false,
        ),
      );

  @override
  String toString() => 'BagItemStruct(${toMap()})';

  @override
  bool operator ==(Object other) {
    return other is BagItemStruct &&
        itemRef == other.itemRef &&
        name == other.name &&
        image == other.image &&
        refValue == other.refValue;
  }

  @override
  int get hashCode =>
      const ListEquality().hash([itemRef, name, image, refValue]);
}

BagItemStruct createBagItemStruct({
  DocumentReference? itemRef,
  String? name,
  String? image,
  double? refValue,
  Map<String, dynamic> fieldValues = const {},
  bool clearUnsetFields = true,
  bool create = false,
  bool delete = false,
}) =>
    BagItemStruct(
      itemRef: itemRef,
      name: name,
      image: image,
      refValue: refValue,
      firestoreUtilData: FirestoreUtilData(
        clearUnsetFields: clearUnsetFields,
        create: create,
        delete: delete,
        fieldValues: fieldValues,
      ),
    );

BagItemStruct? updateBagItemStruct(
  BagItemStruct? bagItem, {
  bool clearUnsetFields = true,
  bool create = false,
}) =>
    bagItem
      ?..firestoreUtilData = FirestoreUtilData(
        clearUnsetFields: clearUnsetFields,
        create: create,
      );

void addBagItemStructData(
  Map<String, dynamic> firestoreData,
  BagItemStruct? bagItem,
  String fieldName, [
  bool forFieldValue = false,
]) {
  firestoreData.remove(fieldName);
  if (bagItem == null) {
    return;
  }
  if (bagItem.firestoreUtilData.delete) {
    firestoreData[fieldName] = FieldValue.delete();
    return;
  }
  final clearFields =
      !forFieldValue && bagItem.firestoreUtilData.clearUnsetFields;
  if (clearFields) {
    firestoreData[fieldName] = <String, dynamic>{};
  }
  final bagItemData = getBagItemFirestoreData(bagItem, forFieldValue);
  final nestedData = bagItemData.map((k, v) => MapEntry('$fieldName.$k', v));

  final mergeFields = bagItem.firestoreUtilData.create || clearFields;
  firestoreData
      .addAll(mergeFields ? mergeNestedFields(nestedData) : nestedData);
}

Map<String, dynamic> getBagItemFirestoreData(
  BagItemStruct? bagItem, [
  bool forFieldValue = false,
]) {
  if (bagItem == null) {
    return {};
  }
  final firestoreData = mapToFirestore(bagItem.toMap());

  // Add any Firestore field values
  bagItem.firestoreUtilData.fieldValues.forEach((k, v) => firestoreData[k] = v);

  return forFieldValue ? mergeNestedFields(firestoreData) : firestoreData;
}

List<Map<String, dynamic>> getBagItemListFirestoreData(
  List<BagItemStruct>? bagItems,
) =>
    bagItems?.map((e) => getBagItemFirestoreData(e, true)).toList() ?? [];
