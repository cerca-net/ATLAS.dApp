import '/flutter_flow/flutter_flow_icon_button.dart';
import '/flutter_flow/flutter_flow_theme.dart';
import '/flutter_flow/flutter_flow_util.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'connection_settings_model.dart';
import '../../../app_state.dart';

class ConnectionSettingsWidget extends StatefulWidget {
  const ConnectionSettingsWidget({super.key});

  static String routeName = 'ConnectionSettings';
  static String routePath = '/connectionSettings';

  @override
  State<ConnectionSettingsWidget> createState() => _ConnectionSettingsWidgetState();
}

class _ConnectionSettingsWidgetState extends State<ConnectionSettingsWidget> {
  late ConnectionSettingsModel _model;
  final scaffoldKey = GlobalKey<ScaffoldState>();

  @override
  void initState() {
    super.initState();
    _model = createModel(context, () => ConnectionSettingsModel());
    
    // Load values from FFAppState
    final appState = FFAppState();
    _model.remoteUrlController.text = appState.remoteNodeUrl;
    _model.localUrlController.text = appState.localNodeUrl;
    _model.apiKeyController.text = appState.adminApiKey;

    WidgetsBinding.instance.addPostFrameCallback((_) => safeSetState(() {}));
  }

  @override
  void dispose() {
    _model.dispose();
    super.dispose();
  }

