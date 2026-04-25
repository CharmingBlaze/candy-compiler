package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"candy/candy_bindgen"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr *os.File) int {
	if len(args) == 0 || args[0] != "wrap" {
		fmt.Fprintln(stderr, "usage: candywrap wrap [header-files...] --name <lib> --output <dir>")
		return 2
	}
	fs := flag.NewFlagSet("wrap", flag.ContinueOnError)
	fs.SetOutput(stderr)
	libName := fs.String("name", "", "library name")
	outputDir := fs.String("output", ".", "output directory")
	includes := fs.String("include", "", "include directories (comma separated)")
	defines := fs.String("define", "", "preprocessor defines (comma separated)")
	configPath := fs.String("config", "", "path to candywrap yaml config")
	profileName := fs.String("profile", "", "built-in profile (raylib|sqlite|curl)")
	whole := fs.Bool("all", false, "scan provided roots and wrap the whole library automatically")
	rootCSV := fs.String("root", "", "library root directories/files (comma separated)")
	parserEngine := fs.String("parser", "auto", "parser engine (auto|libclang|regex)")
	language := fs.String("lang", "", "header language (c|c++)")
	cxxStd := fs.String("cxx-std", "", "C++ standard when --lang c++ (e.g. c++17)")
	namespace := fs.String("namespace", "", "namespace prefix for generated candy extern names")
	stripPrefix := fs.String("strip-prefix", "", "strip function name prefixes (comma separated)")
	ignore := fs.String("ignore", "", "ignore function names by regex/wildcard pattern (comma separated)")
	unsafeABI := fs.Bool("unsafe-abi", false, "allow unsafe ABI signatures in generated manifest (advanced)")
	staticLink := fs.Bool("static", false, "request fully static native link where possible")
	staticLibs := fs.String("static-lib", "", "libraries to link statically (comma separated)")
	linkLibs := fs.String("link-lib", "", "extra -l libraries to add (comma separated)")
	linkLibDirs := fs.String("link-lib-dir", "", "extra library directories to add (comma separated)")
	linkLDFlags := fs.String("link-ldflag", "", "extra linker flags to add (comma separated)")
	cxxShim := fs.Bool("cxx-shim", false, "emit <lib>_cxx_shim.cpp template for C++ wrapping")
	docsOut := fs.Bool("docs", true, "generate library markdown docs")
	simple := fs.Bool("simple", false, "skip risky wrappers and only emit safe externs")
	smart := fs.Bool("smart", true, "generate smart convenience wrappers where available")
	writeStub := fs.Bool("stub", true, "emit optional .candy stub file")
	if err := fs.Parse(normalizeWrapArgs(args[1:])); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	headers := fs.Args()
	explicitHeaderCount := len(headers)
	cfg, err := loadWrapConfig(*configPath)
	if err != nil {
		fmt.Fprintln(stderr, "config:", err)
		return 1
	}
	profile, err := candy_bindgen.ApplyProfile(strings.TrimSpace(*profileName))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if explicitHeaderCount == 0 && len(cfg.Headers) == 0 {
		headers = append(headers, profile.Headers...)
	}
	if strings.TrimSpace(*language) == "" && strings.TrimSpace(cfg.Language) != "" {
		*language = cfg.Language
	}
	if strings.TrimSpace(*language) == "" {
		*language = "c"
	}
	headers = append(headers, cfg.Headers...)
	headers = candy_bindgen.ExpandInputFiles(headers)
	headers = uniqueStrings(headers)
	var discoveredSources []string
	var discoveredIncludes []string
	if *whole {
		roots := mergeCSV(*rootCSV, headers)
		h, s, inc, derr := candy_bindgen.DiscoverLibraryFiles(roots, *language)
		if derr != nil {
			fmt.Fprintln(stderr, "discover:", derr)
			return 1
		}
		headers = h
		discoveredSources = s
		discoveredIncludes = inc
	}
	if len(headers) == 0 {
		fmt.Fprintln(stderr, "wrap: at least one header file required (arg, profile, or config)")
		return 2
	}
	if *libName == "" {
		if cfg.Library != "" {
			*libName = cfg.Library
		} else {
			*libName = candy_bindgen.DefaultLibName(headers[0])
		}
	}
	if cfg.SimpleOnly != nil {
		*simple = *cfg.SimpleOnly
	}
	if cfg.Smart != nil {
		*smart = *cfg.Smart
	}
	if strings.TrimSpace(*namespace) == "" && strings.TrimSpace(cfg.Namespace) != "" {
		*namespace = cfg.Namespace
	}
	if cfg.CXXShim != nil {
		*cxxShim = *cfg.CXXShim
	}
	if cfg.UnsafeABI != nil {
		*unsafeABI = *cfg.UnsafeABI
	}
	if cfg.StaticLink != nil {
		*staticLink = *cfg.StaticLink
	}
	if strings.TrimSpace(*cxxStd) == "" {
		*cxxStd = strings.TrimSpace(cfg.CXXStd)
	}

	includeList := mergeCSV(*includes, profile.Includes, cfg.Includes, discoveredIncludes)
	// Make wrapping work better out-of-the-box by adding header parent dirs as includes.
	headerParentDirs := make([]string, 0, len(headers))
	for _, h := range headers {
		if candy_bindgen.IsHeaderFile(h) {
			headerParentDirs = append(headerParentDirs, filepath.Dir(h))
		}
	}
	includeList = mergeCSV("", includeList, headerParentDirs)
	defineList := mergeCSV(*defines, profile.Defines, cfg.Defines)
	stripList := mergeCSV(*stripPrefix, cfg.StripPrefixes)
	ignoreList := mergeCSV(*ignore, cfg.Ignore)
	linkLibList := mergeCSV(*linkLibs, cfg.LinkLibs)
	linkLibDirList := mergeCSV(*linkLibDirs, cfg.LinkLibDirs)
	linkLDFlagList := mergeCSV(*linkLDFlags, cfg.LinkLDFlags)
	staticLibList := mergeCSV(*staticLibs, cfg.StaticLibs)
	lang := strings.ToLower(strings.TrimSpace(*language))
	if (lang == "c++" || lang == "cpp" || lang == "cxx") && strings.TrimSpace(*cxxStd) == "" {
		*cxxStd = "c++17"
	}

	api, warnings, err := candy_bindgen.ParseHeadersWithEngine(headers, candy_bindgen.ParseOptions{
		IncludeDirs: includeList,
		Defines:     defineList,
		SimpleOnly:  *simple,
		Language:    lang,
		UnsafeABI:   *unsafeABI,
	}, candy_bindgen.ParserEngine(strings.ToLower(strings.TrimSpace(*parserEngine))))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	transformWarnings, err := candy_bindgen.TransformAPI(api, *namespace, stripList, ignoreList)
	if err != nil {
		fmt.Fprintln(stderr, "transform:", err)
		return 1
	}
	warnings = append(warnings, transformWarnings...)
	for _, w := range warnings {
		fmt.Fprintln(stderr, "warn:", w)
	}

	if err := os.MkdirAll(*outputDir, 0o755); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	manifest := candy_bindgen.BuildManifestFromAPI(*libName, *namespace, headers, api)
	manifest.UnsafeABI = *unsafeABI
	manifest.Compile.IncludeDirs = append(manifest.Compile.IncludeDirs, includeList...)
	for _, d := range defineList {
		manifest.Compile.CFlags = append(manifest.Compile.CFlags, "-D"+d)
	}
	if strings.TrimSpace(*cxxStd) != "" {
		manifest.Compile.CFlags = append(manifest.Compile.CFlags, "-std="+strings.TrimSpace(*cxxStd))
	}
	manifest.Link.Static = *staticLink
	manifest.Link.StaticLibs = append(manifest.Link.StaticLibs, staticLibList...)
	manifest.Link.Libs = append(manifest.Link.Libs, linkLibList...)
	manifest.Link.LibDirs = append(manifest.Link.LibDirs, linkLibDirList...)
	manifest.Link.LDFlags = append(manifest.Link.LDFlags, linkLDFlagList...)
	if lang == "c++" || lang == "cpp" || lang == "cxx" {
		manifest.Link.Libs = append(manifest.Link.Libs, "stdc++")
	}
	for _, f := range headers {
		if candy_bindgen.IsSourceFile(f) {
			manifest.Compile.GlueSources = append(manifest.Compile.GlueSources, f)
		}
	}
	manifest.Compile.GlueSources = append(manifest.Compile.GlueSources, discoveredSources...)
	manifestPath := filepath.Join(*outputDir, *libName+".candylib")
	gluePath := filepath.Join(*outputDir, *libName+"_glue.c")
	if err := candy_bindgen.WriteManifest(manifestPath, manifest); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := candy_bindgen.WriteGlue(gluePath, *libName, manifest.Externs); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if *cxxShim {
		shimPath := filepath.Join(*outputDir, *libName+"_cxx_shim.cpp")
		if err := candy_bindgen.WriteCXXShimTemplate(shimPath, *libName, &manifest); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		manifest.Compile.GlueSources = append(manifest.Compile.GlueSources, filepath.Base(shimPath))
		fmt.Fprintln(stdout, "Generated", shimPath)
		if err := candy_bindgen.WriteManifest(manifestPath, manifest); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	if *smart {
		addSmartWrappers(manifest, *libName)
	}
	if *writeStub {
		stubPath := filepath.Join(*outputDir, *libName+".candy")
		if err := candy_bindgen.WriteCandyStub(stubPath, *libName, manifest.Externs); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "Generated", stubPath)
		if strings.TrimSpace(*namespace) != "" {
			nsPath := filepath.Join(*outputDir, *libName+"_namespace.candy")
			if err := candy_bindgen.WriteCandyNamespaceStub(nsPath, *libName, *namespace, manifest.Externs); err != nil {
				fmt.Fprintln(stderr, err)
				return 1
			}
			fmt.Fprintln(stdout, "Generated", nsPath)
		}
	}
	if *docsOut {
		docsPath := filepath.Join(*outputDir, *libName+".md")
		if err := candy_bindgen.WriteLibraryDocs(docsPath, *libName, api, &manifest); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "Generated", docsPath)
	}
	fmt.Fprintln(stdout, "Generated", manifestPath)
	fmt.Fprintln(stdout, "Generated", gluePath)
	return 0
}

