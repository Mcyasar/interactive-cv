# Agent Operating Rules & Definition of Done

## 1. Branching & Lifecycle
- Every task must be developed in a feature branch: `feat/task-<id>-<slug>`.
- Do NOT push directly to `main`.

## 2. Testing & Quality
- All unit tests must pass before opening a Pull Request (`go test ./...`).
- Code must follow the interfaces provided in the Issue description.

## 3. Strict Kanban Boundaries
- Status Flow: `Todo` -> `In Progress` -> `In Review`.
- **CRITICAL:** An agent must NEVER set a card or issue to `Done`. 
- Only create a PR, reference the issue (`Refs #<issue_id>`), and move the item to `In Review`.
- Merging to `main` and moving to `Done` is exclusively handled by the Human Reviewer.