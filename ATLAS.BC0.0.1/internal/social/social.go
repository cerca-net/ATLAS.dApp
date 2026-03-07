package social

import (
	"atlas-blockchain/internal/identity"
	"atlas-blockchain/pkg/database"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"atlas-blockchain/internal/blockchain"
)

// SocialManager handles all social media operations
type SocialManager struct {
	sync.RWMutex
	posts           map[string]*Post
	comments        map[string]*Comment
	likes           map[string]*Like
	feeds           map[string]*Feed
	moderation      *ContentModerator
	trending        *TrendingManager
	hashtags        map[string]*Hashtag
	mentions        map[string][]string
	reports         map[string]*Report
	identityManager *identity.IdentityManager
	db              *database.Database
	stateManager    *blockchain.StateManager
	logicalClock    uint64 // Global Lamport Clock for Causal Time
}

const (
	INFLUENCE_DECAY_RATE    = 100 // Points per hour decay
	FOSSILIZATION_THRESHOLD = 50  // Energy cost per day
	STATUS_FOSSILIZED       = "fossilized"
	STATUS_ACTIVE           = "active"
)

// Post represents a social media post (Smart Contract Object)
type Post struct {
	ID         string   `json:"id"`
	Author     string   `json:"author"`
	Content    string   `json:"content"` // Mass (Storage)
	MediaURLs  []string `json:"media_urls,omitempty"`
	Hashtags   []string `json:"hashtags,omitempty"`
	Mentions   []string `json:"mentions,omitempty"`
	Visibility string   `json:"visibility"`
	Category   string   `json:"category"`

	// Physics / Smart Contract State
	TipBalance     uint64  `json:"tip_balance"`     // Energy (Tokens)
	InfluenceScore float64 `json:"influence_score"` // Velocity
	Upvotes        uint64  `json:"upvotes"`         // Gravity (+)
	Downvotes      uint64  `json:"downvotes"`       // Gravity (-)

	// Legacy/Compat
	Likes       int64                  `json:"likes"` // Kept for backward compat, mapped to Upvotes
	Comments    int64                  `json:"comments"`
	Shares      int64                  `json:"shares"`
	Views       int64                  `json:"views"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	LogicalTime uint64                 `json:"logical_time"` // Causal Time (Lamport Timestamp)
	Metadata    map[string]interface{} `json:"metadata"`
	Moderation  *ModerationInfo        `json:"moderation,omitempty"`
}

// CalculateInfluenceScore calculates the "Velocity" of the post
func (p *Post) CalculateInfluenceScore() float64 {
	ageInHours := time.Since(p.CreatedAt).Hours()
	if ageInHours < 0 {
		ageInHours = 0
	}

	// Formula: (Upvotes * 10) - (Downvotes * 20) + (TipBalance * 5)
	rawScore := (float64(p.Upvotes) * 10) - (float64(p.Downvotes) * 20) + (float64(p.TipBalance) * 5)

	// Decay
	decay := ageInHours * INFLUENCE_DECAY_RATE

	if decay > rawScore {
		return 0
	}
	return rawScore - decay
}

// Comment represents a comment on a post
type Comment struct {
	ID          string                 `json:"id"`
	PostID      string                 `json:"post_id"`
	Author      string                 `json:"author"`
	Content     string                 `json:"content"`
	ParentID    string                 `json:"parent_id,omitempty"` // For nested comments
	Likes       int64                  `json:"likes"`
	Replies     int64                  `json:"replies"`
	Status      string                 `json:"status"` // "active", "hidden", "deleted"
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	LogicalTime uint64                 `json:"logical_time"` // Causal Time
	Metadata    map[string]interface{} `json:"metadata"`
	Moderation  *ModerationInfo        `json:"moderation,omitempty"`
}

// Like represents a like on a post or comment
type Like struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	TargetID   string    `json:"target_id"`   // Post or comment ID
	TargetType string    `json:"target_type"` // "post" or "comment"
	Type       string    `json:"type"`        // "like", "love", "laugh", "wow", "sad", "angry"
	CreatedAt  time.Time `json:"created_at"`
}

// Feed represents a user's personalized feed
type Feed struct {
	UserID    string    `json:"user_id"`
	Posts     []*Post   `json:"posts"`
	LastSync  time.Time `json:"last_sync"`
	Algorithm string    `json:"algorithm"` // "chronological", "relevance", "trending"
}

// ContentModerator handles content moderation
type ContentModerator struct {
	filters    map[string]*ContentFilter
	blacklist  map[string]bool
	whitelist  map[string]bool
	aiModel    *AIModerator
	moderators []string
}

// ContentFilter represents a content filtering rule
type ContentFilter struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // "keyword", "regex", "ai"
	Patterns  []string  `json:"patterns"`
	Action    string    `json:"action"`   // "flag", "hide", "delete", "warn"
	Severity  string    `json:"severity"` // "low", "medium", "high", "critical"
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AIModerator represents AI-based content moderation
type AIModerator struct {
	ModelVersion string   `json:"model_version"`
	Confidence   float64  `json:"confidence_threshold"`
	Categories   []string `json:"categories"`
}

// ModerationInfo contains moderation details
type ModerationInfo struct {
	Status      string    `json:"status"` // "pending", "approved", "flagged", "hidden"
	Reason      string    `json:"reason,omitempty"`
	Moderator   string    `json:"moderator,omitempty"`
	ReviewedAt  time.Time `json:"reviewed_at,omitempty"`
	AutoFlagged bool      `json:"auto_flagged"`
	Score       float64   `json:"score,omitempty"`
}

// TrendingManager manages trending content
type TrendingManager struct {
	trendingPosts []*TrendingItem
	trendingTags  []*TrendingItem
	algorithm     string
	lastUpdate    time.Time
}

// TrendingItem represents a trending post or hashtag
type TrendingItem struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // "post" or "hashtag"
	Score       float64   `json:"score"`
	Trend       string    `json:"trend"` // "rising", "stable", "falling"
	LastUpdated time.Time `json:"last_updated"`
}

// Hashtag represents a hashtag with usage statistics
type Hashtag struct {
	Tag       string    `json:"tag"`
	Count     int64     `json:"count"`
	Posts     []string  `json:"posts"`
	Trending  bool      `json:"trending"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}

