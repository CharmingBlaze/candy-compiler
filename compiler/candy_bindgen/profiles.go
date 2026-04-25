package candy_bindgen

import "fmt"

type Profile struct {
	Headers  []string
	Includes []string
	Defines  []string
}

var builtinProfiles = map[string]Profile{
	"raylib": {
		Headers: []string{"raylib.h"},
	},
	"sqlite": {
		Headers: []string{"sqlite3.h"},
	},
	"curl": {
		Headers: []string{"curl/curl.h"},
	},
}

func ApplyProfile(name string) (Profile, error) {
	if name == "" {
		return Profile{}, nil
	}
	p, ok := builtinProfiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("unknown profile %q", name)
	}
	return p, nil
}
