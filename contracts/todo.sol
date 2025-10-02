pragma solidity ^0.8.0;

contract TodoList {
    struct Task {
        string description;
        bool completed;
        uint256 timestamp;
        uint256 priority;
    }

    mapping(address => Task[]) private userTasks;

    event TaskAdded(address indexed user, uint256 taskId, string description);
    event TaskCompleted(address indexed user, uint256 taskId);
    event TaskDeleted(address indexed user, uint256 taskId);

    function addTask(string memory _description, uint256 _priority) public {
        require(bytes(_description).length > 0, "Description cannot be empty");
        require(_priority >= 1 && _priority <= 3, "Priority must be 1-3");
        
        userTasks[msg.sender].push(Task({
            description: _description,
            completed: false,
            timestamp: block.timestamp,
            priority: _priority
        }));
        
        emit TaskAdded(msg.sender, userTasks[msg.sender].length - 1, _description);
    }

    function markTaskDone(uint256 _taskId) public {
        require(_taskId < userTasks[msg.sender].length, "Invalid task ID");
        require(!userTasks[msg.sender][_taskId].completed, "Task already completed");
        
        userTasks[msg.sender][_taskId].completed = true;
        emit TaskCompleted(msg.sender, _taskId);
    }

    function deleteTask(uint256 _taskId) public {
        require(_taskId < userTasks[msg.sender].length, "Invalid task ID");
        
        Task[] storage tasks = userTasks[msg.sender];
        tasks[_taskId] = tasks[tasks.length - 1];
        tasks.pop();
        
        emit TaskDeleted(msg.sender, _taskId);
    }

    function updateTask(uint256 _taskId, string memory _newDescription) public {
        require(_taskId < userTasks[msg.sender].length, "Invalid task ID");
        require(bytes(_newDescription).length > 0, "Description cannot be empty");
        
        userTasks[msg.sender][_taskId].description = _newDescription;
    }

    function getTasks() public view returns (Task[] memory) {
        return userTasks[msg.sender];
    }

    function getTask(uint256 _taskId) public view returns (Task memory) {
        require(_taskId < userTasks[msg.sender].length, "Invalid task ID");
        return userTasks[msg.sender][_taskId];
    }

    function getTaskCount() public view returns (uint256) {
        return userTasks[msg.sender].length;
    }

    function getPendingTasks() public view returns (Task[] memory) {
        uint256 pendingCount = 0;
        Task[] storage allTasks = userTasks[msg.sender];
        
        for (uint256 i = 0; i < allTasks.length; i++) {
            if (!allTasks[i].completed) {
                pendingCount++;
            }
        }
        
        Task[] memory pending = new Task[](pendingCount);
        uint256 index = 0;
        
        for (uint256 i = 0; i < allTasks.length; i++) {
            if (!allTasks[i].completed) {
                pending[index] = allTasks[i];
                index++;
            }
        }
        
        return pending;
    }

    function getCompletedTasks() public view returns (Task[] memory) {
        uint256 completedCount = 0;
        Task[] storage allTasks = userTasks[msg.sender];
        
        for (uint256 i = 0; i < allTasks.length; i++) {
            if (allTasks[i].completed) {
                completedCount++;
            }
        }
        
        Task[] memory completed = new Task[](completedCount);
        uint256 index = 0;
        
        for (uint256 i = 0; i < allTasks.length; i++) {
            if (allTasks[i].completed) {
                completed[index] = allTasks[i];
                index++;
            }
        }
        
        return completed;
    }
}