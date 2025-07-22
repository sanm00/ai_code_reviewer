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

var (
	flagPrompt string
	flagModel  string
)

// Run 启动 CLI
func Run() {
	var rootCmd = &cobra.Command{
		Use:   "acr",
		Short: "基于 OpenAI 的代码 diff 审查工具",
		Long:  "一个基于 OpenAI API 的本地 Git diff 代码审查命令行工具。",
		Run: func(cmd *cobra.Command, args []string) {
			// 如果没有传递任何 flag 或参数，输出帮助信息
			if len(os.Args) == 1 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	// Review 子命令
	var (
		reviewSourceRef string
		reviewTargetRef string
	)
	var reviewCmd = &cobra.Command{
		Use:   "review",
		Short: "发送diff ai 审查",
		Run: func(cmd *cobra.Command, args []string) {
			// 获取配置
			cfg, err := config.LoadConfig(config.DefaultConfigFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "获取配置失败：%v/n", err)
			}

			// 获取 git diff
			diff, err := gitutil.GetGitDiff(reviewSourceRef, reviewTargetRef)
			if err != nil {
				fmt.Fprintln(os.Stderr, "获取 git diff 失败:", err)
				os.Exit(1)
			}
			if diff == "" {
				fmt.Println("无 diff 变更，无需审查。")
				return
			}

			// 调用 OpenAI 审查
			result, err := openaiutil.Chart(cfg.Token, cfg.Prompt, diff, cfg.Model, cfg.Url)
			if err != nil {
				fmt.Fprintln(os.Stderr, "代码审查失败:", err)
				os.Exit(1)
			}

			// markdown 渲染
			r, _ := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(150),
			)

			out, err := r.Render(result)
			if err != nil {
				fmt.Fprintln(os.Stderr, "输出结果失败:", err)
				os.Exit(1)
			}
			fmt.Print(out)

		},
	}

	reviewCmd.Flags().StringVarP(&reviewSourceRef, "source", "s", "", "源分支")
	reviewCmd.Flags().StringVarP(&reviewTargetRef, "target", "t", "", "目标分支")

	var (
		diffSourceRef string
		diffTargetRef string
	)
	var diffCmd = &cobra.Command{
		Use:   "diff",
		Short: "仅输出本地 git diff 内容",
		Run: func(cmd *cobra.Command, args []string) {
			diff, err := gitutil.GetGitDiff(diffSourceRef, diffTargetRef)
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

	diffCmd.Flags().StringVarP(&diffSourceRef, "source", "s", "", "源分支")
	diffCmd.Flags().StringVarP(&diffTargetRef, "target", "t", "", "目标分支")

	// config 子命令
	var (
		cfgPrint bool
		cfgSet   []string
		cfgInit  bool
	)
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "查看或设置配置文件",
		Run: func(cmd *cobra.Command, args []string) {
			if cfgInit {
				if err := config.InitConfigFile(config.DefaultConfigFile); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				fmt.Println("已初始化配置文件：", config.DefaultConfigFile)
				return
			}

			if cfgPrint {
				cfg, err := config.LoadConfig(config.DefaultConfigFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "获取配置失败：%v/n", err)
				}

				val := reflect.ValueOf(*cfg)
				typ := val.Type()

				for i := 0; i < val.NumField(); i++ {
					field := typ.Field(i)
					value := val.Field(i)
					fmt.Printf("%s: %v\n", field.Name, value.Interface())
				}

			}

			if len(cfgSet) > 0 {
				var updates config.Config

				for _, kv := range cfgSet {
					parts := strings.SplitN(kv, "=", 2)
					if len(parts) != 2 {
						fmt.Fprintf(os.Stderr, "无效参数: %s，应为 key=value 格式\n", kv)
						continue
					}
					key := parts[0]
					val := parts[1]
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
					}
				}

				if err := config.UpdateConfigFile(config.DefaultConfigFile, updates); err != nil {
					fmt.Fprintf(os.Stderr, "写入配置失败: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("配置已更新。")
			}
		},
	}
	configCmd.Flags().BoolVarP(&cfgPrint, "print", "p", false, "查看当前配置")
	configCmd.Flags().StringArrayVarP(&cfgSet, "set", "s", nil, "设置配置项，如 -s key=value，可多次使用; 支持: token，prompt，model，url")
	configCmd.Flags().BoolVarP(&cfgInit, "init", "i", false, "初始化配置文件（如果不存在则新建）")

	rootCmd.Flags().BoolP("help", "h", false, "显示帮助信息")

	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(reviewCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
