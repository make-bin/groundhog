// @AI_GENERATED
package handler

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/groundhog/pkg/interface/http/response"
	"github.com/make-bin/groundhog/pkg/utils/config"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// HealthHandler defines the health check endpoint.
type HealthHandler interface {
	Check(c *gin.Context)
}

type healthHandler struct {
	cfg    *config.AppConfig
	Logger logger.Logger `inject:"logger"`
}

// NewHealthHandler creates a new HealthHandler instance.
func NewHealthHandler(cfg *config.AppConfig) HealthHandler {
	return &healthHandler{cfg: cfg}
}

// Check returns a deep health check response including database and Redis status.
func (h *healthHandler) Check(c *gin.Context) {
	result := map[string]string{
		"status": "ok",
	}

	if h.cfg != nil {
		result["database"] = checkDBHealth(&h.cfg.Database)
		result["redis"] = checkRedisHealth(h.cfg.Redis.Addr)
	}

	// If any component is unhealthy, return 503.
	for _, v := range result {
		if v != "ok" {
			c.JSON(http.StatusServiceUnavailable, result)
			return
		}
	}

	response.Success(c, result)
}

// checkDBHealth pings the PostgreSQL database and returns "ok" or an error string.
func checkDBHealth(cfg *config.DatabaseConfig) string {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=3",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return "ok"
}

// checkRedisHealth sends a PING to Redis and returns "ok" or an error string.
func checkRedisHealth(addr string) string {
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(3 * time.Second))
	if _, err := fmt.Fprintf(conn, "*1\r\n$4\r\nPING\r\n"); err != nil {
		return fmt.Sprintf("error: %v", err)
	}

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}

	if !strings.HasPrefix(strings.TrimSpace(line), "+PONG") {
		return fmt.Sprintf("unexpected: %s", strings.TrimSpace(line))
	}
	return "ok"
}

// @AI_GENERATED: end
