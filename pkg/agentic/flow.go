package agentic

import (
	"strings"

	"google.golang.org/genai"
)

type Flow struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Tasks       []*Task          `json:"tasks"`
	Agents      []*Agent         `json:"agents"`  // Agents involved in the flow
	History     []*genai.Content `json:"history"` // History of the flow
}

func (f *Flow) AddTask(task Task) {
	f.Tasks = append(f.Tasks, &task)
}

func (f *Flow) RemoveTask(taskId string) {
	for i, task := range f.Tasks {
		if task.Id == taskId {
			f.Tasks = append(f.Tasks[:i], f.Tasks[i+1:]...)
			break
		}
	}
}

func NewFlow(id, name, description string, tasks []*Task, agents []*Agent) *Flow {
	flow := &Flow{
		ID:          id,
		Name:        name,
		Description: description,
		Tasks:       tasks,
		Agents:      agents,
		History:     []*genai.Content{},
	}

	// Initialize the common history for each agent
	for _, agent := range agents {
		agent.History = flow.History
	}

	return flow
}

func (f *Flow) AddAgents(agents ...*Agent) {
	// Add all agents to the flow
	if len(agents) == 0 {
		return
	}

	for _, a := range agents {
		f.Agents = append(f.Agents, a)
		a.History = f.History
	}
}

func (f *Flow) Start() string {
	if len(f.Tasks) == 0 {
		return "No tasks to start."
	}
	if len(f.Agents) == 0 {
		return "No agents available to execute the tasks."
	}

	// Prepare agents instructions
	agents_info_instructions := `
  Your goal is to fulfill the tasks assigned to you. 
  You can notify the current task is done by saying #TASK_DONE in your response.
  You can send messages to other agents (or respond to them) using the send_message function.
  Here is the list of currently available agents: 
  `
	for _, agent := range f.Agents {
		agents_info_instructions += "- " + agent.Name + ": " + agent.Description + "\n"
	}
	f.Agents[0].SetInstructions(
		f.Agents[0].GetInstructions() + "\n\n" + agents_info_instructions,
	)

	// Iterate through each task and execute it with the available agents
	for _, t := range f.Tasks {
		for {
			resp, err := f.Agents[0].ChatWithRetry(t.Description, 5)
			if err != nil {
				return "Error executing task: " + err.Error()
			}
			if strings.Contains(resp, "#TASK_DONE") {
				break
			}

		}
	}

	return "Done"
}
