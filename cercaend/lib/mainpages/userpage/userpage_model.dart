import '/backend/backend.dart';
import '/flutter_flow/flutter_flow_util.dart';
import '/flutter_flow/form_field_controller.dart';
import '/flutter_flow/request_manager.dart';

import '/index.dart';
import 'userpage_widget.dart' show UserpageWidget;
import 'package:flutter/material.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import '/services/blockchain/blockchain_service.dart';

class UserpageModel extends FlutterFlowModel<UserpageWidget> {
  ///  Local state fields for this page.

  List<DocumentReference> totalAvgRef = [];
  void addToTotalAvgRef(DocumentReference item) => totalAvgRef.add(item);
  void removeFromTotalAvgRef(DocumentReference item) =>
      totalAvgRef.remove(item);
  void removeAtIndexFromTotalAvgRef(int index) => totalAvgRef.removeAt(index);
  void insertAtIndexInTotalAvgRef(int index, DocumentReference item) =>
      totalAvgRef.insert(index, item);
  void updateTotalAvgRefAtIndex(
          int index, Function(DocumentReference) updateFn) =>
      totalAvgRef[index] = updateFn(totalAvgRef[index]);

  bool trsIsSent = true;

  bool isWalletLoading = false;
  NetworkStatusResponse? networkStatus;
  NodeStatusResponse? nodeStatus;
  List<NodeLogEntry> nodeLogs = [];
  UserIdentity? socialIdentity;

  void addToNodeLogs(NodeLogEntry item) {
    nodeLogs.add(item);
    if (nodeLogs.length > 50) nodeLogs.removeAt(0);
  }

  List<int> ratingsList = [];
  void addToRatingsList(int item) => ratingsList.add(item);
  void removeFromRatingsList(int item) => ratingsList.remove(item);
  void removeAtIndexFromRatingsList(int index) => ratingsList.removeAt(index);
  void insertAtIndexInRatingsList(int index, int item) =>
      ratingsList.insert(index, item);
  void updateRatingsListAtIndex(int index, Function(int) updateFn) =>
      ratingsList[index] = updateFn(ratingsList[index]);

  bool webViewOpen = true;

  ///  State fields for stateful widgets in this page.

  // State field(s) for ChoiceChips widget.
  FormFieldController<List<String>>? choiceChipsValueController1;
  String? get choiceChipsValue1 =>
      choiceChipsValueController1?.value?.firstOrNull;
  set choiceChipsValue1(String? val) =>
      choiceChipsValueController1?.value = val != null ? [val] : [];
  // State field(s) for TabBar widget.
  TabController? tabBarController;
  int get tabBarCurrentIndex =>
      tabBarController != null ? tabBarController!.index : 0;
  int get tabBarPreviousIndex =>
      tabBarController != null ? tabBarController!.previousIndex : 0;

  // State field(s) for ChoiceChipsOrders widget.
  FormFieldController<List<String>>? choiceChipsOrdersValueController;
  String? get choiceChipsOrdersValue =>
      choiceChipsOrdersValueController?.value?.firstOrNull;
  set choiceChipsOrdersValue(String? val) =>
      choiceChipsOrdersValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsAnaytics widget.
  FormFieldController<List<String>>? choiceChipsAnayticsValueController;
  String? get choiceChipsAnayticsValue =>
      choiceChipsAnayticsValueController?.value?.firstOrNull;
  set choiceChipsAnayticsValue(String? val) =>
      choiceChipsAnayticsValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChips widget.
  FormFieldController<List<String>>? choiceChipsValueController2;
  String? get choiceChipsValue2 =>
      choiceChipsValueController2?.value?.firstOrNull;
  set choiceChipsValue2(String? val) =>
      choiceChipsValueController2?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPosts widget.
  FormFieldController<List<String>>? choiceChipsPostsValueController;
  List<String>? get choiceChipsPostsValues =>
      choiceChipsPostsValueController?.value;
  set choiceChipsPostsValues(List<String>? val) =>
      choiceChipsPostsValueController?.value = val;
  // State field(s) for StaggeredView widget.

