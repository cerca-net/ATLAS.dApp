import 'dart:io';

void main() async {
  final file = File(r"c:\Users\beatr\Desktop\ATLAS\cercachain-fix-userpage-syntax-errors\cercachain-fix-userpage-syntax-errors\cercaend\lib\mainpages\orderpage\orderpage_widget.dart");
  List<String> lines = await file.readAsLines();

  // The lines from 1487 to 1680 (0-indexed 1486 to 1679)
  // The lines from 4253 to 4425 (0-indexed 4252 to 4424)
  // Since we delete block 2 first to avoid shifting block 1's indices:

  const b2Start = 4252;
  const b2End = 4424; // inclusive index of the last line to remove.
  // Wait, let's verify line 4424!
  
  String builder2 = """                                                                if (columnOrderRecord.walletMethod != null && columnOrderRecord.paymentMethod != 'Data Unit')
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
                                                                              Icon(Icons.content_copy_rounded, color: FlutterFlowTheme.of(context).secondaryText, size: 18.0),
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
                                                                              Icon(Icons.content_copy_rounded, color: FlutterFlowTheme.of(context).secondaryText, size: 18.0),
                                                                          ],
                                                                        ),
                                                                      ].divide(const SizedBox(height: 4.0)),
                                                                    );
                                                                  },
                                                                ),""";

  lines.removeRange(b2Start, b2End + 1);
  lines.insert(b2Start, builder2);

  const b1Start = 1486;
  const b1End = 1679;
  lines.removeRange(b1Start, b1End + 1);
  lines.insert(b1Start, builder2);

  await file.writeAsString(lines.join('\n'));
}
