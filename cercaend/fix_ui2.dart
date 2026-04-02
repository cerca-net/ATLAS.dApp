import 'dart:io';

void main() async {
  final file = File(r"c:\Users\beatr\Desktop\ATLAS\cercachain-fix-userpage-syntax-errors\cercachain-fix-userpage-syntax-errors\cercaend\lib\mainpages\orderpage\orderpage_widget.dart");
  final lines = await file.readAsLines();

  List<String> newLines = [];
  bool inBlock1 = false;
  bool inBlock2 = false;
  
  // The replacement builder payload:
  List<String> getBuilderLines(int indentLength, bool isBlock1) {
    final indent = ' ' * indentLength;
    final condition = isBlock1 
        ? "if (columnOrderRecord.paymentMethod != 'Data Unit' && columnOrderRecord.walletMethod != null)" 
        : "if (columnOrderRecord.walletMethod != null)";
        
    return """$indent$condition
$indent  FutureBuilder<WalletMethodsRecord>(
$indent    future: WalletMethodsRecord.getDocumentOnce(columnOrderRecord.walletMethod!),
$indent    builder: (context, snapshot) {
$indent      if (!snapshot.hasData) {
$indent        return const SizedBox.shrink();
$indent      }
$indent      final walletMethod = snapshot.data!;
$indent      return Column(
$indent        crossAxisAlignment: CrossAxisAlignment.stretch,
$indent        children: [
$indent          Text(
$indent            walletMethod.methodName,
$indent            style: FlutterFlowTheme.of(context).titleSmall.override(
$indent                  fontFamily: 'Montserrat',
$indent                  fontWeight: FlutterFlowTheme.of(context).titleSmall.fontWeight,
$indent                  color: FlutterFlowTheme.of(context).secondaryText,
$indent                ),
$indent          ),
$indent          Text(
$indent            walletMethod.methodType,
$indent            style: FlutterFlowTheme.of(context).bodySmall.override(
$indent                  fontFamily: 'Inter',
$indent                  fontWeight: FontWeight.w500,
$indent                  color: FlutterFlowTheme.of(context).secondaryText,
$indent                ),
$indent          ),
$indent          Row(
$indent            mainAxisAlignment: MainAxisAlignment.spaceBetween,
$indent            children: [
$indent              RichText(
$indent                textScaler: MediaQuery.of(context).textScaler,
$indent                text: TextSpan(
$indent                  children: [
$indent                    TextSpan(
$indent                      text: '# I.D. : ',
$indent                      style: FlutterFlowTheme.of(context).labelMedium,
$indent                    ),
$indent                    TextSpan(
$indent                      text: walletMethod.methodId,
$indent                      style: const TextStyle(),
$indent                    )
$indent                  ],
$indent                  style: FlutterFlowTheme.of(context).labelMedium,
$indent                ),
$indent              ),
$indent              Icon(Icons.content_copy_rounded, color: FlutterFlowTheme.of(context).secondaryText, size: 18.0),
$indent            ],
$indent          ),
$indent          Row(
$indent            mainAxisAlignment: MainAxisAlignment.spaceBetween,
$indent            children: [
$indent              RichText(
$indent                textScaler: MediaQuery.of(context).textScaler,
$indent                text: TextSpan(
$indent                  children: [
$indent                    TextSpan(
$indent                      text: '# Account : ',
$indent                      style: FlutterFlowTheme.of(context).labelMedium,
$indent                    ),
$indent                    TextSpan(
$indent                      text: walletMethod.methodAccount,
$indent                      style: const TextStyle(),
$indent                    )
$indent                  ],
$indent                  style: FlutterFlowTheme.of(context).labelMedium,
$indent                ),
$indent              ),
$indent              Icon(Icons.content_copy_rounded, color: FlutterFlowTheme.of(context).secondaryText, size: 18.0),
$indent            ],
$indent          ),
$indent        ].divide(const SizedBox(height: 4.0)),
$indent      );
$indent    },
$indent  ),""".split('\n');
  }

  for (int i = 0; i < lines.length; i++) {
    final line = lines[i];

    if (!inBlock1 && !inBlock2) {
      if (line.contains("if (columnOrderRecord") && 
          i + 2 < lines.length && 
          lines[i+1].contains(".paymentMethod !=") && 
          lines[i+2].contains("'Data Unit') ...[")) {
        inBlock1 = true;
        
        final indentLen = line.indexOf('if');
        newLines.addAll(getBuilderLines(indentLen, true));
        i += 2; // skip the condition lines
        continue;
      }
      
      if (line.contains("if (columnOrderRecord") && 
          i + 2 < lines.length && 
          lines[i+1].contains(".walletMethod !=") && 
          lines[i+2].contains("'Token') ...[")) {
        inBlock2 = true;
        
        final indentLen = line.indexOf('if');
        newLines.addAll(getBuilderLines(indentLen, false));
        i += 2; // skip the condition lines
        continue;
      }
      
      newLines.add(line);
    } else {
      // Find the end of the block.
      // Both blocks end similarly:
      //                                                                 ],
      //                                                               ).divide(
      // Wait, let's just count brackets or see the `.divide`
      if (line.contains("].divide(")) {
        inBlock1 = false;
        inBlock2 = false;
        // The `.divide` might be part of the column we injected around it, but our injected `children: [` doesn't have it yet. Wait...
        // My getBuilderLines puts `].divide(const SizedBox(height: 4.0)),` inside the FutureBuilder return Column.
        // So we can just skip down through this `].divide` and its trailing lines!
        // The original divide looks like:
        //      ].divide(
        //          const SizedBox(
        //              height:
        //                  4.0)),
      }
      // If we see `4.0)),` that's the end of my original `.divide(const SizedBox(height: 4.0))`!
      if (line.contains("4.0)),")) {
        inBlock1 = false;
        inBlock2 = false;
      }
    }
  }

  await file.writeAsString(newLines.join('\n'));
}
