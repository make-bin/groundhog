// @AI_GENERATED
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/make-bin/groundhog/pkg/utils/config"
)

// callRPC loads config, builds the server URL, and POSTs a JSON-RPC request to /rpc.
// It returns the raw JSON of the "result" field, or an error from the "error" field.
func callRPC(configPath, method string, params interface{}) (json.RawMessage, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	host := cfg.Server.Host
	if host == "0.0.0.0" {
		host = "127.0.0.1"
	}
	baseURL := fmt.Sprintf("http://%s:%d", host, cfg.Server.Port)

	body := map[string]interface{}{
		"method": method,
		"params": params,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(baseURL+"/rpc", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rpcResp struct {
		Result json.RawMessage `json:"result"`
		Error  string          `json:"error"`
	}
	if err := json.Unmarshal(respBytes, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if rpcResp.Error != "" {
		return nil, fmt.Errorf("%s", rpcResp.Error)
	}

	return rpcResp.Result, nil
}

// NewCronCommand creates the "openclaw cron" command group.
func NewCronCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cron",
		Short: "Manage cron jobs",
	}

	cmd.AddCommand(
		newCronListCmd(),
		newCronStatusCmd(),
		newCronAddCmd(),
		newCronUpdateCmd(),
		newCronRemoveCmd(),
		newCronRunCmd(),
		newCronRunsCmd(),
	)

	return cmd
}

// newCronListCmd implements "cron list [--all] [--json]".
func newCronListCmd() *cobra.Command {
	var configPath string
	var all bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List cron jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{
				"include_disabled": all,
			}

			result, err := callRPC(configPath, "cron.list", params)
			if err != nil {
				return err
			}

			if jsonOutput {
				fmt.Println(string(result))
				return nil
			}

			var resp struct {
				Jobs []struct {
					ID      string `json:"id"`
					Name    string `json:"name"`
					Enabled bool   `json:"enabled"`
					Schedule struct {
						Kind    string `json:"kind"`
						At      string `json:"at,omitempty"`
						EveryMs int64  `json:"every_ms,omitempty"`
						Expr    string `json:"expr,omitempty"`
					} `json:"schedule"`
					State struct {
						NextRunAtMs  *int64 `json:"next_run_at_ms"`
						RunningAtMs  *int64 `json:"running_at_ms"`
						LastRunStatus string `json:"last_run_status"`
					} `json:"state"`
				} `json:"jobs"`
				Total int `json:"total"`
			}
			if err := json.Unmarshal(result, &resp); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			const idW = 36
			const nameW = 24
			const enabledW = 8
			const schedW = 28
			const nextRunW = 24
			const statusW = 12

			fmt.Printf("\n%s%sCron Jobs%s  (total: %d)\n\n", colorBold, colorGreen, colorReset, resp.Total)
			fmt.Printf("  %-*s  %-*s  %-*s  %-*s  %-*s  %s\n",
				idW, "ID", nameW, "Name", enabledW, "Enabled", schedW, "Schedule", nextRunW, "Next Run", "Status")
			fmt.Printf("  %s  %s  %s  %s  %s  %s\n",
				repeatChar('─', idW), repeatChar('─', nameW), repeatChar('─', enabledW),
				repeatChar('─', schedW), repeatChar('─', nextRunW), repeatChar('─', statusW))

			for _, j := range resp.Jobs {
				enabledStr := fmt.Sprintf("%s✓%s", colorGreen, colorReset)
				if !j.Enabled {
					enabledStr = fmt.Sprintf("%s✗%s", colorRed, colorReset)
				}

				schedStr := j.Schedule.Kind
				switch j.Schedule.Kind {
				case "at":
					schedStr = "at:" + j.Schedule.At
				case "every":
					schedStr = fmt.Sprintf("every:%dms", j.Schedule.EveryMs)
				case "cron":
					schedStr = "cron:" + j.Schedule.Expr
				}
				if len(schedStr) > schedW {
					schedStr = schedStr[:schedW-1] + "…"
				}

				nextRunStr := "—"
				if j.State.NextRunAtMs != nil {
					t := time.UnixMilli(*j.State.NextRunAtMs)
					nextRunStr = t.Format("2006-01-02 15:04:05")
				}

				statusStr := j.State.LastRunStatus
				if j.State.RunningAtMs != nil {
					statusStr = fmt.Sprintf("%srunning%s", colorYellow, colorReset)
				}
				if statusStr == "" {
					statusStr = "—"
				}

				fmt.Printf("  %-*s  %-*s  %-*s  %-*s  %-*s  %s\n",
					idW, j.ID, nameW, j.Name, enabledW, enabledStr,
					schedW, schedStr, nextRunW, nextRunStr, statusStr)
			}
			fmt.Println()
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")
	cmd.Flags().BoolVar(&all, "all", false, "include disabled jobs")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "output as JSON")

	return cmd
}

