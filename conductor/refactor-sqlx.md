# Refactor SQLX Queries

## Objective
The goal is to move the query execution functions from the core `sqb` package into a separate `run/sqlx` sub-package. This will decouple the core builder from the `sqlx` dependency, making it more modular and avoiding unnecessary dependencies for users who only need the builder.

## Key Files & Context
- `sqb.go`: The core SQL builder, which currently imports `github.com/jmoiron/sqlx` and `context`.
- `run/sqlx/sqlx.go`: A new file to contain the `sqlx`-specific execution logic.

## Proposed Solution

### 1. Modify `sqb.go`
- **Export `Statement`'s SQL string**: Add a `String() string` method to `Statement[T]`. This allows the execution package to retrieve the generated SQL string without exporting the `stmtString` field (which might be better kept unexported for encapsulation).
- **Remove `sqlx` and `context` imports**: Since `sqb.go` will no longer handle query execution, it no longer needs these dependencies.
- **Remove execution methods**: Delete the `NamedQueryStruct` and `NamedQueryStructContext` methods from the `Statement[T]` struct.

### 2. Create `run/sqlx/sqlx.go`
- **Package Name**: `sqlx` (inside the `run/sqlx` directory).
- **Import Aliasing**: Use an alias for `github.com/jmoiron/sqlx` (e.g., `jsqlx`) to avoid collision with the local package name.
- **Deduplication**: Implement a private helper `scanRows[T any](rows *jsqlx.Rows) ([]*T, error)` that handles the row iteration, `StructScan`, and final `rows.Err()` and `rows.Close()` checks as identified in the earlier code review.
- **Public Functions**: Implement `NamedQueryStruct` and `NamedQueryStructContext` as package-level functions that take an `sqb.Statement[T]` and an `*jsqlx.DB` as arguments.

## Verification
- **Compilation**: Run `go_diagnostics` for both the `sqb` package and the new `run/sqlx` package.
- **Logic Integrity**: Ensure the same error handling and iteration logic is preserved after the move.
