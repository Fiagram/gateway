<font size= "5"> **Table Of Contents** </font>
- [Todo list](#todo-list)
- [Introduction](#introduction)
- [Environment for development](#environment-for-development)
  - [Service](#service)
    - [Preparation](#preparation)
    - [Run service as standalone](#run-service-as-standalone)


# Todo list
- [x] Integrate FX framework
- [x] Add middleware layers to http server
- [x] Implement `users` methods use middleware in gin framework
- [ ] Complete `users` CRUD methods
- [ ] Research NGINX to make a reverse-proxy for microservice
- [ ] Vibe coding a webserver based on `openapi.yml`
- [ ] Deploy the distributed system to a real domain (fiagram.com)
- [ ] (Tech debt) Fork and rewrite openapi to add `RegisterWith<middleware_names>Mids` intent to group middlewares for convenience

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

