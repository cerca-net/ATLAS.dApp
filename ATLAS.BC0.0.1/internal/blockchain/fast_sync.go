package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// SnapshotMetadata contains information about a state snapshot
type SnapshotMetadata struct {
	BlockHeight int64     `json:"block_height"`
	BlockHash   string    `json:"block_hash"`
	StateRoot   string    `json:"state_root"`
	Timestamp   time.Time `json:"timestamp"`
	FileSize    int64     `json:"file_size"`
	Checksum    string    `json:"checksum"`
}

// FastSyncManager handles snapshot-based fast synchronization
type FastSyncManager struct {
	stateManager *StateManager
	blockManager *BlockManager
	snapshotDir  string
}

// NewFastSyncManager creates a new fast sync manager
func NewFastSyncManager(stateManager *StateManager, blockManager *BlockManager, dataDir string) *FastSyncManager {
	snapshotDir := filepath.Join(dataDir, "snapshots")
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		log.Printf("⚠️  Failed to create snapshot directory: %v", err)
	}

	return &FastSyncManager{
		stateManager: stateManager,
		blockManager: blockManager,
		snapshotDir:  snapshotDir,
	}
}

// CreateSnapshot creates a snapshot of the current blockchain state
func (fsm *FastSyncManager) CreateSnapshot() (*SnapshotMetadata, error) {
	log.Printf("📸 [FAST-SYNC] Creating state snapshot...")

	latestBlock := fsm.blockManager.GetLatestBlock()
	if latestBlock == nil {
		return nil, fmt.Errorf("no blocks in chain")
	}

	// Create snapshot metadata
	metadata := &SnapshotMetadata{
		BlockHeight: int64(latestBlock.Index),
		BlockHash:   latestBlock.Hash,
		StateRoot:   fsm.stateManager.GetStateChecksum(),
		Timestamp:   time.Now(),
	}

	// Export state to snapshot file
	snapshotPath := filepath.Join(fsm.snapshotDir, fmt.Sprintf("snapshot_%d_%s.json", metadata.BlockHeight, time.Now().Format("20060102_150405")))

	if err := fsm.stateManager.ExportState(snapshotPath); err != nil {
		return nil, fmt.Errorf("failed to export state: %v", err)
	}

	// Get file size
	fileInfo, err := os.Stat(snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat snapshot file: %v", err)
	}
	metadata.FileSize = fileInfo.Size()

	// Calculate checksum
	metadata.Checksum = fsm.stateManager.GetStateChecksum()

	// Save metadata
	metadataPath := snapshotPath + ".meta"
	metadataData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %v", err)
	}

	if err := ioutil.WriteFile(metadataPath, metadataData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write metadata: %v", err)
	}

	log.Printf("✅ [FAST-SYNC] Snapshot created at height %d (size: %d bytes)", metadata.BlockHeight, metadata.FileSize)
	return metadata, nil
}

// LoadSnapshot loads a state snapshot from file
func (fsm *FastSyncManager) LoadSnapshot(snapshotPath string) error {
	log.Printf("📥 [FAST-SYNC] Loading snapshot from %s...", snapshotPath)

	// Load and verify metadata
	metadataPath := snapshotPath + ".meta"
	metadataData, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %v", err)
	}

	var metadata SnapshotMetadata
	if err := json.Unmarshal(metadataData, &metadata); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	// Import state from snapshot
	if err := fsm.stateManager.ImportState(snapshotPath); err != nil {
		return fmt.Errorf("failed to import state: %v", err)
	}

	// Verify checksum after import
	currentChecksum := fsm.stateManager.GetStateChecksum()
	if currentChecksum != metadata.Checksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", metadata.Checksum, currentChecksum)
	}

	log.Printf("✅ [FAST-SYNC] Snapshot loaded successfully (height: %d)", metadata.BlockHeight)
	return nil
}

// GetLatestSnapshot returns the path to the most recent snapshot
func (fsm *FastSyncManager) GetLatestSnapshot() (string, *SnapshotMetadata, error) {
	files, err := ioutil.ReadDir(fsm.snapshotDir)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read snapshot directory: %v", err)
	}

	var latestSnapshot string
	var latestMetadata *SnapshotMetadata
	var latestTime time.Time

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".meta" {
			metadataPath := filepath.Join(fsm.snapshotDir, file.Name())
			snapshotPath := metadataPath[:len(metadataPath)-5] // Remove .meta extension

			metadataData, err := ioutil.ReadFile(metadataPath)
			if err != nil {
				continue
			}

			var metadata SnapshotMetadata
			if err := json.Unmarshal(metadataData, &metadata); err != nil {
				continue
			}

			if metadata.Timestamp.After(latestTime) {
				latestTime = metadata.Timestamp
				latestSnapshot = snapshotPath
				latestMetadata = &metadata
			}
		}
	}

	if latestSnapshot == "" {
		return "", nil, fmt.Errorf("no snapshots found")
	}

	return latestSnapshot, latestMetadata, nil
}

// CleanupOldSnapshots removes old snapshots, keeping only the most recent N
func (fsm *FastSyncManager) CleanupOldSnapshots(keepCount int) error {
	files, err := ioutil.ReadDir(fsm.snapshotDir)
	if err != nil {
		return fmt.Errorf("failed to read snapshot directory: %v", err)
	}

	// Collect all snapshot metadata
	type snapshotInfo struct {
		path     string
		metadata SnapshotMetadata
	}

	snapshots := make([]snapshotInfo, 0)
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".meta" {
			metadataPath := filepath.Join(fsm.snapshotDir, file.Name())
			snapshotPath := metadataPath[:len(metadataPath)-5]

			metadataData, err := ioutil.ReadFile(metadataPath)
			if err != nil {
				continue
			}

			var metadata SnapshotMetadata
			if err := json.Unmarshal(metadataData, &metadata); err != nil {
				continue
			}

			snapshots = append(snapshots, snapshotInfo{
				path:     snapshotPath,
				metadata: metadata,
			})
		}
	}

	// Sort by timestamp (newest first)
	for i := 0; i < len(snapshots); i++ {
		for j := i + 1; j < len(snapshots); j++ {
			if snapshots[j].metadata.Timestamp.After(snapshots[i].metadata.Timestamp) {
				snapshots[i], snapshots[j] = snapshots[j], snapshots[i]
			}
		}
	}

	// Remove old snapshots
	for i := keepCount; i < len(snapshots); i++ {
		if err := os.Remove(snapshots[i].path); err != nil {
			log.Printf("⚠️  Failed to remove snapshot %s: %v", snapshots[i].path, err)
		}
		if err := os.Remove(snapshots[i].path + ".meta"); err != nil {
			log.Printf("⚠️  Failed to remove metadata %s: %v", snapshots[i].path+".meta", err)
		}
		log.Printf("🗑️  Removed old snapshot at height %d", snapshots[i].metadata.BlockHeight)
	}

	return nil
}
