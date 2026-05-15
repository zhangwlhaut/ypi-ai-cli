package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/studyzy/yapi-ai-cli/internal/model"
	"github.com/studyzy/yapi-ai-cli/internal/output"
)

// stripMSYSPathPrefix removes the MSYS path conversion prefix that Git Bash on Windows
// may add to paths starting with / (e.g., "C:/Program Files/Git/api/foo" -> "/api/foo").
func stripMSYSPathPrefix(s string) string {
	// Common MSYS prefix patterns: "C:/Program Files/Git/" or similar
	if strings.Contains(s, "/Git/") && !strings.HasPrefix(s, "/") {
		if idx := strings.Index(s, "/Git/"); idx >= 0 {
			rest := s[idx+5:]
			// Ensure the result starts with /
			if !strings.HasPrefix(rest, "/") {
				rest = "/" + rest
			}
			return rest
		}
	}
	return s
}

var interfaceCmd = &cobra.Command{
	Use:     "interface",
	Short:   "接口操作",
	Aliases: []string{"api"},
}

// --- list ---
var (
	ifaceListCatID int
	ifaceListPage  int
	ifaceListLimit int
)

var interfaceListCmd = &cobra.Command{
	Use:   "list [--catid <id>] [--page <n>] [--limit <n>]",
	Short: "获取接口列表",
	RunE: func(cmd *cobra.Command, args []string) error {
		if ifaceListCatID > 0 {
			list, err := apiClient.ListInterfacesByCat(ifaceListCatID, ifaceListPage, ifaceListLimit)
			if err != nil {
				output.PrintError(os.Stderr, "api_error", err.Error(), "")
				os.Exit(output.ExitAPIError)
			}
			return output.PrintJSON(os.Stdout, list, !flagPretty)
		}

		list, err := apiClient.ListInterfaces(ifaceListPage, ifaceListLimit)
		if err != nil {
			output.PrintError(os.Stderr, "api_error", err.Error(), "")
			os.Exit(output.ExitAPIError)
		}
		return output.PrintJSON(os.Stdout, list, !flagPretty)
	},
}

// --- show ---
var interfaceShowCmd = &cobra.Command{
	Use:  "show <id>",
	Short: "获取接口详情（含完整请求/响应定义）",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			output.PrintError(os.Stderr, "invalid_id", "interface id must be a number", "")
			os.Exit(output.ExitParamError)
		}

		iface, err := apiClient.GetInterface(id)
		if err != nil {
			output.PrintError(os.Stderr, "api_error", err.Error(), "")
			os.Exit(output.ExitAPIError)
		}

		if useJSONOutput() {
			return output.PrintJSON(os.Stdout, iface, !flagPretty)
		}
		return output.PrintInterfaceMarkdown(os.Stdout, iface)
	},
}

// --- add (with upsert by path+method) ---
var (
	ifaceAddPath         string
	ifaceAddMethod       string
	ifaceAddCatID        int
	ifaceAddTitle        string
	ifaceAddDesc         string
	ifaceAddReqBodyType  string
	ifaceAddReqBodyOther string
	ifaceAddResBodyType  string
	ifaceAddResBody      string
	ifaceAddStatus       string
	ifaceAddYes          bool
	ifaceAddNewOnly      bool
)

