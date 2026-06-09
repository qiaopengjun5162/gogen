# gogen
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

![GitHub release (latest by date)](https://img.shields.io/github/v/release/qiaopengjun5162/gogen)
![GitHub license](https://img.shields.io/github/license/qiaopengjun5162/gogen)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue)

## Chinese Documentation

中文文档请参阅 [README.zh.md](README.zh.md)。

## Overview

`gogen` is a lightweight project generation tool designed to scaffold new projects from Git repositories or local
templates. It streamlines the process by cloning from a Git URL (with optional branch specification) or copying from a
local path, while supporting basic template variable substitution.

### Features

- Clone project templates from Git repositories (e.g., GitHub, GitLab).
- Copy templates from local directories.
- Specify a Git branch for cloning (optional).
- Replace template variables in file contents, file names, and directory names.
- Interactive CLI with ANSI-colored output and progress tracking.
- Non-interactive generation with explicit project names and confirmation skipping.
- Version output with embedded Git commit metadata.

### Quick Start

Clone a template and generate a project in one command:

```bash
gogen --git=https://github.com/example/template --name=myproject --yes
```

### Installation

1. **Using Go**:

   ```bash
   go install github.com/qiaopengjun5162/gogen@latest
   ```

2. **Manual Build**:

   ```bash
   git clone https://github.com/qiaopengjun5162/gogen.git
   cd gogen
   go build -v -ldflags "-X main.GitCommit=$(git rev-parse --short HEAD)" -o gogen .
   ```

   This embeds the Git commit hash into the binary, visible through `gogen --version`.

### Usage

Generate a project from a Git repository:

```bash
gogen --git=https://github.com/example/template --branch=main
```

Generate a project from a local template:

```bash
gogen --local=/path/to/template
```

Generate without interactive prompts:

```bash
gogen --local=/path/to/template --name=myproject --yes
```

Pass custom template variables:

```bash
gogen --local=/path/to/template --name=myproject --var module=github.com/example/myproject --var license=MIT --yes
```

Template variables are replaced in both file contents and paths. For example, a template file named
`{{project_name}}/{{module_name}}.txt` can generate `myproject/core.txt`.

Print build information:

```bash
gogen --version
```

Options:

| Flag | Description |
| --- | --- |
| `--git` | Git repository URL for the template. |
| `--local` | Local template directory path. |
| `--branch` | Git branch to clone; only valid with `--git`. |
| `--name` | Project name and destination directory. |
| `--var` | Template variable in `key=value` form; can be repeated. |
| `--yes`, `-y` | Skip confirmation prompts. |
| `--version` | Print version information and exit. |

Without `--name`, `gogen` derives the default project name from the template source.
`project_name` is reserved and always comes from the project name.

### Requirements

- **Go**: 1.24 or later (for building from source).
- **Git**: Required only for cloning remote templates.

### Example Output

Below is an example of running `gogen` with a Git repository:

```bash
$ gogen --git=https://github.com/example/template --name=myproject --yes
[INFO] Validating input...
[INFO] Generating project 'myproject'...
[PROGRESS] Cloning Git repository from 'https://github.com/example/template'...
[SUCCESS] Project 'myproject' generated successfully!
```

### Project Structure

```
gogen/
├── CHANGELOG.md    # Version history
├── LICENSE         # Project license
├── Makefile        # Compatibility wrapper for justfile
├── README.md       # This file
├── README.zh.md    # Chinese documentation
├── config.go       # CLI configuration and validation patterns
├── files.go        # Local template copy and variable replacement
├── generator.go    # Project generation orchestration
├── go.mod          # Go module file
├── input.go        # Interactive prompt handling
├── justfile        # Canonical development tasks
├── main.go         # CLI entrypoint
├── processor.go    # Git and local template processors
├── logger.go       # CLI output helpers
├── main_test.go    # Behavior tests
```

### Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to report issues, submit
pull requests, or improve the project.

## Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/qiaopengjun5162"><img src="https://avatars.githubusercontent.com/u/124650229?v=4?s=100" width="100px;" alt="Paxon Qiao 乔鹏军"/><br /><sub><b>Paxon Qiao 乔鹏军</b></sub></a><br /><a href="#content-qiaopengjun5162" title="Content">🖋</a></td>
    </tr>
  </tbody>
  <tfoot>
    <tr>
      <td align="center" size="13px" colspan="7">
        <img src="https://raw.githubusercontent.com/all-contributors/all-contributors-cli/1b8533af435da9854653492b1327a23a4dbd0a10/assets/logo-small.svg">
          <a href="https://all-contributors.js.org/docs/en/bot/usage">Add your contributions</a>
        </img>
      </td>
    </tr>
  </tfoot>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

### License

This project is licensed under the [MIT License](LICENSE).
