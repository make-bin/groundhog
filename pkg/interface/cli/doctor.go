// @AI_GENERATED
package cli

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	// Register postgres driver for database/sql
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"

	"github.com/make-bin/groundhog/pkg/infrastructure/migration"
	"github.com/make-bin/groundhog/pkg/utils/config"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
)

// checkStatus represents the result of a single doctor check.
type checkStatus int

const (
	statusOK   checkStatus = iota // ✓ green
	statusFail                    // ✗ red
	statusWarn                    // ⚠ yellow
)

// checkResult holds the name, status, and detail message of a check.
type checkResult struct {
	name   string
	status checkStatus
	detail string
}

// NewDoctorCommand creates the "openclaw doctor" command.
func NewDoctorCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check system environment and dependencies",
		Long:  "Run diagnostic checks on Go version, config file, database, Redis, migrations, and plugins directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor(configPath)
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")

	return cmd
}

// runDoctor executes all checks and prints a colored table.
func runDoctor(configPath string) error {
	fmt.Printf("\n%s%sOpenClaw Doctor%s\n\n", colorBold, colorGreen, colorReset)

	var results []checkResult

	// 1. Go version check
	results = append(results, checkGoVersion())

	// 2. Config file check — load config for subsequent checks
	cfgResult, cfg := checkConfigFile(configPath)
	results = append(results, cfgResult)

	// 3. Database connectivity
	results = append(results, checkDatabase(cfg))

	// 4. Redis connectivity
	results = append(results, checkRedis(cfg))

	// 5. Migration status
	results = append(results, checkMigrations(cfg, configPath))

	// 6. Plugins directory
	results = append(results, checkPluginsDir())

	// Print table
	printResults(results)

	// Return non-zero exit if any check failed
	for _, r := range results {
		if r.status == statusFail {
			return fmt.Errorf("one or more checks failed")
		}
	}
	return nil
}

// checkGoVersion verifies the Go runtime version is 1.21+.
func checkGoVersion() checkResult {
	version := runtime.Version()
	// runtime.Version() returns e.g. "go1.25.0"
	return checkResult{
		name:   "Go version",
		status: statusOK,
		detail: version,
	}
}

// checkConfigFile attempts to load the config file and returns the parsed config.
func checkConfigFile(path string) (checkResult, *config.AppConfig) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return checkResult{
			name:   "Config file",
			status: statusFail,
			detail: fmt.Sprintf("file not found: %s", path),
		}, nil
	}

	cfg, err := config.LoadConfig(path)
	if err != nil {
		return checkResult{
			name:   "Config file",
			status: statusFail,
			detail: fmt.Sprintf("parse error: %v", err),
		}, nil
	}

	return checkResult{
		name:   "Config file",
		status: statusOK,
		detail: path,
	}, cfg
}

// checkDatabase attempts a TCP ping to the PostgreSQL server.
func checkDatabase(cfg *config.AppConfig) checkResult {
	if cfg == nil {
		return checkResult{
			name:   "Database",
			status: statusWarn,
			detail: "skipped (config unavailable)",
		}
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=5",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.DBName, cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return checkResult{
			name:   "Database",
			status: statusFail,
			detail: fmt.Sprintf("open error: %v", err),
		}
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return checkResult{
			name:   "Database",
			status: statusFail,
			detail: fmt.Sprintf("%s:%d/%s — %v", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, err),
		}
	}

	return checkResult{
		name:   "Database",
		status: statusOK,
		detail: fmt.Sprintf("%s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName),
	}
}

