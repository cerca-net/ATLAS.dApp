import '/backend/backend.dart';
import '/flutter_flow/flutter_flow_timer.dart';
import '/flutter_flow/flutter_flow_util.dart';
import '/flutter_flow/form_field_controller.dart';
import '/index.dart';
import 'newpage_widget.dart' show NewpageWidget;
import 'package:stop_watch_timer/stop_watch_timer.dart';
import 'package:flutter/material.dart';
import 'package:record/record.dart';

class NewpageModel extends FlutterFlowModel<NewpageWidget> {
  ///  Local state fields for this page.

  bool isNewItem = false;

  List<String> imagesIn = [];
  void addToImagesIn(String item) => imagesIn.add(item);
  void removeFromImagesIn(String item) => imagesIn.remove(item);
  void removeAtIndexFromImagesIn(int index) => imagesIn.removeAt(index);
  void insertAtIndexInImagesIn(int index, String item) =>
      imagesIn.insert(index, item);
  void updateImagesInAtIndex(int index, Function(String) updateFn) =>
      imagesIn[index] = updateFn(imagesIn[index]);

  bool isrecord = true;

  ThreadRecord? submissionthread;

  ///  State fields for stateful widgets in this page.

  final formKey3 = GlobalKey<FormState>();
  final formKey7 = GlobalKey<FormState>();
  final formKey5 = GlobalKey<FormState>();
  final formKey1 = GlobalKey<FormState>();
  final formKey6 = GlobalKey<FormState>();
  final formKey4 = GlobalKey<FormState>();
  final formKey2 = GlobalKey<FormState>();
  // State field(s) for ChoiceIPOST widget.
  FormFieldController<List<String>>? choiceIPOSTValueController;
  String? get choiceIPOSTValue =>
      choiceIPOSTValueController?.value?.firstOrNull;
  set choiceIPOSTValue(String? val) =>
      choiceIPOSTValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceITEM widget.
  FormFieldController<List<String>>? choiceITEMValueController;
  String? get choiceITEMValue => choiceITEMValueController?.value?.firstOrNull;
  set choiceITEMValue(String? val) =>
      choiceITEMValueController?.value = val != null ? [val] : [];
  // State field(s) for PageView widget.
  PageController? pageViewController1;

  int get pageViewCurrentIndex1 => pageViewController1 != null &&
          pageViewController1!.hasClients &&
          pageViewController1!.page != null
      ? pageViewController1!.page!.round()
      : 0;
  bool isDataUploading_uploadDataBnz = false;
  List<FFUploadedFile> uploadedLocalFiles_uploadDataBnz = [];
  List<String> uploadedFileUrls_uploadDataBnz = [];

  // State field(s) for ChoiceChipsPI widget.
  FormFieldController<List<String>>? choiceChipsPIValueController;
  String? get choiceChipsPIValue =>
      choiceChipsPIValueController?.value?.firstOrNull;
  set choiceChipsPIValue(String? val) =>
      choiceChipsPIValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPIDg widget.
  FormFieldController<List<String>>? choiceChipsPIDgValueController;
  String? get choiceChipsPIDgValue =>
      choiceChipsPIDgValueController?.value?.firstOrNull;
  set choiceChipsPIDgValue(String? val) =>
      choiceChipsPIDgValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPIDt widget.
  FormFieldController<List<String>>? choiceChipsPIDtValueController;
  String? get choiceChipsPIDtValue =>
      choiceChipsPIDtValueController?.value?.firstOrNull;
  set choiceChipsPIDtValue(String? val) =>
      choiceChipsPIDtValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPIM widget.
  FormFieldController<List<String>>? choiceChipsPIMValueController;
  String? get choiceChipsPIMValue =>
      choiceChipsPIMValueController?.value?.firstOrNull;
  set choiceChipsPIMValue(String? val) =>
      choiceChipsPIMValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPIA widget.
  FormFieldController<List<String>>? choiceChipsPIAValueController;
  String? get choiceChipsPIAValue =>
      choiceChipsPIAValueController?.value?.firstOrNull;
  set choiceChipsPIAValue(String? val) =>
      choiceChipsPIAValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPIP widget.
  FormFieldController<List<String>>? choiceChipsPIPValueController;
  String? get choiceChipsPIPValue =>
      choiceChipsPIPValueController?.value?.firstOrNull;
  set choiceChipsPIPValue(String? val) =>
      choiceChipsPIPValueController?.value = val != null ? [val] : [];
  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode1;
  TextEditingController? textController1;
  String? Function(BuildContext, String?)? textController1Validator;
  String? _textController1Validator(BuildContext context, String? val) {
    if (val == null || val.isEmpty) {
      return FFLocalizations.of(context).getText(
        'm7u2wvq5' /* Field is required */,
      );
    }

    if (val.isEmpty) {
      return 'Requires at least 1 characters.';
    }
    if (val.length > 240) {
      return 'Maximum 240 characters allowed, currently ${val.length}.';
    }

    return null;
  }

