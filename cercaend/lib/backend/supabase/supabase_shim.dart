import 'dart:async';
import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:supabase_flutter/supabase_flutter.dart' as supabase;
import 'supabase.dart';

// Shim to replace Firebase Firestore types, bridging them to Supabase.
// This allows UI components to compile while routing data to Supabase.

// ─────────────────────────── Internal Filter Types ───────────────────────────

enum _WhereOp {
  isEqualTo,
  isLessThan,
  isLessThanOrEqualTo,
  isGreaterThan,
  isGreaterThanOrEqualTo,
  arrayContains,
  arrayContainsAny,
  whereIn,
  whereNotIn,
  isNull,
}

class _WhereClause {
  final String field;
  final _WhereOp op;
  final dynamic value;

  const _WhereClause(this.field, this.op, this.value);
}

class _OrderClause {
  final String field;
  final bool descending;

  const _OrderClause(this.field, this.descending);
}

// ─────────────────── Value Serialization Helpers ──────────────────────────

/// Converts a filter value to something Supabase/PostgREST understands.
/// DocumentReferences are serialized to their `.id` string.
dynamic _serializeValue(dynamic value) {
  if (value is DocumentReference) return value.id;
  if (value is Timestamp) return value.toDate().toIso8601String();
  if (value is DateTime) return value.toIso8601String();
  if (value is List) return value.map(_serializeValue).toList();
  return value;
}

/// Serializes an entire data map for Supabase writes (update/set/insert).
/// Handles DocumentReference, DateTime, Timestamp, FieldValue sentinels, and nested maps.
Map<String, dynamic> _serializeDataMap(Map<String, dynamic> data) {
  final result = <String, dynamic>{};
  for (final entry in data.entries) {
    final key = entry.key;
    var value = entry.value;

    // Skip FieldValue sentinels — handled separately in update()
    if (value is _FieldValueSentinel) continue;
    // Remove null-delete sentinels
    if (value is _FieldValueDelete) continue;

    // Serialize known types
    if (value is DocumentReference) {
      value = value.id;
    } else if (value is Timestamp) {
      value = value.toDate().toIso8601String();
    } else if (value is DateTime) {
      value = value.toIso8601String();
    } else if (value is Map<String, dynamic>) {
      value = _serializeDataMap(value);
    } else if (value is List) {
      value = value.map((e) {
        if (e is DocumentReference) return e.id;
        if (e is Timestamp) return e.toDate().toIso8601String();
        if (e is DateTime) return e.toIso8601String();
        if (e is Map<String, dynamic>) return _serializeDataMap(e);
        return e;
      }).toList();
    }

    result[key] = value;
  }
  return result;
}

/// Extracts the comparable value from a row field for in-memory operations.
dynamic _extractField(Map<String, dynamic> row, String field) {
  return row[field];
}

/// Compares two dynamic values for ordering.
int _compareValues(dynamic a, dynamic b) {
  if (a == null && b == null) return 0;
  if (a == null) return -1;
  if (b == null) return 1;
  if (a is num && b is num) return a.compareTo(b);
  return a.toString().compareTo(b.toString());
}

/// Safely converts a possibly-null/string value to a List for FieldValue array ops.
List<dynamic> _parseList(dynamic value) {
  if (value == null) return [];
  if (value is List) return List<dynamic>.from(value);
  // Supabase text columns may store arrays as JSON strings
  if (value is String) {
    final trimmed = value.trim();
    if (trimmed.startsWith('[')) {
      try {
        final decoded = json.decode(trimmed);
        if (decoded is List) return List<dynamic>.from(decoded);
      } catch (_) {}
    }
    // Comma-separated fallback
    if (trimmed.isNotEmpty) return trimmed.split(',').map((e) => e.trim()).toList();
  }
  return [];
}

