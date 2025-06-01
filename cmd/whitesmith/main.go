package main

import (
	"whitesmith/internal/agents"
	"whitesmith/pkg/agentic"
)

func main() {
	TASK_1 := agentic.NewTask(
		"Epic Space Fantasy",
		`
    Décrire un monde fantastique épique avec des personnages et des intrigues captivantes.
    Chaque élément de l'histoire doit être soigneusement développé pour créer une expérience immersive.
    Ce qui va être créé sera sauvegardé dans des fichiers textes afin de servir de base pour l'écriture d'un livre.
    Il faut créer des personnages, bons et mauvais avec des nuances, des créatures, des lieux, des intrigues.
    Chaque personnage doit avoir une personnalité, un rôle et des interactions avec les autres personnages.
    Chaque créature doit être unique, avec des caractéristiques et des comportements distincts.
    Chaque lieu doit avoir une description, une histoire et un rôle dans l'intrigue.
    Chaque intrigue doit être développée avec des rebondissements, des conflits et des résolutions.
    Chaque élément doit être cohérent et contribuer à l'ensemble de l'histoire.

    Dans les fichiers textes, il faut utiliser des sections pour organiser les informations.
    Chaque section doit être clairement définie et structurée.
    Donne le plus de détails possible sur chaque éléments.
  `,
	)

	// Make the flow
	flow := agentic.NewFlow(
		"flow-1",
		"Epic Space Fantasy",
		"Description of the epic space fantasy story to be written.",
		[]*agentic.Task{TASK_1},
		[]*agentic.Agent{
			agents.AUTHOR,
			agents.WORLD_SHAPER,
			agents.CHARACTER_SHAPER,
		},
	)

	flow.Start()
}
