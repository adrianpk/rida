# ADR 0001: Event Processing vs. Event Streaming

## Context

This is a home assignment for a backend developer position. The system must process scooter events (trip start, trip end, location update) and update scooter state accordingly.

## Considered Options

- **Event Streaming (e.g., NATS):**
  - Would allow for decoupling, scalability, and reliability.
  - Adds complexity and infrastructure overhead.
  - Out of scope for a simple assignment.

- **Transactional Event Processing (Chosen):**
  - Simple, direct, and easy to reason about.
  - Events are processed synchronously and update scooter state immediately.
  - Fits the assignment's scope and keeps the codebase focused.

## Decision

For this assignment, a transactional approach is used: events are handled as they arrive and directly update scooter state. Event streaming was considered and would be a natural next step for a production system, but is not implemented here to keep things simple and relevant to the task.

## Consequences

- The code is easy to follow and maintain for the assignment.
- Future work could introduce event streaming for more advanced needs.
