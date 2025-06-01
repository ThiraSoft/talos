package agentic

import "github.com/google/uuid"

type TaskState string

const (
	TO_PLAN   TaskState = "TO_PLAN"
	TO_DO     TaskState = "TO_DO"
	TO_REVIEW TaskState = "TO_REVIEW"
	TO_TEST   TaskState = "TO_TEST"
	DONE      TaskState = "DONE"
)

type Task struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	StateStage  int         `json:"state_stage"`          // Index of the current state in the StateFlow
	StateFlow   []TaskState `json:"state_flow,omitempty"` // Optional field to track state transitions
}

func NewTask(name, description string) *Task {
	id := uuid.New().String()
	return &Task{
		Id:          id,
		Name:        name,
		Description: description,
		StateStage:  0,
		StateFlow:   []TaskState{TO_PLAN, TO_DO, DONE},
	}
}

// func (task *Task) AddAgent(agent *Agent) {
// 	task.Agents = append(task.Agents, agent)
// }

func (t *Task) IsDone() bool {
	return t.StateStage == len(t.StateFlow)-1
}

// IsValidTaskstate checks if the given TaskState is valid.
func IsValidTaskState(t TaskState) bool {
	switch t {
	case TO_PLAN, TO_DO, TO_REVIEW, TO_TEST, DONE:
		return true
	default:
		return false
	}
}