  // State field(s) for ChoiceChipsThreadI widget.
  FormFieldController<List<String>>? choiceChipsThreadIValueController;
  List<String>? get choiceChipsThreadIValues =>
      choiceChipsThreadIValueController?.value;
  set choiceChipsThreadIValues(List<String>? val) =>
      choiceChipsThreadIValueController?.value = val;
  bool isDataUploading_uploadDataVideo = false;
  FFUploadedFile uploadedLocalFile_uploadDataVideo =
      FFUploadedFile(bytes: Uint8List.fromList([]));
  String uploadedFileUrl_uploadDataVideo = '';

  // State field(s) for ChoiceChipsPV widget.
  FormFieldController<List<String>>? choiceChipsPVValueController;
  String? get choiceChipsPVValue =>
      choiceChipsPVValueController?.value?.firstOrNull;
  set choiceChipsPVValue(String? val) =>
      choiceChipsPVValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPVM widget.
  FormFieldController<List<String>>? choiceChipsPVMValueController;
  String? get choiceChipsPVMValue =>
      choiceChipsPVMValueController?.value?.firstOrNull;
  set choiceChipsPVMValue(String? val) =>
      choiceChipsPVMValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPVS widget.
  FormFieldController<List<String>>? choiceChipsPVSValueController;
  String? get choiceChipsPVSValue =>
      choiceChipsPVSValueController?.value?.firstOrNull;
  set choiceChipsPVSValue(String? val) =>
      choiceChipsPVSValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPVC widget.
  FormFieldController<List<String>>? choiceChipsPVCValueController;
  String? get choiceChipsPVCValue =>
      choiceChipsPVCValueController?.value?.firstOrNull;
  set choiceChipsPVCValue(String? val) =>
      choiceChipsPVCValueController?.value = val != null ? [val] : [];
  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode2;
  TextEditingController? textController2;
  String? Function(BuildContext, String?)? textController2Validator;
  String? _textController2Validator(BuildContext context, String? val) {
    if (val == null || val.isEmpty) {
      return FFLocalizations.of(context).getText(
        '5w2eua8h' /* Field is required */,
      );
    }

    if (val.isEmpty) {
      return 'Requires at least 1 characters.';
    }
    if (val.length > 240) {
      return 'Maximum 240 characters allowed, currently ${val.length}.';
    }

    return null;
  }

  // State field(s) for ChoiceChipsThreadV widget.
  FormFieldController<List<String>>? choiceChipsThreadVValueController;
  List<String>? get choiceChipsThreadVValues =>
      choiceChipsThreadVValueController?.value;
  set choiceChipsThreadVValues(List<String>? val) =>
      choiceChipsThreadVValueController?.value = val;
  // State field(s) for Timer widget.
  final timerInitialTimeMs = 60000;
  int timerMilliseconds = 60000;
  String timerValue = StopWatchTimer.getDisplayTime(
    60000,
    hours: false,
    milliSecond: false,
  );
  FlutterFlowTimerController timerController =
      FlutterFlowTimerController(StopWatchTimer(mode: StopWatchMode.countDown));

