package candy_bindgen

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Library       string   `yaml:"library"`
	Namespace     string   `yaml:"namespace"`
	Headers       []string `yaml:"headers"`
	Includes      []string `yaml:"includes"`
	Defines       []string `yaml:"defines"`
	StripPrefixes []string `yaml:"strip_prefixes"`
	Ignore        []string `yaml:"ignore"`
	Profile       string   `yaml:"profile"`
	Language      string   `yaml:"language"`
	CXXStd        string   `yaml:"cxx_std"`
	LinkLibs      []string `yaml:"link_libs"`
	LinkLibDirs   []string `yaml:"link_lib_dirs"`
	LinkLDFlags   []string `yaml:"link_ldflags"`
	StaticLibs    []string `yaml:"static_libs"`
	StaticLink    *bool    `yaml:"static_link"`
	UnsafeABI     *bool    `yaml:"unsafe_abi"`
	Smart         *bool    `yaml:"smart"`
	SimpleOnly    *bool    `yaml:"simple"`
	CXXShim       *bool    `yaml:"cxx_shim"`
}

func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
