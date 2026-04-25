# Using Candy in VS Code

This repository includes ready-to-use VS Code workspace configuration so you can program and run Candy directly from the editor.

## Included workspace setup

- `.vscode/tasks.json`
  - `candy: run current file`
  - `candy: build native current file`
  - `candy: doctor`
  - `candy: build cli (raylib)`
  - `candywrap: wrap header`
- `.vscode/launch.json`
  - one-click launch entries for run/build/doctor
- `.vscode/settings.json`
  - file associations for `.candy`, `.cdy`, and `.candylib`
- `.vscode/extensions.json`
  - recommended extensions for Go/tooling workflows

## Quick workflow

1. Open the repository root in VS Code.
2. Open any `.candy` file from `examples/` or your own project.
3. Use one of these:
   - **Run and Debug** -> `Candy: Run current file`
   - **Terminal** -> **Run Task** -> `candy: run current file`

## Build native binary

- Run task: `candy: build native current file`
- Or launch config: `Candy: Build native current file`

## Toolchain health

- Run task or launch config: `candy: doctor`

This checks toolchain resolution and prints diagnostics.

## Wrapping C/C++ libraries

Run task: `candywrap: wrap header`

It prompts for:

- header path
- library name
- output directory

Then generates manifest/glue/docs/stubs for import in Candy.

