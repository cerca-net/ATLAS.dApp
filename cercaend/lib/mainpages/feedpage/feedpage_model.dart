import '/backend/backend.dart';
import '/flutter_flow/flutter_flow_util.dart';
import '/flutter_flow/form_field_controller.dart';
import 'feedpage_widget.dart' show FeedpageWidget;
import 'package:flutter/material.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';

class FeedpageModel extends FlutterFlowModel<FeedpageWidget> {
  ///  Local state fields for this page.

  String? choiceChipsSelected;

  ///  State fields for stateful widgets in this page.

  final formKey = GlobalKey<FormState>();
  // State field(s) for StaggeredView widget.

  PagingController<DocumentSnapshot<Object?>?, SubmissionRecord>?
      staggeredViewPagingController;
  Query? staggeredViewPagingQuery;
  List<StreamSubscription?> staggeredViewStreamSubscriptions = [];

  // State field(s) for ChoiceChips widget.
  FormFieldController<List<String>>? choiceChipsValueController;
  List<String>? get choiceChipsValues => choiceChipsValueController?.value;
  set choiceChipsValues(List<String>? val) =>
      choiceChipsValueController?.value = val;
  // State field(s) for wordSearcher widget.
  FocusNode? wordSearcherFocusNode;
  TextEditingController? wordSearcherTextController;
  String? Function(BuildContext, String?)? wordSearcherTextControllerValidator;
  List<SubmissionRecord> simpleSearchResults = [];

  @override
  void initState(BuildContext context) {}

  @override
  void dispose() {
    for (var s in staggeredViewStreamSubscriptions) {
      s?.cancel();
    }
    staggeredViewPagingController?.dispose();

    wordSearcherFocusNode?.dispose();
    wordSearcherTextController?.dispose();
  }

  /// Additional helper methods.
  PagingController<DocumentSnapshot<Object?>?, SubmissionRecord>
      setStaggeredViewController(
    Query query, {
    DocumentReference<Object?>? parent,
  }) {
    staggeredViewPagingController ??=
        _createStaggeredViewController(query, parent);
    if (staggeredViewPagingQuery != query) {
      staggeredViewPagingQuery = query;
      staggeredViewPagingController?.refresh();
    }
    return staggeredViewPagingController!;
  }

  PagingController<DocumentSnapshot<Object?>?, SubmissionRecord>
      _createStaggeredViewController(
    Query query,
    DocumentReference<Object?>? parent,
  ) {
    final controller =
        PagingController<DocumentSnapshot<Object?>?, SubmissionRecord>(
            firstPageKey: null);
    return controller
      ..addPageRequestListener(
        (nextPageMarker) => querySubmissionRecordPage(
          queryBuilder: (_) => staggeredViewPagingQuery ??= query,
          nextPageMarker: nextPageMarker,
          streamSubscriptions: staggeredViewStreamSubscriptions,
          controller: controller,
          pageSize: 25,
          isStream: true,
        ),
      );
  }
}
