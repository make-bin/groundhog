# Groundhog 产品功能文档

> 面向前端 UI 设计的功能说明，版本 v1.0

---

## 产品定位

Groundhog 是一个 AI Agent 管理平台，支持多渠道消息接入、多模型对话、定时任务调度、记忆管理和安全审计。

---

## 页面与功能模块

### 1. Dashboard（仪表盘）
系统整体状态概览。

- 健康检查状态（`GET /api/v1/health`）
- 活跃 Session 数量
- 已连接 Channel 数量
- Cron 调度器运行状态（启用任务数 / 运行中任务数 / 下次执行时间）
- 快速入口导航

---

### 2. Sessions（会话管理）
管理 AI Agent 对话会话。

**列表页**
- 按 UserID / AgentID / 状态筛选
- 分页展示（offset / limit）
- 每条显示：SessionID、AgentID、模型、状态、创建时间

**详情 / 聊天页**（`/sessions/:id/chat`）
- 发送消息（普通模式 `POST /sessions/:id/messages`）
- 流式对话（SSE 模式 `POST /sessions/:id/messages/stream`）
  - 事件类型：`chunk`（token 增量）、`tool_start`、`tool_done`、`approval_required`、`done`、`error`
- 工具调用展示（工具名、参数、结果、耗时）
- 人工审批（Approval）
  - 查看待审批工具调用列表
  - 批准 / 拒绝操作

**操作**
- 创建 Session（指定 AgentID、模型、Provider、系统提示词、Skills）
- 删除 Session

**支持的 AI Provider**
`openai` / `anthropic` / `gemini` / `groq` / `ollama` / `openai_compat`

---

### 3. Channels（渠道管理）
管理外部消息渠道（Discord、Telegram、Slack、WhatsApp、Signal 等插件）。

- 渠道列表（名称、类型、状态）
- 创建渠道（选择插件类型、填写配置）
- 删除渠道
- 查看渠道实时状态（`GET /channels/:id/status`）

---

### 4. Cron Jobs（定时任务）
通过 JSON-RPC（`POST /rpc`）管理定时任务调度。

**任务列表**
- 搜索 / 筛选（启用状态、关键词）
- 分页、排序
- 每条显示：名称、调度类型、下次执行时间、运行状态、启用开关

**调度类型**
| 类型 | 说明 |
|------|------|
| `at` | 指定时间点执行一次 |
| `every` | 固定间隔循环执行（毫秒） |
| `cron` | Cron 表达式（支持时区、随机抖动） |

**任务载荷类型**
| 类型 | 说明 |
|------|------|
| `agentTurn` | 触发 Agent 执行一轮对话（可指定模型、超时、消息内容） |
| `systemEvent` | 发送系统事件文本 |

**投递配置（Delivery）**
- 模式：`channel`（发送到渠道）/ `session`（发送到会话）
- 目标：Channel ID、收件人、AccountID
- 失败降级（best_effort）

**失败告警（FailureAlert）**
- 连续失败 N 次后告警
- 告警渠道、冷却时间

**操作**
- 新建 / 编辑 / 删除任务
- 手动触发（`force` 立即执行 / `due` 仅到期时执行）
- 查看调度器状态（运行中任务数、下次执行时间）
- 查看执行日志（按任务 / 全局，支持分页）

---

### 5. Memory（记忆管理）
管理用户的长期记忆，支持向量混合检索。

- 记忆列表（当前用户，分页）
- 新建记忆（输入文本内容，自动生成向量嵌入）
- 查看 / 编辑 / 删除记忆
- 语义搜索（输入查询词，返回相关记忆及相似度分数）

> 记忆与 Session 联动：当 Session 开启记忆功能时，Agent 会在对话前自动检索相关记忆并在对话后保存关键信息。

---

### 6. Security（安全审计）
查看系统操作审计日志。

- 按操作类型（action）筛选
- 按用户 ID（principal_id）筛选
- 分页展示（page / page_size）
- 每条显示：操作类型、操作人、时间、详情

---

### 7. Config（配置）
查看当前系统运行配置（只读展示）。

- 服务器配置（host、port、超时）
- 数据库配置
- Redis 配置
- JWT 配置
- 模型配置（默认 Provider / 模型、各 Provider 参数）
- MCP 工具服务器列表（名称、命令、危险工具、是否需要审批）
- 记忆配置（嵌入模型、向量维度、混合检索权重）
- Skills 目录配置

---

## 全局说明

**认证**：请求头携带 `Authorization: Bearer <token>`，用户标识通过 `X-User-ID` 传递。

**响应格式**：
```json
{ "code": 200, "data": { ... } }
```

**WebSocket**：`/ws` 端点用于实时事件推送（Session 状态变更、工具调用事件等）。

**SSE 流式对话**：`POST /api/v1/sessions/:id/messages/stream`，返回 `text/event-stream`，每条事件格式为 `data: <json>\n\n`。
