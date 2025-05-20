# Project Structure

```
nexen/
├── go.work
├── services
│   ├── gateway/
│   │   ├── go.mod
│   │   ├── cmd/
│   │   ├── pkg/
│   │   └── README.md
│   ├── auth/            # Auth
│   ├── sessions/        # Sessions
│   ├── context/         # Context manager
│   ├── orchestration/   # Orchestration Engine 
│   ├── selection/       # Model Selector
│   ├── detection        # Degredation Detection
│   ├── evaluation       # Evaluator
│   ├── conntectors      # Providor-specific LLM code
│   │   ├── go.mod
│   │   ├── common/
│   │   ├── 
│   └── ...
│
├── models
│   ├── llm-request.go
│   ├── llm-response.go
│   └── registry.go 
│
├── libs
│   ├── logging
│   ├── metrics
│   ├── otel
│   └── ...
│ 
├── tests
│   ├── integration
│   ├── fixtures
│   └── config
│
├── infrastructure
│   ├── k8s
│   ├── helm
│   └── terraform
│
├── .github
├── docs
└── README.md
```
