package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Database represents the blockchain database
type Database struct {
	db *sql.DB
}

// Account represents a blockchain account
type Account struct {
	Address      string    `json:"address"`
	Balance      int64     `json:"balance"`
	Nonce        uint64    `json:"nonce"`
	IsValidator  bool      `json:"is_validator"`
	StakedAmount int64     `json:"staked_amount"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Contract represents a smart contract
type Contract struct {
	Address   string    `json:"address"`
	Code      string    `json:"code"`
	Storage   string    `json:"storage"` // JSON encoded
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Proposal represents a governance proposal
type Proposal struct {
	ID           string    `json:"id"`
	Proposer     string    `json:"proposer"`
	Description  string    `json:"description"`
	Actions      string    `json:"actions"`
	State        string    `json:"state"`
	VotesFor     int64     `json:"votes_for"`
	VotesAgainst int64     `json:"votes_against"`
	StartBlock   int64     `json:"start_block"`
	EndBlock     int64     `json:"end_block"`
	Voters       string    `json:"voters"` // JSON encoded map
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Vote represents a governance vote
type Vote struct {
	ID         int64     `json:"id"`
	ProposalID string    `json:"proposal_id"`
	Voter      string    `json:"voter"`
	Choice     string    `json:"choice"`
	Weight     int64     `json:"weight"`
	CreatedAt  time.Time `json:"created_at"`
}

// OracleData represents oracle data
type OracleData struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Timestamp int64     `json:"timestamp"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

// PostModel represents a post in the database
type PostModel struct {
	ID             string    `json:"id"`
	Author         string    `json:"author"`
	Content        string    `json:"content"`
	MediaURLs      string    `json:"media_urls"` // JSON array
	Hashtags       string    `json:"hashtags"`   // JSON array
	Mentions       string    `json:"mentions"`   // JSON array
	Visibility     string    `json:"visibility"`
	Category       string    `json:"category"`
	TipBalance     uint64    `json:"tip_balance"`
	InfluenceScore float64   `json:"influence_score"`
	Upvotes        uint64    `json:"upvotes"`
	Downvotes      uint64    `json:"downvotes"`
	Likes          int64     `json:"likes"` // legacy
	Comments       int64     `json:"comments"`
	Shares         int64     `json:"shares"`
	Views          int64     `json:"views"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CommentModel represents a comment in the database
type CommentModel struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	ParentID  string    `json:"parent_id"`
	Likes     int64     `json:"likes"`
	Replies   int64     `json:"replies"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// LikeModel represents a like in the database
type LikeModel struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	TargetID   string    `json:"target_id"`
	TargetType string    `json:"target_type"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	database := &Database{db: db}

	// Initialize schema
	if err := database.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %v", err)
	}

	return database, nil
}

// initSchema creates the database tables
func (d *Database) initSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS accounts (
			address TEXT PRIMARY KEY,
			balance INTEGER NOT NULL DEFAULT 0,
			nonce INTEGER NOT NULL DEFAULT 0,
			is_validator BOOLEAN NOT NULL DEFAULT FALSE,
			staked_amount INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS contracts (
			address TEXT PRIMARY KEY,
			code TEXT NOT NULL,
			storage TEXT NOT NULL DEFAULT '{}',
			owner TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS proposals (
			id TEXT PRIMARY KEY,
			proposer TEXT NOT NULL,
			description TEXT NOT NULL,
			actions TEXT NOT NULL,
			state TEXT NOT NULL DEFAULT 'pending',
			votes_for INTEGER NOT NULL DEFAULT 0,
			votes_against INTEGER NOT NULL DEFAULT 0,
			start_block INTEGER NOT NULL,
			end_block INTEGER NOT NULL,
			voters TEXT NOT NULL DEFAULT '{}',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS votes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			proposal_id TEXT NOT NULL,
			voter TEXT NOT NULL,
			choice TEXT NOT NULL,
			weight INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (proposal_id) REFERENCES proposals(id)
		)`,
		`CREATE TABLE IF NOT EXISTS oracle_data (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			timestamp INTEGER NOT NULL,
			source TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS state_snapshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			block_height INTEGER NOT NULL,
			checksum TEXT NOT NULL,
			data TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_accounts_validator ON accounts(is_validator)`,
		`CREATE INDEX IF NOT EXISTS idx_proposals_state ON proposals(state)`,
		`CREATE INDEX IF NOT EXISTS idx_votes_proposal ON votes(proposal_id)`,
		`CREATE INDEX IF NOT EXISTS idx_snapshots_height ON state_snapshots(block_height)`,
		`CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			author TEXT NOT NULL,
			content TEXT NOT NULL,
			media_urls TEXT DEFAULT '[]',
			hashtags TEXT DEFAULT '[]',
			mentions TEXT DEFAULT '[]',
			visibility TEXT DEFAULT 'public',
			category TEXT DEFAULT 'general',
			tip_balance INTEGER DEFAULT 0,
			influence_score REAL DEFAULT 0,
			upvotes INTEGER DEFAULT 0,
			downvotes INTEGER DEFAULT 0,
			likes INTEGER DEFAULT 0,
			comments INTEGER DEFAULT 0,
			shares INTEGER DEFAULT 0,
			views INTEGER DEFAULT 0,
			status TEXT DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id TEXT PRIMARY KEY,
			post_id TEXT NOT NULL,
			author TEXT NOT NULL,
			content TEXT NOT NULL,
			parent_id TEXT DEFAULT '',
			likes INTEGER DEFAULT 0,
			replies INTEGER DEFAULT 0,
			status TEXT DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS likes (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			target_id TEXT NOT NULL,
			target_type TEXT NOT NULL,
			type TEXT DEFAULT 'like',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(author)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id)`,
		`CREATE INDEX IF NOT EXISTS idx_likes_target ON likes(target_id)`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	log.Printf("✅ Database schema initialized successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// Account operations
func (d *Database) GetAccount(address string) (*Account, error) {
	query := `SELECT address, balance, nonce, is_validator, staked_amount, created_at, updated_at
			  FROM accounts WHERE address = ?`

	var account Account
	err := d.db.QueryRow(query, address).Scan(
		&account.Address, &account.Balance, &account.Nonce,
		&account.IsValidator, &account.StakedAmount,
		&account.CreatedAt, &account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %v", err)
	}

	return &account, nil
}

func (d *Database) SetAccount(account *Account) error {
	query := `INSERT OR REPLACE INTO accounts
			  (address, balance, nonce, is_validator, staked_amount, updated_at)
			  VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := d.db.Exec(query, account.Address, account.Balance, account.Nonce,
		account.IsValidator, account.StakedAmount)

	if err != nil {
		return fmt.Errorf("failed to set account: %v", err)
	}

	return nil
}

func (d *Database) GetAllAccounts() ([]*Account, error) {
	query := `SELECT address, balance, nonce, is_validator, staked_amount, created_at, updated_at
			  FROM accounts ORDER BY address`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %v", err)
	}
	defer rows.Close()

	var accounts []*Account
	for rows.Next() {
		var account Account
		err := rows.Scan(&account.Address, &account.Balance, &account.Nonce,
			&account.IsValidator, &account.StakedAmount,
			&account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %v", err)
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

func (d *Database) GetValidators() ([]*Account, error) {
	query := `SELECT address, balance, nonce, is_validator, staked_amount, created_at, updated_at
			  FROM accounts WHERE is_validator = TRUE ORDER BY staked_amount DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query validators: %v", err)
	}
	defer rows.Close()

	var validators []*Account
	for rows.Next() {
		var account Account
		err := rows.Scan(&account.Address, &account.Balance, &account.Nonce,
			&account.IsValidator, &account.StakedAmount,
			&account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan validator: %v", err)
		}
		validators = append(validators, &account)
	}

	return validators, nil
}

// Contract operations
func (d *Database) GetContract(address string) (*Contract, error) {
	query := `SELECT address, code, storage, owner, created_at, updated_at
			  FROM contracts WHERE address = ?`

	var contract Contract
	err := d.db.QueryRow(query, address).Scan(
		&contract.Address, &contract.Code, &contract.Storage,
		&contract.Owner, &contract.CreatedAt, &contract.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %v", err)
	}

	return &contract, nil
}

func (d *Database) SetContract(contract *Contract) error {
	query := `INSERT OR REPLACE INTO contracts
			  (address, code, storage, owner, updated_at)
			  VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := d.db.Exec(query, contract.Address, contract.Code,
		contract.Storage, contract.Owner)

	if err != nil {
		return fmt.Errorf("failed to set contract: %v", err)
	}

	return nil
}

// Proposal operations
func (d *Database) GetProposal(id string) (*Proposal, error) {
	query := `SELECT id, proposer, description, actions, state, votes_for, votes_against,
			  start_block, end_block, voters, created_at, updated_at
			  FROM proposals WHERE id = ?`

	var proposal Proposal
	err := d.db.QueryRow(query, id).Scan(
		&proposal.ID, &proposal.Proposer, &proposal.Description, &proposal.Actions,
		&proposal.State, &proposal.VotesFor, &proposal.VotesAgainst,
		&proposal.StartBlock, &proposal.EndBlock, &proposal.Voters,
		&proposal.CreatedAt, &proposal.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal: %v", err)
	}

	return &proposal, nil
}

func (d *Database) SetProposal(proposal *Proposal) error {
	query := `INSERT OR REPLACE INTO proposals
			  (id, proposer, description, actions, state, votes_for, votes_against,
			   start_block, end_block, voters, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := d.db.Exec(query, proposal.ID, proposal.Proposer, proposal.Description,
		proposal.Actions, proposal.State, proposal.VotesFor, proposal.VotesAgainst,
		proposal.StartBlock, proposal.EndBlock, proposal.Voters)

	if err != nil {
		return fmt.Errorf("failed to set proposal: %v", err)
	}

	return nil
}

func (d *Database) GetAllProposals() ([]*Proposal, error) {
	query := `SELECT id, proposer, description, actions, state, votes_for, votes_against,
			  start_block, end_block, voters, created_at, updated_at
			  FROM proposals ORDER BY created_at DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query proposals: %v", err)
	}
	defer rows.Close()

	var proposals []*Proposal
	for rows.Next() {
		var proposal Proposal
		err := rows.Scan(&proposal.ID, &proposal.Proposer, &proposal.Description,
			&proposal.Actions, &proposal.State, &proposal.VotesFor, &proposal.VotesAgainst,
			&proposal.StartBlock, &proposal.EndBlock, &proposal.Voters,
			&proposal.CreatedAt, &proposal.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan proposal: %v", err)
		}
		proposals = append(proposals, &proposal)
	}

	return proposals, nil
}

// Vote operations
func (d *Database) AddVote(vote *Vote) error {
	query := `INSERT INTO votes (proposal_id, voter, choice, weight) VALUES (?, ?, ?, ?)`

	_, err := d.db.Exec(query, vote.ProposalID, vote.Voter, vote.Choice, vote.Weight)
	if err != nil {
		return fmt.Errorf("failed to add vote: %v", err)
	}

	return nil
}

func (d *Database) GetVotesForProposal(proposalID string) ([]*Vote, error) {
	query := `SELECT id, proposal_id, voter, choice, weight, created_at
			  FROM votes WHERE proposal_id = ? ORDER BY created_at`

	rows, err := d.db.Query(query, proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to query votes: %v", err)
	}
	defer rows.Close()

	var votes []*Vote
	for rows.Next() {
		var vote Vote
		err := rows.Scan(&vote.ID, &vote.ProposalID, &vote.Voter,
			&vote.Choice, &vote.Weight, &vote.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote: %v", err)
		}
		votes = append(votes, &vote)
	}

	return votes, nil
}

// Oracle operations
func (d *Database) SetOracleData(data *OracleData) error {
	query := `INSERT OR REPLACE INTO oracle_data (key, value, timestamp, source)
			  VALUES (?, ?, ?, ?)`

	_, err := d.db.Exec(query, data.Key, data.Value, data.Timestamp, data.Source)
	if err != nil {
		return fmt.Errorf("failed to set oracle data: %v", err)
	}

	return nil
}

func (d *Database) GetOracleData(key string) (*OracleData, error) {
	query := `SELECT key, value, timestamp, source, created_at
			  FROM oracle_data WHERE key = ?`

	var data OracleData
	err := d.db.QueryRow(query, key).Scan(
		&data.Key, &data.Value, &data.Timestamp, &data.Source, &data.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get oracle data: %v", err)
	}

	return &data, nil
}

// State snapshot operations
func (d *Database) SaveSnapshot(blockHeight int64, checksum string, data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot data: %v", err)
	}

	query := `INSERT INTO state_snapshots (block_height, checksum, data) VALUES (?, ?, ?)`

	_, err = d.db.Exec(query, blockHeight, checksum, string(jsonData))
	if err != nil {
		return fmt.Errorf("failed to save snapshot: %v", err)
	}

	return nil
}

func (d *Database) GetLatestSnapshot() (int64, string, map[string]interface{}, error) {
	query := `SELECT block_height, checksum, data FROM state_snapshots
			  ORDER BY block_height DESC LIMIT 1`

	var blockHeight int64
	var checksum, dataStr string

	err := d.db.QueryRow(query).Scan(&blockHeight, &checksum, &dataStr)
	if err == sql.ErrNoRows {
		return 0, "", nil, nil
	}
	if err != nil {
		return 0, "", nil, fmt.Errorf("failed to get latest snapshot: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return 0, "", nil, fmt.Errorf("failed to unmarshal snapshot data: %v", err)
	}

	return blockHeight, checksum, data, nil
}

// Backup and recovery
func (d *Database) Backup(backupPath string) error {
	// TODO: Implement proper SQLite backup using CGO or file copy
	log.Printf("📝 Backup method not yet implemented")
	return fmt.Errorf("backup method not yet implemented")
}

// BackupToWriter creates a backup of the database to a writer (for compression)
func (d *Database) BackupToWriter(writer io.Writer) error {
	// Create a temporary file for the backup
	tempFile, err := os.CreateTemp("", "backup_*.db")
	if err != nil {
		return fmt.Errorf("failed to create temp file for backup: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close() // Close the file so that VACUUM INTO can use it

	// Use SQLite's VACUUM INTO command to create a backup
	_, err = d.db.Exec("VACUUM INTO ?", tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	// Read the backup file and write to the provided writer
	backupFile, err := os.Open(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open backup file: %v", err)
	}
	defer backupFile.Close()

	_, err = io.Copy(writer, backupFile)
	if err != nil {
		return fmt.Errorf("failed to write backup: %v", err)
	}

	return nil
}

// Restore restores the database from a reader
func (d *Database) Restore(reader io.Reader) error {
	// Create a temporary file for the backup
	tempFile, err := os.CreateTemp("", "restore_*.db")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write the backup data to the temp file
	_, err = io.Copy(tempFile, reader)
	if err != nil {
		return fmt.Errorf("failed to write backup to temp file: %v", err)
	}

	// Close the temp file
	tempFile.Close()

	// Close current database connection
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("failed to close current database: %v", err)
	}

	// Remove the current database file
	currentDBPath := d.db.Driver().(interface{ Path() string }).Path()
	if err := os.Remove(currentDBPath); err != nil {
		return fmt.Errorf("failed to remove current database: %v", err)
	}

	// Copy the backup to the current database location
	backupFile, err := os.Open(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open backup file: %v", err)
	}
	defer backupFile.Close()

	currentDBFile, err := os.Create(currentDBPath)
	if err != nil {
		return fmt.Errorf("failed to create new database file: %v", err)
	}
	defer currentDBFile.Close()

	_, err = io.Copy(currentDBFile, backupFile)
	if err != nil {
		return fmt.Errorf("failed to copy backup to database: %v", err)
	}

	// Reopen the database connection
	db, err := sql.Open("sqlite3", currentDBPath)
	if err != nil {
		return fmt.Errorf("failed to reopen database: %v", err)
	}

	d.db = db

	return nil
}

// Migration helper
func (d *Database) MigrateFromJSONSnapshots(snapshotDir string) error {
	log.Printf("🔄 Starting migration from JSON snapshots...")

	// This would read existing JSON snapshots and migrate them to the database
	// For now, we'll just log that this is where migration would happen

	log.Printf("📝 Migration from JSON snapshots would be implemented here")
	log.Printf("📁 Snapshot directory: %s", snapshotDir)

	return nil
}

// Social Media Operations

func (d *Database) CreatePost(post *PostModel) error {
	query := `INSERT INTO posts (
		id, author, content, media_urls, hashtags, mentions, visibility, category,
		tip_balance, influence_score, upvotes, downvotes, likes, comments, shares, views,
		status, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.db.Exec(query,
		post.ID, post.Author, post.Content, post.MediaURLs, post.Hashtags, post.Mentions,
		post.Visibility, post.Category, post.TipBalance, post.InfluenceScore,
		post.Upvotes, post.Downvotes, post.Likes, post.Comments, post.Shares, post.Views,
		post.Status, post.CreatedAt, post.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create post: %v", err)
	}
	return nil
}

func (d *Database) UpdatePost(post *PostModel) error {
	query := `UPDATE posts SET
		content = ?, media_urls = ?, hashtags = ?, mentions = ?, visibility = ?, category = ?,
		tip_balance = ?, influence_score = ?, upvotes = ?, downvotes = ?, likes = ?,
		comments = ?, shares = ?, views = ?, status = ?, updated_at = ?
		WHERE id = ?`

	_, err := d.db.Exec(query,
		post.Content, post.MediaURLs, post.Hashtags, post.Mentions, post.Visibility, post.Category,
		post.TipBalance, post.InfluenceScore, post.Upvotes, post.Downvotes, post.Likes,
		post.Comments, post.Shares, post.Views, post.Status, time.Now(),
		post.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update post: %v", err)
	}
	return nil
}

func (d *Database) GetPost(id string) (*PostModel, error) {
	query := `SELECT id, author, content, media_urls, hashtags, mentions, visibility, category,
		tip_balance, influence_score, upvotes, downvotes, likes, comments, shares, views,
		status, created_at, updated_at FROM posts WHERE id = ?`

	var p PostModel
	err := d.db.QueryRow(query, id).Scan(
		&p.ID, &p.Author, &p.Content, &p.MediaURLs, &p.Hashtags, &p.Mentions,
		&p.Visibility, &p.Category, &p.TipBalance, &p.InfluenceScore,
		&p.Upvotes, &p.Downvotes, &p.Likes, &p.Comments, &p.Shares, &p.Views,
		&p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (d *Database) GetAllPosts(limit int) ([]*PostModel, error) {
	query := `SELECT id, author, content, media_urls, hashtags, mentions, visibility, category,
		tip_balance, influence_score, upvotes, downvotes, likes, comments, shares, views,
		status, created_at, updated_at FROM posts ORDER BY created_at DESC LIMIT ?`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*PostModel
	for rows.Next() {
		var p PostModel
		err := rows.Scan(
			&p.ID, &p.Author, &p.Content, &p.MediaURLs, &p.Hashtags, &p.Mentions,
			&p.Visibility, &p.Category, &p.TipBalance, &p.InfluenceScore,
			&p.Upvotes, &p.Downvotes, &p.Likes, &p.Comments, &p.Shares, &p.Views,
			&p.Status, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	return posts, nil
}

func (d *Database) CreateComment(comment *CommentModel) error {
	query := `INSERT INTO comments (
		id, post_id, author, content, parent_id, likes, replies, status, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.db.Exec(query,
		comment.ID, comment.PostID, comment.Author, comment.Content, comment.ParentID,
		comment.Likes, comment.Replies, comment.Status, comment.CreatedAt, time.Now(),
	)
	return err
}

func (d *Database) GetCommentsByPost(postID string) ([]*CommentModel, error) {
	query := `SELECT id, post_id, author, content, parent_id, likes, replies, status, created_at
		FROM comments WHERE post_id = ? ORDER BY created_at ASC`

	rows, err := d.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*CommentModel
	for rows.Next() {
		var c CommentModel
		err := rows.Scan(
			&c.ID, &c.PostID, &c.Author, &c.Content, &c.ParentID,
			&c.Likes, &c.Replies, &c.Status, &c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}

func (d *Database) CreateLike(like *LikeModel) error {
	query := `INSERT INTO likes (id, user_id, target_id, target_type, type, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := d.db.Exec(query, like.ID, like.UserID, like.TargetID, like.TargetType, like.Type, like.CreatedAt)
	return err
}

func (d *Database) RemoveLike(id string) error {
	_, err := d.db.Exec("DELETE FROM likes WHERE id = ?", id)
	return err
}

func (d *Database) GetLikes(targetID string) ([]*LikeModel, error) {
	query := `SELECT id, user_id, target_id, target_type, type, created_at
		FROM likes WHERE target_id = ?`
	rows, err := d.db.Query(query, targetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var likes []*LikeModel
	for rows.Next() {
		var l LikeModel
		err := rows.Scan(&l.ID, &l.UserID, &l.TargetID, &l.TargetType, &l.Type, &l.CreatedAt)
		if err != nil {
			return nil, err
		}
		likes = append(likes, &l)
	}
	return likes, nil
}
