# AI Code Reviewer (ACR)

基于 OpenAI 的本地 Git 代码 diff 审查 CLI 工具

## ✨ 功能特性

- 🔍 **智能代码审查**: 获取本地 Git 仓库的 diff 变更，调用 OpenAI 大模型 API 自动审查代码
- 📊 **实时进度显示**: 支持进度条、旋转指示器和状态消息，提供良好的用户体验
- 🎨 **美观输出**: 支持 Markdown 渲染，输出格式化的审查结果
- ⚙️ **灵活配置**: 支持命令行参数、环境变量和配置文件多种配置方式
- 🔒 **安全存储**: API Token 本地加密存储，安全便捷
- 🚀 **高性能**: 模块化架构，支持并发处理和进度跟踪

## 📦 项目结构

```
ai_code_reviewer/
├── bin/                    # 编译输出目录
├── internal/
│   ├── cli/               # CLI 相关模块
│   │   ├── commands/      # 命令实现
│   │   │   ├── review.go  # 代码审查命令
│   │   │   ├── diff.go    # 差异查看命令
│   │   │   ├── config.go  # 配置管理命令
│   │   │   └── version.go # 版本信息命令
│   │   ├── progress/      # 进度显示模块
│   │   │   └── progress.go # 进度条、旋转指示器等
│   │   ├── renderer/      # 输出渲染模块
│   │   │   └── renderer.go # Markdown渲染、格式化输出
│   │   ├── root/          # 根命令管理
│   │   │   └── root.go    # CLI根命令和子命令管理
│   │   └── cli.go         # CLI入口
│   ├── config/            # 配置管理
│   │   └── config.go      # 配置文件读写
│   ├── gitutil/           # Git工具
│   │   └── git.go         # Git diff获取
│   └── openaiutil/        # OpenAI工具
│       └── openai.go      # OpenAI API调用
├── main.go                # 程序入口
├── go.mod                 # Go模块文件
├── go.sum                 # 依赖校验文件
└── README.md              # 项目说明
```

## 🚀 安装方法

### 1. 克隆项目

```bash
git clone <your_repo_url>
cd ai_code_reviewer
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 编译可执行文件

```bash
go build -o bin/acr main.go
```

### 4. 添加到PATH（可选）

```bash
# macOS/Linux
sudo cp bin/acr /usr/local/bin/

# 或者添加到 ~/.bashrc 或 ~/.zshrc
export PATH="$PATH:$(pwd)/bin"
```

## ⚙️ 配置说明

### 方式一：配置文件（推荐）

使用 `config` 命令初始化配置文件：

```bash
# 初始化配置文件
acr config --init

# 设置配置项
acr config --set token=sk-your-openai-token
acr config --set model=gpt-3.5-turbo
acr config --set prompt="请帮我审查以下代码变更，指出潜在问题并给出改进建议"

# 查看当前配置
acr config --print
```

配置文件位置：`~/.acr/config.yaml`


## 📖 使用指南

### 基本命令

```bash
# 查看帮助
acr -h

# 查看版本
acr version

# 查看所有命令
acr --help
```

### 代码审查

```bash
# 审查当前工作区的未提交变更
acr review

# 审查指定分支间的差异
acr review master dev

# 使用标志参数
acr review --source master --target dev

# 混合使用位置参数和标志参数
acr review master --target dev
```

### 查看差异

```bash
# 查看当前工作区的未提交变更
acr diff

# 查看与指定分支的差异
acr diff master

# 查看两个分支间的差异
acr diff master dev

# 使用标志参数
acr diff --source main --target feature-branch
```

### 配置管理

```bash
# 初始化配置文件
acr config --init

# 查看当前配置
acr config --print

# 设置单个配置项
acr config --set token=sk-your-token

# 设置多个配置项
acr config --set token=sk-your-token --set model=gpt-4 --set prompt="自定义提示词"
```

## 🎯 使用示例

### 示例1：审查功能分支

```bash
# 切换到功能分支
git checkout feature/new-feature

# 审查与主分支的差异
acr review main

# 输出示例：
# ℹ️  代码审查: 加载配置...
# ✅ 代码审查: 配置加载完成
# ℹ️  代码审查: 获取Git差异...
# ✅ 代码审查: Git差异获取完成
# ℹ️  代码审查: 发送给AI进行代码审查...
# ⠋ AI正在分析代码
# ✅ 代码审查: AI代码审查完成
# ℹ️  代码审查: 渲染审查结果...
# ✅ 代码审查: 审查结果渲染完成
```

### 示例2：查看差异内容

```bash
# 查看与远程主分支的差异
acr diff origin/main

# 输出示例：
# ℹ️  Git差异: 获取Git差异...
# ✅ Git差异: Git差异获取完成
# diff --git a/src/main.go b/src/main.go
# index abc123..def456 100644
# --- a/src/main.go
# +++ b/src/main.go
# @@ -10,6 +10,7 @@ func main() {
#      fmt.Println("Hello, World!")
# +    fmt.Println("New feature added!")
#  }
```

### 示例3：配置管理

```bash
# 初始化配置
acr config --init

# 设置API Token
acr config --set token=sk-your-openai-api-token

# 自定义审查提示词
acr config --set prompt="请从代码质量、安全性、性能等方面审查以下代码变更"

# 查看配置
acr config --print

# 输出示例：
# ℹ️  配置查看: 加载配置文件...
# ✅ 配置查看: 配置加载完成
# Token: sk-your-openai-api-token
# Prompt: 请从代码质量、安全性、性能等方面审查以下代码变更
# Model: gpt-3.5-turbo
# Url: https://api.openai.com/v1
```

## 🔧 进度显示功能

ACR 提供了多种进度显示方式：

### 1. 简单进度显示
- ✅ 成功消息
- ❌ 错误消息  
- ℹ️ 信息消息
- ⚠️ 警告消息

### 2. 旋转指示器
- 用于长时间运行的操作（如AI分析）
- 显示动态旋转动画

### 3. 进度条（预留功能）
- 支持百分比显示
- 显示预计完成时间
- 适用于多步骤操作

## 🛠️ 开发说明

### 依赖要求

- Go 1.21 及以上
- [github.com/sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) - OpenAI API 客户端
- [github.com/spf13/cobra](https://github.com/spf13/cobra) - CLI 框架
- [github.com/spf13/viper](https://github.com/spf13/viper) - 配置管理
- [github.com/charmbracelet/glamour](https://github.com/charmbracelet/glamour) - Markdown 渲染

### 架构设计

项目采用模块化设计，主要模块包括：

- **CLI模块**: 负责命令行界面和用户交互
- **进度显示模块**: 提供多种进度反馈方式
- **渲染模块**: 负责输出格式化和美化
- **配置模块**: 管理应用配置
- **Git工具模块**: 处理Git相关操作
- **OpenAI工具模块**: 处理AI API调用

## ❓ 常见问题

### 1. 提示 token 未配置？
- 请确保已在配置文件、环境变量或命令行参数中正确设置 OpenAI API Token
- 使用 `acr config --init` 初始化配置文件
- 使用 `acr config --set token=your-token` 设置Token

### 2. 无 diff 变更，无需审查？
- 当前工作区没有未提交的代码变更
- 请先进行代码修改或切换到有差异的分支

### 3. OpenAI API 调用失败？
- 检查网络连接是否正常
- 确认API Token是否有效
- 检查模型名称是否正确
- 确认API配额是否充足

### 4. 进度显示不工作？
- 确保终端支持ANSI转义序列
- 在Windows上可能需要使用支持ANSI的终端

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

---

**享受智能代码审查的便利！** 🚀 