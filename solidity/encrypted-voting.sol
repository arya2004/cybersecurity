// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

contract EncryptedVoting {

    address public admin;


    uint256 public registrationEnd;
    uint256 public commitEnd;
    uint256 public revealEnd;

    enum Phase { Registration, Commit, Reveal, Ended }


    struct Candidate {
        string name;
        uint256 voteCount;
    }

    struct Voter {
        bool registered;
        bool committed;
        bool revealed;
        bytes32 commitHash;
        uint256 revealedVote; 
    }

    Candidate[] public candidates;
    mapping(address => Voter) public voters;

    event VoterRegistered(address voter);
    event VoteCommitted(address voter);
    event VoteRevealed(address voter, uint256 candidateId);
    event CandidateAdded(uint256 candidateId, string name);

    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin");
        _;
    }

    modifier inPhase(Phase ph) {
        require(currentPhase() == ph, "Invalid phase for this action");
        _;
    }

    constructor(
        uint256 _registrationDuration,
        uint256 _commitDuration,
        uint256 _revealDuration,
        string[] memory _candidateNames
    ) {
        require(_registrationDuration > 0 && _commitDuration > 0 && _revealDuration > 0, "Durations must be >0");
        admin = msg.sender;

        registrationEnd = block.timestamp + _registrationDuration;
        commitEnd = registrationEnd + _commitDuration;
        revealEnd = commitEnd + _revealDuration;


        for (uint256 i = 0; i < _candidateNames.length; i++) {
            candidates.push(Candidate({ name: _candidateNames[i], voteCount: 0 }));
            emit CandidateAdded(i, _candidateNames[i]);
        }
    }

    function currentPhase() public view returns (Phase) {
        if (block.timestamp <= registrationEnd) {
            return Phase.Registration;
        } else if (block.timestamp <= commitEnd) {
            return Phase.Commit;
        } else if (block.timestamp <= revealEnd) {
            return Phase.Reveal;
        } else {
            return Phase.Ended;
        }
    }

    function addCandidate(string calldata name) external onlyAdmin inPhase(Phase.Registration) {
        uint256 id = candidates.length;
        candidates.push(Candidate({ name: name, voteCount: 0 }));
        emit CandidateAdded(id, name);
    }

    function register() external inPhase(Phase.Registration) {
        Voter storage v = voters[msg.sender];
        require(!v.registered, "Already registered");
        v.registered = true;
        emit VoterRegistered(msg.sender);
    }

    function commitVote(bytes32 commitHash) external inPhase(Phase.Commit) {
        Voter storage v = voters[msg.sender];
        require(v.registered, "Not registered");
        require(!v.committed, "Already committed");
        v.commitHash = commitHash;
        v.committed = true;
        emit VoteCommitted(msg.sender);
    }

    function revealVote(uint256 candidateId, bytes32 salt) external inPhase(Phase.Reveal) {
        Voter storage v = voters[msg.sender];
        require(v.registered, "Not registered");
        require(v.committed, "Did not commit");
        require(!v.revealed, "Already revealed");
        require(candidateId < candidates.length, "Invalid candidate");

        bytes32 computed = keccak256(abi.encodePacked(candidateId, salt));
        require(computed == v.commitHash, "Reveal does not match commit");

        v.revealed = true;
        v.revealedVote = candidateId;

        candidates[candidateId].voteCount += 1;

        emit VoteRevealed(msg.sender, candidateId);
    }


    function candidateCount() external view returns (uint256) {
        return candidates.length;
    }

    function getCandidate(uint256 idx) external view returns (string memory name, uint256 votes) {
        require(idx < candidates.length, "Candidate index out of range");
        Candidate storage c = candidates[idx];
        return (c.name, c.voteCount);
    }

    function getWinners() external view returns (uint256[] memory) {
        require(currentPhase() == Phase.Ended, "Voting not ended yet");

        uint256 maxVotes = 0;
        for (uint256 i = 0; i < candidates.length; i++) {
            if (candidates[i].voteCount > maxVotes) {
                maxVotes = candidates[i].voteCount;
            }
        }

        uint256 count = 0;
        for (uint256 i = 0; i < candidates.length; i++) {
            if (candidates[i].voteCount == maxVotes) {
                count++;
            }
        }

        uint256[] memory winners = new uint256[](count);
        uint256 j = 0;
        for (uint256 i = 0; i < candidates.length; i++) {
            if (candidates[i].voteCount == maxVotes) {
                winners[j] = i;
                j++;
            }
        }
        return winners;
    }
}
