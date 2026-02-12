# Tasks API - Project Documentation & Architecture

## Overview
This project serves as a backend API for managing tasks. It is built with Go and SQLite, demonstrating a clean, layered architecture designed for scalability, maintainability, and testability.

## Architectural Evolution
The codebase has evolved through several stages of refactoring to adhere to Separation of Concerns (SoC) principles.

### Phase 1: The Monolith (`main.go`)
Initially, the entire application resided in a single `main.go` file. This included the HTTP server setup, database connection logic, data models, and business logic. While simple to start, this approach is difficult to scale and test.

### Phase 2: separation of Data & Models
We extracted the core data structures and database interactions into their own packages:
- **`models/`**: Defines the shape of data (e.g., `Task` struct) independent of how it's stored or served.
- **`repos/`**: Implements the Repository pattern. This layer handles all direct database operations (SQL queries), abstracting the storage implementation details from the rest of the app.

### Phase 3: Handler Abstraction (Current State)
To further decouple the transport layer (HTTP) from the application logic, we introduced the **`handlers/`** package.
- **`handlers/`**: Contains HTTP handlers (`List`, `Create`, `Delete`). These handlers parse requests, validate input, call the appropriate repository methods, and format responses.
- **Dependency Injection**: The `handlers` are initialized with dependencies (like the repository), allowing for easier testing and mocking in the future.

## Component Responsibilities

Think of this architecture like a restaurant:

### 1. `handlers/` (The Waiter)
This is the **Interface Layer**. It's the first place requests (orders) land.
- **Role**: This code talks directly to the user (via HTTP).
- **What it does**:
  - Checks if the request is valid (e.g., is the JSON correct?).
  - Asks the `repos` layer to do the heavy lifting (like "Get me all tasks").
  - Formats the result as JSON and sends it back with the right status code (e.g., 200 OK, 404 Not Found).
- **Why**: It keeps the "business logic" separate from the "database logic". The handler doesn't know *how* data is stored, only *how* to ask for it.

### 2. `repos/` (The Kitchen)
This is the **Data Access Layer**. It's where the actual work happens.
- **Role**: The only part of the code allowed to touch the database.
- **What it does**:
  - Runs the actual SQL queries (`SELECT * FROM tasks...`).
  - Converts raw database rows into Go objects (`models`).
  - Handles database connection errors.
- **Why**: If we change databases later (e.g., from SQLite to Postgres), we only change files in this folder. The rest of the app doesn't need to know!

### 3. `models/` (The Menu)
This is the **Domain Layer**. It defines what things *are*.
- **Role**: Simple blueprints for our data.
- **What it does**:
  - Contains strictly `struct` definitions (e.g., `Task { ID, Title, Completed }`).
  - Uses tags (like `json:"title"`) to tell Go how to convert them to/from JSON.
- **Why**: By keeping these definitions separate, both `handlers` and `repos` can agree on what a "Task" looks like without depending on each other directly (avoiding circular dependency errors).

## Roadmap & Next Steps

The immediate focus for the next iteration is **Reliability and Testing**.

### Immediate To-Do List
- [ ] **Unit Testing**: Implement unit tests for `handlers` and `repos`.
  - Mock the database dependency to test handlers in isolation.
  - Use an in-memory SQLite instance to test repositories.
- [ ] **Integration Testing**: End-to-end tests ensuring the API endpoints work correctly with a live database.
- [ ] **Configuration Management**: Move hardcoded values (like database paths) to environment variables or a config file.
- [ ] **Logging & Instrumentation**: Add structured logging for better observability.
