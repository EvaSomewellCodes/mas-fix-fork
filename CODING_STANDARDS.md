# MAS Coding Standards & Philosophy

**Version:** 1.0  
**Last Updated:** 2025-06-15

---

## Overview

This document defines the coding standards, best practices, and guiding philosophy for all MAS backend code. Our aim is to produce code that is robust, readable, maintainable, and a joy for engineers (and agents!) to work with.

---

## Core Principles

- **Well-Commented & Readable:**
  - Every significant function, struct, and module must have clear doc comments—preferably in the Norwegian Scala style: rich in metaphor, intent, and future-facing TODOs.
  - Comments should clarify *why* as well as *how*.

- **DRY (Don't Repeat Yourself):**
  - Avoid code duplication through clear abstractions, helpers, and idiomatic Go interfaces.

- **Separation of Concerns:**
  - Each package, file, and function should have a single, clear responsibility.
  - Business logic, data access, and API layers must be cleanly separated.

- **TDD (Test-Driven Development):**
  - Write tests before or alongside code. All business logic must be covered by meaningful unit and integration tests.
  - Favor table-driven tests and clear test naming.

- **SOLID Principles:**
  - Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, Dependency Inversion.
  - Favor composition over inheritance.

- **Idiomatic Go:**
  - Use Go naming conventions, error handling, and package structure.
  - Keep functions short, favor explicitness, and return errors early.
  - Prefer slices and maps over arrays, and interfaces over concrete types where appropriate.

- **Pythonated Reasoning:**
  - Write code that is beautiful, functional, and intuitive—embracing the "Zen of Python" (explicit is better than implicit, simple is better than complex) even in Go.
  - When in doubt, optimize for clarity and maintainability.

- **Beautiful, Functional, Intuitive:**
  - Code should be easy to read, easy to use, and a delight to maintain.
  - If you’re proud to show it to a new teammate (or a Norwegian bicycle-riding, coffee-drinking, silent-coding engineer), you’re on the right track.

---

## Documentation

- All exported types, functions, and packages must have doc comments.
- Use Norwegian Scala-style doc blocks for major modules and architectural components.
- Keep README, API_REFERENCE, and this document up to date.

---

## Code Review & Collaboration

- Every change must be peer reviewed.
- Reviewers should check for adherence to these standards, clarity, and test coverage.
- Feedback should be constructive, specific, and focused on improvement.

---

## Humor & Culture

- We believe great code is written with focus, humility, and the occasional cup of strong coffee.
- If you find yourself riding a bicycle to work, coding in silence, and writing doc blocks about fjords, you’re in the right place.

---

## Verbose, Contextually Rich Errors

- All errors must be loud, clear, and provide as much context as possible:
  - **What** happened
  - **Where** it happened (file, function, line, or operation)
  - **Why** it happened (root cause, inputs, or preconditions)
- Use Go best practices:
  - Wrap errors with `fmt.Errorf` and `%w` to preserve stack/context
  - Include key variables and operation details in error messages
  - Never swallow errors—bubble them up with added context
- Error messages should help both users and developers debug quickly, even if they are new to the codebase.
- Prefer errors that are actionable and suggest next steps when possible.

---

## Appendix: Example Norwegian Scala-Style Doc Block

```go
/**
 * Plan represents the orchestrated set of steps for a MAS project.
 * Norwegian-style doc: Like a map drawn by candlelight before a journey through the fjords, a Plan guides agents, records memory, and promises deliverables yet unseen.
 * Fields: ID, Name, Steps, Agents, Deliverables, Manifest, etc.
 */
```
