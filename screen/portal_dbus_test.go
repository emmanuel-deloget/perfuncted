package screen

import "testing"

func TestFileURIPath(t *testing.T) {
	path, err := fileURIPath("file:///tmp/portal%20shot.png")
	if err != nil {
		t.Fatal(err)
	}
	if path != "/tmp/portal shot.png" {
		t.Fatalf("path = %q", path)
	}
}

func TestFileURIPathRejectsUnsupportedHost(t *testing.T) {
	if _, err := fileURIPath("file://remotehost/tmp/portal.png"); err == nil {
		t.Fatal("expected error for unsupported host")
	}
}
