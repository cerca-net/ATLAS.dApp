import '/auth/auth_util.dart';
import '/backend/backend.dart';
import '/backend/schema/enums/enums.dart';
import '/components/catalogue/catalogue_widget.dart';
import '/components/object/object_widget.dart';
import '/flutter_flow/flutter_flow_animations.dart';
import '/flutter_flow/flutter_flow_audio_player.dart';
import '/flutter_flow/flutter_flow_button_tabbar.dart';
import '/flutter_flow/flutter_flow_choice_chips.dart';
import '/flutter_flow/flutter_flow_icon_button.dart';
import '/flutter_flow/flutter_flow_media_display.dart';
import '/flutter_flow/flutter_flow_theme.dart';
import '/flutter_flow/flutter_flow_util.dart';
import '/flutter_flow/flutter_flow_video_player.dart';
import '/flutter_flow/flutter_flow_widgets.dart';
import '/flutter_flow/form_field_controller.dart';
import 'dart:async';
import '/flutter_flow/custom_functions.dart' as functions;
import '/index.dart';
import 'package:badges/badges.dart' as badges;
import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_map/flutter_map.dart';
import 'package:latlong2/latlong.dart' as latlong2;
import 'package:flutter_staggered_grid_view/flutter_staggered_grid_view.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import 'package:provider/provider.dart';

import 'publicpage_model.dart';
export 'publicpage_model.dart';

class PublicpageWidget extends StatefulWidget {
  const PublicpageWidget({
    super.key,
    required this.profileparameter,
    this.objectparameter,
    this.objectrefvalue,
    this.ordermethod,
  });

  final DocumentReference? profileparameter;
  final DocumentReference? objectparameter;
  final double? objectrefvalue;
  final String? ordermethod;

  static String routeName = 'Publicpage';
  static String routePath = '/Publicpage';

  @override
  State<PublicpageWidget> createState() => _PublicpageWidgetState();
}

