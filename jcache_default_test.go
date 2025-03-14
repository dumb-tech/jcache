// Created by ChatGPT o3-mini
package jcache

import (
	"errors"
	"testing"
	"time"
)

func TestDefaultCacheSetAndGet(t *testing.T) {
	jc := Default()
	defer func() {
		_ = jc.Close()
	}()

	err := jc.Set("defaultKey", "defaultValue", 1*time.Second)
	if err != nil {
		t.Fatalf("failed to set value: %v", err)
	}

	val := jc.Get("defaultKey")
	if val != "defaultValue" {
		t.Errorf("expected 'defaultValue', got '%v'", val)
	}
}

func TestDefaultCacheHasAndDel(t *testing.T) {
	jc := Default()
	defer func() {
		_ = jc.Close()
	}()

	_ = jc.Set("defaultKey", 456, 1*time.Second)

	if !jc.Has("defaultKey") {
		t.Error("expected key 'defaultKey' to exist")
	}

	jc.Del("defaultKey")
	if jc.Has("defaultKey") {
		t.Error("expected key 'defaultKey' to be deleted")
	}
}

func TestDefaultCacheClear(t *testing.T) {
	jc := Default()
	defer func() {
		_ = jc.Close()
	}()

	_ = jc.Set("key1", "value1", 1*time.Second)
	_ = jc.Set("key2", "value2", 1*time.Second)

	jc.Clear()
	if len(jc.Keys()) != 0 {
		t.Error("expected empty cache after calling Clear()")
	}
}

func TestDefaultCacheExpirationCleanup(t *testing.T) {
	jc := Default().WithStrategy(CleanupStrategyOnTheFly)
	defer func() {
		_ = jc.Close()
	}()

	_ = jc.Set("temp", "data", 30*time.Millisecond)
	time.Sleep(100 * time.Millisecond)

	jc.Clean(time.Now())
	if jc.Has("temp") {
		t.Error("expected expired item to be removed")
	}
}

func TestDefaultCacheCapacity(t *testing.T) {
	jc := Default().WithCapacity(2)
	defer func() {
		_ = jc.Close()
	}()

	if err := jc.Set("a", 1, 1*time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := jc.Set("b", 2, 1*time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err := jc.Set("c", 3, 1*time.Second)
	if err == nil {
		t.Error("expected error when cache capacity is exceeded")
	}
	if !errors.Is(err, ErrorCacheIsFull) {
		t.Errorf("expected error '%v', got: %v", ErrorCacheIsFull, err)
	}
}
