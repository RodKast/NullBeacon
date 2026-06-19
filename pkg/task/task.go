package task

import "github.com/google/uuid"

type Task struct {
	ID      string
	AgentID string
	Command string
	Output  string
	Status  string
}

func NewTask(agentID, command string) *Task {
	return &Task{
		ID:      uuid.New().String(),
		AgentID: agentID,
		Command: command,
		Status:  "pending",
	}
}