  AudioRecorder? audioRecorder;
  String? salidaAudio;
  FFUploadedFile recordedFileBytes =
      FFUploadedFile(bytes: Uint8List.fromList([]));
  bool isDataUploading_submissionAudio = false;
  FFUploadedFile uploadedLocalFile_submissionAudio =
      FFUploadedFile(bytes: Uint8List.fromList([]));
  String uploadedFileUrl_submissionAudio = '';

  bool isDataUploading_uploadData3dt = false;
  FFUploadedFile uploadedLocalFile_uploadData3dt =
      FFUploadedFile(bytes: Uint8List.fromList([]));
  String uploadedFileUrl_uploadData3dt = '';

  // State field(s) for ChoiceChipsPA widget.
  FormFieldController<List<String>>? choiceChipsPAValueController;
  String? get choiceChipsPAValue =>
      choiceChipsPAValueController?.value?.firstOrNull;
  set choiceChipsPAValue(String? val) =>
      choiceChipsPAValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPAM widget.
  FormFieldController<List<String>>? choiceChipsPAMValueController;
  String? get choiceChipsPAMValue =>
      choiceChipsPAMValueController?.value?.firstOrNull;
  set choiceChipsPAMValue(String? val) =>
      choiceChipsPAMValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPAP widget.
  FormFieldController<List<String>>? choiceChipsPAPValueController;
  String? get choiceChipsPAPValue =>
      choiceChipsPAPValueController?.value?.firstOrNull;
  set choiceChipsPAPValue(String? val) =>
      choiceChipsPAPValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPAN widget.
  FormFieldController<List<String>>? choiceChipsPANValueController;
  String? get choiceChipsPANValue =>
      choiceChipsPANValueController?.value?.firstOrNull;
  set choiceChipsPANValue(String? val) =>
      choiceChipsPANValueController?.value = val != null ? [val] : [];
  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode3;
  TextEditingController? textController3;
  String? Function(BuildContext, String?)? textController3Validator;
  // State field(s) for ChoiceChipsThreadA widget.
  FormFieldController<List<String>>? choiceChipsThreadAValueController;
  List<String>? get choiceChipsThreadAValues =>
      choiceChipsThreadAValueController?.value;
  set choiceChipsThreadAValues(List<String>? val) =>
      choiceChipsThreadAValueController?.value = val;
  // State field(s) for TextFieldPoT widget.
  FocusNode? textFieldPoTFocusNode1;
  TextEditingController? textFieldPoTTextController1;
  String? Function(BuildContext, String?)? textFieldPoTTextController1Validator;
  // State field(s) for TextFieldPoT widget.
  FocusNode? textFieldPoTFocusNode2;
  TextEditingController? textFieldPoTTextController2;
  String? Function(BuildContext, String?)? textFieldPoTTextController2Validator;
  // State field(s) for ChoiceChipsPT widget.
  FormFieldController<List<String>>? choiceChipsPTValueController;
  String? get choiceChipsPTValue =>
      choiceChipsPTValueController?.value?.firstOrNull;
  set choiceChipsPTValue(String? val) =>
      choiceChipsPTValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPTA widget.
  FormFieldController<List<String>>? choiceChipsPTAValueController;
  String? get choiceChipsPTAValue =>
      choiceChipsPTAValueController?.value?.firstOrNull;
  set choiceChipsPTAValue(String? val) =>
      choiceChipsPTAValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPTQ widget.
  FormFieldController<List<String>>? choiceChipsPTQValueController;
  String? get choiceChipsPTQValue =>
      choiceChipsPTQValueController?.value?.firstOrNull;
  set choiceChipsPTQValue(String? val) =>
      choiceChipsPTQValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsPTS widget.
  FormFieldController<List<String>>? choiceChipsPTSValueController;
  String? get choiceChipsPTSValue =>
      choiceChipsPTSValueController?.value?.firstOrNull;
  set choiceChipsPTSValue(String? val) =>
      choiceChipsPTSValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsThreadT widget.
  FormFieldController<List<String>>? choiceChipsThreadTValueController;
  List<String>? get choiceChipsThreadTValues =>
      choiceChipsThreadTValueController?.value;
  set choiceChipsThreadTValues(List<String>? val) =>
      choiceChipsThreadTValueController?.value = val;
  // State field(s) for PageView widget.
  PageController? pageViewController2;

