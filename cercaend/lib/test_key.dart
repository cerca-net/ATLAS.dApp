import 'dart:math';
import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:hex/hex.dart';
import 'package:elliptic/elliptic.dart' as elliptic;

void main() {
  final mnemonic = "abandon ability able about above absent absorb abstract absurd abuse access accident";
  final mnemonicBytes = utf8.encode(mnemonic);
  final seedHash = sha256.convert(mnemonicBytes).bytes;
  
  try {
    final privateKeyBytes = BigInt.parse(HEX.encode(seedHash.sublist(0, 32)), radix: 16);
    final curve = elliptic.getP256();
    final privateKey = elliptic.PrivateKey(curve, privateKeyBytes);
    final publicKey = privateKey.publicKey;
    
    final xBytes = publicKey.X.toRadixString(16).padLeft(64, '0');
    final yBytes = publicKey.Y.toRadixString(16).padLeft(64, '0');
    
    print("Address logic succeeded");
  } catch(e) {
    print("Error: $e");
  }
}
