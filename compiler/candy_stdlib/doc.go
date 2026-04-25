// Package candy_stdlib documents standard-library binding for the K-Go interpreter
// and any future native backends. Built-in functions (len, Ok, Err) live in
// candy_evaluator and are preloaded in the global environment. Additional
// “stdlib” can be provided as: (1) more builtins registered like Builtins,
// (2) K-Go source files loaded at startup, or (3) thin Go wrappers in a
// dedicated package once FFI or transpile exists.
package candy_stdlib