  int get pageViewCurrentIndex2 => pageViewController2 != null &&
          pageViewController2!.hasClients &&
          pageViewController2!.page != null
      ? pageViewController2!.page!.round()
      : 0;
  bool isDataUploading_uploadDataProduct = false;
  List<FFUploadedFile> uploadedLocalFiles_uploadDataProduct = [];
  List<String> uploadedFileUrls_uploadDataProduct = [];

  // State field(s) for ChoiceChipsItP widget.
  FormFieldController<List<String>>? choiceChipsItPValueController;
  String? get choiceChipsItPValue =>
      choiceChipsItPValueController?.value?.firstOrNull;
  set choiceChipsItPValue(String? val) =>
      choiceChipsItPValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItPT widget.
  FormFieldController<List<String>>? choiceChipsItPTValueController1;
  String? get choiceChipsItPTValue1 =>
      choiceChipsItPTValueController1?.value?.firstOrNull;
  set choiceChipsItPTValue1(String? val) =>
      choiceChipsItPTValueController1?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItPF widget.
  FormFieldController<List<String>>? choiceChipsItPFValueController;
  String? get choiceChipsItPFValue =>
      choiceChipsItPFValueController?.value?.firstOrNull;
  set choiceChipsItPFValue(String? val) =>
      choiceChipsItPFValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItPW widget.
  FormFieldController<List<String>>? choiceChipsItPWValueController;
  String? get choiceChipsItPWValue =>
      choiceChipsItPWValueController?.value?.firstOrNull;
  set choiceChipsItPWValue(String? val) =>
      choiceChipsItPWValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItPHO widget.
  FormFieldController<List<String>>? choiceChipsItPHOValueController;
  String? get choiceChipsItPHOValue =>
      choiceChipsItPHOValueController?.value?.firstOrNull;
  set choiceChipsItPHOValue(String? val) =>
      choiceChipsItPHOValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItPA widget.
  FormFieldController<List<String>>? choiceChipsItPAValueController;
  String? get choiceChipsItPAValue =>
      choiceChipsItPAValueController?.value?.firstOrNull;
  set choiceChipsItPAValue(String? val) =>
      choiceChipsItPAValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItPV widget.
  FormFieldController<List<String>>? choiceChipsItPVValueController;
  String? get choiceChipsItPVValue =>
      choiceChipsItPVValueController?.value?.firstOrNull;
  set choiceChipsItPVValue(String? val) =>
      choiceChipsItPVValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItPT widget.
  FormFieldController<List<String>>? choiceChipsItPTValueController2;
  String? get choiceChipsItPTValue2 =>
      choiceChipsItPTValueController2?.value?.firstOrNull;
  set choiceChipsItPTValue2(String? val) =>
      choiceChipsItPTValueController2?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItPO widget.
  FormFieldController<List<String>>? choiceChipsItPOValueController;
  String? get choiceChipsItPOValue =>
      choiceChipsItPOValueController?.value?.firstOrNull;
  set choiceChipsItPOValue(String? val) =>
      choiceChipsItPOValueController?.value = val != null ? [val] : [];
  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode4;
  TextEditingController? textController6;
  String? Function(BuildContext, String?)? textController6Validator;
  String? _textController6Validator(BuildContext context, String? val) {
    if (val == null || val.isEmpty) {
      return FFLocalizations.of(context).getText(
        'k5mcha95' /* Field is required */,
      );
    }

    if (val.isEmpty) {
      return 'Requires at least 1 characters.';
    }
    if (val.length > 240) {
      return 'Maximum 240 characters allowed, currently ${val.length}.';
    }

    return null;
  }

  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode5;
  TextEditingController? textController7;
  String? Function(BuildContext, String?)? textController7Validator;
  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode6;
  TextEditingController? textController8;
  String? Function(BuildContext, String?)? textController8Validator;
  // Stores action output result for [Backend Call - Create Document] action in Button widget.
  SubmissionRecord? createdObject;
  // State field(s) for ChoiceChipsThreadP widget.
  FormFieldController<List<String>>? choiceChipsThreadPValueController;
  List<String>? get choiceChipsThreadPValues =>
      choiceChipsThreadPValueController?.value;
  set choiceChipsThreadPValues(List<String>? val) =>
      choiceChipsThreadPValueController?.value = val;
  // State field(s) for PageView widget.
  PageController? pageViewController3;

