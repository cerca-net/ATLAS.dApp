import 'dart:io';

void main() async {
  final file = File(r"c:\Users\beatr\Desktop\ATLAS\cercachain-fix-userpage-syntax-errors\cercachain-fix-userpage-syntax-errors\cercaend\lib\mainpages\orderpage\orderpage_widget.dart");
  String text = await file.readAsString();

  String replaceBlock(String txt, String startStr, String endStr, String replacement) {
    int startIdx = txt.indexOf(startStr);
    if (startIdx == -1) {
      print("Could not find start String...");
      return txt;
    }
    int endIdx = txt.indexOf(endStr, startIdx);
    if (endIdx == -1) {
      print("Could not find end String...");
      return txt;
    }
    endIdx += endStr.length;
    return txt.substring(0, startIdx) + replacement + txt.substring(endIdx);
  }

  const startStr1 = '''                                                                if (columnOrderRecord
                                                                        .paymentMethod !=
                                                                    'Data Unit') ...[''';
  
  const endStr1 = '''                                                                  ],
                                                                ),
                                                                ],''';

  const replacement1 = '''                                                                if (columnOrderRecord.paymentMethod != 'Data Unit' && columnOrderRecord.walletMethod != null)
                                                                  FutureBuilder<WalletMethodsRecord>(
                                                                    future: WalletMethodsRecord.getDocumentOnce(columnOrderRecord.walletMethod!),
                                                                    builder: (context, snapshot) {
                                                                      if (!snapshot.hasData) {
                                                                        return const SizedBox.shrink();
                                                                      }
                                                                      final walletMethod = snapshot.data!;
                                                                      return Column(
                                                                        crossAxisAlignment: CrossAxisAlignment.stretch,
                                                                        children: [
                                                                          Text(
                                                                            walletMethod.methodName,
                                                                            style: FlutterFlowTheme.of(context).titleSmall.override(
                                                                                  fontFamily: 'Montserrat',
                                                                                  fontWeight: FlutterFlowTheme.of(context).titleSmall.fontWeight,
                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                ),
                                                                          ),
                                                                          Text(
                                                                            walletMethod.methodType,
                                                                            style: FlutterFlowTheme.of(context).bodySmall.override(
                                                                                  fontFamily: 'Inter',
                                                                                  fontWeight: FontWeight.w500,
                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                ),
                                                                          ),
                                                                          Row(
                                                                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                                                            children: [
                                                                              RichText(
                                                                                textScaler: MediaQuery.of(context).textScaler,
                                                                                text: TextSpan(
                                                                                  children: [
                                                                                    TextSpan(
                                                                                      text: '# I.D. : ',
                                                                                      style: FlutterFlowTheme.of(context).labelMedium,
                                                                                    ),
                                                                                    TextSpan(
                                                                                      text: walletMethod.methodId,
                                                                                      style: const TextStyle(),
                                                                                    )
                                                                                  ],
                                                                                  style: FlutterFlowTheme.of(context).labelMedium,
                                                                                ),
                                                                              ),
                                                                            ],
                                                                          ),
                                                                          Row(
                                                                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                                                            children: [
                                                                              RichText(
                                                                                textScaler: MediaQuery.of(context).textScaler,
                                                                                text: TextSpan(
                                                                                  children: [
                                                                                    TextSpan(
                                                                                      text: '# Account : ',
                                                                                      style: FlutterFlowTheme.of(context).labelMedium,
                                                                                    ),
                                                                                    TextSpan(
                                                                                      text: walletMethod.methodAccount,
                                                                                      style: const TextStyle(),
                                                                                    )
                                                                                  ],
                                                                                  style: FlutterFlowTheme.of(context).labelMedium,
                                                                                ),
                                                                              ),
                                                                            ],
                                                                          ),
                                                                        ].divide(const SizedBox(height: 4.0)),
                                                                      );
                                                                    },
                                                                  ),''';

  text = replaceBlock(text, startStr1, endStr1, replacement1);

  const startStr2 = '''                                                                if (columnOrderRecord
                                                                        .walletMethod !=
                                                                    'Token') ...[''';

  const replacement2 = '''                                                                if (columnOrderRecord.walletMethod != null)
                                                                  FutureBuilder<WalletMethodsRecord>(
                                                                    future: WalletMethodsRecord.getDocumentOnce(columnOrderRecord.walletMethod!),
                                                                    builder: (context, snapshot) {
                                                                      if (!snapshot.hasData) {
                                                                        return const SizedBox.shrink();
                                                                      }
                                                                      final walletMethod = snapshot.data!;
                                                                      return Column(
                                                                        crossAxisAlignment: CrossAxisAlignment.stretch,
                                                                        children: [
                                                                          Text(
                                                                            walletMethod.methodName,
                                                                            style: FlutterFlowTheme.of(context).titleSmall.override(
                                                                                  fontFamily: 'Montserrat',
                                                                                  fontWeight: FlutterFlowTheme.of(context).titleSmall.fontWeight,
                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                ),
                                                                          ),
                                                                          Text(
                                                                            walletMethod.methodType,
                                                                            style: FlutterFlowTheme.of(context).bodySmall.override(
                                                                                  fontFamily: 'Inter',
                                                                                  fontWeight: FontWeight.w500,
                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                ),
                                                                          ),
                                                                          Row(
                                                                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                                                            children: [
                                                                              RichText(
                                                                                textScaler: MediaQuery.of(context).textScaler,
                                                                                text: TextSpan(
                                                                                  children: [
                                                                                    TextSpan(
                                                                                      text: '# I.D. : ',
                                                                                      style: FlutterFlowTheme.of(context).labelMedium,
                                                                                    ),
                                                                                    TextSpan(
                                                                                      text: walletMethod.methodId,
                                                                                      style: const TextStyle(),
                                                                                    )
                                                                                  ],
                                                                                  style: FlutterFlowTheme.of(context).labelMedium,
                                                                                ),
                                                                              ),
                                                                            ],
                                                                          ),
                                                                          Row(
                                                                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                                                            children: [
                                                                              RichText(
                                                                                textScaler: MediaQuery.of(context).textScaler,
                                                                                text: TextSpan(
                                                                                  children: [
                                                                                    TextSpan(
                                                                                      text: '# Account : ',
                                                                                      style: FlutterFlowTheme.of(context).labelMedium,
                                                                                    ),
                                                                                    TextSpan(
                                                                                      text: walletMethod.methodAccount,
                                                                                      style: const TextStyle(),
                                                                                    )
                                                                                  ],
                                                                                  style: FlutterFlowTheme.of(context).labelMedium,
                                                                                ),
                                                                              ),
                                                                            ],
                                                                          ),
                                                                        ].divide(const SizedBox(height: 4.0)),
                                                                      );
                                                                    },
                                                                  ),''';

  text = replaceBlock(text, startStr2, endStr1, replacement2);

  await file.writeAsString(text);
  print("Replaced both blocks.");
}
