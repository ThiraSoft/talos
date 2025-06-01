llm-orchestrator/
├── cmd/ # Point d'entrée de ton app (main.go)
│ └── whitesmith/
│   └── main.go
├── internal/
│ ├── core/ # Logique métier (routing, decision, agent manager)
│ │ ├── agent.go
│ │ ├── router.go
│ │ └── delegator.go
│ ├── llm/ # Abstraction LLM (OpenAI, LM Studio, etc.)
│ │ └── client.go
│ ├── infra/ # Intégrations externes (HTTP clients, fichiers YAML, DB)
│ │ ├── openai/
│ │ ├── ollama/
│ │ └── http/
│ └── config/ # Gestion de la config (agents, clés API, etc.)
│   └── config.go
├── pkg/ # Code réutilisable/exportable (ex: utils, logging, etc.)
├── api/ # Interfaces d’entrée : REST, gRPC, WebSocket
│ └── rest/
│    ├── handler.go
│    └── server.go
├── scripts/ # Scripts d’aide, génération, init de modèles
├── examples/ # Cas d’usage ou exécution simple d’agents
├── agent-protocol/ # Spécification du protocole d’agent (README + schema)
├── go.mod
├── README.md
└── LICENSE
