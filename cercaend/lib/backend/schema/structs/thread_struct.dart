// ignore_for_file: unnecessary_getters_setters

import '/backend/supabase/supabase_shim.dart';

import '/backend/schema/util/firestore_util.dart';
import '/backend/schema/util/schema_util.dart';

import 'index.dart';
import '/flutter_flow/flutter_flow_util.dart';

class ThreadStruct extends FFFirebaseStruct {
  ThreadStruct({
    DocumentReference? userRef,
    DocumentReference? subissionRef,
    List<String>? objectThread,
    String? thread1,
    String? thread2,
    String? thread3,
    String? thread4,
    String? thread5,
    FirestoreUtilData firestoreUtilData = const FirestoreUtilData(),
  })  : _userRef = userRef,
        _subissionRef = subissionRef,
        _objectThread = objectThread,
        _thread1 = thread1,
        _thread2 = thread2,
        _thread3 = thread3,
        _thread4 = thread4,
        _thread5 = thread5,
        super(firestoreUtilData);

  // "user_ref" field.
  DocumentReference? _userRef;
  DocumentReference? get userRef => _userRef;
  set userRef(DocumentReference? val) => _userRef = val;

  bool hasUserRef() => _userRef != null;

  // "subission_ref" field.
  DocumentReference? _subissionRef;
  DocumentReference? get subissionRef => _subissionRef;
  set subissionRef(DocumentReference? val) => _subissionRef = val;

  bool hasSubissionRef() => _subissionRef != null;

  // "object_thread" field.
  List<String>? _objectThread;
  List<String> get objectThread => _objectThread ?? const [];
  set objectThread(List<String>? val) => _objectThread = val;

  void updateObjectThread(Function(List<String>) updateFn) {
    updateFn(_objectThread ??= []);
  }

  bool hasObjectThread() => _objectThread != null;

  // "thread_1" field.
  String? _thread1;
  String get thread1 => _thread1 ?? '';
  set thread1(String? val) => _thread1 = val;

  bool hasThread1() => _thread1 != null;

  // "thread_2" field.
  String? _thread2;
  String get thread2 => _thread2 ?? '';
  set thread2(String? val) => _thread2 = val;

  bool hasThread2() => _thread2 != null;

  // "thread_3" field.
  String? _thread3;
  String get thread3 => _thread3 ?? '';
  set thread3(String? val) => _thread3 = val;

  bool hasThread3() => _thread3 != null;

  // "thread_4" field.
  String? _thread4;
  String get thread4 => _thread4 ?? '';
  set thread4(String? val) => _thread4 = val;

  bool hasThread4() => _thread4 != null;

  // "thread_5" field.
  String? _thread5;
  String get thread5 => _thread5 ?? '';
  set thread5(String? val) => _thread5 = val;

  bool hasThread5() => _thread5 != null;

  static ThreadStruct fromMap(Map<String, dynamic> data) => ThreadStruct(
        userRef: safeDocRef(data['user_ref']),
        subissionRef: safeDocRef(data['subission_ref']),
        objectThread: getDataList(data['object_thread']),
        thread1: data['thread_1'] as String?,
        thread2: data['thread_2'] as String?,
        thread3: data['thread_3'] as String?,
        thread4: data['thread_4'] as String?,
        thread5: data['thread_5'] as String?,
      );

  static ThreadStruct? maybeFromMap(dynamic data) =>
      data is Map ? ThreadStruct.fromMap(data.cast<String, dynamic>()) : null;

  Map<String, dynamic> toMap() => {
        'user_ref': _userRef,
        'subission_ref': _subissionRef,
        'object_thread': _objectThread,
        'thread_1': _thread1,
        'thread_2': _thread2,
        'thread_3': _thread3,
        'thread_4': _thread4,
        'thread_5': _thread5,
      }.withoutNulls;

  @override
  Map<String, dynamic> toSerializableMap() => {
        'user_ref': serializeParam(
          _userRef,
          ParamType.DocumentReference,
        ),
        'subission_ref': serializeParam(
          _subissionRef,
          ParamType.DocumentReference,
        ),
        'object_thread': serializeParam(
          _objectThread,
          ParamType.String,
          isList: true,
        ),
        'thread_1': serializeParam(
          _thread1,
          ParamType.String,
        ),
        'thread_2': serializeParam(
          _thread2,
          ParamType.String,
        ),
        'thread_3': serializeParam(
          _thread3,
          ParamType.String,
        ),
        'thread_4': serializeParam(
          _thread4,
          ParamType.String,
        ),
        'thread_5': serializeParam(
          _thread5,
          ParamType.String,
        ),
      }.withoutNulls;