var interfaceAddCmd = &cobra.Command{
	Use:   "add --path <path> --method <method> --catid <id> [--title <title>] [-y|--new-only]",
	Short: "新增接口（自动按 path+method 查重，匹配到则确认后更新）",
	Long: `新增接口。如果项目中已存在相同 path+method 的接口：
  - 默认：提示用户确认是否更新
  - -y：自动确认更新，不询问（适合 AI Agent）
  - --new-only：有重复则报错，不更新`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ifaceAddPath == "" {
			output.PrintError(os.Stderr, "param_required", "--path is required", "")
			os.Exit(output.ExitParamError)
		}
		ifaceAddPath = stripMSYSPathPrefix(ifaceAddPath)
		if ifaceAddMethod == "" {
			output.PrintError(os.Stderr, "param_required", "--method is required", "")
			os.Exit(output.ExitParamError)
		}
		if ifaceAddCatID == 0 {
			output.PrintError(os.Stderr, "param_required", "--catid is required", "")
			os.Exit(output.ExitParamError)
		}

		req := &model.AddInterfaceRequest{
			CatID:        ifaceAddCatID,
			Title:        ifaceAddTitle,
			Path:         ifaceAddPath,
			Method:       ifaceAddMethod,
			Desc:         ifaceAddDesc,
			ReqBodyType:  ifaceAddReqBodyType,
			ReqBodyOther: ifaceAddReqBodyOther,
			ResBodyType:  ifaceAddResBodyType,
			ResBody:      ifaceAddResBody,
			Status:       ifaceAddStatus,
		}

		// Check for existing interface with same path+method
		existing, err := apiClient.FindInterfaceByPathMethod(ifaceAddPath, ifaceAddMethod)
		if err != nil {
			// If lookup fails, just proceed with add
			fmt.Fprintf(os.Stderr, "Warning: could not check for duplicates: %v\n", err)
		}

		if existing != nil {
			if ifaceAddNewOnly {
				output.PrintError(os.Stderr, "duplicate_found",
					fmt.Sprintf("Interface with path=%s method=%s already exists (id: %d)", ifaceAddPath, ifaceAddMethod, existing.ID),
					"Remove --new-only to allow update, or use -y to auto-confirm.")
				os.Exit(output.ExitParamError)
			}

			if !ifaceAddYes {
				if !output.ConfirmDuplication(existing) {
					fmt.Fprintln(os.Stderr, "Aborted.")
					return nil
				}
			}

			// Update existing
			req.ID = existing.ID
			fmt.Fprintf(os.Stderr, "Updating existing interface (id: %d)...\n", existing.ID)
			result, err := apiClient.UpdateInterface(req)
			if err != nil {
				output.PrintError(os.Stderr, "api_error", err.Error(), "")
				os.Exit(output.ExitAPIError)
			}
			fmt.Fprintf(os.Stderr, "Interface updated (id: %d)\n", existing.ID)
			return output.PrintJSON(os.Stdout, json.RawMessage(result), !flagPretty)
		}

		// No duplicate, create new
		result, err := apiClient.AddInterface(req)
		if err != nil {
			output.PrintError(os.Stderr, "api_error", err.Error(), "")
			os.Exit(output.ExitAPIError)
		}
		fmt.Fprintf(os.Stderr, "Interface created\n")
		return output.PrintJSON(os.Stdout, json.RawMessage(result), !flagPretty)
	},
}

// --- update (explicit by ID) ---
var (
	ifaceUpdPath         string
	ifaceUpdMethod       string
	ifaceUpdCatID        int
	ifaceUpdTitle        string
	ifaceUpdDesc         string
	ifaceUpdReqBodyType  string
	ifaceUpdReqBodyOther string
	ifaceUpdResBodyType  string
	ifaceUpdResBody      string
	ifaceUpdStatus       string
)

var interfaceUpdateCmd = &cobra.Command{
	Use:  "update <id>",
	Short: "按 ID 更新接口",
	Long: "按 ID 更新接口。YApi 更新接口要求 path 和 method 为必填字段，\n如果未指定则自动获取当前值填充。",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			output.PrintError(os.Stderr, "invalid_id", "interface id must be a number", "")
			os.Exit(output.ExitParamError)
		}

		// Fetch current interface to fill required fields
		current, err := apiClient.GetInterface(id)
		if err != nil {
			output.PrintError(os.Stderr, "api_error", "Failed to fetch current interface: "+err.Error(), "")
			os.Exit(output.ExitAPIError)
		}

		req := &model.AddInterfaceRequest{
			ID:           id,
			Path:         current.Path,
			Method:       current.Method,
			CatID:        current.CatID,
			Title:        current.Title,
			Desc:         current.Desc,
			ReqBodyType:  current.ReqBodyType,
			ReqBodyOther: current.ReqBodyOther,
			ResBodyType:  current.ResBodyType,
			ResBody:      current.ResBody,
			Status:       current.Status,
			ReqQuery:     current.ReqQuery,
			ReqHeaders:   current.ReqHeaders,
			ReqBodyForm:  current.ReqBodyForm,
			ReqParams:    current.ReqParams,
		}

		// Override with explicitly provided flags
		if cmd.Flags().Changed("path") {
			req.Path = stripMSYSPathPrefix(ifaceUpdPath)
		}
		if cmd.Flags().Changed("method") {
			req.Method = ifaceUpdMethod
		}
		if cmd.Flags().Changed("catid") {
			req.CatID = ifaceUpdCatID
		}
		if cmd.Flags().Changed("title") {
			req.Title = ifaceUpdTitle
		}
		if cmd.Flags().Changed("desc") {
			req.Desc = ifaceUpdDesc
		}
		if cmd.Flags().Changed("req-body-type") {
			req.ReqBodyType = ifaceUpdReqBodyType
		}
		if cmd.Flags().Changed("req-body") {
			req.ReqBodyOther = ifaceUpdReqBodyOther
		}
		if cmd.Flags().Changed("res-body-type") {
			req.ResBodyType = ifaceUpdResBodyType
		}
		if cmd.Flags().Changed("res-body") {
			req.ResBody = ifaceUpdResBody
		}
		if cmd.Flags().Changed("status") {
			req.Status = ifaceUpdStatus
		}

		result, err := apiClient.UpdateInterface(req)
		if err != nil {
			output.PrintError(os.Stderr, "api_error", err.Error(), "")
			os.Exit(output.ExitAPIError)
		}
		fmt.Fprintf(os.Stderr, "Interface updated (id: %d)\n", id)
		return output.PrintJSON(os.Stdout, json.RawMessage(result), !flagPretty)
	},
}

