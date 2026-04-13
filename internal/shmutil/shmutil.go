// Package shmutil provides shared-memory file creation for Wayland protocols.
// Both the screen capture and virtual keyboard backends need anonymous temp
// files in XDG_RUNTIME_DIR for wl_shm buffers and XKB keymaps. This package
// deduplicates that logic.
package shmutil

import (
	"fmt"
	"os"
)

// CreateFile creates an anonymous temp file of the given size in
// XDG_RUNTIME_DIR, suitable for wl_shm or XKB keymap file descriptors.
// The file is unlinked immediately so it disappears when the last fd closes.
func CreateFile(size int64) (*os.File, error) {
	dir := os.Getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		return nil, fmt.Errorf("XDG_RUNTIME_DIR not set")
	}
	f, err := os.CreateTemp(dir, "perfuncted-shm-*")
	if err != nil {
		return nil, err
	}
	if err := f.Truncate(size); err != nil {
		f.Close()
		return nil, err
	}
	os.Remove(f.Name()) //nolint:errcheck
	return f, nil
}
