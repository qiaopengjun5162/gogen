# gogen

![GitHub release (latest by date)](https://img.shields.io/github/v/release/yourusername/gogen)
![GitHub license](https://img.shields.io/github/license/yourusername/gogen)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue)

## 概览

`gogen` 是一个轻量级的项目生成工具，旨在从 Git 仓库或本地模板快速搭建新项目。它通过从 Git
URL（支持可选的分支指定）克隆或从本地路径复制来简化流程，并支持基本的模板变量替换。

### 功能

- 从 Git 仓库（例如 GitHub、GitLab）克隆项目模板。
- 从本地目录复制模板。
- 指定 Git 分支进行克隆（可选）。
- 将模板变量（例如 `{{project_name}}`）替换为用户提供的值。
- 交互式命令行界面，带有 ANSI 彩色输出和进度跟踪。
- 嵌入 Git 提交哈希以追踪版本。

### 快速开始

使用一条命令克隆模板并生成项目：

```bash
gogen --git=https://github.com/example/template
```

### 安装

1. **使用 Go**：
   ```bash
   go install github.com/yourusername/gogen@latest
   ```
2. **手动构建**：
   ```bash
   git clone https://github.com/yourusername/gogen.git
   cd gogen
   go build -v -ldflags "-X main.GitCommit=$(git rev-parse --short HEAD)" -o gogen ./main.go
   ```
   此命令会将 Git 提交哈希嵌入二进制文件中，运行工具时可见。

### 使用方法

从 Git 仓库生成项目：

```bash
gogen --git=https://github.com/example/template --branch=main
```

从本地模板生成项目：

```bash
gogen --local=/path/to/template
```

按照交互式提示输入项目名称（或接受从模板源推导出的默认名称）。

### 依赖要求

- **Go**：1.24 或更高版本（用于从源码构建）。
- **Git**：仅在克隆远程模板时需要。

### 示例输出

以下是使用 Git 仓库运行 `gogen` 的示例：

```bash
$ gogen --git=https://github.com/example/template
[INFO] Built from Git commit: abc1234
[INFO] 正在验证输入...
[INPUT] 输入项目名称（默认: template）: myproject
[INPUT] 从 https://github.com/example/template 生成项目 'myproject'？(Y/n): Y
[PROGRESS] 正在从 'https://github.com/example/template' 克隆 Git 仓库...
[SUCCESS] 项目 'myproject' 生成成功！

➜ gogen --git=https://github.com/qiaopengjun5162/gotcha
[INFO] 正在验证输入...
[INPUT] 输入项目名称（默认: gotcha）:
[INPUT] 从 https://github.com/qiaopengjun5162/gotcha 生成项目 'gotcha'？(Y/n):
[INFO] 正在生成项目 'gotcha'...
[PROGRESS] 正在从 'https://github.com/qiaopengjun5162/gotcha' 克隆 Git 仓库...
正克隆到 'gotcha'...
remote: Enumerating objects: 26, done.
remote: Counting objects: 100% (26/26), done.
remote: Compressing objects: 100% (22/22), done.
remote: Total 26 (delta 1), reused 22 (delta 1), pack-reused 0 (from 0)
接收对象中: 100% (26/26), 10.39 KiB | 5.20 MiB/s, 完成.
处理 delta 中: 100% (1/1), 完成.
[SUCCESS] 项目 'gotcha' 生成成功！

➜ gogen --git=https://github.com/qiaopengjun5162/gotcha
[INFO] 正在验证输入...
[INPUT] 输入项目名称（默认: gotcha）: myproject
[INPUT] 从 https://github.com/qiaopengjun5162/gotcha 生成项目 'myproject'？(Y/n): y
[INFO] 正在生成项目 'myproject'...
[PROGRESS] 正在从 'https://github.com/qiaopengjun5162/gotcha' 克隆 Git 仓库...
正克隆到 'myproject'...
remote: Enumerating objects: 26, done.
remote: Counting objects: 100% (26/26), done.
remote: Compressing objects: 100% (22/22), done.
remote: Total 26 (delta 1), reused 22 (delta 1), pack-reused 0 (from 0)
接收对象中: 100% (26/26), 10.39 KiB | 5.20 MiB/s, 完成.
处理 delta 中: 100% (1/1), 完成.
[SUCCESS] 项目 'myproject' 生成成功！
```

### 项目结构

```
gogen/
├── CHANGELOG.md    # 版本历史
├── LICENSE         # 项目许可证
├── Makefile        # 构建自动化
├── README.md       # 英文文档
├── README.zh.md    # 本文件（中文文档）
├── go.mod          # Go 模块文件
├── main.go         # 主源文件
```

### 贡献

我们欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 获取报告问题、提交拉取请求或改进项目的指南。

### 许可证

本项目采用 [MIT 许可证](LICENSE)。
