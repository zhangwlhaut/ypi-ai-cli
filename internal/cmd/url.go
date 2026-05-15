package cmd

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/studyzy/yapi-ai-cli/internal/output"
)

var urlCmd = &cobra.Command{
	Use:   "url <yapi-url>",
	Short: "通过 YApi 前端 URL 查询接口信息",
	Long: `从 YApi 前端 URL 中提取项目/接口信息并查询。
支持格式：
  http://yapi.example.com/project/<project_id>/interface/api/cat_<catid>
  http://yapi.example.com/project/<project_id>/interface/api/<interface_id>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		yapiURL := args[0]

		parsed, err := url.Parse(yapiURL)
		if err != nil {
			output.PrintError(os.Stderr, "invalid_url", "Invalid URL: "+err.Error(), "")
			os.Exit(output.ExitParamError)
		}

		// Extract project_id and interface/cat id from path
		// Expected patterns:
		//   /project/<id>/interface/api/cat_<catid>
		//   /project/<id>/interface/api/<interface_id>
		path := parsed.Path

		// Extract project ID
		projectRe := regexp.MustCompile(`/project/(\d+)`)
		projectMatches := projectRe.FindStringSubmatch(path)
		if len(projectMatches) < 2 {
			output.PrintError(os.Stderr, "invalid_url",
				"Cannot extract project_id from URL",
				"Expected format: /project/<project_id>/interface/api/...")
			os.Exit(output.ExitParamError)
		}
		projectID, _ := strconv.Atoi(projectMatches[1])

		// Extract interface or category ID
		catRe := regexp.MustCompile(`/interface/api/cat_(\d+)`)
		catMatches := catRe.FindStringSubmatch(path)
		ifaceRe := regexp.MustCompile(`/interface/api/(\d+)$`)
		ifaceMatches := ifaceRe.FindStringSubmatch(path)

		// Temporarily override project ID
		origProjectID := apiClient.ProjectID
		apiClient.ProjectID = projectID
		defer func() { apiClient.ProjectID = origProjectID }()

		if len(ifaceMatches) >= 2 {
			// Interface detail
			ifaceID, _ := strconv.Atoi(ifaceMatches[1])
			iface, err := apiClient.GetInterface(ifaceID)
			if err != nil {
				output.PrintError(os.Stderr, "api_error", err.Error(), "")
				os.Exit(output.ExitAPIError)
			}
			if useJSONOutput() {
				return output.PrintJSON(os.Stdout, iface, !flagPretty)
			}
			return output.PrintInterfaceMarkdown(os.Stdout, iface)
		}

		if len(catMatches) >= 2 {
			// Category list
			catID, _ := strconv.Atoi(catMatches[1])
			list, err := apiClient.ListInterfacesByCat(catID, 1, 1000)
			if err != nil {
				output.PrintError(os.Stderr, "api_error", err.Error(), "")
				os.Exit(output.ExitAPIError)
			}
			return output.PrintJSON(os.Stdout, list, !flagPretty)
		}

		// Just project — show project info
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
	_ = fmt.Sprintf // avoid unused import
	rootCmd.AddCommand(urlCmd)
}
