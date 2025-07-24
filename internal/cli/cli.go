package cli

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"ai_code_reviewer/internal/config"
	"ai_code_reviewer/internal/gitutil"
	"ai_code_reviewer/internal/openaiutil"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

const (
	appName    = "acr"
	appVersion = "v1.0.0"
)

// Run 启动 CLI
func Run() {
	rootCmd := createRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func createRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   appName,
		Short: "基于 OpenAI 的代码 diff 审查工具",
		Long:  "一个基于 OpenAI API 的本地 Git diff 代码审查命令行工具。",
		Run: func(cmd *cobra.Command, args []string) {
			if len(os.Args) == 1 {
				_ = cmd.Help()
				os.Exit(0)
			}
		},
	}

	rootCmd.Flags().BoolP("help", "h", false, "显示帮助信息")
	rootCmd.AddCommand(
		createDiffCommand(),
		createConfigCommand(),
		createReviewCommand(),
		createVersionCommand(),
	)

	return rootCmd
}

func createReviewCommand() *cobra.Command {
	var (
		sourceRef string
		targetRef string
	)

	cmd := &cobra.Command{
		Use:   "review",
		Short: "发送diff给AI审查",
		Run:   runReviewCommand(&sourceRef, &targetRef),
	}

	cmd.Flags().StringVarP(&sourceRef, "source", "s", "", "源分支")
	cmd.Flags().StringVarP(&targetRef, "target", "t", "", "目标分支")

	return cmd
}

func runReviewCommand(sourceRef, targetRef *string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(config.DefaultConfigFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取配置失败：%v\n", err)
			os.Exit(1)
		}

		diff, err := gitutil.GetGitDiff(*sourceRef, *targetRef)
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

func createDiffCommand() *cobra.Command {
	var (
		sourceRef string
		targetRef string
	)

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "仅输出本地 git diff 内容",
		Run: func(cmd *cobra.Command, args []string) {
			diff, err := gitutil.GetGitDiff(sourceRef, targetRef)
			if err != nil {
				fmt.Fprintln(os.Stderr, "获取 git diff 失败:", err)
				os.Exit(1)
			}
			if diff == "" {
				fmt.Println("无 diff 变更。")
				return
			}
			fmt.Println(diff)
		},
	}

	cmd.Flags().StringVarP(&sourceRef, "source", "s", "", "源分支")
	cmd.Flags().StringVarP(&targetRef, "target", "t", "", "目标分支")

	return cmd
}

func createConfigCommand() *cobra.Command {
	var (
		print bool
		set   []string
		init  bool
	)

	cmd := &cobra.Command{
		Use:   "config",
		Short: "查看或设置配置文件",
		Run: func(cmd *cobra.Command, args []string) {
			if init {
				if err := handleConfigInit(); err != nil {
					os.Exit(1)
				}
				return
			}

			if print {
				if err := handleConfigPrint(); err != nil {
					os.Exit(1)
				}
			}

			if len(set) > 0 {
				if err := handleConfigSet(set); err != nil {
					os.Exit(1)
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&print, "print", "p", false, "查看当前配置")
	cmd.Flags().StringArrayVarP(&set, "set", "s", nil, "设置配置项，如 -s key=value，可多次使用; 支持: token，prompt，model，url")
	cmd.Flags().BoolVarP(&init, "init", "i", false, "初始化配置文件（如果不存在则新建）")

	return cmd
}

func handleConfigInit() error {
	if err := config.InitConfigFile(config.DefaultConfigFile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	fmt.Println("已初始化配置文件：", config.DefaultConfigFile)
	return nil
}

func handleConfigPrint() error {
	cfg, err := config.LoadConfig(config.DefaultConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取配置失败：%v\n", err)
		return err
	}

	val := reflect.ValueOf(*cfg)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		fmt.Printf("%s: %v\n", field.Name, value.Interface())
	}
	return nil
}

func handleConfigSet(kvPairs []string) error {
	var updates config.Config

	for _, kv := range kvPairs {
		key, val, err := parseConfigKeyValue(kv)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

		switch key {
		case "token":
			updates.Token = val
		case "prompt":
			updates.Prompt = val
		case "model":
			updates.Model = val
		case "url":
			updates.Url = val
		default:
			fmt.Fprintf(os.Stderr, "不支持的配置项: %s\n", key)
			return fmt.Errorf("invalid config key")
		}
	}

	if err := config.UpdateConfigFile(config.DefaultConfigFile, updates); err != nil {
		fmt.Fprintf(os.Stderr, "写入配置失败: %v\n", err)
		return err
	}
	fmt.Println("配置已更新。")
	return nil
}

func parseConfigKeyValue(kv string) (string, string, error) {
	parts := strings.SplitN(kv, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("无效参数: %s，应为 key=value 格式", kv)
	}
	return parts[0], parts[1], nil
}

func createVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "查看版本",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version %s\n", appName, appVersion)
		},
	}
}