  int get pageViewCurrentIndex3 => pageViewController3 != null &&
          pageViewController3!.hasClients &&
          pageViewController3!.page != null
      ? pageViewController3!.page!.round()
      : 0;
  bool isDataUploading_uploadDataService = false;
  List<FFUploadedFile> uploadedLocalFiles_uploadDataService = [];
  List<String> uploadedFileUrls_uploadDataService = [];

  // State field(s) for ChoiceChipsItS widget.
  FormFieldController<List<String>>? choiceChipsItSValueController;
  String? get choiceChipsItSValue =>
      choiceChipsItSValueController?.value?.firstOrNull;
  set choiceChipsItSValue(String? val) =>
      choiceChipsItSValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItST widget.
  FormFieldController<List<String>>? choiceChipsItSTValueController1;
  String? get choiceChipsItSTValue1 =>
      choiceChipsItSTValueController1?.value?.firstOrNull;
  set choiceChipsItSTValue1(String? val) =>
      choiceChipsItSTValueController1?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItSE widget.
  FormFieldController<List<String>>? choiceChipsItSEValueController;
  String? get choiceChipsItSEValue =>
      choiceChipsItSEValueController?.value?.firstOrNull;
  set choiceChipsItSEValue(String? val) =>
      choiceChipsItSEValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItSH widget.
  FormFieldController<List<String>>? choiceChipsItSHValueController;
  String? get choiceChipsItSHValue =>
      choiceChipsItSHValueController?.value?.firstOrNull;
  set choiceChipsItSHValue(String? val) =>
      choiceChipsItSHValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItSP widget.
  FormFieldController<List<String>>? choiceChipsItSPValueController;
  String? get choiceChipsItSPValue =>
      choiceChipsItSPValueController?.value?.firstOrNull;
  set choiceChipsItSPValue(String? val) =>
      choiceChipsItSPValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItST widget.
  FormFieldController<List<String>>? choiceChipsItSTValueController2;
  String? get choiceChipsItSTValue2 =>
      choiceChipsItSTValueController2?.value?.firstOrNull;
  set choiceChipsItSTValue2(String? val) =>
      choiceChipsItSTValueController2?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItSC widget.
  FormFieldController<List<String>>? choiceChipsItSCValueController;
  String? get choiceChipsItSCValue =>
      choiceChipsItSCValueController?.value?.firstOrNull;
  set choiceChipsItSCValue(String? val) =>
      choiceChipsItSCValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItSD widget.
  FormFieldController<List<String>>? choiceChipsItSDValueController;
  String? get choiceChipsItSDValue =>
      choiceChipsItSDValueController?.value?.firstOrNull;
  set choiceChipsItSDValue(String? val) =>
      choiceChipsItSDValueController?.value = val != null ? [val] : [];
  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode7;
  TextEditingController? textController9;
  String? Function(BuildContext, String?)? textController9Validator;
  String? _textController9Validator(BuildContext context, String? val) {
    if (val == null || val.isEmpty) {
      return FFLocalizations.of(context).getText(
        'kacwg20p' /* Field is required */,
      );
    }

    if (val.isEmpty) {
      return 'Requires at least 1 characters.';
    }
    if (val.length > 240) {
      return 'Maximum 240 characters allowed, currently ${val.length}.';
    }

    return null;
  }

  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode8;
  TextEditingController? textController10;
  String? Function(BuildContext, String?)? textController10Validator;
  String? _textController10Validator(BuildContext context, String? val) {
    if (val == null || val.isEmpty) {
      return FFLocalizations.of(context).getText(
        't9prjjt8' /* Field is required */,
      );
    }

    return null;
  }

