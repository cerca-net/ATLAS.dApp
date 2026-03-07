import 'dart:ui' as ui;
import 'package:flutter/material.dart';

class GoogleFonts {
  static final config = GoogleFontsConfig();

  static TextStyle getFont(
    String fontFamily, {
    TextStyle? textStyle,
    Color? color,
    double? fontSize,
    FontWeight? fontWeight,
    FontStyle? fontStyle,
    TextDecoration? decoration,
    double? height,
    double? letterSpacing,
    double? wordSpacing,
    TextBaseline? textBaseline,
    List<Shadow>? shadows,
    List<ui.FontFeature>? fontFeatures,
    List<ui.FontVariation>? fontVariations,
    double? decorationThickness,
    Color? decorationColor,
    TextDecorationStyle? decorationStyle,
  }) {
    return _safeStyle(color, fontSize, fontWeight, fontStyle, decoration,
        height, letterSpacing, wordSpacing, shadows);
  }

  // Define the fonts used in your project with the full signature
  static TextStyle inter(
          {TextStyle? textStyle,
          Color? color,
          double? fontSize,
          FontWeight? fontWeight,
          FontStyle? fontStyle,
          TextDecoration? decoration,
          double? height,
          double? letterSpacing,
          double? wordSpacing,
          List<Shadow>? shadows}) =>
      _safeStyle(color, fontSize, fontWeight, fontStyle, decoration, height,
          letterSpacing, wordSpacing, shadows);
  static TextStyle sora(
          {TextStyle? textStyle,
          Color? color,
          double? fontSize,
          FontWeight? fontWeight,
          FontStyle? fontStyle,
          TextDecoration? decoration,
          double? height,
          double? letterSpacing,
          double? wordSpacing,
          List<Shadow>? shadows}) =>
      _safeStyle(color, fontSize, fontWeight, fontStyle, decoration, height,
          letterSpacing, wordSpacing, shadows);
  static TextStyle outfit(
          {TextStyle? textStyle,
          Color? color,
          double? fontSize,
          FontWeight? fontWeight,
          FontStyle? fontStyle,
          TextDecoration? decoration,
          double? height,
          double? letterSpacing,
          double? wordSpacing,
          List<Shadow>? shadows}) =>
      _safeStyle(color, fontSize, fontWeight, fontStyle, decoration, height,
          letterSpacing, wordSpacing, shadows);
  static TextStyle robotoMono(
          {TextStyle? textStyle,
          Color? color,
          double? fontSize,
          FontWeight? fontWeight,
          FontStyle? fontStyle,
          TextDecoration? decoration,
          double? height,
          double? letterSpacing,
          double? wordSpacing,
          List<Shadow>? shadows}) =>
      _safeStyle(color, fontSize, fontWeight, fontStyle, decoration, height,
          letterSpacing, wordSpacing, shadows);
  static TextStyle montserrat(
          {TextStyle? textStyle,
          Color? color,
          double? fontSize,
          FontWeight? fontWeight,
          FontStyle? fontStyle,
          TextDecoration? decoration,
          double? height,
          double? letterSpacing,
          double? wordSpacing,
          List<Shadow>? shadows}) =>
      _safeStyle(color, fontSize, fontWeight, fontStyle, decoration, height,
          letterSpacing, wordSpacing, shadows);
  static TextStyle plusJakartaSans(
          {TextStyle? textStyle,
          Color? color,
          double? fontSize,
          FontWeight? fontWeight,
          FontStyle? fontStyle,
          TextDecoration? decoration,
          double? height,
          double? letterSpacing,
          double? wordSpacing,
          List<Shadow>? shadows}) =>
      _safeStyle(color, fontSize, fontWeight, fontStyle, decoration, height,
          letterSpacing, wordSpacing, shadows);
  static TextStyle nunitoSans(
          {TextStyle? textStyle,
          Color? color,
          double? fontSize,
          FontWeight? fontWeight,
          FontStyle? fontStyle,
          TextDecoration? decoration,
          double? height,
          double? letterSpacing,
          double? wordSpacing,
          List<Shadow>? shadows}) =>
      _safeStyle(color, fontSize, fontWeight, fontStyle, decoration, height,
          letterSpacing, wordSpacing, shadows);
  static TextStyle poppins(
          {TextStyle? textStyle,
          Color? color,
          double? fontSize,
          FontWeight? fontWeight,
          FontStyle? fontStyle,
          TextDecoration? decoration,
          double? height,
          double? letterSpacing,
          double? wordSpacing,
          List<Shadow>? shadows}) =>
      _safeStyle(color, fontSize, fontWeight, fontStyle, decoration, height,
          letterSpacing, wordSpacing, shadows);

  static TextStyle _safeStyle(
      Color? color,
      double? fontSize,
      FontWeight? fontWeight,
      FontStyle? fontStyle,
      TextDecoration? decoration,
      double? height,
      double? letterSpacing,
      double? wordSpacing,
      List<Shadow>? shadows) {
    return TextStyle(
      fontFamily: 'sans-serif',
      color: color,
      fontSize: fontSize,
      fontWeight: fontWeight,
      fontStyle: fontStyle,
      decoration: decoration,
      height: height,
      letterSpacing: letterSpacing,
      wordSpacing: wordSpacing,
      shadows: shadows,
    );
  }
}

class GoogleFontsConfig {
  bool allowRuntimeFetching = false;
}
