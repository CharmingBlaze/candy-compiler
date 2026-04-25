package candy_bindgen

import "fmt"

type ParserEngine string

const (
	ParserAuto     ParserEngine = "auto"
	ParserRegex    ParserEngine = "regex"
	ParserLibclang ParserEngine = "libclang"
)

func ParseHeadersWithEngine(headerFiles []string, opts ParseOptions, engine ParserEngine) (*API, []string, error) {
	isCXX := opts.Language == "c++" || opts.Language == "cpp" || opts.Language == "cxx"
	switch engine {
	case ParserRegex:
		api, warns, err := ParseHeaders(headerFiles, opts)
		if isCXX {
			warns = append(warns, "regex parser has limited C++ support; prefer --parser libclang for C++ headers")
		}
		return api, warns, err
	case ParserLibclang:
		return parseHeadersLibclang(headerFiles, opts)
	case ParserAuto:
		api, warns, err := parseHeadersLibclang(headerFiles, opts)
		if err == nil {
			return api, warns, nil
		}
		api, warns2, err2 := ParseHeaders(headerFiles, opts)
		if err2 != nil {
			return nil, nil, fmt.Errorf("auto parser failed: libclang=%v; regex=%w", err, err2)
		}
		warns = append(warns, "libclang unavailable, fell back to regex parser")
		warns = append(warns, warns2...)
		return api, warns, nil
	default:
		return nil, nil, fmt.Errorf("unknown parser engine %q", engine)
	}
}
