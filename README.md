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
- Replace template variables (e.g., `{{project_name}}`) with user-provided values.
- Interactive CLI with ANSI-colored output and progress tracking.
- Embedded Git commit hash for version tracking.

### Quick Start

Clone a template and generate a project in one command:

```bash
gogen --git=https://github.com/example/template
```

### Installation

1. **Using Go**:

   ```bash
   go install github.com/yourusername/gogen@latest
   ```

2. **Manual Build**:

   ```bash
   git clone https://github.com/yourusername/gogen.git
   cd gogen
   go build -v -ldflags "-X main.GitCommit=$(git rev-parse --short HEAD)" -o gogen ./main.go
   ```

   This embeds the Git commit hash into the binary, visible when running the tool.

### Usage

Generate a project from a Git repository:

```bash
gogen --git=https://github.com/example/template --branch=main
```

Generate a project from a local template:

```bash
gogen --local=/path/to/template
```

Follow the interactive prompts to specify a project name (or accept the default derived from the template source).

### Requirements

- **Go**: 1.24 or later (for building from source).
- **Git**: Required only for cloning remote templates.

### Example Output

Below is an example of running `gogen` with a Git repository:

```bash
$ gogen --git=https://github.com/example/template
[INFO] Built from Git commit: abc1234
[INFO] Validating input...
[INPUT] Enter project name (default: template): myproject
[INPUT] Generate project 'myproject' from https://github.com/example/template? (Y/n): Y
[PROGRESS] Cloning Git repository from 'https://github.com/example/template'...
[SUCCESS] Project 'myproject' generated successfully!

➜ gogen --git=https://github.com/qiaopengjun5162/gotcha
[INFO] Validating input...
[INPUT] Enter project name (default: gotcha):
[INPUT] Generate project 'gotcha' from https://github.com/qiaopengjun5162/gotcha? (Y/n):
[INFO] Generating project 'gotcha'...
[PROGRESS] Cloning Git repository from 'https://github.com/qiaopengjun5162/gotcha'...
正克隆到 'gotcha'...
remote: Enumerating objects: 26, done.
remote: Counting objects: 100% (26/26), done.
remote: Compressing objects: 100% (22/22), done.
remote: Total 26 (delta 1), reused 22 (delta 1), pack-reused 0 (from 0)
接收对象中: 100% (26/26), 10.39 KiB | 5.20 MiB/s, 完成.
处理 delta 中: 100% (1/1), 完成.
[SUCCESS] Project 'gotcha' generated successfully!

➜ gogen --git=https://github.com/qiaopengjun5162/gotcha
[INFO] Validating input...
[INPUT] Enter project name (default: gotcha): myproject
[INPUT] Generate project 'myproject' from https://github.com/qiaopengjun5162/gotcha? (Y/n): y
[INFO] Generating project 'myproject'...
[PROGRESS] Cloning Git repository from 'https://github.com/qiaopengjun5162/gotcha'...
正克隆到 'myproject'...
remote: Enumerating objects: 26, done.
remote: Counting objects: 100% (26/26), done.
remote: Compressing objects: 100% (22/22), done.
remote: Total 26 (delta 1), reused 22 (delta 1), pack-reused 0 (from 0)
接收对象中: 100% (26/26), 10.39 KiB | 5.20 MiB/s, 完成.
处理 delta 中: 100% (1/1), 完成.
[SUCCESS] Project 'myproject' generated successfully!

```

### Project Structure

```
gogen/
├── CHANGELOG.md    # Version history
├── LICENSE         # Project license
├── Makefile        # Build automation
├── README.md       # This file
├── README.zh.md    # Chinese documentation
├── go.mod          # Go module file
├── main.go         # Main source file
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
