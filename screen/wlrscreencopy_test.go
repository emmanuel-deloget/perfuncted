package screen

import (
	"fmt"
	"testing"
	"time"

	"github.com/nskaggs/perfuncted/internal/wl"
)

func TestWithWlrContextCachingAndReset(t *testing.T) {
	origConnect := wlrConnect
	defer func() { wlrConnect = origConnect }()

	// Simple fake connector that returns an empty *wl.Context
	wlrConnect = func(sock string) (*wl.Context, error) { return &wl.Context{}, nil }

	sock := "/tmp/fake-wl-sock"

	// first call should create a context
	var firstPtr, secondPtr *wl.Context
	if err := withWlrContext(sock, func(ctx *wl.Context) error {
		firstPtr = ctx
		return nil
	}); err != nil {
		t.Fatalf("first withWlrContext failed: %v", err)
	}

	// second call should reuse same pointer
	if err := withWlrContext(sock, func(ctx *wl.Context) error {
		secondPtr = ctx
		return nil
	}); err != nil {
		t.Fatalf("second withWlrContext failed: %v", err)
	}

	if firstPtr != secondPtr {
		t.Fatalf("expected same ctx pointer, got different: %p vs %p", firstPtr, secondPtr)
	}

	// simulate failure during fn; cached context should be closed and reset
	if err := withWlrContext(sock, func(ctx *wl.Context) error {
		return fmt.Errorf("simulated")
	}); err == nil {
		t.Fatalf("expected error from simulated fn")
	}

	wlrCacheMu.Lock()
	c := wlrCaches[sock]
	wlrCacheMu.Unlock()
	if c != nil && c.ctx != nil {
		t.Fatalf("expected cached ctx to be nil after error, got %v", c.ctx)
	}
}

func TestWlrCacheJanitorEvicts(t *testing.T) {
	origConnect := wlrConnect
	origTTL := wlrCacheTTL
	defer func() { wlrConnect = origConnect; wlrCacheTTL = origTTL }()

	wlrConnect = func(sock string) (*wl.Context, error) { return &wl.Context{}, nil }
	wlrCacheTTL = 50 * time.Millisecond

	sock := "/tmp/fake-wl-sock-evict"
	if err := withWlrContext(sock, func(ctx *wl.Context) error { return nil }); err != nil {
		t.Fatalf("setup withWlrContext failed: %v", err)
	}

	// mark lastUsed as old
	wlrCacheMu.Lock()
	if c, ok := wlrCaches[sock]; ok {
		c.lastUsed = time.Now().Add(-time.Hour)
	}
	wlrCacheMu.Unlock()

	// wait for janitor to run (it ticks every second)
	time.Sleep(150 * time.Millisecond)

	wlrCacheMu.Lock()
	_, exists := wlrCaches[sock]
	wlrCacheMu.Unlock()
	if exists {
		t.Fatalf("expected cache entry to be evicted by janitor")
	}
}
