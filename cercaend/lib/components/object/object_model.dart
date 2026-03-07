import '/backend/backend.dart';
import '/flutter_flow/flutter_flow_util.dart';
import '/flutter_flow/form_field_controller.dart';
import '/flutter_flow/request_manager.dart';

import 'object_widget.dart' show ObjectWidget;
import 'package:flutter/material.dart';

class ObjectModel extends FlutterFlowModel<ObjectWidget> {
  ///  Local state fields for this component.

  bool iteminfo = false;

  DocumentReference? vote;

  SubmissionRecord? pin;

  DocumentReference? user;

  bool up = true;

  bool down = true;

  ///  State fields for stateful widgets in this component.

  // State field(s) for OImages widget.
  PageController? oImagesController;

  int get oImagesCurrentIndex => oImagesController != null &&
          oImagesController!.hasClients &&
          oImagesController!.page != null
      ? oImagesController!.page!.round()
      : 0;
  // State field(s) for ChoiceChips widget.
  FormFieldController<List<String>>? choiceChipsValueController;
  String? get choiceChipsValue =>
      choiceChipsValueController?.value?.firstOrNull;
  set choiceChipsValue(String? val) =>
      choiceChipsValueController?.value = val != null ? [val] : [];

  /// Query cache managers for this widget.

  final _objectinformationManager = StreamRequestManager<List<UsersRecord>>();
  Stream<List<UsersRecord>> objectinformation({
    String? uniqueQueryKey,
    bool? overrideCache,
    required Stream<List<UsersRecord>> Function() requestFn,
  }) =>
      _objectinformationManager.performRequest(
        uniqueQueryKey: uniqueQueryKey,
        overrideCache: overrideCache,
        requestFn: requestFn,
      );
  void clearObjectinformationCache() => _objectinformationManager.clear();
  void clearObjectinformationCacheKey(String? uniqueKey) =>
      _objectinformationManager.clearRequest(uniqueKey);

  final _objectdataManager = StreamRequestManager<SubmissionRecord>();
  Stream<SubmissionRecord> objectdata({
    String? uniqueQueryKey,
    bool? overrideCache,
    required Stream<SubmissionRecord> Function() requestFn,
  }) =>
      _objectdataManager.performRequest(
        uniqueQueryKey: uniqueQueryKey,
        overrideCache: overrideCache,
        requestFn: requestFn,
      );
  void clearObjectdataCache() => _objectdataManager.clear();
  void clearObjectdataCacheKey(String? uniqueKey) =>
      _objectdataManager.clearRequest(uniqueKey);

  @override
  void initState(BuildContext context) {}

  @override
  void dispose() {
    /// Dispose query cache managers for this widget.

    clearObjectinformationCache();

    clearObjectdataCache();
  }
}
