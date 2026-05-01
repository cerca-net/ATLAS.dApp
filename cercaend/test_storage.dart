import 'dart:io';
import 'package:http/http.dart' as http;

void main() async {
  final url = Uri.parse('https://epawttrarbrpzmdbmxyn.supabase.co/storage/v1/object/uploads/test_image.png');
  final anonKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImVwYXd0dHJhcmJycHptZGJteHluIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzY2MzA0NzMsImV4cCI6MjA5MjIwNjQ3M30.zhjusW5PGl4dRDsi38FtuFwCsfFw_wNSAG7oUuM_Dds';
  
  final response = await http.post(
    url,
    headers: {
      'apikey': anonKey,
      'Authorization': 'Bearer $anonKey',
      'Content-Type': 'image/png',
    },
    body: [137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 1], // Dummy PNG bytes
  );
  
  print('Storage Status Code: ${response.statusCode}');
  print('Storage Response: ${response.body}');
}
