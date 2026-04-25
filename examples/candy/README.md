# Example Candy programs

All scripts here run with the `candy` executable. **Graphics, mouse, and sound** need a build of `candy` with **`-tags raylib`** (see [../docs/GETTING_STARTED.md](../docs/GETTING_STARTED.md)).

## Kid / teaching spec ([CANDY_KID.md](../docs/CANDY_KID.md))

| File | What it does |
|------|----------------|
| [kid_moving_ball.candy](kid_moving_ball.candy) | Arrow keys move a circle. |
| [kid_clicker.candy](kid_clicker.candy) | Clicker / score. |
| [kid_catch.candy](kid_catch.candy) | Catch a falling object. |
| [kid_bounce.candy](kid_bounce.candy) | Bouncing ball off screen edges. |

**Run (from the `compiler` directory):**  
`go run -tags raylib ./cmd/candy ../examples/candy/kid_moving_ball.candy`

## Other examples

| File | Notes |
|------|--------|
| [bounce.candy](bounce.candy) | Older style: `while not shouldClose`, `flip`, `text(x,y,msg,…)`; `window` is `width, height, "title"`. |
| [candy_rain.candy](candy_rain.candy) | 3D-style helpers (`Graphics3D`, …) — see file header. |
| [candy_crusher.candy](candy_crusher.candy) | Game sample. |
| [candy_shop.candy](candy_shop.candy) | Game sample. |
| [spec_test.candy](spec_test.candy) | Internal or language checks. |
