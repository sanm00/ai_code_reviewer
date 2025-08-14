package root

import (
	"fmt"
	"os"

	"ai_code_reviewer/internal/cli/commands"

	"github.com/spf13/cobra"
)

const (
	VERSION string = "1.0.0"
	NAME    string = "acr"
)

// RootCommand 根命令结构
type RootCommand struct {
	cmd *cobra.Command
}

// NewRootCommand 创建新的根命令
func NewRootCommand() *RootCommand {
	root := &RootCommand{}
	root.cmd = root.createRootCommand()
	root.addSubCommands()
	return root
}

// createRootCommand 创建根命令
func (root *RootCommand) createRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   NAME,
		Short: "基于 OpenAI 的代码 diff 审查工具",
		Long: `一个基于 OpenAI API 的本地 Git diff 代码审查命令行工具。

主要功能：
  • review    - 发送diff给AI进行代码审查
  • diff      - 仅输出本地 git diff 内容
  • config    - 查看或设置配置文件
  • version   - 查看版本信息

使用示例：
  acr review master dev          # 审查从master到dev的变更
  acr diff --source main         # 查看与main分支的差异
  acr config --print             # 查看当前配置
  acr config --init              # 初始化配置文件`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(os.Args) == 1 {
				_ = cmd.Help()
				os.Exit(0)
			}
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
}

// addSubCommands 添加子命令
func (root *RootCommand) addSubCommands() {
	root.cmd.AddCommand(
		commands.CreateDiffCommand(),
		commands.CreateConfigCommand(),
		commands.CreateReviewCommand(),
		commands.CreateVersionCommand(NAME, VERSION),
	)
}

// Execute 执行根命令
func (root *RootCommand) Execute() error {
	return root.cmd.Execute()
}

// GetCommand 获取根命令
func (root *RootCommand) GetCommand() *cobra.Command {
	return root.cmd
}

// Run 启动CLI应用
func Run() {
	root := NewRootCommand()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
