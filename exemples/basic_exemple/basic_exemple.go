package main

import "github.com/ThiraSoft/talos"

// =================
// This example shows how to use Talos to create a flow with multiple agents and tasks.
// =================

func main() {
	BOB := talos.NewAgent(
		"BOB",
		"Question asking agent",
		`
    Be very concise.
    You are the best agent to ask questions.
    You ask all the agents you know questions about the topic you choose.
    `,
		talos.PROVIDER_GOOGLE,
		talos.DEFAULT_MODEL, // Use the default model defined in agent.go
	)
	BOB.AddFunctionDeclarations(talos.Tool_Definition_SendMessage)

	PATRICK := talos.NewAgent(
		"PATRICK",
		"Dumb agent",
		`
    Be very concise.
    You have low IQ.
    Your answers have no actual sense.
    You are creative, but not in a good way.
    `,
		talos.PROVIDER_GOOGLE,
		talos.DEFAULT_MODEL, // Use the default model defined in agent.go
	)

	SANDY := talos.NewAgent(
		"SANDY",
		"Smart agent",
		`
    Be very concise.
    You know how to answer all questions in a witty way.
    `,
		talos.PROVIDER_GOOGLE,
		talos.DEFAULT_MODEL, // Use the default model defined in agent.go
	)

	TASK_1 := talos.NewTask(
		"TASK_EXEMPLE_QUESTION",
		`
    Get informations about a topic of your choice.
    You have to summarize the answers of all agents at the end.
  `,
	)

	TASK_2 := talos.NewTask(
		"TASK_EXEMPLE_HAIKU",
		`
    Make all agents write a haiku on a topic of your choice.
    You have to give back the answers of all agents at the end.
  `,
	)

	// Make the flow
	flow := talos.NewFlow(
		"Name of the flow",
		"Description of the flow",
		[]*talos.Task{TASK_1, TASK_2},
		[]*talos.Agent{
			BOB, // The first agent is the one who is asked to execute the task first
			PATRICK,
			SANDY,
		},
	)

	flow.Start()
}
