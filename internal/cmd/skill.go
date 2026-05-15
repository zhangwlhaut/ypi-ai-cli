package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "AI 工具集成",
}

var skillInitCmd = &cobra.Command{
	Use:   "init",
	Short: "为主流 AI Coding 工具生成 SKILL.md 指令文件",
	Long: "为主流 AI Coding 工具生成 SKILL.md 指令文件，使 AI Agent 可以自动发现和使用 yapi CLI。\n\n" +
		"支持的工具：Claude Code、Cursor、Windsurf、Trae、Codex、Gemini CLI、Cline、Roo Code、Augment。\n\n" +
		"命令会自动检测当前目录下已有的工具配置文件夹并默认选中。",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// Detect tool config directories
		toolDirs := map[string]string{
			"Claude Code": ".claude",
			"Cursor":      ".cursor",
			"Windsurf":    ".windsurf",
			"Trae":        ".trae",
			"Cline":       ".cline",
			"Roo Code":    ".roo",
		}

		var detected []string
		for name, dir := range toolDirs {
			if _, err := os.Stat(filepath.Join(cwd, dir)); err == nil {
				detected = append(detected, name)
			}
		}

		if len(detected) == 0 {
			detected = []string{"Claude Code"}
		}

		// Generate SKILL.md
		skillContent := generateSkillMD()

		for _, name := range detected {
			dir := toolDirs[name]
			targetDir := filepath.Join(cwd, dir)
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not create %s: %v\n", targetDir, err)
				continue
			}
			targetPath := filepath.Join(targetDir, "SKILL.md")
			if err := os.WriteFile(targetPath, []byte(skillContent), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not write %s: %v\n", targetPath, err)
				continue
			}
			fmt.Fprintf(os.Stderr, "Generated %s for %s\n", targetPath, name)
		}

		return nil
	},
}

func generateSkillMD() string {
	return "# YApi CLI\n\n" +
		"YApi 接口文档管理命令行工具。\n\n" +
		"## 认证\n\n" +
		"```bash\n" +
		"yapi auth login --token <project_token> [--server <url>] [--local]\n" +
		"```\n\n" +
		"环境变量：YAPI_SERVER, YAPI_TOKEN, YAPI_PROJECT_ID\n\n" +
		"## 命令参考\n\n" +
		"```bash\n" +
		"# 项目\n" +
		"yapi project show                          # 获取项目基本信息\n\n" +
		"# 分类\n" +
		"yapi category list                         # 获取接口分类列表\n" +
		"yapi category add --name <name> [--desc]   # 新增接口分类\n\n" +
		"# 接口\n" +
		"yapi interface list [--catid <id>] [--page] [--limit]  # 获取接口列表\n" +
		"yapi interface show <id>                               # 获取接口详情\n" +
		"yapi interface add --path <path> --method <method> --catid <id> [-y|--new-only]  # 新增接口（自动查重）\n" +
		"yapi interface update <id> --path <path> ...           # 按 ID 更新接口\n" +
		"yapi interface menu                                    # 获取接口菜单树\n\n" +
		"# 导入\n" +
		"yapi import --type <swagger|json> --merge <normal|good|merge> [--json <data>] [--url <url>]\n\n" +
		"# URL 查询\n" +
		"yapi url <yapi-frontend-url>                          # 从 YApi 前端 URL 查询\n" +
		"```\n\n" +
		"## 全局标志\n\n" +
		"| 标志 | 说明 |\n" +
		"|------|------|\n" +
		"| --server <url> | YApi 服务地址 |\n" +
		"| --token <token> | 项目 token |\n" +
		"| --project-id <id> | 项目 ID |\n" +
		"| --pretty | 格式化 JSON 输出 |\n" +
		"| --json | 强制 JSON 输出（详情默认 Markdown） |\n\n" +
		"## add 命令的 upsert 逻辑\n\n" +
		"当执行 `yapi interface add` 时，会自动按 path+method 查重：\n" +
		"- 无重复 → 新增\n" +
		"- 有重复 → 提示用户确认更新\n" +
		"  - `-y` 自动确认（AI Agent 场景推荐）\n" +
		"  - `--new-only` 严格新增，有重复则报错\n"
}

func init() {
	skillCmd.AddCommand(skillInitCmd)
	rootCmd.AddCommand(skillCmd)
}