class _PublicpageWidgetState extends State<PublicpageWidget>
    with TickerProviderStateMixin {
  late PublicpageModel _model;

  final scaffoldKey = GlobalKey<ScaffoldState>();
  var hasIconTriggered2 = false;
  final animationsMap = <String, AnimationInfo>{};

  @override
  void initState() {
    super.initState();
    _model = createModel(context, () => PublicpageModel());

    // On page load action.
    SchedulerBinding.instance.addPostFrameCallback((_) async {
      FFAppState().isUserinPublicpage =
          !(FFAppState().isUserinPublicpage ?? true);
      safeSetState(() {});
      if (widget.objectparameter != null) {
        _model.addToBagitems(widget.objectparameter!);
      }
      if (widget.objectrefvalue != null) {
        _model.addToRefvaluesinbagitems(widget.objectrefvalue!);
      }
      if (widget.ordermethod != null) {
        _model.addToOrdermethodsinbag(widget.ordermethod!);
      }

      if (widget.profileparameter != null &&
          currentUserReference != widget.profileparameter) {
        unawaited(
          () async {
            final userRec =
                await UsersRecord.getDocumentOnce(widget.profileparameter!);
            if (userRec.analyticsRef != null) {
              await userRec.analyticsRef!.update({
                'user_impressions': FieldValue.increment(1),
              });
            }
          }(),
        );
      }

      safeSetState(() {});
    });

    _model.tabBarController = TabController(
      vsync: this,
      length: 3,
      initialIndex: 2,
    )..addListener(() => safeSetState(() {}));

    animationsMap.addAll({
      'containerOnActionTriggerAnimation': AnimationInfo(
        trigger: AnimationTrigger.onActionTrigger,
        applyInitialState: true,
        effectsBuilder: () => [
          ScaleEffect(
            curve: Curves.easeInOut,
            delay: 0.0.ms,
            duration: 600.0.ms,
            begin: const Offset(0.0, 0.0),
            end: const Offset(1.0, 1.0),
          ),
          MoveEffect(
            curve: Curves.easeInOut,
            delay: 0.0.ms,
            duration: 600.0.ms,
            begin: const Offset(100.0, 100.0),
            end: const Offset(0.0, 0.0),
          ),
        ],
      ),
      'iconOnActionTriggerAnimation1': AnimationInfo(
        trigger: AnimationTrigger.onActionTrigger,
        applyInitialState: true,
        effectsBuilder: () => [
          ScaleEffect(
            curve: Curves.easeInOut,
            delay: 0.0.ms,
            duration: 600.0.ms,
            begin: const Offset(0.0, 1.0),
            end: const Offset(0.0, 1.0),
          ),
        ],
      ),
      'iconOnActionTriggerAnimation2': AnimationInfo(
        trigger: AnimationTrigger.onActionTrigger,
        applyInitialState: true,
        effectsBuilder: () => [
          VisibilityEffect(duration: 1.ms),
          ScaleEffect(
            curve: Curves.easeInOut,
            delay: 0.0.ms,
            duration: 600.0.ms,
            begin: const Offset(0.0, 1.0),
            end: const Offset(0.0, 1.0),
          ),
        ],
      ),
    });
    setupAnimations(
      animationsMap.values.where((anim) =>
          anim.trigger == AnimationTrigger.onActionTrigger ||
          !anim.applyInitialState),
      this,
    );

    WidgetsBinding.instance.addPostFrameCallback((_) => safeSetState(() {}));
  }

  @override
  void dispose() {
    _model.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    context.watch<FFAppState>();

    return StreamBuilder<UsersRecord>(
      stream: widget.profileparameter != null
          ? UsersRecord.getDocument(widget.profileparameter!)
          : null,
      builder: (context, snapshot) {
        // Customize what your widget looks like when it's loading.
        if (!snapshot.hasData) {
          return Scaffold(
            backgroundColor: FlutterFlowTheme.of(context).secondaryBackground,
            body: Center(
              child: SizedBox(
                width: 50.0,
                height: 50.0,
                child: CircularProgressIndicator(
                  valueColor: AlwaysStoppedAnimation<Color>(
                    FlutterFlowTheme.of(context).primary,
                  ),
                ),
              ),
            ),
          );
        }

        final publicpageUsersRecord = snapshot.data!;

        return GestureDetector(
          onTap: () {
            FocusScope.of(context).unfocus();
            FocusManager.instance.primaryFocus?.unfocus();
          },
          child: Scaffold(
            key: scaffoldKey,
            backgroundColor: FlutterFlowTheme.of(context).secondaryBackground,
            body: Align(
              alignment: const AlignmentDirectional(0.0, 0.0),
              child: Padding(
                padding:
                    const EdgeInsetsDirectional.fromSTEB(8.0, 0.0, 8.0, 0.0),
                child: SafeArea(
                  child: Container(
                    width: double.infinity,
                    height: double.infinity,
                    constraints: BoxConstraints(
                      maxWidth: MediaQuery.sizeOf(context).width,
                    ),
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(12.0),
                    ),
                    child: StreamBuilder<List<AnalyticsRecord>>(
                      stream: queryAnalyticsRecord(
                        queryBuilder: (analyticsRecord) =>
                            analyticsRecord.where(
                          'user_ref',
                          isEqualTo: publicpageUsersRecord.reference,
                        ),
                        singleRecord: true,
                      ),
                      builder: (context, snapshot) {
                        // Customize what your widget looks like when it's loading.
                        if (!snapshot.hasData) {
                          return Center(
                            child: SizedBox(
                              width: 50.0,
                              height: 50.0,
                              child: CircularProgressIndicator(
                                valueColor: AlwaysStoppedAnimation<Color>(
                                  FlutterFlowTheme.of(context).primary,
                                ),
                              ),
                            ),
                          );
                        }
                        List<AnalyticsRecord> stackAnalyticsRecordList =
                            snapshot.data!;
                        final stackAnalyticsRecord =
                            stackAnalyticsRecordList.isNotEmpty
                                ? stackAnalyticsRecordList.first
                                : null;

                        return Stack(
                          children: [
                            Column(
                              mainAxisSize: MainAxisSize.min,
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Container(
                                  width: double.infinity,
                                  height: 240.0,
                                  decoration: BoxDecoration(
                                    color:
                                        FlutterFlowTheme.of(context).alternate,
                                    borderRadius: BorderRadius.circular(12.0),
                                    shape: BoxShape.rectangle,
                                  ),
                                  child: Stack(
                                    children: [
                                      Column(
                                        mainAxisSize: MainAxisSize.max,
                                        crossAxisAlignment:
                                            CrossAxisAlignment.stretch,
                                        children: [
                                          Container(
                                            width: 100.0,
                                            height: 100.0,
                                            decoration: BoxDecoration(
                                              color:
                                                  FlutterFlowTheme.of(context)
                                                      .secondaryBackground,
                                              borderRadius:
                                                  const BorderRadius.only(
                                                bottomLeft:
                                                    Radius.circular(0.0),
                                                bottomRight:
                                                    Radius.circular(0.0),
                                                topLeft: Radius.circular(12.0),
                                                topRight: Radius.circular(12.0),
                                              ),
                                            ),
                                            child: ClipRRect(
                                              borderRadius:
                                                  const BorderRadius.only(
                                                bottomLeft:
                                                    Radius.circular(0.0),
                                                bottomRight:
                                                    Radius.circular(0.0),
                                                topLeft: Radius.circular(12.0),
                                                topRight: Radius.circular(12.0),
                                              ),
                                              child: Image.network(
                                                publicpageUsersRecord.banner,
                                                width: 200.0,
                                                height: 200.0,
                                                fit: BoxFit.cover,
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                      Padding(
                                        padding: const EdgeInsetsDirectional
                                            .fromSTEB(8.0, 8.0, 8.0, 8.0),
                                        child: Row(
                                          mainAxisSize: MainAxisSize.max,
                                          mainAxisAlignment:
                                              MainAxisAlignment.spaceEvenly,
                                          crossAxisAlignment:
                                              CrossAxisAlignment.stretch,
                                          children: [
                                            Expanded(
                                              child: Padding(
                                                padding:
                                                    const EdgeInsets.all(4.0),
                                                child: Column(
                                                  mainAxisSize:
                                                      MainAxisSize.min,
                                                  mainAxisAlignment:
                                                      MainAxisAlignment.end,
                                                  crossAxisAlignment:
                                                      CrossAxisAlignment.center,
                                                  children: [
                                                    StreamBuilder<UsersRecord>(
                                                      stream: UsersRecord
                                                          .getDocument(
                                                              currentUserReference!),
                                                      builder:
                                                          (context, snapshot) {
                                                        // Customize what your widget looks like when it's loading.
                                                        if (!snapshot.hasData) {
                                                          return Center(
                                                            child: SizedBox(
                                                              width: 50.0,
                                                              height: 50.0,
                                                              child:
                                                                  CircularProgressIndicator(
                                                                valueColor:
                                                                    AlwaysStoppedAnimation<
                                                                        Color>(
                                                                  FlutterFlowTheme.of(
                                                                          context)
                                                                      .primary,
                                                                ),
                                                              ),
                                                            ),
                                                          );
                                                        }

                                                        final containerUsersRecord =
                                                            snapshot.data!;

                                                        return Container(
                                                          width: 100.0,
                                                          height: 100.0,
                                                          decoration:
                                                              BoxDecoration(
                                                            color: FlutterFlowTheme
                                                                    .of(context)
                                                                .alternate,
                                                            shape:
                                                                BoxShape.circle,
                                                            border: Border.all(
                                                              color: FlutterFlowTheme
                                                                      .of(context)
                                                                  .primaryBackground,
                                                              width: 1.0,
                                                            ),
                                                          ),
                                                          child: Padding(
                                                            padding:
                                                                const EdgeInsets
                                                                    .all(4.0),
                                                            child: Container(
                                                              width: 100.0,
                                                              height: 100.0,
                                                              clipBehavior: Clip
                                                                  .antiAlias,
                                                              decoration:
                                                                  const BoxDecoration(
                                                                shape: BoxShape
                                                                    .circle,
                                                              ),
                                                              child:
                                                                  Image.network(
                                                                publicpageUsersRecord
                                                                    .profilePicture,
                                                                fit: BoxFit
                                                                    .cover,
                                                              ),
                                                            ),
                                                          ),
                                                        );
                                                      },
                                                    ),
                                                    Align(
                                                      alignment:
                                                          const AlignmentDirectional(
                                                              0.0, 1.0),
                                                      child: RichText(
                                                        textScaler:
                                                            MediaQuery.of(
                                                                    context)
                                                                .textScaler,
                                                        text: TextSpan(
                                                          children: [
                                                            TextSpan(
                                                              text:
                                                                  publicpageUsersRecord
                                                                      .name,
                                                              style: FlutterFlowTheme
                                                                      .of(context)
                                                                  .bodyMedium
                                                                  .override(
                                                                    font: GoogleFonts
                                                                        .inter(
                                                                      fontWeight:
                                                                          FontWeight
                                                                              .bold,
                                                                      fontStyle: FlutterFlowTheme.of(
                                                                              context)
                                                                          .bodyMedium
                                                                          .fontStyle,
                                                                    ),
                                                                    letterSpacing:
                                                                        0.0,
                                                                    fontWeight:
                                                                        FontWeight
                                                                            .bold,
                                                                    fontStyle: FlutterFlowTheme.of(
                                                                            context)
                                                                        .bodyMedium
                                                                        .fontStyle,
                                                                  ),
                                                            ),
                                                            TextSpan(
                                                              text: FFLocalizations
                                                                      .of(context)
                                                                  .getText(
                                                                'd9hwpxrz' /*   */,
                                                              ),
                                                              style: FlutterFlowTheme
                                                                      .of(context)
                                                                  .bodyMedium
                                                                  .override(
                                                                    font: GoogleFonts
                                                                        .inter(
                                                                      fontWeight: FlutterFlowTheme.of(
                                                                              context)
                                                                          .bodyMedium
                                                                          .fontWeight,
                                                                      fontStyle: FlutterFlowTheme.of(
                                                                              context)
                                                                          .bodyMedium
                                                                          .fontStyle,
                                                                    ),
                                                                    letterSpacing:
                                                                        0.0,
                                                                    fontWeight: FlutterFlowTheme.of(
                                                                            context)
                                                                        .bodyMedium
                                                                        .fontWeight,
                                                                    fontStyle: FlutterFlowTheme.of(
                                                                            context)
                                                                        .bodyMedium
                                                                        .fontStyle,
                                                                  ),
                                                            ),
                                                            TextSpan(
                                                              text:
                                                                  publicpageUsersRecord
                                                                      .surname,
                                                              style:
                                                                  const TextStyle(),
                                                            )
                                                          ],
                                                          style: FlutterFlowTheme
                                                                  .of(context)
                                                              .bodyMedium
                                                              .override(
                                                                font:
                                                                    GoogleFonts
                                                                        .inter(
                                                                  fontWeight: FlutterFlowTheme.of(
                                                                          context)
                                                                      .bodyMedium
                                                                      .fontWeight,
                                                                  fontStyle: FlutterFlowTheme.of(
                                                                          context)
                                                                      .bodyMedium
                                                                      .fontStyle,
                                                                ),
                                                                letterSpacing:
                                                                    0.0,
                                                                fontWeight: FlutterFlowTheme.of(
                                                                        context)
                                                                    .bodyMedium
                                                                    .fontWeight,
                                                                fontStyle: FlutterFlowTheme.of(
                                                                        context)
                                                                    .bodyMedium
                                                                    .fontStyle,
                                                              ),
                                                        ),
                                                        textAlign:
                                                            TextAlign.start,
                                                      ),
                                                    ),
                                                    Align(
                                                      alignment:
                                                          const AlignmentDirectional(
                                                              0.0, -1.0),
                                                      child: Text(
                                                        publicpageUsersRecord
                                                            .username,
                                                        style:
                                                            FlutterFlowTheme.of(
                                                                    context)
                                                                .bodyMedium
                                                                .override(
                                                                  font:
                                                                      GoogleFonts
                                                                          .inter(
                                                                    fontWeight: FlutterFlowTheme.of(
                                                                            context)
                                                                        .bodyMedium
                                                                        .fontWeight,
                                                                    fontStyle: FlutterFlowTheme.of(
                                                                            context)
                                                                        .bodyMedium
                                                                        .fontStyle,
                                                                  ),
                                                                  letterSpacing:
                                                                      0.0,
                                                                  fontWeight: FlutterFlowTheme.of(
                                                                          context)
                                                                      .bodyMedium
                                                                      .fontWeight,
                                                                  fontStyle: FlutterFlowTheme.of(
                                                                          context)
                                                                      .bodyMedium
                                                                      .fontStyle,
                                                                ),
                                                      ),
                                                    ),
                                                    Flexible(
                                                      child: ClipRRect(
                                                        child: Container(
                                                          decoration:
                                                              const BoxDecoration(),
                                                        ),
                                                      ),
                                                    ),
                                                    if (publicpageUsersRecord
                                                            .userVerified ==
                                                        true)
                                                      Padding(
                                                        padding:
                                                            const EdgeInsetsDirectional
                                                                .fromSTEB(12.0,
                                                                0.0, 12.0, 0.0),
                                                        child: Container(
                                                          width: 120.0,
                                                          height: 40.0,
                                                          decoration:
                                                              BoxDecoration(
                                                            color: FlutterFlowTheme
                                                                    .of(context)
                                                                .secondaryBackground,
                                                            borderRadius:
                                                                BorderRadius
                                                                    .circular(
                                                                        12.0),
                                                          ),
                                                          child: Row(
                                                            mainAxisSize:
                                                                MainAxisSize
                                                                    .max,
                                                            mainAxisAlignment:
                                                                MainAxisAlignment
                                                                    .center,
                                                            children: [
                                                              Text(
                                                                FFLocalizations.of(
                                                                        context)
                                                                    .getText(
                                                                  'am409q5i' /* Verified */,
                                                                ),
                                                                style: FlutterFlowTheme.of(
                                                                        context)
                                                                    .bodyMedium
                                                                    .override(
                                                                      font: GoogleFonts
                                                                          .montserrat(
                                                                        fontWeight:
                                                                            FontWeight.w500,
                                                                        fontStyle: FlutterFlowTheme.of(context)
                                                                            .bodyMedium
                                                                            .fontStyle,
                                                                      ),
                                                                      color: FlutterFlowTheme.of(
                                                                              context)
                                                                          .secondaryText,
                                                                      letterSpacing:
                                                                          0.0,
                                                                      fontWeight:
                                                                          FontWeight
                                                                              .w500,
                                                                      fontStyle: FlutterFlowTheme.of(
                                                                              context)
                                                                          .bodyMedium
                                                                          .fontStyle,
                                                                    ),
                                                              ),
                                                              Icon(
                                                                Icons
                                                                    .verified_rounded,
                                                                color: FlutterFlowTheme.of(
                                                                        context)
                                                                    .secondary,
                                                                size: 24.0,
                                                              ),
                                                            ].divide(
                                                                const SizedBox(
                                                                    width:
                                                                        8.0)),
                                                          ),
                                                        ),
                                                      ),
                                                  ]
                                                      .divide(const SizedBox(
                                                          height: 4.0))
                                                      .addToStart(
                                                          const SizedBox(
                                                              height: 8.0)),
                                                ),
                                              ),
                                            ),
                                            Expanded(
                                              child: Column(
                                                mainAxisSize: MainAxisSize.max,
                                                mainAxisAlignment:
                                                    MainAxisAlignment.end,
                                                crossAxisAlignment:
                                                    CrossAxisAlignment.stretch,
                                                children: [
                                                  Container(
                                                    width: 100.0,
                                                    height: 120.0,
                                                    decoration:
                                                        const BoxDecoration(),
                                                    child: Column(
                                                      mainAxisSize:
                                                          MainAxisSize.max,
                                                      mainAxisAlignment:
                                                          MainAxisAlignment
                                                              .spaceEvenly,
                                                      crossAxisAlignment:
                                                          CrossAxisAlignment
                                                              .stretch,
                                                      children: [
                                                        Container(
                                                          width: 100.0,
                                                          height: 35.0,
                                                          decoration:
                                                              BoxDecoration(
                                                            color: FlutterFlowTheme
                                                                    .of(context)
                                                                .secondaryBackground,
                                                            borderRadius:
                                                                BorderRadius
                                                                    .circular(
                                                                        12.0),
                                                          ),
                                                          child: Row(
                                                            mainAxisSize:
                                                                MainAxisSize
                                                                    .max,
                                                            mainAxisAlignment:
                                                                MainAxisAlignment
                                                                    .spaceEvenly,
                                                            children: [
                                                              Expanded(
                                                                child: InkWell(
                                                                  splashColor:
                                                                      Colors
                                                                          .transparent,
                                                                  focusColor: Colors
                                                                      .transparent,
                                                                  hoverColor: Colors
                                                                      .transparent,
                                                                  highlightColor:
                                                                      Colors
                                                                          .transparent,
                                                                  onTap:
                                                                      () async {
                                                                    safeSetState(
                                                                        () {
                                                                      _model
                                                                          .tabBarController!
                                                                          .animateTo(
                                                                        2,
                                                                        duration:
                                                                            const Duration(milliseconds: 300),
                                                                        curve: Curves
                                                                            .ease,
                                                                      );
                                                                    });
                                                                  },
                                                                  child:
                                                                      Container(
                                                                    width:
                                                                        100.0,
                                                                    height:
                                                                        100.0,
                                                                    decoration:
                                                                        const BoxDecoration(),
                                                                    child: Row(
                                                                      mainAxisSize:
                                                                          MainAxisSize
                                                                              .max,
                                                                      mainAxisAlignment:
                                                                          MainAxisAlignment
                                                                              .spaceEvenly,
                                                                      children: [
                                                                        Icon(
                                                                          Icons
                                                                              .video_collection,
                                                                          color:
                                                                              FlutterFlowTheme.of(context).secondaryText,
                                                                          size:
                                                                              20.0,
                                                                        ),
                                                                        Align(
                                                                          alignment: const AlignmentDirectional(
                                                                              -1.0,
                                                                              0.0),
                                                                          child:
                                                                              FutureBuilder<int>(
                                                                            future:
                                                                                _model.postsCount(
                                                                              requestFn: () => querySubmissionRecordCount(
                                                                                queryBuilder: (submissionRecord) => submissionRecord
                                                                                    .where(
                                                                                      'poster',
                                                                                      isEqualTo: publicpageUsersRecord.reference,
                                                                                    )
                                                                                    .where(
                                                                                      'type_object',
                                                                                      isEqualTo: 'Post',
                                                                                    ),
                                                                              ),
                                                                            ),
                                                                            builder:
                                                                                (context, snapshot) {
                                                                              // Customize what your widget looks like when it's loading.
                                                                              if (!snapshot.hasData) {
                                                                                return Center(
                                                                                  child: SizedBox(
                                                                                    width: 50.0,
                                                                                    height: 50.0,
                                                                                    child: CircularProgressIndicator(
                                                                                      valueColor: AlwaysStoppedAnimation<Color>(
                                                                                        FlutterFlowTheme.of(context).primary,
                                                                                      ),
                                                                                    ),
                                                                                  ),
                                                                                );
                                                                              }
                                                                              int textCount = snapshot.data!;

                                                                              return Text(
                                                                                formatNumber(
                                                                                  textCount,
                                                                                  formatType: FormatType.compact,
                                                                                ),
                                                                                style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                      font: GoogleFonts.inter(
                                                                                        fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                      ),
                                                                                      fontSize: 16.0,
                                                                                      letterSpacing: 0.0,
                                                                                      fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                      fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                    ),
                                                                              );
                                                                            },
                                                                          ),
                                                                        ),
                                                                      ],
                                                                    ),
                                                                  ),
                                                                ),
                                                              ),
                                                              SizedBox(
                                                                height: 100.0,
                                                                child:
                                                                    VerticalDivider(
                                                                  thickness:
                                                                      1.0,
                                                                  indent: 8.0,
                                                                  endIndent:
                                                                      8.0,
                                                                  color: FlutterFlowTheme.of(
                                                                          context)
                                                                      .alternate,
                                                                ),
                                                              ),
                                                              Expanded(
                                                                child: InkWell(
                                                                  splashColor:
                                                                      Colors
                                                                          .transparent,
                                                                  focusColor: Colors
                                                                      .transparent,
                                                                  hoverColor: Colors
                                                                      .transparent,
                                                                  highlightColor:
                                                                      Colors
                                                                          .transparent,
                                                                  onTap:
                                                                      () async {
                                                                    safeSetState(
                                                                        () {
                                                                      _model
                                                                          .tabBarController!
                                                                          .animateTo(
                                                                        3,
                                                                        duration:
                                                                            const Duration(milliseconds: 300),
                                                                        curve: Curves
                                                                            .ease,
                                                                      );
                                                                    });
                                                                  },
                                                                  child:
                                                                      Container(
                                                                    width:
                                                                        100.0,
                                                                    height:
                                                                        100.0,
                                                                    decoration:
                                                                        const BoxDecoration(),
                                                                    child: Row(
                                                                      mainAxisSize:
                                                                          MainAxisSize
                                                                              .max,
                                                                      mainAxisAlignment:
                                                                          MainAxisAlignment
                                                                              .spaceEvenly,
                                                                      children: [
                                                                        Icon(
                                                                          Icons
                                                                              .shopping_bag_rounded,
                                                                          color:
                                                                              FlutterFlowTheme.of(context).secondaryText,
                                                                          size:
                                                                              20.0,
                                                                        ),
                                                                        FutureBuilder<
                                                                            int>(
                                                                          future:
                                                                              _model.itemsCount(
                                                                            requestFn: () =>
                                                                                querySubmissionRecordCount(
                                                                              queryBuilder: (submissionRecord) => submissionRecord
                                                                                  .where(
                                                                                    'poster',
                                                                                    isEqualTo: publicpageUsersRecord.reference,
                                                                                  )
                                                                                  .where(
                                                                                    'type_object',
                                                                                    isEqualTo: 'Item',
                                                                                  ),
                                                                            ),
                                                                          ),
                                                                          builder:
                                                                              (context, snapshot) {
                                                                            // Customize what your widget looks like when it's loading.
                                                                            if (!snapshot.hasData) {
                                                                              return Center(
                                                                                child: SizedBox(
                                                                                  width: 50.0,
                                                                                  height: 50.0,
                                                                                  child: CircularProgressIndicator(
                                                                                    valueColor: AlwaysStoppedAnimation<Color>(
                                                                                      FlutterFlowTheme.of(context).primary,
                                                                                    ),
                                                                                  ),
                                                                                ),
                                                                              );
                                                                            }
                                                                            int textCount =
                                                                                snapshot.data!;

                                                                            return Text(
                                                                              formatNumber(
                                                                                textCount,
                                                                                formatType: FormatType.compact,
                                                                              ),
                                                                              style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                    font: GoogleFonts.inter(
                                                                                      fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                      fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                    ),
                                                                                    fontSize: 16.0,
                                                                                    letterSpacing: 0.0,
                                                                                    fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                    fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                  ),
                                                                            );
                                                                          },
                                                                        ),
                                                                      ],
                                                                    ),
                                                                  ),
                                                                ),
                                                              ),
                                                            ],
                                                          ),
                                                        ),
                                                        Expanded(
                                                          flex: 2,
                                                          child: Container(
                                                            width: 100.0,
                                                            height: 35.0,
                                                            decoration:
                                                                BoxDecoration(
                                                              color: FlutterFlowTheme
                                                                      .of(context)
                                                                  .secondaryBackground,
                                                              borderRadius:
                                                                  BorderRadius
                                                                      .circular(
                                                                          12.0),
                                                            ),
                                                            child: Padding(
                                                              padding:
                                                                  const EdgeInsets
                                                                      .all(8.0),
                                                              child: Column(
                                                                mainAxisSize:
                                                                    MainAxisSize
                                                                        .max,
                                                                crossAxisAlignment:
                                                                    CrossAxisAlignment
                                                                        .stretch,
                                                                children: [
                                                                  Text(
                                                                    publicpageUsersRecord
                                                                        .bio,
                                                                    style: FlutterFlowTheme.of(
                                                                            context)
                                                                        .bodyMedium
                                                                        .override(
                                                                          font:
                                                                              GoogleFonts.inter(
                                                                            fontWeight:
                                                                                FontWeight.w500,
                                                                            fontStyle:
                                                                                FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                          ),
                                                                          color:
                                                                              FlutterFlowTheme.of(context).secondaryText,
                                                                          letterSpacing:
                                                                              0.0,
                                                                          fontWeight:
                                                                              FontWeight.w500,
                                                                          fontStyle: FlutterFlowTheme.of(context)
                                                                              .bodyMedium
                                                                              .fontStyle,
                                                                        ),
                                                                  ),
                                                                ],
                                                              ),
                                                            ),
                                                          ),
                                                        ),
                                                      ].divide(const SizedBox(
                                                          height: 8.0)),
                                                    ),
                                                  ),
                                                ],
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                      Align(
                                        alignment: const AlignmentDirectional(
                                            -1.0, -1.0),
                                        child: Padding(
                                          padding: const EdgeInsets.all(12.0),
                                          child: FlutterFlowIconButton(
                                            borderColor: Colors.transparent,
                                            borderRadius: 100.0,
                                            buttonSize: 35.0,
                                            fillColor:
                                                FlutterFlowTheme.of(context)
                                                    .secondaryBackground,
                                            icon: Icon(
                                              Icons.chevron_left_sharp,
                                              color:
                                                  FlutterFlowTheme.of(context)
                                                      .secondaryText,
                                              size: 18.0,
                                            ),
                                            onPressed: () async {
                                              context.safePop();
                                              FFAppState().isUserinPublicpage =
                                                  false;
                                              safeSetState(() {});
                                            },
                                          ),
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                                Flexible(
                                  child: Stack(
                                    children: [
                                      Column(
                                        children: [
                                          Align(
                                            alignment: const Alignment(-1.0, 0),
                                            child: FlutterFlowButtonTabBar(
                                              useToggleButtonStyle: true,
                                              isScrollable: true,
                                              labelStyle:
                                                  FlutterFlowTheme.of(context)
                                                      .labelSmall
                                                      .override(
                                                        font: GoogleFonts.inter(
                                                          fontWeight:
                                                              FlutterFlowTheme.of(
                                                                      context)
                                                                  .labelSmall
                                                                  .fontWeight,
                                                          fontStyle:
                                                              FlutterFlowTheme.of(
                                                                      context)
                                                                  .labelSmall
                                                                  .fontStyle,
                                                        ),
                                                        fontSize: 12.0,
                                                        letterSpacing: 0.0,
                                                        fontWeight:
                                                            FlutterFlowTheme.of(
                                                                    context)
                                                                .labelSmall
                                                                .fontWeight,
                                                        fontStyle:
                                                            FlutterFlowTheme.of(
                                                                    context)
                                                                .labelSmall
                                                                .fontStyle,
                                                      ),
                                              unselectedLabelStyle:
                                                  const TextStyle(),
                                              labelColor:
                                                  FlutterFlowTheme.of(context)
                                                      .primary,
                                              unselectedLabelColor:
                                                  FlutterFlowTheme.of(context)
                                                      .secondaryText,
                                              backgroundColor:
                                                  FlutterFlowTheme.of(context)
                                                      .accent1,
                                              unselectedBackgroundColor:
                                                  FlutterFlowTheme.of(context)
                                                      .alternate,
                                              borderWidth: 0.0,
                                              borderRadius: 12.0,
                                              elevation: 0.0,
                                              labelPadding:
                                                  const EdgeInsetsDirectional
                                                      .fromSTEB(
                                                      20.0, 0.0, 20.0, 0.0),
                                              buttonMargin:
                                                  const EdgeInsetsDirectional
                                                      .fromSTEB(
                                                      8.0, 0.0, 0.0, 0.0),
                                              tabs: [
                                                Tab(
                                                  text: FFLocalizations.of(
                                                          context)
                                                      .getText(
                                                    'ekd3suu5' /* Activity */,
                                                  ),
                                                ),
                                                Tab(
                                                  text: FFLocalizations.of(
                                                          context)
                                                      .getText(
                                                    'f3fd1gxk' /* Posts */,
                                                  ),
                                                ),
                                                Tab(
                                                  text: FFLocalizations.of(
                                                          context)
                                                      .getText(
                                                    '0pttijto' /* Items */,
                                                  ),
                                                ),
                                              ],
                                              controller:
                                                  _model.tabBarController,
                                              onTap: (i) async {
                                                [
                                                  () async {},
                                                  () async {},
                                                  () async {}
                                                ][i]();
                                              },
                                            ),
                                          ),
                                          Expanded(
                                            child: TabBarView(
                                              controller:
                                                  _model.tabBarController,
                                              children: [
                                                KeepAliveWidgetWrapper(
                                                  builder: (context) =>
                                                      Container(
                                                    width: double.infinity,
                                                    height: 400.0,
                                                    decoration: BoxDecoration(
                                                      borderRadius:
                                                          BorderRadius.circular(
                                                              0.0),
                                                    ),
                                                    child: Padding(
                                                      padding:
                                                          const EdgeInsetsDirectional
                                                              .fromSTEB(4.0,
                                                              0.0, 4.0, 0.0),
                                                      child: Column(
                                                        mainAxisSize:
                                                            MainAxisSize.min,
                                                        mainAxisAlignment:
                                                            MainAxisAlignment
                                                                .start,
                                                        crossAxisAlignment:
                                                            CrossAxisAlignment
                                                                .stretch,
                                                        children: [
                                                          Container(
                                                            width: 100.0,
                                                            decoration:
                                                                BoxDecoration(
                                                              color: FlutterFlowTheme
                                                                      .of(context)
                                                                  .alternate,
                                                              borderRadius:
                                                                  BorderRadius
                                                                      .circular(
                                                                          12.0),
                                                            ),
                                                            child: Padding(
                                                              padding:
                                                                  const EdgeInsets
                                                                      .all(4.0),
                                                              child: Column(
                                                                mainAxisSize:
                                                                    MainAxisSize
                                                                        .min,
                                                                mainAxisAlignment:
                                                                    MainAxisAlignment
                                                                        .start,
                                                                crossAxisAlignment:
                                                                    CrossAxisAlignment
                                                                        .stretch,
                                                                children: [
                                                                  Align(
                                                                    alignment:
                                                                        const AlignmentDirectional(
                                                                            -1.0,
                                                                            1.0),
                                                                    child:
                                                                        Padding(
                                                                      padding: const EdgeInsetsDirectional
                                                                          .fromSTEB(
                                                                          24.0,
                                                                          0.0,
                                                                          0.0,
                                                                          0.0),
                                                                      child:
                                                                          Text(
                                                                        FFLocalizations.of(context)
                                                                            .getText(
                                                                          'rrvc7b9p' /* Activities & Interests : */,
                                                                        ),
                                                                        style: FlutterFlowTheme.of(context)
                                                                            .bodyMedium
                                                                            .override(
                                                                              font: GoogleFonts.montserrat(
                                                                                fontWeight: FontWeight.w500,
                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                              ),
                                                                              color: FlutterFlowTheme.of(context).secondaryText,
                                                                              letterSpacing: 0.0,
                                                                              fontWeight: FontWeight.w500,
                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                            ),
                                                                      ),
                                                                    ),
                                                                  ),
                                                                  Padding(
                                                                    padding: const EdgeInsetsDirectional
                                                                        .fromSTEB(
                                                                        12.0,
                                                                        0.0,
                                                                        12.0,
                                                                        0.0),
                                                                    child:
                                                                        FlutterFlowChoiceChips(
                                                                      options: publicpageUsersRecord
                                                                          .userInterests
                                                                          .take(
                                                                              3)
                                                                          .toList()
                                                                          .map((label) =>
                                                                              ChipData(label))
                                                                          .toList(),
                                                                      onChanged:
                                                                          (val) =>
                                                                              safeSetState(() => _model.choiceChipsValue1 = val?.firstOrNull),
                                                                      selectedChipStyle:
                                                                          ChipStyle(
                                                                        backgroundColor:
                                                                            FlutterFlowTheme.of(context).primary,
                                                                        textStyle: FlutterFlowTheme.of(context)
                                                                            .bodyMedium
                                                                            .override(
                                                                              font: GoogleFonts.inter(
                                                                                fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                              ),
                                                                              color: FlutterFlowTheme.of(context).info,
                                                                              letterSpacing: 0.0,
                                                                              fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                            ),
                                                                        iconColor:
                                                                            FlutterFlowTheme.of(context).info,
                                                                        iconSize:
                                                                            16.0,
                                                                        elevation:
                                                                            0.0,
                                                                        borderRadius:
                                                                            BorderRadius.circular(8.0),
                                                                      ),
                                                                      unselectedChipStyle:
                                                                          ChipStyle(
                                                                        backgroundColor:
                                                                            FlutterFlowTheme.of(context).secondaryBackground,
                                                                        textStyle: FlutterFlowTheme.of(context)
                                                                            .bodyMedium
                                                                            .override(
                                                                              font: GoogleFonts.inter(
                                                                                fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                              ),
                                                                              color: FlutterFlowTheme.of(context).secondaryText,
                                                                              fontSize: 10.0,
                                                                              letterSpacing: 0.0,
                                                                              fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                            ),
                                                                        iconColor:
                                                                            FlutterFlowTheme.of(context).secondaryText,
                                                                        iconSize:
                                                                            16.0,
                                                                        elevation:
                                                                            0.0,
                                                                        borderRadius:
                                                                            BorderRadius.circular(8.0),
                                                                      ),
                                                                      chipSpacing:
                                                                          4.0,
                                                                      rowSpacing:
                                                                          4.0,
                                                                      multiselect:
                                                                          false,
                                                                      alignment:
                                                                          WrapAlignment
                                                                              .spaceAround,
                                                                      controller: _model
                                                                              .choiceChipsValueController1 ??=
                                                                          FormFieldController<
                                                                              List<String>>(
                                                                        [],
                                                                      ),
                                                                      wrapped:
                                                                          false,
                                                                    ),
                                                                  ),
                                                                  Padding(
                                                                    padding: const EdgeInsetsDirectional
                                                                        .fromSTEB(
                                                                        12.0,
                                                                        0.0,
                                                                        12.0,
                                                                        0.0),
                                                                    child:
                                                                        FlutterFlowChoiceChips(
                                                                      options: publicpageUsersRecord
                                                                          .userOccupations
                                                                          .take(
                                                                              3)
                                                                          .toList()
                                                                          .map((label) =>
                                                                              ChipData(label))
                                                                          .toList(),
                                                                      onChanged:
                                                                          (val) =>
                                                                              safeSetState(() => _model.choiceChipsValue2 = val?.firstOrNull),
                                                                      selectedChipStyle:
                                                                          ChipStyle(
                                                                        backgroundColor:
                                                                            FlutterFlowTheme.of(context).primary,
                                                                        textStyle: FlutterFlowTheme.of(context)
                                                                            .bodyMedium
                                                                            .override(
                                                                              font: GoogleFonts.inter(
                                                                                fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                              ),
                                                                              color: FlutterFlowTheme.of(context).info,
                                                                              letterSpacing: 0.0,
                                                                              fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                            ),
                                                                        iconColor:
                                                                            FlutterFlowTheme.of(context).info,
                                                                        iconSize:
                                                                            16.0,
                                                                        elevation:
                                                                            0.0,
                                                                        borderRadius:
                                                                            BorderRadius.circular(8.0),
                                                                      ),
                                                                      unselectedChipStyle:
                                                                          ChipStyle(
                                                                        backgroundColor:
                                                                            FlutterFlowTheme.of(context).secondaryBackground,
                                                                        textStyle: FlutterFlowTheme.of(context)
                                                                            .bodyMedium
                                                                            .override(
                                                                              font: GoogleFonts.inter(
                                                                                fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                              ),
                                                                              color: FlutterFlowTheme.of(context).secondaryText,
                                                                              fontSize: 10.0,
                                                                              letterSpacing: 0.0,
                                                                              fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                            ),
                                                                        iconColor:
                                                                            FlutterFlowTheme.of(context).secondaryText,
                                                                        iconSize:
                                                                            16.0,
                                                                        elevation:
                                                                            0.0,
                                                                        borderRadius:
                                                                            BorderRadius.circular(8.0),
                                                                      ),
                                                                      chipSpacing:
                                                                          4.0,
                                                                      rowSpacing:
                                                                          4.0,
                                                                      multiselect:
                                                                          false,
                                                                      alignment:
                                                                          WrapAlignment
                                                                              .spaceAround,
                                                                      controller: _model
                                                                              .choiceChipsValueController2 ??=
                                                                          FormFieldController<
                                                                              List<String>>(
                                                                        [],
                                                                      ),
                                                                      wrapped:
                                                                          false,
                                                                    ),
                                                                  ),
                                                                ]
                                                                    .divide(const SizedBox(
                                                                        height:
                                                                            8.0))
                                                                    .addToStart(
                                                                        const SizedBox(
                                                                            height:
                                                                                8.0))
                                                                    .addToEnd(const SizedBox(
                                                                        height:
                                                                            8.0)),
                                                              ),
                                                            ),
                                                          ),
                                                          Expanded(
                                                            child: Container(
                                                              width: 100.0,
                                                              height: 200.0,
                                                              decoration:
                                                                  BoxDecoration(
                                                                color: FlutterFlowTheme.of(
                                                                        context)
                                                                    .alternate,
                                                                borderRadius:
                                                                    BorderRadius
                                                                        .circular(
                                                                            12.0),
                                                              ),
                                                              child: Padding(
                                                                padding:
                                                                    const EdgeInsets
                                                                        .all(
                                                                        12.0),
                                                                child: Column(
                                                                  mainAxisSize:
                                                                      MainAxisSize
                                                                          .max,
                                                                  crossAxisAlignment:
                                                                      CrossAxisAlignment
                                                                          .stretch,
                                                                  children: [
                                                                    Expanded(
                                                                      child:
                                                                          Container(
                                                                        width: double
                                                                            .infinity,
                                                                        height:
                                                                            160.0,
                                                                        decoration:
                                                                            BoxDecoration(
                                                                          color:
                                                                              FlutterFlowTheme.of(context).secondaryBackground,
                                                                          borderRadius:
                                                                              BorderRadius.circular(12.0),
                                                                        ),
                                                                        child:
                                                                            Row(
                                                                          mainAxisSize:
                                                                              MainAxisSize.max,
                                                                          mainAxisAlignment:
                                                                              MainAxisAlignment.end,
                                                                          children:
                                                                              [
                                                                            Expanded(
                                                                              child: Padding(
                                                                                padding: const EdgeInsets.all(8.0),
                                                                                child: Container(
                                                                                  width: 100.0,
                                                                                  height: double.infinity,
                                                                                  decoration: BoxDecoration(
                                                                                    color: FlutterFlowTheme.of(context).secondaryBackground,
                                                                                    borderRadius: BorderRadius.circular(12.0),
                                                                                  ),
                                                                                  child: Column(
                                                                                    mainAxisSize: MainAxisSize.max,
                                                                                    crossAxisAlignment: CrossAxisAlignment.stretch,
                                                                                    children: [
                                                                                      Padding(
                                                                                        padding: const EdgeInsets.all(8.0),
                                                                                        child: Text(
                                                                                          FFLocalizations.of(context).getText(
                                                                                            'ieiybayx' /* Places : */,
                                                                                          ),
                                                                                          style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                font: GoogleFonts.montserrat(
                                                                                                  fontWeight: FontWeight.w600,
                                                                                                  fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                ),
                                                                                                color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                letterSpacing: 0.0,
                                                                                                fontWeight: FontWeight.w600,
                                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                              ),
                                                                                        ),
                                                                                      ),
                                                                                      Padding(
                                                                                        padding: const EdgeInsetsDirectional.fromSTEB(12.0, 4.0, 12.0, 0.0),
                                                                                        child: FlutterFlowChoiceChips(
                                                                                          options: publicpageUsersRecord.userPlaces.take(3).toList().map((label) => ChipData(label)).toList(),
                                                                                          onChanged: (val) => safeSetState(() => _model.choiceChipsValue3 = val?.firstOrNull),
                                                                                          selectedChipStyle: ChipStyle(
                                                                                            backgroundColor: FlutterFlowTheme.of(context).secondary,
                                                                                            textStyle: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                  font: GoogleFonts.inter(
                                                                                                    fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                                    fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                  ),
                                                                                                  color: FlutterFlowTheme.of(context).info,
                                                                                                  fontSize: 10.0,
                                                                                                  letterSpacing: 0.0,
                                                                                                  fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                                  fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                ),
                                                                                            iconColor: FlutterFlowTheme.of(context).info,
                                                                                            iconSize: 16.0,
                                                                                            elevation: 0.0,
                                                                                            borderRadius: BorderRadius.circular(8.0),
                                                                                          ),
                                                                                          unselectedChipStyle: ChipStyle(
                                                                                            backgroundColor: FlutterFlowTheme.of(context).alternate,
                                                                                            textStyle: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                  font: GoogleFonts.inter(
                                                                                                    fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                                    fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                  ),
                                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                  fontSize: 10.0,
                                                                                                  letterSpacing: 0.0,
                                                                                                  fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                                  fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                ),
                                                                                            iconColor: FlutterFlowTheme.of(context).secondaryText,
                                                                                            iconSize: 16.0,
                                                                                            elevation: 0.0,
                                                                                            borderRadius: BorderRadius.circular(8.0),
                                                                                          ),
                                                                                          chipSpacing: 4.0,
                                                                                          rowSpacing: 4.0,
                                                                                          multiselect: false,
                                                                                          alignment: WrapAlignment.center,
                                                                                          controller: _model.choiceChipsValueController3 ??= FormFieldController<List<String>>(
                                                                                            [],
                                                                                          ),
                                                                                          wrapped: true,
                                                                                        ),
                                                                                      ),
                                                                                    ],
                                                                                  ),
                                                                                ),
                                                                              ),
                                                                            ),
                                                                            SizedBox(
                                                                              height: 100.0,
                                                                              child: VerticalDivider(
                                                                                thickness: 2.0,
                                                                                color: FlutterFlowTheme.of(context).alternate,
                                                                              ),
                                                                            ),
                                                                            Expanded(
                                                                              child: Padding(
                                                                                padding: const EdgeInsets.all(8.0),
                                                                                child: Column(
                                                                                  mainAxisSize: MainAxisSize.max,
                                                                                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                                                                                  crossAxisAlignment: CrossAxisAlignment.end,
                                                                                  children: [
                                                                                    Expanded(
                                                                                      child: Align(
                                                                                        alignment: const AlignmentDirectional(1.0, 0.0),
                                                                                        child: RichText(
                                                                                          textScaler: MediaQuery.of(context).textScaler,
                                                                                          text: TextSpan(
                                                                                            children: [
                                                                                              TextSpan(
                                                                                                text: FFLocalizations.of(context).getText(
                                                                                                  'p1qie27l' /* USER PIN :  */,
                                                                                                ),
                                                                                                style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                      font: GoogleFonts.montserrat(
                                                                                                        fontWeight: FontWeight.w500,
                                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                      ),
                                                                                                      color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                      fontSize: 12.0,
                                                                                                      letterSpacing: 0.0,
                                                                                                      fontWeight: FontWeight.w500,
                                                                                                      fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                    ),
                                                                                              ),
                                                                                              TextSpan(
                                                                                                text: valueOrDefault<String>(
                                                                                                  formatNumber(
                                                                                                    publicpageUsersRecord.userPins.length,
                                                                                                    formatType: FormatType.compact,
                                                                                                  ),
                                                                                                  '0',
                                                                                                ),
                                                                                                style: TextStyle(
                                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                  fontSize: 16.0,
                                                                                                ),
                                                                                              )
                                                                                            ],
                                                                                            style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                  font: GoogleFonts.montserrat(
                                                                                                    fontWeight: FontWeight.w500,
                                                                                                    fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                  ),
                                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                  fontSize: 12.0,
                                                                                                  letterSpacing: 0.0,
                                                                                                  fontWeight: FontWeight.w500,
                                                                                                  fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                ),
                                                                                          ),
                                                                                          textAlign: TextAlign.center,
                                                                                        ),
                                                                                      ),
                                                                                    ),
                                                                                    Expanded(
                                                                                      child: Align(
                                                                                        alignment: const AlignmentDirectional(1.0, 0.0),
                                                                                        child: RichText(
                                                                                          textScaler: MediaQuery.of(context).textScaler,
                                                                                          text: TextSpan(
                                                                                            children: [
                                                                                              TextSpan(
                                                                                                text: FFLocalizations.of(context).getText(
                                                                                                  'ewrxuzmi' /* OBJECT PIN :  */,
                                                                                                ),
                                                                                                style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                      font: GoogleFonts.montserrat(
                                                                                                        fontWeight: FontWeight.w500,
                                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                      ),
                                                                                                      color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                      fontSize: 12.0,
                                                                                                      letterSpacing: 0.0,
                                                                                                      fontWeight: FontWeight.w500,
                                                                                                      fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                    ),
                                                                                              ),
                                                                                              TextSpan(
                                                                                                text: FFLocalizations.of(context).getText(
                                                                                                  'fbyjwip0' /* 212 */,
                                                                                                ),
                                                                                                style: TextStyle(
                                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                  fontSize: 16.0,
                                                                                                ),
                                                                                              )
                                                                                            ],
                                                                                            style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                  font: GoogleFonts.montserrat(
                                                                                                    fontWeight: FontWeight.w500,
                                                                                                    fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                  ),
                                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                  fontSize: 12.0,
                                                                                                  letterSpacing: 0.0,
                                                                                                  fontWeight: FontWeight.w500,
                                                                                                  fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                ),
                                                                                          ),
                                                                                          textAlign: TextAlign.center,
                                                                                        ),
                                                                                      ),
                                                                                    ),
                                                                                    Expanded(
                                                                                      child: Align(
                                                                                        alignment: const AlignmentDirectional(1.0, 0.0),
                                                                                        child: RichText(
                                                                                          textScaler: MediaQuery.of(context).textScaler,
                                                                                          text: TextSpan(
                                                                                            children: [
                                                                                              TextSpan(
                                                                                                text: FFLocalizations.of(context).getText(
                                                                                                  'j6onlmua' /* VOTE RATE :  */,
                                                                                                ),
                                                                                                style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                      font: GoogleFonts.montserrat(
                                                                                                        fontWeight: FontWeight.w500,
                                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                      ),
                                                                                                      color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                      fontSize: 12.0,
                                                                                                      letterSpacing: 0.0,
                                                                                                      fontWeight: FontWeight.w500,
                                                                                                      fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                    ),
                                                                                              ),
                                                                                              TextSpan(
                                                                                                text: FFLocalizations.of(context).getText(
                                                                                                  'avby2qrn' /* 213 */,
                                                                                                ),
                                                                                                style: TextStyle(
                                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                  fontSize: 16.0,
                                                                                                ),
                                                                                              )
                                                                                            ],
                                                                                            style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                  font: GoogleFonts.montserrat(
                                                                                                    fontWeight: FontWeight.w500,
                                                                                                    fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                  ),
                                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                  fontSize: 12.0,
                                                                                                  letterSpacing: 0.0,
                                                                                                  fontWeight: FontWeight.w500,
                                                                                                  fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                ),
                                                                                          ),
                                                                                          textAlign: TextAlign.center,
                                                                                        ),
                                                                                      ),
                                                                                    ),
                                                                                    Expanded(
                                                                                      child: Align(
                                                                                        alignment: const AlignmentDirectional(1.0, 0.0),
                                                                                        child: RichText(
                                                                                          textScaler: MediaQuery.of(context).textScaler,
                                                                                          text: TextSpan(
                                                                                            children: [
                                                                                              TextSpan(
                                                                                                text: FFLocalizations.of(context).getText(
                                                                                                  '6u4ml2fe' /* ORDER RATE :  */,
                                                                                                ),
                                                                                                style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                      font: GoogleFonts.montserrat(
                                                                                                        fontWeight: FontWeight.w500,
                                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                      ),
                                                                                                      color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                      fontSize: 12.0,
                                                                                                      letterSpacing: 0.0,
                                                                                                      fontWeight: FontWeight.w500,
                                                                                                      fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                    ),
                                                                                              ),
                                                                                              TextSpan(
                                                                                                text: FFLocalizations.of(context).getText(
                                                                                                  'gz1x8wb4' /* 213 */,
                                                                                                ),
                                                                                                style: TextStyle(
                                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                  fontSize: 16.0,
                                                                                                ),
                                                                                              )
                                                                                            ],
                                                                                            style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                  font: GoogleFonts.montserrat(
                                                                                                    fontWeight: FontWeight.w500,
                                                                                                    fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                  ),
                                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                  fontSize: 12.0,
                                                                                                  letterSpacing: 0.0,
                                                                                                  fontWeight: FontWeight.w500,
                                                                                                  fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                ),
                                                                                          ),
                                                                                          textAlign: TextAlign.center,
                                                                                        ),
                                                                                      ),
                                                                                    ),
                                                                                    Expanded(
                                                                                      child: RichText(
                                                                                        textScaler: MediaQuery.of(context).textScaler,
                                                                                        text: TextSpan(
                                                                                          children: [
                                                                                            TextSpan(
                                                                                              text: FFLocalizations.of(context).getText(
                                                                                                '8jz7qus4' /*  # ORDERS :  */,
                                                                                              ),
                                                                                              style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                    font: GoogleFonts.montserrat(
                                                                                                      fontWeight: FontWeight.w500,
                                                                                                      fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                    ),
                                                                                                    color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                    fontSize: 12.0,
                                                                                                    letterSpacing: 0.0,
                                                                                                    fontWeight: FontWeight.w500,
                                                                                                    fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                  ),
                                                                                            ),
                                                                                            TextSpan(
                                                                                              text: FFLocalizations.of(context).getText(
                                                                                                'yl1fcbft' /* 212 */,
                                                                                              ),
                                                                                              style: TextStyle(
                                                                                                color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                fontSize: 16.0,
                                                                                              ),
                                                                                            )
                                                                                          ],
                                                                                          style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                font: GoogleFonts.montserrat(
                                                                                                  fontWeight: FontWeight.w500,
                                                                                                  fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                ),
                                                                                                color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                fontSize: 12.0,
                                                                                                letterSpacing: 0.0,
                                                                                                fontWeight: FontWeight.w500,
                                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                              ),
                                                                                        ),
                                                                                        textAlign: TextAlign.center,
                                                                                      ),
                                                                                    ),
                                                                                  ].divide(const SizedBox(height: 8.0)).addToStart(const SizedBox(height: 4.0)).addToEnd(const SizedBox(height: 4.0)),
                                                                                ),
                                                                              ),
                                                                            ),
                                                                          ].addToStart(const SizedBox(width: 8.0)).addToEnd(const SizedBox(width: 8.0)),
                                                                        ),
                                                                      ),
                                                                    ),
                                                                    Expanded(
                                                                      child:
                                                                          Container(
                                                                        width: double
                                                                            .infinity,
                                                                        height:
                                                                            160.0,
                                                                        decoration:
                                                                            BoxDecoration(
                                                                          color:
                                                                              FlutterFlowTheme.of(context).secondaryBackground,
                                                                          borderRadius:
                                                                              BorderRadius.circular(12.0),
                                                                        ),
                                                                        child:
                                                                            Padding(
                                                                          padding: const EdgeInsets
                                                                              .all(
                                                                              8.0),
                                                                          child:
                                                                              Column(
                                                                            mainAxisSize:
                                                                                MainAxisSize.max,
                                                                            crossAxisAlignment:
                                                                                CrossAxisAlignment.stretch,
                                                                            children:
                                                                                [
                                                                              Expanded(
                                                                                child: ClipRRect(
                                                                                  borderRadius: BorderRadius.circular(12.0),
                                                                                  child: FlutterMap(
                                                                                    options: const MapOptions(
                                                                                      initialCenter: latlong2.LatLng(0.0, 0.0), // Default generic location
                                                                                      initialZoom: 1.0,
                                                                                      interactionOptions: InteractionOptions(
                                                                                        flags: InteractiveFlag.all & ~InteractiveFlag.rotate,
                                                                                      ),
                                                                                    ),
                                                                                    children: [
                                                                                      TileLayer(
                                                                                        urlTemplate: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
                                                                                        userAgentPackageName: 'com.mycompany.cercaend',
                                                                                      ),
                                                                                    ],
                                                                                  ),
                                                                                ),
                                                                              ),
                                                                            ].divide(const SizedBox(height: 8.0)),
                                                                          ),
                                                                        ),
                                                                      ),
                                                                    ),
                                                                  ].divide(const SizedBox(
                                                                      height:
                                                                          8.0)),
                                                                ),
                                                              ),
                                                            ),
                                                          ),
                                                        ]
                                                            .divide(
                                                                const SizedBox(
                                                                    height:
                                                                        8.0))
                                                            .addToStart(
                                                                const SizedBox(
                                                                    height:
                                                                        8.0))
                                                            .addToEnd(
                                                                const SizedBox(
                                                                    height:
                                                                        8.0)),
                                                      ),
                                                    ),
                                                  ),
                                                ),
                                                KeepAliveWidgetWrapper(
                                                  builder: (context) =>
                                                      Container(
                                                    width: double.infinity,
                                                    height: double.infinity,
                                                    decoration: BoxDecoration(
                                                      borderRadius:
                                                          BorderRadius.circular(
                                                              8.0),
                                                    ),
                                                    child: Padding(
                                                      padding:
                                                          const EdgeInsetsDirectional
                                                              .fromSTEB(4.0,
                                                              0.0, 4.0, 0.0),
                                                      child: Column(
                                                        mainAxisSize:
                                                            MainAxisSize.max,
                                                        children: [
                                                          Padding(
                                                            padding:
                                                                const EdgeInsetsDirectional
                                                                    .fromSTEB(
                                                                    0.0,
                                                                    8.0,
                                                                    0.0,
                                                                    8.0),
                                                            child: Container(
                                                              width: double
                                                                  .infinity,
                                                              height: 50.0,
                                                              decoration:
                                                                  BoxDecoration(
                                                                color: FlutterFlowTheme.of(
                                                                        context)
                                                                    .alternate,
                                                                borderRadius:
                                                                    BorderRadius
                                                                        .circular(
                                                                            8.0),
                                                              ),
                                                              child: Row(
                                                                mainAxisSize:
                                                                    MainAxisSize
                                                                        .max,
                                                                mainAxisAlignment:
                                                                    MainAxisAlignment
                                                                        .spaceEvenly,
                                                                children: [
                                                                  Align(
                                                                    alignment:
                                                                        const AlignmentDirectional(
                                                                            0.0,
                                                                            0.0),
                                                                    child:
                                                                        FlutterFlowChoiceChips(
                                                                      options: [
                                                                        ChipData(
                                                                            FFLocalizations.of(context).getText(
                                                                              'wk8go3z8' /* Image */,
                                                                            ),
                                                                            Icons.image),
                                                                        ChipData(
                                                                            FFLocalizations.of(context).getText(
                                                                              '3xfhvvjj' /* Video */,
                                                                            ),
                                                                            Icons.play_circle),
                                                                        ChipData(
                                                                            FFLocalizations.of(context).getText(
                                                                              'hy0alzf3' /* Audio */,
                                                                            ),
                                                                            Icons.audiotrack_rounded),
                                                                        ChipData(
                                                                            FFLocalizations.of(context).getText(
                                                                              '3h6upgka' /* Text */,
                                                                            ),
                                                                            Icons.text_fields)
                                                                      ],
                                                                      onChanged:
                                                                          (val) =>
                                                                              safeSetState(() => _model.choiceChipsPostsValues = val),
                                                                      selectedChipStyle:
                                                                          ChipStyle(
                                                                        backgroundColor:
                                                                            FlutterFlowTheme.of(context).secondary,
                                                                        textStyle: FlutterFlowTheme.of(context)
                                                                            .bodyMedium
                                                                            .override(
                                                                              font: GoogleFonts.inter(
                                                                                fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                              ),
                                                                              color: FlutterFlowTheme.of(context).info,
                                                                              fontSize: 12.0,
                                                                              letterSpacing: 0.0,
                                                                              fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                            ),
                                                                        iconColor:
                                                                            FlutterFlowTheme.of(context).info,
                                                                        iconSize:
                                                                            16.0,
                                                                        elevation:
                                                                            0.0,
                                                                        borderRadius:
                                                                            BorderRadius.circular(8.0),
                                                                      ),
                                                                      unselectedChipStyle:
                                                                          ChipStyle(
                                                                        backgroundColor:
                                                                            FlutterFlowTheme.of(context).secondaryBackground,
                                                                        textStyle: FlutterFlowTheme.of(context)
                                                                            .bodyMedium
                                                                            .override(
                                                                              font: GoogleFonts.inter(
                                                                                fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                              ),
                                                                              color: FlutterFlowTheme.of(context).secondaryText,
                                                                              letterSpacing: 0.0,
                                                                              fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                            ),
                                                                        iconColor:
                                                                            FlutterFlowTheme.of(context).secondaryText,
                                                                        iconSize:
                                                                            16.0,
                                                                        elevation:
                                                                            0.0,
                                                                        borderRadius:
                                                                            BorderRadius.circular(8.0),
                                                                      ),
                                                                      chipSpacing:
                                                                          6.0,
                                                                      rowSpacing:
                                                                          8.0,
                                                                      multiselect:
                                                                          true,
                                                                      initialized:
                                                                          _model.choiceChipsPostsValues !=
                                                                              null,
                                                                      alignment:
                                                                          WrapAlignment
                                                                              .start,
                                                                      controller: _model
                                                                              .choiceChipsPostsValueController ??=
                                                                          FormFieldController<
                                                                              List<String>>(
                                                                        [
                                                                          FFLocalizations.of(context)
                                                                              .getText(
                                                                            'i32agl68' /* Image */,
                                                                          )
                                                                        ],
                                                                      ),
                                                                      wrapped:
                                                                          false,
                                                                    ),
                                                                  ),
                                                                ],
                                                              ),
                                                            ),
                                                          ),
                                                          SizedBox(
                                                            child: Container(
                                                              width: double
                                                                  .infinity,
                                                              height: 100.0,
                                                              decoration:
                                                                  const BoxDecoration(),
                                                              child: PagedMasonryGridView<
                                                                  DocumentSnapshot<
                                                                      Object?>?,
                                                                  SubmissionRecord>.count(
                                                                pagingController:
                                                                    _model
                                                                        .setStaggeredViewController1(
                                                                  SubmissionRecord
                                                                      .collection
                                                                      .where(
                                                                        'poster',
                                                                        isEqualTo:
                                                                            publicpageUsersRecord.reference,
                                                                      )
                                                                      .where(
                                                                        'type_object',
                                                                        isEqualTo:
                                                                            'Post',
                                                                      )
                                                                      .orderBy(
                                                                          'date',
                                                                          descending:
                                                                              true),
                                                                ),
                                                                shrinkWrap:
                                                                    true,
                                                                physics:
                                                                    const NeverScrollableScrollPhysics(),
                                                                crossAxisCount: () {
                                                                  if (MediaQuery.sizeOf(context).width < 430.0) {
                                                                    return 2;
                                                                  } else if (MediaQuery.sizeOf(context).width <= 1024.0) {
                                                                    return 3;
                                                                  } else {
                                                                    return 5;
                                                                  }
                                                                }(),
                                                                crossAxisSpacing:
                                                                    10.0,
                                                                mainAxisSpacing:
                                                                    10.0,
                                                                builderDelegate:
                                                                    PagedChildBuilderDelegate<
                                                                        SubmissionRecord>(
                                                                  // Customize what your widget looks like when it's loading the first page.
                                                                  firstPageProgressIndicatorBuilder:
                                                                      (_) =>
                                                                          Center(
                                                                    child:
                                                                        SizedBox(
                                                                      width:
                                                                          50.0,
                                                                      height:
                                                                          50.0,
                                                                      child:
                                                                          CircularProgressIndicator(
                                                                        valueColor:
                                                                            AlwaysStoppedAnimation<Color>(
                                                                          FlutterFlowTheme.of(context)
                                                                              .primary,
                                                                        ),
                                                                      ),
                                                                    ),
                                                                  ),
                                                                  // Customize what your widget looks like when it's loading another page.
                                                                  newPageProgressIndicatorBuilder:
                                                                      (_) =>
                                                                          Center(
                                                                    child:
                                                                        SizedBox(
                                                                      width:
                                                                          50.0,
                                                                      height:
                                                                          50.0,
                                                                      child:
                                                                          CircularProgressIndicator(
                                                                        valueColor:
                                                                            AlwaysStoppedAnimation<Color>(
                                                                          FlutterFlowTheme.of(context)
                                                                              .primary,
                                                                        ),
                                                                      ),
                                                                    ),
                                                                  ),

                                                                  itemBuilder:
                                                                      (context,
                                                                          _,
                                                                          staggeredViewIndex) {
                                                                    final staggeredViewSubmissionRecord = _model
                                                                        .staggeredViewPagingController1!
                                                                        .itemList![staggeredViewIndex];
                                                                    return Builder(
                                                                      builder: (context) =>
                                                                          FutureBuilder<
                                                                              SubmissionRecord>(
                                                                        future:
                                                                            SubmissionRecord.getDocumentOnce(staggeredViewSubmissionRecord.reference),
                                                                        builder:
                                                                            (context,
                                                                                snapshot) {
                                                                          // Customize what your widget looks like when it's loading.
                                                                          if (!snapshot
                                                                              .hasData) {
                                                                            return Center(
                                                                              child: SizedBox(
                                                                                width: 50.0,
                                                                                height: 50.0,
                                                                                child: CircularProgressIndicator(
                                                                                  valueColor: AlwaysStoppedAnimation<Color>(
                                                                                    FlutterFlowTheme.of(context).primary,
                                                                                  ),
                                                                                ),
                                                                              ),
                                                                            );
                                                                          }

                                                                          final objectSubmissionRecord =
                                                                              snapshot.data!;

                                                                          return InkWell(
                                                                            splashColor:
                                                                                Colors.transparent,
                                                                            focusColor:
                                                                                Colors.transparent,
                                                                            hoverColor:
                                                                                Colors.transparent,
                                                                            highlightColor:
                                                                                Colors.transparent,
                                                                            onTap:
                                                                                () async {
                                                                              await showDialog(
                                                                                context: context,
                                                                                builder: (dialogContext) {
                                                                                  return Dialog(
                                                                                    elevation: 0,
                                                                                    insetPadding: EdgeInsets.zero,
                                                                                    backgroundColor: Colors.transparent,
                                                                                    alignment: const AlignmentDirectional(0.0, 0.0).resolve(Directionality.of(context)),
                                                                                    child: GestureDetector(
                                                                                      onTap: () {
                                                                                        FocusScope.of(dialogContext).unfocus();
                                                                                        FocusManager.instance.primaryFocus?.unfocus();
                                                                                      },
                                                                                      child: ObjectWidget(
                                                                                        object: objectSubmissionRecord,
                                                                                      ),
                                                                                    ),
                                                                                  );
                                                                                },
                                                                              );
                                                                            },
                                                                            child:
                                                                                ClipRRect(
                                                                              borderRadius: BorderRadius.circular(8.0),
                                                                              child: Container(
                                                                                width: double.infinity,
                                                                                constraints: const BoxConstraints(
                                                                                  maxHeight: 320.0,
                                                                                ),
                                                                                decoration: BoxDecoration(
                                                                                  color: FlutterFlowTheme.of(context).secondaryBackground,
                                                                                  boxShadow: const [
                                                                                    BoxShadow(
                                                                                      blurRadius: 5.0,
                                                                                      color: Color(0x3416202A),
                                                                                      offset: Offset(
                                                                                        0.0,
                                                                                        3.0,
                                                                                      ),
                                                                                    )
                                                                                  ],
                                                                                  borderRadius: BorderRadius.circular(8.0),
                                                                                ),
                                                                                child: Column(
                                                                                  mainAxisSize: MainAxisSize.min,
                                                                                  mainAxisAlignment: MainAxisAlignment.start,
                                                                                  crossAxisAlignment: CrossAxisAlignment.stretch,
                                                                                  children: [
                                                                                    if ((objectSubmissionRecord.imagesextra.where((e) => e.isNotEmpty).isNotEmpty) && (objectSubmissionRecord.refvalue == ''))
                                                                                      Builder(
                                                                                        builder: (context) {
                                                                                          final feedobjectimages = objectSubmissionRecord.imagesextra.where((e) => e.isNotEmpty).take(3).toList();

                                                                                          return SingleChildScrollView(
                                                                                            scrollDirection: Axis.horizontal,
                                                                                            child: Row(
                                                                                              mainAxisSize: MainAxisSize.min,
                                                                                              crossAxisAlignment: CrossAxisAlignment.start,
                                                                                              children: List.generate(feedobjectimages.length, (feedobjectimagesIndex) {
                                                                                                final feedobjectimagesItem = feedobjectimages[feedobjectimagesIndex];
                                                                                                return Container(
                                                                                                  width: 200.0,
                                                                                                  height: 200.0,
                                                                                                  decoration: const BoxDecoration(),
                                                                                                  child: Visibility(
                                                                                                    visible: feedobjectimagesItem != '',
                                                                                                    child: ClipRRect(
                                                                                                      borderRadius: BorderRadius.circular(8.0),
                                                                                                      child: Image.network(
                                                                                                        feedobjectimagesItem,
                                                                                                        width: double.infinity,
                                                                                                        height: double.infinity,
                                                                                                        fit: BoxFit.fitHeight,
                                                                                                      ),
                                                                                                    ),
                                                                                                  ),
                                                                                                );
                                                                                              }),
                                                                                            ),
                                                                                          );
                                                                                        },
                                                                                      ),
                                                                                    if (objectSubmissionRecord.video != '')
                                                                                      Expanded(
                                                                                        child: Row(
                                                                                          mainAxisSize: MainAxisSize.max,
                                                                                          crossAxisAlignment: CrossAxisAlignment.start,
                                                                                          children: [
                                                                                            if (objectSubmissionRecord.video != '')
                                                                                              Expanded(
                                                                                                child: FlutterFlowMediaDisplay(
                                                                                                  path: objectSubmissionRecord.video,
                                                                                                  imageBuilder: (path) => ClipRRect(
                                                                                                    borderRadius: BorderRadius.circular(8.0),
                                                                                                    child: Image.network(
                                                                                                      path,
                                                                                                      width: double.infinity,
                                                                                                      height: double.infinity,
                                                                                                      fit: BoxFit.fitHeight,
                                                                                                    ),
                                                                                                  ),
                                                                                                  videoPlayerBuilder: (path) => FlutterFlowVideoPlayer(
                                                                                                    path: path,
                                                                                                    width: double.infinity,
                                                                                                    height: double.infinity,
                                                                                                    autoPlay: true,
                                                                                                    looping: true,
                                                                                                    showControls: false,
                                                                                                    allowFullScreen: false,
                                                                                                    allowPlaybackSpeedMenu: false,
                                                                                                  ),
                                                                                                ),
                                                                                              ),
                                                                                          ],
                                                                                        ),
                                                                                      ),
                                                                                    if (objectSubmissionRecord.audio != '')
                                                                                      Expanded(
                                                                                        child: Row(
                                                                                          mainAxisSize: MainAxisSize.min,
                                                                                          children: [
                                                                                            Expanded(
                                                                                              child: Container(
                                                                                                width: 100.0,
                                                                                                height: 320.0,
                                                                                                decoration: BoxDecoration(
                                                                                                  color: FlutterFlowTheme.of(context).secondaryBackground,
                                                                                                ),
                                                                                                child: Stack(
                                                                                                  children: [
                                                                                                    ClipRRect(
                                                                                                      borderRadius: BorderRadius.circular(8.0),
                                                                                                      child: Image.network(
                                                                                                        objectSubmissionRecord.image,
                                                                                                        width: 200.0,
                                                                                                        height: 200.0,
                                                                                                        fit: BoxFit.cover,
                                                                                                      ),
                                                                                                    ),
                                                                                                    FlutterFlowAudioPlayer(
                                                                                                      audio: Audio.network(
                                                                                                        objectSubmissionRecord.audio,
                                                                                                        metas: Metas(
                                                                                                          title: objectSubmissionRecord.header,
                                                                                                        ),
                                                                                                      ),
                                                                                                      titleTextStyle: FlutterFlowTheme.of(context).titleLarge.override(
                                                                                                            font: GoogleFonts.sora(
                                                                                                              fontWeight: FlutterFlowTheme.of(context).titleLarge.fontWeight,
                                                                                                              fontStyle: FlutterFlowTheme.of(context).titleLarge.fontStyle,
                                                                                                            ),
                                                                                                            color: FlutterFlowTheme.of(context).secondaryBackground,
                                                                                                            fontSize: 20.0,
                                                                                                            letterSpacing: 0.0,
                                                                                                            fontWeight: FlutterFlowTheme.of(context).titleLarge.fontWeight,
                                                                                                            fontStyle: FlutterFlowTheme.of(context).titleLarge.fontStyle,
                                                                                                          ),
                                                                                                      playbackDurationTextStyle: FlutterFlowTheme.of(context).labelMedium.override(
                                                                                                            font: GoogleFonts.inter(
                                                                                                              fontWeight: FlutterFlowTheme.of(context).labelMedium.fontWeight,
                                                                                                              fontStyle: FlutterFlowTheme.of(context).labelMedium.fontStyle,
                                                                                                            ),
                                                                                                            color: FlutterFlowTheme.of(context).secondaryBackground,
                                                                                                            letterSpacing: 0.0,
                                                                                                            fontWeight: FlutterFlowTheme.of(context).labelMedium.fontWeight,
                                                                                                            fontStyle: FlutterFlowTheme.of(context).labelMedium.fontStyle,
                                                                                                          ),
                                                                                                      fillColor: const Color(0x00000000),
                                                                                                      playbackButtonColor: FlutterFlowTheme.of(context).alternate,
                                                                                                      activeTrackColor: FlutterFlowTheme.of(context).primary,
                                                                                                      inactiveTrackColor: FlutterFlowTheme.of(context).alternate,
                                                                                                      elevation: 0.0,
                                                                                                      playInBackground: PlayInBackground.enabled,
                                                                                                    ),
                                                                                                  ],
                                                                                                ),
                                                                                              ),
                                                                                            ),
                                                                                          ],
                                                                                        ),
                                                                                      ),
                                                                                    if ((objectSubmissionRecord.body != '') && !(objectSubmissionRecord.imagesextra.isNotEmpty))
                                                                                      Padding(
                                                                                        padding: const EdgeInsets.all(8.0),
                                                                                        child: Column(
                                                                                          mainAxisSize: MainAxisSize.min,
                                                                                          mainAxisAlignment: MainAxisAlignment.start,
                                                                                          crossAxisAlignment: CrossAxisAlignment.start,
                                                                                          children: [
                                                                                            Align(
                                                                                              alignment: const AlignmentDirectional(-1.0, -1.0),
                                                                                              child: Padding(
                                                                                                padding: const EdgeInsets.all(8.0),
                                                                                                child: Text(
                                                                                                  objectSubmissionRecord.header,
                                                                                                  maxLines: 3,
                                                                                                  style: FlutterFlowTheme.of(context).titleMedium.override(
                                                                                                        font: GoogleFonts.inter(
                                                                                                          fontWeight: FlutterFlowTheme.of(context).titleMedium.fontWeight,
                                                                                                          fontStyle: FlutterFlowTheme.of(context).titleMedium.fontStyle,
                                                                                                        ),
                                                                                                        color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                        letterSpacing: 0.0,
                                                                                                        fontWeight: FlutterFlowTheme.of(context).titleMedium.fontWeight,
                                                                                                        fontStyle: FlutterFlowTheme.of(context).titleMedium.fontStyle,
                                                                                                      ),
                                                                                                ),
                                                                                              ),
                                                                                            ),
                                                                                            Divider(
                                                                                              thickness: 1.0,
                                                                                              indent: 12.0,
                                                                                              endIndent: 12.0,
                                                                                              color: FlutterFlowTheme.of(context).alternate,
                                                                                            ),
                                                                                            Align(
                                                                                              alignment: const AlignmentDirectional(-1.0, -1.0),
                                                                                              child: Padding(
                                                                                                padding: const EdgeInsetsDirectional.fromSTEB(8.0, 0.0, 8.0, 8.0),
                                                                                                child: Text(
                                                                                                  objectSubmissionRecord.body,
                                                                                                  maxLines: 4,
                                                                                                  style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                        font: GoogleFonts.inter(
                                                                                                          fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                                          fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                        ),
                                                                                                        color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                        letterSpacing: 0.0,
                                                                                                        fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                      ),
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
                                                                        },
                                                                      ),
                                                                    );
                                                                  },
                                                                ),
                                                              ),
                                                            ),
                                                          ),
                                                        ],
                                                      ),
                                                    ),
                                                  ),
                                                ),
                                                KeepAliveWidgetWrapper(
                                                  builder: (context) =>
                                                      Container(
                                                    width: double.infinity,
                                                    height: double.infinity,
                                                    decoration: BoxDecoration(
                                                      borderRadius:
                                                          BorderRadius.circular(
                                                              8.0),
                                                    ),
                                                    child: Padding(
                                                      padding:
                                                          const EdgeInsetsDirectional
                                                              .fromSTEB(4.0,
                                                              0.0, 4.0, 0.0),
                                                      child: Column(
                                                        mainAxisSize:
                                                            MainAxisSize.max,
                                                        crossAxisAlignment:
                                                            CrossAxisAlignment
                                                                .stretch,
                                                        children: [
                                                          Padding(
                                                            padding:
                                                                const EdgeInsetsDirectional
                                                                    .fromSTEB(
                                                                    0.0,
                                                                    8.0,
                                                                    0.0,
                                                                    4.0),
                                                            child: Container(
                                                              width: double
                                                                  .infinity,
                                                              height: 50.0,
                                                              decoration:
                                                                  BoxDecoration(
                                                                color: FlutterFlowTheme.of(
                                                                        context)
                                                                    .alternate,
                                                                borderRadius:
                                                                    BorderRadius
                                                                        .circular(
                                                                            8.0),
                                                              ),
                                                              child: Padding(
                                                                padding:
                                                                    const EdgeInsetsDirectional
                                                                        .fromSTEB(
                                                                        6.0,
                                                                        6.0,
                                                                        6.0,
                                                                        6.0),
                                                                child: Row(
                                                                  mainAxisSize:
                                                                      MainAxisSize
                                                                          .max,
                                                                  mainAxisAlignment:
                                                                      MainAxisAlignment
                                                                          .spaceBetween,
                                                                  children: [
                                                                    Expanded(
                                                                      child:
                                                                          Container(
                                                                        width:
                                                                            100.0,
                                                                        height:
                                                                            100.0,
                                                                        decoration:
                                                                            const BoxDecoration(
                                                                          color:
                                                                              Color(0x00FFFFFF),
                                                                        ),
                                                                        alignment: const AlignmentDirectional(
                                                                            0.0,
                                                                            0.0),
                                                                        child:
                                                                            Row(
                                                                          mainAxisSize:
                                                                              MainAxisSize.max,
                                                                          mainAxisAlignment:
                                                                              MainAxisAlignment.spaceEvenly,
                                                                          children: [
                                                                            Align(
                                                                              alignment: const AlignmentDirectional(0.0, 0.0),
                                                                              child: FlutterFlowChoiceChips(
                                                                                options: [
                                                                                  ChipData(FFLocalizations.of(context).getText(
                                                                                    'v1k7vb29' /* All */,
                                                                                  )),
                                                                                  ChipData(
                                                                                      FFLocalizations.of(context).getText(
                                                                                        'ww6dzqw4' /* Popular */,
                                                                                      ),
                                                                                      Icons.star_rounded),
                                                                                  ChipData(
                                                                                      FFLocalizations.of(context).getText(
                                                                                        'jot3aa1n' /* Catalogues */,
                                                                                      ),
                                                                                      Icons.menu_book)
                                                                                ],
                                                                                onChanged: (val) => safeSetState(() => _model.choiceChipsItemsValue = val?.firstOrNull),
                                                                                selectedChipStyle: ChipStyle(
                                                                                  backgroundColor: FlutterFlowTheme.of(context).secondary,
                                                                                  textStyle: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                        font: GoogleFonts.inter(
                                                                                          fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                          fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                        ),
                                                                                        color: FlutterFlowTheme.of(context).info,
                                                                                        letterSpacing: 0.0,
                                                                                        fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                      ),
                                                                                  iconColor: FlutterFlowTheme.of(context).info,
                                                                                  iconSize: 16.0,
                                                                                  elevation: 0.0,
                                                                                  borderRadius: BorderRadius.circular(8.0),
                                                                                ),
                                                                                unselectedChipStyle: ChipStyle(
                                                                                  backgroundColor: FlutterFlowTheme.of(context).secondaryBackground,
                                                                                  textStyle: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                        font: GoogleFonts.inter(
                                                                                          fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                          fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                        ),
                                                                                        color: FlutterFlowTheme.of(context).secondaryText,
                                                                                        letterSpacing: 0.0,
                                                                                        fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                      ),
                                                                                  iconColor: FlutterFlowTheme.of(context).secondaryText,
                                                                                  iconSize: 16.0,
                                                                                  elevation: 0.0,
                                                                                  borderRadius: BorderRadius.circular(8.0),
                                                                                ),
                                                                                chipSpacing: 6.0,
                                                                                rowSpacing: 4.0,
                                                                                multiselect: false,
                                                                                initialized: _model.choiceChipsItemsValue != null,
                                                                                alignment: WrapAlignment.center,
                                                                                controller: _model.choiceChipsItemsValueController ??= FormFieldController<List<String>>(
                                                                                  [
                                                                                    FFLocalizations.of(context).getText(
                                                                                      'wok3m9lz' /* All */,
                                                                                    )
                                                                                  ],
                                                                                ),
                                                                                wrapped: false,
                                                                              ),
                                                                            ),
                                                                          ],
                                                                        ),
                                                                      ),
                                                                    ),
                                                                  ],
                                                                ),
                                                              ),
                                                            ),
                                                          ),
                                                          if (_model
                                                                  .choiceChipsItemsValue !=
                                                              'Catalogues')
                                                            SizedBox(
                                                              child: Padding(
                                                                padding:
                                                                    const EdgeInsetsDirectional
                                                                        .fromSTEB(
                                                                        8.0,
                                                                        12.0,
                                                                        8.0,
                                                                        0.0),
                                                                child: StreamBuilder<
                                                                    List<
                                                                        SubmissionRecord>>(
                                                                  stream:
                                                                      querySubmissionRecord(
                                                                    queryBuilder: (submissionRecord) => submissionRecord
                                                                        .where(
                                                                          'poster',
                                                                          isEqualTo:
                                                                              publicpageUsersRecord.reference,
                                                                        )
                                                                        .where(
                                                                          'type_object',
                                                                          isEqualTo:
                                                                              'Item',
                                                                        )
                                                                        .orderBy('date', descending: true),
                                                                  ),
                                                                  builder: (context,
                                                                      snapshot) {
                                                                    // Customize what your widget looks like when it's loading.
                                                                    if (!snapshot
                                                                        .hasData) {
                                                                      return Center(
                                                                        child:
                                                                            SizedBox(
                                                                          width:
                                                                              50.0,
                                                                          height:
                                                                              50.0,
                                                                          child:
                                                                              CircularProgressIndicator(
                                                                            valueColor:
                                                                                AlwaysStoppedAnimation<Color>(
                                                                              FlutterFlowTheme.of(context).primary,
                                                                            ),
                                                                          ),
                                                                        ),
                                                                      );
                                                                    }
                                                                    List<SubmissionRecord>
                                                                        staggeredViewSubmissionRecordList =
                                                                        snapshot
                                                                            .data!;

                                                                    return MasonryGridView
                                                                        .builder(
                                                                      shrinkWrap:
                                                                          true,
                                                                      physics:
                                                                          const NeverScrollableScrollPhysics(),
                                                                      gridDelegate:
                                                                          SliverSimpleGridDelegateWithFixedCrossAxisCount(
                                                                        crossAxisCount: () {
                                                                          if (MediaQuery.sizeOf(context).width < 430.0) {
                                                                            return 2;
                                                                          } else if (MediaQuery.sizeOf(context).width <= 1024.0) {
                                                                            return 3;
                                                                          } else {
                                                                            return 5;
                                                                          }
                                                                        }(),
                                                                      ),
                                                                      crossAxisSpacing:
                                                                          12.0,
                                                                      mainAxisSpacing:
                                                                          8.0,
                                                                      itemCount:
                                                                          staggeredViewSubmissionRecordList
                                                                              .length,
                                                                      itemBuilder:
                                                                          (context,
                                                                              staggeredViewIndex) {
                                                                        final staggeredViewSubmissionRecord =
                                                                            staggeredViewSubmissionRecordList[staggeredViewIndex];
                                                                        return Builder(
                                                                          builder: (context) =>
                                                                              StreamBuilder<SubmissionRecord>(
                                                                            stream:
                                                                                SubmissionRecord.getDocument(staggeredViewSubmissionRecord.reference),
                                                                            builder:
                                                                                (context, snapshot) {
                                                                              // Customize what your widget looks like when it's loading.
                                                                              if (!snapshot.hasData) {
                                                                                return Center(
                                                                                  child: SizedBox(
                                                                                    width: 50.0,
                                                                                    height: 50.0,
                                                                                    child: CircularProgressIndicator(
                                                                                      valueColor: AlwaysStoppedAnimation<Color>(
                                                                                        FlutterFlowTheme.of(context).primary,
                                                                                      ),
                                                                                    ),
                                                                                  ),
                                                                                );
                                                                              }

                                                                              final objectSubmissionRecord = snapshot.data!;

                                                                              return InkWell(
                                                                                splashColor: Colors.transparent,
                                                                                focusColor: Colors.transparent,
                                                                                hoverColor: Colors.transparent,
                                                                                highlightColor: Colors.transparent,
                                                                                onTap: () async {
                                                                                  await showDialog(
                                                                                    context: context,
                                                                                    builder: (dialogContext) {
                                                                                      return Dialog(
                                                                                        elevation: 0,
                                                                                        insetPadding: EdgeInsets.zero,
                                                                                        backgroundColor: Colors.transparent,
                                                                                        alignment: const AlignmentDirectional(0.0, 0.0).resolve(Directionality.of(context)),
                                                                                        child: GestureDetector(
                                                                                          onTap: () {
                                                                                            FocusScope.of(dialogContext).unfocus();
                                                                                            FocusManager.instance.primaryFocus?.unfocus();
                                                                                          },
                                                                                          child: ObjectWidget(
                                                                                            object: objectSubmissionRecord,
                                                                                          ),
                                                                                        ),
                                                                                      );
                                                                                    },
                                                                                  );
                                                                                },
                                                                                child: ClipRRect(
                                                                                  borderRadius: BorderRadius.circular(8.0),
                                                                                  child: Container(
                                                                                    width: 100.0,
                                                                                    height: 240.0,
                                                                                    constraints: const BoxConstraints(
                                                                                      maxHeight: 320.0,
                                                                                    ),
                                                                                    decoration: BoxDecoration(
                                                                                      color: FlutterFlowTheme.of(context).secondaryBackground,
                                                                                      boxShadow: const [
                                                                                        BoxShadow(
                                                                                          blurRadius: 5.0,
                                                                                          color: Color(0x3416202A),
                                                                                          offset: Offset(
                                                                                            0.0,
                                                                                            3.0,
                                                                                          ),
                                                                                        )
                                                                                      ],
                                                                                      borderRadius: BorderRadius.circular(8.0),
                                                                                    ),
                                                                                    child: Stack(
                                                                                      children: [
                                                                                        if (staggeredViewSubmissionRecord.imagesextra.where((e) => e.isNotEmpty).isNotEmpty)
                                                                                          Builder(
                                                                                            builder: (context) {
                                                                                              final images = staggeredViewSubmissionRecord.imagesextra.where((e) => e.isNotEmpty).toList();

                                                                                              return SingleChildScrollView(
                                                                                                scrollDirection: Axis.horizontal,
                                                                                                child: Row(
                                                                                                  mainAxisSize: MainAxisSize.max,
                                                                                                  crossAxisAlignment: CrossAxisAlignment.stretch,
                                                                                                  children: List.generate(images.length, (imagesIndex) {
                                                                                                    final imagesItem = images[imagesIndex];
                                                                                                    return Visibility(
                                                                                                      visible: imagesItem != '',
                                                                                                      child: ClipRRect(
                                                                                                        borderRadius: BorderRadius.circular(8.0),
                                                                                                        child: Image.network(
                                                                                                          imagesItem,
                                                                                                          width: 200.0,
                                                                                                          height: 100.0,
                                                                                                          fit: BoxFit.cover,
                                                                                                        ),
                                                                                                      ),
                                                                                                    );
                                                                                                  }),
                                                                                                ),
                                                                                              );
                                                                                            },
                                                                                          ),
                                                                                        Align(
                                                                                          alignment: const AlignmentDirectional(1.0, 1.0),
                                                                                          child: Padding(
                                                                                            padding: const EdgeInsets.all(8.0),
                                                                                            child: Container(
                                                                                              height: 40.0,
                                                                                              decoration: BoxDecoration(
                                                                                                color: FlutterFlowTheme.of(context).alternate,
                                                                                                borderRadius: BorderRadius.circular(12.0),
                                                                                              ),
                                                                                              child: Padding(
                                                                                                padding: const EdgeInsets.all(4.0),
                                                                                                child: Row(
                                                                                                  mainAxisSize: MainAxisSize.max,
                                                                                                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                                                                                                  children: [
                                                                                                    Padding(
                                                                                                      padding: const EdgeInsetsDirectional.fromSTEB(8.0, 0.0, 0.0, 0.0),
                                                                                                      child: Icon(
                                                                                                        Icons.toll,
                                                                                                        color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                        size: 16.0,
                                                                                                      ),
                                                                                                    ),
                                                                                                    Expanded(
                                                                                                      child: Text(
                                                                                                        staggeredViewSubmissionRecord.refvalue,
                                                                                                        textAlign: TextAlign.center,
                                                                                                        style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                              font: GoogleFonts.inter(
                                                                                                                fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                              ),
                                                                                                              color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                              letterSpacing: 0.0,
                                                                                                              fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                            ),
                                                                                                      ),
                                                                                                    ),
                                                                                                    Padding(
                                                                                                      padding: const EdgeInsetsDirectional.fromSTEB(0.0, 0.0, 4.0, 0.0),
                                                                                                      child: FlutterFlowIconButton(
                                                                                                        borderRadius: 12.0,
                                                                                                        buttonSize: 30.0,
                                                                                                        fillColor: FlutterFlowTheme.of(context).secondaryBackground,
                                                                                                        icon: Icon(
                                                                                                          Icons.add,
                                                                                                          color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                          size: 14.0,
                                                                                                        ),
                                                                                                        onPressed: () async {
                                                                                                          _model.addToBagitems(objectSubmissionRecord.reference);
                                                                                                          _model.addToRefvaluesinbagitems(objectSubmissionRecord.objectRefvalue);
                                                                                                        },
                                                                                                      ),
                                                                                                    ),
                                                                                                  ].divide(const SizedBox(width: 8.0)),
                                                                                                ),
                                                                                              ),
                                                                                            ),
                                                                                          ),
                                                                                        ),
                                                                                      ],
                                                                                    ),
                                                                                  ),
                                                                                ),
                                                                              );
                                                                            },
                                                                          ),
                                                                        );
                                                                      },
                                                                    );
                                                                  },
                                                                ),
                                                              ),
                                                            ),
                                                          if (_model
                                                                  .choiceChipsItemsValue ==
                                                              'Catalogues')
                                                            Expanded(
                                                              child: Column(
                                                                mainAxisSize:
                                                                    MainAxisSize
                                                                        .max,
                                                                crossAxisAlignment:
                                                                    CrossAxisAlignment
                                                                        .stretch,
                                                                children: [
                                                                  StreamBuilder<
                                                                      List<
                                                                          CatalogueRecord>>(
                                                                    stream:
                                                                        queryCatalogueRecord(
                                                                      queryBuilder:
                                                                          (catalogueRecord) =>
                                                                              catalogueRecord.where(
                                                                        'user_ref',
                                                                        isEqualTo:
                                                                            currentUserReference,
                                                                      ),
                                                                      limit: 8,
                                                                    ),
                                                                    builder:
                                                                        (context,
                                                                            snapshot) {
                                                                      // Customize what your widget looks like when it's loading.
                                                                      if (!snapshot
                                                                          .hasData) {
                                                                        return Center(
                                                                          child:
                                                                              SizedBox(
                                                                            width:
                                                                                50.0,
                                                                            height:
                                                                                50.0,
                                                                            child:
                                                                                CircularProgressIndicator(
                                                                              valueColor: AlwaysStoppedAnimation<Color>(
                                                                                FlutterFlowTheme.of(context).primary,
                                                                              ),
                                                                            ),
                                                                          ),
                                                                        );
                                                                      }
                                                                      List<CatalogueRecord>
                                                                          listViewCatalogueRecordList =
                                                                          snapshot
                                                                              .data!;

                                                                      return ListView
                                                                          .separated(
                                                                        padding:
                                                                            const EdgeInsets.fromLTRB(
                                                                          0,
                                                                          4.0,
                                                                          0,
                                                                          0,
                                                                        ),
                                                                        shrinkWrap:
                                                                            true,
                                                                        scrollDirection:
                                                                            Axis.vertical,
                                                                        itemCount:
                                                                            listViewCatalogueRecordList.length,
                                                                        separatorBuilder:
                                                                            (_, __) =>
                                                                                const SizedBox(height: 8.0),
                                                                        itemBuilder:
                                                                            (context,
                                                                                listViewIndex) {
                                                                          final listViewCatalogueRecord =
                                                                              listViewCatalogueRecordList[listViewIndex];
                                                                          return Builder(
                                                                            builder: (context) =>
                                                                                InkWell(
                                                                              splashColor: Colors.transparent,
                                                                              focusColor: Colors.transparent,
                                                                              hoverColor: Colors.transparent,
                                                                              highlightColor: Colors.transparent,
                                                                              onTap: () async {
                                                                                await showDialog(
                                                                                  context: context,
                                                                                  builder: (dialogContext) {
                                                                                    return Dialog(
                                                                                      elevation: 0,
                                                                                      insetPadding: EdgeInsets.zero,
                                                                                      backgroundColor: Colors.transparent,
                                                                                      alignment: const AlignmentDirectional(0.0, 0.0).resolve(Directionality.of(context)),
                                                                                      child: GestureDetector(
                                                                                        onTap: () {
                                                                                          FocusScope.of(dialogContext).unfocus();
                                                                                          FocusManager.instance.primaryFocus?.unfocus();
                                                                                        },
                                                                                        child: const CatalogueWidget(),
                                                                                      ),
                                                                                    );
                                                                                  },
                                                                                );
                                                                              },
                                                                              child: Container(
                                                                                width: 100.0,
                                                                                height: 157.0,
                                                                                decoration: BoxDecoration(
                                                                                  color: FlutterFlowTheme.of(context).alternate,
                                                                                  borderRadius: BorderRadius.circular(12.0),
                                                                                ),
                                                                                child: Padding(
                                                                                  padding: const EdgeInsets.all(12.0),
                                                                                  child: Column(
                                                                                    mainAxisSize: MainAxisSize.min,
                                                                                    crossAxisAlignment: CrossAxisAlignment.stretch,
                                                                                    children: [
                                                                                      Row(
                                                                                        mainAxisSize: MainAxisSize.min,
                                                                                        mainAxisAlignment: MainAxisAlignment.start,
                                                                                        crossAxisAlignment: CrossAxisAlignment.center,
                                                                                        children: [
                                                                                          Padding(
                                                                                            padding: const EdgeInsetsDirectional.fromSTEB(8.0, 0.0, 0.0, 0.0),
                                                                                            child: Text(
                                                                                              listViewCatalogueRecord.catalalogueName,
                                                                                              style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                    font: GoogleFonts.inter(
                                                                                                      fontWeight: FontWeight.w600,
                                                                                                      fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                    ),
                                                                                                    color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                    fontSize: 18.0,
                                                                                                    letterSpacing: 0.0,
                                                                                                    fontWeight: FontWeight.w600,
                                                                                                    fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                  ),
                                                                                            ),
                                                                                          ),
                                                                                          Padding(
                                                                                            padding: const EdgeInsetsDirectional.fromSTEB(12.0, 0.0, 0.0, 0.0),
                                                                                            child: Text(
                                                                                              FFLocalizations.of(context).getText(
                                                                                                'a7rmp25t' /* # : */,
                                                                                              ),
                                                                                              style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                    font: GoogleFonts.inter(
                                                                                                      fontWeight: FontWeight.normal,
                                                                                                      fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                    ),
                                                                                                    color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                    fontSize: 14.0,
                                                                                                    letterSpacing: 0.0,
                                                                                                    fontWeight: FontWeight.normal,
                                                                                                    fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                  ),
                                                                                            ),
                                                                                          ),
                                                                                          Text(
                                                                                            listViewCatalogueRecord.totalItems.toString(),
                                                                                            style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                  font: GoogleFonts.inter(
                                                                                                    fontWeight: FontWeight.w600,
                                                                                                    fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                  ),
                                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                  fontSize: 16.0,
                                                                                                  letterSpacing: 0.0,
                                                                                                  fontWeight: FontWeight.w600,
                                                                                                  fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                ),
                                                                                          ),
                                                                                          Expanded(
                                                                                            child: Align(
                                                                                              alignment: const AlignmentDirectional(1.0, 0.0),
                                                                                              child: Padding(
                                                                                                padding: const EdgeInsets.all(8.0),
                                                                                                child: Icon(
                                                                                                  Icons.edit,
                                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                  size: 24.0,
                                                                                                ),
                                                                                              ),
                                                                                            ),
                                                                                          ),
                                                                                        ],
                                                                                      ),
                                                                                      Expanded(
                                                                                        child: Container(
                                                                                          width: 100.0,
                                                                                          decoration: BoxDecoration(
                                                                                            borderRadius: BorderRadius.circular(8.0),
                                                                                            border: Border.all(
                                                                                              color: FlutterFlowTheme.of(context).secondaryBackground,
                                                                                              width: 1.0,
                                                                                            ),
                                                                                          ),
                                                                                          child: const Padding(
                                                                                            padding: EdgeInsets.all(8.0),
                                                                                            child: Column(
                                                                                              mainAxisSize: MainAxisSize.max,
                                                                                              mainAxisAlignment: MainAxisAlignment.center,
                                                                                              crossAxisAlignment: CrossAxisAlignment.stretch,
                                                                                              children: [],
                                                                                            ),
                                                                                          ),
                                                                                        ),
                                                                                      ),
                                                                                    ].divide(const SizedBox(height: 8.0)),
                                                                                  ),
                                                                                ),
                                                                              ),
                                                                            ),
                                                                          );
                                                                        },
                                                                      );
                                                                    },
                                                                  ),
                                                                ].divide(
                                                                    const SizedBox(
                                                                        height:
                                                                            8.0)),
                                                              ),
                                                            ),
                                                        ].divide(const SizedBox(
                                                            height: 4.0)),
                                                      ),
                                                    ),
                                                  ),
                                                ),
                                              ],
                                            ),
                                          ),
                                        ],
                                      ),
                                      Align(
                                        alignment: const AlignmentDirectional(
                                            1.0, -1.0),
                                        child: Container(
                                          width: 120.0,
                                          height: 46.0,
                                          decoration: BoxDecoration(
                                            color: FlutterFlowTheme.of(context)
                                                .alternate,
                                            borderRadius:
                                                BorderRadius.circular(12.0),
                                            border: Border.all(
                                              color:
                                                  FlutterFlowTheme.of(context)
                                                      .alternate,
                                              width: 1.0,
                                            ),
                                          ),
                                          child: Stack(
                                            children: [
                                              if (responsiveVisibility(
                                                context: context,
                                                phone: false,
                                              ))
                                                Builder(
                                                  builder: (context) {
                                                    if (!publicpageUsersRecord
                                                        .pinnedUsers
                                                        .contains(
                                                            currentUserReference)) {
                                                      return FFButtonWidget(
                                                        onPressed: () async {
                                                          await publicpageUsersRecord
                                                              .reference
                                                              .update({
                                                            ...mapToFirestore(
                                                              {
                                                                'pinned_users':
                                                                    FieldValue
                                                                        .arrayUnion([
                                                                  currentUserReference
                                                                ]),
                                                              },
                                                            ),
                                                          });

                                                          await currentUserReference!
                                                              .update({
                                                            ...mapToFirestore(
                                                              {
                                                                'user_pins':
                                                                    FieldValue
                                                                        .arrayUnion([
                                                                  publicpageUsersRecord
                                                                      .reference
                                                                ]),
                                                              },
                                                            ),
                                                          });
                                                          // Removed manual isPinned state change.
                                                          safeSetState(() {});
                                                        },
                                                        text:
                                                            FFLocalizations.of(
                                                                    context)
                                                                .getText(
                                                          '7nv9h663' /* unpin */,
                                                        ),
                                                        icon: const FaIcon(
                                                          FontAwesomeIcons
                                                              .mapPin,
                                                          size: 14.0,
                                                        ),
                                                        options:
                                                            FFButtonOptions(
                                                          width:
                                                              double.infinity,
                                                          height:
                                                              double.infinity,
                                                          padding:
                                                              const EdgeInsetsDirectional
                                                                  .fromSTEB(
                                                                  16.0,
                                                                  0.0,
                                                                  16.0,
                                                                  0.0),
                                                          iconPadding:
                                                              const EdgeInsetsDirectional
                                                                  .fromSTEB(
                                                                  0.0,
                                                                  0.0,
                                                                  0.0,
                                                                  0.0),
                                                          iconColor:
                                                              FlutterFlowTheme.of(
                                                                      context)
                                                                  .secondaryText,
                                                          color: FlutterFlowTheme
                                                                  .of(context)
                                                              .alternate,
                                                          textStyle:
                                                              FlutterFlowTheme.of(
                                                                      context)
                                                                  .titleSmall
                                                                  .override(
                                                                    font: GoogleFonts
                                                                        .inter(
                                                                      fontWeight: FlutterFlowTheme.of(
                                                                              context)
                                                                          .titleSmall
                                                                          .fontWeight,
                                                                      fontStyle: FlutterFlowTheme.of(
                                                                              context)
                                                                          .titleSmall
                                                                          .fontStyle,
                                                                    ),
                                                                    color: FlutterFlowTheme.of(
                                                                            context)
                                                                        .secondaryText,
                                                                    letterSpacing:
                                                                        0.0,
                                                                    fontWeight: FlutterFlowTheme.of(
                                                                            context)
                                                                        .titleSmall
                                                                        .fontWeight,
                                                                    fontStyle: FlutterFlowTheme.of(
                                                                            context)
                                                                        .titleSmall
                                                                        .fontStyle,
                                                                  ),
                                                          elevation: 0.0,
                                                          borderRadius:
                                                              BorderRadius
                                                                  .circular(
                                                                      12.0),
                                                        ),
                                                      );
                                                    } else {
                                                      return FFButtonWidget(
                                                        onPressed: () async {
                                                          await publicpageUsersRecord
                                                              .reference
                                                              .update({
                                                            ...mapToFirestore(
                                                              {
                                                                'pinned_users':
                                                                    FieldValue
                                                                        .arrayRemove([
                                                                  currentUserReference
                                                                ]),
                                                              },
                                                            ),
                                                          });

                                                          await currentUserReference!
                                                              .update({
                                                            ...mapToFirestore(
                                                              {
                                                                'user_pins':
                                                                    FieldValue
                                                                        .arrayRemove([
                                                                  publicpageUsersRecord
                                                                      .reference
                                                                ]),
                                                              },
                                                            ),
                                                          });
                                                          // Removed manual isPinned state change.
                                                          safeSetState(() {});
                                                        },
                                                        text:
                                                            FFLocalizations.of(
                                                                    context)
                                                                .getText(
                                                          'orfoq8sw' /* unpin */,
                                                        ),
                                                        icon: const FaIcon(
                                                          FontAwesomeIcons
                                                              .mapPin,
                                                          size: 14.0,
                                                        ),
                                                        options:
                                                            FFButtonOptions(
                                                          width:
                                                              double.infinity,
                                                          height:
                                                              double.infinity,
                                                          padding:
                                                              const EdgeInsetsDirectional
                                                                  .fromSTEB(
                                                                  16.0,
                                                                  0.0,
                                                                  16.0,
                                                                  0.0),
                                                          iconPadding:
                                                              const EdgeInsetsDirectional
                                                                  .fromSTEB(
                                                                  0.0,
                                                                  0.0,
                                                                  0.0,
                                                                  0.0),
                                                          iconColor:
                                                              FlutterFlowTheme.of(
                                                                      context)
                                                                  .secondary,
                                                          color: FlutterFlowTheme
                                                                  .of(context)
                                                              .secondaryBackground,
                                                          textStyle:
                                                              FlutterFlowTheme.of(
                                                                      context)
                                                                  .titleSmall
                                                                  .override(
                                                                    font: GoogleFonts
                                                                        .inter(
                                                                      fontWeight: FlutterFlowTheme.of(
                                                                              context)
                                                                          .titleSmall
                                                                          .fontWeight,
                                                                      fontStyle: FlutterFlowTheme.of(
                                                                              context)
                                                                          .titleSmall
                                                                          .fontStyle,
                                                                    ),
                                                                    color: FlutterFlowTheme.of(
                                                                            context)
                                                                        .secondaryText,
                                                                    letterSpacing:
                                                                        0.0,
                                                                    fontWeight: FlutterFlowTheme.of(
                                                                            context)
                                                                        .titleSmall
                                                                        .fontWeight,
                                                                    fontStyle: FlutterFlowTheme.of(
                                                                            context)
                                                                        .titleSmall
                                                                        .fontStyle,
                                                                  ),
                                                          elevation: 0.0,
                                                          borderRadius:
                                                              BorderRadius
                                                                  .circular(
                                                                      12.0),
                                                        ),
                                                      );
                                                    }
                                                  },
                                                ),
                                              Row(
                                                mainAxisSize: MainAxisSize.max,
                                                children: [
                                                  if (!publicpageUsersRecord
                                                      .pinnedUsers
                                                      .contains(
                                                          currentUserReference))
                                                    Flexible(
                                                      child: FFButtonWidget(
                                                        onPressed: () async {
                                                          await publicpageUsersRecord
                                                              .reference
                                                              .update({
                                                            ...mapToFirestore(
                                                              {
                                                                'pinned_users':
                                                                    FieldValue
                                                                        .arrayUnion([
                                                                  currentUserReference
                                                                ]),
                                                              },
                                                            ),
                                                          });

                                                          await currentUserReference!
                                                              .update({
                                                            ...mapToFirestore(
                                                              {
                                                                'user_pins':
                                                                    FieldValue
                                                                        .arrayUnion([
                                                                  publicpageUsersRecord
                                                                      .reference
                                                                ]),
                                                              },
                                                            ),
                                                          });
                                                        },
                                                        text:
                                                            FFLocalizations.of(
                                                                    context)
                                                                .getText(
                                                          'c7daje6u' /* pin */,
                                                        ),
                                                        icon: const FaIcon(
                                                          FontAwesomeIcons
                                                              .mapPin,
                                                          size: 14.0,
                                                        ),
                                                        options:
                                                            FFButtonOptions(
                                                          width:
                                                              double.infinity,
                                                          height:
                                                              double.infinity,
                                                          padding:
                                                              const EdgeInsetsDirectional
                                                                  .fromSTEB(
                                                                  16.0,
                                                                  0.0,
                                                                  16.0,
                                                                  0.0),
                                                          iconPadding:
                                                              const EdgeInsetsDirectional
                                                                  .fromSTEB(
                                                                  0.0,
                                                                  0.0,
                                                                  0.0,
                                                                  0.0),
                                                          iconColor:
                                                              FlutterFlowTheme.of(
                                                                      context)
                                                                  .secondary,
                                                          color: FlutterFlowTheme
                                                                  .of(context)
                                                              .secondaryBackground,
                                                          textStyle:
                                                              FlutterFlowTheme.of(
                                                                      context)
                                                                  .titleSmall
                                                                  .override(
                                                                    font: GoogleFonts
                                                                        .inter(
                                                                      fontWeight: FlutterFlowTheme.of(
                                                                              context)
                                                                          .titleSmall
                                                                          .fontWeight,
                                                                      fontStyle: FlutterFlowTheme.of(
                                                                              context)
                                                                          .titleSmall
                                                                          .fontStyle,
                                                                    ),
                                                                    color: FlutterFlowTheme.of(
                                                                            context)
                                                                        .secondaryText,
                                                                    letterSpacing:
                                                                        0.0,
                                                                    fontWeight: FlutterFlowTheme.of(
                                                                            context)
                                                                        .titleSmall
                                                                        .fontWeight,
                                                                    fontStyle: FlutterFlowTheme.of(
                                                                            context)
                                                                        .titleSmall
                                                                        .fontStyle,
                                                                  ),
                                                          elevation: 0.0,
                                                          borderRadius:
                                                              BorderRadius
                                                                  .circular(
                                                                      12.0),
                                                        ),
                                                      ),
                                                    ),
                                                  if (publicpageUsersRecord
                                                      .pinnedUsers
                                                      .contains(
                                                          currentUserReference))
                                                    Flexible(
                                                      child: FFButtonWidget(
                                                        onPressed: () async {
                                                          await publicpageUsersRecord
                                                              .reference
                                                              .update({
                                                            ...mapToFirestore(
                                                              {
                                                                'pinned_users':
                                                                    FieldValue
                                                                        .arrayRemove([
                                                                  currentUserReference
                                                                ]),
                                                              },
                                                            ),
                                                          });

                                                          await currentUserReference!
                                                              .update({
                                                            ...mapToFirestore(
                                                              {
                                                                'user_pins':
                                                                    FieldValue
                                                                        .arrayRemove([
                                                                  publicpageUsersRecord
                                                                      .reference
                                                                ]),
                                                              },
                                                            ),
                                                          });
                                                        },
                                                        text:
                                                            FFLocalizations.of(
                                                                    context)
                                                                .getText(
                                                          'orfoq8sw' /* unpin */,
                                                        ),
                                                        icon: const FaIcon(
                                                          FontAwesomeIcons
                                                              .mapPin,
                                                          size: 14.0,
                                                        ),
                                                        options:
                                                            FFButtonOptions(
                                                          width:
                                                              double.infinity,
                                                          height:
                                                              double.infinity,
                                                          padding:
                                                              const EdgeInsetsDirectional
                                                                  .fromSTEB(
                                                                  16.0,
                                                                  0.0,
                                                                  16.0,
                                                                  0.0),
                                                          iconPadding:
                                                              const EdgeInsetsDirectional
                                                                  .fromSTEB(
                                                                  0.0,
                                                                  0.0,
                                                                  0.0,
                                                                  0.0),
                                                          iconColor:
                                                              FlutterFlowTheme.of(
                                                                      context)
                                                                  .secondary,
                                                          color: FlutterFlowTheme
                                                                  .of(context)
                                                              .secondaryBackground,
                                                          textStyle:
                                                              FlutterFlowTheme.of(
                                                                      context)
                                                                  .titleSmall
                                                                  .override(
                                                                    font: GoogleFonts
                                                                        .inter(
                                                                      fontWeight: FlutterFlowTheme.of(
                                                                              context)
                                                                          .titleSmall
                                                                          .fontWeight,
                                                                      fontStyle: FlutterFlowTheme.of(
                                                                              context)
                                                                          .titleSmall
                                                                          .fontStyle,
                                                                    ),
                                                                    color: FlutterFlowTheme.of(
                                                                            context)
                                                                        .secondaryText,
                                                                    letterSpacing:
                                                                        0.0,
                                                                    fontWeight: FlutterFlowTheme.of(
                                                                            context)
                                                                        .titleSmall
                                                                        .fontWeight,
                                                                    fontStyle: FlutterFlowTheme.of(
                                                                            context)
                                                                        .titleSmall
                                                                        .fontStyle,
                                                                  ),
                                                          elevation: 0.0,
                                                          borderRadius:
                                                              BorderRadius
                                                                  .circular(
                                                                      12.0),
                                                        ),
                                                      ),
                                                    ),
                                                ],
                                              ),
                                            ],
                                          ),
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              ].divide(const SizedBox(height: 8.0)),
                            ),
                            if ((_model.bagitems.length != null) &&
                                responsiveVisibility(
                                  context: context,
                                  phone: false,
                                ))
                              Align(
                                alignment: const AlignmentDirectional(0.0, 1.0),
                                child: Padding(
                                  padding: const EdgeInsetsDirectional.fromSTEB(
                                      0.0, 0.0, 0.0, 24.0),
                                  child: Column(
                                    mainAxisSize: MainAxisSize.min,
                                    mainAxisAlignment: MainAxisAlignment.end,
                                    crossAxisAlignment:
                                        CrossAxisAlignment.stretch,
                                    children: [
                                      Padding(
                                        padding: const EdgeInsetsDirectional
                                            .fromSTEB(8.0, 0.0, 8.0, 0.0),
                                        child: Container(
                                          width: 100.0,
                                          height: 560.0,
                                          decoration: BoxDecoration(
                                            color: FlutterFlowTheme.of(context)
                                                .alternate,
                                            boxShadow: const [
                                              BoxShadow(
                                                blurRadius: 4.0,
                                                color: Color(0x33000000),
                                                offset: Offset(
                                                  0.0,
                                                  2.0,
                                                ),
                                              )
                                            ],
                                            borderRadius:
                                                BorderRadius.circular(12.0),
                                          ),
                                          child: Form(
                                            key: _model.formKey,
                                            autovalidateMode:
                                                AutovalidateMode.disabled,
                                            child: Padding(
                                              padding:
                                                  const EdgeInsets.all(8.0),
                                              child: SizedBox(
                                                width: double.infinity,
                                                height: 500.0,
                                                child: PageView(
                                                  controller: _model
                                                          .bagPageViewController ??=
                                                      PageController(
                                                          initialPage: 0),
                                                  scrollDirection:
                                                      Axis.horizontal,
                                                  children: [
                                                    Column(
                                                      mainAxisSize:
                                                          MainAxisSize.max,
                                                      children: [
                                                        Padding(
                                                          padding:
                                                              const EdgeInsetsDirectional
                                                                  .fromSTEB(
                                                                  0.0,
                                                                  8.0,
                                                                  0.0,
                                                                  0.0),
                                                          child: Row(
                                                            mainAxisSize:
                                                                MainAxisSize
                                                                    .max,
                                                            mainAxisAlignment:
                                                                MainAxisAlignment
                                                                    .spaceBetween,
                                                            children: [
                                                              Padding(
                                                                padding:
                                                                    const EdgeInsetsDirectional
                                                                        .fromSTEB(
                                                                        24.0,
                                                                        0.0,
                                                                        0.0,
                                                                        0.0),
                                                                child: Text(
                                                                  FFLocalizations.of(
                                                                          context)
                                                                      .getText(
                                                                    'qn31det6' /* Your bag */,
                                                                  ),
                                                                  style: FlutterFlowTheme.of(
                                                                          context)
                                                                      .bodyMedium
                                                                      .override(
                                                                        font: GoogleFonts
                                                                            .montserrat(
                                                                          fontWeight:
                                                                              FontWeight.w500,
                                                                          fontStyle: FlutterFlowTheme.of(context)
                                                                              .bodyMedium
                                                                              .fontStyle,
                                                                        ),
                                                                        color: FlutterFlowTheme.of(context)
                                                                            .secondaryText,
                                                                        fontSize:
                                                                            20.0,
                                                                        letterSpacing:
                                                                            0.0,
                                                                        fontWeight:
                                                                            FontWeight.w500,
                                                                        fontStyle: FlutterFlowTheme.of(context)
                                                                            .bodyMedium
                                                                            .fontStyle,
                                                                      ),
                                                                ),
                                                              ),
                                                            ],
                                                          ),
                                                        ),
                                                        Expanded(
                                                          child: Padding(
                                                            padding:
                                                                const EdgeInsetsDirectional
                                                                    .fromSTEB(
                                                                    8.0,
                                                                    0.0,
                                                                    8.0,
                                                                    0.0),
                                                            child: Container(
                                                              width: double
                                                                  .infinity,
                                                              height: 100.0,
                                                              decoration:
                                                                  BoxDecoration(
                                                                color: FlutterFlowTheme.of(
                                                                        context)
                                                                    .secondaryBackground,
                                                                borderRadius:
                                                                    BorderRadius
                                                                        .circular(
                                                                            12.0),
                                                              ),
                                                              child: Column(
                                                                mainAxisSize:
                                                                    MainAxisSize
                                                                        .min,
                                                                mainAxisAlignment:
                                                                    MainAxisAlignment
                                                                        .spaceBetween,
                                                                children: [
                                                                  SingleChildScrollView(
                                                                    child:
                                                                        Column(
                                                                      mainAxisSize:
                                                                          MainAxisSize
                                                                              .min,
                                                                      children: [
                                                                        Padding(
                                                                          padding: const EdgeInsets
                                                                              .all(
                                                                              12.0),
                                                                          child:
                                                                              Builder(
                                                                            builder:
                                                                                (context) {
                                                                              final bag = _model.bagitems.unique((e) => e).toList();

                                                                              return ListView.builder(
                                                                                padding: EdgeInsets.zero,
                                                                                shrinkWrap: true,
                                                                                scrollDirection: Axis.vertical,
                                                                                itemCount: bag.length,
                                                                                itemBuilder: (context, bagIndex) {
                                                                                  final bagItem = bag[bagIndex];
                                                                                  return Padding(
                                                                                    padding: const EdgeInsetsDirectional.fromSTEB(0.0, 0.0, 0.0, 8.0),
                                                                                    child: Container(
                                                                                      decoration: BoxDecoration(
                                                                                        borderRadius: BorderRadius.circular(12.0),
                                                                                        border: Border.all(
                                                                                          color: FlutterFlowTheme.of(context).alternate,
                                                                                          width: 1.0,
                                                                                        ),
                                                                                      ),
                                                                                      child: Padding(
                                                                                        padding: const EdgeInsets.all(8.0),
                                                                                        child: StreamBuilder<SubmissionRecord>(
                                                                                          stream: SubmissionRecord.getDocument(bagItem),
                                                                                          builder: (context, snapshot) {
                                                                                            // Customize what your widget looks like when it's loading.
                                                                                            if (!snapshot.hasData) {
                                                                                              return Center(
                                                                                                child: SizedBox(
                                                                                                  width: 50.0,
                                                                                                  height: 50.0,
                                                                                                  child: CircularProgressIndicator(
                                                                                                    valueColor: AlwaysStoppedAnimation<Color>(
                                                                                                      FlutterFlowTheme.of(context).primary,
                                                                                                    ),
                                                                                                  ),
                                                                                                ),
                                                                                              );
                                                                                            }

                                                                                            final rowSubmissionRecord = snapshot.data!;

                                                                                            return Row(
                                                                                              mainAxisSize: MainAxisSize.max,
                                                                                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                                                                              children: [
                                                                                                Expanded(
                                                                                                  child: Padding(
                                                                                                    padding: const EdgeInsetsDirectional.fromSTEB(12.0, 0.0, 0.0, 0.0),
                                                                                                    child: RichText(
                                                                                                      textScaler: MediaQuery.of(context).textScaler,
                                                                                                      text: TextSpan(
                                                                                                        children: [
                                                                                                          TextSpan(
                                                                                                            text: valueOrDefault<String>(
                                                                                                              _model.bagitems.length.toString(),
                                                                                                              '1',
                                                                                                            ),
                                                                                                            style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                                  font: GoogleFonts.inter(
                                                                                                                    fontWeight: FontWeight.w600,
                                                                                                                    fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                                  ),
                                                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                                  fontSize: 16.0,
                                                                                                                  letterSpacing: 8.0,
                                                                                                                  fontWeight: FontWeight.w600,
                                                                                                                  fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                                ),
                                                                                                          ),
                                                                                                          TextSpan(
                                                                                                            text: rowSubmissionRecord.header,
                                                                                                            style: const TextStyle(),
                                                                                                          )
                                                                                                        ],
                                                                                                        style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                              font: GoogleFonts.inter(
                                                                                                                fontWeight: FontWeight.w500,
                                                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                              ),
                                                                                                              color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                              fontSize: 16.0,
                                                                                                              letterSpacing: 0.0,
                                                                                                              fontWeight: FontWeight.w500,
                                                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                            ),
                                                                                                      ),
                                                                                                      textAlign: TextAlign.start,
                                                                                                    ),
                                                                                                  ),
                                                                                                ),
                                                                                                Padding(
                                                                                                  padding: const EdgeInsetsDirectional.fromSTEB(0.0, 0.0, 8.0, 0.0),
                                                                                                  child: Text(
                                                                                                    rowSubmissionRecord.objectRefvalue.toString(),
                                                                                                    style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                          font: GoogleFonts.inter(
                                                                                                            fontWeight: FontWeight.w500,
                                                                                                            fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                          ),
                                                                                                          color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                          fontSize: 16.0,
                                                                                                          letterSpacing: 0.0,
                                                                                                          fontWeight: FontWeight.w500,
                                                                                                          fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                        ),
                                                                                                  ),
                                                                                                ),
                                                                                                Padding(
                                                                                                  padding: const EdgeInsetsDirectional.fromSTEB(0.0, 0.0, 8.0, 0.0),
                                                                                                  child: Text(
                                                                                                    _model.ordermethodsinbag.length > bagIndex ? _model.ordermethodsinbag[bagIndex] : 'Standard',
                                                                                                    style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                          font: GoogleFonts.inter(
                                                                                                            fontWeight: FontWeight.w600,
                                                                                                            fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                          ),
                                                                                                          color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                          fontSize: 16.0,
                                                                                                          letterSpacing: 0.0,
                                                                                                          fontWeight: FontWeight.w600,
                                                                                                          fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                        ),
                                                                                                  ),
                                                                                                ),
                                                                                                FlutterFlowIconButton(
                                                                                                  borderRadius: 100.0,
                                                                                                  buttonSize: 30.0,
                                                                                                  fillColor: FlutterFlowTheme.of(context).alternate,
                                                                                                  icon: Icon(
                                                                                                    Icons.close,
                                                                                                    color: FlutterFlowTheme.of(context).secondary,
                                                                                                    size: 14.0,
                                                                                                  ),
                                                                                                  onPressed: () async {
                                                                                                    _model.removeAtIndexFromBagitems(bagIndex);
                                                                                                    _model.removeAtIndexFromRefvaluesinbagitems(bagIndex);
                                                                                                    if (bagIndex < _model.ordermethodsinbag.length) {
                                                                                                      _model.removeAtIndexFromOrdermethodsinbag(bagIndex);
                                                                                                    }
                                                                                                    safeSetState(() {});
                                                                                                  },
                                                                                                ),
                                                                                              ],
                                                                                            );
                                                                                          },
                                                                                        ),
                                                                                      ),
                                                                                    ),
                                                                                  );
                                                                                },
                                                                              );
                                                                            },
                                                                          ),
                                                                        ),
                                                                      ],
                                                                    ),
                                                                  ),
                                                                ],
                                                              ),
                                                            ),
                                                          ),
                                                        ),
                                                        Padding(
                                                          padding:
                                                              const EdgeInsetsDirectional
                                                                  .fromSTEB(
                                                                  8.0,
                                                                  0.0,
                                                                  8.0,
                                                                  0.0),
                                                          child: Row(
                                                            mainAxisSize:
                                                                MainAxisSize
                                                                    .max,
                                                            mainAxisAlignment:
                                                                MainAxisAlignment
                                                                    .spaceEvenly,
                                                            children: [
                                                              Expanded(
                                                                child:
                                                                    FFButtonWidget(
                                                                  onPressed:
                                                                      () async {
                                                                    await _model
                                                                        .bagPageViewController
                                                                        ?.nextPage(
                                                                      duration: const Duration(
                                                                          milliseconds:
                                                                              300),
                                                                      curve: Curves
                                                                          .ease,
                                                                    );
                                                                  },
                                                                  text: FFLocalizations.of(
                                                                          context)
                                                                      .getText(
                                                                    'w2nk774u' /* Order */,
                                                                  ),
                                                                  icon:
                                                                      const Icon(
                                                                    Icons
                                                                        .login_sharp,
                                                                    size: 20.0,
                                                                  ),
                                                                  options:
                                                                      FFButtonOptions(
                                                                    height:
                                                                        40.0,
                                                                    padding: const EdgeInsetsDirectional
                                                                        .fromSTEB(
                                                                        16.0,
                                                                        0.0,
                                                                        16.0,
                                                                        0.0),
                                                                    iconAlignment:
                                                                        IconAlignment
                                                                            .end,
                                                                    iconPadding:
                                                                        const EdgeInsetsDirectional
                                                                            .fromSTEB(
                                                                            0.0,
                                                                            0.0,
                                                                            0.0,
                                                                            0.0),
                                                                    color: FlutterFlowTheme.of(
                                                                            context)
                                                                        .secondaryBackground,
                                                                    textStyle: FlutterFlowTheme.of(
                                                                            context)
                                                                        .titleSmall
                                                                        .override(
                                                                          font:
                                                                              GoogleFonts.montserrat(
                                                                            fontWeight:
                                                                                FlutterFlowTheme.of(context).titleSmall.fontWeight,
                                                                            fontStyle:
                                                                                FlutterFlowTheme.of(context).titleSmall.fontStyle,
                                                                          ),
                                                                          color:
                                                                              FlutterFlowTheme.of(context).secondaryText,
                                                                          fontSize:
                                                                              14.0,
                                                                          letterSpacing:
                                                                              0.0,
                                                                          fontWeight: FlutterFlowTheme.of(context)
                                                                              .titleSmall
                                                                              .fontWeight,
                                                                          fontStyle: FlutterFlowTheme.of(context)
                                                                              .titleSmall
                                                                              .fontStyle,
                                                                        ),
                                                                    elevation:
                                                                        0.0,
                                                                    borderRadius:
                                                                        BorderRadius.circular(
                                                                            8.0),
                                                                  ),
                                                                ),
                                                              ),
                                                            ].divide(
                                                                const SizedBox(
                                                                    width:
                                                                        8.0)),
                                                          ),
                                                        ),
                                                      ]
                                                          .divide(
                                                              const SizedBox(
                                                                  height: 8.0))
                                                          .addToEnd(
                                                              const SizedBox(
                                                                  height: 8.0)),
                                                    ),
                                                    Column(
                                                      mainAxisSize:
                                                          MainAxisSize.max,
                                                      children: [
                                                        Expanded(
                                                          child: Padding(
                                                            padding:
                                                                const EdgeInsets
                                                                    .all(8.0),
                                                            child: Container(
                                                              decoration:
                                                                  BoxDecoration(
                                                                color: FlutterFlowTheme.of(
                                                                        context)
                                                                    .secondaryBackground,
                                                                borderRadius:
                                                                    BorderRadius
                                                                        .circular(
                                                                            12.0),
                                                              ),
                                                              child: Padding(
                                                                padding:
                                                                    const EdgeInsets
                                                                        .all(
                                                                        8.0),
                                                                child: Column(
                                                                  mainAxisSize:
                                                                      MainAxisSize
                                                                          .max,
                                                                  children: [
                                                                    Padding(
                                                                      padding: const EdgeInsetsDirectional
                                                                          .fromSTEB(
                                                                          8.0,
                                                                          8.0,
                                                                          8.0,
                                                                          0.0),
                                                                      child:
                                                                          Row(
                                                                        mainAxisSize:
                                                                            MainAxisSize.max,
                                                                        mainAxisAlignment:
                                                                            MainAxisAlignment.spaceBetween,
                                                                        children: [
                                                                          Container(
                                                                            width:
                                                                                32.0,
                                                                            height:
                                                                                32.0,
                                                                            decoration:
                                                                                const BoxDecoration(
                                                                              color: Color(0xADFF6A73),
                                                                              shape: BoxShape.circle,
                                                                            ),
                                                                            child:
                                                                                Icon(
                                                                              Icons.description_rounded,
                                                                              color: FlutterFlowTheme.of(context).secondaryBackground,
                                                                              size: 20.0,
                                                                            ),
                                                                          ),
                                                                          Expanded(
                                                                            child:
                                                                                Padding(
                                                                              padding: const EdgeInsetsDirectional.fromSTEB(12.0, 0.0, 0.0, 0.0),
                                                                              child: Text(
                                                                                FFLocalizations.of(context).getText(
                                                                                  'n04chxo4' /* Order description */,
                                                                                ),
                                                                                style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                      font: GoogleFonts.plusJakartaSans(
                                                                                        fontWeight: FontWeight.bold,
                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                      ),
                                                                                      color: const Color(0xFF15161E),
                                                                                      fontSize: 14.0,
                                                                                      letterSpacing: 0.0,
                                                                                      fontWeight: FontWeight.bold,
                                                                                      fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                    ),
                                                                              ),
                                                                            ),
                                                                          ),
                                                                        ],
                                                                      ),
                                                                    ),
                                                                    Expanded(
                                                                      child:
                                                                          Padding(
                                                                        padding: const EdgeInsets
                                                                            .all(
                                                                            12.0),
                                                                        child:
                                                                            Builder(
                                                                          builder:
                                                                              (context) {
                                                                            final bag =
                                                                                _model.bagitems.unique((e) => e).toList();

                                                                            return ListView.separated(
                                                                              padding: EdgeInsets.zero,
                                                                              shrinkWrap: true,
                                                                              scrollDirection: Axis.vertical,
                                                                              itemCount: bag.length,
                                                                              separatorBuilder: (_, __) => const SizedBox(height: 4.0),
                                                                              itemBuilder: (context, bagIndex) {
                                                                                final bagItem = bag[bagIndex];
                                                                                return Container(
                                                                                  decoration: BoxDecoration(
                                                                                    color: FlutterFlowTheme.of(context).alternate,
                                                                                    borderRadius: BorderRadius.circular(12.0),
                                                                                  ),
                                                                                  child: Padding(
                                                                                    padding: const EdgeInsets.all(8.0),
                                                                                    child: StreamBuilder<SubmissionRecord>(
                                                                                      stream: SubmissionRecord.getDocument(bagItem),
                                                                                      builder: (context, snapshot) {
                                                                                        // Customize what your widget looks like when it's loading.
                                                                                        if (!snapshot.hasData) {
                                                                                          return Center(
                                                                                            child: SizedBox(
                                                                                              width: 50.0,
                                                                                              height: 50.0,
                                                                                              child: CircularProgressIndicator(
                                                                                                valueColor: AlwaysStoppedAnimation<Color>(
                                                                                                  FlutterFlowTheme.of(context).primary,
                                                                                                ),
                                                                                              ),
                                                                                            ),
                                                                                          );
                                                                                        }

                                                                                        final rowSubmissionRecord = snapshot.data!;

                                                                                        return Row(
                                                                                          mainAxisSize: MainAxisSize.max,
                                                                                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                                                                          children: [
                                                                                            Expanded(
                                                                                              child: Padding(
                                                                                                padding: const EdgeInsetsDirectional.fromSTEB(12.0, 0.0, 0.0, 0.0),
                                                                                                child: RichText(
                                                                                                  textScaler: MediaQuery.of(context).textScaler,
                                                                                                  text: TextSpan(
                                                                                                    children: [
                                                                                                      TextSpan(
                                                                                                        text: bagIndex == bagIndex
                                                                                                            ? valueOrDefault<String>(
                                                                                                                _model.bagitems.length.toString(),
                                                                                                                '1',
                                                                                                              )
                                                                                                            : '1',
                                                                                                        style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                              font: GoogleFonts.inter(
                                                                                                                fontWeight: FontWeight.w600,
                                                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                              ),
                                                                                                              color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                              fontSize: 16.0,
                                                                                                              letterSpacing: 8.0,
                                                                                                              fontWeight: FontWeight.w600,
                                                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                            ),
                                                                                                      ),
                                                                                                      TextSpan(
                                                                                                        text: rowSubmissionRecord.header,
                                                                                                        style: const TextStyle(),
                                                                                                      )
                                                                                                    ],
                                                                                                    style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                          font: GoogleFonts.inter(
                                                                                                            fontWeight: FontWeight.w500,
                                                                                                            fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                          ),
                                                                                                          color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                          fontSize: 16.0,
                                                                                                          letterSpacing: 0.0,
                                                                                                          fontWeight: FontWeight.w500,
                                                                                                          fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                        ),
                                                                                                  ),
                                                                                                  textAlign: TextAlign.start,
                                                                                                ),
                                                                                              ),
                                                                                            ),
                                                                                            Padding(
                                                                                              padding: const EdgeInsetsDirectional.fromSTEB(0.0, 0.0, 8.0, 0.0),
                                                                                              child: Text(
                                                                                                rowSubmissionRecord.objectRefvalue.toString(),
                                                                                                style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                      font: GoogleFonts.inter(
                                                                                                        fontWeight: FontWeight.w500,
                                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                      ),
                                                                                                      color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                      fontSize: 16.0,
                                                                                                      letterSpacing: 0.0,
                                                                                                      fontWeight: FontWeight.w500,
                                                                                                      fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                    ),
                                                                                              ),
                                                                                            ),
                                                                                            Padding(
                                                                                              padding: const EdgeInsetsDirectional.fromSTEB(0.0, 0.0, 8.0, 0.0),
                                                                                              child: Text(
                                                                                                _model.ordermethodsinbag.length > bagIndex ? _model.ordermethodsinbag[bagIndex] : 'Standard',
                                                                                                style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                                      font: GoogleFonts.inter(
                                                                                                        fontWeight: FontWeight.w600,
                                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                      ),
                                                                                                      color: FlutterFlowTheme.of(context).secondaryText,
                                                                                                      fontSize: 16.0,
                                                                                                      letterSpacing: 0.0,
                                                                                                      fontWeight: FontWeight.w600,
                                                                                                      fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                                    ),
                                                                                              ),
                                                                                            ),
                                                                                          ],
                                                                                        );
                                                                                      },
                                                                                    ),
                                                                                  ),
                                                                                );
                                                                              },
                                                                            );
                                                                          },
                                                                        ),
                                                                      ),
                                                                    ),
                                                                    Padding(
                                                                      padding: const EdgeInsetsDirectional
                                                                          .fromSTEB(
                                                                          12.0,
                                                                          0.0,
                                                                          12.0,
                                                                          0.0),
                                                                      child:
                                                                          Container(
                                                                        decoration:
                                                                            BoxDecoration(
                                                                          color:
                                                                              FlutterFlowTheme.of(context).alternate,
                                                                          borderRadius:
                                                                              BorderRadius.circular(12.0),
                                                                        ),
                                                                        child:
                                                                            Padding(
                                                                          padding: const EdgeInsets
                                                                              .all(
                                                                              8.0),
                                                                          child:
                                                                              Row(
                                                                            mainAxisSize:
                                                                                MainAxisSize.max,
                                                                            mainAxisAlignment:
                                                                                MainAxisAlignment.spaceBetween,
                                                                            children: [
                                                                              Expanded(
                                                                                child: Padding(
                                                                                  padding: const EdgeInsetsDirectional.fromSTEB(12.0, 0.0, 0.0, 0.0),
                                                                                  child: Text(
                                                                                    FFLocalizations.of(context).getText(
                                                                                      'b42n7wje' /* Total : */,
                                                                                    ),
                                                                                    style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                          font: GoogleFonts.inter(
                                                                                            fontWeight: FontWeight.w500,
                                                                                            fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                          ),
                                                                                          color: FlutterFlowTheme.of(context).secondaryText,
                                                                                          fontSize: 16.0,
                                                                                          letterSpacing: 0.0,
                                                                                          fontWeight: FontWeight.w500,
                                                                                          fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                        ),
                                                                                  ),
                                                                                ),
                                                                              ),
                                                                              Flexible(
                                                                                child: Padding(
                                                                                  padding: const EdgeInsetsDirectional.fromSTEB(0.0, 0.0, 8.0, 0.0),
                                                                                  child: Text(
                                                                                    functions.sumofBagitemsListRefvalues(_model.refvaluesinbagitems.sortedList(keyOf: (e) => e.toString(), desc: false).toList()).toString(),
                                                                                    style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                          font: GoogleFonts.inter(
                                                                                            fontWeight: FontWeight.w500,
                                                                                            fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                          ),
                                                                                          color: FlutterFlowTheme.of(context).secondaryText,
                                                                                          fontSize: 16.0,
                                                                                          letterSpacing: 0.0,
                                                                                          fontWeight: FontWeight.w500,
                                                                                          fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                        ),
                                                                                  ),
                                                                                ),
                                                                              ),
                                                                              Padding(
                                                                                padding: const EdgeInsetsDirectional.fromSTEB(0.0, 0.0, 8.0, 0.0),
                                                                                child: Text(
                                                                                  FFLocalizations.of(context).getText(
                                                                                    'fqagdjt7' /* Total */,
                                                                                  ),
                                                                                  style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                        font: GoogleFonts.inter(
                                                                                          fontWeight: FontWeight.w600,
                                                                                          fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                        ),
                                                                                        color: FlutterFlowTheme.of(context).secondaryText,
                                                                                        fontSize: 16.0,
                                                                                        letterSpacing: 0.0,
                                                                                        fontWeight: FontWeight.w600,
                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                      ),
                                                                                ),
                                                                              ),
                                                                            ],
                                                                          ),
                                                                        ),
                                                                      ),
                                                                    ),
                                                                    Padding(
                                                                      padding: const EdgeInsetsDirectional
                                                                          .fromSTEB(
                                                                          12.0,
                                                                          8.0,
                                                                          12.0,
                                                                          0.0),
                                                                      child:
                                                                          Container(
                                                                        width: double
                                                                            .infinity,
                                                                        decoration:
                                                                            BoxDecoration(
                                                                          color:
                                                                              FlutterFlowTheme.of(context).alternate,
                                                                          borderRadius:
                                                                              BorderRadius.circular(12.0),
                                                                        ),
                                                                        child:
                                                                            Padding(
                                                                          padding: const EdgeInsetsDirectional
                                                                              .fromSTEB(
                                                                              0.0,
                                                                              8.0,
                                                                              0.0,
                                                                              8.0),
                                                                          child:
                                                                              Row(
                                                                            mainAxisSize:
                                                                                MainAxisSize.max,
                                                                            mainAxisAlignment:
                                                                                MainAxisAlignment.spaceBetween,
                                                                            children: [
                                                                              Padding(
                                                                                padding: const EdgeInsetsDirectional.fromSTEB(12.0, 0.0, 0.0, 0.0),
                                                                                child: Text(
                                                                                  FFLocalizations.of(context).getText(
                                                                                    'i4u5yua3' /* Comment */,
                                                                                  ),
                                                                                  style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                        font: GoogleFonts.montserrat(
                                                                                          fontWeight: FontWeight.w500,
                                                                                          fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                        ),
                                                                                        color: FlutterFlowTheme.of(context).secondaryText,
                                                                                        letterSpacing: 0.0,
                                                                                        fontWeight: FontWeight.w500,
                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                      ),
                                                                                ),
                                                                              ),
                                                                              Padding(
                                                                                padding: const EdgeInsetsDirectional.fromSTEB(0.0, 0.0, 16.0, 0.0),
                                                                                child: Icon(
                                                                                  Icons.edit,
                                                                                  color: FlutterFlowTheme.of(context).secondaryText,
                                                                                  size: 18.0,
                                                                                ),
                                                                              ),
                                                                            ],
                                                                          ),
                                                                        ),
                                                                      ),
                                                                    ),
                                                                    Padding(
                                                                      padding: const EdgeInsets
                                                                          .all(
                                                                          12.0),
                                                                      child:
                                                                          Container(
                                                                        width: double
                                                                            .infinity,
                                                                        decoration:
                                                                            BoxDecoration(
                                                                          color:
                                                                              FlutterFlowTheme.of(context).alternate,
                                                                          borderRadius:
                                                                              BorderRadius.circular(12.0),
                                                                        ),
                                                                        child:
                                                                            Padding(
                                                                          padding: const EdgeInsetsDirectional
                                                                              .fromSTEB(
                                                                              8.0,
                                                                              0.0,
                                                                              8.0,
                                                                              0.0),
                                                                          child:
                                                                              Column(
                                                                            mainAxisSize:
                                                                                MainAxisSize.max,
                                                                            crossAxisAlignment:
                                                                                CrossAxisAlignment.stretch,
                                                                            children:
                                                                                [
                                                                              Padding(
                                                                                padding: const EdgeInsetsDirectional.fromSTEB(16.0, 4.0, 4.0, 0.0),
                                                                                child: Text(
                                                                                  FFLocalizations.of(context).getText(
                                                                                    'rbv2478s' /* Select order method : */,
                                                                                  ),
                                                                                  style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                        font: GoogleFonts.montserrat(
                                                                                          fontWeight: FontWeight.w500,
                                                                                          fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                        ),
                                                                                        color: FlutterFlowTheme.of(context).secondaryText,
                                                                                        letterSpacing: 0.0,
                                                                                        fontWeight: FontWeight.w500,
                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                      ),
                                                                                ),
                                                                              ),
                                                                              Padding(
                                                                                padding: const EdgeInsets.all(8.0),
                                                                                child: StreamBuilder<List<OrderMethodsRecord>>(
                                                                                  stream: queryOrderMethodsRecord(
                                                                                    queryBuilder: (orderMethodsRecord) => orderMethodsRecord.where(
                                                                                      'method_poster',
                                                                                      isEqualTo: publicpageUsersRecord.reference,
                                                                                    ),
                                                                                  ),
                                                                                  builder: (context, snapshot) {
                                                                                    // Customize what your widget looks like when it's loading.
                                                                                    if (!snapshot.hasData) {
                                                                                      return Center(
                                                                                        child: SizedBox(
                                                                                          width: 50.0,
                                                                                          height: 50.0,
                                                                                          child: CircularProgressIndicator(
                                                                                            valueColor: AlwaysStoppedAnimation<Color>(
                                                                                              FlutterFlowTheme.of(context).primary,
                                                                                            ),
                                                                                          ),
                                                                                        ),
                                                                                      );
                                                                                    }
                                                                                    List<OrderMethodsRecord> choiceChipsORDERMETHODSOrderMethodsRecordList = snapshot.data!;

                                                                                    return FlutterFlowChoiceChips(
                                                                                      options: (choiceChipsORDERMETHODSOrderMethodsRecordList.map((e) => e.methodType).toList()..add('Standard')).unique((e) => e).map((label) => ChipData(label)).toList(),
                                                                                      onChanged: (val) => safeSetState(() => _model.choiceChipsORDERMETHODSValue = val?.firstOrNull),
                                                                                      selectedChipStyle: ChipStyle(
                                                                                        backgroundColor: FlutterFlowTheme.of(context).secondary,
                                                                                        textStyle: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                              font: GoogleFonts.inter(
                                                                                                fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                              ),
                                                                                              color: FlutterFlowTheme.of(context).info,
                                                                                              fontSize: 12.0,
                                                                                              letterSpacing: 0.0,
                                                                                              fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                            ),
                                                                                        iconColor: FlutterFlowTheme.of(context).info,
                                                                                        iconSize: 16.0,
                                                                                        elevation: 0.0,
                                                                                        borderRadius: BorderRadius.circular(8.0),
                                                                                      ),
                                                                                      unselectedChipStyle: ChipStyle(
                                                                                        backgroundColor: FlutterFlowTheme.of(context).secondaryBackground,
                                                                                        textStyle: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                              font: GoogleFonts.inter(
                                                                                                fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                              ),
                                                                                              color: FlutterFlowTheme.of(context).secondaryText,
                                                                                              fontSize: 12.0,
                                                                                              letterSpacing: 0.0,
                                                                                              fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                            ),
                                                                                        iconColor: FlutterFlowTheme.of(context).secondaryText,
                                                                                        iconSize: 16.0,
                                                                                        elevation: 0.0,
                                                                                        borderRadius: BorderRadius.circular(8.0),
                                                                                      ),
                                                                                      chipSpacing: 8.0,
                                                                                      rowSpacing: 8.0,
                                                                                      multiselect: false,
                                                                                      alignment: WrapAlignment.center,
                                                                                      controller: _model.choiceChipsORDERMETHODSValueController ??= FormFieldController<List<String>>(
                                                                                        [],
                                                                                      ),
                                                                                      wrapped: false,
                                                                                    );
                                                                                  },
                                                                                ),
                                                                              ),
                                                                              Padding(
                                                                                padding: const EdgeInsetsDirectional.fromSTEB(16.0, 0.0, 4.0, 0.0),
                                                                                child: Text(
                                                                                  FFLocalizations.of(context).getText(
                                                                                    'zxia3f0n' /* Select wallet method : */,
                                                                                  ),
                                                                                  style: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                        font: GoogleFonts.montserrat(
                                                                                          fontWeight: FontWeight.w500,
                                                                                          fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                        ),
                                                                                        color: FlutterFlowTheme.of(context).secondaryText,
                                                                                        letterSpacing: 0.0,
                                                                                        fontWeight: FontWeight.w500,
                                                                                        fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                      ),
                                                                                ),
                                                                              ),
                                                                              Padding(
                                                                                padding: const EdgeInsets.all(8.0),
                                                                                child: StreamBuilder<List<WalletMethodsRecord>>(
                                                                                  stream: queryWalletMethodsRecord(
                                                                                    queryBuilder: (walletMethodsRecord) => walletMethodsRecord.where(
                                                                                      'method_poster',
                                                                                      isEqualTo: publicpageUsersRecord.reference,
                                                                                    ),
                                                                                  ),
                                                                                  builder: (context, snapshot) {
                                                                                    List<String> options = [];
                                                                                    bool currentUserHasWallet = currentUserDocument?.walletAddress.isNotEmpty == true;
                                                                                    bool publicUserHasWallet = publicpageUsersRecord.walletAddress.isNotEmpty;
                                                                                    if (currentUserHasWallet && publicUserHasWallet) {
                                                                                      options.add('Token');
                                                                                    }
                                                                                    if (snapshot.hasData) {
                                                                                      options.addAll(snapshot.data!.map((e) => e.methodType));
                                                                                    }
                                                                                    options = options.unique((e) => e).toList();

                                                                                    _model.choiceChipsWALLETMETHODSValueController ??= FormFieldController<List<String>>(
                                                                                      options.contains('Token') ? ['Token'] : (options.isNotEmpty ? [options.first] : []),
                                                                                    );

                                                                                    return FlutterFlowChoiceChips(
                                                                                      options: options.map((label) => ChipData(label)).toList(),
                                                                                      onChanged: (val) => safeSetState(() => _model.choiceChipsWALLETMETHODSValue = val?.firstOrNull),
                                                                                      selectedChipStyle: ChipStyle(
                                                                                        backgroundColor: FlutterFlowTheme.of(context).secondary,
                                                                                        textStyle: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                              font: GoogleFonts.inter(
                                                                                                fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                              ),
                                                                                              color: FlutterFlowTheme.of(context).info,
                                                                                              fontSize: 12.0,
                                                                                              letterSpacing: 0.0,
                                                                                              fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                            ),
                                                                                        iconColor: FlutterFlowTheme.of(context).info,
                                                                                        iconSize: 16.0,
                                                                                        elevation: 0.0,
                                                                                        borderRadius: BorderRadius.circular(8.0),
                                                                                      ),
                                                                                      unselectedChipStyle: ChipStyle(
                                                                                        backgroundColor: FlutterFlowTheme.of(context).secondaryBackground,
                                                                                        textStyle: FlutterFlowTheme.of(context).bodyMedium.override(
                                                                                              font: GoogleFonts.inter(
                                                                                                fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                                fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                              ),
                                                                                              color: FlutterFlowTheme.of(context).secondaryText,
                                                                                              fontSize: 12.0,
                                                                                              letterSpacing: 0.0,
                                                                                              fontWeight: FlutterFlowTheme.of(context).bodyMedium.fontWeight,
                                                                                              fontStyle: FlutterFlowTheme.of(context).bodyMedium.fontStyle,
                                                                                            ),
                                                                                        iconColor: FlutterFlowTheme.of(context).secondaryText,
                                                                                        iconSize: 16.0,
                                                                                        elevation: 0.0,
                                                                                        borderRadius: BorderRadius.circular(8.0),
                                                                                      ),
                                                                                      chipSpacing: 8.0,
                                                                                      rowSpacing: 8.0,
                                                                                      multiselect: false,
                                                                                      alignment: WrapAlignment.center,
                                                                                      controller: _model.choiceChipsWALLETMETHODSValueController!,
                                                                                      wrapped: false,
                                                                                    );
                                                                                  },
                                                                                ),
                                                                              ),
                                                                            ].divide(const SizedBox(height: 8.0)).addToStart(const SizedBox(height: 8.0)).addToEnd(const SizedBox(height: 8.0)),
                                                                          ),
                                                                        ),
                                                                      ),
                                                                    ),
                                                                  ],
                                                                ),
                                                              ),
                                                            ),
                                                          ),
                                                        ),
                                                        Padding(
                                                          padding:
                                                              const EdgeInsetsDirectional
                                                                  .fromSTEB(
                                                                  8.0,
                                                                  0.0,
                                                                  8.0,
                                                                  0.0),
                                                          child: Row(
                                                            mainAxisSize:
                                                                MainAxisSize
                                                                    .max,
                                                            mainAxisAlignment:
                                                                MainAxisAlignment
                                                                    .center,
                                                            children: [
                                                              Expanded(
                                                                flex: 2,
                                                                child:
                                                                    FFButtonWidget(
                                                                  onPressed:
                                                                      () async {
                                                                    _model.bagitems =
                                                                        [];
                                                                    safeSetState(
                                                                        () {});
                                                                  },
                                                                  text: FFLocalizations.of(
                                                                          context)
                                                                      .getText(
                                                                    'vd8j3btf' /* Cancel */,
                                                                  ),
                                                                  options:
                                                                      FFButtonOptions(
                                                                    height:
                                                                        40.0,
                                                                    padding: const EdgeInsetsDirectional
                                                                        .fromSTEB(
                                                                        16.0,
                                                                        0.0,
                                                                        16.0,
                                                                        0.0),
                                                                    iconPadding:
                                                                        const EdgeInsetsDirectional
                                                                            .fromSTEB(
                                                                            0.0,
                                                                            0.0,
                                                                            0.0,
                                                                            0.0),
                                                                    color: FlutterFlowTheme.of(
                                                                            context)
                                                                        .secondaryBackground,
                                                                    textStyle: FlutterFlowTheme.of(
                                                                            context)
                                                                        .titleSmall
                                                                        .override(
                                                                          font:
                                                                              GoogleFonts.inter(
                                                                            fontWeight:
                                                                                FlutterFlowTheme.of(context).titleSmall.fontWeight,
                                                                            fontStyle:
                                                                                FlutterFlowTheme.of(context).titleSmall.fontStyle,
                                                                          ),
                                                                          color:
                                                                              FlutterFlowTheme.of(context).secondary,
                                                                          letterSpacing:
                                                                              0.0,
                                                                          fontWeight: FlutterFlowTheme.of(context)
                                                                              .titleSmall
                                                                              .fontWeight,
                                                                          fontStyle: FlutterFlowTheme.of(context)
                                                                              .titleSmall
                                                                              .fontStyle,
                                                                        ),
                                                                    elevation:
                                                                        0.0,
                                                                    borderRadius:
                                                                        BorderRadius.circular(
                                                                            8.0),
                                                                  ),
                                                                ),
                                                              ),
                                                              Expanded(
                                                                flex: 3,
                                                                child:
                                                                    FFButtonWidget(
                                                                  onPressed:
                                                                      () async {
                                                                    bool
                                                                        currentUserHasWallet =
                                                                        currentUserDocument?.walletAddress.isNotEmpty ==
                                                                            true;
                                                                    bool
                                                                        publicUserHasWallet =
                                                                        publicpageUsersRecord
                                                                            .walletAddress
                                                                            .isNotEmpty;

                                                                    if (!currentUserHasWallet ||
                                                                        !publicUserHasWallet) {
                                                                      ScaffoldMessenger.of(
                                                                              context)
                                                                          .showSnackBar(
                                                                        SnackBar(
                                                                          content:
                                                                              Text(
                                                                            'Both users must have a registered wallet before placing an order on the network.',
                                                                            style:
                                                                                TextStyle(
                                                                              color: FlutterFlowTheme.of(context).primaryText,
                                                                            ),
                                                                          ),
                                                                          duration:
                                                                              const Duration(milliseconds: 4000),
                                                                          backgroundColor:
                                                                              FlutterFlowTheme.of(context).error,
                                                                        ),
                                                                      );
                                                                      return;
                                                                    }

                                                                    await OrderRecord
                                                                        .collection
                                                                        .doc()
                                                                        .set({
                                                                      ...createOrderRecordData(
                                                                        userRef:
                                                                            currentUserReference,
                                                                        date:
                                                                            getCurrentTimestamp,
                                                                        orderStats:
                                                                            OrderStatuses.Pending,
                                                                        publicuserRef:
                                                                            publicpageUsersRecord.reference,
                                                                        refValue: functions.sumofBagitemsListRefvalues(_model
                                                                            .refvaluesinbagitems
                                                                            .toList()),
                                                                        username: valueOrDefault(
                                                                            currentUserDocument?.username,
                                                                            ''),
                                                                        publicusername:
                                                                            publicpageUsersRecord.username,
                                                                        totalRefValue: functions.sumofBagitemsListRefvalues(_model
                                                                            .refvaluesinbagitems
                                                                            .toList()),
                                                                        paymentMethod:
                                                                            _model.choiceChipsWALLETMETHODSValue,
                                                                        key1: _model
                                                                            .ordermethodsinbag
                                                                            .toString(),
                                                                        key2: _model
                                                                            .choiceChipsORDERMETHODSValue,
                                                                      ),
                                                                      ...mapToFirestore(
                                                                        {
                                                                          'items_inorder': _model
                                                                              .bagitems
                                                                              .unique((e) => e),
                                                                        },
                                                                      ),
                                                                    });
                                                                    unawaited(
                                                                      () async {
                                                                        await currentUserDocument!
                                                                            .analyticsRef!
                                                                            .update({
                                                                          ...mapToFirestore(
                                                                            {
                                                                              'order_accum': FieldValue.increment(1),
                                                                              'orderAvgReference': FieldValue.arrayUnion([
                                                                                functions.sumofBagitemsListRefvalues(_model.refvaluesinbagitems.toList())
                                                                              ]),
                                                                            },
                                                                          ),
                                                                        });
                                                                      }(),
                                                                    );

                                                                    await stackAnalyticsRecord!
                                                                        .reference
                                                                        .update({
                                                                      ...mapToFirestore(
                                                                        {
                                                                          'order_accum':
                                                                              FieldValue.increment(1),
                                                                          'orderAvgReference':
                                                                              FieldValue.arrayUnion([
                                                                            functions.sumofBagitemsListRefvalues(_model.refvaluesinbagitems.toList())
                                                                          ]),
                                                                        },
                                                                      ),
                                                                    });

                                                                    context
                                                                        .pushNamed(
                                                                      UserpageWidget
                                                                          .routeName,
                                                                      extra: <String,
                                                                          dynamic>{
                                                                        kTransitionInfoKey:
                                                                            const TransitionInfo(
                                                                          hasTransition:
                                                                              true,
                                                                          transitionType:
                                                                              PageTransitionType.rightToLeft,
                                                                          duration:
                                                                              Duration(milliseconds: 450),
                                                                        ),
                                                                      },
                                                                    );
                                                                  },
                                                                  text: FFLocalizations.of(
                                                                          context)
                                                                      .getText(
                                                                    'ofwx8axq' /* Generate */,
                                                                  ),
                                                                  options:
                                                                      FFButtonOptions(
                                                                    height:
                                                                        40.0,
                                                                    padding: const EdgeInsetsDirectional
                                                                        .fromSTEB(
                                                                        16.0,
                                                                        0.0,
                                                                        16.0,
                                                                        0.0),
                                                                    iconPadding:
                                                                        const EdgeInsetsDirectional
                                                                            .fromSTEB(
                                                                            0.0,
                                                                            0.0,
                                                                            0.0,
                                                                            0.0),
                                                                    color: FlutterFlowTheme.of(
                                                                            context)
                                                                        .secondaryBackground,
                                                                    textStyle: FlutterFlowTheme.of(
                                                                            context)
                                                                        .titleSmall
                                                                        .override(
                                                                          font:
                                                                              GoogleFonts.inter(
                                                                            fontWeight:
                                                                                FlutterFlowTheme.of(context).titleSmall.fontWeight,
                                                                            fontStyle:
                                                                                FlutterFlowTheme.of(context).titleSmall.fontStyle,
                                                                          ),
                                                                          color:
                                                                              FlutterFlowTheme.of(context).secondaryText,
                                                                          letterSpacing:
                                                                              0.0,
                                                                          fontWeight: FlutterFlowTheme.of(context)
                                                                              .titleSmall
                                                                              .fontWeight,
                                                                          fontStyle: FlutterFlowTheme.of(context)
                                                                              .titleSmall
                                                                              .fontStyle,
                                                                        ),
                                                                    elevation:
                                                                        0.0,
                                                                    borderRadius:
                                                                        BorderRadius.circular(
                                                                            8.0),
                                                                  ),
                                                                ),
                                                              ),
                                                            ].divide(
                                                                const SizedBox(
                                                                    width:
                                                                        8.0)),
                                                          ),
                                                        ),
                                                      ].addToEnd(const SizedBox(
                                                          height: 8.0)),
                                                    ),
                                                  ],
                                                ),
                                              ),
                                            ),
                                          ),
                                        ).animateOnActionTrigger(
                                          animationsMap[
                                              'containerOnActionTriggerAnimation']!,
                                        ),
                                      ),
                                      Padding(
                                        padding: const EdgeInsetsDirectional
                                            .fromSTEB(0.0, 0.0, 12.0, 0.0),
                                        child: Row(
                                          mainAxisSize: MainAxisSize.max,
                                          mainAxisAlignment:
                                              MainAxisAlignment.end,
                                          children: [
                                            badges.Badge(
                                              badgeContent: Text(
                                                valueOrDefault<String>(
                                                  _model.bagitems.length
                                                      .toString(),
                                                  '0',
                                                ),
                                                style:
                                                    FlutterFlowTheme.of(context)
                                                        .titleSmall
                                                        .override(
                                                          font:
                                                              GoogleFonts.inter(
                                                            fontWeight:
                                                                FlutterFlowTheme.of(
                                                                        context)
                                                                    .titleSmall
                                                                    .fontWeight,
                                                            fontStyle:
                                                                FlutterFlowTheme.of(
                                                                        context)
                                                                    .titleSmall
                                                                    .fontStyle,
                                                          ),
                                                          color: Colors.white,
                                                          fontSize: 10.0,
                                                          letterSpacing: 0.0,
                                                          fontWeight:
                                                              FlutterFlowTheme.of(
                                                                      context)
                                                                  .titleSmall
                                                                  .fontWeight,
                                                          fontStyle:
                                                              FlutterFlowTheme.of(
                                                                      context)
                                                                  .titleSmall
                                                                  .fontStyle,
                                                        ),
                                              ),
                                              showBadge: true,
                                              shape: badges.BadgeShape.circle,
                                              badgeColor:
                                                  FlutterFlowTheme.of(context)
                                                      .primary,
                                              elevation: 4.0,
                                              padding:
                                                  const EdgeInsetsDirectional
                                                      .fromSTEB(
                                                      8.0, 8.0, 8.0, 8.0),
                                              position:
                                                  badges.BadgePosition.topEnd(),
                                              animationType: badges
                                                  .BadgeAnimationType.scale,
                                              toAnimate: true,
                                              child: InkWell(
                                                splashColor: Colors.transparent,
                                                focusColor: Colors.transparent,
                                                hoverColor: Colors.transparent,
                                                highlightColor:
                                                    Colors.transparent,
                                                onTap: () async {
                                                  if (FFAppState().isBagOpen) {
                                                    if (animationsMap[
                                                            'containerOnActionTriggerAnimation'] !=
                                                        null) {
                                                      animationsMap[
                                                              'containerOnActionTriggerAnimation']!
                                                          .controller
                                                          .reverse();
                                                    }
                                                    if (animationsMap[
                                                            'iconOnActionTriggerAnimation1'] !=
                                                        null) {
                                                      animationsMap[
                                                              'iconOnActionTriggerAnimation1']!
                                                          .controller
                                                          .reverse();
                                                    }
                                                    if (animationsMap[
                                                            'iconOnActionTriggerAnimation2'] !=
                                                        null) {
                                                      animationsMap[
                                                              'iconOnActionTriggerAnimation2']!
                                                          .controller
                                                          .reverse();
                                                    }
                                                    FFAppState().isBagOpen =
                                                        false;
                                                  } else {
                                                    if (animationsMap[
                                                            'containerOnActionTriggerAnimation'] !=
                                                        null) {
                                                      animationsMap[
                                                              'containerOnActionTriggerAnimation']!
                                                          .controller
                                                          .forward(from: 0.0);
                                                    }
                                                    if (animationsMap[
                                                            'iconOnActionTriggerAnimation1'] !=
                                                        null) {
                                                      animationsMap[
                                                              'iconOnActionTriggerAnimation1']!
                                                          .controller
                                                          .forward(from: 0.0);
                                                    }
                                                    if (animationsMap[
                                                            'iconOnActionTriggerAnimation2'] !=
                                                        null) {
                                                      safeSetState(() =>
                                                          hasIconTriggered2 =
                                                              true);
                                                      SchedulerBinding.instance
                                                          .addPostFrameCallback(
                                                              (_) async => animationsMap[
                                                                      'iconOnActionTriggerAnimation2']!
                                                                  .controller
                                                                  .forward(
                                                                      from:
                                                                          0.0));
                                                    }
                                                    FFAppState().isBagOpen =
                                                        true;
                                                  }
                                                },
                                                child: Container(
                                                  width: 52.0,
                                                  height: 52.0,
                                                  decoration: BoxDecoration(
                                                    color: FlutterFlowTheme.of(
                                                            context)
                                                        .secondary,
                                                    shape: BoxShape.circle,
                                                  ),
                                                  child: Padding(
                                                    padding:
                                                        const EdgeInsets.all(
                                                            4.0),
                                                    child: Stack(
                                                      children: [
                                                        Align(
                                                          alignment:
                                                              const AlignmentDirectional(
                                                                  0.0, 0.0),
                                                          child: Icon(
                                                            Icons
                                                                .shopping_bag_outlined,
                                                            color: FlutterFlowTheme
                                                                    .of(context)
                                                                .secondaryBackground,
                                                            size: 26.0,
                                                          ).animateOnActionTrigger(
                                                            animationsMap[
                                                                'iconOnActionTriggerAnimation1']!,
                                                          ),
                                                        ),
                                                        Align(
                                                          alignment:
                                                              const AlignmentDirectional(
                                                                  0.0, 0.0),
                                                          child: Icon(
                                                            Icons.close_sharp,
                                                            color: FlutterFlowTheme
                                                                    .of(context)
                                                                .secondaryBackground,
                                                            size: 26.0,
                                                          ).animateOnActionTrigger(
                                                              animationsMap[
                                                                  'iconOnActionTriggerAnimation2']!,
                                                              hasBeenTriggered:
                                                                  hasIconTriggered2),
                                                        ),
                                                      ],
                                                    ),
                                                  ),
                                                ),
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                    ].divide(const SizedBox(height: 12.0)),
                                  ),
                                ),
                              ),
                          ],
                        );
                      },
                    ),
                  ),
                ),
              ),
            ),
          ),
        );
      },
    );
  }
}
