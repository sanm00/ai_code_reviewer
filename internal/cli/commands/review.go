package commands

import (
	"fmt"
	"os"

	"ai_code_reviewer/internal/config"
	"ai_code_reviewer/internal/gitutil"
	"ai_code_reviewer/internal/openaiutil"

	"github.com/charmbracelet/glamour"
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

		cfg, err := config.LoadConfig(config.DefaultConfigFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取配置失败：%v\n", err)
			os.Exit(1)
		}

		diff, err := gitutil.GetGitDiff(opts.SourceRef, opts.TargetRef)
		if err != nil {
			fmt.Fprintln(os.Stderr, "获取 git diff 失败:", err)
			os.Exit(1)
		}
		if diff == "" {
			fmt.Println("无 diff 变更，无需审查。")
			return
		}

		result, err := openaiutil.Chart(cfg.Token, cfg.Prompt, diff, cfg.Model, cfg.Url)
		if err != nil {
			fmt.Fprintln(os.Stderr, "代码审查失败:", err)
			os.Exit(1)
		}

		if err := renderMarkdown(result); err != nil {
			fmt.Fprintln(os.Stderr, "输出结果失败:", err)
			os.Exit(1)
		}
	}
}

func renderMarkdown(content string) error {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(150),
	)
	if err != nil {
		return err
	}

	out, err := r.Render(content)
	if err != nil {
		return err
	}

	fmt.Print(out)
	return nil
}
