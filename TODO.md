# Junior Engineer Learning Path

This roadmap is designed to help you master backend development concepts step-by-step. Each level introduces new challenges that simulate real-world tasks.

## Level 1: The Basics (Validation & Logic)
**Goal:** Understand how to handle user input safely.
- [ ] **Validate Creating Tasks**: Currently, you can create a task with an empty title. Add a check in `handlers/tasks.go` to return a `400 Bad Request` if `title` is empty.
- [ ] **Prevent Duplicate Titles**: BEFORE inserting into the database, check if a task with the same title already exists. Return a helpful error message.
  - *Hint*: You'll need to add a method to `repos/tasks.go` first!

## Level 2: Expanded Functionality (CRUD Mastery)
**Goal:** Practice adding new endpoints and database interaction.
- [ ] **Get Single Task**: Implement `GET /tasks/{id}`.
  - Add `GetByID(id int)` to the Repository.
  - Add the handler.
  - If the ID doesn't exist, return `404 Not Found`.
- [ ] **Update Task**: Implement `PUT /tasks/{id}` to mark a task as completed without deleting it.
  - You'll need an `Update(id int, completed bool)` function in the repository.

## Level 3: Database & Data Integrity
**Goal:** Learn how to manage database schemas and migrations.
- [ ] **Add `updated_at` Column**:
  - Modify the `CREATE TABLE` SQL in `repos/tasks.go`.
  - Update the `Task` struct in `models/tasks.go`.
  - Update the `Update` function to set `updated_at = CURRENT_TIMESTAMP`.
- [ ] **Soft Delete**: instead of actually deleting the row (hard delete), add a `deleted_at` column.
  - Update `List()` to only show tasks where `deleted_at` is NULL.

## Level 4: Reliability (Testing)
**Goal:** Ensure your code works and doesn't break in the future.
- [ ] **Test the "Business Logic"**: Write a unit test for your validation logic in Level 1.
- [ ] **Test the "Database Interction"**: Use a test database (file: `test.db`) to verify that `Create` actually saves data and `List` returns it.

## Level 5: Professional Grade (Structure & Config)
**Goal:** Make the app production-ready.
- [ ] **Configuration**: Stop hardcoding `"./tasks.db"` and `":8080"`.
  - Use a package like `os` or `godotenv` to read these from Environment Variables.
- [ ] **Logging Middleware**: Create a middleware that logs every request:
  - `[INFO] GET /tasks - 200 OK (12ms)`
  - *Hint*: Learn about `http.Handler` chaining.
