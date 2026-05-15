# yapi-ai-cli

面向 AI Agent 的 YApi 命令行工具 / AI-Agent-oriented YApi CLI

通过 YApi Open API 实现接口文档管理操作，让 Claude Code、Cursor 等 AI Coding 工具可以方便地查询、创建、更新 YApi 接口文档。

---

## 安装 / Installation

### 方式一：go install（推荐）

```bash
go install github.com/studyzy/yapi-ai-cli/cmd/yapi@latest
```

### 方式二：从源码构建

```bash
git clone git@github.com:zhangwlhaut/ypi-ai-cli.git
cd ypi-ai-cli
make build     # 在当前目录生成 ./yapi
make install   # 编译并安装到 $GOPATH/bin
```

---

## 认证 / Authentication

```bash
yapi auth login --token <project_token> [--server <url>] [--project-id <id>] [--local]
```

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--token` | YApi 项目 token（必填） | - |
| `--server` | YApi 服务地址 | `http://127.0.0.1:3000` |
| `--project-id` | 项目 ID | - |
| `--local` | 保存凭据到当前目录 `.yapi.json`（默认保存到 `~/.yapi.json`） | - |

**凭据优先级** / Credential priority：CLI flags > 环境变量 > `./.yapi.json` > `~/.yapi.json`

### 环境变量 / Environment Variables

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `YAPI_SERVER` | YApi 服务地址 | `http://127.0.0.1:3000` |
| `YAPI_TOKEN` | 项目 token | - |
| `YAPI_PROJECT_ID` | 项目 ID | - |

### 配置文件 / Config File

`~/.yapi.json` 或 `./.yapi.json`：

```json
{
  "server": "https://yapi.example.com",
  "token": "your-project-token",
  "project_id": 550
}
```

---

## 命令一览 / Command Reference

```
yapi
├── auth        login --token <token> [--server <url>] [--project-id <id>] [--local]
├── project     show
├── category    list | add --name <name> [--desc <desc>]
├── interface   list [--catid <id>] [--page <n>] [--limit <n>]
│               show <id>
│               add --path <path> --method <method> --catid <id> [-y|--new-only]
│               update <id> [--path] [--method] [--title] [--status] ...
│               menu
├── import      --type <type> --merge <mode> [--json <data>] [--url <url>]
├── skill       init
└── url         <yapi-frontend-url>
```

---

## 基本用法 / Basic Usage

### 项目 / Project

```bash
# 获取项目基本信息 / Get project info
yapi project show
```

输出示例（默认 Markdown）：

```
# UMPlayer异步卡

- **ID**: 550
- **Base Path**:
- **Group ID**: 245
```

加 `--json` 或 `--pretty` 可切换为 JSON 输出。

### 分类 / Category

```bash
# 获取接口分类列表 / List interface categories
yapi category list

# 新增接口分类 / Add a new category
yapi category add --name "用户服务" --desc "用户相关接口"
```

### 接口 / Interface

#### 获取接口列表 / List interfaces

```bash
# 获取项目全部接口 / List all interfaces in project
yapi interface list --limit 1000

# 获取指定分类下接口 / List interfaces under a category
yapi interface list --catid 8613 --page 1 --limit 20
```

#### 获取接口详情 / Show interface detail

```bash
yapi interface show 33881
```

输出示例：