// newCronStatusCmd implements "cron status".
func newCronStatusCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show cron scheduler status",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := callRPC(configPath, "cron.status", map[string]interface{}{})
			if err != nil {
				return err
			}

			var resp struct {
				Running     bool   `json:"running"`
				EnabledJobs int    `json:"enabled_jobs"`
				RunningJobs int    `json:"running_jobs"`
				NextRunAtMs *int64 `json:"next_run_at_ms"`
				HeartbeatMs int64  `json:"heartbeat_ms"`
			}
			if err := json.Unmarshal(result, &resp); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			fmt.Printf("\n%s%sCron Scheduler Status%s\n\n", colorBold, colorGreen, colorReset)

			const nameW = 20
			const valW = 30
			fmt.Printf("  %-*s  %s\n", nameW, "Item", "Value")
			fmt.Printf("  %s  %s\n", repeatChar('─', nameW), repeatChar('─', valW))

			runningStr := fmt.Sprintf("%s✗ stopped%s", colorRed, colorReset)
			if resp.Running {
				runningStr = fmt.Sprintf("%s✓ running%s", colorGreen, colorReset)
			}
			fmt.Printf("  %-*s  %s\n", nameW, "Scheduler", runningStr)
			fmt.Printf("  %-*s  %d\n", nameW, "Enabled jobs", resp.EnabledJobs)
			fmt.Printf("  %-*s  %d\n", nameW, "Running jobs", resp.RunningJobs)

			nextRunStr := "—"
			if resp.NextRunAtMs != nil {
				nextRunStr = time.UnixMilli(*resp.NextRunAtMs).Format("2006-01-02 15:04:05")
			}
			fmt.Printf("  %-*s  %s\n", nameW, "Next run at", nextRunStr)
			fmt.Printf("  %-*s  %dms\n", nameW, "Heartbeat", resp.HeartbeatMs)
			fmt.Println()
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")
	return cmd
}