  PagingController<DocumentSnapshot?, SubmissionRecord>?
      staggeredViewPagingController1;
  Query? staggeredViewPagingQuery1;
  List<StreamSubscription?> staggeredViewStreamSubscriptions1 = [];

  // State field(s) for ChoiceChipsItems widget.
  FormFieldController<List<String>>? choiceChipsItemsValueController;
  String? get choiceChipsItemsValue =>
      choiceChipsItemsValueController?.value?.firstOrNull;
  set choiceChipsItemsValue(String? val) =>
      choiceChipsItemsValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsWallet widget.
  FormFieldController<List<String>>? choiceChipsWalletValueController;
  String? get choiceChipsWalletValue =>
      choiceChipsWalletValueController?.value?.firstOrNull;
  set choiceChipsWalletValue(String? val) =>
      choiceChipsWalletValueController?.value = val != null ? [val] : [];

  // State field(s) for Validator Staking.
  TextEditingController? stakeAmountController;

  /// Query cache managers for this widget.

  final _itemPostsManager = FutureRequestManager<int>();
  Future<int> itemPosts({
    String? uniqueQueryKey,
    bool? overrideCache,
    required Future<int> Function() requestFn,
  }) =>
      _itemPostsManager.performRequest(
        uniqueQueryKey: uniqueQueryKey,
        overrideCache: overrideCache,
        requestFn: requestFn,
      );
  void clearItemPostsCache() => _itemPostsManager.clear();
  void clearItemPostsCacheKey(String? uniqueKey) =>
      _itemPostsManager.clearRequest(uniqueKey);

  final _itemsCountManager = FutureRequestManager<int>();
  Future<int> itemsCount({
    String? uniqueQueryKey,
    bool? overrideCache,
    required Future<int> Function() requestFn,
  }) =>
      _itemsCountManager.performRequest(
        uniqueQueryKey: uniqueQueryKey,
        overrideCache: overrideCache,
        requestFn: requestFn,
      );
  void clearItemsCountCache() => _itemsCountManager.clear();
  void clearItemsCountCacheKey(String? uniqueKey) =>
      _itemsCountManager.clearRequest(uniqueKey);

  @override
  void initState(BuildContext context) {
    stakeAmountController = TextEditingController();
  }

  @override
  void dispose() {
    tabBarController?.dispose();
    stakeAmountController?.dispose();
    for (var s in staggeredViewStreamSubscriptions1) {
      s?.cancel();
    }
    staggeredViewPagingController1?.dispose();

    /// Dispose query cache managers for this widget.

    clearItemPostsCache();

    clearItemsCountCache();
  }

  /// Additional helper methods.
  PagingController<DocumentSnapshot?, SubmissionRecord>
      setStaggeredViewController1(
    Query query, {
    DocumentReference<Object?>? parent,
  }) {
    staggeredViewPagingController1 ??=
        _createStaggeredViewController1(query, parent);
    if (staggeredViewPagingQuery1 != query) {
      staggeredViewPagingQuery1 = query;
      staggeredViewPagingController1?.refresh();
    }
    return staggeredViewPagingController1!;
  }

  PagingController<DocumentSnapshot?, SubmissionRecord>
      _createStaggeredViewController1(
    Query query,
    DocumentReference<Object?>? parent,
  ) {
    final controller = PagingController<DocumentSnapshot?, SubmissionRecord>(
        firstPageKey: null);
    return controller
      ..addPageRequestListener(
        (nextPageMarker) => querySubmissionRecordPage(
          queryBuilder: (_) => staggeredViewPagingQuery1 ??= query,
          nextPageMarker: nextPageMarker,
          streamSubscriptions: staggeredViewStreamSubscriptions1,
          controller: controller,
          pageSize: 25,
          isStream: true,
        ),
      );
  }
}
