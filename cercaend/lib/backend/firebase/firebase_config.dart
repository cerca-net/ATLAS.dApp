import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/foundation.dart';

Future initFirebase() async {
  if (kIsWeb) {
    await Firebase.initializeApp(
        options: const FirebaseOptions(
            apiKey: "AIzaSyB8W1785yBnfMAJZtQcJZ1GnPE97zAmL_Y",
            authDomain: "cercaend-vk8icg.firebaseapp.com",
            projectId: "cercaend-vk8icg",
            storageBucket: "cercaend-vk8icg.appspot.com",
            messagingSenderId: "618504134340",
            appId: "1:618504134340:web:d3e2ff5b7908d248107391"));
  } else {
    await Firebase.initializeApp();
  }
}