// newCronAddCmd implements "cron add" with all required flags.
func newCronAddCmd() *cobra.Command {
	var configPath string
	var name string
	var at string
	var every string
	var cronExpr string
	var systemEvent string
	var message string
	var session string
	var tz string
	var stagger int64
	var exact int64
	var agent string
	var disabled bool
	var announce bool
	var noDeliver bool
	var channel string
	var to string
	var timeoutSeconds int
	var deleteAfterRun bool

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new cron job",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate mutually exclusive flags
			if cmd.Flags().Changed("stagger") && cmd.Flags().Changed("exact") {
				return fmt.Errorf("Choose either --stagger or --exact, not both")
			}

			// Validate exactly one schedule
			schedCount := 0
			if cmd.Flags().Changed("at") {
				schedCount++
			}
			if cmd.Flags().Changed("every") {
				schedCount++
			}
			if cmd.Flags().Changed("cron") {
				schedCount++
			}
			if schedCount != 1 {
				return fmt.Errorf("Choose exactly one schedule: --at, --every, or --cron")
			}

			// Validate exactly one payload
			payloadCount := 0
			if cmd.Flags().Changed("system-event") {
				payloadCount++
			}
			if cmd.Flags().Changed("message") {
				payloadCount++
			}
			if payloadCount != 1 {
				return fmt.Errorf("Choose exactly one payload: --system-event or --message")
			}

			// Build schedule DTO
			scheduleDTO := map[string]interface{}{}
			if cmd.Flags().Changed("at") {
				scheduleDTO["kind"] = "at"
				scheduleDTO["at"] = at
			} else if cmd.Flags().Changed("every") {
				scheduleDTO["kind"] = "every"
				scheduleDTO["every_ms"] = every
			} else {
				scheduleDTO["kind"] = "cron"
				scheduleDTO["expr"] = cronExpr
				if tz != "" {
					scheduleDTO["tz"] = tz
				}
				if cmd.Flags().Changed("stagger") {
					scheduleDTO["stagger_ms"] = stagger
				}
				if cmd.Flags().Changed("exact") {
					scheduleDTO["stagger_ms"] = exact
				}
			}

			// Build payload DTO
			payloadDTO := map[string]interface{}{}
			if cmd.Flags().Changed("system-event") {
				payloadDTO["kind"] = "systemEvent"
				payloadDTO["text"] = systemEvent
			} else {
				payloadDTO["kind"] = "agentTurn"
				payloadDTO["message"] = message
				if timeoutSeconds > 0 {
					payloadDTO["timeout_seconds"] = timeoutSeconds
				}
			}

			// Build delivery DTO
			var deliveryDTO map[string]interface{}
			if announce {
				deliveryDTO = map[string]interface{}{"mode": "announce"}
				if channel != "" {
					deliveryDTO["channel"] = channel
				}
				if to != "" {
					deliveryDTO["to"] = to
				}
			} else if noDeliver {
				deliveryDTO = map[string]interface{}{"mode": "none"}
			}

			sessionTarget := session
			if sessionTarget == "" {
				sessionTarget = "isolated"
			}

			enabled := !disabled

			params := map[string]interface{}{
				"name":           name,
				"schedule":       scheduleDTO,
				"session_target": sessionTarget,
				"payload":        payloadDTO,
				"enabled":        enabled,
			}
			if agent != "" {
				params["agent_id"] = agent
			}
			if deliveryDTO != nil {
				params["delivery"] = deliveryDTO
			}
			if deleteAfterRun {
				params["delete_after_run"] = true
			}

			result, err := callRPC(configPath, "cron.add", params)
			if err != nil {
				return err
			}

			var resp struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
			if err := json.Unmarshal(result, &resp); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			fmt.Printf("%s✓%s  Created cron job %s%s%s (id: %s)\n",
				colorGreen, colorReset, colorBold, resp.Name, colorReset, resp.ID)
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")
	cmd.Flags().StringVar(&name, "name", "", "job name (required)")
	cmd.Flags().StringVar(&at, "at", "", "ISO8601 timestamp for one-time schedule")
	cmd.Flags().StringVar(&every, "every", "", "interval in ms (e.g. 60000)")
	cmd.Flags().StringVar(&cronExpr, "cron", "", "cron expression")
	cmd.Flags().StringVar(&systemEvent, "system-event", "", "text for system event payload")
	cmd.Flags().StringVar(&message, "message", "", "text for agentTurn payload")
	cmd.Flags().StringVar(&session, "session", "", "session target (default: isolated)")
	cmd.Flags().StringVar(&tz, "tz", "", "timezone for cron schedule")
	cmd.Flags().Int64Var(&stagger, "stagger", 0, "stagger ms for cron schedule")
	cmd.Flags().Int64Var(&exact, "exact", 0, "exact ms offset (mutually exclusive with --stagger)")
	cmd.Flags().StringVar(&agent, "agent", "", "agent ID")
	cmd.Flags().BoolVar(&disabled, "disabled", false, "create job as disabled")
	cmd.Flags().BoolVar(&announce, "announce", false, "delivery mode: announce")
	cmd.Flags().BoolVar(&noDeliver, "no-deliver", false, "delivery mode: none")
	cmd.Flags().StringVar(&channel, "channel", "", "channel for announce delivery")
	cmd.Flags().StringVar(&to, "to", "", "target for announce/webhook delivery")
	cmd.Flags().IntVar(&timeoutSeconds, "timeout-seconds", 0, "timeout for agentTurn (seconds)")
	cmd.Flags().BoolVar(&deleteAfterRun, "delete-after-run", false, "delete job after first run")

	_ = cmd.MarkFlagRequired("name")

	return cmd
}