/// Tests whether a single row passes an in-memory where clause.
bool _matchesClause(Map<String, dynamic> row, _WhereClause clause) {
  final rawField = _extractField(row, clause.field);
  final expected = _serializeValue(clause.value);

  switch (clause.op) {
    case _WhereOp.isEqualTo:
      if (expected == null) return rawField == null;
      // Support comparing DocumentReference IDs stored as text
      return rawField?.toString() == expected.toString();

    case _WhereOp.isLessThan:
      if (rawField == null || expected == null) return false;
      return _compareValues(rawField, expected) < 0;

    case _WhereOp.isLessThanOrEqualTo:
      if (rawField == null || expected == null) return false;
      return _compareValues(rawField, expected) <= 0;

    case _WhereOp.isGreaterThan:
      if (rawField == null || expected == null) return false;
      return _compareValues(rawField, expected) > 0;

    case _WhereOp.isGreaterThanOrEqualTo:
      if (rawField == null || expected == null) return false;
      return _compareValues(rawField, expected) >= 0;

    case _WhereOp.arrayContains:
      if (rawField == null) return false;
      final list = _parseList(rawField);
      return list.any((e) => e.toString() == expected.toString());

    case _WhereOp.arrayContainsAny:
      if (rawField == null || expected == null) return false;
      final list = _parseList(rawField);
      final targets = (expected as List).map((e) => e.toString()).toSet();
      return list.any((e) => targets.contains(e.toString()));

    case _WhereOp.whereIn:
      if (expected == null) return true; // null whereIn means no filter
      final targets = (expected as List).map((e) => e.toString()).toSet();
      return targets.contains(rawField?.toString());

    case _WhereOp.whereNotIn:
      if (expected == null) return true;
      final excluded = (expected as List).map((e) => e.toString()).toSet();
      return !excluded.contains(rawField?.toString());

    case _WhereOp.isNull:
      final checkNull = expected as bool? ?? true;
      return checkNull ? rawField == null : rawField != null;
  }
}




// ─────────────────────────── DocumentSnapshot ─────────────────────────────

class DocumentSnapshot<T> {
  final String id;
  final T? _data;
  final DocumentReference<T> reference;

  DocumentSnapshot(this.id, this._data, this.reference);

  bool get exists => _data != null;
  T? data() => _data;
  
  dynamic get(String field) {
    if (_data == null) return null;
    if (_data is Map<String, dynamic>) {
      return (_data as Map<String, dynamic>)[field];
    }
    return null;
  }
}

// ─────────────────────────── DocumentReference ────────────────────────────

class DocumentReference<T> {
  final String path;
  final String id;

  DocumentReference(this.path, this.id);

  String get collectionName => path.split('/').first;

  CollectionReference<Map<String, dynamic>> collection(String collectionPath) =>
      CollectionReference<Map<String, dynamic>>('$path/$id/$collectionPath');

  CollectionReference<T> get parent => CollectionReference<T>(collectionName);

  Future<void> update(Map<String, dynamic> data) async {
    try {
      // Handle FieldValue sentinels that require read-then-write
      final hasFieldValueOps = data.values.any((v) => v is _FieldValueSentinel);
      Map<String, dynamic>? currentRow;

      if (hasFieldValueOps) {
        final rows = await SupaFlow.client.from(collectionName).select().eq('id', id);
        if (rows.isNotEmpty) currentRow = rows.first;
      }

      final serialized = _serializeDataMap(data);

      // Apply FieldValue operations
      for (final entry in data.entries) {
        if (entry.value is _FieldValueIncrement) {
          final inc = entry.value as _FieldValueIncrement;
          final current = currentRow?[entry.key];
          final currentNum = (current is num) ? current : num.tryParse(current?.toString() ?? '0') ?? 0;
          serialized[entry.key] = currentNum + inc.value;
        } else if (entry.value is _FieldValueArrayUnion) {
          final au = entry.value as _FieldValueArrayUnion;
          final current = currentRow?[entry.key];
          final currentList = _parseList(current);
          final newItems = au.elements.map(_serializeValue).toList();
          final merged = [...currentList];
          for (final item in newItems) {
            if (!merged.any((e) => e.toString() == item.toString())) {
              merged.add(item);
            }
          }
          serialized[entry.key] = merged;
        } else if (entry.value is _FieldValueArrayRemove) {
          final ar = entry.value as _FieldValueArrayRemove;
          final current = currentRow?[entry.key];
          final currentList = _parseList(current);
          final removeItems = ar.elements.map(_serializeValue).map((e) => e.toString()).toSet();
          serialized[entry.key] = currentList.where((e) => !removeItems.contains(e.toString())).toList();
        } else if (entry.value is _FieldValueServerTimestamp) {
          serialized[entry.key] = DateTime.now().toIso8601String();
        }
      }

      if (serialized.isNotEmpty) {
        await SupaFlow.client.from(collectionName).update(serialized).eq('id', id);
      }
    } catch (e) {
      debugPrint('DocumentReference.update() error on $collectionName/$id: $e');
    }
  }

