// @AI_GENERATED
package adk

import (
	"context"
	"fmt"
	"sync"

	"github.com/make-bin/groundhog/pkg/domain/conversation/aggregate/agent_session"
)

// WorkflowType defines the type of workflow agent.
type WorkflowType string

const (
	WorkflowSequential WorkflowType = "sequential"
	WorkflowLoop       WorkflowType = "loop"
	WorkflowParallel   WorkflowType = "parallel"
)

// SubAgentConfig defines a sub-agent in a workflow.
type SubAgentConfig struct {
	AgentID      string
	WorkflowType WorkflowType
	MaxLoops     int // for Loop workflow
}

// WorkflowRunner orchestrates multiple sub-agents.
type WorkflowRunner struct {
	RunnerAdapter *RunnerAdapter `inject:""`
}

// NewWorkflowRunner creates a new WorkflowRunner.
func NewWorkflowRunner() *WorkflowRunner {
	return &WorkflowRunner{}
}

// RunSequential runs sub-agents sequentially, passing output of each as input to the next.
func (w *WorkflowRunner) RunSequential(ctx context.Context, sessions []*agent_session.AgentSession, initialInput string) (string, error) {
	input := initialInput
	for _, sess := range sessions {
		ch, err := w.RunnerAdapter.Run(ctx, sess, input)
		if err != nil {
			return "", fmt.Errorf("sequential workflow: agent %s: %w", sess.ID().Value(), err)
		}
		var lastResponse string
		for turn := range ch {
			lastResponse = turn.Response()
		}
		input = lastResponse
	}
	return input, nil
}

// RunLoop runs a single agent in a loop until maxLoops or context cancellation.
func (w *WorkflowRunner) RunLoop(ctx context.Context, sess *agent_session.AgentSession, initialInput string, maxLoops int) (string, error) {
	input := initialInput
	for i := 0; i < maxLoops; i++ {
		select {
		case <-ctx.Done():
			return input, ctx.Err()
		default:
		}
		ch, err := w.RunnerAdapter.Run(ctx, sess, input)
		if err != nil {
			return "", fmt.Errorf("loop workflow iteration %d: %w", i, err)
		}
		var lastResponse string
		for turn := range ch {
			lastResponse = turn.Response()
		}
		input = lastResponse
	}
	return input, nil
}

// RunParallel runs multiple agents in parallel and collects all results.
func (w *WorkflowRunner) RunParallel(ctx context.Context, sessions []*agent_session.AgentSession, input string) ([]string, error) {
	results := make([]string, len(sessions))
	errs := make([]error, len(sessions))
	var wg sync.WaitGroup

	for i, sess := range sessions {
		wg.Add(1)
		go func(idx int, s *agent_session.AgentSession) {
			defer wg.Done()
			ch, err := w.RunnerAdapter.Run(ctx, s, input)
			if err != nil {
				errs[idx] = err
				return
			}
			var lastResponse string
			for turn := range ch {
				lastResponse = turn.Response()
			}
			results[idx] = lastResponse
		}(i, sess)
	}

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, fmt.Errorf("parallel workflow: %w", err)
		}
	}
	return results, nil
}

// @AI_GENERATED: end
