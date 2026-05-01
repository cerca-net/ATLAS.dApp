import 'dart:convert';
import 'package:http/http.dart' as http;

void main() async {
  final url = Uri.parse('https://epawttrarbrpzmdbmxyn.supabase.co/rest/v1/catalogue');
  final anonKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImVwYXd0dHJhcmJycHptZGJteHluIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzY2MzA0NzMsImV4cCI6MjA5MjIwNjQ3M30.zhjusW5PGl4dRDsi38FtuFwCsfFw_wNSAG7oUuM_Dds';
  
  final response = await http.post(
    url,
    headers: {
      'apikey': anonKey,
      'Authorization': 'Bearer $anonKey',
      'Content-Type': 'application/json',
      'Prefer': 'return=representation',
    },
    body: json.encode({
      'id': 'test_script_id_1234567',
      'catalogue_buzzwords': ['test1', 'test2'],
    }),
  );
  
  print('Status Code: ${response.statusCode}');
  print('Response Body: ${response.body}');
}
