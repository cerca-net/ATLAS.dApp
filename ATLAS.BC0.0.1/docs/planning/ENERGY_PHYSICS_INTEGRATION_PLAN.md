# ⚡ Energy Physics Integration Plan

## The Problem (Concepts Got Mixed)

We currently have **two separate display systems** for the same content:

| System | Source | Widget | Where | Energy Physics? |
|--------|--------|--------|-------|----------------|
| **Firebase Feed** | `SubmissionRecord` (Firestore) | `ObjectWidget` | Everywhere | ❌ No |
| **Blockchain Feed** | `SocialPost` (Go Backend) | `BlockchainPostWidget` | Public Page only, behind "Blockchain Mode" toggle | ✅ Yes |

**What the design doc says:** Every submission IS a Smart Object. Energy physics applies to ALL objects, not just a special "blockchain mode" view.

## The Solution: Unify Into One System

### Architecture Decision: **Firebase + Blockchain Overlay**

Instead of replacing Firebase with the blockchain for content storage (which would break the existing app), we **bridge them**:

1. **`SubmissionRecord` (Firebase)** remains the source of truth for content (images, text, video, metadata, poster info).
2. **Go Backend (Blockchain)** is the source of truth for **energy state** (tip_balance, influence_score, status, fossilization).
3. **`ObjectWidget`** renders BOTH — displaying the Firebase content WITH the blockchain energy overlay.

### The Bridge: Firebase `reference.id` → Blockchain `objectId`

Every `SubmissionRecord` in Firebase has a unique `reference.id`. When the `ObjectWidget` loads, it uses this ID to query the Go backend for the object's energy state. If the object doesn't exist in the blockchain yet, it gets auto-registered with default energy (100 TCOIN grace period).

---

## Implementation Steps

### Phase 1: Backend — Object Energy Lookup Endpoint

**File:** `ATLAS.BC0.0.1/internal/api/api.go`

Add a new endpoint: `GET /social/object/energy?objectId={firebase_doc_id}`

This endpoint:
- Looks up the object in the `SocialManager` by its Firebase document ID.
- If not found, auto-creates it with default energy (100 TCOIN).
- Returns: `{ tip_balance, influence_score, status, upvotes, downvotes }`

Also add: `POST /social/object/energize` — To send energy (tokens) to any object by its Firebase doc ID.

### Phase 2: Flutter — Energy Service

**New File:** `cercaend/lib/services/blockchain/energy_service.dart`

```dart
class ObjectEnergyState {
  final int tipBalance;
  final double influenceScore;
  final String status; // "active", "fossilized"
  final int upvotes;
  final int downvotes;
}

class EnergyService {
  // Singleton
  
  /// Get energy state for a Firebase submission
  Future<ObjectEnergyState?> getObjectEnergy(String firebaseDocId);
  
  /// Send energy to an object (Energize/Revive)
  Future<bool> energizeObject(String firebaseDocId, String walletAddress, int amount);
}
```

### Phase 3: ObjectWidget — Integrate Energy Bar

**File:** `cercaend/lib/components/object/object_widget.dart`

Add an **Energy Bar** to EVERY ObjectWidget, positioned between the action buttons row and the info section. The bar shows:

```
┌──────────────────────────────────────────────────┐
│  ⚡ 85 TCOIN    ████████████░░░░░░░    🔥 Score: 450  │
│  [Energize ⚡]     [Tip 💰]                            │
└──────────────────────────────────────────────────┘
```

**Implementation details:**
- On `initState()`, call `EnergyService.getObjectEnergy(widget.object!.reference.id)` 
- Store the result in `ObjectModel` state
- Display the energy bar with:
  - Linear progress indicator showing energy level (0-100 scale, capped at 100 for display)
  - Tip balance number
  - Influence score badge
  - "Energize" button that opens a dialog to send TCOIN
  - If `status == "fossilized"` → show a grayscale overlay with a "Revive" button
- The upvote/downvote buttons ALSO trigger blockchain energy transfer (1 TCOIN per vote)

### Phase 4: Remove Blockchain Mode Toggle & BlockchainPostWidget

**File:** `cercaend/lib/mainpages/publicpage/publicpage_widget.dart`

- Remove the "Blockchain Mode" switch
- Remove the `BlockchainPostWidget` rendering section
- The `ObjectWidget` now handles everything

**File:** `cercaend/lib/components/blockchain_post_widget.dart`

- Can be archived/deleted — its functionality is now in `ObjectWidget`

---

## Visual Design: Energy Bar in ObjectWidget

```
┌─────────────────────────────────────────┐
│  [Image/Video/Audio Content]            │  ← Existing
│                                         │
├─────────────────────────────────────────┤
│  ↑ 5  │  ↓ 1  │  📌 Pin │  Share       │  ← Existing action row
├─────────────────────────────────────────┤
│  ⚡ Energy: 85 TCOIN                    │  ← NEW: Energy bar
│  ████████████████░░░░░░  🔥 Score: 450  │
│  [⚡ Energize]   [💰 Tip 10]            │  ← NEW: Energy actions
├─────────────────────────────────────────┤
│  🖼️ @username  •  2h ago               │  ← Existing info row
│  Title / Description                    │  ← Existing
│  Tags: [art] [music]                    │  ← Existing
└─────────────────────────────────────────┘
```

**Fossilized State:**
```
┌─────────────────────────────────────────┐
│  [Grayscale Content - Dimmed]           │
│          🗿 FOSSILIZED                  │
│  This object ran out of energy          │
│  [⚡ Revive - 50 TCOIN]                │
├─────────────────────────────────────────┤
│  Interactions disabled                  │
└─────────────────────────────────────────┘
```

---

## Data Flow

```
User opens feed
  → Firebase returns SubmissionRecord list
  → ObjectWidget renders each one
  → ObjectWidget calls EnergyService.getObjectEnergy(docId)
  → Go Backend returns { tip_balance, influence, status }
  → ObjectWidget renders energy bar
  
User taps "Energize"
  → Dialog asks for amount
  → EnergyService.energizeObject(docId, walletAddr, amount)
  → Go Backend: deducts from user wallet, adds to object energy
  → ObjectWidget updates UI optimistically
  
User upvotes
  → Firebase: arrayUnion (existing behavior, kept)
  → ALSO: EnergyService sends 1 TCOIN to object (blockchain)
  → Both systems updated
```

## Files to Create/Modify

| File | Action | Description |
|------|--------|-------------|
| `api.go` | Modify | Add `/social/object/energy` and `/social/object/energize` endpoints |
| `social.go` | Modify | Add `GetOrCreateObjectEnergy()` and `EnergizeObject()` methods |
| `energy_service.dart` | Create | Flutter service to bridge ObjectWidget ↔ Blockchain energy |
| `object_widget.dart` | Modify | Add energy bar UI, energy state loading, energize/tip actions |
| `object_model.dart` | Modify | Add energy state fields |
| `publicpage_widget.dart` | Modify | Remove Blockchain Mode toggle & BlockchainPostWidget usage |
| `blockchain_post_widget.dart` | Archive | No longer needed |