// Report represents a content report
type Report struct {
	ID          string     `json:"id"`
	Reporter    string     `json:"reporter"`
	TargetID    string     `json:"target_id"`
	TargetType  string     `json:"target_type"` // "post" or "comment"
	Reason      string     `json:"reason"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status"`   // "pending", "reviewed", "resolved", "dismissed"
	Priority    string     `json:"priority"` // "low", "medium", "high", "urgent"
	CreatedAt   time.Time  `json:"created_at"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy  string     `json:"reviewed_by,omitempty"`
}

// NewSocialManager creates a new social media manager
func NewSocialManager(identityManager *identity.IdentityManager, db *database.Database, stateManager *blockchain.StateManager) *SocialManager {
	sm := &SocialManager{
		posts:           make(map[string]*Post),
		comments:        make(map[string]*Comment),
		likes:           make(map[string]*Like),
		feeds:           make(map[string]*Feed),
		moderation:      NewContentModerator(),
		trending:        NewTrendingManager(),
		hashtags:        make(map[string]*Hashtag),
		mentions:        make(map[string][]string),
		reports:         make(map[string]*Report),
		identityManager: identityManager,
		db:              db,
		stateManager:    stateManager,
		logicalClock:    0,
	}

	// Load data from DB if available
	if db != nil {
		sm.loadFromDB()
	}

	return sm
}

func (sm *SocialManager) loadFromDB() {
	if sm.db == nil {
		return
	}

	// Load posts
	posts, err := sm.db.GetAllPosts(1000) // Load last 1000 posts
	if err != nil {
		fmt.Printf("Failed to load posts from DB: %v\n", err)
		return
	}

	for _, pModel := range posts {
		// Parse JSON fields
		var mediaURLs, hashtags, mentions []string
		_ = json.Unmarshal([]byte(pModel.MediaURLs), &mediaURLs)
		_ = json.Unmarshal([]byte(pModel.Hashtags), &hashtags)
		_ = json.Unmarshal([]byte(pModel.Mentions), &mentions)

		post := &Post{
			ID:             pModel.ID,
			Author:         pModel.Author,
			Content:        pModel.Content,
			MediaURLs:      mediaURLs,
			Hashtags:       hashtags,
			Mentions:       mentions,
			Visibility:     pModel.Visibility,
			Category:       pModel.Category,
			TipBalance:     pModel.TipBalance,
			InfluenceScore: pModel.InfluenceScore,
			Upvotes:        pModel.Upvotes,
			Downvotes:      pModel.Downvotes,
			Likes:          pModel.Likes,
			Comments:       pModel.Comments,
			Shares:         pModel.Shares,
			Views:          pModel.Views,
			Status:         pModel.Status,
			CreatedAt:      pModel.CreatedAt,
			UpdatedAt:      pModel.UpdatedAt,
			Metadata:       make(map[string]interface{}),
		}
		sm.posts[post.ID] = post

		// Load comments for this post
		if comments, err := sm.db.GetCommentsByPost(post.ID); err == nil {
			for _, cModel := range comments {
				comment := &Comment{
					ID:         cModel.ID,
					PostID:     cModel.PostID,
					Author:     cModel.Author,
					Content:    cModel.Content,
					ParentID:   cModel.ParentID,
					Likes:      cModel.Likes,
					Replies:    cModel.Replies,
					Status:     cModel.Status,
					CreatedAt:  cModel.CreatedAt,
					UpdatedAt:  cModel.CreatedAt,
					Metadata:   make(map[string]interface{}),
					Moderation: &ModerationInfo{Status: "approved"},
				}
				sm.comments[comment.ID] = comment
			}
		}

		// Load likes for this post
		if likes, err := sm.db.GetLikes(post.ID); err == nil {
			for _, lModel := range likes {
				like := &Like{
					ID:         lModel.ID,
					UserID:     lModel.UserID,
					TargetID:   lModel.TargetID,
					TargetType: lModel.TargetType,
					Type:       lModel.Type,
					CreatedAt:  lModel.CreatedAt,
				}
				sm.likes[like.ID] = like
			}
		}

		// Populate indices
		for _, tag := range hashtags {
			sm.updateHashtag(tag, post.ID)
		}
		for _, mention := range mentions {
			sm.updateMentions(mention, post.ID)
		}
	}
	fmt.Printf("✅ Loaded %d posts from database\n", len(posts))
	fmt.Printf("✅ Loaded %d comments from database\n", len(sm.comments))
	fmt.Printf("✅ Loaded %d likes from database\n", len(sm.likes))
}

func (sm *SocialManager) savePostToDB(post *Post) {
	if sm.db == nil {
		return
	}

	mediaJSON, _ := json.Marshal(post.MediaURLs)
	hashtagsJSON, _ := json.Marshal(post.Hashtags)
	mentionsJSON, _ := json.Marshal(post.Mentions)

	model := &database.PostModel{
		ID:             post.ID,
		Author:         post.Author,
		Content:        post.Content,
		MediaURLs:      string(mediaJSON),
		Hashtags:       string(hashtagsJSON),
		Mentions:       string(mentionsJSON),
		Visibility:     post.Visibility,
		Category:       post.Category,
		TipBalance:     post.TipBalance,
		InfluenceScore: post.InfluenceScore,
		Upvotes:        post.Upvotes,
		Downvotes:      post.Downvotes,
		Likes:          post.Likes,
		Comments:       post.Comments,
		Shares:         post.Shares,
		Views:          post.Views,
		Status:         post.Status,
		CreatedAt:      post.CreatedAt,
		UpdatedAt:      post.UpdatedAt,
	}

	if err := sm.db.CreatePost(model); err != nil {
		fmt.Printf("Failed to save post to DB: %v\n", err)
	}
}

func (sm *SocialManager) updatePostInDB(post *Post) {
	if sm.db == nil {
		return
	}

	mediaJSON, _ := json.Marshal(post.MediaURLs)
	hashtagsJSON, _ := json.Marshal(post.Hashtags)
	mentionsJSON, _ := json.Marshal(post.Mentions)

	model := &database.PostModel{
		ID:             post.ID,
		Content:        post.Content,
		MediaURLs:      string(mediaJSON),
		Hashtags:       string(hashtagsJSON),
		Mentions:       string(mentionsJSON),
		Visibility:     post.Visibility,
		Category:       post.Category,
		TipBalance:     post.TipBalance,
		InfluenceScore: post.InfluenceScore,
		Upvotes:        post.Upvotes,
		Downvotes:      post.Downvotes,
		Likes:          post.Likes,
		Comments:       post.Comments,
		Shares:         post.Shares,
		Views:          post.Views,
		Status:         post.Status,
	}

	// Use UpdatePost method (assumed to exist or will be added)
	if err := sm.db.UpdatePost(model); err != nil {
		fmt.Printf("Failed to update post in DB: %v\n", err)
	}
}

// CreatePost creates a new social media post
func (sm *SocialManager) CreatePost(author, content string, mediaURLs []string, visibility, category string) (*Post, error) {
	sm.Lock()
	defer sm.Unlock()

	// Validate content
	if len(strings.TrimSpace(content)) == 0 {
		return nil, fmt.Errorf("post content cannot be empty")
	}

	// Extract hashtags and mentions
	hashtags := extractHashtags(content)
	mentions := extractMentions(content)

	// Moderate content
	moderationInfo := sm.moderation.ModerateContent(content, author)

	postID := generatePostID(author)
	// CAUSAL TIME (Lamport Clock)
	sm.logicalClock++
	logicalTime := sm.logicalClock

	// Physics initialization (New objects have mass and start with default energy)
	initialEnergy := uint64(100)

	post := &Post{
		ID:             postID,
		Author:         author,
		Content:        content,
		MediaURLs:      mediaURLs,
		Hashtags:       hashtags,
		Mentions:       mentions,
		Visibility:     visibility,
		Category:       category,
		Likes:          0,
		Upvotes:        0,             // Gravity (+)
		Downvotes:      0,             // Gravity (-)
		TipBalance:     initialEnergy, // Energy: Starter Pack (Grace Period)
		InfluenceScore: 0,             // Velocity
		Comments:       0,
		Shares:         0,
		Views:          0,
		Status:         STATUS_ACTIVE,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LogicalTime:    logicalTime,
		Metadata:       make(map[string]interface{}),
		Moderation:     moderationInfo,
	}

	// Set status based on moderation
	if moderationInfo.Status == "flagged" {
		post.Status = "hidden"
	}

	// Initial Physics State
	post.InfluenceScore = 100 // Boost for new object

	sm.posts[postID] = post

	// Update hashtags
	for _, tag := range hashtags {
		sm.updateHashtag(tag, postID)
	}

	// Update mentions
	for _, mention := range mentions {
		sm.updateMentions(mention, postID)
	}

	// Update user activity
	if sm.identityManager != nil {
		sm.identityManager.UpdateActivity(author, "post", 1)
	}

	// Save to DB
	sm.savePostToDB(post)

	return post, nil
}

// TipPost adds energy (tokens) to a post
func (sm *SocialManager) TipPost(postID, userID string, amount uint64) (*Post, error) {
	sm.Lock()
	defer sm.Unlock()

	post, exists := sm.posts[postID]
	if !exists {
		return nil, fmt.Errorf("post not found")
	}

	// Revival Logic
	if post.Status == STATUS_FOSSILIZED {
		if amount < FOSSILIZATION_THRESHOLD {
			return nil, fmt.Errorf("insufficient energy to revive post (need %d)", FOSSILIZATION_THRESHOLD)
		}
		post.Status = STATUS_ACTIVE
		// fmt.Printf("Interaction: Post %s revived from fossilization!\n", postID)
	} else if post.Status != STATUS_ACTIVE {
		return nil, fmt.Errorf("cannot tip %s post", post.Status)
	}

	// ENERGY TRANSFER (Core Physics)
	if sm.stateManager != nil {
		balance := sm.stateManager.GetBalance(userID)
		if balance < int64(amount) {
			return nil, fmt.Errorf("insufficient energy (tokens)")
		}
		// Transfer energy: User -> System/Post
		sm.stateManager.SetBalance(userID, balance-int64(amount))
	}

	// Update Energy
	post.TipBalance += amount

	// Update Velocity (Influence)
	post.InfluenceScore = post.CalculateInfluenceScore()
	post.UpdatedAt = time.Now()

	// Identity update
	if sm.identityManager != nil {
		sm.identityManager.UpdateActivity(userID, "tip_given", int64(amount))
		sm.identityManager.UpdateActivity(post.Author, "tip_received", int64(amount))
	}

	// Update DB
	sm.updatePostInDB(post)

	return post, nil
}

// GetPost returns a post by ID
func (sm *SocialManager) GetPost(postID string) (*Post, error) {
	post, exists := sm.posts[postID]
	if !exists {
		return nil, fmt.Errorf("post %s not found", postID)
	}

	// Increment view count
	post.Views++

	// If fossilized, content might be restricted or marked
	if post.Status == STATUS_FOSSILIZED {
		// return nil, fmt.Errorf("post is fossilized") // Strict mode
		// Or separate handling
	}

	post.UpdatedAt = time.Now()

	return post, nil
}

// CheckFossilization iterates through posts and archives those with low energy
func (sm *SocialManager) CheckFossilization() {
	sm.Lock()
	defer sm.Unlock()

	fossilizedCount := 0
	for _, post := range sm.posts {
		if post.Status != STATUS_ACTIVE {
			continue
		}

		// Calculate daily cost (simplified logic)
		// Assume constant burn rate for now

		// If balance drops below threshold, fossilize
		// In a real system, we'd deduct daily rent.
		// Here we check if they have enough "Life Force" to remain active
		if post.TipBalance < FOSSILIZATION_THRESHOLD {
			post.Status = STATUS_FOSSILIZED
			post.InfluenceScore = 0 // Dead objects have no velocity
			sm.updatePostInDB(post)
			fossilizedCount++
		}
	}
	if fossilizedCount > 0 {
		fmt.Printf("🗿 Fossilized %d posts due to low energy.\n", fossilizedCount)
	}
}

// ObjectEnergyState represents the energy physics state of any object (post or item)
// bridged from Firebase via its document ID
type ObjectEnergyState struct {
	ObjectID       string  `json:"object_id"`       // Firebase document ID
	TipBalance     uint64  `json:"tip_balance"`     // Energy (Tokens)
	InfluenceScore float64 `json:"influence_score"` // Velocity
	Status         string  `json:"status"`          // "active", "fossilized"
	Upvotes        uint64  `json:"upvotes"`         // Gravity (+)
	Downvotes      uint64  `json:"downvotes"`       // Gravity (-)
	ObjectType     string  `json:"object_type"`     // "post" or "item"
}

// GetOrCreateObjectEnergy returns the energy state for a Firebase object.
// If the object doesn't exist in the blockchain yet, it auto-registers it
// with default energy (100 TCOIN grace period).
// Both posts and items are treated as "objects" — the only difference is
// that items have a purchase flow, but energy physics applies equally.
func (sm *SocialManager) GetOrCreateObjectEnergy(objectID string, objectType string) (*ObjectEnergyState, error) {
	sm.Lock()
	defer sm.Unlock()

	// Look up existing object by Firebase document ID
	post, exists := sm.posts[objectID]
	if !exists {
		// Auto-register this Firebase object in the blockchain with default energy
		sm.logicalClock++

		post = &Post{
			ID:             objectID, // Use Firebase doc ID directly
			Author:         "",       // Will be set when interactions happen
			Content:        "",       // Content lives in Firebase, not here
			Visibility:     "public",
			Category:       objectType, // "post" or "item"
			TipBalance:     100,        // Grace period energy
			InfluenceScore: 100,        // Initial boost
			Upvotes:        0,
			Downvotes:      0,
			Likes:          0,
			Comments:       0,
			Shares:         0,
			Views:          0,
			Status:         STATUS_ACTIVE,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			LogicalTime:    sm.logicalClock,
			Metadata:       make(map[string]interface{}),
		}
		sm.posts[objectID] = post
		sm.savePostToDB(post)
		fmt.Printf("⚡ Auto-registered Firebase object %s as blockchain object (type: %s)\n", objectID, objectType)
	}

	// Recalculate influence
	post.InfluenceScore = post.CalculateInfluenceScore()

	return &ObjectEnergyState{
		ObjectID:       objectID,
		TipBalance:     post.TipBalance,
		InfluenceScore: post.InfluenceScore,
		Status:         post.Status,
		Upvotes:        post.Upvotes,
		Downvotes:      post.Downvotes,
		ObjectType:     objectType,
	}, nil
}

// EnergizeObject transfers energy (tokens) from a user's wallet to an object.
// Works for both posts and items. Handles revival of fossilized objects.
func (sm *SocialManager) EnergizeObject(objectID string, userID string, amount uint64) (*ObjectEnergyState, error) {
	sm.Lock()
	defer sm.Unlock()

	post, exists := sm.posts[objectID]
	if !exists {
		return nil, fmt.Errorf("object %s not found — fetch energy state first to auto-register", objectID)
	}

	// Revival Logic
	if post.Status == STATUS_FOSSILIZED {
		if amount < FOSSILIZATION_THRESHOLD {
			return nil, fmt.Errorf("insufficient energy to revive object (need %d TCOIN)", FOSSILIZATION_THRESHOLD)
		}
		post.Status = STATUS_ACTIVE
		fmt.Printf("⚡ Object %s revived from fossilization!\n", objectID)
	} else if post.Status != STATUS_ACTIVE {
		return nil, fmt.Errorf("cannot energize %s object", post.Status)
	}

	// ENERGY TRANSFER (Core Physics)
	if sm.stateManager != nil {
		balance := sm.stateManager.GetBalance(userID)
		if balance < int64(amount) {
			return nil, fmt.Errorf("insufficient energy (tokens): need %d, have %d", amount, balance)
		}
		// Transfer energy: User -> Object
		sm.stateManager.SetBalance(userID, balance-int64(amount))
	}

	// Update Energy
	post.TipBalance += amount

	// Update Velocity (Influence)
	post.InfluenceScore = post.CalculateInfluenceScore()
	post.UpdatedAt = time.Now()

	// Identity update
	if sm.identityManager != nil {
		sm.identityManager.UpdateActivity(userID, "energize_given", int64(amount))
		if post.Author != "" {
			sm.identityManager.UpdateActivity(post.Author, "energy_received", int64(amount))
		}
	}

	// Update DB
	sm.updatePostInDB(post)

	objectType := "post"
	if post.Category == "item" {
		objectType = "item"
	}

	return &ObjectEnergyState{
		ObjectID:       objectID,
		TipBalance:     post.TipBalance,
		InfluenceScore: post.InfluenceScore,
		Status:         post.Status,
		Upvotes:        post.Upvotes,
		Downvotes:      post.Downvotes,
		ObjectType:     objectType,
	}, nil
}

// CreateComment creates a comment on a post
func (sm *SocialManager) CreateComment(postID, author, content, parentID string) (*Comment, error) {
	// Validate post exists
	post, exists := sm.posts[postID]
	if !exists {
		return nil, fmt.Errorf("post %s not found", postID)
	}

	if post.Status == STATUS_FOSSILIZED {
		return nil, fmt.Errorf("cannot interact with fossilized post; revive it first")
	}
	if post.Status != STATUS_ACTIVE {
		return nil, fmt.Errorf("cannot comment on %s post", post.Status)
	}

	// ENERGY CHECK & CONSUMPTION
	if sm.stateManager != nil {
		balance := sm.stateManager.GetBalance(author)
		cost := int64(2) // Comments are heavier than likes (mass)

		if balance < cost {
			return nil, fmt.Errorf("insufficient energy: comment costs %d tokens", cost)
		}

		// Deduct from author
		sm.stateManager.SetBalance(author, balance-cost)

		// Feed the Post (Engagement = Life)
		post.TipBalance += uint64(cost)
	}

	// Validate content
	if len(strings.TrimSpace(content)) == 0 {
		return nil, fmt.Errorf("comment content cannot be empty")
	}

	// Moderate content
	moderationInfo := sm.moderation.ModerateContent(content, author)

	// CAUSAL TIME for Comments
	// A comment is causally happened after the post it replies to.
	// Logic: max(Post.LogicalTime, Local.LogicalTime) + 1
	sm.logicalClock++
	if post.LogicalTime >= sm.logicalClock {
		sm.logicalClock = post.LogicalTime + 1
	}
	logicalTime := sm.logicalClock

	commentID := generateCommentID(author, postID)
	comment := &Comment{
		ID:          commentID,
		PostID:      postID,
		Author:      author,
		Content:     content,
		ParentID:    parentID,
		Likes:       0,
		Replies:     0,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		LogicalTime: logicalTime,
		Metadata:    make(map[string]interface{}),
		Moderation:  moderationInfo,
	}

	// Set status based on moderation
	if moderationInfo.Status == "flagged" {
		comment.Status = "hidden"
	}

	sm.comments[commentID] = comment

	// Update post comment count
	post.Comments++
	post.UpdatedAt = time.Now()

	// Update parent comment reply count if nested
	if parentID != "" {
		if parentComment, exists := sm.comments[parentID]; exists {
			parentComment.Replies++
			parentComment.UpdatedAt = time.Now()
		}
	}

	// Update user activity
	if sm.identityManager != nil {
		sm.identityManager.UpdateActivity(author, "comment", 1)
	}

	// Save comment to DB
	if sm.db != nil {
		cModel := &database.CommentModel{
			ID:        comment.ID,
			PostID:    comment.PostID,
			Author:    comment.Author,
			Content:   comment.Content,
			ParentID:  comment.ParentID,
			Likes:     comment.Likes,
			Replies:   comment.Replies,
			Status:    comment.Status,
			CreatedAt: comment.CreatedAt,
		}
		if err := sm.db.CreateComment(cModel); err != nil {
			fmt.Printf("Failed to save comment to DB: %v\n", err)
		}
		sm.updatePostInDB(post) // Update post comment count
	}

	return comment, nil
}

// LikePost likes a post
func (sm *SocialManager) LikePost(postID, userID, likeType string) error {
	sm.Lock()
	defer sm.Unlock()

	// Validate post exists
	post, exists := sm.posts[postID]
	if !exists {
		return fmt.Errorf("post %s not found", postID)
	}

	if post.Status == STATUS_FOSSILIZED {
		return fmt.Errorf("cannot interact with fossilized post; revive it first")
	}
	if post.Status != STATUS_ACTIVE {
		return fmt.Errorf("cannot like %s post", post.Status)
	}

	// ENERGY-BASED VOTING (PHYSICS)
	// Every interaction requires energy (tokens)
	if sm.stateManager != nil {
		balance := sm.stateManager.GetBalance(userID)
		cost := int64(1) // Base cost for interaction

		if balance < cost {
			return fmt.Errorf("insufficient energy (tokens): need %d, have %d", cost, balance)
		}

		// Deduct energy from user (Work performed)
		sm.stateManager.SetBalance(userID, balance-cost)

		// Apply Energy Physics to Object
		isEntropy := likeType == "angry" || likeType == "sad" // Downvote/Attack

		if isEntropy {
			// Entropy: Destructive interference.
			// User pays cost (burned/spent). Post LOSES energy.
			// Total system energy decreases (2 units lost effectively).
			if post.TipBalance >= uint64(cost) {
				post.TipBalance -= uint64(cost)
			} else {
				post.TipBalance = 0
			}
			// Check if this attack killed the post immediately
			if post.TipBalance < FOSSILIZATION_THRESHOLD/10 { // Critical energy failure check could go here
				// For now, we wait for CheckFossilization cycle
			}
		} else {
			// Constructive interference (Upvote/Like).
			// User pays cost. Post GAINS energy.
			// Energy checks out (Transferred).
			post.TipBalance += uint64(cost)

			// If post was implicitly revived (edge case if we allow interaction with fossilized via direct tip, but here we blocked it)
		}
	}

	// Check if already liked
	likeID := generateLikeID(userID, postID)
	if _, exists := sm.likes[likeID]; exists {
		return fmt.Errorf("user already liked this post")
	}

	like := &Like{
		ID:         likeID,
		UserID:     userID,
		TargetID:   postID,
		TargetType: "post",
		Type:       likeType,
		CreatedAt:  time.Now(),
	}

	sm.likes[likeID] = like

	// Update post like count and Physics
	post.Likes++ // Legacy

	// Determine "Gravity" impact
	if likeType == "angry" || likeType == "sad" {
		post.Downvotes++
	} else {
		post.Upvotes++
	}

	post.InfluenceScore = post.CalculateInfluenceScore()
	post.UpdatedAt = time.Now()

	// Update user activity
	if sm.identityManager != nil {
		sm.identityManager.UpdateActivity(userID, "like_given", 1)
		sm.identityManager.UpdateActivity(post.Author, "like_received", 1)
	}

	// Save like to DB
	if sm.db != nil {
		lModel := &database.LikeModel{
			ID:         like.ID,
			UserID:     like.UserID,
			TargetID:   like.TargetID,
			TargetType: like.TargetType,
			Type:       like.Type,
			CreatedAt:  like.CreatedAt,
		}
		if err := sm.db.CreateLike(lModel); err != nil {
			fmt.Printf("Failed to save like to DB: %v\n", err)
		}
		sm.updatePostInDB(post) // Update post like count
	}

	return nil
}

// UnlikePost removes a like from a post
func (sm *SocialManager) UnlikePost(postID, userID string) error {
	sm.Lock()
	defer sm.Unlock()

	likeID := generateLikeID(userID, postID)
	like, exists := sm.likes[likeID]
	if !exists {
		return fmt.Errorf("like not found")
	}

	// Remove like
	delete(sm.likes, likeID)

	// Update post like count and Physics
	if post, exists := sm.posts[postID]; exists {
		post.Likes--
		if post.Likes < 0 {
			post.Likes = 0
		}

		// Revert Gravity impact
		if like.Type == "angry" || like.Type == "sad" {
			if post.Downvotes > 0 {
				post.Downvotes--
			}
		} else {
			if post.Upvotes > 0 {
				post.Upvotes--
			}
		}

		post.InfluenceScore = post.CalculateInfluenceScore()
		post.UpdatedAt = time.Now()

		// Save changes to DB
		if sm.db != nil {
			if err := sm.db.RemoveLike(likeID); err != nil {
				fmt.Printf("Failed to remove like from DB: %v\n", err)
			}
			sm.updatePostInDB(post)
		}
	}

	return nil
}

// GetFeed returns a user's personalized feed
func (sm *SocialManager) GetFeed(userID string, limit int) ([]*Post, error) {
	feed, exists := sm.feeds[userID]
	if !exists {
		// Create new feed
		feed = &Feed{
			UserID:    userID,
			Posts:     make([]*Post, 0),
			LastSync:  time.Now(),
			Algorithm: "relevance",
		}
		sm.feeds[userID] = feed
	}

	// Get posts based on algorithm
	var posts []*Post
	switch feed.Algorithm {
	case "chronological":
		posts = sm.getChronologicalFeed(userID, limit)
	case "relevance":
		posts = sm.getRelevanceFeed(userID, limit)
	case "trending":
		posts = sm.getTrendingFeed(limit)
	case "causal":
		posts = sm.getCausalFeed(userID, limit)
	default:
		posts = sm.getChronologicalFeed(userID, limit)
	}

	// Update feed
	feed.Posts = posts
	feed.LastSync = time.Now()

	return posts, nil
}

// getChronologicalFeed returns posts in chronological order
func (sm *SocialManager) getChronologicalFeed(userID string, limit int) []*Post {
	var posts []*Post
	for _, post := range sm.posts {
		if post.Status == "active" && (post.Visibility == "public" || post.Author == userID) {
			posts = append(posts, post)
		}
	}

	// Sort by creation time (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	// Limit results
	if len(posts) > limit {
		posts = posts[:limit]
	}

	return posts
}

// getCausalFeed returns posts in causal order (Lamport Clock)
func (sm *SocialManager) getCausalFeed(userID string, limit int) []*Post {
	var posts []*Post
	for _, post := range sm.posts {
		// Only include active posts
		if post.Status == STATUS_ACTIVE && (post.Visibility == "public" || post.Author == userID) {
			posts = append(posts, post)
		}
	}

	// Sort by LogicalTime (Causal Order), descending
	sort.Slice(posts, func(i, j int) bool {
		if posts[i].LogicalTime == posts[j].LogicalTime {
			// Fallback to CreatedAt if logical times match (shouldn't happen often)
			return posts[i].CreatedAt.After(posts[j].CreatedAt)
		}
		return posts[i].LogicalTime > posts[j].LogicalTime
	})

	// Limit results
	if len(posts) > limit {
		posts = posts[:limit]
	}

	return posts
}

// getRelevanceFeed returns posts based on relevance to user
func (sm *SocialManager) getRelevanceFeed(userID string, limit int) []*Post {
	// Get user identity for relevance calculation
	var userIdentity *identity.UserIdentity
	if sm.identityManager != nil {
		userIdentity, _ = sm.identityManager.GetIdentity(userID)
	}

	var posts []*Post
	for _, post := range sm.posts {
		if post.Status == "active" && (post.Visibility == "public" || post.Author == userID) {
			// Calculate relevance score
			relevance := sm.calculateRelevance(post, userIdentity)
			post.Metadata["relevance_score"] = relevance
			posts = append(posts, post)
		}
	}

	// Sort by relevance score
	sort.Slice(posts, func(i, j int) bool {
		scoreI := posts[i].Metadata["relevance_score"].(float64)
		scoreJ := posts[j].Metadata["relevance_score"].(float64)
		return scoreI > scoreJ
	})

	// Limit results
	if len(posts) > limit {
		posts = posts[:limit]
	}

	return posts
}

// getTrendingFeed returns trending posts
func (sm *SocialManager) getTrendingFeed(limit int) []*Post {
	var posts []*Post
	for _, post := range sm.posts {
		if post.Status == "active" {
			// Calculate trending score
			trendingScore := sm.calculateTrendingScore(post)
			post.Metadata["trending_score"] = trendingScore
			posts = append(posts, post)
		}
	}

	// Sort by trending score
	sort.Slice(posts, func(i, j int) bool {
		scoreI := posts[i].Metadata["trending_score"].(float64)
		scoreJ := posts[j].Metadata["trending_score"].(float64)
		return scoreI > scoreJ
	})

	// Limit results
	if len(posts) > limit {
		posts = posts[:limit]
	}

	return posts
}

// calculateRelevance calculates relevance score for a post
func (sm *SocialManager) calculateRelevance(post *Post, userIdentity *identity.UserIdentity) float64 {
	score := 0.0

	// Base score from engagement
	score += float64(post.Likes) * 0.1
	score += float64(post.Comments) * 0.2
	score += float64(post.Shares) * 0.3

	// Recency bonus
	hoursSinceCreation := time.Since(post.CreatedAt).Hours()
	if hoursSinceCreation < 24 {
		score += 10.0
	} else if hoursSinceCreation < 168 { // 1 week
		score += 5.0
	}

	// User-specific factors
	if userIdentity != nil {
		// Category preference
		if post.Category == "governance" && userIdentity.Activity.ProposalsCreated > 0 {
			score += 5.0
		}
		if post.Category == "commerce" && userIdentity.Activity.Transactions > 0 {
			score += 3.0
		}

		// Author reputation
		if sm.identityManager != nil {
			if authorIdentity, err := sm.identityManager.GetIdentity(post.Author); err == nil {
				score += authorIdentity.Reputation.Overall * 0.1
			}
		}
	}

	return score
}

// calculateTrendingScore calculates trending score for a post
func (sm *SocialManager) calculateTrendingScore(post *Post) float64 {
	score := 0.0

	// Engagement velocity
	hoursSinceCreation := time.Since(post.CreatedAt).Hours()
	if hoursSinceCreation > 0 {
		engagementRate := float64(post.Likes+post.Comments*2+post.Shares*3) / hoursSinceCreation
		score += engagementRate * 10.0
	}

	// Absolute engagement
	score += float64(post.Likes) * 0.1
	score += float64(post.Comments) * 0.2
	score += float64(post.Shares) * 0.5
	score += float64(post.Views) * 0.01

	// Recency penalty
	if hoursSinceCreation > 168 { // 1 week
		score *= 0.5
	}

	return score
}

// ReportContent reports content for moderation
func (sm *SocialManager) ReportContent(reporter, targetID, targetType, reason, description string) error {
	// Validate target exists
	if targetType == "post" {
		if _, exists := sm.posts[targetID]; !exists {
			return fmt.Errorf("post %s not found", targetID)
		}
	} else if targetType == "comment" {
		if _, exists := sm.comments[targetID]; !exists {
			return fmt.Errorf("comment %s not found", targetID)
		}
	} else {
		return fmt.Errorf("invalid target type: %s", targetType)
	}

	reportID := generateReportID(reporter, targetID)
	report := &Report{
		ID:          reportID,
		Reporter:    reporter,
		TargetID:    targetID,
		TargetType:  targetType,
		Reason:      reason,
		Description: description,
		Status:      "pending",
		Priority:    "medium",
		CreatedAt:   time.Now(),
	}

	sm.reports[reportID] = report

	return nil
}

// GetTrendingHashtags returns trending hashtags
func (sm *SocialManager) GetTrendingHashtags(limit int) []*Hashtag {
	var hashtags []*Hashtag
	for _, hashtag := range sm.hashtags {
		hashtags = append(hashtags, hashtag)
	}

	// Sort by count
	sort.Slice(hashtags, func(i, j int) bool {
		return hashtags[i].Count > hashtags[j].Count
	})

	// Limit results
	if len(hashtags) > limit {
		hashtags = hashtags[:limit]
	}

	return hashtags
}

// SearchPosts searches posts by content, hashtags, or author
func (sm *SocialManager) SearchPosts(query string, limit int) []*Post {
	query = strings.ToLower(query)
	var results []*Post

	for _, post := range sm.posts {
		if post.Status != "active" {
			continue
		}

		// Search in content
		if strings.Contains(strings.ToLower(post.Content), query) {
			results = append(results, post)
			continue
		}

		// Search in hashtags
		for _, hashtag := range post.Hashtags {
			if strings.Contains(strings.ToLower(hashtag), query) {
				results = append(results, post)
				break
			}
		}

		// Search in author
		if strings.Contains(strings.ToLower(post.Author), query) {
			results = append(results, post)
		}
	}

	// Sort by relevance (simplified)
	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.After(results[j].CreatedAt)
	})

	// Limit results
	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

// Helper functions
func extractHashtags(content string) []string {
	var hashtags []string
	words := strings.Fields(content)
	for _, word := range words {
		if strings.HasPrefix(word, "#") && len(word) > 1 {
			hashtag := strings.ToLower(strings.TrimPrefix(word, "#"))
			hashtag = strings.Trim(hashtag, ".,!?;:")
			if hashtag != "" {
				hashtags = append(hashtags, hashtag)
			}
		}
	}
	return hashtags
}

func extractMentions(content string) []string {
	var mentions []string
	words := strings.Fields(content)
	for _, word := range words {
		if strings.HasPrefix(word, "@") && len(word) > 1 {
			mention := strings.TrimPrefix(word, "@")
			mention = strings.Trim(mention, ".,!?;:")
			if mention != "" {
				mentions = append(mentions, mention)
			}
		}
	}
	return mentions
}

func (sm *SocialManager) updateHashtag(tag, postID string) {
	hashtag, exists := sm.hashtags[tag]
	if !exists {
		hashtag = &Hashtag{
			Tag:       tag,
			Count:     0,
			Posts:     make([]string, 0),
			Trending:  false,
			CreatedAt: time.Now(),
		}
		sm.hashtags[tag] = hashtag
	}

	hashtag.Count++
	hashtag.Posts = append(hashtag.Posts, postID)
	hashtag.LastUsed = time.Now()

	// Check if trending
	if hashtag.Count >= 10 {
		hashtag.Trending = true
	}
}

func (sm *SocialManager) updateMentions(mention, postID string) {
	sm.mentions[mention] = append(sm.mentions[mention], postID)
}

// Content moderation functions
func NewContentModerator() *ContentModerator {
	return &ContentModerator{
		filters:    make(map[string]*ContentFilter),
		blacklist:  make(map[string]bool),
		whitelist:  make(map[string]bool),
		aiModel:    &AIModerator{ModelVersion: "1.0", Confidence: 0.8},
		moderators: make([]string, 0),
	}
}

func (cm *ContentModerator) ModerateContent(content, author string) *ModerationInfo {
	// Simple keyword-based moderation for now
	// In production, this would use AI/ML models

	contentLower := strings.ToLower(content)

	// Check for obvious violations
	if strings.Contains(contentLower, "spam") || strings.Contains(contentLower, "scam") {
		return &ModerationInfo{
			Status:      "flagged",
			Reason:      "Potential spam content",
			AutoFlagged: true,
			Score:       0.8,
		}
	}

	// Check for excessive caps
	upperCount := 0
	for _, char := range content {
		if char >= 'A' && char <= 'Z' {
			upperCount++
		}
	}
	if float64(upperCount)/float64(len(content)) > 0.7 {
		return &ModerationInfo{
			Status:      "flagged",
			Reason:      "Excessive capitalization",
			AutoFlagged: true,
			Score:       0.6,
		}
	}

	// Default: approved
	return &ModerationInfo{
		Status:      "approved",
		AutoFlagged: false,
		Score:       0.1,
	}
}

// Trending manager functions
func NewTrendingManager() *TrendingManager {
	return &TrendingManager{
		trendingPosts: make([]*TrendingItem, 0),
		trendingTags:  make([]*TrendingItem, 0),
		algorithm:     "engagement_velocity",
		lastUpdate:    time.Now(),
	}
}

// Helper ID generation functions
func generatePostID(author string) string {
	data := fmt.Sprintf("post_%s_%d", author, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

func generateCommentID(author, postID string) string {
	data := fmt.Sprintf("comment_%s_%s_%d", author, postID, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

func generateLikeID(userID, targetID string) string {
	data := fmt.Sprintf("like_%s_%s", userID, targetID)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

func generateReportID(reporter, targetID string) string {
	data := fmt.Sprintf("report_%s_%s_%d", reporter, targetID, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}
