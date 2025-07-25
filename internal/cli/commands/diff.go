package commands

import (
	"fmt"
	"os"

	"ai_code_reviewer/internal/gitutil"

	"github.com/spf13/cobra"
)

type DiffOptions struct {
	SourceRef string
	TargetRef string
}

func CreateDiffCommand() *cobra.Command {
	opts := &DiffOptions{}

	cmd := &cobra.Command{
		Use:     "diff [args] |",
		Short:   "仅输出本地 git diff 内容",
		Args:    cobra.MaximumNArgs(2), // 允许 0-2 个位置参数
		Example: "  # 标志参数用法\n  diff --source master --target dev\n\n  # 位置参数用法\n  diff master dev\n\n  # 混合用法\n  diff master --target dev",
		Run:     runDiff(opts),
	}

	cmd.Flags().StringVarP(&opts.SourceRef, "source", "s", "", "源分支")
	cmd.Flags().StringVarP(&opts.TargetRef, "target", "t", "", "目标分支")

	return cmd
}

func runDiff(opts *DiffOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			opts.SourceRef = args[0]
		}
		if len(args) > 1 {
			opts.TargetRef = args[1]
		}

		diff, err := gitutil.GetGitDiff(opts.SourceRef, opts.TargetRef)
		if err != nil {
			fmt.Fprintln(os.Stderr, "获取 git diff 失败:", err)
			os.Exit(1)
		}
		if diff == "" {
			fmt.Println("无 diff 变更。")
			return
		}
		fmt.Println(diff)
	}
}
