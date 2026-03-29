// @AI_GENERATED
package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/make-bin/groundhog/pkg/application/approval"
	"github.com/make-bin/groundhog/pkg/application/hook"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/adk"
	agentinfra "github.com/make-bin/groundhog/pkg/infrastructure/agent"
	cronfra "github.com/make-bin/groundhog/pkg/infrastructure/cron"
	"github.com/make-bin/groundhog/pkg/infrastructure/daemon"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence"
	skillpkg "github.com/make-bin/groundhog/pkg/infrastructure/skill"
	"github.com/make-bin/groundhog/pkg/utils/pprof"
	"github.com/spf13/cobra"

	"github.com/make-bin/groundhog/pkg/application/eventbus"
	appservice "github.com/make-bin/groundhog/pkg/application/service"
	cron_service "github.com/make-bin/groundhog/pkg/domain/cron/service"
	msgservice "github.com/make-bin/groundhog/pkg/domain/messaging/service"
	"github.com/make-bin/groundhog/pkg/infrastructure/datastore"
	"github.com/make-bin/groundhog/pkg/infrastructure/migration"
	"github.com/make-bin/groundhog/pkg/infrastructure/service"
	"github.com/make-bin/groundhog/pkg/interface/http/handler"
	"github.com/make-bin/groundhog/pkg/interface/http/middleware"
	wshandler "github.com/make-bin/groundhog/pkg/interface/ws"
	"github.com/make-bin/groundhog/pkg/server"
	"github.com/make-bin/groundhog/pkg/utils/config"
	"github.com/make-bin/groundhog/pkg/utils/container"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)


// NewGatewayCommand creates the "openclaw gateway" command group.
func NewGatewayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "Gateway server commands",
	}
	cmd.AddCommand(NewGatewayRunCommand())
	cmd.AddCommand(NewGatewayStopCommand())
	return cmd
}

// NewGatewayRunCommand creates the "openclaw gateway run" subcommand.
func NewGatewayRunCommand() *cobra.Command {
	var configPath string
	var daemonMode bool
	var pidFile string
	var logFile string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the gateway HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if daemonMode {
				exe, err := os.Executable()
				if err != nil {
					return fmt.Errorf("cannot determine executable path: %w", err)
				}
				childArgs := []string{"gateway", "run", "--config", configPath}
				d := daemon.New(pidFile, logFile)
				return d.Start(exe, childArgs)
			}
			return runGateway(configPath)
		},
	}

	cmd.PersistentFlags().StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")
	cmd.Flags().BoolVar(&daemonMode, "daemon", false, "run as background daemon process")
	cmd.Flags().StringVar(&pidFile, "pid-file", "openclaw.pid", "path to PID file (daemon mode)")
	cmd.Flags().StringVar(&logFile, "log-file", "openclaw.log", "path to log file (daemon mode)")

	return cmd
}

// NewGatewayStopCommand creates the "openclaw gateway stop" subcommand.
func NewGatewayStopCommand() *cobra.Command {
	var pidFile string

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a running gateway daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			d := daemon.New(pidFile, "")
			return d.Stop()
		},
	}

	cmd.Flags().StringVar(&pidFile, "pid-file", "openclaw.pid", "path to PID file")
	return cmd
}