// newCronUpdateCmd implements "cron update <id>" with optional flags.
func newCronUpdateCmd() *cobra.Command {
	var configPath string
	var name string
	var at string
	var every string
	var cronExpr string
	var systemEvent string
	var message string
	var session string
	var tz string
	var stagger int64
	var exact int64
	var agent string
	var enabled bool
	var disabled bool
	var announce bool
	var noDeliver bool
	var channel string
	var to string
	var timeoutSeconds int
	var deleteAfterRun bool

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an existing cron job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			if cmd.Flags().Changed("stagger") && cmd.Flags().Changed("exact") {
				return fmt.Errorf("Choose either --stagger or --exact, not both")
			}

			schedCount := 0
			if cmd.Flags().Changed("at") {
				schedCount++
			}
			if cmd.Flags().Changed("every") {
				schedCount++
			}
			if cmd.Flags().Changed("cron") {
				schedCount++
			}
			if schedCount > 1 {
				return fmt.Errorf("Choose exactly one schedule: --at, --every, or --cron")
			}

			params := map[string]interface{}{
				"id": id,
			}

			if cmd.Flags().Changed("name") {
				params["name"] = name
			}
			if cmd.Flags().Changed("agent") {
				params["agent_id"] = agent
			}
			if cmd.Flags().Changed("session") {
				params["session_target"] = session
			}
			if cmd.Flags().Changed("enabled") {
				params["enabled"] = enabled
			}
			if cmd.Flags().Changed("disabled") {
				params["enabled"] = !disabled
			}
			if cmd.Flags().Changed("delete-after-run") {
				params["delete_after_run"] = deleteAfterRun
			}

			// Build schedule if any schedule flag was set
			if schedCount == 1 {
				scheduleDTO := map[string]interface{}{}
				if cmd.Flags().Changed("at") {
					scheduleDTO["kind"] = "at"
					scheduleDTO["at"] = at
				} else if cmd.Flags().Changed("every") {
					scheduleDTO["kind"] = "every"
					scheduleDTO["every_ms"] = every
				} else {
					scheduleDTO["kind"] = "cron"
					scheduleDTO["expr"] = cronExpr
					if tz != "" {
						scheduleDTO["tz"] = tz
					}
					if cmd.Flags().Changed("stagger") {
						scheduleDTO["stagger_ms"] = stagger
					}
					if cmd.Flags().Changed("exact") {
						scheduleDTO["stagger_ms"] = exact
					}
				}
				params["schedule"] = scheduleDTO
			}

			// Build payload if any payload flag was set
			if cmd.Flags().Changed("system-event") {
				params["payload"] = map[string]interface{}{
					"kind": "systemEvent",
					"text": systemEvent,
				}
			} else if cmd.Flags().Changed("message") {
				p := map[string]interface{}{
					"kind":    "agentTurn",
					"message": message,
				}
				if timeoutSeconds > 0 {
					p["timeout_seconds"] = timeoutSeconds
				}
				params["payload"] = p
			}

			// Build delivery if any delivery flag was set
			if announce {
				d := map[string]interface{}{"mode": "announce"}
				if channel != "" {
					d["channel"] = channel
				}
				if to != "" {
					d["to"] = to
				}
				params["delivery"] = d
			} else if noDeliver {
				params["delivery"] = map[string]interface{}{"mode": "none"}
			}

			result, err := callRPC(configPath, "cron.update", params)
			if err != nil {
				return err
			}

			var resp struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
			if err := json.Unmarshal(result, &resp); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			fmt.Printf("%s✓%s  Updated cron job %s%s%s (id: %s)\n",
				colorGreen, colorReset, colorBold, resp.Name, colorReset, resp.ID)
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")
	cmd.Flags().StringVar(&name, "name", "", "job name")
	cmd.Flags().StringVar(&at, "at", "", "ISO8601 timestamp for one-time schedule")
	cmd.Flags().StringVar(&every, "every", "", "interval in ms (e.g. 60000)")
	cmd.Flags().StringVar(&cronExpr, "cron", "", "cron expression")
	cmd.Flags().StringVar(&systemEvent, "system-event", "", "text for system event payload")
	cmd.Flags().StringVar(&message, "message", "", "text for agentTurn payload")
	cmd.Flags().StringVar(&session, "session", "", "session target")
	cmd.Flags().StringVar(&tz, "tz", "", "timezone for cron schedule")
	cmd.Flags().Int64Var(&stagger, "stagger", 0, "stagger ms for cron schedule")
	cmd.Flags().Int64Var(&exact, "exact", 0, "exact ms offset (mutually exclusive with --stagger)")
	cmd.Flags().StringVar(&agent, "agent", "", "agent ID")
	cmd.Flags().BoolVar(&enabled, "enabled", false, "enable the job")
	cmd.Flags().BoolVar(&disabled, "disabled", false, "disable the job")
	cmd.Flags().BoolVar(&announce, "announce", false, "delivery mode: announce")
	cmd.Flags().BoolVar(&noDeliver, "no-deliver", false, "delivery mode: none")
	cmd.Flags().StringVar(&channel, "channel", "", "channel for announce delivery")
	cmd.Flags().StringVar(&to, "to", "", "target for announce/webhook delivery")
	cmd.Flags().IntVar(&timeoutSeconds, "timeout-seconds", 0, "timeout for agentTurn (seconds)")
	cmd.Flags().BoolVar(&deleteAfterRun, "delete-after-run", false, "delete job after first run")

	return cmd
}

// newCronRemoveCmd implements "cron remove <id>".
func newCronRemoveCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "remove <id>",
		Short: "Remove a cron job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			result, err := callRPC(configPath, "cron.remove", map[string]interface{}{"id": id})
			if err != nil {
				return err
			}

			var resp struct {
				Removed bool `json:"removed"`
			}
			if err := json.Unmarshal(result, &resp); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			if resp.Removed {
				fmt.Printf("%s✓%s  Removed cron job %s\n", colorGreen, colorReset, id)
			} else {
				fmt.Printf("%s⚠%s  Job %s was not found\n", colorYellow, colorReset, id)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")
	return cmd
}

