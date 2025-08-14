package commands

import (
	"ai_code_reviewer/internal/cli/renderer"

	"github.com/spf13/cobra"
)

func CreateVersionCommand(name, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "查看版本",
		Run: func(cmd *cobra.Command, args []string) {
			renderer, err := renderer.NewRenderer()
			if err != nil {
				// 如果渲染器初始化失败，使用简单的输出
				cmd.Printf("%s version %s\n", name, version)
				return
			}
			renderer.RenderVersion(name, version)
		},
	}
}
