package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zhangwlhaut/ypi-ai-cli/internal/client"
	"github.com/zhangwlhaut/ypi-ai-cli/internal/config"
	"github.com/zhangwlhaut/ypi-ai-cli/internal/output"
)

var (
	flagServer    string
	flagToken     string
	flagProjectID int
	flagPretty    bool
	flagJSON      bool

	apiClient *client.Client
	appConfig *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "yapi",
	Short: "面向 AI Agent 的 YApi 命令行工具",
	Long:  "yapi-ai-cli 是一个面向 AI Agent 的 YApi 命令行工具，通过 YApi Open API 实现接口文档管理操作。",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// auth login / skill init don't need config
		if cmd.Name() == "login" || cmd.Name() == "init" {
			return nil
		}
		return initClient(cmd)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagServer, "server", "", "YApi 服务地址（覆盖配置）")
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "YApi 项目 token（覆盖配置）")
	rootCmd.PersistentFlags().IntVar(&flagProjectID, "project-id", 0, "YApi 项目 ID（覆盖配置）")
	rootCmd.PersistentFlags().BoolVar(&flagPretty, "pretty", false, "输出格式化 JSON（便于人类阅读；默认紧凑 JSON 节省 token）")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "强制 JSON 输出（详情默认 Markdown 更省 token）")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initClient(cmd *cobra.Command) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	appConfig = cfg

	// CLI flags override config
	server := flagServer
	token := flagToken
	projectID := flagProjectID

	if server == "" {
		server = cfg.Server
	}
	if token == "" {
		token = cfg.Token
	}
	if projectID == 0 {
		projectID = cfg.ProjectID
	}

	if token == "" {
		output.PrintError(os.Stderr, "auth_required",
			"No YApi token found",
			"Run 'yapi auth login --token <your_token>' or set YAPI_TOKEN environment variable.")
		os.Exit(output.ExitAuthError)
	}

	// Commands that don't need project_id
	skipProjectCheck := map[string]bool{"auth": true, "login": true, "init": true, "show": true}
	parentName := ""
	if cmd.Parent() != nil {
		parentName = cmd.Parent().Name()
	}
	needsProject := !skipProjectCheck[parentName] && !skipProjectCheck[cmd.Name()] && cmd.Name() != "url"
	if needsProject && projectID == 0 {
		output.PrintError(os.Stderr, "project_id_required",
			"No project ID configured",
			"Run 'yapi auth login --token <token>' with a project that has project_id, or use --project-id flag.")
		os.Exit(output.ExitParamError)
	}

	apiClient = client.New(server, token, projectID)
	return nil
}

func useJSONOutput() bool {
	return flagJSON || flagPretty
}
