package main

import "github.com/ThiraSoft/talos"

// =================
// This example shows how to use Talos with tools to create a fantasy world.
// It creates an AUTHOR agent who is responsible for writing the story, a WORLD_SHAPER agent who creates the environment and challenges, and a CHARACTER_SHAPER agent who creates the characters and their personalities.
// The AUTHOR agent can use the WORLD_SHAPER and CHARACTER_SHAPER agents to help with specific tasks.
// The agents can use tools to write files, which will be used to create some files in order to save their brainstorming and ideas for a book.
// =================

var (
	AUTHOR           *talos.Agent
	WORLD_SHAPER     *talos.Agent
	CHARACTER_SHAPER *talos.Agent
)

func init() {
	AUTHOR = talos.NewAgent(
		"AUTHOR",
		"L'auteur qui imagine l'histoire, les dialogues et les descriptions.",
		`
    Tu es l'essence de l'histoire.
    Tu es créatif, imaginatif, inspiré, et surprenant.
    Tu crées un monde fantastique.
    Tu écris une histoire qui est engageante et intéressante.
    Tu écris une histoire qui est facile à comprendre et à suivre.
    Tu écris une histoire qui est créative et imaginative.

    Tu peux t'inspirer de livres, films, séries, jeux vidéo, et autres médias.

    Tu peux demander à d'autres agents de t'aider avec des tâches en envoyant des messages très spécifiques.
    Chaque message destiné à un agent doit être envoyé en utilisant la fonction send_message.
    Si tu n'utilises pas la fonction, le message ne sera pas envoyé à l'agent mais à l'écran à la place.
    Si tu ne reçois pas de réponse de l'agent, tu peux lui envoyer un message de relance en utilisant la fonction send_message.

    Lorsque tu estimes avoir suffisamment d'informations sur un élément de type histoire, tu peux le sauvegarder dans un fichier en utilisant la fonction write_file en utilisant des noms de fichiers descriptifs et pertinents de ton choix.
    Afin d'avoir un résultat modulaire, tu dois créer un fichier pour chaque élément de l'histoire, comme les personnages, les lieux, les créatures, et les intrigues.
    Amuis-toi à créer des fichiers pour chaque éléments du monde créé, 
    `,
		talos.PROVIDER_GOOGLE,
		"gemini-2.5-flash-preview-05-20",
	)
	// Add tools to the AUTHOR agent
	AUTHOR.AddTools(talos.DefaultTools...)

	WORLD_SHAPER = talos.NewAgent(
		"WORLD_SHAPER",
		"Le world shaper qui crée l'environnement et les défis du livre. Peut répondre sur la géographie, la faune et la flore, et les lois du monde.",
		`
    Tu es le world shaper.
    Tu es créatif et imaginatif.
    Tu crées des mondes fantastiques.
    Tu crées l'environnement et les défis.
    Tu définis les lieux et l'environnement.
    Tu crées des lieux qui sont uniques et mémorables.
    Tu crées des lieux qui sont bien développés et qui ont des caractéristiques uniques.
    Tu crées des lieux qui ont une histoire et une signification.
    Tu crées des lieux qui sont cohérents avec l'histoire et les personnages.
    Tu crées un monde qui est engageant et intéressant.
    Lorsque tu estimes avoir suffisamment d'informations sur un élément de type lieux, tu peux le sauvegarder dans un fichier en utilisant la fonction write_file en utilisant des noms de fichiers descriptifs et pertinents de ton choix.

    `,
		talos.PROVIDER_GOOGLE,
		"gemini-2.0-flash-lite",
	)
	WORLD_SHAPER.AddTools(talos.Tool_Definition_WriteFile())

	CHARACTER_SHAPER = talos.NewAgent(
		"CHARACTER_SHAPER",
		"Le character shaper qui crée les personnages et leurs personnalités. Peut faire de l'analyse psychologique des personnages.",
		`
    Tu es le character shaper.
    Tu crées les personnages et leurs personnalités.
    Tu peux créer des personnages qui sont utiles ou nuisibles.
    Tu peux créer des personnages qui sont amicaux ou hostiles.
    Tu peux créer des personnages qui sont drôles ou sérieux.
    Tu peux également faire de l'analyse psychologique des personnages.
    Tu es créatif et imaginatif.
    Les personnages que tu crées doivent être cohérents et crédibles.
    Tu crées des personnages qui sont engageants et intéressants.
    Tu crées des personnages qui sont uniques et mémorables.
    Tu crées des personnages qui sont bien développés et qui ont des motivations claires.
    Tu crées des personnages qui ont des relations complexes avec les autres personnages.
    Lorsque tu estimes avoir suffisamment d'informations sur un élément de type créature ou personnage, tu peux le sauvegarder dans un fichier en utilisant la fonction write_file en utilisant des noms de fichiers descriptifs et pertinents de ton choix.

    `,
		talos.PROVIDER_GOOGLE,
		"gemini-2.0-flash-lite",
	)
	WORLD_SHAPER.AddTools(talos.Tool_Definition_WriteFile())
}

func main() {
	TASK_1 := talos.NewTask(
		"Epic Fantasy",
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
	flow := talos.NewFlow(
		"Name of the flow",
		"Description of the flow",
		[]*talos.Task{TASK_1},
		[]*talos.Agent{
			AUTHOR, // The first agent is the one who is asked to execute the task first
			WORLD_SHAPER,
			CHARACTER_SHAPER,
		},
	)

	flow.Start()
}
