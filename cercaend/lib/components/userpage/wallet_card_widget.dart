import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../app_state.dart';
import '../../flutter_flow/flutter_flow_theme.dart';
import '../../flutter_flow/flutter_flow_util.dart';

class WalletCardWidget extends StatelessWidget {
  final VoidCallback onSend;
  final VoidCallback onReceive;
  final VoidCallback onFaucet;
  final VoidCallback onHistory;

  const WalletCardWidget({
    super.key,
    required this.onSend,
    required this.onReceive,
    required this.onFaucet,
    required this.onHistory,
  });

  Widget _buildActionButton(
      BuildContext context, String label, IconData icon, VoidCallback onTap) {
    return Column(
      children: [
        InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(12),
          child: Container(
            width: 50,
            height: 50,
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.2),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Icon(icon, color: Colors.white, size: 24),
          ),
        ),
        const SizedBox(height: 8),
        Text(
          label,
          style: const TextStyle(color: Colors.white70, fontSize: 12),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(20.0),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [
            FlutterFlowTheme.of(context).primary,
            FlutterFlowTheme.of(context).tertiary,
          ],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(16.0),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.2),
            blurRadius: 10,
            offset: const Offset(0, 4),
          )
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Icon(Icons.account_balance_wallet,
                  color: Colors.white, size: 28),
              Text(
                'Personal Wallet',
                style: GoogleFonts.outfit(color: Colors.white70, fontSize: 14),
              ),
            ],
          ),
          const SizedBox(height: 20),
          Text(
            'Available Balance',
            style: GoogleFonts.outfit(color: Colors.white70, fontSize: 12),
          ),
          Text(
            '${formatNumber(FFAppState().walletBalance, formatType: FormatType.decimal, decimalType: DecimalType.automatic)} TCOIN',
            style: GoogleFonts.outfit(
              color: Colors.white,
              fontSize: 36,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 10),
          Row(
            children: [
              InkWell(
                onTap: () async {
                  await Clipboard.setData(
                      ClipboardData(text: FFAppState().walletAddress));

                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(
                        content: Text('Address copied to clipboard'),
                        duration: Duration(seconds: 1)),
                  );
                },
                child: const Icon(Icons.copy, color: Colors.white54, size: 14),
              ),
              const SizedBox(width: 4),
              Expanded(
                child: Text(
                  FFAppState().walletAddress,
                  style: GoogleFonts.robotoMono(
                      color: Colors.white54, fontSize: 12),
                  overflow: TextOverflow.ellipsis,
                ),
              ),
            ],
          ),
          const SizedBox(height: 20),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              _buildActionButton(context, 'Send', Icons.arrow_outward_rounded, onSend),
              _buildActionButton(context, 'Receive', Icons.arrow_downward_rounded, onReceive),
              _buildActionButton(context, 'Faucet', Icons.water_drop, onFaucet),
              _buildActionButton(context, 'History', Icons.history, onHistory),
            ],
          ),
        ],
      ),
    );
  }
}
