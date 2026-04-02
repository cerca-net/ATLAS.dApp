import re
import os

filepath = r"c:\Users\beatr\Desktop\ATLAS\cercachain-fix-userpage-syntax-errors\cercachain-fix-userpage-syntax-errors\cercaend\lib\mainpages\orderpage\orderpage_widget.dart"

with open(filepath, 'r', encoding='utf-8') as f:
    text = f.read()

def replace_block(text, start_str, end_str, replacement):
    start_idx = text.find(start_str)
    if start_idx == -1:
        print(f"Could not find start: {start_str[:50]}")
        return text
    end_idx = text.find(end_str, start_idx)
    if end_idx == -1:
        print("Could not find end")
        return text
    end_idx += len(end_str)
    return text[:start_idx] + replacement + text[end_idx:]

start_str_1 = """                                                                if (columnOrderRecord
                                                                        .paymentMethod !=
                                                                    'Data Unit') ...["""
end_str_1 = """                                                                  ],
                                                                ),
                                                                ],"""

replacement_1 = """                                                                if (columnOrderRecord.paymentMethod != 'Data Unit' && columnOrderRecord.walletMethod != null)
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
                                                                                  font: GoogleFonts.montserrat(
                                                                                    fontWeight: FlutterFlowTheme.of(context).titleSmall.fontWeight,
                                                                                  ),
                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                ),
                                                                          ),
                                                                          Text(
                                                                            walletMethod.methodType,
                                                                            style: FlutterFlowTheme.of(context).bodySmall.override(
                                                                                  font: GoogleFonts.inter(
                                                                                    fontWeight: FontWeight.w500,
                                                                                  ),
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
                                                                  ),"""

text = replace_block(text, start_str_1, end_str_1, replacement_1)

# Now block 2
start_str_2 = """                                                                if (columnOrderRecord
                                                                        .walletMethod !=
                                                                    'Token') ...["""

# Wait, the compiler probably complained if I don't use 'Data Unit'. Let's replace 'Token' with a proper check or just duplicate logic. Let's find exactly what's there:
start_str_2_exact = """                                                                if (columnOrderRecord
                                                                        .walletMethod !=
                                                                    'Token') ...["""

replacement_2 = """                                                                if (columnOrderRecord.walletMethod != null)
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
                                                                                  font: GoogleFonts.montserrat(
                                                                                    fontWeight: FlutterFlowTheme.of(context).titleSmall.fontWeight,
                                                                                  ),
                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                ),
                                                                          ),
                                                                          Text(
                                                                            walletMethod.methodType,
                                                                            style: FlutterFlowTheme.of(context).bodySmall.override(
                                                                                  font: GoogleFonts.inter(
                                                                                    fontWeight: FontWeight.w500,
                                                                                  ),
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
                                                                  ),"""

text = replace_block(text, start_str_2_exact, end_str_1, replacement_2)

with open(filepath, 'w', encoding='utf-8') as f:
    f.write(text)
print("Done")
