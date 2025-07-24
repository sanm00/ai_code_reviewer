# ai_code_reviewer

基于 OpenAI 的本地 Git 代码 diff 审查 CLI 工具

## 功能简介

- 获取本地 Git 仓库的 diff 变更
- 调用 OpenAI 大模型 API 自动审查代码变更，输出潜在问题和建议
- 支持自定义审查提示词（prompt）
- 支持 API Token 一次性配置，安全便捷
- 支持命令行参数、环境变量和配置文件多种配置方式

## 安装方法

1. **克隆项目**

```bash
git clone <your_repo_url>
cd ai_code_reviewer
```

2. **拉取依赖**

```bash
go mod tidy
```

3. **编译可执行文件**

```bash
go build -o acr
```

## 配置说明

### 方式一：配置文件（推荐）

在用户主目录下创建 `~/.acr/config.yaml`，内容示例：

```yaml
token: sk-xxxxxxx   # 你的 OpenAI API Token
prompt: 请帮我审查以下代码变更，指出潜在问题并给出建议
model: gpt-3.5-turbo
```

### 方式二：环境变量

```bash
export AI_CODE_REVIEWER_TOKEN=sk-xxxxxxx
export AI_CODE_REVIEWER_PROMPT="请帮我审查以下代码变更，指出潜在问题并给出建议"
export AI_CODE_REVIEWER_MODEL=gpt-3.5-turbo
```

## 使用示例
```
acr -h
```

## 常见问题

1. **提示 token 未配置？**
   - 请确保已在配置文件、环境变量或命令行参数中正确设置 OpenAI API Token。

2. **无 diff 变更，无需审查？**
   - 当前工作区没有未提交的代码变更。

3. **OpenAI API 调用失败？**
   - 检查网络、token 是否有效、模型名称是否正确。

## 依赖说明

- Go 1.23 及以上
- [github.com/sashabaranov/go-openai](https://github.com/sashabaranov/go-openai)
- [github.com/spf13/cobra](https://github.com/spf13/cobra)
- [github.com/spf13/viper](https://github.com/spf13/viper)

---

如有问题或建议，欢迎提 issue！ 