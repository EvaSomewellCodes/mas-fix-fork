package agency

import (
	"context"
	"testing"
	"github.com/voocel/mas/agent"

)

/**
 * Norwegian-style doc: Like a conductor gathering musicians before a symphony, this suite ensures every Agency method performs in harmony, revealing flaws before they become cacophony in production. All edge cases and errors are tested, so no mystery lingers in the fjords of our codebase.
 */

func TestNewAgency_Defaults(t *testing.T) {
	ag := New(Config{Name: "TestAgency"})
	if ag.Name != "TestAgency" {
		t.Errorf("expected agency name to be 'TestAgency', got '%s'", ag.Name)
	}
	if ag.FlowChart == nil {
		t.Errorf("expected default flowchart to be initialized")
	}
	if ag.Orchestrator == nil {
		t.Errorf("expected default orchestrator to be initialized")
	}
}

func TestAgency_AddAndGetAgent(t *testing.T) {
	ag := New(Config{Name: "TestAgency"})
	mockAgent := agent.NewBaseAgent("A1")
	err := ag.AddAgent(mockAgent)
	if err != nil {
		t.Fatalf("failed to add agent: %v", err)
	}
	retrieved, err := ag.GetAgent("A1")
	if err != nil {
		t.Fatalf("failed to get agent: %v", err)
	}
	if retrieved.Name() != "A1" {
		t.Errorf("expected agent name 'A1', got '%s'", retrieved.Name())
	}
	// Adding same agent again should error
	err = ag.AddAgent(mockAgent)
	if err == nil {
		t.Error("expected error when adding duplicate agent, got nil")
	}
}

func TestAgency_ListAgents(t *testing.T) {
	ag := New(Config{Name: "TestAgency"})
	if len(ag.ListAgents()) != 0 {
		t.Error("expected no agents initially")
	}
	ag.AddAgent(agent.NewBaseAgent("A1"))
	ag.AddAgent(agent.NewBaseAgent("A2"))
	agents := ag.ListAgents()
	if len(agents) != 2 {
		t.Errorf("expected 2 agents, got %d", len(agents))
	}
}

func TestAgency_FlowChart(t *testing.T) {
	ag := New(Config{Name: "TestAgency"})
	f := NewFlowChart()
	f.AddEntryPoint("A1")
	ag.SetFlowChart(f)
	if ag.FlowChart.EntryPoints[0] != "A1" {
		t.Errorf("expected entry point to be 'A1'")
	}
}

func TestAgency_DefineFlowChart_Invalid(t *testing.T) {
	ag := New(Config{Name: "TestAgency"})
	err := ag.DefineFlowChart([]Flow{{}})
	if err == nil {
		t.Error("expected error for invalid flow definition")
	}
}

func TestAgency_Execute_NoEntryPoint(t *testing.T) {
	ag := New(Config{Name: "TestAgency"})
	_, err := ag.Execute(context.Background(), "input")
	if err == nil || err.Error() != "no entry point defined in the agency" {
		t.Errorf("expected 'no entry point defined' error, got %v", err)
	}
}

func TestAgency_Execute_EntryPointNotFound(t *testing.T) {
	ag := New(Config{Name: "TestAgency"})
	f := NewFlowChart()
	f.AddEntryPoint("A1")
	ag.SetFlowChart(f)
	_, err := ag.Execute(context.Background(), "input")
	if err == nil || err.Error() != "agent with ID A1 not found" {
		t.Errorf("expected 'agent with ID A1 not found' error, got %v", err)
	}
}

func TestAgency_RegisterWorkflow_NotImplemented(t *testing.T) {
	ag := New(Config{Name: "TestAgency"})
	err := ag.RegisterWorkflow(nil)
	if err == nil || err.Error() != "workflow support not yet implemented" {
		t.Errorf("expected not implemented error, got %v", err)
	}
}
