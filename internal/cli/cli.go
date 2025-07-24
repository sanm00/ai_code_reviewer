// cli/cli.go
package cli

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

// Run 启动 CLI
func Run() {
	rootCmd := &cobra.Command{
		Use:   "acr",
		Short: "基于 OpenAI 的代码 diff 审查工具",
		Long:  "一个基于 OpenAI API 的本地 Git diff 代码审查命令行工具。",
		Run: func(cmd *cobra.Command, args []string) {
			if len(os.Args) == 1 {
				_ = cmd.Help()
				os.Exit(0)
			}
		},
	}

	rootCmd.AddCommand(
		commands.CreateDiffCommand(),
		commands.CreateConfigCommand(),
		commands.CreateReviewCommand(),
		commands.CreateVersionCommand(NAME, VERSION),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
