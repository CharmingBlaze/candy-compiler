# Candy (kid) — extended edition

This is an **optional second layer** on top of the core kid dialect ([CANDY_KID.md](CANDY_KID.md): variables, `loop`, drawing, `random`, etc.). It is aimed at small games: timers, collisions, a friendlier `draw` for textures, and common string/math helpers.

Everything below assumes you run the interpreter with **Raylib** (see [GETTING_STARTED.md](GETTING_STARTED.md)) unless noted.

---

## Philosophy (same as the base kid spec)

Keep ideas teachable. When the full engine already had a more detailed name, the extended spec adds **short aliases** so lessons read like the table below, without hiding the real implementation in [GAME_HELPERS.md](GAME_HELPERS.md) and `compiler/candy_raylib/`.

---

## Time

| Spec idea | In KGO |
|----------|--------|
| `seconds()` | `seconds` → Raylib time in seconds (alias of `getTime`). |
| `deltaTime()` | `deltaTime` (alias of `getFrameTime`). |
| `wait(2.5)` | In the Raylib build, `wait` is **wall-clock seconds** (`waitTime` / `rl.WaitTime`), matching “pause 2.5s”. |
| `every 2 { ... }` | The implementation is **`every(interval, fn)`** and **`after(delay, fn)`**, both taking a **callable** (a `fun` value, etc.); there is no separate `every` block form in the parser. Name your callback with `fun` and pass the identifier, e.g. the runnable [kid_every_after.candy](../examples/candy/kid_every_after.candy). As of the main loop, **`show()`** / `endDrawing` and **`flip()`** run **`helperTick`** so intervals advance once per **completed frame** in normal `clear` → draw → `show` loops. The older **`updateScene()`** path no longer runs a second tick to avoid double-firing. |
| Countdown in a `loop` | Use `seconds()` / `deltaTime()` and normal variables (example in the spec works). |

---

## Collision (2D)

| Spec | KGO / Raylib |
|------|----------------|
| `touching(x1, y1, r1, x2, y2, r2)` | `touching` → `circleCollision` (same 6 float arguments). |
| `boxHit(x1,y1,w1,h1, x2,y2,w2,h2)` | `boxHit` → `boxCollision` (8 args). |
| `inside(mx, my, x, y, w, h)` | `inside` → `pointInBox` (px, py, rect). |

---

## “Simple objects”

The language uses **`object Name { }`** (same as other declaration forms), **`fun`**, and classes/structs in the full compiler. The extended spec’s `object Player { x = 400 ... }` is a **teaching** shape: use the real Candy syntax for fields and methods from [LANGUAGE.md](LANGUAGE.md) if a construct fails to parse.

Instancing like `player = Player()` is only valid where the class/object system in your build supports that constructor form.

---

## Lists

| Spec | KGO |
|------|-----|
| `enemies.remove(0)` | Index 0: removes first element (same as kid list rules). |
| `enemies.removeLast()` / `pop` | Implemented as **`removeLast`** and **`pop`**. |
| `enemies.clear()` | `clear` on the array. |
| `enemies.empty()` | Use **`is_empty`** / **`empty` / `isempty`**. |
| `for` + `remove` while iterating | Be careful: removing by index in a `for` loop needs the usual off-by-one care; the spec’s pattern is intentional teaching. |

---

## Sprites & `draw`

| Spec | KGO |
|------|-----|
| `sprite("a.png")` | `sprite` → `loadTexture` (returns a numeric **texture id**). |
| `draw(id, x, y)` | **`draw`**: 3-tuple draws unscaled. |
| `draw(id, x, y, w, h)` + optional color | **5+ args** draw scaled with `DrawTexturePro`. 4 args: `(id, x, y, color)` for tint. |
| `drawRotated(id, x, y, angle, [scale], [color])` | **`drawRotated`**, maps to `DrawTextureEx` semantics. |
| `drawFlipped(id, x, y, flipH, flipV, [color])` | **`drawFlipped`** (uses `DrawTexturePro` and flipped source rect). |
| `unload(t)` | `unload` → `unloadTexture`. |

`player.width` / `player.height` for a **bare id** is not a special sprite object; use `texture` APIs or the Raylib `measure` / texture queries from the [EXTENSIONS](EXTENSIONS.md) / register file if you need dimensions.

`image("path", x, y)` from the kid spec remains: load-by-path cache; **`sprite`/`draw`** are for id-based drawing.

---

## Animation, particles, camera, save/load

Large parts of this space already exist under different names (see [GAME_HELPERS.md](GAME_HELPERS.md)):

