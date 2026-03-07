import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

const _kLocaleStorageKey = '__locale_key__';

class FFLocalizations {
  FFLocalizations(this.locale);

  final Locale locale;

  static FFLocalizations of(BuildContext context) =>
      Localizations.of<FFLocalizations>(context, FFLocalizations)!;

  static List<String> languages() => ['en', 'es'];

  static late SharedPreferences _prefs;
  static Future initialize() async =>
      _prefs = await SharedPreferences.getInstance();
  static Future storeLocale(String locale) =>
      _prefs.setString(_kLocaleStorageKey, locale);
  static Locale? getStoredLocale() {
    final locale = _prefs.getString(_kLocaleStorageKey);
    return locale != null && locale.isNotEmpty ? createLocale(locale) : null;
  }

  String get languageCode => locale.toString();
  String? get languageShortCode =>
      _languagesWithShortCode.contains(locale.toString())
          ? '${locale.toString()}_short'
          : null;
  int get languageIndex => languages().contains(languageCode)
      ? languages().indexOf(languageCode)
      : 0;

  String getText(String key) {
    final translations = kTranslationsMap[key] ?? {};
    final translatedText = translations[locale.toString()];
    if (translatedText == null || translatedText.isEmpty) {
      return translations['en'] ?? '';
    }
    return translatedText;
  }

  String getVariableText({
    String? enText = '',
    String? esText = '',
  }) =>
      [enText, esText][languageIndex] ?? '';

  static const Set<String> _languagesWithShortCode = {
    'ar',
    'az',
    'ca',
    'cs',
    'da',
    'de',
    'dv',
    'en',
    'es',
    'et',
    'fi',
    'fr',
    'gr',
    'he',
    'hi',
    'hu',
    'it',
    'km',
    'ku',
    'mn',
    'ms',
    'no',
    'pt',
    'ro',
    'ru',
    'rw',
    'sv',
    'th',
    'uk',
    'vi',
  };
}

/// Used if the locale is not supported by GlobalMaterialLocalizations.
class FallbackMaterialLocalizationDelegate
    extends LocalizationsDelegate<MaterialLocalizations> {
  const FallbackMaterialLocalizationDelegate();

  @override
  bool isSupported(Locale locale) => _isSupportedLocale(locale);

  @override
  Future<MaterialLocalizations> load(Locale locale) async =>
      SynchronousFuture<MaterialLocalizations>(
        const DefaultMaterialLocalizations(),
      );

  @override
  bool shouldReload(FallbackMaterialLocalizationDelegate old) => false;
}

/// Used if the locale is not supported by GlobalCupertinoLocalizations.
class FallbackCupertinoLocalizationDelegate
    extends LocalizationsDelegate<CupertinoLocalizations> {
  const FallbackCupertinoLocalizationDelegate();

  @override
  bool isSupported(Locale locale) => _isSupportedLocale(locale);

  @override
  Future<CupertinoLocalizations> load(Locale locale) =>
      SynchronousFuture<CupertinoLocalizations>(
        const DefaultCupertinoLocalizations(),
      );

  @override
  bool shouldReload(FallbackCupertinoLocalizationDelegate old) => false;
}

class FFLocalizationsDelegate extends LocalizationsDelegate<FFLocalizations> {
  const FFLocalizationsDelegate();

  @override
  bool isSupported(Locale locale) => _isSupportedLocale(locale);

  @override
  Future<FFLocalizations> load(Locale locale) =>
      SynchronousFuture<FFLocalizations>(FFLocalizations(locale));

  @override
  bool shouldReload(FFLocalizationsDelegate old) => false;
}

Locale createLocale(String language) => language.contains('_')
    ? Locale.fromSubtags(
        languageCode: language.split('_').first,
        scriptCode: language.split('_').last,
      )
    : Locale(language);

bool _isSupportedLocale(Locale locale) {
  final language = locale.toString();
  return FFLocalizations.languages().contains(
    language.endsWith('_')
        ? language.substring(0, language.length - 1)
        : language,
  );
}

