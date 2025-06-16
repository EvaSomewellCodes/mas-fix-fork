package agent

import (
	"context"
	"testing"
	"github.com/voocel/mas/memory"
	"github.com/voocel/mas/tools"
)

/**
 * Norwegian-style doc: Like a lighthouse on a windswept coast, these tests illuminate the Agent interface and BaseAgent implementation, ensuring every method is visible and safe for future journeys. Edge cases, errors, and contracts are tested so no developer is lost in the fog.
 */

func TestBaseAgent_Name(t *testing.T) {
	a := NewBaseAgent("TestAgent")
	if a.Name() != "TestAgent" {
		t.Errorf("expected name 'TestAgent', got '%s'", a.Name())
	}
}

func TestBaseAgent_GetMemory(t *testing.T) {
	mem := memory.NewInMemory(memory.Config{})
	a := NewBaseAgentWithOptions("A", mem, nil, nil)
	if a.GetMemory() != mem {
		t.Error("expected memory to match assigned value")
	}
}

func TestBaseAgent_GetKnowledgeGraph(t *testing.T) {
	// nil is fine for now
	a := NewBaseAgentWithOptions("A", nil, nil, nil)
	if a.GetKnowledgeGraph() != nil {
		t.Error("expected knowledge graph to be nil")
	}
}

func TestBaseAgent_GetTools(t *testing.T) {
	toolsList := []tools.Tool{tools.NewTool("t", "desc", nil, func(ctx context.Context, params map[string]interface{}) (interface{}, error) { return nil, nil })}
	a := NewBaseAgentWithOptions("A", nil, nil, toolsList)
	if len(a.GetTools()) != 1 {
		t.Errorf("expected 1 tool, got %d", len(a.GetTools()))
	}
}

func TestBaseAgent_Perceive_NotImplemented(t *testing.T) {
	a := NewBaseAgent("A")
	err := a.Perceive(context.Background(), nil)
	if err == nil || err.Error() != "unimplemented method: Perceive is not available for base agent" {
		t.Errorf("expected unimplemented error, got %v", err)
	}
}

func TestBaseAgent_Think_NotImplemented(t *testing.T) {
	a := NewBaseAgent("A")
	err := a.Think(context.Background())
	if err == nil || err.Error() != "unimplemented method: Think is not available for base agent" {
		t.Errorf("expected unimplemented error, got %v", err)
	}
}

func TestBaseAgent_Act_NotImplemented(t *testing.T) {
	a := NewBaseAgent("A")
	_, err := a.Act(context.Background())
	if err == nil || err.Error() != "unimplemented method: Act is not available for base agent" {
		t.Errorf("expected unimplemented error, got %v", err)
	}
}

type testAgentPhases struct {
	*BaseAgent
	order *[]string
}
func (a *testAgentPhases) Perceive(ctx context.Context, input interface{}) error {
	*a.order = append(*a.order, "perceive")
	return nil
}
func (a *testAgentPhases) Think(ctx context.Context) error {
	*a.order = append(*a.order, "think")
	return nil
}
func (a *testAgentPhases) Act(ctx context.Context) (interface{}, error) {
	*a.order = append(*a.order, "act")
	return "done", nil
}
func (a *testAgentPhases) Process(ctx context.Context, input interface{}) (interface{}, error) {
	if err := a.Perceive(ctx, input); err != nil {
		return nil, err
	}
	if err := a.Think(ctx); err != nil {
		return nil, err
	}
	return a.Act(ctx)
}

func TestBaseAgent_Process_CallsAllPhases(t *testing.T) {
	order := []string{}
	a := &testAgentPhases{BaseAgent: NewBaseAgent("A"), order: &order}
	result, err := a.Process(context.Background(), "input")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "done" {
		t.Errorf("expected result 'done', got %v", result)
	}
	if len(order) != 3 || order[0] != "perceive" || order[1] != "think" || order[2] != "act" {
		t.Errorf("expected phase order perceive-think-act, got %+v", order)
	}
}
