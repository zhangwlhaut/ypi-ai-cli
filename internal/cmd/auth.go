package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zhangwlhaut/yapi-ai-cli/internal/config"
)

var authLocal bool

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "管理 YApi 认证",
}

var authLoginCmd = &cobra.Command{
	Use:   "login --token <token> [--server <url>] [--project-id <id>] [--local]",
	Short: "登录 YApi，保存凭据",
	Long:  "登录 YApi，保存项目 token、服务地址和项目 ID 到配置文件。\n凭据优先级：CLI flags > 环境变量 > ./.yapi.json > ~/.yapi.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := cmd.Flags().GetString("token")
		server, _ := cmd.Flags().GetString("server")
		projectID, _ := cmd.Flags().GetInt("project-id")

		if token == "" {
			return fmt.Errorf("--token is required")
		}

		// Load existing config first (to preserve values)
		cfg, _ := config.LoadConfig()

		// Override with provided values
		cfg.Token = token
		if server != "" {
			cfg.Server = server
		}
		if projectID > 0 {
			cfg.ProjectID = projectID
		}

		if err := config.SaveConfig(cfg, authLocal); err != nil {
			return err
		}

		path := configFileName(false)
		if authLocal {
			path = configFileName(true)
		}
		fmt.Fprintf(os.Stderr, "Credentials saved to %s\n", path)
		fmt.Fprintf(os.Stdout, `{"success":true,"server":"%s","project_id":%d}`+"\n", cfg.Server, cfg.ProjectID)
		return nil
	},
}

func init() {
	authLoginCmd.Flags().String("token", "", "YApi 项目 token (required)")
	authLoginCmd.Flags().String("server", "", "YApi 服务地址 (default: http://127.0.0.1:3000)")
	authLoginCmd.Flags().Int("project-id", 0, "YApi 项目 ID")
	authLoginCmd.Flags().BoolVar(&authLocal, "local", false, "保存凭据到当前目录 .yapi.json（默认保存到 ~/.yapi.json）")

	authCmd.AddCommand(authLoginCmd)
	rootCmd.AddCommand(authCmd)
}

func configFileName(local bool) string {
	if local {
		return ".yapi.json"
	}
	home, _ := os.UserHomeDir()
	return home + "/.yapi.json"
}
