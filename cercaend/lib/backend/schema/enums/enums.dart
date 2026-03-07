import 'package:collection/collection.dart';

enum OrderStatuses {
  Created,
  Generated,
  Pending,
  Accepted,
  Document_Uploaded,
  Document_confirmed,
  Key_Swapped,
  Completed,
  Order_Reviewed,
  Paid,
}

enum InteractionTypes {
  Generation,
  Participation,
  Transition,
}

extension FFEnumExtensions<T extends Enum> on T {
  String serialize() => name;
}

extension FFEnumListExtensions<T extends Enum> on Iterable<T> {
  T? deserialize(String? value) =>
      firstWhereOrNull((e) => e.serialize() == value);
}

T? deserializeEnum<T>(String? value) {
  switch (T) {
    case (OrderStatuses):
      return OrderStatuses.values.deserialize(value) as T?;
    case (InteractionTypes):
      return InteractionTypes.values.deserialize(value) as T?;
    default:
      return null;
  }
}
