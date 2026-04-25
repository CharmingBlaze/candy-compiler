# Raylib extensions (reference)

This reference tracks the host wrapper implementation in:
- `compiler/candy_raylib/register.go`
- `compiler/candy_raylib/rcore.go`
- `compiler/candy_raylib/rshapes.go`
- `compiler/candy_raylib/rtextures.go`
- `compiler/candy_raylib/rtext.go`
- `compiler/candy_raylib/rshaders.go`
- `compiler/candy_raylib/rmodels.go`
- `compiler/candy_raylib/raudio.go`
- `compiler/candy_raylib/reasing.go`
- `compiler/candy_raylib/utils.go`

Raylib cheatsheet (v6): https://www.raylib.com/cheatsheet/cheatsheet.html

## Registration model

All wrappers are flat builtins registered by `RegisterBuiltins()` in `register.go`.
Compatibility aliases remain (`window`, `clear`, `circle`, `text`, `flip`, `key`, `shouldClose`).

## Handle-backed resources

These builtins return integer handles stored in Go maps and validated on use:
- Textures (`textureId`)
- Render textures (`rtId`)
- Images (`imageId`)
- Fonts (`fontId`)
- Shaders (`shaderId`)
- Models (`modelId`)
- Model animations (`animId`)
- Sounds (`soundId`)
- Music streams (`musicId`)

Invalid handle usage returns runtime errors.

## Value conversion policy

Wrapper boundary is `*candy_evaluator.Value` only:
- numbers accept int/float
- colors use string names
- struct-like values return as maps (for example vectors, mouse position, ray collision data)
- no raw pointers are exposed to Candy scripts

## Coverage

Current wrapper covers broad portions of `rcore`, `rshapes`, `rtextures`, `rtext`, `rshaders`, `rmodels`, `raudio`, and `reasing`.
See `docs/EXTENSIONS.md` for usage-level function lists and examples.
For high-level gameplay helper commands, see `docs/GAME_HELPERS.md`.

Recent additions include:
- gamepad helpers: `getGamepadName`, `isGamepadButtonDown`
- image helpers: `genImageColor`, `exportImage`