  Future<void> set(Map<String, dynamic> data) async {
    try {
      final serialized = _serializeDataMap(data);
      if (!serialized.containsKey('id')) {
        serialized['id'] = id;
      }
      await SupaFlow.client.from(collectionName).upsert(serialized);
    } catch (e) {
      debugPrint('DocumentReference.set() error on $collectionName/$id: $e');
    }
  }

  Future<void> delete() async {
    await SupaFlow.client.from(collectionName).delete().eq('id', id);
  }

  Future<DocumentSnapshot<T>> get() async {
    try {
      final response = await SupaFlow.client.from(collectionName).select().eq('id', id).maybeSingle();
      if (response == null) {
        debugPrint('Supabase: No record found for $collectionName/$id');
        return DocumentSnapshot<T>(id, <String, dynamic>{'id': id} as T, this);
      }
      return DocumentSnapshot<T>(id, response as T, this);
    } catch (e) {
      debugPrint('Supabase Get Error ($collectionName/$id): $e');
      rethrow;
    }
  }

  Stream<DocumentSnapshot<T>> snapshots() {
    debugPrint('Supabase: Opening stream for $collectionName/$id');
    return SupaFlow.client
        .from(collectionName)
        .stream(primaryKey: ['id'])
        .eq('id', id)
        .map((rows) {
          if (rows.isEmpty) {
            debugPrint('Supabase: Stream empty for $collectionName/$id');
            return DocumentSnapshot<T>(id, <String, dynamic>{'id': id} as T, this);
          }
          debugPrint('Supabase: Stream update for $collectionName/$id');
          return DocumentSnapshot<T>(id, rows.first as T, this);
        });
  }
}

// ─────────────────────────── Helper Types ─────────────────────────────────

class GeoPoint {
  final double latitude;
  final double longitude;
  GeoPoint(this.latitude, this.longitude);
}

class Timestamp {
  final DateTime _date;
  Timestamp(this._date);
  DateTime toDate() => _date;
  
  static Timestamp now() => Timestamp(DateTime.now());
  static Timestamp fromDate(DateTime date) => Timestamp(date);
}

// ─────────────────────────── FirebaseFirestore ────────────────────────────

class FirebaseFirestore {
  static final instance = FirebaseFirestore._();
  FirebaseFirestore._();

  CollectionReference<Map<String, dynamic>> collection(String path) => CollectionReference<Map<String, dynamic>>(path);
  
  Query<Map<String, dynamic>> collectionGroup(String path) => Query<Map<String, dynamic>>(path);
  
  DocumentReference doc(String path) {
    var parts = path.split('/');
    var collectionLoc = parts.length > 1 ? parts.sublist(0, parts.length - 1).join('/') : 'unknown';
    var id = parts.isNotEmpty ? parts.last : '';
    return DocumentReference(collectionLoc, id);
  }
}

// ─────────────────────────── CollectionReference ──────────────────────────

class CollectionReference<T> extends Query<T> {
  CollectionReference(super.path);
  
  DocumentReference<T>? get parent {
    var parts = path.split('/');
    if (parts.length > 1) {
      return DocumentReference<T>(parts.sublist(0, parts.length - 2).join('/'), parts[parts.length - 2]);
    }
    return null;
  }

  DocumentReference<T> doc([String? id]) {
    return DocumentReference<T>(path, id ?? DateTime.now().millisecondsSinceEpoch.toString());
  }
  
