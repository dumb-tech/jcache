// Created by ChatGPT o3-mini
package jcache

import (
	"errors"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	jc := New(50*time.Millisecond, 10000)
	defer func(jc *JustCache) {
		_ = jc.Close()
	}(jc)

	err := jc.Set("key1", "value1", 1*time.Second)
	if err != nil {
		t.Fatalf("set value error: %v", err)
	}

	val := jc.Get("key1")
	if val != "value1" {
		t.Errorf("wants 'value1', got '%v'", val)
	}
}

func TestHasAndDel(t *testing.T) {
	jc := New(50*time.Millisecond, 10000)
	defer func(jc *JustCache) {
		_ = jc.Close()
	}(jc)

	_ = jc.Set("key1", 123, 1*time.Second)

	if !jc.Has("key1") {
		t.Error("expected key 'key1' is exists")
	}

	jc.Del("key1")
	if jc.Has("key1") {
		t.Error("expected key 'key1' has been deleted")
	}
}

func TestKeysAndItems(t *testing.T) {
	jc := New(50*time.Millisecond, 10000)
	defer func(jc *JustCache) {
		_ = jc.Close()
	}(jc)

	_ = jc.Set("a", "alpha", 1*time.Second)
	_ = jc.Set("b", "beta", 1*time.Second)

	keys := jc.Keys()
	if len(keys) != 2 {
		t.Errorf("wants keys count 2, got %d", len(keys))
	}

	items := jc.Items()
	if len(items) != 2 {
		t.Errorf("wants items count 2, got %d", len(items))
	}

	foundA, foundB := false, false
	for _, item := range items {
		if item.Key == "a" && item.Value == "alpha" {
			foundA = true
		}
		if item.Key == "b" && item.Value == "beta" {
			foundB = true
		}
	}

	if !foundA || !foundB {
		t.Error("one or both items not found in Items()")
	}
}

func TestClear(t *testing.T) {
	jc := New(50*time.Millisecond, 10000)
	defer func(jc *JustCache) {
		_ = jc.Close()
	}(jc)

	_ = jc.Set("key1", "value1", 1*time.Second)
	_ = jc.Set("key2", "value2", 1*time.Second)

	jc.Clear()
	if len(jc.Keys()) != 0 {
		t.Error("expected empty cache after calling Clear()")
	}
}

func TestExpirationCleanup_OnTheFly(t *testing.T) {
	jc := New(20*time.Millisecond, 10000).WithStrategy(CleanupStrategyOnTheFly)
	defer func(jc *JustCache) {
		_ = jc.Close()
	}(jc)

	_ = jc.Set("temp", "data", 30*time.Millisecond)
	time.Sleep(100 * time.Millisecond)

	jc.Clean(time.Now())

	if jc.Has("temp") {
		t.Error("expected expired item will be removed")
	}
}

func TestExpirationCleanup_Collect(t *testing.T) {
	jc := New(20*time.Millisecond, 10000).WithStrategy(CleanupStrategyCollect)
	defer func(jc *JustCache) {
		_ = jc.Close()
	}(jc)

	_ = jc.Set("temp", "data", 30*time.Millisecond)
	time.Sleep(100 * time.Millisecond)

	jc.Clean(time.Now())

	if jc.Has("temp") {
		t.Error("the expired item was expected to be removed")
	}
}

func TestCapacity(t *testing.T) {
	jc := New(50*time.Millisecond, 10000).WithCapacity(2)
	defer func(jc *JustCache) {
		_ = jc.Close()
	}(jc)

	if err := jc.Set("a", 1, 1*time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := jc.Set("b", 2, 1*time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err := jc.Set("c", 3, 1*time.Second)
	if err == nil {
		t.Error("Error expected when cache capacity exceeded")
	}
	if !errors.Is(err, ErrorCacheIsFull) {
		t.Errorf("wants error '%v', got: %v", ErrorCacheIsFull, err)
	}
}
