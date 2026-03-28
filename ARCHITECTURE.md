# OpenClaw-Go 架构图

## 整体架构（DDD 分层）

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              外部访问入口                                         │
│   Browser/Client          Plugin (gRPC)          CLI Terminal                   │
│   HTTP / WebSocket        Discord/Telegram        cobra commands                 │
│                           Slack/WhatsApp/Signal                                  │
└────────────┬──────────────────────┬──────────────────────┬───────────────────────┘
             │                      │                      │
┌────────────▼──────────────────────▼──────────────────────▼───────────────────────┐
│                           Interface Layer  (pkg/interface/)                       │
│                                                                                   │
│  ┌─────────────────────────────────────────────────────────────────────────────┐  │
│  │  HTTP (pkg/interface/http/)                                                 │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌────────────────┐  │  │
│  │  │HealthHandler │  │SessionHandler│  │ChannelHandler│  │SecurityHandler │  │  │
│  │  │  GET /health │  │GET  /sessions│  │GET  /channels│  │GET /audit      │  │  │
│  │  │              │  │GET  /:id     │  │POST /channels│  │                │  │  │
│  │  │              │  │DELETE /:id   │  │DELETE /:id   │  │                │  │  │
│  │  │              │  │POST /:id/msg │  │GET /:id/status│ │                │  │  │
│  │  └──────────────┘  └──────────────┘  └──────────────┘  └────────────────┘  │  │
│  │                                                                              │  │
│  │  Middleware: AuthMiddleware(JWT) │ CORSMiddleware │ LoggerMiddleware         │  │
│  │  Router: /api/v1/* │ /ws (WebSocket) │ /metrics (Prometheus)               │  │
│  └─────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                   │
│  ┌──────────────────────────┐   ┌──────────────────────────────────────────────┐  │
│  │  WebSocket               │   │  CLI (pkg/interface/cli/)                    │  │
│  │  (pkg/interface/ws/)     │   │  gateway │ config │ status │ onboard │ doctor│  │
│  │  WSEventHandler          │   │  (Cobra commands)                            │  │
│  └──────────────────────────┘   └──────────────────────────────────────────────┘  │
└────────────────────────────────────────┬──────────────────────────────────────────┘
                                         │ depends on
┌────────────────────────────────────────▼──────────────────────────────────────────┐
│                          Application Layer  (pkg/application/)                     │
│                                                                                    │
│  ┌──────────────────────────────────────────────────────────────────────────────┐  │
│  │  Services (pkg/application/service/)                                         │  │
│  │  ┌─────────────────────┐  ┌─────────────────────┐  ┌──────────────────────┐ │  │
│  │  │  AgentAppService    │  │  GatewayAppService  │  │  ChannelAppService   │ │  │
│  │  │  - CreateSession    │  │  - HandleInbound    │  │  - CreateChannel     │ │  │
│  │  │  - GetSession       │  │    Message          │  │  - DeleteChannel     │ │  │
│  │  │  - ListSessions     │  │                     │  │  - ListChannels      │ │  │
│  │  │  - DeleteSession    │  │                     │  │  - GetChannelStatus  │ │  │
│  │  │  - ExecuteTurn      │  │                     │  │                      │ │  │
│  │  └─────────────────────┘  └─────────────────────┘  └──────────────────────┘ │  │
│  └──────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                    │
│  ┌──────────────────────────────────────────────────────────────────────────────┐  │
│  │  DTOs (pkg/application/dto/)                                                 │  │
│  │  SessionDTO │ MessageDTO │ ChannelDTO │ InboundMessageRequest                │  │
│  │  CreateSessionRequest │ SendMessageRequest │ SessionListRequest              │  │
│  └──────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                    │
│  ┌──────────────────────────────────────────────────────────────────────────────┐  │
│  │  Assemblers (pkg/application/assembler/)                                     │  │
│  │  SessionAssembler │ MessageAssembler                                         │  │
│  └──────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                    │
│  ┌──────────────────────────┐   ┌──────────────────────────────────────────────┐  │
│  │  EventBus                │   │  Hook Registry                               │  │
│  │  (pkg/application/       │   │  (pkg/application/hook/)                     │  │
│  │   eventbus/)             │   │  BeforeMessageReceive                        │  │
│  │  Publish / Subscribe     │   │  AfterMessageReceive                         │  │
│  └──────────────────────────┘   └──────────────────────────────────────────────┘  │
└────────────────────────────────────────┬──────────────────────────────────────────┘
                                         │ depends on
┌────────────────────────────────────────▼──────────────────────────────────────────┐
│                            Domain Layer  (pkg/domain/)                             │
│                                                                                    │
│  ┌──────────────────────────────────────────────────────────────────────────────┐  │
│  │  Conversation Domain                                                         │  │
│  │  Aggregate: AgentSession                                                     │  │
│  │    id │ agentID │ userID │ turns │ activeModel │ tools │ systemPrompt        │  │
│  │    skills │ state │ createdAt │ lastActiveAt │ metadata                      │  │
│  │  Entities: Turn │ Agent │ ToolDefinition │ ToolCall                          │  │
│  │  VOs: SessionID │ AgentID │ ModelConfig │ Prompt │ SessionState              │  │
│  │       TokenUsage │ ToolPolicy │ ToolResult                                   │  │
│  │  Services: CompactionService │ ModelSelectionService │ SystemPromptService   │  │
│  │            ToolPolicyService                                                 │  │
│  │  Events: SessionCreated │ SessionArchived │ TurnCompleted │ ToolCallApproved │  │
│  │  Repository: SessionRepository                                               │  │
│  └──────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                    │
│  ┌──────────────────────────────────────────────────────────────────────────────┐  │
│  │  Messaging Domain                                                            │  │
│  │  Aggregate: InboundMessage                                                   │  │
│  │    id │ channelID │ accountID │ content │ receivedAt │ routedTo              │  │
│  │    status │ chunks                                                           │  │
│  │  Entities: Channel │ Route │ AutoReplyRule                                   │  │
│  │  VOs: MessageID │ ChannelID │ AccountID │ MessageContent │ MessageChunk      │  │
│  │       MessageStatus                                                          │  │
│  │  Services: RoutingService │ CommandDetectionService │ MessageChunkingService │  │
│  │  Events: MessageReceived │ MessageRouted │ MessageDelivered │ MessageFailed  │  │
│  │  Repository: MessageRepository │ ChannelRepository                          │  │
│  └──────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                    │
│  ┌──────────────────────────────────────────────────────────────────────────────┐  │
│  │  Plugin Domain                                                               │  │
│  │  Aggregate: PluginInstance                                                   │  │
│  │    id │ manifest │ process │ status │ capabilities │ startedAt │ restartCount│  │
│  │  Entities: PluginManifest │ PluginProcess                                    │  │
│  │  VOs: PluginID │ Capability                                                  │  │
│  │  Services: PluginLifecycleService                                            │  │
│  │  Events: PluginStarted │ PluginStopped │ PluginCrashed                       │  │
│  │  Repository: PluginRepository                                                │  │
│  └──────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                    │
│  ┌────────────────────────┐  ┌────────────────────────┐  ┌──────────────────────┐ │
│  │  Media Domain          │  │  Identity Domain       │  │  Configuration Domain│ │
│  │  Aggregate: MediaAsset │  │  Aggregate: Principal  │  │  Aggregate: ConfigRoot│ │
│  │  VOs: MimeType         │  │  Entities: Device      │  │  VOs: SchemaVersion  │ │
│  │       MediaRef         │  │           AuditEntry   │  │       SecretRef      │ │
│  │  Repo: MediaRepository │  │  VOs: PrincipalID      │  │       ValidationError│ │
│  │                        │  │       Credential       │  │  Services: Validation│ │
│  │                        │  │       RateLimitPolicy  │  │           Migration  │ │
│  │                        │  │  Services: AuthService │  │                      │ │
│  │                        │  │           AuditService │  │                      │ │
│  │                        │  │           RateLimitSvc │  │                      │ │
│  └────────────────────────┘  └────────────────────────┘  └──────────────────────┘ │
└────────────────────────────────────────┬──────────────────────────────────────────┘
                                         │ implemented by
┌────────────────────────────────────────▼──────────────────────────────────────────┐
│                        Infrastructure Layer  (pkg/infrastructure/)                 │
│                                                                                    │
│  ┌──────────────────────────────────────────────────────────────────────────────┐  │
│  │  Persistence (pkg/infrastructure/persistence/)                               │  │
│  │                                                                              │  │
│  │  Repository Implementations:                                                 │  │
│  │  SessionRepositoryImpl │ MessageRepositoryImpl │ ChannelRepositoryImpl       │  │
│  │  PluginRepositoryImpl  │ MediaAssetRepositoryImpl                            │  │
│  │                                                                              │  │
│  │  Persistence Objects (po/):                                                  │  │
│  │  SessionPO+TurnPO │ MessagePO │ ChannelPO │ PluginPO │ MediaAssetPO         │  │
│  │  AuditLogPO                                                                  │  │
│  │                                                                              │  │
│  │  Mappers (mapper/):                                                          │  │
│  │  SessionMapper │ MessageMapper │ ChannelMapper │ PluginMapper                │  │
│  │  MediaAssetMapper                                                            │  │
│  └──────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                    │
│  ┌──────────────────────┐  ┌──────────────────────┐  ┌──────────────────────────┐ │
│  │  DataStore           │  │  Migration           │  │  gRPC                    │ │
│  │  (datastore/)        │  │  (migration/)        │  │  (grpc/)                 │ │
│  │  GORM + PostgreSQL   │  │  golang-migrate      │  │  PluginHost              │ │
│  │  Connection Pool     │  │  000001~000006       │  │  PluginBridge            │ │
│  │  DB() *gorm.DB       │  │  .up/.down SQL       │  │  ChannelPlugin gRPC      │ │
│  └──────────────────────┘  └──────────────────────┘  └──────────────────────────┘ │
│                                                                                    │
│  ┌──────────────────────────────────────────────────────────────────────────────┐  │
│  │  ADK - AI Development Kit (adk/)                                             │  │
│  │  OpenAI │ Anthropic │ Gemini │ Groq │ Mistral │ Ollama │ OpenAI-Compat       │  │
│  │  WorkflowEngine │ FallbackChain │ AuthProfile                                │  │
│  └──────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                    │
│  ┌──────────────────────┐  ┌──────────────────────┐  ┌──────────────────────────┐ │
│  │  Services            │  │  Telemetry           │  │  Daemon                  │ │
│  │  (service/)          │  │  (telemetry/)        │  │  (daemon/)               │ │
│  │  JWTService          │  │  OpenTelemetry       │  │  Background process      │ │
│  │  AuditServiceImpl    │  │  Prometheus metrics  │  │  management              │ │
│  └──────────────────────┘  └──────────────────────┘  └──────────────────────────┘ │
└────────────────────────────────────────┬──────────────────────────────────────────┘
                                         │ uses
┌────────────────────────────────────────▼──────────────────────────────────────────┐
│                            Utils Layer  (pkg/utils/)                               │
│                                                                                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │  Config      │  │  Logger      │  │  Container   │  │  BCode               │  │
│  │  (config/)   │  │  (logger/)   │  │  (container/)│  │  (bcode/)            │  │
│  │  Viper-based │  │  Zap-based   │  │  barnettZQG/ │  │  ErrValidationFailed │  │
│  │  AppConfig   │  │  Logger      │  │  inject DI   │  │  ErrUnauthorized     │  │
│  │  ServerConfig│  │  interface   │  │  Provides()  │  │  ErrNotFound         │  │
│  │  DBConfig    │  │              │  │  ProvideWith │  │  ErrInternal         │  │
│  │  RedisConfig │  │              │  │  Name()      │  │  ErrTokenExpired     │  │
│  │  JWTConfig   │  │              │  │  Populate()  │  │  ErrTokenInvalid     │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│                                                                                    │
│  ┌──────────────┐  ┌──────────────┐                                               │
│  │  Errors      │  │  Pprof       │                                               │
│  │  (errors/)   │  │  (pprof/)    │                                               │
│  │  Error types │  │  CPU/Memory  │                                               │
│  │  & helpers   │  │  profiling   │                                               │
│  └──────────────┘  └──────────────┘                                               │
└────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 数据库 Schema（6 张迁移表）

```
migrations/
├── 000001_create_sessions.up.sql      → sessions + turns 表
├── 000002_create_channels.up.sql      → channels 表
├── 000003_create_messages.up.sql      → messages 表
├── 000004_create_plugins.up.sql       → plugins 表
├── 000005_create_audit_logs.up.sql    → audit_logs 表
└── 000006_create_media_assets.up.sql  → media_assets 表
```

---

## gRPC Plugin 协议

```
proto/channel/v1/channel.proto
  ChannelPlugin service:
    rpc Start(StartRequest) → StartResponse
    rpc Stop(StopRequest) → StopResponse
    rpc Status(StatusRequest) → StatusResponse
    rpc SendMessage(SendMessageRequest) → SendMessageResponse
    rpc OnMessage(stream InboundMessageProto) → stream Ack
    rpc GetCapabilities(Empty) → Capabilities
    rpc SetTyping(TypingRequest) → Empty

plugins/
├── discord/    → Discord 频道插件
├── telegram/   → Telegram 频道插件
├── slack/      → Slack 频道插件
├── whatsapp/   → WhatsApp 频道插件
└── signal/     → Signal 频道插件
```

---

## 前端 Web 应用

```
web/
├── control-ui/   → 管理控制台 (React + Vite + Tailwind)
│   └── Pages: Dashboard │ Channels │ Sessions │ Config │ Models │ Plugins │ Security
├── chat-ui/      → 实时对话界面 (React + Vite)
│   └── Components: MessageStream │ ToolExecution │ AgentStatus │ ApprovalDialog
└── canvas/       → Markdown 渲染画布 (React + Vite + Marked)
```

---

## 依赖注入初始化顺序

```
cmd/server/main.go
  └── cli.NewRootCommand() (Cobra)
        └── gateway start
              ├── 1. config.LoadConfig()          → AppConfig (Viper)
              ├── 2. logger.NewLogger()            → Logger (Zap)
              ├── 3. container.NewContainer()      → DI Container
              ├── 4. ProvideWithName("datastore")  → PostgreSQL (GORM)
              ├── 5. ProvideWithName("logger")     → Logger
              ├── 6. Provides(NewSessionRepository, NewMessageRepository, ...)
              ├── 7. Provides(NewRoutingService, NewJWTService, ...)
              ├── 8. Provides(NewAgentAppService, NewGatewayAppService, ...)
              ├── 9. Provides(NewSessionHandler, NewChannelHandler, ...)
              ├── 10. container.Populate()         → 注入所有依赖
              ├── 11. server.NewServer()           → Gin Engine
              ├── 12. server.RegisterMiddleware()  → Auth │ CORS │ Logger
              ├── 13. server.RegisterRoutes()      → /api/v1/* │ /ws │ /metrics
              └── 14. server.Run()                 → ListenAndServe
```

---

## 消息处理流程

```
外部 Plugin (gRPC)
  └── InboundMessageProto
        └── pkg/infrastructure/grpc/PluginHost
              └── GatewayAppService.HandleInboundMessage()
                    ├── HookRegistry.Execute(BeforeMessageReceive)
                    ├── InboundMessage.NewInboundMessage()     [Domain Aggregate]
                    ├── MessageRepository.Create()             [Persist]
                    ├── RoutingService.Resolve()               [Domain Service]
                    ├── InboundMessage.RouteTo(sessionID)      [Domain Method]
                    ├── MessageRepository.Update()             [Persist]
                    ├── EventBus.Publish(MessageReceived)      [Domain Event]
                    └── HookRegistry.Execute(AfterMessageReceive)
```

---

## AI 模型调用流程

```
AgentAppService.ExecuteTurn()
  └── ADK WorkflowEngine
        ├── AuthProfile (API Key 管理)
        ├── Provider 路由:
        │   ├── OpenAI
        │   ├── Anthropic
        │   ├── Gemini
        │   ├── Groq
        │   ├── Mistral
        │   ├── Ollama (本地)
        │   └── OpenAI-Compatible
        └── FallbackChain (故障转移)
```

---

## 技术栈汇总

| 组件 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.25 |
| HTTP 框架 | Gin | v1.12 |
| ORM | GORM | v1.31 |
| 数据库 | PostgreSQL | 14+ |
| 缓存 | Redis | 7+ |
| gRPC | google.golang.org/grpc | v1.58+ |
| 日志 | Zap | v1.27 |
| 配置 | Viper | v1.21 |
| JWT | golang-jwt | v5.3 |
| 迁移 | golang-migrate | v4.19 |
| CLI | Cobra | v1.10 |
| DI 容器 | barnettZQG/inject | latest |
| 可观测性 | OpenTelemetry | - |
| 前端 | React 18 + Vite | - |
