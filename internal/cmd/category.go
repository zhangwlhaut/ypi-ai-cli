package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zhangwlhaut/ypi-ai-cli/internal/model"
	"github.com/zhangwlhaut/ypi-ai-cli/internal/output"
)

var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "接口分类操作",
	Aliases: []string{"cat"},
}

var categoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取接口分类列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		cats, err := apiClient.GetCatMenu()
		if err != nil {
			output.PrintError(os.Stderr, "api_error", err.Error(), "")
			os.Exit(output.ExitAPIError)
		}

		if useJSONOutput() {
			return output.PrintJSON(os.Stdout, cats, !flagPretty)
		}
		// Compact list: just id + name + count
		type catSummary struct {
			ID       int    `json:"_id"`
			Name     string `json:"name"`
			Count    int    `json:"count"`
		}
		var summaries []catSummary
		for _, cat := range cats {
			summaries = append(summaries, catSummary{
				ID:    cat.ID,
				Name:  cat.Name,
				Count: len(cat.List),
			})
		}
		return output.PrintJSON(os.Stdout, summaries, true)
	},
}

var (
	catName string
	catDesc string
)

var categoryAddCmd = &cobra.Command{
	Use:   "add --name <name> [--desc <desc>]",
	Short: "新增接口分类",
	RunE: func(cmd *cobra.Command, args []string) error {
		if catName == "" {
			output.PrintError(os.Stderr, "param_required", "--name is required", "")
			os.Exit(output.ExitParamError)
		}

		req := &model.AddCatRequest{
			Name: catName,
			Desc: catDesc,
		}
		result, err := apiClient.AddCat(req)
		if err != nil {
			output.PrintError(os.Stderr, "api_error", err.Error(), "")
			os.Exit(output.ExitAPIError)
		}

		fmt.Fprintf(os.Stderr, "Category created\n")
		return output.PrintJSON(os.Stdout, json.RawMessage(result), !flagPretty)
	},
}

func init() {
	categoryAddCmd.Flags().StringVar(&catName, "name", "", "分类名称 (required)")
	categoryAddCmd.Flags().StringVar(&catDesc, "desc", "", "分类描述")

	categoryCmd.AddCommand(categoryListCmd)
	categoryCmd.AddCommand(categoryAddCmd)
	rootCmd.AddCommand(categoryCmd)
}
