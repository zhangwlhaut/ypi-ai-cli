package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zhangwlhaut/yapi-ai-cli/internal/output"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "项目操作",
}

var projectShowCmd = &cobra.Command{
	Use:   "show",
	Short: "获取项目基本信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		project, err := apiClient.GetProject()
		if err != nil {
			output.PrintError(os.Stderr, "api_error", err.Error(), "")
			os.Exit(output.ExitAPIError)
		}

		if useJSONOutput() {
			return output.PrintJSON(os.Stdout, project, !flagPretty)
		}
		return output.PrintProjectMarkdown(os.Stdout, project)
	},
}

func init() {
	projectCmd.AddCommand(projectShowCmd)
	rootCmd.AddCommand(projectCmd)
}
