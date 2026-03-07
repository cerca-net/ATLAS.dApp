import '../database.dart';

class ProfilesTable extends SupabaseTable<ProfilesRow> {
  @override
  String get tableName => 'profiles';

  @override
  ProfilesRow createRow(Map<String, dynamic> data) => ProfilesRow(data);
}

class ProfilesRow extends SupabaseDataRow {
  ProfilesRow(super.data);

  @override
  SupabaseTable get table => ProfilesTable();

  int get id => getField<int>('id')!;
  set id(int value) => setField<int>('id', value);

  String? get walletAddress => getField<String>('wallet_address');
  set walletAddress(String? value) => setField<String>('wallet_address', value);

  String? get userId => getField<String>('user_id');
  set userId(String? value) => setField<String>('user_id', value);

  DateTime? get createdAt => getField<DateTime>('created_at');
  set createdAt(DateTime? value) => setField<DateTime>('created_at', value);
}
