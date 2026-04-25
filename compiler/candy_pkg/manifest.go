package candy_pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Manifest struct {
	Name    string
	Version string
	Deps    map[string]string
}

func LoadManifest(path string) (Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return Manifest{}, err
	}
	defer f.Close()
	m := Manifest{Deps: map[string]string{}}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "name":
			m.Name = v
		case "version":
			m.Version = v
		default:
			if strings.HasPrefix(k, "dep.") {
				m.Deps[strings.TrimPrefix(k, "dep.")] = v
			}
		}
	}
	if err := sc.Err(); err != nil {
		return Manifest{}, err
	}
	if m.Name == "" {
		return Manifest{}, fmt.Errorf("manifest missing name")
	}
	return m, nil
}