func mergeCSV(csv string, groups ...[]string) []string {
	var out []string
	if strings.TrimSpace(csv) != "" {
		for _, it := range strings.Split(csv, ",") {
			it = strings.TrimSpace(it)
			if it != "" {
				out = append(out, it)
			}
		}
	}
	for _, g := range groups {
		for _, it := range g {
			it = strings.TrimSpace(it)
			if it != "" {
				out = append(out, it)
			}
		}
	}
	return uniqueStrings(out)
}

func uniqueStrings(items []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, it := range items {
		key := strings.ToLower(strings.TrimSpace(it))
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, strings.TrimSpace(it))
	}
	return out
}

func loadWrapConfig(path string) (*candy_bindgen.Config, error) {
	if strings.TrimSpace(path) == "" {
		return &candy_bindgen.Config{}, nil
	}
	cfg, err := candy_bindgen.LoadConfig(path)
	if err != nil {
		return nil, err
	}
	base := filepath.Dir(path)
	resolve := func(items []string) []string {
		out := make([]string, 0, len(items))
		for _, it := range items {
			it = strings.TrimSpace(it)
			if it == "" || filepath.IsAbs(it) {
				out = append(out, it)
				continue
			}
			out = append(out, filepath.Clean(filepath.Join(base, it)))
		}
		return out
	}
	cfg.Headers = resolve(cfg.Headers)
	cfg.Includes = resolve(cfg.Includes)
	return cfg, nil
}

func addSmartWrappers(_ candy_bindgen.Manifest, _ string) {
	// Placeholder for future library-specific smart wrapper generation.
}

func normalizeWrapArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}
	boolFlags := map[string]struct{}{
		"--all": {}, "--unsafe-abi": {}, "--static": {}, "--cxx-shim": {},
		"--docs": {}, "--simple": {}, "--smart": {}, "--stub": {},
	}
	var flags []string
	var positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !strings.HasPrefix(a, "-") {
			positional = append(positional, a)
			continue
		}
		flags = append(flags, a)
		if strings.Contains(a, "=") {
			continue
		}
		if _, isBool := boolFlags[strings.ToLower(strings.TrimSpace(a))]; isBool {
			continue
		}
		if i+1 < len(args) {
			flags = append(flags, args[i+1])
			i++
		}
	}
	return append(flags, positional...)
}
