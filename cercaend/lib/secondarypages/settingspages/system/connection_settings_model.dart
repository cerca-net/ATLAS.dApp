import '/flutter_flow/flutter_flow_util.dart';
import 'connection_settings_widget.dart' show ConnectionSettingsWidget;
import 'package:flutter/material.dart';

class ConnectionSettingsModel extends FlutterFlowModel<ConnectionSettingsWidget> {
  /// State fields for stateful widgets in this page.
  final remoteUrlController = TextEditingController();
  final localUrlController = TextEditingController();
  final apiKeyController = TextEditingController();

  @override
  void initState(BuildContext context) {}

  @override
  void dispose() {
    remoteUrlController.dispose();
    localUrlController.dispose();
    apiKeyController.dispose();
  }
}
