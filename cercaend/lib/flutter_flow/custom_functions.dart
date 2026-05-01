import '/backend/backend.dart';
import '../backend/supabase/supabase_shim.dart';
import '/backend/schema/custom_records.dart';

double sumofBagitemsListRefvalues(List<double>? refvaluefromiteminlist) {
  // sum the refvalue fields of the bagitem list
  if (refvaluefromiteminlist == null || refvaluefromiteminlist.isEmpty) {
    return 0.0;
  }

  double sum = 0.0;
  for (double value in refvaluefromiteminlist) {
    sum += value;
  }

  return sum;
}

double? averageOrderValue(
  List<double>? totalrefvalue,
  DocumentReference? user,
) {
  // create a function to calculate de average total_ref_value of a user's orders
  if (totalrefvalue == null || user == null) {
    return null;
  }

  double total = 0;
  int count = 0;

  for (double value in totalrefvalue) {
    total += value;
    count++;
  }

  if (count == 0) {
    return null;
  }

  return total / count;
}

int? averageRating(List<int>? ratings) {
  if (ratings == null || ratings.isEmpty) {
    return null;
  }
  int sum = 0;
  for (int rating in ratings) {
    sum += rating;
  }
  return (sum / ratings.length).round();
}

double? generatedcredits(
  double? generation,
  double? participation,
  double? transition,
) {
// sum the arguments
  if (generation == null || participation == null || transition == null) {
    return null;
  }

  return generation + participation + transition;
}

List<UserPostsRecord> getPostsFromUser(String userId) {
  return [];
}

List<UserItemsRecord> getItemsFromUser(String userId) {
  return [];
}

Stream<List<MethodRecord>> queryMethodRecord({
  Query Function(Query)? queryBuilder,
  int limit = -1,
  bool singleRecord = false,
}) {
  return Stream.value([]);
}
