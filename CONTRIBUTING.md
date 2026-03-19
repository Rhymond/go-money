# Contributing to go-money

Thank you for taking the time to contribute! This document explains how to
report bugs, request features, and submit pull requests.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Reporting Bugs](#reporting-bugs)
- [Requesting Features](#requesting-features)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Pull Request Guidelines](#pull-request-guidelines)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Commit Messages](#commit-messages)

---

## Code of Conduct

Please be respectful and constructive in all interactions. Harassment of any
kind is not tolerated.

---

## Reporting Bugs

1. Search [existing issues](https://github.com/im-adarsh/go-money/issues) first
   to avoid duplicates.
2. Open a [new issue](https://github.com/im-adarsh/go-money/issues/new) with:
   - A clear, descriptive title.
   - A minimal reproducible example (Go Playground link preferred).
   - Expected vs. actual behaviour.
   - Go version (`go version`) and OS.

---

## Requesting Features

Open a [new issue](https://github.com/im-adarsh/go-money/issues/new) tagged
**enhancement**. Describe:

- The problem you are trying to solve.
- Your proposed API (function signatures, types).
- Any alternative approaches you considered.

---

## Development Setup

### Prerequisites

- Go 1.21 or later (`go version`)
- Git

### Clone and test

```bash
git clone https://github.com/im-adarsh/go-money.git
cd go-money
go test -race ./...
```

No external tools are required — the library has zero dependencies.

### Linting (optional but recommended)

```bash
# Install golangci-lint — https://golangci-lint.run/usage/install/
golangci-lint run
```

---

## Making Changes

1. **Fork** the repository on GitHub.
2. **Create a branch** from `master`:

   ```bash
   git checkout -b feat/my-feature
   # or
   git checkout -b fix/my-bug-fix
   ```

3. Make your changes (see [Coding Standards](#coding-standards)).
4. Add or update **tests** (see [Testing](#testing)).
5. Run the full test suite:

   ```bash
   go test -race ./...
   ```

6. Push your branch and open a pull request.

---

## Pull Request Guidelines

- Keep PRs focused — one logical change per PR.
- Reference the related issue (e.g. `Closes #42`).
- Ensure all CI checks pass before requesting review.
- Write a clear description of **what** the change does and **why**.
- Add entries to the table in README.md if you add new currencies or public
  API surface.

---

## Coding Standards

- Follow standard Go style (`gofmt`, `go vet`).
- All exported types, functions, and methods **must** have doc comments.
- Keep the library dependency-free — do not add `go.sum` entries.
- New monetary operations must be **immutable** — return a new `*Money`
  instead of modifying the receiver.
- Error messages should be lowercase, no trailing period
  (e.g. `"currencies don't match"`).

### Adding a new currency (formatting)

Edit `currency.go` and add an entry to the `currencies` map:

```go
"XYZ": {
    Decimal:  ".",
    Thousand: ",",
    Code:     "XYZ",
    Fraction: 2,
    Grapheme: "X$",
    Template: "$1",
},
```

### Adding a new currency (words)

Edit `currencyToWords.go` and add an entry to `CountryCurrencyMeta`:

```go
"XYZ": {MainWord: "zorkmid", SubWord: "zork"},
```

Or at runtime:

```go
money.AddCurrencyMeta("XYZ", "zorkmid", "zork")
```

---

## Testing

- All new code must be covered by tests.
- Tests live in `*_test.go` files alongside the source they test.
- Use table-driven tests wherever possible.
- Run with the race detector:

  ```bash
  go test -race ./...
  ```

- Check coverage:

  ```bash
  go test -coverprofile=cover.out ./...
  go tool cover -html=cover.out
  ```

---

## Commit Messages

Use the conventional format:

```
<type>(<scope>): <short summary>

<optional body>
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`.

Examples:

```
feat(money): add Percentage method
fix(formatter): handle zero-fraction currencies correctly
docs(readme): add JSON serialization example
test(money): add AsParts table-driven tests
```

---

Thank you for contributing!
