# API Directory

This directory contains all API-related files.

## Structure

```shell
api
├── README.md
├── generated
│   └── alarms.generated.go           # Generated Go code (do not edit manually)
├── openapi.yaml                      # OpenAPI specification
└── tools
    ├── generate.go                   # Code generation script
    └── oapi-codegen.yaml             # Config file for `oapi-codegen` 
```

## How to Generate Code

```shell
go generate ./...
```