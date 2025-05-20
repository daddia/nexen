# nexen

A high-performance, stateless, vendor-agnostic LLM orchestration gateway with context-aware routing and adaptive model selection, written in Go.

Nexen Platform provides a unified API for Large Language Models (LLMs). Nexen handles token management, dynamic model selection, context management, cost and reliability optimization, and fail‑over, so internal clients can focus on their business logic without worrying about provider details.

## Key Features

* **Unified Gateway**: Single REST/gRPC API to interact with any supported LLM.
* **Dynamic Model Selection**: Automatic routing to the optimal model based on task profile, cost, latency, and reliability.
* **Context Management**: Session history storage with intelligent pruning and semantic recall.
* **Cost & Reliability Optimiser**: Real‑time monitoring, budgeting, and SLA‑driven routing.
* **Degradation Detection**: Automatic fail‑over when a model’s performance degrades.
* **Extensible Connectors**: Plug‑in adapters for OpenAI, Anthropic, Vertex AI, and more via regex‑based registry.

## Repository Structure

```
nexen/
├── go.work                     # Go workspace for multi‑module monorepo
├── services/
│   ├── gateway/                # API Gateway service
│   ├── auth/                   # Authentication & IAM service
│   ├── sessions/               # Session state service
│   ├── context/                # Context management service
│   ├── orchestrator/           # Core request orchestration engine
│   ├── selection/              # Model selector service
│   ├── detection/              # Degradation detection service
│   ├── evaluation/             # Complexity evaluator service
│   └── connectors/             # Provider adapters & registry
│       ├── common/             # Base LLM interfaces & types
│       ├── openai/             # OpenAI adapter
│       └── anthropic/          # Anthropic adapter
│
├── models/                     # Shared DTOs & model metadata registry
│
├── libs/                       # Shared libraries (logging, metrics, telemetry)
│
├── tests/                      # End‑to‑end and integration tests
│
├── infrastructure/             # IaC (Kubernetes manifests, Helm, Terraform)
│
├── .github/                    # CI/CD workflows
└── docs/                       # Architecture and contribution guides
```

## Getting Started

1. **Clone the repository** and ensure you have Go 1.21+ installed.
2. **Initialize workspace**: run `go work sync` at the repo root.
3. **Build & Test** each service independently via `cd services/<service>` and `go test ./...`.
4. **Run the API Gateway** locally to begin sending LLM requests through Nexen.

For detailed build, test, and deployment instructions, refer to each module’s README under `services/` and the documentation in `docs/`.

## Contributing

We welcome contributions! Please read the `docs/`.

## License

This project is licensed under the [Apache 2.0 License](LICENSE).
