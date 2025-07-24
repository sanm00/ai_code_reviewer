package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CreateVersionCommand(name, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "查看版本",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version %s\n", name, version)
		},
	}
}
