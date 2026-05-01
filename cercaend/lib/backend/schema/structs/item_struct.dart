// ignore_for_file: unnecessary_getters_setters

import '/backend/supabase/supabase_shim.dart';

import '/backend/schema/util/firestore_util.dart';
import '/backend/schema/util/schema_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class ItemStruct extends FFFirebaseStruct {
  ItemStruct({
    String? name,
    String? description,
    double? refValue,
    List<DocumentReference>? itemImages,
    FirestoreUtilData firestoreUtilData = const FirestoreUtilData(),
  })  : _name = name,
        _description = description,
        _refValue = refValue,
        _itemImages = itemImages,
        super(firestoreUtilData);

  // "name" field.
  String? _name;
  String get name => _name ?? '';
  set name(String? val) => _name = val;

  bool hasName() => _name != null;

  // "description" field.
  String? _description;
  String get description => _description ?? '';
  set description(String? val) => _description = val;

  bool hasDescription() => _description != null;

  // "ref_value" field.
  double? _refValue;
  double get refValue => _refValue ?? 0.0;
  set refValue(double? val) => _refValue = val;

  void incrementRefValue(double amount) => refValue = refValue + amount;

  bool hasRefValue() => _refValue != null;

  // "item_images" field.
  List<DocumentReference>? _itemImages;
  List<DocumentReference> get itemImages => _itemImages ?? const [];
  set itemImages(List<DocumentReference>? val) => _itemImages = val;

  void updateItemImages(Function(List<DocumentReference>) updateFn) {
    updateFn(_itemImages ??= []);
  }

  bool hasItemImages() => _itemImages != null;

  static ItemStruct fromMap(Map<String, dynamic> data) => ItemStruct(
        name: data['name'] as String?,
        description: data['description'] as String?,
        refValue: castToType<double>(data['ref_value']),
        itemImages: getDataList(data['item_images']),
      );

  static ItemStruct? maybeFromMap(dynamic data) =>
      data is Map ? ItemStruct.fromMap(data.cast<String, dynamic>()) : null;

  Map<String, dynamic> toMap() => {
        'name': _name,
        'description': _description,
        'ref_value': _refValue,
        'item_images': _itemImages,
      }.withoutNulls;

  @override
  Map<String, dynamic> toSerializableMap() => {
        'name': serializeParam(
          _name,
          ParamType.String,
        ),
        'description': serializeParam(
          _description,
          ParamType.String,
        ),
        'ref_value': serializeParam(
          _refValue,
          ParamType.double,
        ),
        'item_images': serializeParam(
          _itemImages,
          ParamType.DocumentReference,
          isList: true,
        ),
      }.withoutNulls;

  static ItemStruct fromSerializableMap(Map<String, dynamic> data) =>
      ItemStruct(
        name: deserializeParam(
          data['name'],
          ParamType.String,
          false,
        ),
        description: deserializeParam(
          data['description'],
          ParamType.String,
          false,
        ),
        refValue: deserializeParam(
          data['ref_value'],
          ParamType.double,
          false,
        ),
        itemImages: deserializeParam<DocumentReference>(
          data['item_images'],
          ParamType.DocumentReference,
          true,
          collectionNamePath: ['submission'],
        ),
      );

  @override
  String toString() => 'ItemStruct(${toMap()})';

  @override
  bool operator ==(Object other) {
    const listEquality = ListEquality();
    return other is ItemStruct &&
        name == other.name &&
        description == other.description &&
        refValue == other.refValue &&
        listEquality.equals(itemImages, other.itemImages);
  }

  @override
  int get hashCode =>
      const ListEquality().hash([name, description, refValue, itemImages]);
}

ItemStruct createItemStruct({
  String? name,
  String? description,
  double? refValue,
  Map<String, dynamic> fieldValues = const {},
  bool clearUnsetFields = true,
  bool create = false,
  bool delete = false,
}) =>
    ItemStruct(
      name: name,
      description: description,
      refValue: refValue,
      firestoreUtilData: FirestoreUtilData(
        clearUnsetFields: clearUnsetFields,
        create: create,
        delete: delete,
        fieldValues: fieldValues,
      ),
    );

ItemStruct? updateItemStruct(
  ItemStruct? item, {
  bool clearUnsetFields = true,
  bool create = false,
}) =>
    item
      ?..firestoreUtilData = FirestoreUtilData(
        clearUnsetFields: clearUnsetFields,
        create: create,
      );

void addItemStructData(
  Map<String, dynamic> firestoreData,
  ItemStruct? item,
  String fieldName, [
  bool forFieldValue = false,
]) {
  firestoreData.remove(fieldName);
  if (item == null) {
    return;
  }
  if (item.firestoreUtilData.delete) {
    firestoreData[fieldName] = FieldValue.delete();
    return;
  }
  final clearFields = !forFieldValue && item.firestoreUtilData.clearUnsetFields;
  if (clearFields) {
    firestoreData[fieldName] = <String, dynamic>{};
  }
  final itemData = getItemFirestoreData(item, forFieldValue);
  final nestedData = itemData.map((k, v) => MapEntry('$fieldName.$k', v));

  final mergeFields = item.firestoreUtilData.create || clearFields;
  firestoreData
      .addAll(mergeFields ? mergeNestedFields(nestedData) : nestedData);
}

Map<String, dynamic> getItemFirestoreData(
  ItemStruct? item, [
  bool forFieldValue = false,
]) {
  if (item == null) {
    return {};
  }
  final firestoreData = mapToFirestore(item.toMap());

  // Add any Firestore field values
  item.firestoreUtilData.fieldValues.forEach((k, v) => firestoreData[k] = v);

  return forFieldValue ? mergeNestedFields(firestoreData) : firestoreData;
}

List<Map<String, dynamic>> getItemListFirestoreData(
  List<ItemStruct>? items,
) =>
    items?.map((e) => getItemFirestoreData(e, true)).toList() ?? [];
