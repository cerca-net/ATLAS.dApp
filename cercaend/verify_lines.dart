import 'dart:io';

void main() async {
  final file = File(r"c:\Users\beatr\Desktop\ATLAS\cercachain-fix-userpage-syntax-errors\cercachain-fix-userpage-syntax-errors\cercaend\lib\mainpages\orderpage\orderpage_widget.dart");
  List<String> lines = await file.readAsLines();

  print("Line 1484 (index 1483): " + lines[1483]);
  print("Line 1669 (index 1668): " + lines[1668]);
  
  print("Line 4240 (index 4239): " + lines[4239]);
  print("Line 4412 (index 4411): " + lines[4411]);
}
