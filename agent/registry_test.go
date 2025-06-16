package agent

import (
	"context"
	"testing"
	"github.com/voocel/mas/memory"
	"github.com/voocel/mas/knowledge"
	"github.com/voocel/mas/tools"
)

/**
 * Norwegian-style doc: These tests ensure the agent registry is a safe harbor for all agents—no role is lost, no registration is forgotten, and every lookup is a beacon in the night.
 */

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	a := &mockAgent{name: "captain"}
	r.Register("captain", a)
	got, ok := r.Get("captain")
	if !ok || got.Name() != a.Name() {
		t.Errorf("expected to get registered agent with same name, got %v (ok=%v)", got, ok)
	}
}

func TestRegistry_List(t *testing.T) {
	r := NewRegistry()
	r.Register("a", &mockAgent{name: "a"})
	r.Register("b", &mockAgent{name: "b"})
	list := r.List()
	if len(list) != 2 {
		t.Errorf("expected 2 agents, got %d", len(list))
	}
	if _, ok := list["a"]; !ok {
		t.Error("expected agent 'a' in list")
	}
}

func TestRegistry_RequireRoles(t *testing.T) {
	r := NewRegistry()
	r.Register("a", &mockAgent{name: "a"})
	err := r.RequireRoles("a", "b")
	if err == nil || err.Error() == "" {
		t.Errorf("expected error for missing role, got %v", err)
	}
}

func TestRegistry_HasRole(t *testing.T) {
	r := NewRegistry()
	r.Register("a", &mockAgent{name: "a"})
	if !r.HasRole("a") {
		t.Error("expected HasRole to be true for registered role")
	}
	if r.HasRole("b") {
		t.Error("expected HasRole to be false for missing role")
	}
}

func TestRegistry_Clear(t *testing.T) {
	r := NewRegistry()
	r.Register("a", &mockAgent{name: "a"})
	r.Clear()
	if len(r.List()) != 0 {
		t.Error("expected registry to be empty after Clear")
	}
}

type mockAgent struct{ name string }
func (a *mockAgent) Name() string { return a.name }
func (a *mockAgent) Perceive(_ context.Context, _ interface{}) error { return nil }
func (a *mockAgent) Think(_ context.Context) error { return nil }
func (a *mockAgent) Act(_ context.Context) (interface{}, error) { return nil, nil }
func (a *mockAgent) Process(_ context.Context, _ interface{}) (interface{}, error) { return nil, nil }
func (a *mockAgent) GetMemory() memory.Memory { return nil }
func (a *mockAgent) GetKnowledgeGraph() knowledge.Graph { return nil }
func (a *mockAgent) GetTools() []tools.Tool { return nil }
