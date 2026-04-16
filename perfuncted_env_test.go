package perfuncted

import (
	"os"
	"testing"
)

func TestResolveSessionEnvPrefersExplicitValues(t *testing.T) {
	t.Setenv("XDG_RUNTIME_DIR", "/run/user/1000")
	t.Setenv("WAYLAND_DISPLAY", "wayland-0")
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", "unix:path=/run/user/1000/bus")

	env, err := resolveSessionEnv(Options{
		XDGRuntimeDir:      "/tmp/perfuncted-xdg-test",
		WaylandDisplay:     "wayland-9",
		DBusSessionAddress: "unix:path=/tmp/perfuncted-xdg-test/bus",
	})
	if err != nil {
		t.Fatal(err)
	}
	if env.xdgRuntimeDir != "/tmp/perfuncted-xdg-test" {
		t.Fatalf("xdg = %q", env.xdgRuntimeDir)
	}
	if env.waylandDisplay != "wayland-9" {
		t.Fatalf("wayland = %q", env.waylandDisplay)
	}
	if env.dbusSessionAddress != "unix:path=/tmp/perfuncted-xdg-test/bus" {
		t.Fatalf("dbus = %q", env.dbusSessionAddress)
	}
}

func TestApplySessionEnvRestoresPreviousProcessEnv(t *testing.T) {
	t.Setenv("XDG_RUNTIME_DIR", "/run/user/1000")
	t.Setenv("WAYLAND_DISPLAY", "wayland-0")
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", "unix:path=/run/user/1000/bus")

	restore := applySessionEnv(sessionEnv{
		xdgRuntimeDir:      "/tmp/perfuncted-xdg-test",
		waylandDisplay:     "wayland-9",
		dbusSessionAddress: "unix:path=/tmp/perfuncted-xdg-test/bus",
	})

	if got := os.Getenv("XDG_RUNTIME_DIR"); got != "/tmp/perfuncted-xdg-test" {
		t.Fatalf("XDG_RUNTIME_DIR = %q", got)
	}
	if got := os.Getenv("WAYLAND_DISPLAY"); got != "wayland-9" {
		t.Fatalf("WAYLAND_DISPLAY = %q", got)
	}
	if got := os.Getenv("DBUS_SESSION_BUS_ADDRESS"); got != "unix:path=/tmp/perfuncted-xdg-test/bus" {
		t.Fatalf("DBUS_SESSION_BUS_ADDRESS = %q", got)
	}

	restore()

	if got := os.Getenv("XDG_RUNTIME_DIR"); got != "/run/user/1000" {
		t.Fatalf("restored XDG_RUNTIME_DIR = %q", got)
	}
	if got := os.Getenv("WAYLAND_DISPLAY"); got != "wayland-0" {
		t.Fatalf("restored WAYLAND_DISPLAY = %q", got)
	}
	if got := os.Getenv("DBUS_SESSION_BUS_ADDRESS"); got != "unix:path=/run/user/1000/bus" {
		t.Fatalf("restored DBUS_SESSION_BUS_ADDRESS = %q", got)
	}
}
