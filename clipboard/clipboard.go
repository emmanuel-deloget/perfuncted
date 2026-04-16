// Package clipboard provides cross-platform clipboard access for Linux desktops.
// On Wayland it uses wl-copy/wl-paste; on X11 it uses xclip. Both approaches
// spawn a subprocess, making this work regardless of compositor or toolkit.
package clipboard

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const commandTimeout = 5 * time.Second

// Clipboard reads and writes the system clipboard.
type Clipboard interface {
	// Get returns the current clipboard text content.
	Get() (string, error)
	// Set sets the clipboard text content.
	Set(text string) error
	// Close releases any resources held by the clipboard implementation.
	Close() error
}

// Open returns the best available Clipboard for the current session.
// On Wayland sessions (WAYLAND_DISPLAY set) it uses wl-copy/wl-paste.
// On X11 sessions it uses xclip.
func Open() (Clipboard, error) {
	env := append([]string(nil), os.Environ()...)
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		if _, err := exec.LookPath("wl-copy"); err == nil {
			if _, err := exec.LookPath("wl-paste"); err == nil {
				return &waylandClipboard{env: env}, nil
			}
		}
	}
	if os.Getenv("DISPLAY") != "" {
		if _, err := exec.LookPath("xclip"); err == nil {
			return &x11Clipboard{env: env}, nil
		}
	}
	return nil, fmt.Errorf("clipboard: no clipboard tool available (install wl-clipboard or xclip)")
}

// waylandClipboard uses wl-copy/wl-paste for Wayland sessions.
type waylandClipboard struct {
	env []string
}

func (c *waylandClipboard) Get() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "wl-paste", "--no-newline")
	cmd.Env = c.env
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("clipboard: wl-paste: %w", err)
	}
	return out.String(), nil
}

func (c *waylandClipboard) Set(text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "wl-copy")
	cmd.Env = c.env
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard: wl-copy: %w", err)
	}
	return nil
}

func (c *waylandClipboard) Close() error { return nil }

// x11Clipboard uses xclip for X11 sessions.
type x11Clipboard struct {
	env []string
}

func (c *x11Clipboard) Get() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "xclip", "-selection", "clipboard", "-o")
	cmd.Env = c.env
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("clipboard: xclip: %w", err)
	}
	return out.String(), nil
}

func (c *x11Clipboard) Set(text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "xclip", "-selection", "clipboard")
	cmd.Env = c.env
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard: xclip: %w", err)
	}
	return nil
}

func (c *x11Clipboard) Close() error { return nil }
