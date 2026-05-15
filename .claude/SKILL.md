# YApi CLI

YApi 接口文档管理命令行工具。

## 认证

```bash
yapi auth login --token <project_token> [--server <url>] [--local]
```

环境变量：YAPI_SERVER, YAPI_TOKEN, YAPI_PROJECT_ID

## 命令参考

```bash
# 项目
yapi project show                          # 获取项目基本信息

# 分类
yapi category list                         # 获取接口分类列表
yapi category add --name <name> [--desc]   # 新增接口分类

# 接口
yapi interface list [--catid <id>] [--page] [--limit]  # 获取接口列表
yapi interface show <id>                               # 获取接口详情
yapi interface add --path <path> --method <method> --catid <id> [-y|--new-only]  # 新增接口（自动查重）
yapi interface update <id> --path <path> ...           # 按 ID 更新接口
yapi interface menu                                    # 获取接口菜单树

# 导入
yapi import --type <swagger|json> --merge <normal|good|merge> [--json <data>] [--url <url>]

# URL 查询
yapi url <yapi-frontend-url>                          # 从 YApi 前端 URL 查询
```

## 全局标志

| 标志 | 说明 |
|------|------|
| --server <url> | YApi 服务地址 |
| --token <token> | 项目 token |
| --project-id <id> | 项目 ID |
| --pretty | 格式化 JSON 输出 |
| --json | 强制 JSON 输出（详情默认 Markdown） |

## add 命令的 upsert 逻辑

当执行 `yapi interface add` 时，会自动按 path+method 查重：
- 无重复 → 新增
- 有重复 → 提示用户确认更新
  - `-y` 自动确认（AI Agent 场景推荐）
  - `--new-only` 严格新增，有重复则报错
