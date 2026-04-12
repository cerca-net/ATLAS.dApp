import 'dart:io';

void main() {
  final file = File('analyze_output_utf8.txt');
  if (!file.existsSync()) {
    print('analyze_output_utf8.txt not found');
    return;
  }

  final lines = file.readAsLinesSync();
  final outputFiles = <String, Set<int>>{};

  for (var line in lines) {
    if (line.contains('use_build_context_synchronously')) {
      final parts = line.split(' - ');
      if (parts.length >= 3) {
        final fileInfo = parts[1].trim();
        final fileParts = fileInfo.split(':');
        if (fileParts.length >= 3) {
          final filepath = fileParts[0];
          final lineNum = int.tryParse(fileParts[1]) ?? 0;
          if (lineNum > 0) {
            outputFiles.putIfAbsent(filepath, () => <int>{});
            outputFiles[filepath]!.add(lineNum - 1);
          }
        }
      }
    }
  }

  for (var entry in outputFiles.entries) {
    var filepath = entry.key;
    final dartFile = File(filepath);
    if (!dartFile.existsSync()) {
      print('File not found: $filepath');
      continue;
    }

    final dartLines = dartFile.readAsLinesSync();
    final lineNumsDesc = entry.value.toList()..sort((a, b) => b.compareTo(a));

    for (var ln in lineNumsDesc) {
      if (ln < dartLines.length) {
        final targetLine = dartLines[ln];
        final indentMatch = RegExp(r'^\s*').firstMatch(targetLine);
        final indentStr = indentMatch?.group(0) ?? '';

        bool alreadyHasMounted = false;
        for (var i = 1; i <= 3; i++) {
          if (ln - i >= 0) {
            final prevLine = dartLines[ln - i].trim();
            if (prevLine.contains('mounted')) {
              alreadyHasMounted = true;
              break;
            }
            if (prevLine.isNotEmpty && prevLine != '}' && prevLine != '{') {
              break;
            }
          }
        }

        if (!alreadyHasMounted) {
          var insertStr = '${indentStr}if (!context.mounted) return;';
          dartLines.insert(ln, insertStr);
        }
      }
    }

    dartFile.writeAsStringSync('${dartLines.join('\n')}\n');
    print('Fixed: $filepath');
  }
}
