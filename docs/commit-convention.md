# Commit Convention

This project uses Conventional Commits as the required standard for all commit messages.

## Rule

All commits must follow this format:

<type>(optional-scope): <description>

Examples:

- feat(api): add portfolio summary endpoint
- fix(investment): correct monthly interest calculation
- docs(readme): update setup instructions
- refactor(account): simplify repository query

## Allowed Types

- feat: new feature
- fix: bug fix
- docs: documentation only changes
- style: formatting changes (no code behavior changes)
- refactor: code change that is neither a fix nor a feature
- test: adding or updating tests
- chore: maintenance tasks (build, tooling, dependencies)

## Additional Guidelines

- Use imperative mood in the description.
- Keep the description short and objective.
- Use scope when it adds context.
- Avoid generic messages like "update" or "changes".

## Examples to Avoid

- update files
- ajustes
- fix stuff
- commit final

## Enforcement

Conventional Commits are the default and expected commit pattern for this repository.
