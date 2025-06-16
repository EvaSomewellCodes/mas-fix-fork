package agent

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"github.com/voocel/mas/llm"
	"github.com/voocel/mas/memory"
	"github.com/voocel/mas/tools"
)

/**
 * Norwegian-style doc: These tests ensure every LLM agent phase is clear as a fjord at sunrise—inputs are remembered, thoughts are formed, tools are called, and errors are never lost in the fog.
 */

func TestNewLLMAgent_Basics(t *testing.T) {
	agent := NewLLMAgent(LLMAgentConfig{Name: "llm"})
	if agent.Name() != "llm" {
		t.Errorf("expected agent name 'llm', got %s", agent.Name())
	}
}

func TestLLMAgent_PerceiveAndMemory(t *testing.T) {
	agent := NewLLMAgent(LLMAgentConfig{Name: "llm"})
	err := agent.Perceive(context.Background(), "input")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLLMAgent_Act_ToolNotFound(t *testing.T) {
	agent := NewLLMAgent(LLMAgentConfig{Name: "llm"})
	_, err := agent.callTool(context.Background(), "missing", nil)
	if err == nil || err.Error() == "" {
		t.Errorf("expected error for missing tool, got %v", err)
	}
}

func TestLLMAgent_ParseToolCall(t *testing.T) {
	name, params, err := parseToolCall("Tool:sum\nParameters:{\"a\":1,\"b\":2}")
	if err != nil || name != "sum" || params["a"] != float64(1) {
		t.Errorf("expected to parse tool call, got name=%s, params=%v, err=%v", name, params, err)
	}
}

func TestLLMAgent_ParseToolCall_MissingName(t *testing.T) {
	_, _, err := parseToolCall("Parameters:{}")
	if err == nil {
		t.Error("expected error for missing tool name")
	}
}

func TestLLMAgent_IsToolCall(t *testing.T) {
	if !isToolCall("Tool:sum") {
		t.Error("expected true for tool call string")
	}
	if isToolCall("no tool here") {
		t.Error("expected false for non-tool string")
	}
}

// --- Advanced Coverage Tests ---

type mockProvider struct{
	fail bool
	resp *llm.ChatCompletionResponse
	err  error
}

func (m *mockProvider) ID() string { return "mock" }
func (m *mockProvider) ChatCompletion(ctx context.Context, req llm.ChatCompletionRequest) (*llm.ChatCompletionResponse, error) {
	if m.fail { return nil, errors.New("provider error") }
	if m.resp != nil { return m.resp, nil }
	return &llm.ChatCompletionResponse{
		Choices: []llm.Choice{{Message: llm.Message{Content: "Thought: use a tool\nTool:adder\nParameters:{\"x\":1,\"y\":2}"}}},
	}, nil
}
func (m *mockProvider) GetModels(ctx context.Context) ([]string, error) { return nil, nil }
func (m *mockProvider) Close() error { return nil }

type mockTool struct{
	fail bool
}
func (m *mockTool) Name() string { return "adder" }
func (m *mockTool) Description() string { return "adds numbers" }
func (m *mockTool) Schema() json.RawMessage { return json.RawMessage(`{"type":"object"}`) }
func (m *mockTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	if m.fail { return nil, errors.New("tool error") }
	return params["x"].(float64) + params["y"].(float64), nil
}

func TestLLMAgent_Think_ProviderError(t *testing.T) {
	agent := NewLLMAgent(LLMAgentConfig{Name: "llm"})
	agent.provider = &mockProvider{fail: true}
	err := agent.Think(context.Background())
	if err == nil || err.Error() == "" {
		t.Error("expected error from provider failure")
	}
}

func TestLLMAgent_Act_ToolSuccess(t *testing.T) {
	agent := NewLLMAgent(LLMAgentConfig{Name: "llm", Tools: []tools.Tool{&mockTool{}}})
	agent.provider = &mockProvider{}
	agent.currentThought = "Tool:adder\nParameters:{\"x\":1,\"y\":2}"
	result, err := agent.Act(context.Background())
	if err != nil || result != float64(3) {
		t.Errorf("expected sum result 3, got %v, err=%v", result, err)
	}
}

func TestLLMAgent_Act_ToolError(t *testing.T) {
	agent := NewLLMAgent(LLMAgentConfig{Name: "llm", Tools: []tools.Tool{&mockTool{fail: true}}})
	agent.provider = &mockProvider{}
	agent.currentThought = "Tool:adder\nParameters:{\"x\":1,\"y\":2}"
	_, err := agent.Act(context.Background())
	if err == nil || err.Error() == "" {
		t.Error("expected error from tool execution failure")
	}
}

func TestLLMAgent_Process_EndToEnd(t *testing.T) {
	agent := NewLLMAgent(LLMAgentConfig{Name: "llm", Tools: []tools.Tool{&mockTool{}}})
	agent.provider = &mockProvider{}
	result, err := agent.Process(context.Background(), "input")
	if err != nil {
		t.Errorf("unexpected error in end-to-end process: %v", err)
	}
	if result == nil {
		t.Error("expected non-nil result from process")
	}
}

func TestLLMAgent_preparePrompt(t *testing.T) {
	agent := NewLLMAgent(LLMAgentConfig{Name: "llm", Tools: []tools.Tool{&mockTool{}}})
	agent.currentInput = map[string]interface{}{ "foo": "bar" }
	prompt := agent.preparePrompt()
	if prompt == "" || prompt == "\n\n" {
		t.Error("expected non-empty prompt")
	}
}

type badMemory struct{ memory.Memory }
func (b *badMemory) Add(ctx context.Context, item memory.MemoryItem) error { return errors.New("memory add error") }

func TestLLMAgent_MemoryError(t *testing.T) {
	agent := NewLLMAgent(LLMAgentConfig{Name: "llm"})
	agent.memory = &badMemory{}
	err := agent.Perceive(context.Background(), "input")
	if err == nil || err.Error() == "" {
		t.Error("expected error from memory add failure")
	}
}
