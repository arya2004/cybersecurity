# Voting Contract with Deadline

This Solidity contract implements a time-limited voting system that automatically disables voting after a specified deadline using `block.timestamp`.

## Features

### Core Functionality
- ✅ **Time-limited voting**: Automatically disables voting after deadline
- ✅ **Candidate registration**: Owner can register candidates before/during voting
- ✅ **Secure voting**: Each address can only vote once
- ✅ **Real-time results**: View results anytime (final after deadline)
- ✅ **Vote tracking**: Track who voted and their choices

### Advanced Features
- 🕐 **Deadline extension**: Owner can extend voting period
- 📊 **Comprehensive results**: Get detailed voting statistics
- 🚨 **Emergency controls**: Owner can end voting immediately
- 📈 **Status monitoring**: Check voting status and time remaining

## Contract Structure

### Main Components

```solidity
struct Candidate {
    uint256 id;
    string name;
    uint256 voteCount;
    bool exists;
}
```

### Key Mappings
- `mapping(uint256 => Candidate) public candidates` - Stores candidate data
- `mapping(address => bool) public hasVoted` - Tracks voting status
- `mapping(address => uint256) public voterChoice` - Records vote choices

### Time Logic
- Uses `block.timestamp` for current time
- Voting deadline set in constructor
- Automatic voting disabling after deadline

## Usage Examples

### 1. Deploy Contract
```solidity
// Deploy with election title and duration in minutes
VotingWithDeadline voting = new VotingWithDeadline("City Mayor Election 2025", 1440); // 24 hours
```

### 2. Register Candidates (Owner Only)
```solidity
voting.registerCandidate("Alice Johnson");
voting.registerCandidate("Bob Smith");
voting.registerCandidate("Carol Davis");
```

### 3. Cast Votes
```solidity
// Users vote by candidate ID (1, 2, 3, etc.)
voting.vote(1); // Vote for Alice Johnson
```

### 4. Check Results
```solidity
(string memory winner, uint256 votes, bool final) = voting.getResults();
```

### 5. Monitor Voting Status
```solidity
(bool active, uint256 timeLeft, uint256 totalVotes) = voting.getVotingStatus();
```

## Function Reference

### Public Functions

#### `registerCandidate(string memory _name)`
- **Access**: Owner only
- **Timing**: Before/during voting
- **Purpose**: Add new candidate to election

#### `vote(uint256 _candidateId)`
- **Access**: Any address (once)
- **Timing**: During voting period only
- **Purpose**: Cast vote for specific candidate

#### `getResults()`
- **Returns**: Winner name, vote count, finality status
- **Purpose**: Get current/final election results

#### `getAllCandidates()`
- **Returns**: Arrays of IDs, names, and vote counts
- **Purpose**: Get complete candidate information

#### `getVotingStatus()`
- **Returns**: Active status, time remaining, total votes
- **Purpose**: Monitor election progress

#### `extendDeadline(uint256 _additionalMinutes)`
- **Access**: Owner only
- **Timing**: During voting only
- **Purpose**: Extend voting period

## Security Features

### Access Control
- Owner-only functions protected by `onlyOwner` modifier
- Voting functions protected by time-based modifiers

### Voting Integrity
- One vote per address enforcement
- Candidate validation for all votes
- Immutable vote records

### Time Security
- Uses `block.timestamp` for reliable time tracking
- Automatic deadline enforcement
- No manual intervention required for ending

## Events

```solidity
event CandidateRegistered(uint256 indexed candidateId, string name);
event VoteCast(address indexed voter, uint256 indexed candidateId);
event VotingEnded(uint256 timestamp);
event DeadlineExtended(uint256 newDeadline);
```

## Error Messages

- `"Only owner can perform this action"` - Unauthorized access
- `"Voting period has ended"` - Action after deadline
- `"You have already voted"` - Duplicate voting attempt
- `"Invalid candidate ID"` - Non-existent candidate
- `"No candidates registered"` - Empty election

## Testing Scenarios

### Basic Flow
1. Deploy contract with 60-minute duration
2. Register 3 candidates
3. Cast votes from different addresses
4. Check intermediate results
5. Wait for deadline to pass
6. Verify final results

### Edge Cases
- Vote after deadline (should fail)
- Vote for non-existent candidate (should fail)
- Double voting (should fail)
- Extend deadline multiple times
- Emergency end voting

## Gas Considerations

- Candidate registration: ~100,000 gas
- Voting: ~50,000 gas
- Result queries: ~30,000 gas (view function)
- Deadline extension: ~25,000 gas

## Deployment Notes

1. Set appropriate voting duration for your use case
2. Consider time zones when setting deadlines
3. Test thoroughly on testnet before mainnet deployment
4. Verify contract source code for transparency

## License

MIT License - See LICENSE file for details.