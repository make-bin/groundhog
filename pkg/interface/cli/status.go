// @AI_GENERATED
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/make-bin/groundhog/pkg/utils/config"
)

// NewStatusCommand creates the "openclaw status" command.
func NewStatusCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current server running status",
		Long:  "Display server running status, active channels, active sessions, loaded plugins, and database connection status.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(configPath)
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")

	return cmd
}

// runStatus loads config, queries the server, and prints a status table.
func runStatus(configPath string) error {
	fmt.Printf("\n%s%sOpenClaw Status%s\n\n", colorBold, colorGreen, colorReset)

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("  %s✗%s  Cannot load config: %v\n\n", colorRed, colorReset, err)
		return err
	}

	baseURL := fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port)
	if cfg.Server.Host == "0.0.0.0" {
		baseURL = fmt.Sprintf("http://127.0.0.1:%d", cfg.Server.Port)
	}

	client := &http.Client{Timeout: 5 * time.Second}

	// Check server health
	serverRunning, dbStatus := checkServerHealth(client, baseURL)

	// Fetch counts (only if server is running)
	channelCount := fetchCount(client, baseURL+"/api/v1/channels")
	sessionCount := fetchCount(client, baseURL+"/api/v1/sessions")
	pluginCount := fetchCount(client, baseURL+"/api/v1/plugins")

	// Print table
	const nameWidth = 26
	const valueWidth = 30

	fmt.Printf("  %-*s  %s\n", nameWidth, "Item", "Value")
	fmt.Printf("  %s  %s\n", repeatChar('─', nameWidth), repeatChar('─', valueWidth))

	// Server running status
	if serverRunning {
		fmt.Printf("  %-*s  %s✓ running%s\n", nameWidth, "Server", colorGreen, colorReset)
	} else {
		fmt.Printf("  %-*s  %s✗ not running%s\n", nameWidth, "Server", colorRed, colorReset)
	}

	// Database status
	if dbStatus == "ok" {
		fmt.Printf("  %-*s  %s✓ connected%s\n", nameWidth, "Database", colorGreen, colorReset)
	} else if dbStatus == "unknown" {
		fmt.Printf("  %-*s  %s⚠ unknown (server offline)%s\n", nameWidth, "Database", colorYellow, colorReset)
	} else {
		fmt.Printf("  %-*s  %s✗ %s%s\n", nameWidth, "Database", colorRed, dbStatus, colorReset)
	}

	// Counts
	printCountRow("Active channels", channelCount, serverRunning, nameWidth)
	printCountRow("Active sessions", sessionCount, serverRunning, nameWidth)
	printCountRow("Loaded plugins", pluginCount, serverRunning, nameWidth)

	fmt.Println()
	return nil
}

// checkServerHealth calls /api/v1/health and returns (serverRunning, dbStatus).
func checkServerHealth(client *http.Client, baseURL string) (bool, string) {
	resp, err := client.Get(baseURL + "/api/v1/health")
	if err != nil {
		return false, "unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "unknown"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return true, "unknown"
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return true, "unknown"
	}

	// Try to extract database status from health response
	dbStatus := "unknown"
	if db, ok := result["database"]; ok {
		if dbStr, ok := db.(string); ok {
			dbStatus = dbStr
		}
	} else if status, ok := result["status"]; ok {
		if statusStr, ok := status.(string); ok && statusStr == "ok" {
			dbStatus = "ok"
		}
	}

	return true, dbStatus
}

// fetchCount calls a list endpoint and returns the count of items, or -1 on error.
func fetchCount(client *http.Client, url string) int {
	resp, err := client.Get(url)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return -1
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1
	}

	// Try array response
	var arr []interface{}
	if err := json.Unmarshal(body, &arr); err == nil {
		return len(arr)
	}

	// Try object with "data" array
	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err == nil {
		if data, ok := obj["data"]; ok {
			if dataArr, ok := data.([]interface{}); ok {
				return len(dataArr)
			}
		}
		// Try "total" field
		if total, ok := obj["total"]; ok {
			switch v := total.(type) {
			case float64:
				return int(v)
			}
		}
	}

	return 0
}

// printCountRow prints a single count row in the status table.
func printCountRow(label string, count int, serverRunning bool, nameWidth int) {
	if !serverRunning {
		fmt.Printf("  %-*s  %s⚠ n/a (server offline)%s\n", nameWidth, label, colorYellow, colorReset)
		return
	}
	if count < 0 {
		fmt.Printf("  %-*s  %s⚠ unavailable%s\n", nameWidth, label, colorYellow, colorReset)
		return
	}
	fmt.Printf("  %-*s  %d\n", nameWidth, label, count)
}

// @AI_GENERATED: end
