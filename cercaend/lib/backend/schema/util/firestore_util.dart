import '/backend/supabase/supabase_shim.dart';

import '/backend/schema/util/schema_util.dart';
import '/flutter_flow/flutter_flow_util.dart';

typedef RecordBuilder<T> = T Function(DocumentSnapshot snapshot);

abstract class FirestoreRecord {
  FirestoreRecord(this.reference, this.snapshotData);
  Map<String, dynamic> snapshotData;
  DocumentReference reference;
}

abstract class FFFirebaseStruct extends BaseStruct {
  FFFirebaseStruct(this.firestoreUtilData);

  /// Utility class for Firestore updates
  FirestoreUtilData firestoreUtilData = const FirestoreUtilData();
}

class FirestoreUtilData {
  const FirestoreUtilData({
    this.fieldValues = const {},
    this.clearUnsetFields = true,
    this.create = false,
    this.delete = false,
  });
  final Map<String, dynamic> fieldValues;
  final bool clearUnsetFields;
  final bool create;
  final bool delete;
  static String get name => 'firestoreUtilData';
}

Map<String, dynamic> mapFromFirestore(Map<String, dynamic> data) =>
    mergeNestedFields(data)
        .where((k, _) => k != FirestoreUtilData.name)
        .map((key, value) {
      // Handle Timestamp (Firebase)
      if (value is Timestamp) {
        value = value.toDate();
      }
      // Handle boolean strings from Supabase text columns
      if (value is String) {
        final lower = value.toLowerCase();
        if (lower == 'true') {
          value = true;
        } else if (lower == 'false') {
          value = false;
        }
      }
      // Handle stringified DateTime (Supabase) — only for date-like patterns
      // Must look like a date: starts with 4 digits, contains dashes
      if (value is String && value.length >= 10 && RegExp(r'^\d{4}-\d{2}').hasMatch(value)) {
        final parsed = DateTime.tryParse(value);
        if (parsed != null) {
          value = parsed;
        }
      }
      // Handle list of Timestamp
      if (value is Iterable && value.isNotEmpty && value.first is Timestamp) {
        value = value.map((v) => (v as Timestamp).toDate()).toList();
      }
      // Handle GeoPoint
      if (value is GeoPoint) {
        value = value.toLatLng();
      }
      // Handle list of GeoPoint
      if (value is Iterable && value.isNotEmpty && value.first is GeoPoint) {
        value = value.map((v) => (v as GeoPoint).toLatLng()).toList();
      }
      // Handle nested data.
      if (value is Map) {
        value = mapFromFirestore(value as Map<String, dynamic>);
      }
      // Handle list of nested data.
      if (value is Iterable && value.isNotEmpty && value.first is Map) {
        value = value
            .map((v) => mapFromFirestore(v as Map<String, dynamic>))
            .toList();
      }
      return MapEntry(key, value);
    });

Map<String, dynamic> mapToFirestore(Map<String, dynamic> data) =>
    data.where((k, v) => k != FirestoreUtilData.name).map((key, value) {
      // Handle DocumentReference → serialize to ID string for Supabase
      if (value is DocumentReference) {
        value = value.id;
      }
      // Handle list of DocumentReference
      if (value is Iterable && value.isNotEmpty && value.first is DocumentReference) {
        value = value.map((v) => (v as DocumentReference).id).toList();
      }
      // Handle DateTime → ISO string for Supabase text columns
      if (value is DateTime) {
        value = value.toIso8601String();
      }
      // Handle Timestamp
      if (value is Timestamp) {
        value = value.toDate().toIso8601String();
      }
      // Handle GeoPoint
      if (value is LatLng) {
        value = value.toGeoPoint();
      }
      // Handle list of GeoPoint
      if (value is Iterable && value.isNotEmpty && value.first is LatLng) {
        value = value.map((v) => (v as LatLng).toGeoPoint()).toList();
      }
      // Handle Color
      if (value is Color) {
        value = value.toCssString();
      }
      // Handle list of Color
      if (value is Iterable && value.isNotEmpty && value.first is Color) {
        value = value.map((v) => (v as Color).toCssString()).toList();
      } // Handle Enums.
      if (value is Enum) {
        value = value.serialize();
      }
      // Handle list of Enums.
      if (value is Iterable && value.isNotEmpty && value.first is Enum) {
        value = value.map((v) => (v as Enum).serialize()).toList();
      }
      // Handle nested data.
      if (value is Map) {
        value = mapToFirestore(value as Map<String, dynamic>);
      }
      // Handle list of nested data.
      if (value is Iterable && value.isNotEmpty && value.first is Map) {
        value = value
            .map((v) => mapToFirestore(v as Map<String, dynamic>))
            .toList();
      }
      return MapEntry(key, value);
    });

List<GeoPoint>? convertToGeoPointList(List<LatLng>? list) =>
    list?.map((e) => e.toGeoPoint()).toList();

extension GeoPointExtension on LatLng {
  GeoPoint toGeoPoint() => GeoPoint(latitude, longitude);
}

extension LatLngExtension on GeoPoint {
  LatLng toLatLng() => LatLng(latitude, longitude);
}

/// Safely converts a value to DocumentReference, handling String paths from Supabase
DocumentReference? safeDocRef(dynamic value) {
  if (value == null) return null;
  if (value is DocumentReference) return value;
  if (value is String && value.isNotEmpty) {
    return FirebaseFirestore.instance.doc(value.contains('/') ? value : 'unknown/$value');
  }
  return null;
}

DocumentReference toRef(String ref) => FirebaseFirestore.instance.doc(ref);

T? safeGet<T>(T Function() func, [Function(dynamic)? reportError]) {
  try {
    return func();
  } catch (e) {
    reportError?.call(e);
  }
  return null;
}

Map<String, dynamic> mergeNestedFields(Map<String, dynamic> data) {
  final nestedData = data.where((k, _) => k.contains('.'));
  final fieldNames = nestedData.keys.map((k) => k.split('.').first).toSet();
  // Remove nested values (e.g. 'foo.bar') and merge them into a map.
  data.removeWhere((k, _) => k.contains('.'));
  for (var name in fieldNames) {
    final mergedValues = mergeNestedFields(
      nestedData
          .where((k, _) => k.split('.').first == name)
          .map((k, v) => MapEntry(k.split('.').skip(1).join('.'), v)),
    );
    final existingValue = data[name];
    data[name] = {
      if (existingValue != null && existingValue is Map)
        ...existingValue as Map<String, dynamic>,
      ...mergedValues,
    };
  }
  // Merge any nested maps inside any of the fields as well.
  data.where((_, v) => v is Map).forEach((k, v) {
    data[k] = mergeNestedFields(v as Map<String, dynamic>);
  });

  return data;
}

extension _WhereMapExtension<K, V> on Map<K, V> {
  Map<K, V> where(bool Function(K, V) test) =>
      Map.fromEntries(entries.where((e) => test(e.key, e.value)));
}
