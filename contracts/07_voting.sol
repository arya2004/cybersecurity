// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract VotingSystem {
    enum ElectionStage {
        Created,
        Voting,
        Ended
    }

    struct Candidate {
        uint256 id;
        string name;
        uint256 voteCount;
    }

    mapping(address => bool) public hasVoted;

    address public electionCommissioner;

    ElectionStage public stage;

    uint256 public candidateCount;

    mapping(uint256 => Candidate) public candidates;

    event ElectionStageChanged(ElectionStage newStage);
    event Voted(address voter, uint256 candidateId);

    constructor() {
        electionCommissioner = msg.sender;
        stage = ElectionStage.Created;
    }

    modifier onlyCommissioner() {
        require(
            msg.sender == electionCommissioner,
            "Only the election commissioner can perform this action."
        );
        _;
    }

    modifier atStage(ElectionStage _stage) {
        require(stage == _stage, "Function cannot be called at this time.");
        _;
    }

    function registerCandidate(string memory _name)
        public
        onlyCommissioner
        atStage(ElectionStage.Created)
    {
        candidateCount++;
        candidates[candidateCount] = Candidate(candidateCount, _name, 0);
    }

    function startVoting()
        public
        onlyCommissioner
        atStage(ElectionStage.Created)
    {
        stage = ElectionStage.Voting;
        emit ElectionStageChanged(stage);
    }

    function vote(uint256 _candidateId) public atStage(ElectionStage.Voting) {
        require(!hasVoted[msg.sender], "You have already voted.");
        require(
            _candidateId > 0 && _candidateId <= candidateCount,
            "Invalid candidate ID."
        );

        hasVoted[msg.sender] = true;
        candidates[_candidateId].voteCount++;

        emit Voted(msg.sender, _candidateId);
    }

    function endVoting() public onlyCommissioner atStage(ElectionStage.Voting) {
        stage = ElectionStage.Ended;
        emit ElectionStageChanged(stage);
    }

    function getWinner()
        public
        view
        atStage(ElectionStage.Ended)
        returns (string memory winnerName, uint256 voteCount)
    {
        uint256 maxVotes = 0;
        uint256 winnerId = 0;

        for (uint256 i = 1; i <= candidateCount; i++) {
            if (candidates[i].voteCount > maxVotes) {
                maxVotes = candidates[i].voteCount;
                winnerId = i;
            }
        }

        winnerName = candidates[winnerId].name;
        voteCount = candidates[winnerId].voteCount;
    }
}
