import '../database.dart';

class MoralisStreamsTable extends SupabaseTable<MoralisStreamsRow> {
  @override
  String get tableName => 'moralisStreams';

  @override
  MoralisStreamsRow createRow(Map<String, dynamic> data) =>
      MoralisStreamsRow(data);
}

class MoralisStreamsRow extends SupabaseDataRow {
  MoralisStreamsRow(super.data);

  @override
  SupabaseTable get table => MoralisStreamsTable();

  String get hash => getField<String>('hash')!;
  set hash(String value) => setField<String>('hash', value);

  String get value => getField<String>('value')!;
  set value(String value) => setField<String>('value', value);

  String? get from => getField<String>('from');
  set from(String? value) => setField<String>('from', value);

  String? get to => getField<String>('to');
  set to(String? value) => setField<String>('to', value);
}