  Future<DocumentReference<T>> add(Map<String, dynamic> data) async {
    final id = DateTime.now().millisecondsSinceEpoch.toString();
    try {
      final serialized = _serializeDataMap(data);
      serialized['id'] = id;
      await SupaFlow.client.from(path).insert(serialized);
    } catch (e) {
      debugPrint('CollectionReference.add() error on $path: $e');
    }
    return DocumentReference<T>(path, id);
  }
}

// ━━━━━━━━━━━━━━━━━━━━━ Query<T> — THE CORE FIX ━━━━━━━━━━━━━━━━━━━━━━━━

class Query<T> {
  final String path;
  final List<_WhereClause> _filters;
  final List<_OrderClause> _orders;
  final int? _limitCount;
  final DocumentSnapshot? _startAfterDoc;

  Query(this.path,
      {List<_WhereClause>? filters,
      List<_OrderClause>? orders,
      int? limitCount,
      DocumentSnapshot? startAfterDoc})
      : _filters = filters ?? const [],
        _orders = orders ?? const [],
        _limitCount = limitCount,
        _startAfterDoc = startAfterDoc;

  /// Clone with additional state.
  Query<T> _copyWith({
    List<_WhereClause>? filters,
    List<_OrderClause>? orders,
    int? limitCount,
    DocumentSnapshot? startAfterDoc,
  }) {
    return Query<T>(
      path,
      filters: filters ?? _filters,
      orders: orders ?? _orders,
      limitCount: limitCount ?? _limitCount,
      startAfterDoc: startAfterDoc ?? _startAfterDoc,
    );
  }

  // ─────────────── Filter builders (accumulate, return new Query) ────────

  Query<T> where(String field, {
    dynamic isEqualTo,
    dynamic isLessThan,
    dynamic isLessThanOrEqualTo,
    dynamic isGreaterThan,
    dynamic isGreaterThanOrEqualTo,
    dynamic arrayContains,
    List<dynamic>? arrayContainsAny,
    List<dynamic>? whereIn,
    List<dynamic>? whereNotIn,
    bool? isNull,
  }) {
    final newFilters = List<_WhereClause>.from(_filters);

    if (isEqualTo != null) {
      newFilters.add(_WhereClause(field, _WhereOp.isEqualTo, isEqualTo));
    }
    if (isLessThan != null) {
      newFilters.add(_WhereClause(field, _WhereOp.isLessThan, isLessThan));
    }
    if (isLessThanOrEqualTo != null) {
      newFilters.add(_WhereClause(field, _WhereOp.isLessThanOrEqualTo, isLessThanOrEqualTo));
    }
    if (isGreaterThan != null) {
      newFilters.add(_WhereClause(field, _WhereOp.isGreaterThan, isGreaterThan));
    }
    if (isGreaterThanOrEqualTo != null) {
      newFilters.add(_WhereClause(field, _WhereOp.isGreaterThanOrEqualTo, isGreaterThanOrEqualTo));
    }
    if (arrayContains != null) {
      newFilters.add(_WhereClause(field, _WhereOp.arrayContains, arrayContains));
    }
    if (arrayContainsAny != null) {
      newFilters.add(_WhereClause(field, _WhereOp.arrayContainsAny, arrayContainsAny));
    }
    if (whereIn != null) {
      newFilters.add(_WhereClause(field, _WhereOp.whereIn, whereIn));
    }
    if (whereNotIn != null) {
      newFilters.add(_WhereClause(field, _WhereOp.whereNotIn, whereNotIn));
    }
    if (isNull != null) {
      newFilters.add(_WhereClause(field, _WhereOp.isNull, isNull));
    }

    return _copyWith(filters: newFilters);
  }

  Query<T> orderBy(String field, {bool descending = false}) {
    return _copyWith(orders: [..._orders, _OrderClause(field, descending)]);
  }

  Query<T> limit(int count) {
    return _copyWith(limitCount: count);
  }

  Query<T> startAfterDocument(DocumentSnapshot? doc) {
    return _copyWith(startAfterDoc: doc);
  }

  // ─────────────── Execution: .get() — server-side via PostgREST ────────