  // State field(s) for DropDown widget.
  String? dropDownValue1;
  FormFieldController<String>? dropDownValueController1;
  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode9;
  TextEditingController? textController11;
  String? Function(BuildContext, String?)? textController11Validator;
  String? _textController11Validator(BuildContext context, String? val) {
    if (val == null || val.isEmpty) {
      return FFLocalizations.of(context).getText(
        'gz66q4ds' /* Field is required */,
      );
    }

    return null;
  }

  // State field(s) for ChoiceChipsThreadS widget.
  FormFieldController<List<String>>? choiceChipsThreadSValueController;
  List<String>? get choiceChipsThreadSValues =>
      choiceChipsThreadSValueController?.value;
  set choiceChipsThreadSValues(List<String>? val) =>
      choiceChipsThreadSValueController?.value = val;
  // State field(s) for PageView widget.
  PageController? pageViewController4;

  int get pageViewCurrentIndex4 => pageViewController4 != null &&
          pageViewController4!.hasClients &&
          pageViewController4!.page != null
      ? pageViewController4!.page!.round()
      : 0;
  bool isDataUploading_uploadDataEvent = false;
  List<FFUploadedFile> uploadedLocalFiles_uploadDataEvent = [];
  List<String> uploadedFileUrls_uploadDataEvent = [];

  // State field(s) for ChoiceChipsItE widget.
  FormFieldController<List<String>>? choiceChipsItEValueController;
  String? get choiceChipsItEValue =>
      choiceChipsItEValueController?.value?.firstOrNull;
  set choiceChipsItEValue(String? val) =>
      choiceChipsItEValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItEC widget.
  FormFieldController<List<String>>? choiceChipsItECValueController;
  String? get choiceChipsItECValue =>
      choiceChipsItECValueController?.value?.firstOrNull;
  set choiceChipsItECValue(String? val) =>
      choiceChipsItECValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItEF widget.
  FormFieldController<List<String>>? choiceChipsItEFValueController;
  String? get choiceChipsItEFValue =>
      choiceChipsItEFValueController?.value?.firstOrNull;
  set choiceChipsItEFValue(String? val) =>
      choiceChipsItEFValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItEM widget.
  FormFieldController<List<String>>? choiceChipsItEMValueController;
  String? get choiceChipsItEMValue =>
      choiceChipsItEMValueController?.value?.firstOrNull;
  set choiceChipsItEMValue(String? val) =>
      choiceChipsItEMValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItEE widget.
  FormFieldController<List<String>>? choiceChipsItEEValueController;
  String? get choiceChipsItEEValue =>
      choiceChipsItEEValueController?.value?.firstOrNull;
  set choiceChipsItEEValue(String? val) =>
      choiceChipsItEEValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItES widget.
  FormFieldController<List<String>>? choiceChipsItESValueController;
  String? get choiceChipsItESValue =>
      choiceChipsItESValueController?.value?.firstOrNull;
  set choiceChipsItESValue(String? val) =>
      choiceChipsItESValueController?.value = val != null ? [val] : [];
  // State field(s) for ChoiceChipsItEW widget.
  FormFieldController<List<String>>? choiceChipsItEWValueController;
  String? get choiceChipsItEWValue =>
      choiceChipsItEWValueController?.value?.firstOrNull;
  set choiceChipsItEWValue(String? val) =>
      choiceChipsItEWValueController?.value = val != null ? [val] : [];
  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode10;
  TextEditingController? textController12;
  String? Function(BuildContext, String?)? textController12Validator;
  String? _textController12Validator(BuildContext context, String? val) {
    if (val == null || val.isEmpty) {
      return FFLocalizations.of(context).getText(
        'crv9v00s' /* Field is required */,
      );
    }

    return null;
  }

  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode11;
  TextEditingController? textController13;
  String? Function(BuildContext, String?)? textController13Validator;
  String? _textController13Validator(BuildContext context, String? val) {
    if (val == null || val.isEmpty) {
      return FFLocalizations.of(context).getText(
        'w6yea6v3' /* Field is required */,
      );
    }

    return null;
  }

