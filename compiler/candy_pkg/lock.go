package candy_pkg

import (
	"encoding/json"
	"os"
)

type Lock struct {
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Deps    map[string]string `json:"deps"`
}

func WriteLock(path string, m Manifest) error {
	l := Lock{Name: m.Name, Version: m.Version, Deps: m.Deps}
	b, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