  void _saveSettings() {
    final appState = FFAppState();
    appState.update(() {
      appState.remoteNodeUrl = _model.remoteUrlController.text;
      appState.localNodeUrl = _model.localUrlController.text;
      appState.adminApiKey = _model.apiKeyController.text;
    });

    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Connection settings saved successfully!'),
        backgroundColor: Colors.green,
      ),
    );
    context.safePop();
  }

  @override
  Widget build(BuildContext context) {
    final appState = FFAppState();
    
    return Scaffold(
      key: scaffoldKey,
      backgroundColor: FlutterFlowTheme.of(context).secondaryBackground,
      appBar: AppBar(
        backgroundColor: FlutterFlowTheme.of(context).secondaryBackground,
        automaticallyImplyLeading: false,
        leading: FlutterFlowIconButton(
          borderColor: Colors.transparent,
          borderRadius: 30.0,
          borderWidth: 1.0,
          buttonSize: 60.0,
          icon: Icon(
            Icons.arrow_back_rounded,
            color: FlutterFlowTheme.of(context).primaryText,
            size: 30.0,
          ),
          onPressed: () async {
            context.safePop();
          },
        ),
        title: Text(
          'Connection Settings',
          style: FlutterFlowTheme.of(context).titleMedium.override(
                font: GoogleFonts.inter(),
                color: FlutterFlowTheme.of(context).primaryText,
                fontWeight: FontWeight.bold,
              ),
        ),
        actions: [
          Padding(
            padding: const EdgeInsetsDirectional.fromSTEB(0.0, 8.0, 12.0, 8.0),
            child: FlutterFlowIconButton(
              borderRadius: 12.0,
              buttonSize: 40.0,
              fillColor: FlutterFlowTheme.of(context).secondaryBackground,
              icon: Icon(
                Icons.check_rounded,
                color: FlutterFlowTheme.of(context).primary,
                size: 24.0,
              ),
              onPressed: _saveSettings,
            ),
          ),
        ],
        centerTitle: false,
        elevation: 0.0,
      ),
      body: SafeArea(
        top: true,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.max,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Padding(
                padding: const EdgeInsetsDirectional.fromSTEB(24.0, 16.0, 24.0, 0.0),
                child: Text(
                  'Configure how the Cerca client connects to the ATLAS blockchain network.',
                  style: FlutterFlowTheme.of(context).bodyMedium.override(
                        font: GoogleFonts.inter(),
                        color: FlutterFlowTheme.of(context).secondaryText,
                      ),
                ),
              ),
              Padding(
                padding: const EdgeInsetsDirectional.fromSTEB(0.0, 12.0, 0.0, 0.0),
                child: Material(
                  color: Colors.transparent,
                  child: SwitchListTile.adaptive(
                    value: appState.isLocalNodeMode,
                    onChanged: (newValue) async {
                      appState.update(() {
                        appState.isLocalNodeMode = newValue;
                      });
                    },
                    title: Text(
                      'Local Node Mode',
                      style: FlutterFlowTheme.of(context).bodyLarge.override(
                            font: GoogleFonts.inter(
                              fontWeight: FlutterFlowTheme.of(context).bodyLarge.fontWeight,
                            ),
                          ),
                    ),
                    subtitle: Text(
                      appState.isLocalNodeMode 
                          ? 'Connected to local sidecar node (localhost:8081). Full validator controls enabled.'
                          : 'Connected to public seed node. Node controls are restricted.',
                      style: FlutterFlowTheme.of(context).bodyMedium.override(
                            font: GoogleFonts.inter(),
                            color: const Color(0xFF8B97A2),
                          ),
                    ),
                    tileColor: FlutterFlowTheme.of(context).secondaryBackground,
                    activeColor: FlutterFlowTheme.of(context).primary,
                    activeTrackColor: FlutterFlowTheme.of(context).accent1,
                    dense: false,
                    controlAffinity: ListTileControlAffinity.trailing,
                    contentPadding: const EdgeInsetsDirectional.fromSTEB(24.0, 12.0, 24.0, 12.0),
                  ),
                ),
              ),
              const Divider(height: 1.0, thickness: 1.0, color: Colors.white12),
              
              // Inputs section
              Padding(
                padding: const EdgeInsetsDirectional.fromSTEB(24.0, 24.0, 24.0, 24.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'NODE URL CONFIGURATION',
                      style: FlutterFlowTheme.of(context).bodySmall.override(
                            font: GoogleFonts.inter(),
                            color: FlutterFlowTheme.of(context).secondaryText,
                            fontWeight: FontWeight.bold,
                          ),
                    ),
                    const SizedBox(height: 16.0),
                    
                    // Remote URL Input
                    TextFormField(
                      controller: _model.remoteUrlController,
                      enabled: !appState.isLocalNodeMode,
                      decoration: InputDecoration(
                        labelText: 'Remote Seed Node URL',
                        labelStyle: FlutterFlowTheme.of(context).bodyMedium,
                        hintText: 'http://localhost:8080',
                        hintStyle: FlutterFlowTheme.of(context).bodyMedium.override(
                              font: GoogleFonts.inter(),
                              color: FlutterFlowTheme.of(context).secondaryText,
                            ),
                        enabledBorder: OutlineInputBorder(
                          borderSide: const BorderSide(color: Colors.white12, width: 1.0),
                          borderRadius: BorderRadius.circular(8.0),
                        ),
                        focusedBorder: OutlineInputBorder(
                          borderSide: BorderSide(color: FlutterFlowTheme.of(context).primary, width: 2.0),
                          borderRadius: BorderRadius.circular(8.0),
                        ),
                        disabledBorder: OutlineInputBorder(
                          borderSide: const BorderSide(color: Colors.white12, width: 1.0),
                          borderRadius: BorderRadius.circular(8.0),
                        ),
                        filled: true,
                        fillColor: appState.isLocalNodeMode 
                            ? const Color(0xFF13131A)
                            : FlutterFlowTheme.of(context).primaryBackground,
                      ),
                      style: FlutterFlowTheme.of(context).bodyMedium,
                    ),
                    const SizedBox(height: 20.0),
                    
                    // Local URL Input
                    TextFormField(
                      controller: _model.localUrlController,
                      enabled: appState.isLocalNodeMode,
                      decoration: InputDecoration(
                        labelText: 'Local Node URL',
                        labelStyle: FlutterFlowTheme.of(context).bodyMedium,
                        hintText: 'http://localhost:8081',
                        hintStyle: FlutterFlowTheme.of(context).bodyMedium.override(
                              font: GoogleFonts.inter(),
                              color: FlutterFlowTheme.of(context).secondaryText,
                            ),
                        enabledBorder: OutlineInputBorder(
                          borderSide: const BorderSide(color: Colors.white12, width: 1.0),
                          borderRadius: BorderRadius.circular(8.0),
                        ),
                        focusedBorder: OutlineInputBorder(
                          borderSide: BorderSide(color: FlutterFlowTheme.of(context).primary, width: 2.0),
                          borderRadius: BorderRadius.circular(8.0),
                        ),
                        disabledBorder: OutlineInputBorder(
                          borderSide: const BorderSide(color: Colors.white12, width: 1.0),
                          borderRadius: BorderRadius.circular(8.0),
                        ),
                        filled: true,
                        fillColor: !appState.isLocalNodeMode 
                            ? const Color(0xFF13131A)
                            : FlutterFlowTheme.of(context).primaryBackground,
                      ),
                      style: FlutterFlowTheme.of(context).bodyMedium,
                    ),
                    const SizedBox(height: 20.0),
                    
                    // Admin API Key Input
                    TextFormField(
                      controller: _model.apiKeyController,
                      decoration: InputDecoration(
                        labelText: 'Admin API Key (Protected endpoints)',
                        labelStyle: FlutterFlowTheme.of(context).bodyMedium,
                        hintText: 'cerca-dev-admin-secret-key',
                        hintStyle: FlutterFlowTheme.of(context).bodyMedium.override(
                              font: GoogleFonts.inter(),
                              color: FlutterFlowTheme.of(context).secondaryText,
                            ),
                        enabledBorder: OutlineInputBorder(
                          borderSide: const BorderSide(color: Colors.white12, width: 1.0),
                          borderRadius: BorderRadius.circular(8.0),
                        ),
                        focusedBorder: OutlineInputBorder(
                          borderSide: BorderSide(color: FlutterFlowTheme.of(context).primary, width: 2.0),
                          borderRadius: BorderRadius.circular(8.0),
                        ),
                        filled: true,
                        fillColor: FlutterFlowTheme.of(context).primaryBackground,
                      ),
                      style: FlutterFlowTheme.of(context).bodyMedium,
                      obscureText: true,
                    ),
                    
                    const SizedBox(height: 32.0),
                    ElevatedButton(
                      onPressed: _saveSettings,
                      style: ElevatedButton.styleFrom(
                        backgroundColor: FlutterFlowTheme.of(context).primary,
                        minimumSize: const Size(double.infinity, 50.0),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8.0),
                        ),
                      ),
                      child: Text(
                        'Save Configuration',
                        style: FlutterFlowTheme.of(context).titleMedium.override(
                              font: GoogleFonts.inter(),
                              color: Colors.white,
                              fontWeight: FontWeight.bold,
                            ),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
