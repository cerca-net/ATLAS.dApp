import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:path/path.dart' as p;

class LocalNodeManager {
  static final LocalNodeManager _instance = LocalNodeManager._internal();
  factory LocalNodeManager() => _instance;
  LocalNodeManager._internal();

  Process? _nodeProcess;
  bool get isRunning => _nodeProcess != null;

  /// Starts the embedded Go Blockchain Node natively.
  /// Works only on Desktop platforms (Windows, macOS, Linux).
  Future<void> startEmbeddedNode() async {
    if (kIsWeb) {
      debugPrint('LocalNodeManager: Cannot embed node on Web architecture.');
      return;
    }

    if (Platform.isAndroid || Platform.isIOS) {
      debugPrint('LocalNodeManager: Mobile node embedding requires gomobile FFI bindings. Skipping for now.');
      return;
    }

    if (isRunning) {
      debugPrint('LocalNodeManager: Embedded Node is already running.');
      return;
    }

    try {
      debugPrint('LocalNodeManager: Spawning embedded blockchain node process...');
      
      // In a production built app, the executable would be bundled in the assets or same directory.
      // For development, we'll try to use 'go run' from the project root.
      // Note: We assign it to port 8001 / api 8081 specifically for the App User.
      
      // Determine paths. Assuming flutter app is inside /cercaend and Go backend is at /ATLAS.BC0.0.1
      String executable = 'go';
      List<String> arguments = [
        'run', 
        'cmd/main.go', 
        '--port', '8001', 
        '--api', '8081', 
        '--datadir', '.data_app'
      ];

      // Assuming the working directory is near the go project.
      // If deployed, this would be an executed binary: e.g., 'atlas_node.exe'
      // String executable = 'atlas_node.exe';
      // List<String> arguments = ['--port', '8001', '--api', '8081'];
      
      // Need to walk up one dir from /cercaend for development pathing
      String workingDir = p.join(Directory.current.path, '..', 'ATLAS.BC0.0.1');

      _nodeProcess = await Process.start(
        executable,
        arguments,
        workingDirectory: workingDir,
      );

      debugPrint('LocalNodeManager: Node spawned with PID \${_nodeProcess!.pid}');

      _nodeProcess!.stdout.listen((data) {
         // Optionally capture node stdout
         // debugPrint('NODE OUT: \${String.fromCharCodes(data)}');
      });

      _nodeProcess!.stderr.listen((data) {
         // Optionally capture node stderr
         // debugPrint('NODE ERR: \${String.fromCharCodes(data)}');
      });

      _nodeProcess!.exitCode.then((code) {
        debugPrint('LocalNodeManager: Embedded Node exited with code $code');
        _nodeProcess = null;
      });

    } catch (e) {
      debugPrint('LocalNodeManager: Failed to start embedded node: $e');
    }
  }

  /// Kills the embedded node process.
  Future<void> stopEmbeddedNode() async {
    if (_nodeProcess != null) {
      debugPrint('LocalNodeManager: Shutting down embedded node process...');
      _nodeProcess!.kill();
      _nodeProcess = null;
    }
  }
}