  DateTime? datePicked;
  // State field(s) for DropDown widget.
  String? dropDownValue2;
  FormFieldController<String>? dropDownValueController2;
  // State field(s) for PlacePicker widget.
  FFPlace placePickerValue = const FFPlace();
  // State field(s) for TextField widget.
  FocusNode? textFieldFocusNode12;
  TextEditingController? textController14;
  String? Function(BuildContext, String?)? textController14Validator;
  String? _textController14Validator(BuildContext context, String? val) {
    if (val == null || val.isEmpty) {
      return FFLocalizations.of(context).getText(
        'xdx6g8k1' /* Field is required */,
      );
    }

    return null;
  }

  // State field(s) for ChoiceChipsThreadE widget.
  FormFieldController<List<String>>? choiceChipsThreadEValueController;
  List<String>? get choiceChipsThreadEValues =>
      choiceChipsThreadEValueController?.value;
  set choiceChipsThreadEValues(List<String>? val) =>
      choiceChipsThreadEValueController?.value = val;

  @override
  void initState(BuildContext context) {
    textController1Validator = _textController1Validator;
    textController2Validator = _textController2Validator;
    textController6Validator = _textController6Validator;
    textController9Validator = _textController9Validator;
    textController10Validator = _textController10Validator;
    textController11Validator = _textController11Validator;
    textController12Validator = _textController12Validator;
    textController13Validator = _textController13Validator;
    textController14Validator = _textController14Validator;
  }

  @override
  void dispose() {
    textFieldFocusNode1?.dispose();
    textController1?.dispose();

    textFieldFocusNode2?.dispose();
    textController2?.dispose();

    timerController.dispose();
    textFieldFocusNode3?.dispose();
    textController3?.dispose();

    textFieldPoTFocusNode1?.dispose();
    textFieldPoTTextController1?.dispose();

    textFieldPoTFocusNode2?.dispose();
    textFieldPoTTextController2?.dispose();

    textFieldFocusNode4?.dispose();
    textController6?.dispose();

    textFieldFocusNode5?.dispose();
    textController7?.dispose();

    textFieldFocusNode6?.dispose();
    textController8?.dispose();

    textFieldFocusNode7?.dispose();
    textController9?.dispose();

    textFieldFocusNode8?.dispose();
    textController10?.dispose();

    textFieldFocusNode9?.dispose();
    textController11?.dispose();

    textFieldFocusNode10?.dispose();
    textController12?.dispose();

    textFieldFocusNode11?.dispose();
    textController13?.dispose();

    textFieldFocusNode12?.dispose();
    textController14?.dispose();
  }
}
