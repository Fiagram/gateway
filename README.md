<font size= "5"> **Table Of Contents** </font>
- [Todo list](#todo-list)
- [Introduction](#introduction)
- [Environment for development](#environment-for-development)
  - [Service](#service)
    - [Preparation](#preparation)
    - [Run service as standalone](#run-service-as-standalone)


# Todo list
- [x] Integrate FX framework

# Introduction
- A gateway for microservices.

# Environment for development

## Service
### Preparation
- Download tools for the service
```
make init
```
- Generate openapi (by oapi-codegen tool) and others
```
make generate
```
### Run service as standalone
- Run service with the make command
```
make run-standalone-server
```

