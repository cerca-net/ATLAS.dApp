import 'dart:io';

void main() async {
  final schemaDir = Directory('lib/backend/schema');
  final files = schemaDir.listSync().whereType<File>().where((f) => f.path.endsWith('_record.dart'));

  final sqlBuffer = StringBuffer();
  sqlBuffer.writeln('-- Supabase Schema Migration Script');
  sqlBuffer.writeln('-- Auto-generated to match CercaChain FlutterFlow schema\\n');

  for (final file in files) {
    final content = await file.readAsString();
    final filename = file.uri.pathSegments.last;
    if (filename == 'index.dart') continue;
    
    // Convert e.g., 'users_record.dart' -> 'users'
    final tableName = filename.replaceAll('_record.dart', '');
    
    sqlBuffer.writeln('CREATE TABLE IF NOT EXISTS public."$tableName" (');
    sqlBuffer.writeln('  "id" text PRIMARY KEY,');
    
    // Regex to match: // "field_name" field.
    final fieldRegex = RegExp(r'// "([^"]+)" field\.');
    final matches = fieldRegex.allMatches(content);
    
    final fields = <String>{};
    for (final match in matches) {
      final fieldName = match.group(1);
      if (fieldName != null && fieldName != 'id' && fieldName != 'reference') {
        fields.add(fieldName);
      }
    }
    
    // Add columns (defaulting to jsonb for complex types, text for simple things)
    // For simplicity we will set them as text, but might need typing.
    // jsonb covers most dynamic arrays/maps if we aren't perfectly sure.
    // Let's at least extract if it's an int, double, etc.
    // We can grep the getter `int get <name> =>`
    
    final fieldTypes = <String, String>{};
    for(final fieldName in fields) {
       // Look for `Type get getterName => _getterName`
       // This can be tricky due to dart naming, let's use a simpler heuristic.
       final lines = content.split('\\n');
       String dbType = 'text'; // Default
       for(int i = 0; i < lines.length; i++) {
          if (lines[i].contains('// "$fieldName" field.')) {
             if (i+1 < lines.length) {
                final typeLine = lines[i+1].trim();
                if (typeLine.startsWith('int')) dbType = 'integer';
                else if (typeLine.startsWith('double')) dbType = 'double precision';
                else if (typeLine.startsWith('bool')) dbType = 'boolean';
                else if (typeLine.startsWith('DateTime')) dbType = 'timestamp with time zone';
                else if (typeLine.startsWith('List')) dbType = 'jsonb';
                else if (typeLine.startsWith('DocumentReference')) dbType = 'text';
                else if (typeLine.startsWith('LatLng')) dbType = 'jsonb';
                break;
             }
          }
       }
       sqlBuffer.writeln('  "$fieldName" $dbType,');
    }
    
    sqlBuffer.writeln('  "created_at" timestamp with time zone DEFAULT timezone(\\\'utc\\\'::text, now())');
    sqlBuffer.writeln(');\\n');
  }

  print(sqlBuffer.toString());
}
