import 'package:cloud_firestore/cloud_firestore.dart';

class UserPostsRecord {
  final DocumentReference reference;

  UserPostsRecord(this.reference);

  static UserPostsRecord fromSnapshot(DocumentSnapshot snapshot) {
    return UserPostsRecord(snapshot.reference);
  }
}

class UserItemsRecord {
  final DocumentReference reference;

  UserItemsRecord(this.reference);

  static UserItemsRecord fromSnapshot(DocumentSnapshot snapshot) {
    return UserItemsRecord(snapshot.reference);
  }
}

class MethodRecord {
  final DocumentReference reference;
  final String methodName;
  final double price;
  final DocumentReference? userRef;
  final DateTime? date;
  final OrderStats? methodType;

  MethodRecord({
    required this.reference,
    required this.methodName,
    required this.price,
    this.userRef,
    this.date,
    this.methodType,
  });

  static MethodRecord fromSnapshot(DocumentSnapshot snapshot) {
    final data = snapshot.data() as Map<String, dynamic>;
    return MethodRecord(
      reference: snapshot.reference,
      methodName: data['method_name'],
      price: (data['price'] as num).toDouble(),
      userRef: data['user_ref'],
      date: data['date']?.toDate(),
      methodType: data['method_type'] != null
          ? OrderStats.values.byName(data['method_type'])
          : null,
    );
  }
}

enum OrderStats {
  Created,
  Generated,
  Pending,
  Accepted,
  Document_Uploaded,
  Document_confirmed,
  Key_Swapped,
  Completed,
  Order_Reviewed,
}