final kTranslationsMap = <Map<String, Map<String, String>>>[
  // Signup
  {
    '24rmtw49': {
      'en': 'Sign up',
      'es': '',
    },
    'ydpo41yy': {
      'en': 'Create an account by using the form below.',
      'es': '',
    },
    'fztojcj7': {
      'en': 'Email',
      'es': '',
    },
    'ynalme06': {
      'en': 'Username',
      'es': '',
    },
    'om2787my': {
      'en': 'Password',
      'es': '',
    },
    't2wf7c91': {
      'en': 'Create Account',
      'es': '',
    },
    'ldrckxrf': {
      'en': 'Or sign up with',
      'es': '',
    },
    'l0vok5gn': {
      'en': 'Continue with Google',
      'es': '',
    },
    '2j06gm58': {
      'en': 'Continue with Metamask',
      'es': '',
    },
    '6y6i5t19': {
      'en': 'Already have an account? ',
      'es': '',
    },
    'l9ms2xxq': {
      'en': 'Log in here',
      'es': '',
    },
    'dgp1ub85': {
      'en': 'Home',
      'es': '',
    },
  },
  // ForgotPassword
  {
    '41allicc': {
      'en': 'Forgot Password',
      'es': '',
    },
    'if6uj6h3': {
      'en':
          'Please fill out your email belo in order to recieve a reset password link.',
      'es': '',
    },
    '2ipplna8': {
      'en': 'Email',
      'es': '',
    },
    'e0c05vmd': {
      'en': 'Send Reset Link',
      'es': '',
    },
    '9rdhbjme': {
      'en': 'Home',
      'es': '',
    },
  },
  // Login
  {
    '66c84yjf': {
      'en': 'Log in',
      'es': '',
    },
    'dzwgwcti': {
      'en': 'Email',
      'es': '',
    },
    'sh9ibwbz': {
      'en': 'Password',
      'es': '',
    },
    'j5tlci9v': {
      'en': 'Continue',
      'es': '',
    },
    '0j23omtc': {
      'en': 'Continue with Google',
      'es': '',
    },
    'nope6ond': {
      'en': 'Forgot password?',
      'es': '',
    },
    'bi4p1qnt': {
      'en': 'Home',
      'es': '',
    },
  },
  // StartingPage
  {
    '73qngijf': {
      'en': 'EXPLORE',
      'es': '',
    },
    'szwv3n45': {
      'en': 'Discover your surroundings like never before',
      'es': '',
    },
    '2naolsoc': {
      'en': 'CONNECT',
      'es': '',
    },
    'fz1t8fpi': {
      'en': 'Share and discover new potential',
      'es': '',
    },
    'wrjzgjvf': {
      'en': 'INTERACT',
      'es': '',
    },
    'p2bvpaxe': {
      'en': 'Automated transaction processing',
      'es': '',
    },
    '8towxatg': {
      'en': 'Login',
      'es': '',
    },
    '6hmrfrd8': {
      'en': 'Sign up',
      'es': '',
    },
    'no1fu3db': {
      'en': 'Explore',
      'es': '',
    },
    'mfn33mid': {
      'en': 'Home',
      'es': '',
    },
  },
  // UserSurvey
  {
    'y3cebvy3': {
      'en': 'User & Profile',
      'es': '',
    },
    'hw57vcox': {
      'en': 'Upload Banner',
      'es': '',
    },
    'xoeqk0hy': {
      'en': 'Upload Profile \nPicture',
      'es': '',
    },
    '6277moed': {
      'en': 'Person',
      'es': '',
    },
    'o06nweqf': {
      'en': 'Business',
      'es': '',
    },
    'gywi2pjw': {
      'en': 'Person',
      'es': '',
    },
    '9utw0jk2': {
      'en': 'Name',
      'es': '',
    },
    '4jwhpr5h': {
      'en': 'Surname',
      'es': '',
    },
    'anx5jv0q': {
      'en': 'Bio',
      'es': '',
    },
    'z4pso157': {
      'en': 'Select your users occupation',
      'es': '',
    },
    '9498i41r': {
      'en': 'Healthcare',
      'es': '',
    },
    'lpmcl9f3': {
      'en': 'Retail',
      'es': '',
    },
    '3cno6qyl': {
      'en': 'Wholesale',
      'es': '',
    },
    'bj916lkf': {
      'en': 'Manufaturing',
      'es': '',
    },
    '3takcnep': {
      'en': 'Finance',
      'es': '',
    },
    'qhmbbk09': {
      'en': 'Technology',
      'es': '',
    },
    'u4l06res': {
      'en': 'Energy',
      'es': '',
    },
    'xs4qgq4j': {
      'en': 'Transportation',
      'es': '',
    },
    '9cm00fvq': {
      'en': 'Construction',
      'es': '',
    },
    'omwzhj0q': {
      'en': 'Education',
      'es': '',
    },
    'd7qp50tz': {
      'en': 'Hospality',
      'es': '',
    },
    'lkgpngzb': {
      'en': 'Government',
      'es': '',
    },
    'gk6io9uh': {
      'en': 'Entretainment',
      'es': '',
    },
    'aghudzi4': {
      'en': 'Design',
      'es': '',
    },
    'xu470d1f': {
      'en': 'Agriculture',
      'es': '',
    },
    '064nfovl': {
      'en': 'Natural Resources',
      'es': '',
    },
    'k4zngo00': {
      'en': 'Utilities',
      'es': '',
    },
    'l90ugofi': {
      'en': 'Professional Services',
      'es': '',
    },
    '11wwxkau': {
      'en': 'Personal Services',
      'es': '',
    },
    'qu99e17z': {
      'en': 'Security',
      'es': '',
    },
    '48aksrxo': {
      'en': 'Real Estate',
      'es': '',
    },
    '3h9kkf1k': {
      'en': 'Food',
      'es': '',
    },
    'su48hcw0': {
      'en': 'Travel',
      'es': '',
    },
    'xcrkmowp': {
      'en': 'Sports',
      'es': '',
    },
    '90bavbxd': {
      'en': ' Tell us about your interests',
      'es': '',
    },
    'no5jpekm': {
      'en': 'Travel',
      'es': '',
    },
    '4ufg73x4': {
      'en': 'Reading',
      'es': '',
    },
    'cvw50ugt': {
      'en': 'Cooking',
      'es': '',
    },
    'rr5zgldv': {
      'en': 'Outdoor',
      'es': '',
    },
    'wvk4c858': {
      'en': 'Art',
      'es': '',
    },
    'wqly0a13': {
      'en': 'Craft',
      'es': '',
    },
    'gqomjujw': {
      'en': 'Entretainment',
      'es': '',
    },
    'bw0eqkha': {
      'en': 'Wellness',
      'es': '',
    },
    '0azskocy': {
      'en': 'Media',
      'es': '',
    },
    'i7s51dmb': {
      'en': 'Sports',
      'es': '',
    },
    'dl9ijqgv': {
      'en': 'Exploring',
      'es': '',
    },
    'sdcyn9v4': {
      'en': 'Meeting',
      'es': '',
    },
    'vdi4gncr': {
      'en': 'Music',
      'es': '',
    },
    'aruqfke0': {
      'en': 'Events',
      'es': '',
    },
    '1iz0much': {
      'en': 'Learning',
      'es': '',
    },
    '5ykefu5j': {
      'en': 'Design',
      'es': '',
    },
    '2dfjm9lv': {
      'en': '',
      'es': '',
    },
    'hsjsapae': {
      'en': 'Home',
      'es': '',
    },
  },
  // UserProfile
  {
    '83nti0r2': {
      'en': 'Upload Banner',
      'es': '',
    },
    'uwkta84s': {
      'en': 'Upload Profile \nPicture',
      'es': '',
    },
    'ntjeme37': {
      'en': 'Name',
      'es': '',
    },
    '8abdhrml': {
      'en': 'Surname',
      'es': '',
    },
    'bpcd2ifp': {
      'en': 'Bio',
      'es': '',
    },
    '1iwh4ihg': {
      'en': 'Save Changes',
      'es': '',
    },
    'l649sg61': {
      'en': 'Field is required',
      'es': '',
    },
    'sprsxbyq': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'i4ufhhof': {
      'en': 'Field is required',
      'es': '',
    },
    'oyrnknux': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'divith4m': {
      'en': 'Field is required',
      'es': '',
    },
    '2jjhpd4u': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'gusnjbpa': {
      'en': 'Home',
      'es': '',
    },
  },
  // Userpage
  {
    'x28kyi58': {
      'en': ' ',
      'es': '',
    },
    'ahih9kc1': {
      'en': 'Option 1',
      'es': '',
    },
    'm7w9n0s9': {
      'en': 'Option 2',
      'es': '',
    },
    'kglak0b4': {
      'en': 'Option 3',
      'es': '',
    },
    'g67c2sxv': {
      'en': '0',
      'es': '',
    },
    'imxgzxx9': {
      'en': 'Orders',
      'es': '',
    },
    'obmdky6k': {
      'en': 'Inbox',
      'es': '',
    },
    'f1esz34w': {
      'en': 'Sent',
      'es': '',
    },
    'tmh1sj9i': {
      'en': 'Methods',
      'es': '',
    },
    'isnxbf2n': {
      'en': 'Inbox',
      'es': '',
    },
    'f4huv1ys': {
      'en': 'Method: ',
      'es': '',
    },
    'yqgfpnuc': {
      'en': 'Items in order: ',
      'es': '',
    },
    'wdxw6z6m': {
      'en': 'Method: ',
      'es': '',
    },
    'bspe9z7u': {
      'en': 'Total: ',
      'es': '',
    },
    'lymz2vnh': {
      'en': 'P2P Order Method(s)',
      'es': '',
    },
    '6ixhtfyw': {
      'en': '# I.D. : ',
      'es': '',
    },
    '8kbtl5np': {
      'en': 'Quanity: 1',
      'es': '',
    },
    '1iy5supo': {
      'en': '# Index : ',
      'es': '',
    },
    'qy61koig': {
      'en': 'Quanity: 1',
      'es': '',
    },
    'uyq9twej': {
      'en': 'Analytics',
      'es': '',
    },
    'ffcxey0n': {
      'en': 'Order',
      'es': '',
    },
    '0kvddroj': {
      'en': 'Objects',
      'es': '',
    },
    'uzenmhah': {
      'en': 'User',
      'es': '',
    },
    'ebcookjc': {
      'en': 'Objects',
      'es': '',
    },
    'hhee5bpa': {
      'en': 'Completed orders :',
      'es': '',
    },
    'dhd7d4si': {
      'en': 'Indicators :',
      'es': '',
    },
    'th5s4iq4': {
      'en': ' #',
      'es': '',
    },
    '6e6bmne8': {
      'en': 'Indicators :',
      'es': '',
    },
    'hwhozquw': {
      'en': 'Average ref. :',
      'es': '',
    },
    'yqt6a4ck': {
      'en': 'Indicators :',
      'es': '',
    },
    'dnmhg8ju': {
      'en': ' T.',
      'es': '',
    },
    '4bipyjr2': {
      'en': 'Indicators :',
      'es': '',
    },
    'c6tij4kv': {
      'en': 'Average time :',
      'es': '',
    },
    'ckavxt4e': {
      'en': 'Indicators :',
      'es': '',
    },
    'd76thnm5': {
      'en': ' min.',
      'es': '',
    },
    'agcjsr9v': {
      'en': 'Indicators :',
      'es': '',
    },
    'esjfutcy': {
      'en': 'Bag to order ratio:',
      'es': '',
    },
    'k7u8fbbv': {
      'en': 'Indicators :',
      'es': '',
    },
    'u1gcrqzf': {
      'en': ' / 1',
      'es': '',
    },
    'vcugyrt9': {
      'en': 'Indicators :',
      'es': '',
    },
    '3dqpn6vc': {
      'en': 'Order review:',
      'es': '',
    },
    '0kygx8zb': {
      'en': 'Indicators :',
      'es': '',
    },
    '79rwpb9j': {
      'en': ' / 5',
      'es': '',
    },
    'q3sqxbk5': {
      'en': 'Indicators :',
      'es': '',
    },
    'may6qt2k': {
      'en': 'Accumulated ref:',
      'es': '',
    },
    'btaq9axv': {
      'en': 'Indicators :',
      'es': '',
    },
    'w7yaa7ef': {
      'en': ' T.',
      'es': '',
    },
    'rhqpyuty': {
      'en': 'Indicators :',
      'es': '',
    },
    'swwh73oe': {
      'en': 'Vote rate',
      'es': '',
    },
    'dvp6f0s1': {
      'en': 'Pin index',
      'es': '',
    },
    'hrutzr9j': {
      'en': 'Shares',
      'es': '',
    },
    'pcm2o3qc': {
      'en': 'Impressions:  ',
      'es': '',
    },
    'llb516rt': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    '0fzxiwqw': {
      'en': '(  ',
      'es': '',
    },
    'gnfjwfai': {
      'en': '14.6',
      'es': '',
    },
    'cicxaux3': {
      'en': '  % )',
      'es': '',
    },
    '9dlnj3k9': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    '7wpu4vev': {
      'en': 'Reach:  ',
      'es': '',
    },
    'd419l5uf': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'r98tze10': {
      'en': '(  ',
      'es': '',
    },
    'havut3rl': {
      'en': '14.6',
      'es': '',
    },
    '5wjcyvnc': {
      'en': '  % )',
      'es': '',
    },
    '2e53c1n2': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'k0zon12l': {
      'en': 'Interactivity:  ',
      'es': '',
    },
    'u8cs0oba': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'kyg0mcpi': {
      'en': '(  ',
      'es': '',
    },
    '4wrgn0d6': {
      'en': '14.6',
      'es': '',
    },
    '1tv2yx6k': {
      'en': '  % )',
      'es': '',
    },
    'z9nw6iqo': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    '0xqapb7s': {
      'en': 'Top performer',
      'es': '',
    },
    'td5b4780': {
      'en': 'Interaction:  ',
      'es': '',
    },
    'dz75j3mt': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    '3cb5ipdy': {
      'en': '',
      'es': '',
    },
    'z93v6apx': {
      'en': '69',
      'es': '',
    },
    'j1km0lc3': {
      'en': '  %',
      'es': '',
    },
    'p1zd100e': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    '2p6ptbo0': {
      'en': 'User verification status:   ',
      'es': '',
    },
    'oc3qm24w': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    '1ys6gyt7': {
      'en': 'User impressions:  ',
      'es': '',
    },
    'uueru30u': {
      'en': '',
      'es': '',
    },
    'mkjwvoll': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'zmfg75ax': {
      'en': 'User pin: ',
      'es': '',
    },
    'tvtdjw3x': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'r9cppg5y': {
      'en': 'Object activity:  ',
      'es': '',
    },
    '4g6c6vc6': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'ir7m9euc': {
      'en': '',
      'es': '',
    },
    'rxvx4t4s': {
      'en': 'Avg Rating   ',
      'es': '',
    },
    'zjqwd7xk': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'gq9s9t12': {
      'en': 'MARKET INDEX :  ',
      'es': '',
    },
    'tp2mb201': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'ili4vcsr': {
      'en': '( ',
      'es': '',
    },
    'ym3fo7xk': {
      'en': '+ ',
      'es': '',
    },
    'scpa7cuh': {
      'en': '  )',
      'es': '',
    },
    'zvbdk9qm': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'ca5v0ewa': {
      'en': 'PERFORMANCE : ',
      'es': '',
    },
    '0j4q491u': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'm6odsaaz': {
      'en': '( ',
      'es': '',
    },
    'wf5k2qk2': {
      'en': '+ ',
      'es': '',
    },
    '22l45ua0': {
      'en': '  % )',
      'es': '',
    },
    'g75ls3xa': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    't5rt5r2j': {
      'en': 'AVERAGE HASH  :',
      'es': '',
    },
    'bsbm4rbq': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'rrl37m0m': {
      'en': '   H/S',
      'es': '',
    },
    'by8dafib': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'nuvjz7is': {
      'en': 'Posts',
      'es': '',
    },
    'e708w00f': {
      'en': 'Image',
      'es': '',
    },
    'rvozhkq7': {
      'en': 'Video',
      'es': '',
    },
    'um6wh61y': {
      'en': 'Audio',
      'es': '',
    },
    '3vgcopam': {
      'en': 'Text',
      'es': '',
    },
    'zv7h1wfa': {
      'en': '',
      'es': '',
    },
    'o8yrnf4n': {
      'en': 'Image',
      'es': '',
    },
    'vrwv3ph2': {
      'en': 'Items',
      'es': '',
    },
    'i9o8t90g': {
      'en': 'All',
      'es': '',
    },
    's8rcrry3': {
      'en': 'Popular',
      'es': '',
    },
    'n9c1xdrq': {
      'en': 'Catalogues',
      'es': '',
    },
    'dr0no1fn': {
      'en': 'Popular',
      'es': '',
    },
    'zn5p428u': {
      'en': '# :',
      'es': '',
    },
    '7c1vetmk': {
      'en': 'Wallet',
      'es': '',
    },
    'matvykmi': {
      'en': 'Book',
      'es': '',
    },
    'ynwxm69q': {
      'en': 'Methods',
      'es': '',
    },
    'gunbpyhb': {
      'en': 'Account',
      'es': '',
    },
    '7803th34': {
      'en': 'Book',
      'es': '',
    },
    '5en137ce': {
      'en': 'Transaction history',
      'es': '',
    },
    'vn7r2tpm': {
      'en': 'A list of historical transactions',
      'es': '',
    },
    'bdfckyy1': {
      'en': 'Reciever: ',
      'es': '',
    },
    'glhf9f71': {
      'en': 'Sender: ',
      'es': '',
    },
    'czwg08jc': {
      'en': 'ORDER I.D.  #',
      'es': '',
    },
    'rzrzheb1': {
      'en': 'REF : ',
      'es': '',
    },
    '8o1mx6c6': {
      'en': 'Reciever: ',
      'es': '',
    },
    'ufvv30se': {
      'en': 'Sender: ',
      'es': '',
    },
    'lmanv682': {
      'en': 'ORDER I.D.  #',
      'es': '',
    },
    'icfn8zgw': {
      'en': 'REF : ',
      'es': '',
    },
    'y15ubtt5': {
      'en': 'Other P2P Transaction Method(s)',
      'es': '',
    },
    '3i3onmpd': {
      'en': '# I.D. : ',
      'es': '',
    },
    '1lk9df3u': {
      'en': 'Quanity: 1',
      'es': '',
    },
    'gn9oxkhn': {
      'en': '# Account : ',
      'es': '',
    },
    'yirspgij': {
      'en': 'Quanity: 1',
      'es': '',
    },
    'q0vr7c4a': {
      'en': 'Connect Wallet',
      'es': '',
    },
    '374b7k11': {
      'en': 'SEND',
      'es': '',
    },
    '231pdd7e': {
      'en': 'RECIEVE',
      'es': '',
    },
    'uqrzw76w': {
      'en': 'HISTORY',
      'es': '',
    },
    '4plssfel': {
      'en': 'Tcoin balance',
      'es': '',
    },
    'h109n1ac': {
      'en': 'Generated D.U.',
      'es': '',
    },
    'c4yw2f3u': {
      'en': 'Available :',
      'es': '',
    },
    '41f05aue': {
      'en': 'Available :',
      'es': '',
    },
    '4y2ugccp': {
      'en': 'Mint',
      'es': '',
    },
    'lhbl6vre': {
      'en': '',
      'es': '',
    },
  },
  // Parameters
  {
    'u110qpli': {
      'en': 'Object Notifications',
      'es': '',
    },
    'igkyw93p': {
      'en':
          'Receive Push notifications from our application on a semi regular basis.',
      'es': '',
    },
    'mdkbf5i3': {
      'en': 'Order Notifications',
      'es': '',
    },
    'cxxpnl3g': {
      'en':
          'Receive email notifications from our marketing team about new features.',
      'es': '',
    },
    '8sot5dum': {
      'en': 'Order Location ',
      'es': '',
    },
    'oueeplo5': {
      'en':
          'Allow Tom to track your instant location when swapping keys in an order, this helps keep track of orders.',
      'es': '',
    },
    'mei7n1bv': {
      'en': 'Home',
      'es': '',
    },
  },
  // Newpage
  {
    'iw6r7zs6': {
      'en': 'New',
      'es': '',
    },
    '0l25dry9': {
      'en': 'Image',
      'es': '',
    },
    'gqhlc2p5': {
      'en': 'Video',
      'es': '',
    },
    'rxhumk4d': {
      'en': 'Audio',
      'es': '',
    },
    'vdfupv5r': {
      'en': 'Text',
      'es': '',
    },
    'xdbfprng': {
      'en': 'Image',
      'es': '',
    },
    'uv09qkko': {
      'en': 'Product',
      'es': '',
    },
    '10uvwcqx': {
      'en': 'Service',
      'es': '',
    },
    'dr4nqxmn': {
      'en': 'Event',
      'es': '',
    },
    'zow2pkdx': {
      'en': 'Product',
      'es': '',
    },
    '62dm787f': {
      'en': 'Photo',
      'es': '',
    },
    'e5ht3pvr': {
      'en': 'Art',
      'es': '',
    },
    '1cgupov9': {
      'en': 'Media',
      'es': '',
    },
    'wju6arih': {
      'en': 'Data',
      'es': '',
    },
    '62tg68d4': {
      'en': 'Design',
      'es': '',
    },
    's0x3reum': {
      'en': 'Graphic',
      'es': '',
    },
    'juzoggio': {
      'en': 'Fashion',
      'es': '',
    },
    'p0x8zynu': {
      'en': 'Spaces',
      'es': '',
    },
    'stvk92xs': {
      'en': 'Conceptual',
      'es': '',
    },
    'r6rejbzu': {
      'en': 'Analysis',
      'es': '',
    },
    '4ot7y919': {
      'en': 'Graph',
      'es': '',
    },
    'jdg44iu9': {
      'en': 'Statistics',
      'es': '',
    },
    '6xmmq3m8': {
      'en': 'Science',
      'es': '',
    },
    'u1m24pgf': {
      'en': 'Content',
      'es': '',
    },
    'phemfzne': {
      'en': 'News',
      'es': '',
    },
    'l5adn476': {
      'en': 'Social',
      'es': '',
    },
    'dmf1kp71': {
      'en': 'Abstract',
      'es': '',
    },
    'a33goo8l': {
      'en': 'Surreal',
      'es': '',
    },
    'haf9nb94': {
      'en': 'Conceptual',
      'es': '',
    },
    '9bve97cs': {
      'en': 'Realism',
      'es': '',
    },
    'pfcxy9ln': {
      'en': 'Social',
      'es': '',
    },
    '9ck5jmjk': {
      'en': 'Nature',
      'es': '',
    },
    'y2bxmgdk': {
      'en': 'Lens',
      'es': '',
    },
    'lvayql2v': {
      'en': 'City',
      'es': '',
    },
    's0vqvxh5': {
      'en': 'Lifestyle',
      'es': '',
    },
    '26ihtbl0': {
      'en': 'Add footnote here...',
      'es': '',
    },
    'j1aca0wc': {
      'en': 'Submit',
      'es': '',
    },
    'm7u2wvq5': {
      'en': 'Field is required',
      'es': '',
    },
    'zop2zy01': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'gus79ymc': {
      'en': 'Movie',
      'es': '',
    },
    'm2eqjyyy': {
      'en': 'Short',
      'es': '',
    },
    'lieadc2q': {
      'en': 'Clip',
      'es': '',
    },
    'bsubmgai': {
      'en': 'Documentary',
      'es': '',
    },
    'mwd4txia': {
      'en': 'Entretainment',
      'es': '',
    },
    'l555j5xw': {
      'en': 'Education',
      'es': '',
    },
    '7gerbmeg': {
      'en': 'Creativity',
      'es': '',
    },
    'kp71z61w': {
      'en': 'Series',
      'es': '',
    },
    'r65j9frg': {
      'en': 'Sketch',
      'es': '',
    },
    'taxkrbfl': {
      'en': 'Science',
      'es': '',
    },
    'hrkc7p0m': {
      'en': 'Social',
      'es': '',
    },
    'owx6wuww': {
      'en': 'Animation',
      'es': '',
    },
    'zct4gjs2': {
      'en': 'Entretainment',
      'es': '',
    },
    'rm9qq4a6': {
      'en': 'Entretainment',
      'es': '',
    },
    'fwt79vza': {
      'en': 'Info',
      'es': '',
    },
    'cq19rk71': {
      'en': 'Talk',
      'es': '',
    },
    'vj30b74v': {
      'en': 'Opinion',
      'es': '',
    },
    'nee80c5b': {
      'en': 'Add footnote here...',
      'es': '',
    },
    'w2wab6rb': {
      'en': 'Submit',
      'es': '',
    },
    '5w2eua8h': {
      'en': 'Field is required',
      'es': '',
    },
    '3crt0hpu': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'obrg3as0': {
      'en': 'Music',
      'es': '',
    },
    'nitq11u6': {
      'en': 'Podcast',
      'es': '',
    },
    'phwiqlql': {
      'en': 'Note',
      'es': '',
    },
    'wfbmm3ev': {
      'en': 'Electronic',
      'es': '',
    },
    '3ni1na80': {
      'en': 'Acoustic',
      'es': '',
    },
    'ohf4qwi4': {
      'en': 'Record',
      'es': '',
    },
    'q1kuhemd': {
      'en': 'Interview',
      'es': '',
    },
    '5hmzhwpm': {
      'en': 'Show',
      'es': '',
    },
    'jnchvg9e': {
      'en': 'Comment',
      'es': '',
    },
    'y68yukkj': {
      'en': 'Comment',
      'es': '',
    },
    '9ann79ao': {
      'en': 'Opinion',
      'es': '',
    },
    'cl4h2vv8': {
      'en': 'Header...',
      'es': '',
    },
    'u5xn9nlx': {
      'en': 'Submit',
      'es': '',
    },
    'zn3fgz8h': {
      'en': 'Field is required',
      'es': '',
    },
    'xihceedh': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'qguz9xc2': {
      'en': 'Header...',
      'es': '',
    },
    'or2yy3mb': {
      'en': 'Body...',
      'es': '',
    },
    'u3cn6s21': {
      'en': 'Article',
      'es': '',
    },
    '024zcif4': {
      'en': 'Quote',
      'es': '',
    },
    'hrxtxkcz': {
      'en': 'Story',
      'es': '',
    },
    '7vk8bvbr': {
      'en': 'Report',
      'es': '',
    },
    'qgfkon9n': {
      'en': 'Review',
      'es': '',
    },
    'ngw8xueh': {
      'en': 'Insight',
      'es': '',
    },
    'i06fcxls': {
      'en': 'Comment',
      'es': '',
    },
    's2ryrk39': {
      'en': 'Note',
      'es': '',
    },
    'k1zk0u8f': {
      'en': 'Fiction',
      'es': '',
    },
    'g9681hte': {
      'en': 'Real',
      'es': '',
    },
    '9y3mo43e': {
      'en': 'Submit',
      'es': '',
    },
    'wcle3am4': {
      'en': 'Field is required',
      'es': '',
    },
    'aje0m3gw': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    '17a369bh': {
      'en': 'Tech',
      'es': '',
    },
    '1ausvomf': {
      'en': 'Food',
      'es': '',
    },
    'zkjxj7oh': {
      'en': 'Wellness',
      'es': '',
    },
    'kevv3htw': {
      'en': 'Home & Office',
      'es': '',
    },
    'ctnxaf2w': {
      'en': 'Apparel',
      'es': '',
    },
    'lvkib4l8': {
      'en': 'Vehicle',
      'es': '',
    },
    'k7mjfxy4': {
      'en': 'Tools',
      'es': '',
    },
    'p3gu0uxj': {
      'en': 'Others',
      'es': '',
    },
    'hfr416q6': {
      'en': 'Electronics',
      'es': '',
    },
    '6zqmufdw': {
      'en': 'Mechanics',
      'es': '',
    },
    't3vn8wgb': {
      'en': 'Ingridients',
      'es': '',
    },
    'ncwspkng': {
      'en': 'Dish',
      'es': '',
    },
    'qsf4maj8': {
      'en': 'Beverage',
      'es': '',
    },
    '9wtq21v4': {
      'en': 'Snack',
      'es': '',
    },
    '5fylhjlf': {
      'en': 'Sports',
      'es': '',
    },
    '3gyv2wsd': {
      'en': 'Beauty',
      'es': '',
    },
    'krijnrnh': {
      'en': 'Wellbeing',
      'es': '',
    },
    'pc17r7kl': {
      'en': 'Appliances',
      'es': '',
    },
    '7n7l6evt': {
      'en': 'Furniture',
      'es': '',
    },
    'i5abhkqy': {
      'en': 'Outdoors',
      'es': '',
    },
    '9dxewyth': {
      'en': 'Supplies',
      'es': '',
    },
    'ahucjs2v': {
      'en': 'Clothing',
      'es': '',
    },
    'tm2irqc5': {
      'en': 'Footweat',
      'es': '',
    },
    'sp20xmti': {
      'en': 'Accesories',
      'es': '',
    },
    '2uiqma4e': {
      'en': 'Parts',
      'es': '',
    },
    'u9jiilmg': {
      'en': 'Accesories',
      'es': '',
    },
    'o7mdnt9x': {
      'en': 'Manual',
      'es': '',
    },
    'c0x6dh3w': {
      'en': 'Electric',
      'es': '',
    },
    '5g3q6d0g': {
      'en': 'Mechanic',
      'es': '',
    },
    'j0ankzwf': {
      'en': 'Materials',
      'es': '',
    },
    'yycudpnz': {
      'en': 'Utilities',
      'es': '',
    },
    'ouj22ajb': {
      'en': 'Collectible',
      'es': '',
    },
    '1kaajki4': {
      'en': 'Add Title here...',
      'es': '',
    },
    'md7c6owb': {
      'en': 'Ref.',
      'es': '',
    },
    'n2ueuudf': {
      'en': 'Add description here...',
      'es': '',
    },
    'rpbxl6ge': {
      'en': 'Submit',
      'es': '',
    },
    'k5mcha95': {
      'en': 'Field is required',
      'es': '',
    },
    'dsilfr2g': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'r2duy674': {
      'en': 'Transport',
      'es': '',
    },
    'kbe9l697': {
      'en': 'Education',
      'es': '',
    },
    'qyxluk75': {
      'en': 'Home',
      'es': '',
    },
    'vu46hopx': {
      'en': 'Proffesional',
      'es': '',
    },
    'ymcio5io': {
      'en': 'Technical',
      'es': '',
    },
    'tuk5qjtg': {
      'en': 'Care',
      'es': '',
    },
    'i86evncm': {
      'en': 'Design',
      'es': '',
    },
    'oj24xz9o': {
      'en': 'Sharing',
      'es': '',
    },
    'ns9q8vuw': {
      'en': 'Logistics',
      'es': '',
    },
    't0uverv5': {
      'en': 'Particular',
      'es': '',
    },
    'i0iiqqbp': {
      'en': 'Tutoring',
      'es': '',
    },
    '5wwybi2s': {
      'en': 'Course',
      'es': '',
    },
    'fqgywmwa': {
      'en': 'Training',
      'es': '',
    },
    'sprqbp0l': {
      'en': 'Mainteinance',
      'es': '',
    },
    'y4o5zykr': {
      'en': 'Repair',
      'es': '',
    },
    '2ck6qidm': {
      'en': 'Installation',
      'es': '',
    },
    'wactdvsz': {
      'en': 'Legal',
      'es': '',
    },
    'glsrz6d2': {
      'en': 'Accounting',
      'es': '',
    },
    'ehdkn4rj': {
      'en': 'Consulting',
      'es': '',
    },
    'hpmq6b6n': {
      'en': 'I.T.',
      'es': '',
    },
    'kpas26dy': {
      'en': 'Electrical',
      'es': '',
    },
    '12nsbzj8': {
      'en': 'Mechanical',
      'es': '',
    },
    'j5xapt9v': {
      'en': 'Craft',
      'es': '',
    },
    'tfpiceon': {
      'en': 'Self',
      'es': '',
    },
    'bfqc8b3u': {
      'en': 'Nursing',
      'es': '',
    },
    'ygxyv26j': {
      'en': 'Therapy',
      'es': '',
    },
    't6v445s1': {
      'en': 'Digital',
      'es': '',
    },
    '62axivug': {
      'en': 'Spatial',
      'es': '',
    },
    'zysq2snp': {
      'en': 'Fashion',
      'es': '',
    },
    'gjx2e171': {
      'en': 'Industrial',
      'es': '',
    },
    'qu87l10i': {
      'en': 'Add Title here...',
      'es': '',
    },
    'b5pbrunj': {
      'en': 'Ref.',
      'es': '',
    },
    'b1230r4i': {
      'en': 'Activity type...',
      'es': '',
    },
    'msy1y5ik': {
      'en': 'Search...',
      'es': '',
    },
    'h1vzr7zn': {
      'en': 'Open',
      'es': '',
    },
    'kcpdik3i': {
      'en': 'Booking',
      'es': '',
    },
    '7ur6z0n8': {
      'en': 'Add description here...',
      'es': '',
    },
    '8rej0tiv': {
      'en': 'Submit',
      'es': '',
    },
    'kacwg20p': {
      'en': 'Field is required',
      'es': '',
    },
    'rxfqx94z': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    't9prjjt8': {
      'en': 'Field is required',
      'es': '',
    },
    'hw8fdr6s': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'gz66q4ds': {
      'en': 'Field is required',
      'es': '',
    },
    'o5q2ooq8': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    '5yz0ixa6': {
      'en': 'Conference',
      'es': '',
    },
    'w1shffn1': {
      'en': 'Festival',
      'es': '',
    },
    'elfy2yfi': {
      'en': 'Meeting',
      'es': '',
    },
    'qwp6z1zk': {
      'en': 'Expo',
      'es': '',
    },
    'zqelhz0h': {
      'en': 'Sports',
      'es': '',
    },
    'gl5sylu2': {
      'en': 'Workshop',
      'es': '',
    },
    'cyw3ymga': {
      'en': 'Academic',
      'es': '',
    },
    'uwqypozo': {
      'en': 'Business',
      'es': '',
    },
    '61dxvlsd': {
      'en': 'Professional',
      'es': '',
    },
    'lmj5z9fa': {
      'en': 'Food',
      'es': '',
    },
    'd3u08ub2': {
      'en': 'Art',
      'es': '',
    },
    'l10fytbw': {
      'en': 'Culture',
      'es': '',
    },
    'etcpdvds': {
      'en': 'Social',
      'es': '',
    },
    'qo4sz0sv': {
      'en': 'Network',
      'es': '',
    },
    'ij7jduh7': {
      'en': 'Show',
      'es': '',
    },
    'e8l006in': {
      'en': 'Job',
      'es': '',
    },
    'q1vio3o6': {
      'en': 'Product',
      'es': '',
    },
    'r6ke4syq': {
      'en': 'Expo',
      'es': '',
    },
    '8pq7dxn0': {
      'en': 'Sports',
      'es': '',
    },
    'sm401yzm': {
      'en': 'Workshop',
      'es': '',
    },
    'sai3rryr': {
      'en': 'Crafts',
      'es': '',
    },
    'syo37omo': {
      'en': 'Arts',
      'es': '',
    },
    'vx12j7vg': {
      'en': 'Professional',
      'es': '',
    },
    '5ye32ei9': {
      'en': 'Wellness',
      'es': '',
    },
    'bq7mywx1': {
      'en': 'Tech',
      'es': '',
    },
    'd3ftucu4': {
      'en': 'Add Title here...',
      'es': '',
    },
    'tzcs6i18': {
      'en': 'Ref.',
      'es': '',
    },
    'k16kj41i': {
      'en': 'Ticket type...',
      'es': '',
    },
    'jgm15e2a': {
      'en': 'Search...',
      'es': '',
    },
    'zirpqzhy': {
      'en': 'Request',
      'es': '',
    },
    'yberehbf': {
      'en': 'Open',
      'es': '',
    },
    'tb7w0dfy': {
      'en': 'Place',
      'es': '',
    },
    '8ovy98d7': {
      'en': 'Add description here...',
      'es': '',
    },
    'vo01x1c3': {
      'en': 'Submit',
      'es': '',
    },
    'crv9v00s': {
      'en': 'Field is required',
      'es': '',
    },
    '0qvxfwnn': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'w6yea6v3': {
      'en': 'Field is required',
      'es': '',
    },
    'xndt0xn6': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'xdx6g8k1': {
      'en': 'Field is required',
      'es': '',
    },
    'hzcwg0a4': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    '8u9ahb6r': {
      'en': '',
      'es': '',
    },
  },
  // Feedpage
  {
    '18f80jiq': {
      'en': 'Post',
      'es': '',
    },
    'ighi3yyj': {
      'en': 'Item',
      'es': '',
    },
    'qcfd3902': {
      'en': 'Users',
      'es': '',
    },
    'qll7j0w6': {
      'en': 'Image',
      'es': '',
    },
    'nf7vus5v': {
      'en': 'Video',
      'es': '',
    },
    'ixs63p7v': {
      'en': 'Audio',
      'es': '',
    },
    'c6gh8s5a': {
      'en': 'Text',
      'es': '',
    },
    'tw138i86': {
      'en': 'Product',
      'es': '',
    },
    'hr2wsnpt': {
      'en': 'Service',
      'es': '',
    },
    'qky01nka': {
      'en': 'Event',
      'es': '',
    },
    '01osuru4': {
      'en': 'Place',
      'es': '',
    },
    'yjjqv62d': {
      'en': 'People',
      'es': '',
    },
    'ieh3d638': {
      'en': '',
      'es': '',
    },
    'cnzvxg2r': {
      'en': 'TextField',
      'es': '',
    },
    'go0mkojt': {
      'en': '',
      'es': '',
    },
  },
  // NewCatalogue
  {
    'tt5rlnz3': {
      'en': 'New Catalogue',
      'es': '',
    },
    '3k85oz1k': {
      'en': 'Search terms',
      'es': '',
    },
    'kvbue6ir': {
      'en': 'Select items to include',
      'es': '',
    },
    'xldej51d': {
      'en': '2',
      'es': '',
    },
    'nbwgtnjm': {
      'en': ' /6',
      'es': '',
    },
    'n6dmdj2n': {
      'en': 'Add ID words',
      'es': '',
    },
    's6387byi': {
      'en': 'Catalogue title...',
      'es': '',
    },
    '5m3j3ego': {
      'en': 'Hello World',
      'es': '',
    },
    'bi5xuenz': {
      'en': 'Add catalogue',
      'es': '',
    },
    'zjy5ctrm': {
      'en': 'Home',
      'es': '',
    },
  },
  // MethodWallet
  {
    'k3hk08r9': {
      'en': 'New wallet method',
      'es': '',
    },
    'ljtopag5': {
      'en': 'Digital Payment',
      'es': '',
    },
    '3fnekfvg': {
      'en': 'Cash',
      'es': '',
    },
    '2hswau1r': {
      'en': '3rd Party Transfer',
      'es': '',
    },
    'dk5mqfqd': {
      'en': '',
      'es': '',
    },
    'd6holmek': {
      'en': 'Method Name here...',
      'es': '',
    },
    'ngnweltt': {
      'en': 'Crypto',
      'es': '',
    },
    'nd1iwm50': {
      'en': 'Mobile app',
      'es': '',
    },
    'xwi3gpgr': {
      'en': '',
      'es': '',
    },
    '4j4fkjzj': {
      'en': 'Local Currency',
      'es': '',
    },
    'sbbbcqfi': {
      'en': 'USD',
      'es': '',
    },
    '3tzksghr': {
      'en': '',
      'es': '',
    },
    'crx0e3vw': {
      'en': 'Bank Transfer',
      'es': '',
    },
    '6om6sfou': {
      'en': '',
      'es': '',
    },
    'yzksqmdb': {
      'en': 'I.D. number',
      'es': '',
    },
    'dgi81017': {
      'en': 'Terms & Conditions',
      'es': '',
    },
    'unrmnvh1': {
      'en': '',
      'es': '',
    },
    'thwkci8f': {
      'en': 'Digital Payment',
      'es': '',
    },
    'f1o6ktxi': {
      'en': 'Cash',
      'es': '',
    },
    'fu1ts0bw': {
      'en': '',
      'es': '',
    },
    'duyiik77': {
      'en': 'Mar 25 • 3 hr, 32 min',
      'es': '',
    },
    'ndp3vdpu': {
      'en': 'Add to methods',
      'es': '',
    },
    'ojt5j8kx': {
      'en': 'Method Name here... is required',
      'es': '',
    },
    'u6dle9q8': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'g1qqg58u': {
      'en': 'I.D. number is required',
      'es': '',
    },
    '8qwtjvtl': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    '05gv2pw2': {
      'en': 'Account number is required',
      'es': '',
    },
    'qrpm42wq': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'rs69303z': {
      'en': 'Home',
      'es': '',
    },
  },
  // Publicpage
  {
    'd9hwpxrz': {
      'en': ' ',
      'es': '',
    },
    'am409q5i': {
      'en': 'Verified',
      'es': '',
    },
    'ekd3suu5': {
      'en': 'Activity',
      'es': '',
    },
    'rrvc7b9p': {
      'en': 'Activities & Interests :',
      'es': '',
    },
    'xuootcwx': {
      'en': 'Option 1',
      'es': '',
    },
    'hhv6w6n6': {
      'en': 'Option 2',
      'es': '',
    },
    'jdbrstmi': {
      'en': 'Option 3',
      'es': '',
    },
    'y5d2jfe1': {
      'en': 'Option 1',
      'es': '',
    },
    '76fapnhl': {
      'en': 'Option 2',
      'es': '',
    },
    'zin84te0': {
      'en': 'Option 3',
      'es': '',
    },
    'ieiybayx': {
      'en': 'Places :',
      'es': '',
    },
    '1ddw779v': {
      'en': 'Option 1',
      'es': '',
    },
    'r427bm1x': {
      'en': 'Option 2',
      'es': '',
    },
    'aasaau02': {
      'en': 'Option 3',
      'es': '',
    },
    'p1qie27l': {
      'en': 'USER PIN : ',
      'es': '',
    },
    'amnvvjou': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'ewrxuzmi': {
      'en': 'OBJECT PIN : ',
      'es': '',
    },
    'fbyjwip0': {
      'en': '212',
      'es': '',
    },
    '9sthnl63': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'j6onlmua': {
      'en': 'VOTE RATE : ',
      'es': '',
    },
    'avby2qrn': {
      'en': '213',
      'es': '',
    },
    'xmtajz41': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    '6u4ml2fe': {
      'en': 'ORDER RATE : ',
      'es': '',
    },
    'gz1x8wb4': {
      'en': '213',
      'es': '',
    },
    'sckon5rj': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    '8jz7qus4': {
      'en': ' # ORDERS : ',
      'es': '',
    },
    'yl1fcbft': {
      'en': '212',
      'es': '',
    },
    'xpke3ic9': {
      'en': 'COMPLETED ORDERS',
      'es': '',
    },
    'f3fd1gxk': {
      'en': 'Posts',
      'es': '',
    },
    'wk8go3z8': {
      'en': 'Image',
      'es': '',
    },
    '3xfhvvjj': {
      'en': 'Video',
      'es': '',
    },
    'hy0alzf3': {
      'en': 'Audio',
      'es': '',
    },
    '3h6upgka': {
      'en': 'Text',
      'es': '',
    },
    'm6d4xas2': {
      'en': '',
      'es': '',
    },
    'i32agl68': {
      'en': 'Image',
      'es': '',
    },
    '0pttijto': {
      'en': 'Items',
      'es': '',
    },
    'v1k7vb29': {
      'en': 'All',
      'es': '',
    },
    'ww6dzqw4': {
      'en': 'Popular',
      'es': '',
    },
    'jot3aa1n': {
      'en': 'Catalogues',
      'es': '',
    },
    'wok3m9lz': {
      'en': 'All',
      'es': '',
    },
    'a7rmp25t': {
      'en': '# :',
      'es': '',
    },
    '7nv9h663': {
      'en': 'unpin',
      'es': '',
    },
    'c7daje6u': {
      'en': 'pin',
      'es': '',
    },
    'orfoq8sw': {
      'en': 'unpin',
      'es': '',
    },
    '37ex09zx': {
      'en': 'pin',
      'es': '',
    },
    'qn31det6': {
      'en': 'Your bag',
      'es': '',
    },
    '6hukamju': {
      'en': 'T.',
      'es': '',
    },
    'w2nk774u': {
      'en': 'Order',
      'es': '',
    },
    'n04chxo4': {
      'en': 'Order description',
      'es': '',
    },
    'q5u81hxr': {
      'en': 'T.',
      'es': '',
    },
    'b42n7wje': {
      'en': 'Total :',
      'es': '',
    },
    'fqagdjt7': {
      'en': 'T.',
      'es': '',
    },
    'i4u5yua3': {
      'en': 'Comment',
      'es': '',
    },
    'rbv2478s': {
      'en': 'Select order method :',
      'es': '',
    },
    'oc4jddhx': {
      'en': 'Option 1',
      'es': '',
    },
    'ldnaya7j': {
      'en': 'Option 2',
      'es': '',
    },
    '6rekko00': {
      'en': 'Option 3',
      'es': '',
    },
    'zxia3f0n': {
      'en': 'Select wallet method :',
      'es': '',
    },
    'i07t352l': {
      'en': 'Option 1',
      'es': '',
    },
    'ff5iiqvd': {
      'en': 'Option 2',
      'es': '',
    },
    '7bqru2i2': {
      'en': 'Option 3',
      'es': '',
    },
    'vd8j3btf': {
      'en': 'Cancel',
      'es': '',
    },
    'ofwx8axq': {
      'en': 'Generate',
      'es': '',
    },
    'ysoumgzz': {
      'en': '',
      'es': '',
    },
  },
  // PinnedUsers
  {
    'pg44gxt0': {
      'en': 'Users',
      'es': '',
    },
    '8b8ed1cb': {
      'en': 'Pinned',
      'es': '',
    },
    'nuvyfpvo': {
      'en': 'Messages',
      'es': '',
    },
    '7n0mr3w0': {
      'en': 'Search pinned users',
      'es': '',
    },
    '6tbzstp4': {
      'en': 'status',
      'es': '',
    },
    'gju0ezw7': {
      'en': 'status',
      'es': '',
    },
    'pwwoj152': {
      'en': '1',
      'es': '',
    },
    'bdqb9j7p': {
      'en': 'Home',
      'es': '',
    },
  },
  // Orderpage
  {
    'ptv4zo36': {
      'en': 'Order ',
      'es': '',
    },
    'fq3mtffe': {
      'en': '',
      'es': '',
    },
    '1gjup55i': {
      'en': 'Order',
      'es': '',
    },
    '0ed9ruqm': {
      'en': 'I.D. # ',
      'es': '',
    },
    '67bpyave': {
      'en': 'All Activity from this past month.',
      'es': '',
    },
    't6kzh6w6': {
      'en': 'User: ',
      'es': '',
    },
    'qd6c4f9n': {
      'en': 'All Activity from this past month.',
      'es': '',
    },
    'dul8o5xo': {
      'en': 'User: ',
      'es': '',
    },
    'cy024eta': {
      'en': 'All Activity from this past month.',
      'es': '',
    },
    'fdajd74w': {
      'en': 'Status :',
      'es': '',
    },
    'b3d4lly5': {
      'en': '',
      'es': '',
    },
    'xnqgu47s': {
      'en': 'Pending',
      'es': '',
    },
    't7e0b90a': {
      'en': '4m ago',
      'es': '',
    },
    'xn6znbnn': {
      'en': 'Items :',
      'es': '',
    },
    'vw4ex25c': {
      'en': 'Total :',
      'es': '',
    },
    'sspduwcy': {
      'en': 'T.',
      'es': '',
    },
    'l4of6i0l': {
      'en': 'Proceed with transaction',
      'es': '',
    },
    'l4s92nyd': {
      'en': 'Transaction',
      'es': '',
    },
    'i2aoet1k': {
      'en': '4m ago',
      'es': '',
    },
    'jars4iia': {
      'en': 'Wallet method',
      'es': '',
    },
    'nyizwla2': {
      'en': 'Pago Movil',
      'es': '',
    },
    'e0jdd342': {
      'en': '3rd Party Method',
      'es': '',
    },
    'kzjv21jr': {
      'en': '# I.D. : ',
      'es': '',
    },
    'ynev88b4': {
      'en': 'Hello World ',
      'es': '',
    },
    'xgs2somp': {
      'en': 'Quanity: 1',
      'es': '',
    },
    'xdselok0': {
      'en': '# Account : ',
      'es': '',
    },
    'ghcepavg': {
      'en': 'Hello World ',
      'es': '',
    },
    'cylhat4h': {
      'en': 'Quanity: 1',
      'es': '',
    },
    '4trre1jp': {
      'en': 'Upload',
      'es': '',
    },
    'ygzwaowi': {
      'en': 'Remove',
      'es': '',
    },
    'nac7y66s': {
      'en': 'Notify reciever',
      'es': '',
    },
    '66iqmk0j': {
      'en': 'Confirmation',
      'es': '',
    },
    'qr6squq5': {
      'en': '4m ago',
      'es': '',
    },
    'kv67kckj': {
      'en': 'Management',
      'es': '',
    },
    'whqjep9l': {
      'en': 'Order Processed',
      'es': '',
    },
    'lz2cofye': {
      'en': 'Order boxed',
      'es': '',
    },
    '18xz8tdh': {
      'en': 'Item (s) ready',
      'es': '',
    },
    '2h507tgf': {
      'en': 'Cancel',
      'es': '',
    },
    'yt5930k0': {
      'en': 'OR',
      'es': '',
    },
    'ywt77d5k': {
      'en': 'I.D. number',
      'es': '',
    },
    '3zn0x8tm': {
      'en': 'Button',
      'es': '',
    },
    '9goib2m5': {
      'en': 'Approval',
      'es': '',
    },
    '812citbp': {
      'en': '4m ago',
      'es': '',
    },
    'ar0ewv83': {
      'en': '',
      'es': '',
    },
    'iqpsdv4c': {
      'en': '225',
      'es': '',
    },
    'olfah8sj': {
      'en': 'T.',
      'es': '',
    },
    'ltidf651': {
      'en': 'Comment',
      'es': '',
    },
    'twwgb871': {
      'en': 'Methods :',
      'es': '',
    },
    '1uki59l3': {
      'en': 'Order method',
      'es': '',
    },
    'wxbswwu0': {
      'en': 'Wallet method',
      'es': '',
    },
    'tnzl7v71': {
      'en': 'Mark as relevant',
      'es': '',
    },
    'k3rzaoi9': {
      'en': '',
      'es': '',
    },
    'moz5ntyh': {
      'en': 'Decline',
      'es': '',
    },
    'ow5163kb': {
      'en': 'Accept',
      'es': '',
    },
    '7qkhdi5d': {
      'en': 'Transaction',
      'es': '',
    },
    '2hslcg8a': {
      'en': '4m ago',
      'es': '',
    },
    'eutmuzx9': {
      'en': 'Wallet method',
      'es': '',
    },
    'al16wiwg': {
      'en': 'Pago Movil',
      'es': '',
    },
    'z5qbtuxr': {
      'en': '3rd Party Method',
      'es': '',
    },
    'dkc63i06': {
      'en': '# I.D. : ',
      'es': '',
    },
    'qai4coun': {
      'en': 'Hello World ',
      'es': '',
    },
    'ftd0tyqj': {
      'en': 'Quanity: 1',
      'es': '',
    },
    'nu3esj9t': {
      'en': '# Account : ',
      'es': '',
    },
    'cxd0n0hc': {
      'en': 'Hello World ',
      'es': '',
    },
    's37x1884': {
      'en': 'Quanity: 1',
      'es': '',
    },
    '8fwwgdsa': {
      'en': 'Appeal',
      'es': '',
    },
    '78j7hal4': {
      'en': 'Confirm',
      'es': '',
    },
    'ctf9gphd': {
      'en': 'Confirmation',
      'es': '',
    },
    'qro12c0n': {
      'en': '4m ago',
      'es': '',
    },
    'e5i5ddrc': {
      'en': 'Management',
      'es': '',
    },
    'wts8nbpd': {
      'en': 'Processed',
      'es': '',
    },
    'o3akzmsh': {
      'en': 'Boxed',
      'es': '',
    },
    '386oqqk7': {
      'en': 'Item(s) ready',
      'es': '',
    },
    'clj0r7r5': {
      'en': 'Review',
      'es': '',
    },
    'w9xgg8ca': {
      'en': 'Home',
      'es': '',
    },
  },
  // MethodOrder
  {
    'h2cbvgye': {
      'en': 'New order method',
      'es': '',
    },
    'laks95yc': {
      'en': 'Local',
      'es': '',
    },
    'b4hg2os1': {
      'en': 'Mobile',
      'es': '',
    },
    'wmu2skpg': {
      'en': '',
      'es': '',
    },
    'x12sskfm': {
      'en': 'Pick up',
      'es': '',
    },
    'ukuuo97r': {
      'en': 'Indoor',
      'es': '',
    },
    'n5qsr3ns': {
      'en': '',
      'es': '',
    },
    '1jisxnqm': {
      'en': 'Delivery',
      'es': '',
    },
    'b1r1l4sp': {
      'en': 'Area',
      'es': '',
    },
    '450ur5ru': {
      'en': '',
      'es': '',
    },
    'smfpapcw': {
      'en': 'Method Name here...',
      'es': '',
    },
    'k2fm7n0o': {
      'en': 'Select a place :',
      'es': '',
    },
    'nrsoolrp': {
      'en': 'Add place details...',
      'es': '',
    },
    '9rahx5w8': {
      'en': 'Terms & Conditions',
      'es': '',
    },
    'qymn3ia0': {
      'en': '',
      'es': '',
    },
    'kztfr6tv': {
      'en': 'Option 1',
      'es': '',
    },
    'wfef2681': {
      'en': 'Option 2',
      'es': '',
    },
    '28xovtww': {
      'en': 'Mar 25 • 3 hr, 32 min',
      'es': '',
    },
    'zunbrd0o': {
      'en': 'Add to Order Methods',
      'es': '',
    },
    'no7vmenw': {
      'en': 'Method Name here... is required',
      'es': '',
    },
    '5i75ge9y': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'y0ig19zl': {
      'en': 'I.D. number is required',
      'es': '',
    },
    'e93kvayd': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'w72mz01x': {
      'en': 'Account number is required',
      'es': '',
    },
    '1a1a4jnb': {
      'en': 'Please choose an option from the dropdown',
      'es': '',
    },
    'gimr567a': {
      'en': 'Home',
      'es': '',
    },
  },
  // About
  {
    'f3zu2w4w': {
      'en': 'Home',
      'es': '',
    },
  },
  // System
  {
    'syg54hv6': {
      'en': 'Push Notifications',
      'es': '',
    },
    'kjhboitw': {
      'en':
          'Receive Push notifications from our application on a semi regular basis.',
      'es': '',
    },
    '6db3hml9': {
      'en': 'Object Notifications',
      'es': '',
    },
    'v16cybrp': {
      'en':
          'Receive Push notifications from our application on a semi regular basis.',
      'es': '',
    },
    'a3see93f': {
      'en': 'Order Notifications',
      'es': '',
    },
    'ow977w7y': {
      'en':
          'Receive email notifications from our marketing team about new features.',
      'es': '',
    },
    'q38cwdo8': {
      'en': 'Home',
      'es': '',
    },
  },
  // PinnedObjects
  {
    'bew8sle5': {
      'en': 'Objects',
      'es': '',
    },
    'c1hc8t4b': {
      'en': 'Search pinned objects',
      'es': '',
    },
    'chvmons2': {
      'en': 'Home',
      'es': '',
    },
  },
  // Hashingpage
  {
    'wdamqwvh': {
      'en': 'Welcome',
      'es': '',
    },
    '2li0qa9j': {
      'en': 'Your recent activity is below.',
      'es': '',
    },
    'whsm7nic': {
      'en': 'GENERATED',
      'es': '',
    },
    'bu1jt1dk': {
      'en': 'HASHED',
      'es': '',
    },
    'doaw52vs': {
      'en': 'Last 30 Days',
      'es': '',
    },
    '2ostvmuv': {
      'en': 'Avg. Grade',
      'es': '',
    },
    'ud6yronp': {
      'en': 'Tasks',
      'es': '',
    },
    '2p1i1td9': {
      'en': 'A summary of outstanding tasks.',
      'es': '',
    },
    'uskvzpmt': {
      'en': 'Task Type',
      'es': '',
    },
    '9yo8urjf': {
      'en':
          'Task Description here this one is really long and it goes over maybe? And goes to two lines.',
      'es': '',
    },
    'fv62yqmh': {
      'en': 'Due',
      'es': '',
    },
    'wfmgs8dp': {
      'en': 'Today, 5:30pm',
      'es': '',
    },
    'qcyv31hu': {
      'en': 'Update',
      'es': '',
    },
    'rsgw95rp': {
      'en': '1',
      'es': '',
    },
    '7tmpiycn': {
      'en': 'Task Type',
      'es': '',
    },
    'wtyn9mt4': {
      'en':
          'Task Description here this one is really long and it goes over maybe? And goes to two lines.',
      'es': '',
    },
    'h3e8ln02': {
      'en': 'Due',
      'es': '',
    },
    '2b8cxijg': {
      'en': 'Today, 5:30pm',
      'es': '',
    },
    'rurfl58x': {
      'en': 'Update',
      'es': '',
    },
    '8fiag2l0': {
      'en': '1',
      'es': '',
    },
    '2zvqhfie': {
      'en': 'Home',
      'es': '',
    },
  },
  // Verification
  {
    '39xfv8gr': {
      'en': 'Main contact tags',
      'es': '',
    },
    'wu0rvc5g': {
      'en': 'Secondary text',
      'es': '',
    },
    'bvxyhfyy': {
      'en': 'Enter main mail account.',
      'es': '',
    },
    '7f2kh3ik': {
      'en': 'Enter phone number.',
      'es': '',
    },
    'yzti4w19': {
      'en': 'Get codes',
      'es': '',
    },
    'acpsgxy7': {
      'en': 'Enter code recieved',
      'es': '',
    },
    'vnbwesbi': {
      'en': 'Verify',
      'es': '',
    },
    'ztdch1jj': {
      'en': 'Contact sync',
      'es': '',
    },
    'vb9kdi56': {
      'en': 'Enable authentication by syncing your contact information.',
      'es': '',
    },
    'itdwxm77': {
      'en': 'I.D. Submission',
      'es': '',
    },
    'esiberkj': {
      'en': 'Document',
      'es': '',
    },
    '7cr0nik2': {
      'en': 'Upload',
      'es': '',
    },
    '81ohf845': {
      'en': 'I.D. number',
      'es': '',
    },
    '8onato6t': {
      'en': 'Submit selfie',
      'es': '',
    },
    'iy3hzx83': {
      'en': 'Verify',
      'es': '',
    },
    'peu1x8vq': {
      'en': 'User verification',
      'es': '',
    },
    'r979cvs0': {
      'en': 'Enable authentication by subitting your phone number.',
      'es': '',
    },
    'c5oin773': {
      'en': 'Location Tags',
      'es': '',
    },
    'zlx06ngl': {
      'en': 'Secondary text',
      'es': '',
    },
    '65x6enbv': {
      'en': 'Baruta',
      'es': '',
    },
    'ngyyhg78': {
      'en': 'Chacao',
      'es': '',
    },
    '2dyrnkl0': {
      'en': 'El Hatillo',
      'es': '',
    },
    'pagq8stg': {
      'en': 'Libertador',
      'es': '',
    },
    'xale90ax': {
      'en': 'Sucre',
      'es': '',
    },
    'y7pi43tk': {
      'en': 'Baruta',
      'es': '',
    },
    '0xozg9hy': {
      'en': 'El Cafetal',
      'es': '',
    },
    '3lpa0rgw': {
      'en': 'Las Minas',
      'es': '',
    },
    'w34hznip': {
      'en': 'Altamira',
      'es': '',
    },
    '292k5g0d': {
      'en': 'Bello Campo',
      'es': '',
    },
    'wlot97hf': {
      'en': 'Campo Alegre',
      'es': '',
    },
    'zyu0ahkw': {
      'en': 'Chacao',
      'es': '',
    },
    'rzw1ufe4': {
      'en': 'Chuao',
      'es': '',
    },
    '6jxinyb1': {
      'en': 'Country Club',
      'es': '',
    },
    'ul8grgmr': {
      'en': 'El Bosque',
      'es': '',
    },
    'lox63mg1': {
      'en': 'El Dorado',
      'es': '',
    },
    '0xmvvp5u': {
      'en': 'Pedregal',
      'es': '',
    },
    'mdn44iag': {
      'en': 'El Retiro',
      'es': '',
    },
    '10ba4g90': {
      'en': 'El Rosal',
      'es': '',
    },
    'cx5ahryd': {
      'en': 'Estado Leal',
      'es': '',
    },
    'ca3u9fou': {
      'en': 'La Castellana',
      'es': '',
    },
    '47lkjfl5': {
      'en': 'La Floresta',
      'es': '',
    },
    '9e889g14': {
      'en': 'Los Palos Grandes',
      'es': '',
    },
    '5smfkaaz': {
      'en': 'San Marino',
      'es': '',
    },
    'jry70nh8': {
      'en': 'Chacaito',
      'es': '',
    },
    'keogftap': {
      'en': 'El Hatillo',
      'es': '',
    },
    'n3id31a1': {
      'en': 'La Union',
      'es': '',
    },
    'yno0igh3': {
      'en': 'La Lagunita',
      'es': '',
    },
    't9ctac3z': {
      'en': 'Altagracia',
      'es': '',
    },
    'yjzs2n0g': {
      'en': 'Antimano',
      'es': '',
    },
    '3zjhn5va': {
      'en': 'Caricuao',
      'es': '',
    },
    '9j8wj5fd': {
      'en': 'Catedral',
      'es': '',
    },
    'nnhv8ol3': {
      'en': 'Coche',
      'es': '',
    },
    'ta2cs9c6': {
      'en': 'El Junquito',
      'es': '',
    },
    'dqxlocri': {
      'en': 'El Paraiso',
      'es': '',
    },
    '08i1u7z2': {
      'en': 'El Recreo',
      'es': '',
    },
    'k9j9lxtk': {
      'en': 'El Valle',
      'es': '',
    },
    '8cnjn9r0': {
      'en': 'La Candelaria',
      'es': '',
    },
    '8xbvc7p9': {
      'en': 'Sucre',
      'es': '',
    },
    'vg5kkj61': {
      'en': 'Set location image',
      'es': '',
    },
    '93x5e4ov': {
      'en': 'Verify',
      'es': '',
    },
    '3paa0qf2': {
      'en': 'Activity Sync',
      'es': '',
    },
    '5pbya5zw': {
      'en': 'Enable activity by fixing location tags.',
      'es': '',
    },
    'kdbjn6hm': {
      'en': 'Home',
      'es': '',
    },
  },
  // userrating
  {
    'z393d99m': {
      'en': 'Hello World',
      'es': '',
    },
    'qum9uttj': {
      'en': 'Submit Rating',
      'es': '',
    },
    '1ne00zdy': {
      'en': 'Leave a comment.',
      'es': '',
    },
    'kmmj28py': {
      'en': 'Terms & Conditions',
      'es': '',
    },
    'j9zwje2e': {
      'en': '',
      'es': '',
    },
    'ymckf6ln': {
      'en': 'Home',
      'es': '',
    },
  },
  // MAIN001Copy
  {
    '76f4f4oi': {
      'en': 'Pending',
      'es': '',
    },
    'xb8ca7x7': {
      'en': 'Active',
      'es': '',
    },
    'hktk0piq': {
      'en': 'Recent',
      'es': '',
    },
    '2ian4min': {
      'en': 'Pending',
      'es': '',
    },
    'gpwy80df': {
      'en': 'Pending',
      'es': '',
    },
    'jkzkkiil': {
      'en': 'Active',
      'es': '',
    },
    '5teu75js': {
      'en': 'Recent',
      'es': '',
    },
    'mq1bw9xp': {
      'en': 'Pending',
      'es': '',
    },
    'ee0knoaz': {
      'en': 'Activity',
      'es': '',
    },
    'cmsj2p5r': {
      'en': 'Objects',
      'es': '',
    },
    'r9oq8b4n': {
      'en': 'Orders',
      'es': '',
    },
    'zmbz2tpu': {
      'en': 'Account',
      'es': '',
    },
    '2z55o97c': {
      'en': 'My pins (200)',
      'es': '',
    },
    'o0li9ic8': {
      'en': 'Users',
      'es': '',
    },
    '51trjnw4': {
      'en': 'Secondary text',
      'es': '',
    },
    'b9odq9am': {
      'en': 'Objects',
      'es': '',
    },
    '7shn2c5y': {
      'en': 'Secondary text',
      'es': '',
    },
    'vlbm0x24': {
      'en': 'Places',
      'es': '',
    },
    'lkpr32go': {
      'en': 'Secondary text',
      'es': '',
    },
    'q9v6g5l3': {
      'en': 'Inbox',
      'es': '',
    },
    'xtmhi5px': {
      'en': 'Task Type',
      'es': '',
    },
    'b5xk1o4t': {
      'en':
          'Task Description here this one is really long and it goes over maybe? And goes to two lines.',
      'es': '',
    },
    'xv9753km': {
      'en': 'Due',
      'es': '',
    },
    'vvj7ohpf': {
      'en': 'Today, 5:30pm',
      'es': '',
    },
    'zw4hki33': {
      'en': 'Update',
      'es': '',
    },
    'qc6w6202': {
      'en': '1',
      'es': '',
    },
    'hd3wzqkv': {
      'en': 'Posts (20)',
      'es': '',
    },
    '3g7rk1a8': {
      'en': 'List Item Title',
      'es': '',
    },
    'j49xo8g6': {
      'en': 'Secondary text',
      'es': '',
    },
    'wfth4mzp': {
      'en': 'List Item Title',
      'es': '',
    },
    'dwqd6q0j': {
      'en': 'Secondary text',
      'es': '',
    },
    'axco3bgg': {
      'en': 'My orders',
      'es': '',
    },
    'g604312b': {
      'en': 'Inbox',
      'es': '',
    },
    'ursx2188': {
      'en': 'Sent',
      'es': '',
    },
    'zn36lrld': {
      'en': 'Inbox',
      'es': '',
    },
    'avm9d42c': {
      'en': 'Active',
      'es': '',
    },
    'aip8uml4': {
      'en': 'Hello World',
      'es': '',
    },
    'db74d10f': {
      'en': 'Completed',
      'es': '',
    },
    'sha2lhg7': {
      'en': 'Hello World',
      'es': '',
    },
    'ywv2qql0': {
      'en': 'Reviewed',
      'es': '',
    },
    '6w5k2187': {
      'en': 'Hello World',
      'es': '',
    },
    'md5mw657': {
      'en': 'Method: ',
      'es': '',
    },
    'mered0wn': {
      'en': 'Items in order: ',
      'es': '',
    },
    'yhxy3n23': {
      'en': 'Method: ',
      'es': '',
    },
    'h980xqcb': {
      'en': 'Total: ',
      'es': '',
    },
    'n2tchepm': {
      'en': 'P2P Order Method(s)',
      'es': '',
    },
    'a8tjc3vk': {
      'en': '# I.D. : ',
      'es': '',
    },
    'fktin1o1': {
      'en': 'Quanity: 1',
      'es': '',
    },
    'q7gol38j': {
      'en': '# Index : ',
      'es': '',
    },
    'y6pr92vz': {
      'en': 'Quanity: 1',
      'es': '',
    },
    'ie7qed01': {
      'en': 'Token balance :',
      'es': '',
    },
    'sp8urq16': {
      'en': '123.456.78',
      'es': '',
    },
    '4chhei72': {
      'en': '  T.',
      'es': '',
    },
    '0c4h6ucb': {
      'en': 'Share',
      'es': '',
    },
    '0x1viwfk': {
      'en': 'Share',
      'es': '',
    },
    'cchrjapx': {
      'en': 'Share',
      'es': '',
    },
    '82m1kkkq': {
      'en': 'Pins',
      'es': '',
    },
    '1t4g4fh3': {
      'en': 'Places',
      'es': '',
    },
    'nh6kavof': {
      'en': 'Wallet',
      'es': '',
    },
    'oki807ae': {
      'en': 'Posts',
      'es': '',
    },
    'xywh0hul': {
      'en': 'Items',
      'es': '',
    },
    'u1cfhros': {
      'en': 'Spaces',
      'es': '',
    },
    '03bd1fc4': {
      'en': 'Option 1',
      'es': '',
    },
    '4hmevfro': {
      'en': 'Option 2',
      'es': '',
    },
    'q1d4qbqf': {
      'en': 'Option 3',
      'es': '',
    },
    '56onwfy7': {
      'en': 'Active',
      'es': '',
    },
    'rke65yux': {
      'en': 'Recent',
      'es': '',
    },
    'tn08jb3f': {
      'en': 'Assistant',
      'es': '',
    },
    '9bjcojki': {
      'en': 'User',
      'es': '',
    },
    'emp79606': {
      'en': 'Node',
      'es': '',
    },
    'azymlve0': {
      'en': 'Wallet',
      'es': '',
    },
    't6uuwacw': {
      'en': 'Pins',
      'es': '',
    },
    't9fvhplo': {
      'en': 'Places',
      'es': '',
    },
    'pztqm7sq': {
      'en': 'Wallet',
      'es': '',
    },
    'unnq6u25': {
      'en': 'Posts',
      'es': '',
    },
    'k3m1w1gf': {
      'en': 'Items',
      'es': '',
    },
    '3gfmyjtw': {
      'en': 'Spaces',
      'es': '',
    },
    'vk6e307s': {
      'en': 'Option 1',
      'es': '',
    },
    'nwswpscv': {
      'en': 'Option 2',
      'es': '',
    },
    'k7rvspn1': {
      'en': 'Option 3',
      'es': '',
    },
    'shxdkghb': {
      'en': 'Active',
      'es': '',
    },
    'axaqjx8g': {
      'en': 'Recent',
      'es': '',
    },
    'w1kzrc2i': {
      'en': 'Assistant',
      'es': '',
    },
    'fahsaiau': {
      'en': 'Posts',
      'es': '',
    },
    'rrbjphwe': {
      'en': 'Items',
      'es': '',
    },
    'ejnidnct': {
      'en': 'Spaces',
      'es': '',
    },
    '6jkdvyng': {
      'en': 'Pins',
      'es': '',
    },
    'lokhxy2r': {
      'en': 'Places',
      'es': '',
    },
    'kcrslm75': {
      'en': 'Wallet',
      'es': '',
    },
    'o2617bh0': {
      'en': 'Posts',
      'es': '',
    },
    '4o6n5yva': {
      'en': 'Items',
      'es': '',
    },
    'xqhnbrhy': {
      'en': 'Spaces',
      'es': '',
    },
    'yqgnkeui': {
      'en': 'Option 1',
      'es': '',
    },
    '552g899c': {
      'en': 'Option 2',
      'es': '',
    },
    'w9yzcqru': {
      'en': 'Option 3',
      'es': '',
    },
    'pehdbns6': {
      'en': 'Active',
      'es': '',
    },
    'glrlbi1k': {
      'en': 'Recent',
      'es': '',
    },
    '4r6ywi7q': {
      'en': 'Assistant',
      'es': '',
    },
    'qev9608t': {
      'en': 'Manager',
      'es': '',
    },
    'zd5iujll': {
      'en': 'Metrics',
      'es': '',
    },
    '898lrfup': {
      'en': 'Wallet',
      'es': '',
    },
    'z1dt95t1': {
      'en': 'Pins',
      'es': '',
    },
    'bx62c187': {
      'en': 'Places',
      'es': '',
    },
    '7gadlxxc': {
      'en': 'Wallet',
      'es': '',
    },
    'rtcl8upq': {
      'en': 'Posts',
      'es': '',
    },
    'jlvhmxnf': {
      'en': 'Items',
      'es': '',
    },
    'a4xqbmuh': {
      'en': 'Spaces',
      'es': '',
    },
    'p8809rfv': {
      'en': 'Option 1',
      'es': '',
    },
    'jgi6hfz0': {
      'en': 'Option 2',
      'es': '',
    },
    'zmwpaj3r': {
      'en': 'Option 3',
      'es': '',
    },
    'x44tpy55': {
      'en': 'Active',
      'es': '',
    },
    'ezr6u09g': {
      'en': 'Recent',
      'es': '',
    },
    'pmbsa801': {
      'en': 'Assistant',
      'es': '',
    },
    'fn9yn8fv': {
      'en': 'User',
      'es': '',
    },
    'eqycrsvq': {
      'en': 'Node',
      'es': '',
    },
    'nmj6ozy8': {
      'en': 'Wallet',
      'es': '',
    },
    '8bpshx8y': {
      'en': 'Home',
      'es': '',
    },
  },
  // SETTINGS
  {
    '2rvm4npd': {
      'en': 'User & Profile ',
      'es': '',
    },
    '7g01k86v': {
      'en': 'Verification',
      'es': '',
    },
    '2hvtajr0': {
      'en': 'Parameters',
      'es': '',
    },
    'zp38d7fz': {
      'en': 'Exchange',
      'es': '',
    },
    'c7er8alj': {
      'en': 'System',
      'es': '',
    },
    '66pn5jni': {
      'en': 'Referral',
      'es': '',
    },
    'rdry5m77': {
      'en': 'About',
      'es': '',
    },
    '4r3o4zx1': {
      'en': 'Log out',
      'es': '',
    },
  },
  // OBJECT
  {
    '2u4nwfp6': {
      'en': 'Pin',
      'es': '',
    },
    'u13iz43r': {
      'en': 'Share',
      'es': '',
    },
    'ng55bj8r': {
      'en': 'ADD TO BAG',
      'es': '',
    },
  },
  // CATALOGUE
  {
    '67n3uccs': {
      'en': 'Username',
      'es': '',
    },
    'iy48ih5s': {
      'en': '2h',
      'es': '',
    },
    'qafpjf5h': {
      'en': 'Catalogue name',
      'es': '',
    },
    '0aco7je1': {
      'en': 'Option 1',
      'es': '',
    },
    'amx4paqo': {
      'en': 'Option 2',
      'es': '',
    },
    'j5xjl8cj': {
      'en': 'Option 3',
      'es': '',
    },
  },
  // USER_OBJECT
  {
    'q86mo82g': {
      'en': 'Pins',
      'es': '',
    },
    '7wigayxd': {
      'en': 'Shares',
      'es': '',
    },
    'vfpxuj2w': {
      'en': 'Created:',
      'es': '',
    },
    'kpzd6smo': {
      'en': 'Delete object',
      'es': '',
    },
  },
  // walletconnectTEST
  {
    '9d2kvxka': {
      'en': 'Connect to wallet',
      'es': '',
    },
  },
  // Miscellaneous
  {
    '5d30d5o3': {
      'en': '',
      'es': '',
    },
    '8psauo8g': {
      'en': '',
      'es': '',
    },
    '1s8ug4lx': {
      'en': '',
      'es': '',
    },
    'j71wwb69': {
      'en': '',
      'es': '',
    },
    'vim90xdq': {
      'en': '',
      'es': '',
    },
    '09re2ncf': {
      'en': '',
      'es': '',
    },
    '8e4lcbcv': {
      'en': '',
      'es': '',
    },
    's0emip1o': {
      'en': '',
      'es': '',
    },
    '1a2wlila': {
      'en': '',
      'es': '',
    },
    'r224elko': {
      'en': '',
      'es': '',
    },
    'm0lb8svq': {
      'en': '',
      'es': '',
    },
    'clmykqui': {
      'en': '',
      'es': '',
    },
    '9k7n56f6': {
      'en': '',
      'es': '',
    },
    'bwfb0lzw': {
      'en': '',
      'es': '',
    },
    'b5a14eay': {
      'en': '',
      'es': '',
    },
    'ne56j0ld': {
      'en': '',
      'es': '',
    },
    '7ydwm6nl': {
      'en': '',
      'es': '',
    },
    'saip9pl6': {
      'en': '',
      'es': '',
    },
    'ya5atxc3': {
      'en': '',
      'es': '',
    },
    'um892fm0': {
      'en': '',
      'es': '',
    },
    'x8wk10pm': {
      'en': '',
      'es': '',
    },
    '2jqnds4i': {
      'en': '',
      'es': '',
    },
    '4eryt3u7': {
      'en': '',
      'es': '',
    },
    'wgrppqrt': {
      'en': '',
      'es': '',
    },
    '8ettmxo9': {
      'en': '',
      'es': '',
    },
    '39n6tcqz': {
      'en': '',
      'es': '',
    },
    'gfbofhk9': {
      'en': '',
      'es': '',
    },
    'rihqeann': {
      'en': '',
      'es': '',
    },
    '2q6j92g7': {
      'en': 'Connect Wallet',
      'es': 'Conectar Billetera',
    },
    'booking_history_empty': {
      'en': 'Your booking history is empty',
      'es': 'Tu historial de reservas está vacío',
    },
    'booking_history_desc': {
      'en': 'Book items or services to see them here.',
      'es': 'Reserva artículos o servicios para verlos aquí.',
    },
    'matvykmi': {
      'en': 'Book',
      'es': 'Reservas',
    },
    'ynwxm69q': {
      'en': 'Methods',
      'es': 'Métodos',
    },
    'gunbpyhb': {
      'en': 'Account',
      'es': 'Cuenta',
    },
    '5en137ce': {
      'en': 'Transaction history',
      'es': 'Historial de transacciones',
    },
    'vn7r2tpm': {
      'en': 'A list of historical transactions',
      'es': 'Una lista de transacciones históricas',
    },
    '4q7jxj7j': {
      'en': 'Account details',
      'es': 'Detalles de la cuenta',
    },
    '8j7j7j7j': {
      'en': 'Account information',
      'es': 'Información de la cuenta',
    },
    '7803th34': {
      'en': 'Book',
      'es': 'Reservas',
    },
  },
].reduce((a, b) => a..addAll(b));