- **Animation:** `animation`, `playAnimation`, `updateAnimation`, `animationFrame`, …
- **Particles / FX:** `particles`, `emit`, `drawParticles`, `explosion`, etc.; **camera** follows: `cameraFollow`, `cameraSnapTo`, `cameraBounds`, `cameraZoom`, `cameraShake` (helper model).
- **Save/load:** in-memory and file-style helpers: **`save`**, **`load`**, **`saveExists`**.

The **extended spec** only adds: **`camera(x, y)`** as an alias to **`cameraSnapTo`**, **`shake(duration, intensity)`** (intensity → helper shake field; **duration** is **reserved**), and **`zoom(scale)`** → **`cameraZoom`**.

`flash` / `fadeIn` / `fadeOut` as one-line effects are **not** implemented as dedicated builtins; use `clear`+alpha patterns or the tween/blend tools from the main API when you outgrow the kid layer.

---

## Grid helpers

- **`gridToPixel(col, tileSize)`** → `col * tileSize`
- **`pixelToGrid(px, tileSize)`** → `floor(px / tileSize)`  

These are small math helpers; your tile `map` is still a normal 2D array in Candy.

---

## Scenes and UI

The spec’s `scene = "menu"` and string comparisons work as normal Candy. There is no separate “scene” engine in the extended aliases—**use your own state variable** as in the doc.

---

## Better strings (top-level in the stdlib merge)

| Spec | KGO |
|------|-----|
| `join("Hello", " ", "World")` (many parts) | **`join`** (concatenates all arguments’ string forms). For **`join(list, ",")`** use `string.join(list, sep)` (module) — different API. |
| `split("a,b", ",")` | **`split`** (alias of `string.split` behavior). |
| `replace` / `contains` | **Top-level** `replace` and `contains` (string replace / contains). |
| `toNumber("123")` | **`toNumber`**, **`to_number`**, **`tonumber`**. |

---

## Better math

Top-level **`min`**, **`max`**, **`clamp`**, **`lerp`**, **`round`**, **`floor`**, **`ceil`**, **`sqrt`**, etc., come from the standard merge in the evaluator. The Raylib helper build also adds **`distance`** and **`angleBetween`** (alias of the existing **`distance2D`** and **`angleTo`**, degrees).

**`distance`:** `distance(x1, y1, x2, y2)` in the extended doc matches **`distance2D`**.

**`lerp`:** `lerp(a, b, t)` — not “move current toward target” in one form; the spec’s `x = lerp(x, targetX, 0.1)` is written in user code as three arguments where `a` is the start.

---

## Debugging

| Spec | KGO |
|------|-----|
| `debug("...")` | Stays the **print-style** `debug` in the stdlib (console), unless you only use the screen helpers. For **on-screen** text, use **`debugText`** (game helpers). |
| `debugLine(x1,y1,x2,y2, [c])` | **`debugLine`** = thin wrapper around `DrawLine` (default color lime). |
| `debugBox` / `debug` overlay | See **`debugBox`**, **`debugText`** in [GAME_HELPERS.md](GAME_HELPERS.md). |
| `pause()` | **Not** a global builtin; use a `paused` bool and `continue` the loop. |

---

## 3D (teaching)

The [Blitz-style 3D](examples/candy/candy_rain.candy) layer uses names like `Graphics3D`, `CreateCamera`, `RenderWorld`, etc. The extended doc’s `window3D` / `start3D` / `cube` names are **aspirational**; use the 3D section in [GAME_HELPERS](GAME_HELPERS.md) and the Raylib 3D registrations for what is actually present.

---

## “Full feature list” checklist

- **Core kid language** — [CANDY_KID.md](CANDY_KID.md) + this doc’s aliases.
- **Engine breadth** (particles, full camera, pathfinding, networking, etc.) — [GAME_HELPERS](GAME_HELPERS.md), not repeated here.
- **C# / C / custom transpiler** — not part of KGO’s primary path; see the kid spec’s C *design* section if you need a teaching comparison.

The platformer at the end of the original “extended” marketing text is **pseudocode**: you must line up `sprite`/`play` (path vs. sound id) and your own `object` semantics with the real [LANGUAGE](LANGUAGE.md) and Raylib register.

---

## See also

- [CANDY_KID.md](CANDY_KID.md) — base teaching dialect.  
- [GETTING_STARTED.md](GETTING_STARTED.md) — build and run.  
- [README.md](../README.md) — repo overview.

No separate `candyc` or “extended” compiler: one **interpreter**, optional **raylib** tag.
