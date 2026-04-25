//go:build !libclang

package candy_bindgen

import "fmt"

func parseHeadersLibclangImpl(_ []string, _ ParseOptions) (*API, []string, error) {
	return nil, nil, fmt.Errorf("libclang parser unavailable in this build (rebuild with -tags libclang)")
}