```markdown
# 节目-播放控制

- **Method**: POST
- **Path**: /api/program/control
- **Status**: done
- **Category ID**: 8613

## Request Headers

| Name | Value | Required |
|------|-------|----------|
| Content-Type |  | Yes |
| name |  | Yes |
| deviceId |  | Yes |
| token |  | Yes |

## Request Body (json)

```json
{"type":"object","properties":{"id":{"type":"string","description":"节目id"},"action":{"type":"string","description":"play 播放 stop 停止"}},"required":["action","id"]}
```

## Response (json)

```json
{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","properties":{"code":{"type":"number"},"msg":{"type":"string"}}}
```
```

#### 新增接口（自动查重）/ Add interface (auto-dedup)

```bash
yapi interface add --path /api/user/list --method GET --catid 8613 --title "用户列表"
```

**核心逻辑** / **Core logic**：

1. 按 `path + method` 在项目中查找已有接口 / Search for existing interface by path+method
2. **未找到** → 直接新增 / Not found → create new
3. **找到匹配** → 提示用户确认是否更新 / Found duplicate → prompt user to confirm update

交互模式示例 / Interactive mode：

```
⚠ Found existing interface with same path+method:
  ID: 33881 | /api/program/control (POST) | CatID: 8613 | Status: done

Update this interface? [y/N]: y

✓ Interface updated (id: 33881)
```

| Flag | 说明 | 场景 |
|------|------|------|
| `-y` / `--yes` | 发现重复时自动确认更新，不询问 | AI Agent 场景（推荐） |
| `--new-only` | 严格新增模式，有重复则报错退出 | CI 校验 / 确保不误更新 |

```bash
# AI Agent 中使用 / Use in AI Agent
yapi interface add --path /api/user/list --method GET --catid 8613 --title "用户列表" -y

# 严格新增 / Strict create only
yapi interface add --path /api/user/list --method GET --catid 8613 --title "用户列表" --new-only
```

#### 更新接口 / Update interface

```bash
# 按 ID 更新，只修改指定字段 / Update by ID, only specified fields
yapi interface update 33881 --status done
yapi interface update 33881 --title "新标题" --desc "新描述"
```

> YApi 更新接口要求 path 和 method 为必填字段，`update` 命令会自动获取当前值填充，你只需指定要修改的字段。

#### 获取接口菜单树 / Get interface menu tree

```bash
yapi interface menu
```

输出示例：

```markdown
## 安卓播放器-升级 (id: 8736)

| ID | Title | Path | Method | Status |
|----|-------|------|--------|--------|
| 34511 | 升级-开始升级 | /api/upgrade/start | POST | done |
| 34518 | 升级-检测升级状态 | /api/upgrade/status | GET | done |

## 安卓播放器-播放服务 (id: 8613)

| ID | Title | Path | Method | Status |
|----|-------|------|--------|--------|
| 33881 | 节目-播放控制 | /api/program/control | POST | done |
```

### 数据导入 / Data Import

```bash
yapi import --type swagger --merge good --url https://example.com/swagger.json
```

| 参数 | 说明 |
|------|------|
| `--type` | 导入类型：`swagger`、`json` |
| `--merge` | 同步模式：`normal`（普通模式）、`good`（智能合并）、`merge`（完全覆盖） |
| `--json` | JSON 数据（序列化后的字符串） |
| `--url` | 导入数据 URL（优先于 --json） |

### URL 查询 / URL Query

从 YApi 前端页面 URL 中提取信息并自动查询：

```bash
# 查询接口详情 / Query interface detail
yapi url https://yapi.example.com/project/550/interface/api/33881

# 查询分类下接口 / Query interfaces under category
yapi url https://yapi.example.com/project/550/interface/api/cat_8613

# 查询项目信息 / Query project info
yapi url https://yapi.example.com/project/550
```

### AI 工具集成 / AI Tool Integration

```bash
yapi skill init
```

一键为主流 AI Coding 工具生成 `SKILL.md` 指令文件，使 AI Agent 可以自动发现和使用 yapi CLI。

支持的工具 / Supported tools：Claude Code、Cursor、Windsurf、Trae、Codex、Gemini CLI、Cline、Roo Code、Augment

命令会自动检测当前目录下已有的工具配置文件夹并默认选中。

---

## 全局标志 / Global Flags

| 标志 | 说明 |
|------|------|
| `--server <url>` | YApi 服务地址（覆盖配置） |
| `--token <token>` | 项目 token（覆盖配置） |
| `--project-id <id>` | 项目 ID（覆盖配置） |
| `--pretty` | 输出格式化 JSON（便于人类阅读；默认紧凑 JSON 节省 token） |
| `--json` | 强制 JSON 输出（详情默认 Markdown 更省 token） |

---

## 输出格式 / Output Format

面向 AI Agent 优化，默认输出最省 token 的格式：

| 命令类型 | 默认格式 | `--json` | `--pretty` |
|----------|----------|----------|------------|
| 列表命令 | 紧凑 JSON（无缩进） | 紧凑 JSON | 格式化 JSON |
| 详情命令 | Markdown（可读性好、省 token） | 紧凑 JSON | 格式化 JSON |
| 错误输出 | 结构化 JSON（含 error_code + hint） | - | - |

```bash
# 紧凑 JSON（默认，省 token）/ Compact JSON (default, token-saving)
yapi interface list --limit 2

# Markdown 格式（详情默认）/ Markdown format (default for details)
yapi interface show 33881

# 格式化 JSON（人类阅读）/ Pretty JSON (human reading)
yapi interface show 33881 --pretty
```

---

## Windows Git Bash 注意事项 / Windows Git Bash Notes

在 Git Bash 中，以 `/` 开头的路径可能被 MSYS 自动转换为 Windows 路径（如 `/api/foo` → `C:/Program Files/Git/api/foo`）。本工具已内置自动修复，无需额外处理。

如果遇到路径问题，可设置环境变量：

```bash
export MSYS_NO_PATHCONV=1
```

---

## 开发 / Development

```bash
make build      # 构建 / Build
make install    # 安装到 $GOPATH/bin / Install to $GOPATH/bin
make test       # 运行测试 / Run tests
make clean      # 清理构建产物 / Clean build artifacts
```

---

## YApi Open API 覆盖 / API Coverage

| API | Path | Method | CLI 命令 |
|-----|------|--------|----------|
| 获取项目基本信息 | `/api/project/get` | GET | `yapi project show` |
| 新增接口分类 | `/api/interface/add_cat` | POST | `yapi category add` |
| 获取菜单列表 | `/api/interface/getCatMenu` | GET | `yapi category list` |
| 获取接口详情 | `/api/interface/get` | GET | `yapi interface show` |
| 获取分类下接口列表 | `/api/interface/list_cat` | GET | `yapi interface list --catid` |
| 获取接口列表 | `/api/interface/list` | GET | `yapi interface list` |
| 获取接口菜单列表 | `/api/interface/list_menu` | GET | `yapi interface menu` |
| 新增接口 | `/api/interface/add` | POST | `yapi interface add` |
| 更新接口 | `/api/interface/up` | POST | `yapi interface update` / `add -y` |
| 新增或更新接口 | `/api/interface/save` | POST | `yapi interface add -y` |
| 服务端数据导入 | `/api/open/import_data` | POST | `yapi import` |

---

## 许可证 / License

Apache License 2.0
