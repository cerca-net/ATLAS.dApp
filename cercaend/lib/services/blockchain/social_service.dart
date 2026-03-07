import 'blockchain_service.dart';

class SocialService {
  final BlockchainService _blockchainService;

  static final SocialService _instance = SocialService._internal();

  factory SocialService() {
    return _instance;
  }

  SocialService._internal() : _blockchainService = BlockchainService();

  /// Fetch the main feed
  /// [userId] can be passed to get a specific user's feed, or null for global/mixed
  Future<List<SocialPost>> getFeed({String? userId, int limit = 20}) async {
    try {
      return await _blockchainService.getSocialFeed(
          userId: userId, limit: limit);
    } catch (e) {
      print('SocialService: Error fetching feed: $e');
      return [];
    }
  }

  /// Create a new post
  /// In the future, this will sign a transaction on the client side.
  /// Currently, it sends the data to the node's API to be processed.
  Future<SocialPost?> createPost({
    required String content,
    required String authorId,
    List<String>? mediaUrls,
  }) async {
    try {
      if (content.isEmpty && (mediaUrls == null || mediaUrls.isEmpty)) {
        throw Exception('Post content cannot be empty');
      }

      return await _blockchainService.createPost(authorId, content,
          mediaUrls: mediaUrls);
    } catch (e) {
      print('SocialService: Error creating post: $e');
      return null;
    }
  }

  /// Interact with a post (Like/Unlike)
  /// In the future, this determines "Influence" on-chain.
  Future<bool> likePost(String postId, String userId) async {
    try {
      return await _blockchainService.likePost(postId, userId, type: 'like');
    } catch (e) {
      print('SocialService: Error liking post: $e');
      return false;
    }
  }

  /// Tip a post (Transfer Energy)
  Future<bool> tipPost(String postId, String userId, int amount) async {
    try {
      return await _blockchainService.tipPost(postId, userId, amount);
    } catch (e) {
      print('SocialService: Error tipping post: $e');
      return false;
    }
  }

  // TODO: Add methods for converting Firestore streams to Blockchain streams
  // to support the hybrid migration phase.
}
