---
name: python-senior
description: Skill for Python implementation, refactoring, and test design. Use when class-oriented implementation style, pytest coverage execution, and strict naming conventions are required for coding, review, maintenance, and test additions.
---

# Python Senior

Follow these rules for Python coding.

## Python

- When creating functions, implement them as instance methods of a class by default.
- When creating tests, create a test class and implement tests as instance methods.
- Run module tests with `python -m pytest --cov=src --cov-branch --tb=short -vv`.

## Naming Conventions

Use PascalCase for:
- Classes
- Exceptions
- Test classes

Use lower_snake_case for:
- Directory names
- File names
- Modules
- Methods
- Functions
- Variables

Use lower_snake_case starting with `test_` for:
- Test methods
- Test functions

Use UPPER_SNAKE_CASE for:
- Environment variables
- Constants
- Global configurations

Use PascalCase for:
- Components
- Type definitions
- Interfaces

Use kebab-case for the following:
- Directory names (e.g., components/auth-wizard)
- File names (e.g., user-profile.tsx)

Use camelCase for the following:
- Variables
- Functions
- Methods
- Hooks
- Properties
- Props

Use uppercase for the following:
- Environment variables
- Constants
- Global configurations

For boolean variables, use a verb as prefix: `isLoading`, `hasError`, `canSubmit`
