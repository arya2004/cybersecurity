// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./voting_deadline.sol";

/**
 * @title VotingExample
 * @dev Example contract showing how to interact with VotingWithDeadline
 */
contract VotingExample {
    VotingWithDeadline public votingContract;
    
    event ExampleLog(string message, uint256 value);
    
    /**
     * @dev Deploy and setup an example election
     */
    function createExampleElection() external {
        // Create election with 2 hour duration
        votingContract = new VotingWithDeadline("Example City Council Election", 120);
        
        // Register candidates
        votingContract.registerCandidate("Alice Johnson - Progressive Party");
        votingContract.registerCandidate("Bob Smith - Conservative Party");
        votingContract.registerCandidate("Carol Davis - Independent");
        
        emit ExampleLog("Election created with 3 candidates", 120);
    }
    
    /**
     * @dev Simulate voting process
     */
    function simulateVoting() external {
        require(address(votingContract) != address(0), "Create election first");
        
        // Cast a vote (in real scenario, this would be called by different addresses)
        try votingContract.vote(1) {
            emit ExampleLog("Vote cast successfully", 1);
        } catch {
            emit ExampleLog("Voting failed - might have already voted or deadline passed", 0);
        }
    }
    
    /**
     * @dev Get current election status
     */
    function getElectionStatus() external view returns (
        string memory title,
        uint256 totalCandidates,
        uint256 totalVotes,
        bool isActive,
        uint256 timeRemaining
    ) {
        require(address(votingContract) != address(0), "No election created");
        
        (title, , , totalCandidates, totalVotes, isActive) = votingContract.getElectionInfo();
        (, timeRemaining, ) = votingContract.getVotingStatus();
    }
    
    /**
     * @dev Get election results
     */
    function getElectionResults() external view returns (
        string memory winnerName,
        uint256 winnerVotes,
        bool isFinal
    ) {
        require(address(votingContract) != address(0), "No election created");
        
        (winnerName, winnerVotes, isFinal) = votingContract.getResults();
    }
    
    /**
     * @dev Extend deadline by 30 minutes (for demonstration)
     */
    function extendDeadline() external {
        require(address(votingContract) != address(0), "No election created");
        
        try votingContract.extendDeadline(30) {
            emit ExampleLog("Deadline extended by 30 minutes", 30);
        } catch {
            emit ExampleLog("Failed to extend deadline - might have ended", 0);
        }
    }
}

/**
 * @title VotingTestHelper
 * @dev Helper contract for testing edge cases
 */
contract VotingTestHelper {
    
    /**
     * @dev Test contract with very short duration (for testing)
     */
    function createShortElection() external returns (VotingWithDeadline) {
        // Create election with 1 minute duration for testing
        VotingWithDeadline shortVoting = new VotingWithDeadline("Test Election", 1);
        
        // Add test candidates
        shortVoting.registerCandidate("Test Candidate A");
        shortVoting.registerCandidate("Test Candidate B");
        
        return shortVoting;
    }
    
    /**
     * @dev Test multiple voting scenarios
     */
    function testVotingScenarios(VotingWithDeadline voting) external returns (bool[] memory results) {
        results = new bool[](4);
        
        // Test 1: Valid vote
        try voting.vote(1) {
            results[0] = true;
        } catch {
            results[0] = false;
        }
        
        // Test 2: Double voting (should fail)
        try voting.vote(2) {
            results[1] = false; // Should fail
        } catch {
            results[1] = true; // Expected to fail
        }
        
        // Test 3: Invalid candidate (should fail)
        try voting.vote(999) {
            results[2] = false; // Should fail
        } catch {
            results[2] = true; // Expected to fail
        }
        
        // Test 4: Check if voted
        results[3] = voting.hasVoted(address(this));
        
        return results;
    }
}