// newCronRunCmd implements "cron run <id> [--force]".
func newCronRunCmd() *cobra.Command {
	var configPath string
	var force bool

	cmd := &cobra.Command{
		Use:   "run <id>",
		Short: "Trigger a cron job immediately",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			mode := "due"
			if force {
				mode = "force"
			}

			result, err := callRPC(configPath, "cron.run", map[string]interface{}{
				"id":   id,
				"mode": mode,
			})
			if err != nil {
				return err
			}

			var resp struct {
				Triggered bool `json:"triggered"`
			}
			if err := json.Unmarshal(result, &resp); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			if resp.Triggered {
				fmt.Printf("%s✓%s  Triggered cron job %s\n", colorGreen, colorReset, id)
			} else {
				fmt.Printf("%s⚠%s  Job %s was not triggered (already running or not due)\n", colorYellow, colorReset, id)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")
	cmd.Flags().BoolVar(&force, "force", false, "force run even if not due")
	return cmd
}

// newCronRunsCmd implements "cron runs [--job <id>] [--limit n] [--offset n] [--status s] [--json]".
func newCronRunsCmd() *cobra.Command {
	var configPath string
	var jobID string
	var limit int
	var offset int
	var status string
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "runs",
		Short: "List cron job run logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{
				"offset": offset,
				"limit":  limit,
			}
			if jobID != "" {
				params["job_id"] = jobID
			}
			if status != "" {
				params["statuses"] = []string{status}
			}

			result, err := callRPC(configPath, "cron.runs", params)
			if err != nil {
				return err
			}

			if jsonOutput {
				fmt.Println(string(result))
				return nil
			}

			var resp struct {
				Logs []struct {
					ID         int64  `json:"id"`
					JobID      string `json:"job_id"`
					Status     string `json:"status"`
					Action     string `json:"action"`
					RunAtMs    int64  `json:"run_at_ms"`
					DurationMs int64  `json:"duration_ms"`
					Error      string `json:"error,omitempty"`
					Summary    string `json:"summary,omitempty"`
				} `json:"logs"`
				Total int `json:"total"`
			}
			if err := json.Unmarshal(result, &resp); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			const idW = 8
			const jobW = 36
			const statusW = 8
			const actionW = 12
			const runAtW = 20
			const durW = 12

			fmt.Printf("\n%s%sCron Run Logs%s  (total: %d)\n\n", colorBold, colorGreen, colorReset, resp.Total)
			fmt.Printf("  %-*s  %-*s  %-*s  %-*s  %-*s  %-*s  %s\n",
				idW, "ID", jobW, "Job ID", statusW, "Status", actionW, "Action",
				runAtW, "Run At", durW, "Duration", "Error/Summary")
			fmt.Printf("  %s  %s  %s  %s  %s  %s  %s\n",
				repeatChar('─', idW), repeatChar('─', jobW), repeatChar('─', statusW),
				repeatChar('─', actionW), repeatChar('─', runAtW), repeatChar('─', durW),
				repeatChar('─', 20))

			for _, l := range resp.Logs {
				statusColor := colorReset
				if l.Status == "ok" {
					statusColor = colorGreen
				} else if l.Status == "error" {
					statusColor = colorRed
				}

				runAt := time.UnixMilli(l.RunAtMs).Format("2006-01-02 15:04:05")
				durStr := fmt.Sprintf("%dms", l.DurationMs)

				note := l.Summary
				if l.Error != "" {
					note = l.Error
				}
				if len(note) > 40 {
					note = note[:39] + "…"
				}

				fmt.Printf("  %-*d  %-*s  %s%-*s%s  %-*s  %-*s  %-*s  %s\n",
					idW, l.ID, jobW, l.JobID,
					statusColor, statusW, l.Status, colorReset,
					actionW, l.Action, runAtW, runAt, durW, durStr, note)
			}
			fmt.Println()
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")
	cmd.Flags().StringVar(&jobID, "job", "", "filter by job ID")
	cmd.Flags().IntVar(&limit, "limit", 20, "number of results to return")
	cmd.Flags().IntVar(&offset, "offset", 0, "offset for pagination")
	cmd.Flags().StringVar(&status, "status", "", "filter by status (ok/error)")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "output as JSON")

	return cmd
}

// @AI_GENERATED: end
