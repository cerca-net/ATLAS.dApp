import 'dart:typed_data';

import 'package:mime_type/mime_type.dart';
import '/backend/supabase/supabase.dart';

/// Upload data to Supabase Storage (replaces Firebase Storage).
/// Uses the 'uploads' bucket by default.
Future<String?> uploadData(String path, Uint8List data) async {
  try {
    final storageBucket = SupaFlow.client.storage.from('uploads');
    await storageBucket.uploadBinary(
      path,
      data,
      fileOptions: FileOptions(
        contentType: mime(path),
        upsert: true,
      ),
    );
    return storageBucket.getPublicUrl(path);
  } catch (e) {
    print('Supabase Storage upload error for path "$path": $e');
    return null;
  }
}