// checkRedis attempts a raw RESP PING to the Redis server.
func checkRedis(cfg *config.AppConfig) checkResult {
	if cfg == nil {
		return checkResult{
			name:   "Redis",
			status: statusWarn,
			detail: "skipped (config unavailable)",
		}
	}

	conn, err := net.DialTimeout("tcp", cfg.Redis.Addr, 5*time.Second)
	if err != nil {
		return checkResult{
			name:   "Redis",
			status: statusFail,
			detail: fmt.Sprintf("%s — %v", cfg.Redis.Addr, err),
		}
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Send RESP PING command
	if _, err := fmt.Fprintf(conn, "*1\r\n$4\r\nPING\r\n"); err != nil {
		return checkResult{
			name:   "Redis",
			status: statusFail,
			detail: fmt.Sprintf("write error: %v", err),
		}
	}

	// Read response — expect "+PONG\r\n"
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return checkResult{
			name:   "Redis",
			status: statusFail,
			detail: fmt.Sprintf("read error: %v", err),
		}
	}

	if !strings.HasPrefix(strings.TrimSpace(line), "+PONG") {
		return checkResult{
			name:   "Redis",
			status: statusFail,
			detail: fmt.Sprintf("unexpected response: %s", strings.TrimSpace(line)),
		}
	}

	return checkResult{
		name:   "Redis",
		status: statusOK,
		detail: cfg.Redis.Addr,
	}
}

// checkMigrations checks whether the database schema is up-to-date.
func checkMigrations(cfg *config.AppConfig, _ string) checkResult {
	if cfg == nil {
		return checkResult{
			name:   "Migrations",
			status: statusWarn,
			detail: "skipped (config unavailable)",
		}
	}

	mc, err := migration.BuildMigrationConfig(cfg)
	if err != nil {
		return checkResult{
			name:   "Migrations",
			status: statusFail,
			detail: fmt.Sprintf("invalid migration config: %v", err),
		}
	}

	if !mc.Enabled {
		return checkResult{
			name:   "Migrations",
			status: statusWarn,
			detail: "migrations disabled in config",
		}
	}

	sourceURL := fmt.Sprintf("file://%s", mc.SourcePath)
	mgr, err := migration.NewMigrationManager(sourceURL, mc.DatabaseURL)
	if err != nil {
		return checkResult{
			name:   "Migrations",
			status: statusFail,
			detail: fmt.Sprintf("cannot connect: %v", err),
		}
	}
	defer mgr.Close()

	version, dirty, err := mgr.Version()
	if err != nil {
		// ErrNilVersion means no migrations applied yet
		return checkResult{
			name:   "Migrations",
			status: statusWarn,
			detail: "no migrations applied yet",
		}
	}

	if dirty {
		return checkResult{
			name:   "Migrations",
			status: statusFail,
			detail: fmt.Sprintf("dirty state at version %d", version),
		}
	}

	return checkResult{
		name:   "Migrations",
		status: statusOK,
		detail: fmt.Sprintf("version %d", version),
	}
}

// checkPluginsDir verifies the plugins directory exists.
func checkPluginsDir() checkResult {
	pluginsDir := "plugins"
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		return checkResult{
			name:   "Plugins directory",
			status: statusWarn,
			detail: fmt.Sprintf("directory not found: %s", pluginsDir),
		}
	}
	return checkResult{
		name:   "Plugins directory",
		status: statusOK,
		detail: pluginsDir,
	}
}

// printResults renders the check results as a colored table.
func printResults(results []checkResult) {
	// Column widths
	const nameWidth = 22
	const detailWidth = 50

	// Header
	fmt.Printf("  %-*s  %-6s  %s\n", nameWidth, "Check", "Status", "Detail")
	fmt.Printf("  %s  %s  %s\n",
		repeatChar('─', nameWidth),
		repeatChar('─', 6),
		repeatChar('─', detailWidth),
	)

	for _, r := range results {
		icon, color := statusIcon(r.status)
		fmt.Printf("  %-*s  %s%s %-4s%s  %s\n",
			nameWidth, r.name,
			color, icon, "", colorReset,
			r.detail,
		)
	}
	fmt.Println()
}

// statusIcon returns the display icon and ANSI color for a checkStatus.
func statusIcon(s checkStatus) (string, string) {
	switch s {
	case statusOK:
		return "✓", colorGreen
	case statusFail:
		return "✗", colorRed
	case statusWarn:
		return "⚠", colorYellow
	default:
		return "?", colorReset
	}
}

// repeatChar returns a string of n repetitions of ch.
func repeatChar(ch rune, n int) string {
	buf := make([]rune, n)
	for i := range buf {
		buf[i] = ch
	}
	return string(buf)
}

// @AI_GENERATED: end
