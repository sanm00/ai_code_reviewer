package commands

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"ai_code_reviewer/internal/config"

	"github.com/spf13/cobra"
)

type ConfigOptions struct {
	Print bool
	Set   []string
	Init  bool
}

func CreateConfigCommand() *cobra.Command {
	opts := &ConfigOptions{}

	cmd := &cobra.Command{
		Use:   "config",
		Short: "查看或设置配置文件",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.Init {
				if err := handleConfigInit(); err != nil {
					os.Exit(1)
				}
				return
			}

			if opts.Print {
				if err := handleConfigPrint(); err != nil {
					os.Exit(1)
				}
			}

			if len(opts.Set) > 0 {
				if err := handleConfigSet(opts.Set); err != nil {
					os.Exit(1)
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&opts.Print, "print", "p", false, "查看当前配置")
	cmd.Flags().StringArrayVarP(&opts.Set, "set", "s", nil, "设置配置项，如 -s key=value，可多次使用; 支持: token，prompt，model，url")
	cmd.Flags().BoolVarP(&opts.Init, "init", "i", false, "初始化配置文件（如果不存在则新建）")

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