  static ThreadStruct fromSerializableMap(Map<String, dynamic> data) =>
      ThreadStruct(
        userRef: deserializeParam(
          data['user_ref'],
          ParamType.DocumentReference,
          false,
          collectionNamePath: ['users'],
        ),
        subissionRef: deserializeParam(
          data['subission_ref'],
          ParamType.DocumentReference,
          false,
          collectionNamePath: ['submission'],
        ),
        objectThread: deserializeParam<String>(
          data['object_thread'],
          ParamType.String,
          true,
        ),
        thread1: deserializeParam(
          data['thread_1'],
          ParamType.String,
          false,
        ),
        thread2: deserializeParam(
          data['thread_2'],
          ParamType.String,
          false,
        ),
        thread3: deserializeParam(
          data['thread_3'],
          ParamType.String,
          false,
        ),
        thread4: deserializeParam(
          data['thread_4'],
          ParamType.String,
          false,
        ),
        thread5: deserializeParam(
          data['thread_5'],
          ParamType.String,
          false,
        ),
      );

  @override
  String toString() => 'ThreadStruct(${toMap()})';

  @override
  bool operator ==(Object other) {
    const listEquality = ListEquality();
    return other is ThreadStruct &&
        userRef == other.userRef &&
        subissionRef == other.subissionRef &&
        listEquality.equals(objectThread, other.objectThread) &&
        thread1 == other.thread1 &&
        thread2 == other.thread2 &&
        thread3 == other.thread3 &&
        thread4 == other.thread4 &&
        thread5 == other.thread5;
  }

  @override
  int get hashCode => const ListEquality().hash([
        userRef,
        subissionRef,
        objectThread,
        thread1,
        thread2,
        thread3,
        thread4,
        thread5
      ]);
}

ThreadStruct createThreadStruct({
  DocumentReference? userRef,
  DocumentReference? subissionRef,
  String? thread1,
  String? thread2,
  String? thread3,
  String? thread4,
  String? thread5,
  Map<String, dynamic> fieldValues = const {},
  bool clearUnsetFields = true,
  bool create = false,
  bool delete = false,
}) =>
    ThreadStruct(
      userRef: userRef,
      subissionRef: subissionRef,
      thread1: thread1,
      thread2: thread2,
      thread3: thread3,
      thread4: thread4,
      thread5: thread5,
      firestoreUtilData: FirestoreUtilData(
        clearUnsetFields: clearUnsetFields,
        create: create,
        delete: delete,
        fieldValues: fieldValues,
      ),
    );

ThreadStruct? updateThreadStruct(
  ThreadStruct? thread, {
  bool clearUnsetFields = true,
  bool create = false,
}) =>
    thread
      ?..firestoreUtilData = FirestoreUtilData(
        clearUnsetFields: clearUnsetFields,
        create: create,
      );

void addThreadStructData(
  Map<String, dynamic> firestoreData,
  ThreadStruct? thread,
  String fieldName, [
  bool forFieldValue = false,
]) {
  firestoreData.remove(fieldName);
  if (thread == null) {
    return;
  }
  if (thread.firestoreUtilData.delete) {
    firestoreData[fieldName] = FieldValue.delete();
    return;
  }
  final clearFields =
      !forFieldValue && thread.firestoreUtilData.clearUnsetFields;
  if (clearFields) {
    firestoreData[fieldName] = <String, dynamic>{};
  }
  final threadData = getThreadFirestoreData(thread, forFieldValue);
  final nestedData = threadData.map((k, v) => MapEntry('$fieldName.$k', v));

  final mergeFields = thread.firestoreUtilData.create || clearFields;
  firestoreData
      .addAll(mergeFields ? mergeNestedFields(nestedData) : nestedData);
}

Map<String, dynamic> getThreadFirestoreData(
  ThreadStruct? thread, [
  bool forFieldValue = false,
]) {
  if (thread == null) {
    return {};
  }
  final firestoreData = mapToFirestore(thread.toMap());

  // Add any Firestore field values
  thread.firestoreUtilData.fieldValues.forEach((k, v) => firestoreData[k] = v);

  return forFieldValue ? mergeNestedFields(firestoreData) : firestoreData;
}

List<Map<String, dynamic>> getThreadListFirestoreData(
  List<ThreadStruct>? threads,
) =>
    threads?.map((e) => getThreadFirestoreData(e, true)).toList() ?? [];
