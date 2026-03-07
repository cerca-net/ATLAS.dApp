# 📜 Constitution of Social Physics: The "Post" as a Living Entity

## **1. Core Philosophy: Digital Existence**
In the ATLAS ecosystem, a "Post" is not just a row in a database. It is a **living Smart Contract (Object)** that exists in the blockchain state, subject to the **Laws of Social Physics**.

### **The Three Laws:**
1.  **Energy is Required for Existence:** Data occupies space and requires energy (tokens) to maintain.
2.  **Time is Causal, Not Global:** Events are defined by their relationship to other events (Parent/Child), not by a universal clock.
3.  **Reputation is a Flow:** Influence must be constantly maintained or it decays (Entropy).

---

## **2. The "Object" Model (Matter & Energy)**

### **A. Structure (Matter)**
A Post Object has "Mass" (Storage Size), "State", and "Threads" (Tag arrays that dictate its structural context).
```go
contract Post {
    // Identity
    address public author;
    bytes32 public contentHash; // IPFS/Storage hash (The "Matter")
    string[] public threads;    // Substrate tags for Data Unit propagation
    
    // State
    bool public isFossilized;   // True if Data Units (Energy) drop below threshold
    uint256 public energy;      // The "Life Force" (Data Units [DU] Balance)
}
```

### **B. The Archeological Model (Fossilization)**
*   **Living State:** A post is "Alive" as long as its Data Units (`Energy`) > minimum. It is indexed, searchable, and interactive in the "Hot State" (RAM/Fast DB).
*   **Fossilization (Entropy):**
    *   Unlike static systems, posts lose relevance naturally. When `Energy` drops too low (via decay or downvotes), the post "dies".
    *   It disappears from active algorithmic Feed discovery.
*   **Revival:** A user must spend their own accumulated Data Units (DUs) to revive the post back into the Living State.

---

## **3. The "Interaction" Model (Energy Transfer & Data Units)**

### **A. Data Units (DU) - The Fuel**
*   **Tokens (TCOIN) ≠ Data Units (DU).**
    *   **TCOIN** is the hard-capped (1 Billion) store of value for validators and network security.
    *   **Data Units** are the fluid "cognitive fuel" generated passively by holding TCOIN, or actively by interacting with the environment.
*   **Interactions:** There are no "Likes." There is only **Data Unit Transfer and Generation**.
    *   **Publishing an Object:** Costs a minimum pool of DUs (Anti-spam).
    *   **Upvote:** A user sends Data Units to a Post. This extends the Post's life and validates its "Threads". The system additionally spawns new, behavior-tagged DUs tied to the user's action for algorithm enrichment.
    *   **Downvote:** A user spends DUs to *remove* Influence from the post (Entropy).

### **B. Causal Time (The Relativity of Truth)**
In a decentralized system, "Global Time" is an illusion. "Causal Time" is the truth.
*   **The Chain of Causality:** A Reply (`Y`) cannot exist before its Parent (`X`). 
*   **The Feed Algorithm:**
    *   Do not sort solely by `timestamp` (subjective).
    *   Sort by `ReferenceDepth` and total internal `Energy` (Data Units).

---

## **4. Value Extraction (From Social to Economic)**

### **A. Converting Fuel to Value**
*   **Accumulation:** Objects that harbor high-quality "Threads" accrue massive amounts of Data Units from user interaction.
*   **Extraction:** The author of a highly energized object can extract real value. The protocol tracks the mass concentration of DUs within popular objects, rewarding creators with proportional fractions of the native TCOIN emission pool (connecting the social layer back to the main economic model).

### **B. Reputation as Flow (Not Stock)**
Reputation is not a bank account; it is a muscle.
*   **Decay Function:** Reputation (`R`) decays over time if a user stops participating.
*   **Work Function:** Validating blocks, voting on governance, or creating high-energy posts generates new Reputation.

---

## **5. Integration Plan (From Theory to Code)**

### **Phase 1: The Social Manager (Current State)**
We are currently simulating these laws in the Go `SocialManager`:

1.  **Fossilization:**
    *   `Status` behaves as a function of `Energy` (Data Units).
    *   **Threshold:** Posts needing revival fall into `fossilized` state.

2.  **Data Unit Voting System:**
    *   **Interaction:** We must convert current logic (which directly burns user balances) to track "Data Units" specific to the social context.

3.  **Causal Time (Implemented):**
    *   `LogicalTime` (Lamport Clock) orders the `GetFeed` priority above traditional `CreatedAt` stamps.

### **Phase 2: The Smart Contract (Future)**
We will deploy the `PostContract` to the VM to enforce the separation of TCOIN state arrays explicitly from the highly liquid "Data Unit" behavioral arrays.
```
