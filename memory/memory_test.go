package memory

import (
	"context"
	"testing"
	"time"
)

/**
 * Norwegian-style doc: Like a journal kept by a wise traveler, these tests ensure every memory is recorded, found, and erased as expected. No recollection is lost, and all errors are revealed to the morning sun.
 */

/**
 * Norwegian-style doc: VectorStore, the memory of the future, must be as reliable as the mind of Odin. These tests ensure every vector is stored, found, and forgotten with precision, and that even the empty void returns the right echo.
 */

func TestVectorStore_AddAndGet(t *testing.T) {
	store := NewVectorStore(Config{})
	item := MemoryItem{ID: "v1", Content: "vector hello", Type: TypeObservation, CreatedAt: time.Now()}
	store.Add(context.Background(), item)
	got, err := store.Get(context.Background(), "v1")
	if err != nil || got.ID != "v1" {
		t.Errorf("expected to get item with ID 'v1', got %v (err: %v)", got, err)
	}
}

func TestVectorStore_Search(t *testing.T) {
	store := NewVectorStore(Config{})
	store.Add(context.Background(), MemoryItem{ID: "v1", Content: "vector hello", Type: TypeObservation, CreatedAt: time.Now()})
	store.Add(context.Background(), MemoryItem{ID: "v2", Content: "other content", Type: TypeObservation, CreatedAt: time.Now()})
	results, err := store.Search(context.Background(), "vector", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, item := range results {
		if item.ID == "v1" {
			found = true
			break
		}
	}
	if !found || len(results) != 1 {
		t.Errorf("expected 1 result with ID 'v1', got %+v", results)
	}
}

func TestVectorStore_GetRecent(t *testing.T) {
	store := NewVectorStore(Config{})
	store.Add(context.Background(), MemoryItem{ID: "v1", Content: "a", Type: TypeObservation, CreatedAt: time.Now()})
	store.Add(context.Background(), MemoryItem{ID: "v2", Content: "b", Type: TypeObservation, CreatedAt: time.Now()})
	recent, err := store.GetRecent(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recent) != 1 || recent[0].ID != "v2" {
		t.Errorf("expected most recent ID 'v2', got %+v", recent)
	}
}

func TestVectorStore_Clear(t *testing.T) {
	store := NewVectorStore(Config{})
	store.Add(context.Background(), MemoryItem{ID: "v1", Content: "a", Type: TypeObservation, CreatedAt: time.Now()})
	store.Clear(context.Background())
	if len(store.items) != 0 {
		t.Errorf("expected items to be cleared, got %d", len(store.items))
	}
}

func TestVectorStore_Get_NotFound(t *testing.T) {
	store := NewVectorStore(Config{})
	_, err := store.Get(context.Background(), "missing")
	if err == nil || err.Error() != "memory item not found" {
		t.Errorf("expected 'memory item not found' error, got %v", err)
	}
}

func TestVectorStore_Search_EmptyAndNoMatch(t *testing.T) {
	store := NewVectorStore(Config{})
	results, err := store.Search(context.Background(), "nothing", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty store, got %+v", results)
	}
	store.Add(context.Background(), MemoryItem{ID: "v1", Content: "abc", Type: TypeObservation, CreatedAt: time.Now()})
	results, err = store.Search(context.Background(), "zzz", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for no match, got %+v", results)
	}
}

func TestVectorStore_GetRecent_EmptyAndOverLimit(t *testing.T) {
	store := NewVectorStore(Config{})
	recent, err := store.GetRecent(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recent) != 0 {
		t.Errorf("expected 0 results for empty store, got %+v", recent)
	}
	// Add fewer items than requested
	store.Add(context.Background(), MemoryItem{ID: "v1", Content: "a", Type: TypeObservation, CreatedAt: time.Now()})
	store.Add(context.Background(), MemoryItem{ID: "v2", Content: "b", Type: TypeObservation, CreatedAt: time.Now()})
	recent, err = store.GetRecent(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recent) != 2 {
		t.Errorf("expected 2 results for over-limit request, got %+v", recent)
	}
}

func TestInMemory_AddAndGet(t *testing.T) {
	mem := NewInMemory(Config{Capacity: 2})
	item := MemoryItem{ID: "1", Content: "hello", Type: TypeObservation, CreatedAt: time.Now()}
	mem.Add(context.Background(), item)
	got, err := mem.Get(context.Background(), "1")
	if err != nil || got.ID != "1" {
		t.Errorf("expected to get item with ID '1', got %v (err: %v)", got, err)
	}
}

func TestInMemory_Capacity(t *testing.T) {
	mem := NewInMemory(Config{Capacity: 2})
	mem.Add(context.Background(), MemoryItem{ID: "1", Content: "a", Type: TypeObservation, CreatedAt: time.Now()})
	mem.Add(context.Background(), MemoryItem{ID: "2", Content: "b", Type: TypeObservation, CreatedAt: time.Now()})
	mem.Add(context.Background(), MemoryItem{ID: "3", Content: "c", Type: TypeObservation, CreatedAt: time.Now()})
	if len(mem.items) != 2 {
		t.Errorf("expected capacity 2, got %d", len(mem.items))
	}
	if mem.items[0].ID != "2" || mem.items[1].ID != "3" {
		t.Errorf("expected IDs '2', '3', got %v", []string{mem.items[0].ID, mem.items[1].ID})
	}
}

func TestInMemory_Search(t *testing.T) {
	mem := NewInMemory(Config{})
	mem.Add(context.Background(), MemoryItem{ID: "1", Content: "foo bar", Type: TypeObservation, CreatedAt: time.Now()})
	mem.Add(context.Background(), MemoryItem{ID: "2", Content: "baz qux", Type: TypeObservation, CreatedAt: time.Now()})
	results, err := mem.Search(context.Background(), "foo", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, item := range results {
		if item.ID == "1" {
			found = true
			break
		}
	}
	if !found || len(results) != 1 {
		t.Errorf("expected 1 result with ID '1', got %+v", results)
	}
}

func TestInMemory_GetRecent(t *testing.T) {
	mem := NewInMemory(Config{})
	mem.Add(context.Background(), MemoryItem{ID: "1", Content: "a", Type: TypeObservation, CreatedAt: time.Now()})
	mem.Add(context.Background(), MemoryItem{ID: "2", Content: "b", Type: TypeObservation, CreatedAt: time.Now()})
	recent, err := mem.GetRecent(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recent) != 1 || recent[0].ID != "2" {
		t.Errorf("expected most recent ID '2', got %+v", recent)
	}
}

func TestInMemory_Clear(t *testing.T) {
	mem := NewInMemory(Config{})
	mem.Add(context.Background(), MemoryItem{ID: "1", Content: "a", Type: TypeObservation, CreatedAt: time.Now()})
	mem.Clear(context.Background())
	if len(mem.items) != 0 {
		t.Errorf("expected items to be cleared, got %d", len(mem.items))
	}
}

func TestInMemory_Get_NotFound(t *testing.T) {
	mem := NewInMemory(Config{})
	_, err := mem.Get(context.Background(), "missing")
	if err == nil || err.Error() != "memory item not found" {
		t.Errorf("expected 'memory item not found' error, got %v", err)
	}
}
