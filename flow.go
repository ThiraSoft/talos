package talos

import (
	"strings"

	"github.com/google/uuid"
	"google.golang.org/genai"
)

var Agents = []*Agent{}

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

func NewFlow(name, description string, tasks []*Task, agents []*Agent) *Flow {
	id := uuid.New().String()
	flow := &Flow{
		Id:          id,
		Name:        name,
		Description: description,
		Tasks:       tasks,
		Agents:      agents,
		History:     make([]*genai.Content, 0, 10000),
	}

	return flow
}

func (f *Flow) AddAgents(agents ...*Agent) {
	// Add all agents to the flow
	if len(agents) == 0 {
		return
	}

	f.Agents = append(f.Agents, agents...)
}

func (f *Flow) Start() string {
	Agents = f.Agents // Update the global agents list, fo the tools to access
	if len(f.Tasks) == 0 {
		return "No tasks to start."
	}
	if len(f.Agents) == 0 {
		return "No agents available to execute the tasks."
	}

	// Prepare agents instructions
	agents_info_instructions := `
  Your goal is to fulfill the tasks assigned to you. 
  You MUST notify when the task is completed by saying TASK_DONE in your response.
  You MUST use the send_message tool to contace other agents (or respond to them).
  If you don't use the send_message tool when you talk to another agent, the message will not be sent.
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
		firstCall := true
		for {
			msg := ""
			if firstCall {
				msg = t.Description
			} else {
				msg = "Please continue with the task: " + t.Description
			}

			resp, err := f.Agents[0].ChatWithRetry(msg, 5)
			if err != nil {
				return "Error executing task: " + err.Error()
			}
			if strings.Contains(resp, "TASK_DONE") {
				logger.Debug("Task completed by agent", "task_name", t.Name, "agent_name", f.Agents[0].Name)
				break
			}
		}
	}

	Agents = []*Agent{} // Clear the global agents list after the flow is done
	return "Done"
}