// runGateway contains the actual server startup logic.
func runGateway(configPath string) error {
	// 1. Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize logger
	log, err := logger.NewLogger(logger.LogConfig{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// 2a. Start pprof if enabled
	if cfg.Pprof.Enabled {
		pprof.EnableHTTP(cfg.Pprof.Addr)
	}

	// 3. Run database migrations
	if err := migration.RunMigrations(cfg, log); err != nil {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// 4. Create DataStore
	ds, err := datastore.NewDataStore(&cfg.Database, log)
	if err != nil {
		log.Error("failed to create datastore", "error", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := ds.Close(); closeErr != nil {
			log.Warn("failed to close datastore", "error", closeErr)
		}
	}()

	// 5. Create DI container
	c := container.NewContainer()

	if err := c.ProvideWithName("datastore", ds); err != nil {
		log.Error("failed to register datastore", "error", err)
		os.Exit(1)
	}
	if err := c.ProvideWithName("logger", log); err != nil {
		log.Error("failed to register logger", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(service.NewJWTService(&cfg.JWT)); err != nil {
		log.Error("failed to register jwt service", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(eventbus.NewEventBus()); err != nil {
		log.Error("failed to register event bus", "error", err)
		os.Exit(1)
	}


	// Step 6: Infrastructure repositories
	if err := c.Provides(persistence.NewSessionRepository()); err != nil {
		log.Error("failed to register session repository", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(persistence.NewMemoryRepository()); err != nil {
		log.Error("failed to register memory repository", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(persistence.NewMessageRepository()); err != nil {
		log.Error("failed to register message repository", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(persistence.NewChannelRepository()); err != nil {
		log.Error("failed to register channel repository", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(persistence.NewPluginRepository()); err != nil {
		log.Error("failed to register plugin repository", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(persistence.NewMediaAssetRepository()); err != nil {
		log.Error("failed to register media asset repository", "error", err)
		os.Exit(1)
	}
	// Cron infrastructure repositories (Task 16.1)
	if err := c.Provides(persistence.NewCronJobRepositoryImpl()); err != nil {
		log.Error("failed to register cron job repository", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(persistence.NewCronRunLogRepositoryImpl()); err != nil {
		log.Error("failed to register cron run log repository", "error", err)
		os.Exit(1)
	}

	// Step 7: Domain services
	agentRegistry := agentinfra.NewRegistry(cfg)
	log.Info("agents loaded", "count", len(agentRegistry.List()), "default", agentRegistry.DefaultID())

	agentRoutingSvc := agentinfra.NewAgentRoutingService(agentRegistry, cfg.Bindings, log)
	if err := c.Provides(agentRoutingSvc); err != nil {
		log.Error("failed to register agent routing service", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(msgservice.NewMessageChunkingService()); err != nil {
		log.Error("failed to register message chunking service", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(msgservice.NewCommandDetectionService()); err != nil {
		log.Error("failed to register command detection service", "error", err)
		os.Exit(1)
	}
	// Cron domain scheduler service (Task 16.2)
	if err := c.Provides(cron_service.NewCronSchedulerService()); err != nil {
		log.Error("failed to register cron scheduler service", "error", err)
		os.Exit(1)
	}


	// Step 8: Application services
	if err := c.Provides(appservice.NewGatewayAppService()); err != nil {
		log.Error("failed to register gateway app service", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(appservice.NewChannelAppService()); err != nil {
		log.Error("failed to register channel app service", "error", err)
		os.Exit(1)
	}

	modelAdapter := adk.NewModelAdapter()
	if err := registerModelProviders(modelAdapter, cfg); err != nil {
		log.Error("failed to register model providers", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(modelAdapter); err != nil {
		log.Error("failed to register model adapter", "error", err)
		os.Exit(1)
	}
	toolAdapter := adk.NewToolAdapter()
	if err := c.Provides(toolAdapter); err != nil {
		log.Error("failed to register tool adapter", "error", err)
		os.Exit(1)
	}

	if cfg.MCP.Enabled && len(cfg.MCP.Servers) > 0 {
		mcpMgr := adk.NewMCPManager(log)
		mcpCfgs := make([]adk.MCPServerCfg, 0, len(cfg.MCP.Servers))
		for _, s := range cfg.MCP.Servers {
			mcpCfgs = append(mcpCfgs, adk.MCPServerCfg{
				Name:    s.Name,
				Command: s.Command,
				Args:    s.Args,
				Env:     s.Env,
			})
		}
		mcpMgr.Start(context.Background(), mcpCfgs)
		toolAdapter.RegisterMCPTools(mcpMgr)
		defer mcpMgr.Close()
		log.Info("mcp tools registered", "count", len(mcpMgr.Tools()))
	}
	if err := c.Provides(adk.NewSessionAdapter()); err != nil {
		log.Error("failed to register session adapter", "error", err)
		os.Exit(1)
	}
	runnerAdapter := adk.NewRunnerAdapter()
	skillCfg := skillpkg.DefaultLoadConfig(".", cfg.Skills.ExtraDirs)
	if cfg.Skills.Dir != "" {
		skillCfg.WorkspaceDir = cfg.Skills.Dir
	}
	skillRegistry := skillpkg.NewRegistry(skillCfg)
	log.Info("skills loaded", "count", len(skillRegistry.List()))
	runnerAdapter.SkillRegistry = skillRegistry
	if err := c.Provides(runnerAdapter); err != nil {
		log.Error("failed to register runner adapter", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(adk.NewCompactionService()); err != nil {
		log.Error("failed to register compaction service", "error", err)
		os.Exit(1)
	}
	if cfg.Memory.Enabled {
		embeddingProvider := adk.NewOpenAICompatEmbeddingProvider(
			cfg.Memory.EmbeddingBaseURL,
			cfg.Memory.EmbeddingAPIKey,
			cfg.Memory.EmbeddingModel,
		)
		if err := c.ProvideWithName("embedding_provider", embeddingProvider); err != nil {
			log.Error("failed to register embedding provider", "error", err)
			os.Exit(1)
		}
		if err := c.Provides(appservice.NewMemoryAppService()); err != nil {
			log.Error("failed to register memory app service", "error", err)
			os.Exit(1)
		}
	}
	agentAppSvcVal := appservice.NewAgentAppService(cfg, agentRegistry)
	if err := c.Provides(agentAppSvcVal); err != nil {
		log.Error("failed to register agent app service", "error", err)
		os.Exit(1)
	}
	// Cron application service (Task 16.2)
	if err := c.Provides(appservice.NewCronAppService()); err != nil {
		log.Error("failed to register cron app service", "error", err)
		os.Exit(1)
	}


	// Step 9: Interface handlers
	healthHandler := handler.NewHealthHandler(cfg)
	if err := c.Provides(healthHandler); err != nil {
		log.Error("failed to register health handler", "error", err)
		os.Exit(1)
	}
	sessionHandler := handler.NewSessionHandler()
	if err := c.Provides(sessionHandler); err != nil {
		log.Error("failed to register session handler", "error", err)
		os.Exit(1)
	}
	wsHandler := wshandler.NewWSEventHandler()
	if err := c.Provides(wsHandler); err != nil {
		log.Error("failed to register ws event handler", "error", err)
		os.Exit(1)
	}
	channelHandler := handler.NewChannelHandler()
	if err := c.Provides(channelHandler); err != nil {
		log.Error("failed to register channel handler", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(hook.NewHookRegistry()); err != nil {
		log.Error("failed to register hook registry", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(approval.NewManager()); err != nil {
		log.Error("failed to register approval manager", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(service.NewAuditService()); err != nil {
		log.Error("failed to register audit service", "error", err)
		os.Exit(1)
	}
	if err := c.Provides(adk.NewWorkflowRunner()); err != nil {
		log.Error("failed to register workflow runner", "error", err)
		os.Exit(1)
	}
	securityHandler := handler.NewSecurityHandler()
	if err := c.Provides(securityHandler); err != nil {
		log.Error("failed to register security handler", "error", err)
		os.Exit(1)
	}
	memoryHandler := handler.NewMemoryHandler()
	if err := c.Provides(memoryHandler); err != nil {
		log.Error("failed to register memory handler", "error", err)
		os.Exit(1)
	}
	agentHandler := handler.NewAgentHandler()
	if err := c.Provides(agentHandler); err != nil {
		log.Error("failed to register agent handler", "error", err)
		os.Exit(1)
	}
	// Cron RPC handler (Task 16.3)
	if err := c.Provides(wshandler.NewCronRPCHandler()); err != nil {
		log.Error("failed to register cron rpc handler", "error", err)
		os.Exit(1)
	}
	rpcRouter := wshandler.NewRPCRouter()
	if err := c.Provides(rpcRouter); err != nil {
		log.Error("failed to register rpc router", "error", err)
		os.Exit(1)
	}

	// Step 10: Populate container — wires all inject tags
	if err := c.Populate(); err != nil {
		log.Error("failed to populate container", "error", err)
		os.Exit(1)
	}

	// Wire AgentRoutingService.SessionCreator after Populate (avoids circular DI).
	// agentAppSvcVal holds the concrete value returned by NewAgentAppService;
	// after Populate() its inject fields are fully wired, so it can create sessions.
	agentRoutingSvc.SetSessionCreator(agentAppSvcVal)


	// Step 11: Wire Scheduler and Reaper into CronAppService (Task 16.3)
	// Scheduler and Reaper live in infrastructure/cron which imports application/service,
	// so they cannot be registered in the DI container (circular import).
	// After Populate(), the CronAppService has all its repos/services injected.
	// We retrieve them via getter methods and build the Scheduler/Reaper manually.
	gatewayCtx, gatewayCancel := context.WithCancel(context.Background())
	defer gatewayCancel()

	if rpcRouter.CronHandler != nil && rpcRouter.CronHandler.CronAppSvc != nil {
		cronSvc := rpcRouter.CronHandler.CronAppSvc

		deliveryExec := cronfra.NewDeliveryExecutor(log)
		jobExec := cronfra.NewJobExecutor(
			cronSvc.GetAgentAppSvc(),
			cronSvc.GetSessionRepo(),
			cronSvc.GetChannelAppSvc(),
			log,
		)
		scheduler := cronfra.NewScheduler(
			cronSvc.GetCronRepo(),
			cronSvc.GetRunLogRepo(),
			cronSvc.GetSchedulerSvc(),
			jobExec,
			deliveryExec,
			0,
			log,
		)
		reaper := cronfra.NewReaper(cronSvc.GetSessionRepo(), 0, 0, true, log)

		cronSvc.SetScheduler(scheduler)
		cronSvc.SetReaper(reaper)

		// Step 12: Start CronAppService
		if err := cronSvc.Start(gatewayCtx); err != nil {
			log.Error("failed to start cron app service", "error", err)
			os.Exit(1)
		}
		defer func() {
			stopCtx, stopCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer stopCancel()
			if stopErr := cronSvc.Stop(stopCtx); stopErr != nil {
				log.Warn("cron app service stop error", "error", stopErr)
			}
		}()
	}

	// Step 13: Create server, register middleware and routes
	srv := server.NewServer(&cfg.Server, log)
	srv.RegisterMiddleware(
		middleware.CORSMiddleware(),
		middleware.LoggerMiddleware(log),
		middleware.UserIDMiddleware(),
	)
	srv.RegisterRoutes(healthHandler, sessionHandler, channelHandler, wsHandler, securityHandler, memoryHandler, agentHandler, rpcRouter)

	// Step 14: Start HTTP server
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Run()
	}()

	log.Info("gateway server started", "host", cfg.Server.Host, "port", cfg.Server.Port)

	// Step 15: Graceful shutdown on signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Info("received shutdown signal", "signal", sig.String())
	case err := <-errCh:
		if err != nil {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown error", "error", err)
		os.Exit(1)
	}

	log.Info("gateway server stopped gracefully")
	return nil
}


// registerModelProviders wires provider factories from config into the ModelAdapter.
func registerModelProviders(adapter *adk.ModelAdapter, cfg *config.AppConfig) error {
	providers := cfg.Models.Providers

	if p, ok := providers["openai_compat"]; ok && p.BaseURL != "" {
		apiKey := ""
		if len(p.APIKeys) > 0 {
			apiKey = p.APIKeys[0]
		}
		adapter.RegisterProvider(vo.ProviderOpenAICompat, adk.OpenAICompatProviderFactory(p.BaseURL, apiKey))
	}

	if p, ok := providers["ollama"]; ok && p.BaseURL != "" {
		adapter.RegisterProvider(vo.ProviderOllama, adk.OllamaProviderFactory())
	}

	if p, ok := providers["groq"]; ok && len(p.APIKeys) > 0 {
		adapter.RegisterProvider(vo.ProviderGroq, adk.GroqProviderFactory())
	}

	return nil
}

// @AI_GENERATED: end
