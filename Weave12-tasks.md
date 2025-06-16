# Weave12-tasks.md

## Epic: MAS System Refactor & Dynamic Orchestration

This task list covers all steps needed to implement dynamic agent definition, tool discovery, and GraphQL orchestration in the MAS system. Tasks are grouped by feature and ordered for maximum clarity and traceability. Each task should be checked off as completed.

---

## 1. Project Preparation
- [ ] Review and backup current codebase
- [ ] Set up a development branch for Weave12 refactor
- [ ] Document architectural goals in README

---

## 2. Tool Metadata & Discovery
### 2.1. Tool Metadata Definition
- [ ] Define a `ToolMetadata` struct (name, description, parameter schema, example usage)
- [ ] Update all tool implementations to provide metadata
- [ ] Create a global tool registry for metadata

### 2.2. Tool Registration
- [ ] Refactor tool registration to populate the global registry
- [ ] Ensure all tools in `tools/` and examples are registered

### 2.3. Tool Discovery API
- [ ] Add HTTP server using Go standard library
- [ ] Implement `/tools` endpoint returning JSON array of tool metadata
- [ ] Write unit tests for tool registry and discovery handler
- [ ] Write integration test: curl `/tools`, verify output
- [ ] Document tool discovery in README

### 2.4. (Optional) OpenAPI/Swagger
- [ ] Research Go OpenAPI/Swagger libraries
- [ ] Implement `/openapi.json` endpoint if desired
- [ ] Add Swagger UI if desired

---

## 3. GraphQL Orchestration Layer
### 3.1. Library & Schema Design
- [ ] Evaluate Go GraphQL libraries (gqlgen, graphql-go, etc.)
- [ ] Choose and install preferred library
- [ ] Define GraphQL schema for:
    - [ ] Tool queries
    - [ ] Agent queries
    - [ ] Agent mutations (create, update, delete)
    - [ ] Tool invocation (if desired)
- [ ] Document schema and reasoning in project docs

### 3.2. Agent Model Refactor
- [ ] Refactor agent constructors to accept config structs
- [ ] Define agent config schema (fields, types, validation)
- [ ] Implement unmarshalling from GraphQL input objects to Go structs
- [ ] Ensure agents can be created, updated, deleted at runtime
- [ ] Add persistence (in-memory first, DB/file optional)

### 3.3. GraphQL Server Implementation
- [ ] Scaffold GraphQL server entrypoint
- [ ] Implement resolvers for:
    - [ ] Tool discovery (query all tools)
    - [ ] Agent management (CRUD)
    - [ ] Tool invocation by agent (if needed)
- [ ] Add error handling and validation in resolvers
- [ ] Write unit tests for resolvers
- [ ] Write integration tests for GraphQL API
- [ ] Provide example queries/mutations in docs

### 3.4. Security & Auth (Optional)
- [ ] Decide on authentication/authorization needs
- [ ] Implement basic auth or token-based auth if required
- [ ] Add tests for auth logic

---

## 4. Testing & Validation
- [ ] Add `_test.go` files for all major packages (agent, tools, registry, GraphQL)
- [ ] Write shell scripts to:
    - [ ] Launch server
    - [ ] Curl `/tools` and GraphQL endpoints
    - [ ] Validate outputs
- [ ] Set up CI to run all tests and scripts
- [ ] Measure and document test coverage

---

## 5. Documentation & Examples
- [ ] Update README with:
    - [ ] New architecture overview
    - [ ] Tool discovery instructions
    - [ ] GraphQL usage and schema
    - [ ] Example agent definitions and queries
- [ ] Add example configs and GraphQL queries/mutations
- [ ] Provide migration guide for users of old system

---

## 6. Stretch Goals & Enhancements
- [ ] Dynamic tool registration/unregistration via API
- [ ] Real-time event subscription via GraphQL subscriptions
- [ ] Tool and agent versioning
- [ ] UI for tool/agent management (Swagger UI, GraphQL Playground, or custom)

---

## 7. Project Management
- [ ] Assign owners to each major task
- [ ] Define milestones and deadlines
- [ ] Review progress weekly
- [ ] Refine requirements and tasks as needed

---

## 8. Final Review & Launch
- [ ] Code review for all new/changed components
- [ ] End-to-end test: create agent via GraphQL, discover tools, invoke tool
- [ ] Finalize documentation and changelog
- [ ] Merge to main branch
- [ ] Tag release and announce

---

# End of Weave12-tasks.md
