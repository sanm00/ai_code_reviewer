package commands

import (
	"fmt"
	"os"

	"ai_code_reviewer/internal/cli/progress"
	"ai_code_reviewer/internal/cli/renderer"
	"ai_code_reviewer/internal/config"
	"ai_code_reviewer/internal/gitutil"
	"ai_code_reviewer/internal/openaiutil"

	"github.com/spf13/cobra"
)

type ReviewOptions struct {
	SourceRef string
	TargetRef string
}

func CreateReviewCommand() *cobra.Command {
	opts := &ReviewOptions{}

	cmd := &cobra.Command{
		Use:     "review [args] |",
		Short:   "发送diff给AI审查",
		Args:    cobra.MaximumNArgs(2), // 允许 0-2 个位置参数
		Example: "  # 标志参数用法\n  review --source master --target dev\n\n  # 位置参数用法\n  review master dev\n\n  # 混合用法\n  review master --target dev",
		Run:     runReview(opts),
	}

	cmd.Flags().StringVarP(&opts.SourceRef, "source", "s", "", "源分支")
	cmd.Flags().StringVarP(&opts.TargetRef, "target", "t", "", "目标分支")

	return cmd
}

func runReview(opts *ReviewOptions) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			opts.SourceRef = args[0]
		}
		if len(args) > 1 {
			opts.TargetRef = args[1]
		}

		// 初始化进度显示和渲染器
		progressTracker := progress.NewSimpleProgress("")
		renderer, err := renderer.NewRenderer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "初始化渲染器失败：%v\n", err)
			os.Exit(1)
		}

		// 加载配置
		progressTracker.Show("加载配置...")
		cfg, err := config.LoadConfig(config.DefaultConfigFile)
		if err != nil {
			progressTracker.Error(fmt.Sprintf("获取配置失败：%v", err))
			os.Exit(1)
		}
		progressTracker.Success("配置加载完成")

		// 获取Git diff
		progressTracker.Show("获取Git差异...")
		diff, err := gitutil.GetGitDiff(opts.SourceRef, opts.TargetRef)
		if err != nil {
			progressTracker.Error(fmt.Sprintf("获取 git diff 失败: %v", err))
			os.Exit(1)
		}
		if diff == "" {
			progressTracker.Info("无 diff 变更，无需审查")
			return
		}
		progressTracker.Success("Git差异获取完成")

		// 发送给AI审查
		progressTracker.Show("发送给AI进行代码审查...")
		spinner := progress.NewSpinner("AI正在分析代码")
		spinner.Start()

		result, err := openaiutil.Chart(cfg.Token, cfg.Prompt, diff, cfg.Model, cfg.Url)
		spinner.Stop()

		if err != nil {
			progressTracker.Error(fmt.Sprintf("代码审查失败: %v", err))
			os.Exit(1)
		}
		progressTracker.Success("AI代码审查完成")

		// 渲染结果
		progressTracker.Show("渲染审查结果...")
		if err := renderer.RenderMarkdown(result); err != nil {
			progressTracker.Error(fmt.Sprintf("输出结果失败: %v", err))
			os.Exit(1)
		}
		progressTracker.Success("审查结果渲染完成")
	}
}
