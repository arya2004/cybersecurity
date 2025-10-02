// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title VotingWithDeadline
 * @dev A voting contract that automatically disables voting after a deadline
 * Uses block.timestamp for time-based logic
 */
contract VotingWithDeadline {
    
    // Struct to store candidate information
    struct Candidate {
        uint256 id;
        string name;
        uint256 voteCount;
        bool exists;
    }
    
    // Contract owner (election commissioner)
    address public owner;
    
    // Voting deadline timestamp
    uint256 public votingDeadline;
    
    // Election title/description
    string public electionTitle;
    
    // Total number of candidates
    uint256 public candidateCount;
    
    // Total number of votes cast
    uint256 public totalVotes;
    
    // Mapping to store candidates by ID
    mapping(uint256 => Candidate) public candidates;
    
    // Mapping to track if an address has voted
    mapping(address => bool) public hasVoted;
    
    // Mapping to track which candidate an address voted for
    mapping(address => uint256) public voterChoice;
    
    // Events
    event CandidateRegistered(uint256 indexed candidateId, string name);
    event VoteCast(address indexed voter, uint256 indexed candidateId);
    event VotingEnded(uint256 timestamp);
    event DeadlineExtended(uint256 newDeadline);
    
    // Modifiers
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can perform this action");
        _;
    }
    
    modifier votingActive() {
        require(block.timestamp < votingDeadline, "Voting period has ended");
        _;
    }
    
    modifier votingEnded() {
        require(block.timestamp >= votingDeadline, "Voting is still active");
        _;
    }
    
    modifier validCandidate(uint256 _candidateId) {
        require(
            _candidateId > 0 && _candidateId <= candidateCount && candidates[_candidateId].exists,
            "Invalid candidate ID"
        );
        _;
    }
    
    /**
     * @dev Constructor to initialize the voting contract
     * @param _electionTitle Title of the election
     * @param _votingDurationInMinutes Duration of voting period in minutes
     */
    constructor(string memory _electionTitle, uint256 _votingDurationInMinutes) {
        owner = msg.sender;
        electionTitle = _electionTitle;
        votingDeadline = block.timestamp + (_votingDurationInMinutes * 1 minutes);
        candidateCount = 0;
        totalVotes = 0;
    }
    
    /**
     * @dev Register a new candidate (only owner)
     * @param _name Name of the candidate
     */
    function registerCandidate(string memory _name) external onlyOwner votingActive {
        require(bytes(_name).length > 0, "Candidate name cannot be empty");
        
        candidateCount++;
        candidates[candidateCount] = Candidate({
            id: candidateCount,
            name: _name,
            voteCount: 0,
            exists: true
        });
        
        emit CandidateRegistered(candidateCount, _name);
    }
    
    /**
     * @dev Cast a vote for a candidate
     * @param _candidateId ID of the candidate to vote for
     */
    function vote(uint256 _candidateId) external votingActive validCandidate(_candidateId) {
        require(!hasVoted[msg.sender], "You have already voted");
        require(candidateCount > 0, "No candidates registered");
        
        // Mark as voted
        hasVoted[msg.sender] = true;
        voterChoice[msg.sender] = _candidateId;
        
        // Increment vote count
        candidates[_candidateId].voteCount++;
        totalVotes++;
        
        emit VoteCast(msg.sender, _candidateId);
    }
    
    /**
     * @dev Get voting results (can be called anytime, but final results after deadline)
     * @return winnerName Name of the candidate with most votes
     * @return winnerVotes Number of votes for the winner
     * @return isFinal Whether the results are final (voting ended)
     */
    function getResults() external view returns (
        string memory winnerName,
        uint256 winnerVotes,
        bool isFinal
    ) {
        require(candidateCount > 0, "No candidates registered");
        
        uint256 maxVotes = 0;
        uint256 winnerId = 0;
        
        // Find candidate with most votes
        for (uint256 i = 1; i <= candidateCount; i++) {
            if (candidates[i].voteCount > maxVotes) {
                maxVotes = candidates[i].voteCount;
                winnerId = i;
            }
        }
        
        winnerName = candidates[winnerId].name;
        winnerVotes = candidates[winnerId].voteCount;
        isFinal = block.timestamp >= votingDeadline;
    }
    
    /**
     * @dev Get all candidates and their vote counts
     * @return candidateIds Array of candidate IDs
     * @return candidateNames Array of candidate names
     * @return voteCounts Array of vote counts
     */
    function getAllCandidates() external view returns (
        uint256[] memory candidateIds,
        string[] memory candidateNames,
        uint256[] memory voteCounts
    ) {
        candidateIds = new uint256[](candidateCount);
        candidateNames = new string[](candidateCount);
        voteCounts = new uint256[](candidateCount);
        
        for (uint256 i = 1; i <= candidateCount; i++) {
            candidateIds[i - 1] = candidates[i].id;
            candidateNames[i - 1] = candidates[i].name;
            voteCounts[i - 1] = candidates[i].voteCount;
        }
    }
    
    /**
     * @dev Get voting status information
     * @return isActive Whether voting is still active
     * @return timeRemaining Seconds remaining until deadline (0 if ended)
     * @return totalVotesTotalCast Total number of votes cast so far
     */
    function getVotingStatus() external view returns (
        bool isActive,
        uint256 timeRemaining,
        uint256 totalVotesTotalCast
    ) {
        isActive = block.timestamp < votingDeadline;
        timeRemaining = isActive ? votingDeadline - block.timestamp : 0;
        totalVotesTotalCast = totalVotes;
    }
    
    /**
     * @dev Extend voting deadline (only owner, only if voting is active)
     * @param _additionalMinutes Additional minutes to extend the deadline
     */
    function extendDeadline(uint256 _additionalMinutes) external onlyOwner votingActive {
        require(_additionalMinutes > 0, "Extension must be positive");
        
        votingDeadline += (_additionalMinutes * 1 minutes);
        emit DeadlineExtended(votingDeadline);
    }
    
    /**
     * @dev Check if a specific address has voted
     * @param _voter Address to check
     * @return voted Whether the address has voted
     * @return candidateId ID of candidate voted for (0 if not voted)
     */
    function getVoterInfo(address _voter) external view returns (bool voted, uint256 candidateId) {
        voted = hasVoted[_voter];
        candidateId = voterChoice[_voter];
    }
    
    /**
     * @dev Emergency function to end voting immediately (only owner)
     */
    function emergencyEndVoting() external onlyOwner votingActive {
        votingDeadline = block.timestamp;
        emit VotingEnded(block.timestamp);
    }
    
    /**
     * @dev Get detailed election information
     */
    function getElectionInfo() external view returns (
        string memory title,
        address electionOwner,
        uint256 deadline,
        uint256 totalCandidates,
        uint256 totalVotesCast,
        bool isActive
    ) {
        title = electionTitle;
        electionOwner = owner;
        deadline = votingDeadline;
        totalCandidates = candidateCount;
        totalVotesCast = totalVotes;
        isActive = block.timestamp < votingDeadline;
    }
}