import '/flutter_flow/flutter_flow_util.dart';
import '/flutter_flow/form_field_controller.dart';
import '/index.dart';
import 'method_wallet_widget.dart' show MethodWalletWidget;
import 'package:flutter/material.dart';

class MethodWalletModel extends FlutterFlowModel<MethodWalletWidget> {
  ///  State fields for stateful widgets in this page.

  final formKey = GlobalKey<FormState>();
  // State field(s) for ChoiceChips widget.
  FormFieldController<List<String>>? choiceChipsValueController1;
  String? get choiceChipsValue1 =>
      choiceChipsValueController1?.value?.firstOrNull;
  set choiceChipsValue1(String? val) =>
      choiceChipsValueController1?.value = val != null ? [val] : [];
  // State field(s) for task widget.
  FocusNode? taskFocusNode1;
  TextEditingController? taskTextController1;
  String? Function(BuildContext, String?)? taskTextController1Validator;
  String? _taskTextController1Validator(BuildContext context, String? val) {
    if (val == null || val.isEmpty) {
      return FFLocalizations.of(context).getText(
        'ojt5j8kx' /* Method Name here... is require... */,
      );
    }

    return null;
  }

  // State field(s) for ChoiceChips widget.
  FormFieldController<List<String>>? choiceChipsValueController2;
  String? get choiceChipsValue2 =>
      choiceChipsValueController2?.value?.firstOrNull;
  set choiceChipsValue2(String? val) =>
      choiceChipsValueController2?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChips widget.
  FormFieldController<List<String>>? choiceChipsValueController3;
  String? get choiceChipsValue3 =>
      choiceChipsValueController3?.value?.firstOrNull;
  set choiceChipsValue3(String? val) =>
      choiceChipsValueController3?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChips widget.
  FormFieldController<List<String>>? choiceChipsValueController4;
  String? get choiceChipsValue4 =>
      choiceChipsValueController4?.value?.firstOrNull;
  set choiceChipsValue4(String? val) =>
      choiceChipsValueController4?.value = val != null ? [val] : [];
  // State field(s) for task widget.
  FocusNode? taskFocusNode2;
  TextEditingController? taskTextController2;
  String? Function(BuildContext, String?)? taskTextController2Validator;
  String? _taskTextController2Validator(BuildContext context, String? val) {
    if (val == null || val.isEmpty) {
      return FFLocalizations.of(context).getText(
        'g1qqg58u' /* I.D. number is required */,
      );
    }

    return null;
  }

  // State field(s) for task widget.
  FocusNode? taskFocusNode3;
  TextEditingController? taskTextController3;
  String? Function(BuildContext, String?)? taskTextController3Validator;
  String? _taskTextController3Validator(BuildContext context, String? val) {
    if (val == null || val.isEmpty) {
      return FFLocalizations.of(context).getText(
        '05gv2pw2' /* Account number is required */,
      );
    }

    return null;
  }

  // State field(s) for CheckboxListTile widget.
  bool? checkboxListTileValue;
  // State field(s) for ChoiceChips widget.
  FormFieldController<List<String>>? choiceChipsValueController5;
  String? get choiceChipsValue5 =>
      choiceChipsValueController5?.value?.firstOrNull;
  set choiceChipsValue5(String? val) =>
      choiceChipsValueController5?.value = val != null ? [val] : [];

  @override
  void initState(BuildContext context) {
    taskTextController1Validator = _taskTextController1Validator;
    taskTextController2Validator = _taskTextController2Validator;
    taskTextController3Validator = _taskTextController3Validator;
  }

  @override
  void dispose() {
    taskFocusNode1?.dispose();
    taskTextController1?.dispose();

    taskFocusNode2?.dispose();
    taskTextController2?.dispose();

    taskFocusNode3?.dispose();
    taskTextController3?.dispose();
  }
}
