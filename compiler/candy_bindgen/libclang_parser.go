package candy_bindgen

func parseHeadersLibclang(headerFiles []string, opts ParseOptions) (*API, []string, error) {
	return parseHeadersLibclangImpl(headerFiles, opts)
}
