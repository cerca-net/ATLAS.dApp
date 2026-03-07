import '/flutter_flow/flutter_flow_util.dart';
import '/flutter_flow/form_field_controller.dart';
import '/index.dart';
import 'orderpage_widget.dart' show OrderpageWidget;
import 'package:flutter/material.dart';

class OrderpageModel extends FlutterFlowModel<OrderpageWidget> {
  ///  Local state fields for this page.

  DocumentReference? userparameter;

  DocumentReference? publicuserparameter;

  DocumentReference? qrCodevalue;

  bool qrdisplay = false;

  int rate1 = 0;

  int rate2 = 0;

  bool isPaymentProcessing = false;

  ///  State fields for stateful widgets in this page.

  // State field(s) for orderMaker widget.
  PageController? orderMakerController;

  int get orderMakerCurrentIndex => orderMakerController != null &&
          orderMakerController!.hasClients &&
          orderMakerController!.page != null
      ? orderMakerController!.page!.round()
      : 0;
  bool isDataUploading_uploadData9fw = false;
  FFUploadedFile uploadedLocalFile_uploadData9fw =
      FFUploadedFile(bytes: Uint8List.fromList([]));
  String uploadedFileUrl_uploadData9fw = '';

  var orderQR = '';
  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode;
  TextEditingController? textController;
  String? Function(BuildContext, String?)? textControllerValidator;
  // State field(s) for orderTaker widget.
  PageController? orderTakerController;

  int get orderTakerCurrentIndex => orderTakerController != null &&
          orderTakerController!.hasClients &&
          orderTakerController!.page != null
      ? orderTakerController!.page!.round()
      : 0;
  // State field(s) for ChoiceChips widget.
  FormFieldController<List<String>>? choiceChipsValueController;
  String? get choiceChipsValue =>
      choiceChipsValueController?.value?.firstOrNull;
  set choiceChipsValue(String? val) =>
      choiceChipsValueController?.value = val != null ? [val] : [];
  // State field(s) for CheckboxListTile widget.
  bool? checkboxListTileValue;
  // State field(s) for Switch widget.
  bool? switchValue1;
  // State field(s) for Switch widget.
  bool? switchValue2;
  // State field(s) for Switch widget.
  bool? switchValue3;

  @override
  void initState(BuildContext context) {}

  @override
  void dispose() {
    textFieldFocusNode?.dispose();
    textController?.dispose();
  }
}
