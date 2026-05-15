package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zhangwlhaut/yapi-ai-cli/internal/model"
	"github.com/zhangwlhaut/yapi-ai-cli/internal/output"
)

var (
	importType  string
	importMerge string
	importJSON  string
	importURL   string
)

var importCmd = &cobra.Command{
	Use:   "import --type <type> --merge <mode> [--json <data>] [--url <url>]",
	Short: "服务端数据导入（swagger/json）",
	Long: `导入接口数据到 YApi 项目。
支持类型: swagger, json
同步模式: normal(普通模式), good(智能合并), merge(完全覆盖)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if importType == "" {
			output.PrintError(os.Stderr, "param_required", "--type is required (swagger/json)", "")
			os.Exit(output.ExitParamError)
		}
		if importMerge == "" {
			output.PrintError(os.Stderr, "param_required", "--merge is required (normal/good/merge)", "")
			os.Exit(output.ExitParamError)
		}

		req := &model.ImportRequest{
			Type:  importType,
			Merge: importMerge,
			JSON:  importJSON,
			URL:   importURL,
		}

		result, err := apiClient.ImportData(req)
		if err != nil {
			output.PrintError(os.Stderr, "api_error", err.Error(), "")
			os.Exit(output.ExitAPIError)
		}

		fmt.Fprintf(os.Stderr, "Import completed\n")
		return output.PrintJSON(os.Stdout, json.RawMessage(result), !flagPretty)
	},
}

func init() {
	importCmd.Flags().StringVar(&importType, "type", "", "导入类型: swagger, json (required)")
	importCmd.Flags().StringVar(&importMerge, "merge", "", "同步模式: normal, good, merge (required)")
	importCmd.Flags().StringVar(&importJSON, "json", "", "JSON 数据（序列化后的字符串）")
	importCmd.Flags().StringVar(&importURL, "url", "", "导入数据 URL（存在时通过 URL 获取数据）")

	rootCmd.AddCommand(importCmd)
}
