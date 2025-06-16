package orchestrator

import (
	"context"
	"fmt"
	"testing"
	"time"
	"github.com/voocel/mas/agent"
	"github.com/voocel/mas/knowledge"
	"github.com/voocel/mas/memory"
	"github.com/voocel/mas/tools"
)

/**
 * Norwegian-style doc: Like a harbor master guiding ships, these tests ensure every orchestrator operation is predictable, safe, and free from storms. All error and edge cases are tested, so every task reaches its destination.
 */

/**
 * Norwegian-style doc: These tests ensure the orchestrator handles every storm—from missing agents to failing tasks, from double registration to race conditions. No ship is lost, and every error is a beacon of clarity.
 */

func TestRegisterAgent_Duplicate(t *testing.T) {
	orch := NewBasicOrchestrator(Options{})
	a := agent.NewBaseAgent("A")
	_ = orch.RegisterAgent(a)
	err := orch.RegisterAgent(a)
	if err == nil {
		t.Errorf("expected error on duplicate agent registration, got nil")
	}
}

func TestSubmitTask_MissingAgent(t *testing.T) {
	orch := NewBasicOrchestrator(Options{})
	orch.Start()
	task := Task{AgentIDs: []string{"missing"}}
	_, err := orch.SubmitTask(context.Background(), task)
	if err == nil || err.Error() == "" {
		t.Errorf("expected error for missing agent, got %v", err)
	}
}

func TestCancelTask_NotFoundAndCompleted(t *testing.T) {
	orch := NewBasicOrchestrator(Options{})
	orch.Start()
	err := orch.CancelTask(context.Background(), "notfound")
	if err == nil || err.Error() == "" {
		t.Errorf("expected error for cancelling non-existent task, got %v", err)
	}
	// Completed task cannot be cancelled
	a := agent.NewBaseAgent("A")
	_ = orch.RegisterAgent(a)
	task := Task{AgentIDs: []string{"A"}}
	id, _ := orch.SubmitTask(context.Background(), task)
	time.Sleep(20 * time.Millisecond)
	_ = orch.CancelTask(context.Background(), id) // Should not panic
	got, _ := orch.GetTask(id)
	if got.Status == TaskStatusCancelled {
		t.Errorf("should not be able to cancel completed task")
	}
}