  Future<QuerySnapshot<T>> get() async {
    try {
      dynamic query = SupaFlow.client.from(path).select();

      // Apply filters server-side
      for (final f in _filters) {
        final val = _serializeValue(f.value);
        switch (f.op) {
          case _WhereOp.isEqualTo:
            query = query.eq(f.field, val);
            break;
          case _WhereOp.isLessThan:
            query = query.lt(f.field, val);
            break;
          case _WhereOp.isLessThanOrEqualTo:
            query = query.lte(f.field, val);
            break;
          case _WhereOp.isGreaterThan:
            query = query.gt(f.field, val);
            break;
          case _WhereOp.isGreaterThanOrEqualTo:
            query = query.gte(f.field, val);
            break;
          case _WhereOp.whereIn:
            if (val != null && (val as List).isNotEmpty) {
              query = query.inFilter(f.field, val);
            }
            break;
          case _WhereOp.whereNotIn:
            // PostgREST: use .not().inFilter() pattern
            if (val != null && (val as List).isNotEmpty) {
              query = query.not(f.field, 'in', '(${val.join(",")})');
            }
            break;
          case _WhereOp.isNull:
            if (f.value == true) {
              query = query.isFilter(f.field, null);
            } else {
              query = query.not(f.field, 'is', null);
            }
            break;
          case _WhereOp.arrayContains:
            query = query.contains(f.field, [val]);
            break;
          case _WhereOp.arrayContainsAny:
            if (val != null && (val as List).isNotEmpty) {
              query = query.overlaps(f.field, val);
            }
            break;
        }
      }

      // Apply ordering
      for (final o in _orders) {
        query = query.order(o.field, ascending: !o.descending);
      }

      // Apply limit
      if (_limitCount != null && _limitCount! > 0) {
        query = query.limit(_limitCount!);
      }

      final List<dynamic> response = await query;

      return QuerySnapshot<T>(
        response
            .whereType<Map<String, dynamic>>()
            .map((r) => DocumentSnapshot<T>(
                  r['id']?.toString() ?? '',
                  r as T,
                  DocumentReference<T>(path, r['id']?.toString() ?? ''),
                ))
            .toList(),
      );
    } catch (e) {
      debugPrint('Query.get() error on $path: $e');
      return QuerySnapshot<T>([]);
    }
  }

  // ─────────────── Execution: .snapshots() — realtime with in-memory ─────

  Stream<QuerySnapshot<T>> snapshots() {
    try {
      final rawStream = SupaFlow.client.from(path).stream(primaryKey: ['id']);

      // Optimization: push the first isEqualTo filter to Supabase Realtime.
      // .eq() returns a different type (SupabaseStreamFilterBuilder), so we
      // branch into a Stream<List<Map>> to unify both paths.
      _WhereClause? pushedFilter;
      Stream<List<Map<String, dynamic>>> dataStream;

      final firstEq = _filters.cast<_WhereClause?>().firstWhere(
        (f) => f!.op == _WhereOp.isEqualTo && f.value != null &&
               (_serializeValue(f.value) is String || _serializeValue(f.value) is num),
        orElse: () => null,
      );

      if (firstEq != null) {
        final val = _serializeValue(firstEq.value);
        dataStream = rawStream.eq(firstEq.field, val);
        pushedFilter = firstEq;
      } else {
        dataStream = rawStream;
      }

      return dataStream.map((rows) {
        var results = rows.toList();

        // Apply remaining where clauses in-memory
        for (final f in _filters) {
          if (identical(f, pushedFilter)) continue; // already pushed server-side
          results = results.where((row) => _matchesClause(row, f)).toList();
        }

        // Apply ordering in-memory
        if (_orders.isNotEmpty) {
          results.sort((a, b) {
            for (final o in _orders) {
              final valA = _extractField(a, o.field);
              final valB = _extractField(b, o.field);
              final cmp = _compareValues(valA, valB);
              if (cmp != 0) return o.descending ? -cmp : cmp;
            }
            return 0;
          });
        }

        // Apply limit
        if (_limitCount != null && _limitCount! > 0 && results.length > _limitCount!) {
          results = results.sublist(0, _limitCount!);
        }

        // Apply startAfterDocument (cursor-based pagination)
        if (_startAfterDoc != null) {
          final startId = _startAfterDoc!.id;
          final idx = results.indexWhere((r) => r['id']?.toString() == startId);
          if (idx >= 0 && idx + 1 < results.length) {
            results = results.sublist(idx + 1);
          } else if (idx >= 0) {
            results = [];
          }
        }

        return QuerySnapshot<T>(
          results
              .map((r) => DocumentSnapshot<T>(
                    r['id']?.toString() ?? '',
                    r as T,
                    DocumentReference<T>(path, r['id']?.toString() ?? ''),
                  ))
              .toList(),
        );
      });
    } catch (e) {
      debugPrint('Query.snapshots() error on $path: $e');
      return Stream.value(QuerySnapshot<T>([]));
    }
  }

