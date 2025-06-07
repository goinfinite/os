# Project Specific Rules

- Unit tests CANNOT be executed on your local machine. They MUST be executed using a container built from the Containerfile.test.

# Infinite Standard Agent Guidelines

## @source https://github.com/goinfinite/tk/blob/main/docs/AGENTS.md

## General Rules

- NEVER create new features that weren't explicitly requested, even if they seem like useful additions.
- Ambiguity SHOULD always be questioned as it generates doubts about the intention of the developer.
- Value objects, infrastructure and use cases (with complex logic) MUST have unit tests.
- Unit tests SHOULD use testCases as much as possible.
- During delete operations, attempt to validate all constraints upfront rather than discovering them mid-operation.

## Code Style Rules

- Avoid comments unless strictly necessary.
- Method names should focus on the conceptual purpose rather than implementation details.
- NEVER use 'else' statements unless it's the UI layer.
- NEVER use single letter variable names. Use descriptive names, but avoid long names.
- When naming variables, try to choose names that convey the intention or purpose rather than just describing what the variable stores.
- Variable names should reflect the primary flow, not conditional outcomes.
- Use purposeful named return values whenever the method returns multiple values.
- Use Ptr suffix on variables when parsing optional fields (usually pointers on DTOs).
- Prefer value objects as custom primitive types rather than structs when possible.
- Boolean variables should start with "Is", "Should", "Has" etc prefixes.
- Use PascalCase format for the entire error message whenever possible.
- Prefer "Fail", "Error", "Invalid" suffixes instead of "FailedTo", "Cannot", "UnableTo" prefixes.
- Prefer "Read" prefix or "Factory" suffix instead of "Get" suffix (depending on context).
- Prefer suffixes instead of prefixes for struct and method names to preserve alphabetical order context.
- Avoid redundant prefixes in struct field names when the context is already clear from the struct name or surrounding code.
- Struct fields should be ordered by importance, followed by alphabetical order.
- Struct required fields should be placed before optional (pointer) fields.

## Go(lang) Specific Rules

- Prefer using slog.Error or slog.Debug instead of log.Printf depending on the gravity of the log.
- Value objects accept interface{}/any directly without the need for pre-assertion.
- When using struct constructors (New), use multiple arguments per line.
- Sequential method parameters of the same type should be combined together on the same line.
