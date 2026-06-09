# Production Checklist

`gogen` should be maintained as a production-grade CLI. A change is not done until behavior, validation, documentation, and repository state are all handled.

## Required Gate

Run these before committing code changes:

```bash
just check
```

For CLI behavior changes, also run a real binary smoke test:

```bash
go build -o gogen .
./gogen --version
```

When template generation changes, verify at least one local template generation path with `--local`, `--name`, and `--yes`.

## Documentation Gate

Update these files when relevant:

- `README.md` for English user-facing behavior.
- `README.zh.md` for Chinese user-facing behavior.
- `CHANGELOG.md` for unreleased user-visible changes.
- `DEVLOG.md` for concrete work records, validation, blockers, and resolutions.
- `AGENTS.md` for durable project rules or non-obvious implementation constraints.
- `.github/workflows/build.yml`, `justfile`, `Makefile`, or `.pre-commit-config.yaml` when commands or release flow change.

## Behavior Standards

- CLI flags must be deterministic and script-friendly.
- Interactive prompts must have a non-interactive alternative when practical.
- Error messages should identify the failed input or operation.
- Template processing must avoid corrupting binary files or repository metadata.
- `project_name` is a reserved template variable and must always come from the validated project name.
- Custom template variables must use validated keys and structured `--var key=value` parsing.
- Git operations should pass arguments as structured command arguments, never as shell-concatenated strings.
- Generated project cleanup must leave no partial destination directory after failure.

## Code Review Standards

- Review each change for correctness, error handling, test coverage, user-facing documentation, and release/build impact before committing.
- Prefer targeted refactors that remove duplication or clarify ownership boundaries.
- Do not merge behavior changes that only pass manual smoke tests when they can reasonably be covered by automated tests.
- Keep command construction structured and filesystem writes explicit.
- Treat documentation drift as a product defect.

## Repository Hygiene

- Keep generated binaries ignored and out of commits.
- Do not keep duplicate implementations of the same CLI behavior.
- Keep `justfile` as the canonical local task runner; `Makefile` is a compatibility wrapper.
- Keep Go task commands isolated from stale shell toolchain settings such as an old `GOROOT`.
- Keep `main` protected: require pull requests, require `build-go (ubuntu-latest)`, forbid force pushes, forbid deletions, require linear history, and require resolved conversations.
- Commit only coherent, reviewable changes.
- Push completed commits when the user asks to continue shared work on the remote branch.
- If `--no-verify` is necessary, record the hook failure and manual checks in `DEVLOG.md`.