  // ─────────────── count() ───────────────────────────────────────────────

  AggregateQuery count() => AggregateQuery(path, _filters);
}

// ─────────────────────── QueryDocumentSnapshot ────────────────────────────

class QueryDocumentSnapshot<T> extends DocumentSnapshot<T> {
  QueryDocumentSnapshot(super.id, super.data, super.reference);
}

// ─────────────────────── AggregateQuery (with filters) ───────────────────

class AggregateQuery {
  final String path;
  final List<_WhereClause> _filters;

  AggregateQuery(this.path, this._filters);

  Future<AggregateQuerySnapshot> get() async {
    try {
      dynamic query = SupaFlow.client.from(path).select();

      for (final f in _filters) {
        final val = _serializeValue(f.value);
        switch (f.op) {
          case _WhereOp.isEqualTo:
            query = query.eq(f.field, val);
            break;
          case _WhereOp.isGreaterThan:
            query = query.gt(f.field, val);
            break;
          case _WhereOp.isLessThan:
            query = query.lt(f.field, val);
            break;
          case _WhereOp.whereIn:
            if (val != null && (val as List).isNotEmpty) {
              query = query.inFilter(f.field, val);
            }
            break;
          default:
            break; // Other ops are rare in count contexts
        }
      }

      final List<dynamic> response = await query;
      return AggregateQuerySnapshot(response.length);
    } catch (e) {
      debugPrint('AggregateQuery.get() error on $path: $e');
      return AggregateQuerySnapshot(0);
    }
  }
}

class AggregateQuerySnapshot {
  final int? count;
  AggregateQuerySnapshot(this.count);
}

// ─────────────────────────── QuerySnapshot ────────────────────────────────

class QuerySnapshot<T> {
  final List<DocumentSnapshot<T>> docs;
  QuerySnapshot(this.docs);
}

// ─────────────────────────── Filter ──────────────────────────────────────

class Filter {
  Filter(String field, {List<dynamic>? whereIn, List<dynamic>? arrayContainsAny});
}

// ─────────────────────────── FieldValue ──────────────────────────────────

/// Base class for FieldValue sentinels, which are detected by update().
abstract class _FieldValueSentinel {}

class _FieldValueIncrement extends _FieldValueSentinel {
  final num value;
  _FieldValueIncrement(this.value);
}

class _FieldValueArrayUnion extends _FieldValueSentinel {
  final List<dynamic> elements;
  _FieldValueArrayUnion(this.elements);
}

class _FieldValueArrayRemove extends _FieldValueSentinel {
  final List<dynamic> elements;
  _FieldValueArrayRemove(this.elements);
}

class _FieldValueServerTimestamp extends _FieldValueSentinel {}

class _FieldValueDelete {}

class FieldValue {
  static dynamic increment(num value) => _FieldValueIncrement(value);
  static dynamic arrayUnion(List<dynamic> elements) => _FieldValueArrayUnion(elements);
  static dynamic arrayRemove(List<dynamic> elements) => _FieldValueArrayRemove(elements);
  static dynamic serverTimestamp() => _FieldValueServerTimestamp();
  static dynamic delete() => _FieldValueDelete();
}
