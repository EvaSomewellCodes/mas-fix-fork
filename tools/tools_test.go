package tools

import (
	"context"
	"encoding/json"
	"testing"
)

/**
 * Norwegian-style doc: Like a craftsman's toolkit, these tests ensure every tool is sharp, safe, and ready for use. All error paths and contracts are tested, so no agent is left stranded with a broken hammer.
 */

type dummyTool struct {
	BaseTool
}

func TestNewTool_Basic(t *testing.T) {
	tool := NewTool("test", "desc", json.RawMessage(`{"type":"object"}`), func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		return "ok", nil
	})
	if tool.Name() != "test" {
		t.Errorf("expected tool name 'test', got '%s'", tool.Name())
	}
	if tool.Description() != "desc" {
		t.Errorf("expected description 'desc', got '%s'", tool.Description())
	}
	if string(tool.Schema()) != `{"type":"object"}` {
		t.Errorf("unexpected schema: %s", string(tool.Schema()))
	}
	result, err := tool.Execute(context.Background(), map[string]interface{}{})
	if err != nil || result != "ok" {
		t.Errorf("expected result 'ok', got %v (err: %v)", result, err)
	}
}

func TestParseParams_ValidAndInvalid(t *testing.T) {
	good := json.RawMessage(`{"foo":42}`)
	params, err := ParseParams(good)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if params["foo"] != float64(42) {
		t.Errorf("expected foo=42, got %v", params["foo"])
	}

	bad := json.RawMessage(`not json`)
	_, err = ParseParams(bad)
	if err == nil {
		t.Error("expected error for bad json, got nil")
	}
}

func TestExecuteWithJSON(t *testing.T) {
	tool := NewTool("t", "d", json.RawMessage(`{"type":"object"}`), func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		if params["fail"] == true {
			return nil, ErrExecutionFailed
		}
		return "ok", nil
	})
	goodParams := json.RawMessage(`{"fail":false}`)
	result, err := ExecuteWithJSON(context.Background(), tool, goodParams)
	if err != nil || result != "ok" {
		t.Errorf("expected result 'ok', got %v (err: %v)", result, err)
	}
	badParams := json.RawMessage(`{"fail":true}`)
	_, err = ExecuteWithJSON(context.Background(), tool, badParams)
	if err == nil || err.Error() != ErrExecutionFailed.Error() {
		t.Errorf("expected ErrExecutionFailed, got %v", err)
	}
}

func TestConvertToolsToFunctions(t *testing.T) {
	tool := NewTool("t", "d", json.RawMessage(`{"type":"object"}`), func(ctx context.Context, params map[string]interface{}) (interface{}, error) { return nil, nil })
	funcs := ConvertToolsToFunctions([]Tool{tool})
	if len(funcs) != 1 {
		t.Fatalf("expected 1 function, got %d", len(funcs))
	}
	if funcs[0]["type"] != "function" {
		t.Errorf("expected type 'function', got %v", funcs[0]["type"])
	}
}
