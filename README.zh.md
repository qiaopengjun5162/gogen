# gogen

![GitHub release (latest by date)](https://img.shields.io/github/v/release/qiaopengjun5162/gogen)
![GitHub license](https://img.shields.io/github/license/qiaopengjun5162/gogen)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue)

## 概览

`gogen` 是一个轻量级的项目生成工具，旨在从 Git 仓库或本地模板快速搭建新项目。它通过从 Git
URL（支持可选的分支指定）克隆或从本地路径复制来简化流程，并支持基本的模板变量替换。

### 功能

- 从 Git 仓库（例如 GitHub、GitLab）克隆项目模板。
- 从本地目录复制模板。
- 指定 Git 分支进行克隆（可选）。
- 替换文件内容、文件名和目录名中的模板变量。
- 交互式命令行界面，带有 ANSI 彩色输出和进度跟踪。
- 支持指定项目名并跳过确认提示，便于脚本和 CI 使用。
- 通过版本命令输出嵌入的 Git 提交信息。

### 快速开始

使用一条命令克隆模板并生成项目：

```bash
gogen --git=https://github.com/example/template --name=myproject --yes
```

### 安装

1. **使用 Go**：
   ```bash
   go install github.com/qiaopengjun5162/gogen@latest
   ```
2. **手动构建**：
   ```bash
   git clone https://github.com/qiaopengjun5162/gogen.git
   cd gogen
   go build -v -ldflags "-X main.GitCommit=$(git rev-parse --short HEAD)" -o gogen .
   ```
   此命令会将 Git 提交哈希嵌入二进制文件中，可通过 `gogen --version` 查看。

### 使用方法

从 Git 仓库生成项目：

```bash
gogen --git=https://github.com/example/template --branch=main
```

从本地模板生成项目：

```bash
gogen --local=/path/to/template
```

不进入交互提示直接生成：

```bash
gogen --local=/path/to/template --name=myproject --yes
```

传入自定义模板变量：

```bash
gogen --local=/path/to/template --name=myproject --var module=github.com/example/myproject --var license=MIT --yes
```

模板变量会同时替换文件内容和路径。例如，模板文件 `{{project_name}}/{{module_name}}.txt`
可以生成 `myproject/core.txt`。

查看构建信息：

```bash
gogen --version
```

参数：

| 参数 | 说明 |
| --- | --- |
| `--git` | 模板 Git 仓库 URL。 |
| `--local` | 本地模板目录路径。 |
| `--branch` | 要克隆的 Git 分支；仅可与 `--git` 一起使用。 |
| `--name` | 项目名，也是目标目录名。 |
| `--var` | `key=value` 形式的模板变量，可重复传入。 |
| `--yes`, `-y` | 跳过确认提示。 |
| `--version` | 输出版本信息后退出。 |

未传 `--name` 时，`gogen` 会从模板来源推导默认项目名。
`project_name` 是保留变量，始终来自项目名。

### 依赖要求

- **Go**：1.24 或更高版本（用于从源码构建）。
- **Git**：仅在克隆远程模板时需要。

### 示例输出

以下是使用 Git 仓库运行 `gogen` 的示例：

```bash
$ gogen --git=https://github.com/example/template --name=myproject --yes
[INFO] Validating input...
[INFO] Generating project 'myproject'...
[PROGRESS] Cloning Git repository from 'https://github.com/example/template'...
[SUCCESS] Project 'myproject' generated successfully!
```

### 项目结构

```
gogen/
├── CHANGELOG.md    # 版本历史
├── LICENSE         # 项目许可证
├── Makefile        # justfile 的兼容包装入口
├── README.md       # 英文文档
├── README.zh.md    # 本文件（中文文档）
├── config.go       # CLI 配置和校验规则
├── files.go        # 本地模板复制和变量替换
├── generator.go    # 项目生成编排
├── go.mod          # Go 模块文件
├── input.go        # 交互式输入处理
├── justfile        # 标准开发任务入口
├── main.go         # CLI 入口
├── processor.go    # Git 和本地模板处理器
├── logger.go       # CLI 输出辅助函数
├── main_test.go    # 行为测试
```

### 贡献

我们欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 获取报告问题、提交拉取请求或改进项目的指南。

### 许可证

本项目采用 [MIT 许可证](LICENSE)。