func TestStatusTransitions(t *testing.T) {
	orch := NewBasicOrchestrator(Options{})
	orch.Start()
	a := &successAgent{}
	_ = orch.RegisterAgent(a)
	task := Task{AgentIDs: []string{"success"}, Input: "input"}
	id, err := orch.SubmitTask(context.Background(), task)
	if err != nil {
		t.Fatalf("failed to submit task: %v", err)
	}
	var lastStatus TaskStatus
	for i := 0; i < 10; i++ {
		task, _ := orch.GetTask(id)
		if task.Status != lastStatus {
			lastStatus = task.Status
		}
		if task.Status == TaskStatusCompleted || task.Status == TaskStatusFailed {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if lastStatus != TaskStatusCompleted {
		t.Errorf("expected task to complete, got status %v", lastStatus)
	}
}

type successAgent struct{}
func (a *successAgent) Name() string { return "success" }
func (a *successAgent) Perceive(ctx context.Context, input interface{}) error { return nil }
func (a *successAgent) Think(ctx context.Context) error { return nil }
func (a *successAgent) Act(ctx context.Context) (interface{}, error) { return "done", nil }
func (a *successAgent) Process(ctx context.Context, input interface{}) (interface{}, error) { return "done", nil }
func (a *successAgent) GetMemory() memory.Memory { return nil }
func (a *successAgent) GetKnowledgeGraph() knowledge.Graph { return nil }
func (a *successAgent) GetTools() []tools.Tool { return nil }


func TestAgentErrorPropagation(t *testing.T) {
	orch := NewBasicOrchestrator(Options{})
	orch.Start()
	testAgent := &errorAgent{}
	_ = orch.RegisterAgent(testAgent)
	task := Task{AgentIDs: []string{"err"}, Input: "fail"}
	id, _ := orch.SubmitTask(context.Background(), task)
	time.Sleep(20 * time.Millisecond)
	taskResult, _ := orch.GetTask(id)
	if taskResult.Status != TaskStatusFailed {
		t.Errorf("expected failed status, got %v", taskResult.Status)
	}
	if taskResult.Error == "" {
		t.Errorf("expected error message, got empty string")
	}
}

func TestOrchestrator_Concurrency(t *testing.T) {
	orch := NewBasicOrchestrator(Options{})
	orch.Start()
	a := &successAgent{}
	_ = orch.RegisterAgent(a)
	done := make(chan struct{})
	n := 10
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		go func(idx int) {
			task := Task{AgentIDs: []string{"success"}, Input: idx}
			id, err := orch.SubmitTask(context.Background(), task)
			if err == nil {
				ids[idx] = id
			}
			if idx == n-1 {
				close(done)
			}
		}(i)
	}
	<-done
	// Wait for all tasks to reach a terminal state (completed or failed)
	for _, id := range ids {
		if id == "" {
			continue
		}
		for i := 0; i < 20; i++ {
			task, _ := orch.GetTask(id)
			if task.Status == TaskStatusCompleted || task.Status == TaskStatusFailed {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	// Check that all tasks completed
	for _, id := range ids {
		if id == "" {
			continue
		}
		task, _ := orch.GetTask(id)
		if task.Status != TaskStatusCompleted {
			t.Errorf("expected completed status for task %s, got %v", id, task.Status)
		}
	}
}

type errorAgent struct{}
func (a *errorAgent) Name() string { return "err" }
func (a *errorAgent) Perceive(ctx context.Context, input interface{}) error { return fmt.Errorf("perceive error") }
func (a *errorAgent) Think(ctx context.Context) error { return fmt.Errorf("think error") }
func (a *errorAgent) Act(ctx context.Context) (interface{}, error) { return nil, fmt.Errorf("act error") }
func (a *errorAgent) Process(ctx context.Context, input interface{}) (interface{}, error) { return nil, fmt.Errorf("process error") }
func (a *errorAgent) GetMemory() memory.Memory { return nil }
func (a *errorAgent) GetKnowledgeGraph() knowledge.Graph { return nil }
func (a *errorAgent) GetTools() []tools.Tool { return nil }

func TestNewBasicOrchestrator_Defaults(t *testing.T) {
	orch := NewBasicOrchestrator(Options{})
	if orch == nil {
		t.Fatal("expected orchestrator instance")
	}
	if orch.ttl <= 0 {
		t.Error("expected positive ttl")
	}
	if orch.pollInt <= 0 {
		t.Error("expected positive poll interval")
	}
}

func TestRegisterAndGetAgent(t *testing.T) {
	orch := NewBasicOrchestrator(Options{})
	a := agent.NewBaseAgent("A")
	err := orch.RegisterAgent(a)
	if err != nil {
		t.Fatalf("failed to register agent: %v", err)
	}
	retrieved, err := orch.GetAgent("A")
	if err != nil || retrieved.Name() != "A" {
		t.Errorf("expected agent 'A', got %v (err: %v)", retrieved, err)
	}
}

func TestSubmitAndGetTask(t *testing.T) {
	orch := NewBasicOrchestrator(Options{})
	orch.Start()
	a := agent.NewBaseAgent("A")
	orch.RegisterAgent(a)
	task := Task{
		AgentIDs: []string{"A"},
		Input:    "input",
	}
	id, err := orch.SubmitTask(context.Background(), task)
	if err != nil {
		t.Fatalf("failed to submit task: %v", err)
	}
	got, err := orch.GetTask(id)
	if err != nil {
		t.Fatalf("failed to get task: %v", err)
	}
	if got.ID != id {
		t.Errorf("expected task ID %s, got %s", id, got.ID)
	}
}

func TestSubmitTask_NotRunning(t *testing.T) {
	orch := NewBasicOrchestrator(Options{})
	task := Task{AgentIDs: []string{"A"}}
	_, err := orch.SubmitTask(context.Background(), task)
	if err == nil || err.Error() != "orchestrator not running" {
		t.Errorf("expected 'not running' error, got %v", err)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	orch := NewBasicOrchestrator(Options{})
	_, err := orch.GetTask("missing")
	if err == nil || err.Error() != "task not found" {
		t.Errorf("expected 'task not found' error, got %v", err)
	}
}

func TestCancelTask(t *testing.T) {
	orch := NewBasicOrchestrator(Options{})
	orch.Start()

	// Custom agent whose Act blocks until closed
	done := make(chan struct{})
	testAgent := &blockingAgent{done: done}
	orch.RegisterAgent(testAgent)
	task := Task{AgentIDs: []string{"blocker"}, Input: "input"}
	id, _ := orch.SubmitTask(context.Background(), task)
	// Wait briefly to ensure task is started
	time.Sleep(20 * time.Millisecond)
	err := orch.CancelTask(context.Background(), id)
	close(done) // Unblock Act if needed
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

type blockingAgent struct{ done chan struct{} }
func (a *blockingAgent) Name() string { return "blocker" }
func (a *blockingAgent) Perceive(ctx context.Context, input interface{}) error { return nil }
func (a *blockingAgent) Think(ctx context.Context) error { return nil }
func (a *blockingAgent) Act(ctx context.Context) (interface{}, error) {
	select {
	case <-a.done:
		return "cancelled", nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
func (a *blockingAgent) Process(ctx context.Context, input interface{}) (interface{}, error) {
	if err := a.Perceive(ctx, input); err != nil { return nil, err }
	if err := a.Think(ctx); err != nil { return nil, err }
	return a.Act(ctx)
}
func (a *blockingAgent) GetMemory() memory.Memory { return nil }
func (a *blockingAgent) GetKnowledgeGraph() knowledge.Graph { return nil }
func (a *blockingAgent) GetTools() []tools.Tool { return nil }

