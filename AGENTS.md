# AGENTS.md

## Project Shape

- `gogen` is a Go CLI for generating projects from Git repositories or local template directories.
- The root package is the only implementation. Do not reintroduce duplicate split packages for the same CLI logic.
- `main.go` should stay as the thin entrypoint. Keep parsing, prompting, generation, processors, file operations, and logging in their existing module files.

## Commands

- Format: `gofmt -w *.go`
- Check: `make check`
- Test: `GOCACHE=/private/tmp/gogen-go-cache GOMODCACHE=/private/tmp/gogen-go-mod-cache go test ./...`
- Build: `GOCACHE=/private/tmp/gogen-go-cache GOMODCACHE=/private/tmp/gogen-go-mod-cache go build ./...`
- Release-style build: `go build -v -ldflags "-X main.GitCommit=$(git rev-parse --short HEAD)" -o gogen .`

## Production Standard

- Treat `gogen` as a production CLI, not a demo or throwaway scaffold.
- Every behavior change must have tests or a recorded reason why automated coverage is not practical.
- Review changes for correctness, error handling, test coverage, documentation drift, and build/release impact before committing.
- Any user-facing flag, output, error message, build command, release workflow, or template-processing behavior must be reflected in `README.md`, `README.zh.md`, `DEVLOG.md`, and `PRODUCTION.md` when relevant.
- Prefer small, complete changes that can be formatted, tested, committed, and pushed in the same session.
- Do not leave known broken workflows undocumented. If a local tool is missing, record the exact blocker and the manual verification used instead.
- Before commit, run the production gate in `PRODUCTION.md` or record which parts could not run and why.

## Implementation Notes

- Build the package with `go build .`, not `go build ./main.go`; the CLI is now modular.
- `--local` must preserve the provided template path in `Config.TemplateSrc`.
- `--name` sets `Config.ProjectName`; `--yes` and `-y` skip confirmation prompts.
- `--var key=value` can be repeated for custom template variables; `project_name` is reserved.
- `--version` must not require `--git` or `--local`.
- Local template copying intentionally skips `.git` directories.
- Template variable replacement intentionally skips `.git` directories and binary files.
- Keep tests focused on CLI behavior and filesystem effects before adding new features.

## Change Recording

- Update `DEVLOG.md` whenever code, documentation, configuration, build scripts, or release workflow files change.
- Each record should include the concrete files or areas changed, validation commands, problems encountered, and how they were resolved.
- If a commit uses `--no-verify`, record the exact hook blocker and the checks that were run manually.
