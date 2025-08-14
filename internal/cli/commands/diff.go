package commands

import (
	"fmt"
	"os"

	"ai_code_reviewer/internal/cli/progress"
	"ai_code_reviewer/internal/cli/renderer"
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

		// 初始化进度显示和渲染器
		progressTracker := progress.NewSimpleProgress("Git差异")
		renderer, err := renderer.NewRenderer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "初始化渲染器失败：%v\n", err)
			os.Exit(1)
		}

		// 获取Git diff
		progressTracker.Show("获取Git差异...")
		diff, err := gitutil.GetGitDiff(opts.SourceRef, opts.TargetRef)
		if err != nil {
			progressTracker.Error(fmt.Sprintf("获取 git diff 失败: %v", err))
			os.Exit(1)
		}

		if diff == "" {
			progressTracker.Info("无 diff 变更")
			return
		}

		progressTracker.Success("Git差异获取完成")
		renderer.RenderDiff(diff)
	}
}
