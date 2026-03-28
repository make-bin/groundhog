// @AI_GENERATED
package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/groundhog/pkg/infrastructure/telemetry"
	"github.com/make-bin/groundhog/pkg/interface/http/handler"
	"github.com/make-bin/groundhog/pkg/interface/http/router"
	"github.com/make-bin/groundhog/pkg/interface/ws"
	"github.com/make-bin/groundhog/pkg/utils/config"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// telemetryMetricsHandler returns the metrics HTTP handler (lazy init).
func telemetryMetricsHandler() http.Handler {
	if m := telemetry.Global(); m != nil {
		return telemetry.MetricsHandler()
	}
	return telemetry.MetricsHandler()
}

// Server wraps a Gin engine with configuration and logging.
type Server struct {
	engine     *gin.Engine
	httpServer *http.Server
	cfg        *config.ServerConfig
	logger     logger.Logger
}

// NewServer creates a new Server instance with the given config and logger.
func NewServer(cfg *config.ServerConfig, logger logger.Logger) *Server {
	engine := gin.New()
	return &Server{
		engine: engine,
		cfg:    cfg,
		logger: logger,
	}
}

// RegisterMiddleware registers the given middleware handlers on the engine.
// Recommended order: Recovery → CORS → Logger.
func (s *Server) RegisterMiddleware(middlewares ...gin.HandlerFunc) {
	s.engine.Use(middlewares...)
}

// RegisterRoutes registers application routes on the engine.
func (s *Server) RegisterRoutes(healthHandler handler.HealthHandler, sessionHandler handler.SessionHandler, channelHandler handler.ChannelHandler, wsHandler ws.WSEventHandler, securityHandler handler.SecurityHandler, memoryHandler handler.MemoryHandler, rpcRouter *ws.RPCRouter) {
	v1 := s.engine.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Check)
		router.RegisterSessionRoutes(v1, sessionHandler)
		router.RegisterChannelRoutes(v1, channelHandler)
		router.RegisterSecurityRoutes(v1, securityHandler)
		router.RegisterMemoryRoutes(v1, memoryHandler)
	}
	// WebSocket / SSE event streaming endpoint
	s.engine.GET("/ws", wsHandler.Handle)

	// JSON-RPC endpoint for cron and other RPC methods
	s.engine.POST("/rpc", rpcRouter.Handle)

	// Prometheus-compatible metrics endpoint
	s.engine.GET("/metrics", gin.WrapH(telemetryMetricsHandler()))
}

// RegisterStaticFiles serves the embedded unified frontend app.
// Route: /* → app/dist (SPA with client-side routing)
func (s *Server) RegisterStaticFiles(assets embed.FS) {
	sub, err := fs.Sub(assets, "app/dist")
	if err != nil {
		s.logger.Warn("failed to create sub-filesystem for frontend", "error", err)
		return
	}
	fileServer := http.FileServer(http.FS(sub))
	// Serve static assets directly; fall back to index.html for SPA routing
	s.engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Try to serve the file; if not found, serve index.html
		f, err := sub.(fs.FS).Open(path[1:]) // strip leading /
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}
		// SPA fallback
		c.Request.URL.Path = "/"
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

// Run starts the HTTP server using net/http.Server with configured timeouts.
func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	s.logger.Info("starting HTTP server", "addr", addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the HTTP server within the given context deadline.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	s.logger.Info("shutting down HTTP server")

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
	}

	return s.httpServer.Shutdown(ctx)
}

// @AI_GENERATED: end