// --- menu ---
var interfaceMenuCmd = &cobra.Command{
	Use:   "menu",
	Short: "获取接口菜单树（含分类及分类下接口列表）",
	RunE: func(cmd *cobra.Command, args []string) error {
		cats, err := apiClient.ListMenu()
		if err != nil {
			output.PrintError(os.Stderr, "api_error", err.Error(), "")
			os.Exit(output.ExitAPIError)
		}

		if useJSONOutput() {
			return output.PrintJSON(os.Stdout, cats, !flagPretty)
		}
		return output.PrintCategoryListMarkdown(os.Stdout, cats)
	},
}

func init() {
	// list flags
	interfaceListCmd.Flags().IntVar(&ifaceListCatID, "catid", 0, "分类 ID（指定则只返回该分类下接口）")
	interfaceListCmd.Flags().IntVar(&ifaceListPage, "page", 1, "页码")
	interfaceListCmd.Flags().IntVar(&ifaceListLimit, "limit", 0, "每页数量（默认 10，设大值如 1000 取全量）")

	// add flags
	interfaceAddCmd.Flags().StringVar(&ifaceAddPath, "path", "", "接口路径 (required)")
	interfaceAddCmd.Flags().StringVar(&ifaceAddMethod, "method", "", "请求方法 GET/POST/PUT/DELETE 等 (required)")
	interfaceAddCmd.Flags().IntVar(&ifaceAddCatID, "catid", 0, "分类 ID (required)")
	interfaceAddCmd.Flags().StringVar(&ifaceAddTitle, "title", "", "接口标题")
	interfaceAddCmd.Flags().StringVar(&ifaceAddDesc, "desc", "", "接口描述")
	interfaceAddCmd.Flags().StringVar(&ifaceAddReqBodyType, "req-body-type", "", "请求数据类型: json/form/raw")
	interfaceAddCmd.Flags().StringVar(&ifaceAddReqBodyOther, "req-body", "", "请求体 JSON 字符串")
	interfaceAddCmd.Flags().StringVar(&ifaceAddResBodyType, "res-body-type", "", "返回数据类型: json/raw")
	interfaceAddCmd.Flags().StringVar(&ifaceAddResBody, "res-body", "", "返回数据 JSON 字符串")
	interfaceAddCmd.Flags().StringVar(&ifaceAddStatus, "status", "undone", "接口状态: undone/done")
	interfaceAddCmd.Flags().BoolVarP(&ifaceAddYes, "yes", "y", false, "发现重复时自动确认更新，不询问（AI Agent 场景）")
	interfaceAddCmd.Flags().BoolVar(&ifaceAddNewOnly, "new-only", false, "严格新增模式，有重复则报错")

	// update flags (independent variables, only send explicitly set fields)
	interfaceUpdateCmd.Flags().StringVar(&ifaceUpdPath, "path", "", "接口路径")
	interfaceUpdateCmd.Flags().StringVar(&ifaceUpdMethod, "method", "", "请求方法")
	interfaceUpdateCmd.Flags().IntVar(&ifaceUpdCatID, "catid", 0, "分类 ID")
	interfaceUpdateCmd.Flags().StringVar(&ifaceUpdTitle, "title", "", "接口标题")
	interfaceUpdateCmd.Flags().StringVar(&ifaceUpdDesc, "desc", "", "接口描述")
	interfaceUpdateCmd.Flags().StringVar(&ifaceUpdReqBodyType, "req-body-type", "", "请求数据类型: json/form/raw")
	interfaceUpdateCmd.Flags().StringVar(&ifaceUpdReqBodyOther, "req-body", "", "请求体 JSON 字符串")
	interfaceUpdateCmd.Flags().StringVar(&ifaceUpdResBodyType, "res-body-type", "", "返回数据类型: json/raw")
	interfaceUpdateCmd.Flags().StringVar(&ifaceUpdResBody, "res-body", "", "返回数据 JSON 字符串")
	interfaceUpdateCmd.Flags().StringVar(&ifaceUpdStatus, "status", "", "接口状态: undone/done")

	interfaceCmd.AddCommand(interfaceListCmd)
	interfaceCmd.AddCommand(interfaceShowCmd)
	interfaceCmd.AddCommand(interfaceAddCmd)
	interfaceCmd.AddCommand(interfaceUpdateCmd)
	interfaceCmd.AddCommand(interfaceMenuCmd)
	rootCmd.AddCommand(interfaceCmd)
}
