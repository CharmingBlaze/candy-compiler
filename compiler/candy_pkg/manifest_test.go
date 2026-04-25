package candy_pkg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManifestLoadAndLock(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "candy.pkg")
	src := "name = app\nversion = 0.1.0\ndep.math = 1.0.0\n"
	if err := os.WriteFile(manifestPath, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "app" || m.Version != "0.1.0" {
		t.Fatalf("bad manifest: %#v", m)
	}
	lockPath := filepath.Join(dir, "candy.lock")
	if err := WriteLock(lockPath, m); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(lockPath); err != nil {
		t.Fatal(err)
	}
}
