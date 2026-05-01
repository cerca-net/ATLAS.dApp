// ignore_for_file: unnecessary_getters_setters

import '/backend/supabase/supabase_shim.dart';

import '/backend/schema/util/firestore_util.dart';
import '/backend/schema/util/schema_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class ObjectStruct extends FFFirebaseStruct {
  ObjectStruct({
    List<String>? wordThread,
    FirestoreUtilData firestoreUtilData = const FirestoreUtilData(),
  })  : _wordThread = wordThread,
        super(firestoreUtilData);

  // "word_thread" field.
  List<String>? _wordThread;
  List<String> get wordThread => _wordThread ?? const [];
  set wordThread(List<String>? val) => _wordThread = val;

  void updateWordThread(Function(List<String>) updateFn) {
    updateFn(_wordThread ??= []);
  }

  bool hasWordThread() => _wordThread != null;

  static ObjectStruct fromMap(Map<String, dynamic> data) => ObjectStruct(
        wordThread: getDataList(data['word_thread']),
      );

  static ObjectStruct? maybeFromMap(dynamic data) =>
      data is Map ? ObjectStruct.fromMap(data.cast<String, dynamic>()) : null;

  Map<String, dynamic> toMap() => {
        'word_thread': _wordThread,
      }.withoutNulls;

  @override
  Map<String, dynamic> toSerializableMap() => {
        'word_thread': serializeParam(
          _wordThread,
          ParamType.String,
          isList: true,
        ),
      }.withoutNulls;

  static ObjectStruct fromSerializableMap(Map<String, dynamic> data) =>
      ObjectStruct(
        wordThread: deserializeParam<String>(
          data['word_thread'],
          ParamType.String,
          true,
        ),
      );

  @override
  String toString() => 'ObjectStruct(${toMap()})';

  @override
  bool operator ==(Object other) {
    const listEquality = ListEquality();
    return other is ObjectStruct &&
        listEquality.equals(wordThread, other.wordThread);
  }

  @override
  int get hashCode => const ListEquality().hash([wordThread]);
}

ObjectStruct createObjectStruct({
  Map<String, dynamic> fieldValues = const {},
  bool clearUnsetFields = true,
  bool create = false,
  bool delete = false,
}) =>
    ObjectStruct(
      firestoreUtilData: FirestoreUtilData(
        clearUnsetFields: clearUnsetFields,
        create: create,
        delete: delete,
        fieldValues: fieldValues,
      ),
    );

ObjectStruct? updateObjectStruct(
  ObjectStruct? object, {
  bool clearUnsetFields = true,
  bool create = false,
}) =>
    object
      ?..firestoreUtilData = FirestoreUtilData(
        clearUnsetFields: clearUnsetFields,
        create: create,
      );

void addObjectStructData(
  Map<String, dynamic> firestoreData,
  ObjectStruct? object,
  String fieldName, [
  bool forFieldValue = false,
]) {
  firestoreData.remove(fieldName);
  if (object == null) {
    return;
  }
  if (object.firestoreUtilData.delete) {
    firestoreData[fieldName] = FieldValue.delete();
    return;
  }
  final clearFields =
      !forFieldValue && object.firestoreUtilData.clearUnsetFields;
  if (clearFields) {
    firestoreData[fieldName] = <String, dynamic>{};
  }
  final objectData = getObjectFirestoreData(object, forFieldValue);
  final nestedData = objectData.map((k, v) => MapEntry('$fieldName.$k', v));

  final mergeFields = object.firestoreUtilData.create || clearFields;
  firestoreData
      .addAll(mergeFields ? mergeNestedFields(nestedData) : nestedData);
}

Map<String, dynamic> getObjectFirestoreData(
  ObjectStruct? object, [
  bool forFieldValue = false,
]) {
  if (object == null) {
    return {};
  }
  final firestoreData = mapToFirestore(object.toMap());

  // Add any Firestore field values
  object.firestoreUtilData.fieldValues.forEach((k, v) => firestoreData[k] = v);

  return forFieldValue ? mergeNestedFields(firestoreData) : firestoreData;
}

List<Map<String, dynamic>> getObjectListFirestoreData(
  List<ObjectStruct>? objects,
) =>
    objects?.map((e) => getObjectFirestoreData(e, true)).toList() ?? [];
