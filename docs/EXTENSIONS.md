# Raylib in Candy

Candy has built-in support for [Raylib](https://www.raylib.com/), a simple library for making 2D and 3D games and graphics. You can draw shapes, load images and sounds, read keyboard/mouse input, and build full games — all from Candy.

If you want to import external C libraries via generated manifests, see [CANDYWRAP.md](CANDYWRAP.md).

## Your first program

```candy
window(800, 600, "My Game")   // open a 800×600 window
setTargetFPS(60)               // aim for 60 frames per second

while !shouldClose() {         // loop until the user closes the window
  clear("black")               // wipe the screen black each frame
  drawText(20, 20, "Hello, Raylib!", 24, "white")
  flip()                       // show everything you just drew
}

closeWindow()                  // clean up when done
```

Save that as `game.candy` and run it. You will see a window with white text on a black background.

---

## Game helper layer

Candy also has a higher-level helper surface for common gameplay patterns (camera follow, input helpers, animation controllers, tweens, timers, particles, collision helpers, pathfinding, spatial grids, scene/state, save/load, and debug overlays).

Use this reference for the helper API:

- [GAME_HELPERS.md](GAME_HELPERS.md)

The rest of this document focuses on the lower-level Raylib-style wrappers.

### BlitzBasic-style 3D aliases

When you build with `-tags raylib`, Candy also registers a small **entity + camera** layer inspired by Blitz3D (`Graphics3D`, `CreateCamera`, `CreateSphere`, `PositionEntity`, `RenderWorld`, `flip`, `keyHit`, …). That surface is documented alongside the core language in [LANGUAGE.md](LANGUAGE.md) (section *3D helpers*). Use it for short tutorials; use the full Raylib calls below for production games.

---

## Immediate-mode GUI (raygui)

Candy includes the full [raygui](https://github.com/raysan5/raygui) immediate-mode GUI library — buttons, sliders, checkboxes, dropdowns, text input, color pickers, dialogs, tab bars, list views, and more.

See the full reference:

- [RAYGUI.md](RAYGUI.md)

---

## Key concepts

### Handles
Resources like images, textures, sounds, fonts, and models are not passed around as objects. Instead, `load*` and `gen*` functions return an **integer handle** (e.g. `imageId`, `textureId`, `soundId`). Pass that number to every function that needs the resource, and call the matching `unload*` function when you are done with it.

```candy
let tex = loadTexture("player.png")   // tex is just a number, e.g. 1
drawTexture(tex, 100, 200)            // pass the number to draw
unloadTexture(tex)                    // free GPU memory when done
```

### Colors
Any function that takes a color accepts **either** a lowercase name string **or** a `{r, g, b, a}` map:

```candy
drawCircle(400, 300, 50, "red")               // name string
drawCircle(400, 300, 50, {r:255, g:0, b:0, a:255})  // map
drawCircle(400, 300, 50, color(255, 0, 0))    // constructor
drawCircle(400, 300, 50, COLOR_RED())         // named constant
```

### Struct maps
Vectors, rectangles, and similar types are just Candy maps. You can write them as map literals or use the constructor builtins:

```candy
let pos  = vec2(100, 200)          // {x:100, y:200}
let size = vec3(1, 1, 1)           // {x:1, y:1, z:1}
let box  = rect(10, 10, 80, 40)    // {x:10, y:10, width:80, height:40}
```

---

## Struct constructors

These builtins create the common value types used throughout the API. You can also just write the map literal directly — both styles work everywhere.

| Constructor | Returns | Example |
|---|---|---|
| `vec2(x, y)` | `{x, y}` | `vec2(100, 200)` |
| `vec3(x, y, z)` | `{x, y, z}` | `vec3(0, 1, 0)` |
| `vec4(x, y, z, w)` | `{x, y, z, w}` | `vec4(1, 0, 0, 1)` |
| `rect(x, y, width, height)` | `{x, y, width, height}` | `rect(10, 10, 80, 40)` |
| `color(r, g, b)` or `color(r, g, b, a)` | `{r, g, b, a}` | `color(255, 128, 0)` |
| `color("name")` | `{r, g, b, a}` | `color("skyblue")` |
| `boundingBox(min, max)` | `{min:{x,y,z}, max:{x,y,z}}` | `boundingBox(vec3(-1,-1,-1), vec3(1,1,1))` |
| `boundingBox(minX,minY,minZ, maxX,maxY,maxZ)` | same | flat 6-number form |
| `ray(px,py,pz, dx,dy,dz)` | `{position:{x,y,z}, direction:{x,y,z}}` | ray origin + direction |
| `camera3D(px,py,pz, tx,ty,tz, fovy)` | camera map | position, target, field-of-view |
| `camera2D(offsetX,offsetY, targetX,targetY, [rotation], [zoom])` | camera map | 2D camera with optional rotation and zoom |
| `getNamedColor("name")` | `{r, g, b, a}` | look up any color by name at runtime |

---

## Colors

Every color argument in the entire API accepts **any of these three forms**:

```candy
// 1. Name string (case-insensitive)
drawCircle(400, 300, 50, "red")

// 2. Map literal with r,g,b,a fields (0–255)
drawCircle(400, 300, 50, {r:255, g:0, b:0, a:255})

// 3. color() constructor
drawCircle(400, 300, 50, color(255, 0, 0))
```

**Named color constants** — call with `()` to get the `{r,g,b,a}` map:

`COLOR_LIGHTGRAY()`, `COLOR_GRAY()`, `COLOR_DARKGRAY()`, `COLOR_YELLOW()`, `COLOR_GOLD()`, `COLOR_ORANGE()`, `COLOR_PINK()`, `COLOR_RED()`, `COLOR_MAROON()`, `COLOR_GREEN()`, `COLOR_LIME()`, `COLOR_DARKGREEN()`, `COLOR_SKYBLUE()`, `COLOR_BLUE()`, `COLOR_DARKBLUE()`, `COLOR_PURPLE()`, `COLOR_VIOLET()`, `COLOR_DARKPURPLE()`, `COLOR_BEIGE()`, `COLOR_BROWN()`, `COLOR_DARKBROWN()`, `COLOR_WHITE()`, `COLOR_BLACK()`, `COLOR_BLANK()`, `COLOR_MAGENTA()`, `COLOR_RAYWHITE()`

**All accepted name strings:** `white`, `black`, `red`, `green`, `blue`, `yellow`, `gray`, `darkgray`, `lightgray`, `magenta`, `gold`, `orange`, `lime`, `darkgreen`, `darkblue`, `maroon`, `beige`, `navy`, `purple`, `pink`, `brown`, `sky` / `skyblue`, `violet`, `darkpurple`, `darkbrown`, `blank` / `transparent`, `raywhite`

---

## Window & frame

Every Raylib program follows this structure:

```candy
window(800, 600, "Title")    // create the window
setTargetFPS(60)

while !shouldClose() {
  clear("black")             // must be first inside the loop
  // ... draw things here ...
  flip()                     // must be last — shows the frame
}

closeWindow()
```

| Function | What it does |
|---|---|
| `window(w, h, title)` | Open a window (same as `initWindow`) |
| `initWindow(w, h, title)` | Open a window |
| `closeWindow()` | Close the window and free resources |
| `shouldClose()` | Returns true when the user clicks ✕ or presses Escape |
| `setTargetFPS(fps)` | Limit the frame rate |
| `clear(color)` / `clearBackground(color)` | Fill the screen with a color |
| `flip()` / `beginDrawing()` / `endDrawing()` | Show the current frame (`flip` does both begin+end) |
| `getFrameTime()` | Seconds since last frame — use this for movement speed |
| `getTime()` | Total seconds since the window opened |
| `getFPS()` | Current frames per second |
| `getScreenWidth()` / `getScreenHeight()` | Window size in pixels |
| `getRenderWidth()` / `getRenderHeight()` | Actual render size (may differ on HiDPI) |
| `takeScreenshot(filename)` | Save a PNG screenshot |

**Window control:**
`setWindowTitle(title)`, `setWindowSize(w, h)`, `setWindowPosition(x, y)`, `setWindowOpacity(0.0–1.0)`, `setWindowMinSize(w, h)`, `setWindowMaxSize(w, h)`, `setWindowMonitor(index)`, `setWindowIcon(imageId)`, `setWindowIcons([imageIds])`

**Window state queries:**
`isWindowReady()`, `isWindowFullscreen()`, `isWindowHidden()`, `isWindowMinimized()`, `isWindowMaximized()`, `isWindowFocused()`, `isWindowResized()`

**Window state control:**
`maximizeWindow()`, `minimizeWindow()`, `restoreWindow()`, `setWindowFocused()`, `toggleFullscreen()`, `toggleBorderlessWindowed()`

**Window flags** (for `isWindowState` / `setWindowState` / `clearWindowState`):
Pass a flag string or an array of flag strings:
`"vsync"`, `"fullscreen"`, `"resizable"`, `"undecorated"`, `"hidden"`, `"minimized"`, `"maximized"`, `"unfocused"`, `"topmost"`, `"always_run"`, `"transparent"`, `"highdpi"`, `"mouse_passthrough"`, `"borderless"`, `"msaa4x"`, `"interlaced"`

```candy
setWindowState("resizable")
setWindowState(["resizable", "highdpi"])
```

**Monitor info:**
`getMonitorCount()`, `getCurrentMonitor()`, `getMonitorName(n)`, `getMonitorWidth(n)`, `getMonitorHeight(n)`, `getMonitorPhysicalWidth(n)`, `getMonitorPhysicalHeight(n)`, `getMonitorRefreshRate(n)`, `getMonitorPosition(n)` → `{x,y}`, `getWindowPosition()` → `{x,y}`, `getWindowScaleDPI()` → `{x,y}`

**Misc:**
`enableEventWaiting()` / `disableEventWaiting()`, `setExitKey("keyname")`

---

## Cursor

```candy
hideCursor()        // hide the system mouse cursor
showCursor()        // show it again
isCursorHidden()    // → bool
enableCursor()      // allow cursor movement (default)
disableCursor()     // lock cursor to window — good for FPS games
isCursorOnScreen()  // → bool
```

---

## Input

### Keyboard

Key names are lowercase strings like `"space"`, `"enter"`, `"escape"`, `"left"`, `"right"`, `"up"`, `"down"`, `"shift"`, `"ctrl"`, `"a"` … `"z"`, `"f1"` … `"f12"`.

```candy
if isKeyPressed("space") {   // true on the frame the key is first pressed
  jump()
}
if isKeyDown("right") {      // true every frame the key is held
  player.x = player.x + 200 * getFrameTime()
}
isKeyReleased("space")       // true on the frame the key is released
isKeyUp("shift")             // true while key is NOT held
key("a")                     // returns the key code integer for "a"
```

### Mouse

```candy
let x = getMouseX()
let y = getMouseY()
let pos = getMousePosition()     // → {x, y}
let wheel = getMouseWheelMove()  // → float, positive = scroll up

isMouseButtonDown(0)      // 0=left, 1=right, 2=middle
isMouseButtonPressed(0)   // true only on the first frame clicked
isMouseButtonReleased(0)  // true on the frame released

setMousePosition(x, y)
setMouseOffset(x, y)
setMouseScale(x, y)
```

### Clipboard

```candy
setClipboardText("hello")
let t = getClipboardText()   // → string
let img = getClipboardImage() // → imageId (if clipboard contains an image)
```

### Gamepad

```candy
if isGamepadAvailable(0) {             // 0 = first connected gamepad
  let name = getGamepadName(0)
  if isGamepadButtonPressed(0, 0) { }  // button 0 = A/Cross
  let axis = getGamepadAxisValue(0, 0) // left stick X
}
```

---

## 2D drawing

All drawing must happen between `clear()` and `flip()`.

### Pixels and lines

```candy
drawPixel(x, y, color)
drawLine(x1, y1, x2, y2, color)
drawLineEx(start, end, thickness, color)   // start/end are {x,y} maps
drawLineBezier(start, end, thickness, color)
drawLineStrip([{x,y}, {x,y}, ...], color)
drawLineDashed(start, end, dashSize, spaceSize, color)
```

### Circles and ellipses

```candy
drawCircle(x, y, radius, color)
drawCircleLines(x, y, radius, color)         // outline only
drawCircleV({x,y}, radius, color)            // center as {x,y} map
drawCircleGradient({x,y}, radius, innerColor, outerColor)
drawCircleSector({x,y}, radius, startAngle, endAngle, segments, color)
drawCircleSectorLines(...)                   // same args, outline only
drawEllipse(x, y, radiusH, radiusV, color)
drawEllipseLines(x, y, radiusH, radiusV, color)
drawRing(x, y, innerRadius, outerRadius, startAngle, endAngle, [segments], color)
drawRingLines(...)
```

### Rectangles

```candy
drawRectangle(x, y, width, height, color)
drawRectangleLines(x, y, width, height, color)      // outline
drawRectangleLinesEx(rec, lineThickness, color)     // rec = {x,y,width,height}
drawRectangleRounded(rec, roundness, segments, color)
drawRectangleGradientV(x, y, w, h, topColor, bottomColor)
drawRectangleGradientH(x, y, w, h, leftColor, rightColor)
```

### Triangles and polygons

```candy
drawTriangle(x1,y1, x2,y2, x3,y3, color)
drawTriangleLines(x1,y1, x2,y2, x3,y3, color)
drawTriangleFan([{x,y}, ...], color)
drawTriangleStrip([{x,y}, ...], color)
drawPoly(center, sides, radius, rotation, color)
drawPolyLines(center, sides, radius, rotation, color)
drawPolyLinesEx(center, sides, radius, rotation, lineThick, color)
```

### Splines

Splines draw smooth curves through a list of `{x,y}` control points:

```candy
let pts = [{x:100,y:300}, {x:200,y:100}, {x:300,y:300}, {x:400,y:100}]
drawSplineCatmullRom(pts, 2, "white")       // smooth curve
drawSplineBezierCubic(pts, 2, "yellow")     // cubic bezier
drawSplineLinear(pts, 1, "gray")            // straight segments
```

Individual segment drawing: `drawSplineSegmentLinear`, `drawSplineSegmentBasis`, `drawSplineSegmentCatmullRom`, `drawSplineSegmentBezierQuadratic`, `drawSplineSegmentBezierCubic`

Evaluate a point on a spline (returns `{x,y}`): `getSplinePointLinear(p1,p2,t)`, `getSplinePointBasis(p1,p2,p3,p4,t)`, `getSplinePointCatmullRom(...)`, `getSplinePointBezierQuad(p1,c2,p3,t)`, `getSplinePointBezierCubic(p1,c2,c3,p4,t)`

### Shapes texture

Use a texture atlas as the source for all shape drawing (advanced):
`setShapesTexture(textureId, rec)`, `getShapesTexture()` → `textureId`, `getShapesTextureRectangle()` → `{x,y,width,height}`

### 2D collision detection

```candy
// Rectangle vs rectangle
if checkCollisionRecs(x1,y1,w1,h1, x2,y2,w2,h2) { }

// Circle vs circle
if checkCollisionCircles(cx1,cy1,r1, cx2,cy2,r2) { }

// Circle vs rectangle
if checkCollisionCircleRec(cx,cy,cr, rx,ry,rw,rh) { }

// Point checks
if checkCollisionPointRec({x,y}, {x,y,width,height}) { }
if checkCollisionPointCircle({x,y}, {x,y}, radius) { }
if checkCollisionPointTriangle(point, p1, p2, p3) { }

// Get the overlapping area of two rectangles
let overlap = getCollisionRec(rec1, rec2)   // → {x,y,width,height}
```

Also: `checkCollisionCircleLine`, `checkCollisionPointLine`, `checkCollisionPointPoly`, `checkCollisionLines` → `{hit, point}`

### FPS counter

```candy
drawFPS(10, 10)    // draw the current FPS in the top-left corner
```

---

## Text

### Loading fonts

By default `drawText` uses the built-in font. Load a custom font file for nicer text:

```candy
let font = loadFont("assets/myfont.ttf")
// or with a specific size and character set:
let font = loadFontEx("assets/myfont.ttf", 32, [])

drawTextEx(font, "Hello!", 100, 100, 32, 2, "white")
unloadFont(font)
```

| Function | What it does |
|---|---|
| `getFontDefault()` → `fontId` | Get the built-in font as a handle |
| `loadFont(path)` → `fontId` | Load a font from a .ttf/.otf file |
| `loadFontEx(path, fontSize, [codepoints])` → `fontId` | Load with specific size and optional character list |
| `loadFontFromImage(imageId, keyColor, firstChar)` → `fontId` | Load from a bitmap font image |
| `loadFontFromMemory(fileType, bytesArr, fontSize, [codepoints])` → `fontId` | Load from memory |
| `isFontValid(fontId)` → bool | Check if font loaded successfully |
| `unloadFont(fontId)` | Free font memory |

### Drawing text

```candy
// Simple text — uses the built-in font
drawText(x, y, "Hello!", fontSize, "white")

// Custom font with letter spacing
drawTextEx(fontId, "Hello!", x, y, fontSize, spacing, "white")

// Full control — rotation, origin point
drawTextPro(fontId, "Hello!", position, origin, rotation, fontSize, spacing, tint)
// position, origin are {x,y} maps; rotation is degrees

// Draw a single Unicode character
drawTextCodepoint(fontId, 65, {x:100,y:100}, 32, "white")
```

### Measuring text

```candy
let w = measureText("Hello!", 24)                    // → pixel width (int)
let size = measureTextEx(fontId, "Hello!", 24, 2)    // → {x, y} (width, height)
setTextLineSpacing(4)                                 // extra pixels between lines
```

### Glyph info

```candy
let idx  = getGlyphIndex(fontId, 65)                 // → int index in font atlas
let info = getGlyphInfo(fontId, 65)                  // → {value, offsetX, offsetY, advanceX}
let rec  = getGlyphAtlasRec(fontId, 65)              // → {x, y, width, height} in atlas
```

### Text string utilities

These work on plain Candy strings:

```candy
textIsEqual("abc", "abc")          // → true
textLength("hello")                // → 5
textToUpper("hello")               // → "HELLO"
textToLower("HELLO")               // → "hello"
textToPascal("hello world")        // → "HelloWorld"
textToCamel("hello world")         // → "helloWorld"
textToSnake("HelloWorld")          // → "hello_world"
textSubtext("hello world", 6, 5)   // → "world"
textReplace("hello world", "world", "there")  // → "hello there"
textInsert("helloworld", " ", 5)   // → "hello world"
textRemoveSpaces("h e l l o")      // → "hello"
getTextBetween("(value)", "(", ")")  // → "value"
textFindIndex("hello", "ll")       // → 2 (or -1 if not found)
textJoin(["a","b","c"], ", ")      // → "a, b, c"
textSplit("a,b,c", ",")            // → ["a","b","c"]
textToInteger("42")                // → 42
textToFloat("3.14")                // → 3.14
```

### Codepoint utilities

```candy
getCodepointCount("héllo")         // → 5 (Unicode-aware length)
loadCodepoints("ABC")              // → [65, 66, 67]  (Unicode values)
codepointToUTF8(65)                // → "A"
loadUTF8([72, 101, 108, 108, 111]) // → "Hello"
loadTextLines("line1\nline2")      // → ["line1","line2"]
```

---

## Images

Images live in CPU memory (RAM). You can load, edit, and generate them, then upload to a texture for drawing on screen.

### Loading and saving

```candy
let img = loadImage("player.png")          // load from file
let img = loadImageFromScreen()            // capture what is on screen right now
isImageValid(img)                          // → bool
exportImage(img, "output.png")             // save to disk
unloadImage(img)                           // free memory
```

### Generating images in code

```candy
let img = genImageColor(200, 200, "blue")
let img = genImageGradientLinear(200, 200, 0, "black", "white")  // 0=vertical
let img = genImageGradientRadial(200, 200, 0.5, "white", "black")
let img = genImageChecked(200, 200, 10, 10, "black", "white")
let img = genImageWhiteNoise(200, 200, 0.5)
let img = genImagePerlinNoise(200, 200, 0, 0, 4.0)
```

### Manipulating images (CPU-side)

All of these modify the image in place. None affect the screen — you need to upload to a texture to see the result.

```candy
imageCrop(img, x, y, width, height)        // crop to a region
imageResize(img, newW, newH)               // resize with bilinear filter
imageResizeNN(img, newW, newH)             // resize with nearest-neighbor (pixel art)
imageFlipHorizontal(img)
imageFlipVertical(img)
imageRotateCW(img)                         // rotate 90° clockwise
imageRotateCCW(img)
imageRotate(img, degrees)                  // rotate any angle
imageColorGrayscale(img)
imageColorInvert(img)
imageColorTint(img, color)
imageColorBrightness(img, value)           // -255 to 255
imageColorContrast(img, value)             // -100 to 100
imageColorReplace(img, fromColor, toColor)
imageBlurGaussian(img, blurSize)
imageMipmaps(img)
imageAlphaCrop(img, threshold)
imageAlphaClear(img, color, threshold)
imageAlphaMask(img, maskImageId)
imageAlphaPremultiply(img)
```

Copying and combining:

```candy
let copy = imageCopy(img)
let sub  = imageFromImage(img, x, y, w, h)    // crop without modifying original
imageDraw(dstImg, srcImg, srcRec, dstRec)      // paste src onto dst
```

Reading pixel data:

```candy
let pixels  = loadImageColors(img)             // → [{r,g,b,a}, ...]
let palette = loadImagePalette(img, maxColors) // → [{r,g,b,a}, ...]
let px      = getImageColor(img, x, y)         // → {r,g,b,a}
```

Drawing on an image (CPU):

```candy
imageClearBackground(img, "black")
imageDrawPixel(img, x, y, color)
imageDrawLine(img, x1, y1, x2, y2, color)
imageDrawCircle(img, cx, cy, radius, color)
imageDrawRectangle(img, x, y, w, h, color)
imageDrawText(img, "hello", x, y, fontSize, color)
imageDrawTextEx(img, fontId, "hello", {x,y}, fontSize, spacing, tint)
```

---

## Textures

Textures live in GPU memory (VRAM) and are what you actually draw on screen.

### Loading textures

```candy
let tex = loadTexture("player.png")          // load directly from file
let tex = loadTextureFromImage(imageId)      // upload an Image you already have
isTextureValid(tex)                          // → bool
unloadTexture(tex)                           // free GPU memory
```

### Drawing textures

```candy
drawTexture(tex, x, y)                       // simplest form
drawTexture(tex, x, y, "white")              // with tint color
drawTextureV(tex, {x,y}, tint)
drawTextureEx(tex, {x,y}, rotation, scale, tint)

// Draw a sub-region of the texture (sprite sheets!)
let src = rect(0, 0, 32, 32)                 // first frame of a sprite sheet
let dst = rect(100, 100, 64, 64)             // draw it scaled to 64x64
drawTexturePro(tex, src, dst, {x:0,y:0}, 0, "white")
```

### Texture settings

```candy
genTextureMipmaps(tex)             // generate mipmaps for distant rendering
setTextureFilter(tex, 1)           // 0=point, 1=bilinear, 2=trilinear
setTextureWrap(tex, 0)             // 0=repeat, 1=clamp, 2=mirror
```

### Render textures (draw to texture)

A render texture lets you draw into a texture instead of the screen. Useful for effects, minimaps, etc.

```candy
let rt = loadRenderTexture(400, 300)

// Inside your game loop:
beginTextureMode(rt)
  clear("black")
  drawCircle(200, 150, 50, "red")   // this goes onto the texture
endTextureMode()

// Then draw that texture to the screen like any other texture:
let tex = getRenderTextureTexture(rt)
drawTexture(tex, 0, 0)

unloadRenderTexture(rt)
```

---

## Shaders

Shaders are GPU programs that control how things are drawn. You load a vertex and fragment shader from GLSL files.

```candy
let sh = loadShader("shaders/vert.glsl", "shaders/frag.glsl")

while !shouldClose() {
  clear("black")
  beginShaderMode(sh)
    drawRectangle(100, 100, 200, 200, "white")  // drawn with the shader
  endShaderMode()
  flip()
}

unloadShader(sh)
```

---

## 3D graphics

### Setting up a 3D scene

```candy
window(800, 600, "3D Scene")
setTargetFPS(60)

while !shouldClose() {
  clear("raywhite")

  beginMode3D(
    0, 5, 10,    // camera position: x, y, z
    0, 0, 0,     // looking at point: x, y, z
    45           // field of view in degrees
  )

    drawGrid(10, 1)                             // ground grid
    drawCube(0, 0.5, 0, 1, 1, 1, "red")        // a red cube
    drawCubeWires(0, 0.5, 0, 1, 1, 1, "black") // wireframe outline

  endMode3D()
  flip()
}
closeWindow()
```

### Basic 3D shapes

All positions and sizes are in world-space 3D coordinates.

```candy
drawLine3D(x1,y1,z1, x2,y2,z2, color)
drawPoint3D(x, y, z, color)
drawCircle3D(cx,cy,cz, radius, axisX,axisY,axisZ, angle, color)
drawTriangle3D(x1,y1,z1, x2,y2,z2, x3,y3,z3, color)

drawCube(x, y, z, width, height, depth, color)
drawCubeV({x,y,z}, {x,y,z}, color)         // position and size as vec3 maps
drawCubeWires(x, y, z, w, h, d, color)
drawCubeWiresV({x,y,z}, {x,y,z}, color)

drawSphere(x, y, z, radius, color)
drawSphereEx(x, y, z, radius, rings, slices, color)
drawSphereWires(x, y, z, radius, rings, slices, color)

drawCylinder(x, y, z, radiusTop, radiusBottom, height, color)
drawCylinderEx(x1,y1,z1, x2,y2,z2, startRadius, endRadius, sides, color)
drawCylinderWires(x, y, z, radiusTop, radiusBottom, height, slices, color)
drawCylinderWiresEx(x1,y1,z1, x2,y2,z2, startRadius, endRadius, sides, color)

drawCapsule(x1,y1,z1, x2,y2,z2, radius, slices, rings, color)
drawCapsuleWires(x1,y1,z1, x2,y2,z2, radius, slices, rings, color)

drawPlane(x, y, z, width, length, color)   // flat XZ plane
drawRay(px,py,pz, dx,dy,dz, color)         // ray from position in direction
drawGrid(slices, spacing)                  // centered grid on the XZ plane
drawBoundingBox(minX,minY,minZ, maxX,maxY,maxZ, color)
```

### Loading and drawing 3D models

Load a `.obj`, `.glb`, `.gltf`, or other supported 3D file:

```candy
let model = loadModel("assets/character.glb")
isModelValid(model)    // → bool

// In your draw loop (inside beginMode3D ... endMode3D):
drawModel(model, x, y, z, scale, tint)
drawModelEx(model, x,y,z, rotX,rotY,rotZ, angle, scaleX,scaleY,scaleZ, tint)
drawModelWires(model, x, y, z, scale, tint)

// Bounding box of the whole model
let bb = getModelBoundingBox(model)   // → {min:{x,y,z}, max:{x,y,z}}

unloadModel(model)
```

### Billboards (2D sprites in 3D world)

A billboard always faces the camera — useful for trees, particles, UI elements in 3D:

```candy
drawBillboard(textureId, x, y, z, scale, tint)
drawBillboardRec(textureId, sx,sy,sw,sh, x,y,z, sizeX,sizeY, tint)
```

### Procedural meshes

Generate a mesh in code without loading a file:

```candy
let mesh = genMeshCube(1, 1, 1)             // unit cube
let mesh = genMeshSphere(0.5, 16, 16)       // sphere, 16 rings/slices
let mesh = genMeshPlane(10, 10, 4, 4)       // 10×10 plane, 4×4 subdivisions
let mesh = genMeshCylinder(0.5, 2, 16)
let mesh = genMeshCone(0.5, 2, 16)
let mesh = genMeshTorus(1, 0.3, 16, 16)
let mesh = genMeshKnot(1, 0.3, 16, 16)
let mesh = genMeshPoly(6, 1)                // hexagon
let mesh = genMeshHeightmap(imageId, {x:10,y:1,z:10})
let mesh = genMeshCubicmap(imageId, {x:1,y:1,z:1})

getMeshBoundingBox(meshId)                  // → {min, max}
exportMesh(meshId, "output.obj")
unloadMesh(meshId)
```

### Materials

```candy
let mat = loadMaterialDefault()             // basic white material
let mats = loadMaterials("model.mtl")       // → [materialId, ...]
isMaterialValid(mat)                        // → bool
setMaterialTexture(mat, 0, textureId)       // 0 = DIFFUSE map
setModelMeshMaterial(modelId, meshIdx, materialIdx)
unloadMaterial(mat)
```

### Animations

```candy
let anims = loadModelAnimations("character.glb")   // → animsId (group handle)
isModelAnimationValid(modelId, anims)               // → bool

// In your game loop:
frame = frame + 1
updateModelAnimation(modelId, anims, frame)

unloadModelAnimations(anims)
```

### 2D camera mode

Use `beginMode2D` for camera panning, zooming, and rotation in 2D games:

```candy
beginMode2D(offsetX, offsetY, targetX, targetY, rotation, zoom)
  // draw things here — they are affected by the camera
endMode2D()
```

### Raycasting and 3D collision

```candy
// Cast a ray from the mouse into the 3D scene
let r = getMouseRay()   // → {position:{x,y,z}, direction:{x,y,z}}

// Test what the ray hits
let hit = getRayCollisionSphere(ray, cx, cy, cz, radius)
let hit = getRayCollisionBox(ray, boundingBox)
let hit = getRayCollisionMesh(ray, meshId)
let hit = getRayCollisionTriangle(ray, p1, p2, p3)
let hit = getRayCollisionQuad(ray, p1, p2, p3, p4)
// hit → {hit:bool, distance:float, point:{x,y,z}, normal:{x,y,z}}

// Bounding-box and sphere collision
checkCollisionBoxes(min1x,min1y,min1z, max1x,max1y,max1z,
                    min2x,min2y,min2z, max2x,max2y,max2z) // → bool
checkCollisionSpheres(cx1,cy1,cz1,r1, cx2,cy2,cz2,r2)    // → bool
checkCollisionBoxSphere(boundingBox, cx, cy, cz, radius)  // → bool
```

---

## Audio

Always initialize the audio device before loading any sounds or music:

```candy
initAudioDevice()
// ... load sounds and music ...
closeAudioDevice()   // call at the very end of your program
isAudioDeviceReady() // → bool
setMasterVolume(0.8) // 0.0 = silent, 1.0 = full volume
getMasterVolume()    // → float
```

### Short sounds (Sound)

Sounds are loaded fully into memory and best for short effects like jumps, clicks, explosions:

```candy
let snd = loadSound("jump.wav")
let snd = loadSoundFromWave(waveId)     // create from a Wave handle
let alias = loadSoundAlias(snd)         // share sample data — no extra memory

playSound(snd)
stopSound(snd)
pauseSound(snd)
resumeSound(snd)
isSoundPlaying(snd)     // → bool
isSoundValid(snd)       // → bool

setSoundVolume(snd, 0.5)    // 0.0–1.0
setSoundPitch(snd, 1.5)     // 1.0 = normal pitch
setSoundPan(snd, -1.0)      // -1=left, 0=center, 1=right

unloadSound(snd)
unloadSoundAlias(alias)
```

### Music streams (Music)

Music is streamed from disk — better for long background tracks. **You must call `updateMusicStream` every frame:**

```candy
let bgm = loadMusicStream("music.ogg")
playMusicStream(bgm)

while !shouldClose() {
  updateMusicStream(bgm)    // ← required every frame or music will stutter
  clear("black")
  flip()
}

unloadMusicStream(bgm)
```

Other music functions:

```candy
isMusicStreamPlaying(bgm)       // → bool
isMusicValid(bgm)               // → bool
stopMusicStream(bgm)
pauseMusicStream(bgm)
resumeMusicStream(bgm)
seekMusicStream(bgm, 30.5)      // jump to 30.5 seconds
setMusicVolume(bgm, 0.7)
setMusicPitch(bgm, 1.0)
setMusicPan(bgm, 0.0)
getMusicTimeLength(bgm)         // → total duration in seconds
getMusicTimePlayed(bgm)         // → current position in seconds
```

### Wave data (Wave)

Waves are raw audio data in CPU memory. Load them to inspect or convert before making a Sound:

```candy
let wav = loadWave("sound.wav")
let wav = loadWaveFromMemory(".wav", bytesArray)
isWaveValid(wav)                    // → bool
let copy = waveCopy(wav)
waveCrop(wav, 0, 44100)             // keep only the first second (44100 frames)
waveFormat(wav, 44100, 16, 2)       // resample to 44100Hz, 16-bit stereo
let samples = loadWaveSamples(wav)  // → [float, ...] raw sample data
exportWave(wav, "output.wav")
unloadWave(wav)
```

### AudioStream (raw PCM streaming)

For advanced use — stream your own generated audio data in real time:

```candy
let stream = loadAudioStream(44100, 32, 2)  // sampleRate, sampleSize, channels
playAudioStream(stream)

// In your loop — fill the buffer whenever it is ready:
if isAudioStreamProcessed(stream) {
  updateAudioStream(stream, mySamplesArray)  // float32 samples
}

isAudioStreamValid(stream)    // → bool
isAudioStreamPlaying(stream)  // → bool
pauseAudioStream(stream)
resumeAudioStream(stream)
stopAudioStream(stream)
setAudioStreamVolume(stream, 1.0)
setAudioStreamPitch(stream, 1.0)
setAudioStreamPan(stream, 0.0)
setAudioStreamBufferSizeDefault(4096)  // set before loading streams
unloadAudioStream(stream)
```

Terminal bell: `beep()` (plays the system notification sound, not Raylib audio)

---

## Math utilities

These are helper functions for game math. They work on plain numbers and `{x,y}` / `{x,y,z}` / `{x,y,z,w}` maps.

### Scalar math

```candy
mathClamp(value, min, max)               // clamp value between min and max
mathLerp(start, end, 0.5)               // interpolate — 0.0=start, 1.0=end
mathNormalize(value, rangeStart, rangeEnd)  // map to 0.0–1.0
mathRemap(v, inStart, inEnd, outStart, outEnd)
mathWrap(value, min, max)               // wrap around like a modulo
floatEquals(a, b)                       // → bool (safe float comparison)
```

### Vector2 — `{x, y}`

```candy
vector2Zero()          // → {x:0, y:0}
vector2One()           // → {x:1, y:1}
vector2Add(a, b)
vector2Subtract(a, b)
vector2Scale(v, f)     // multiply by scalar
vector2Multiply(a, b)  // element-wise multiply
vector2Normalize(v)
vector2Length(v)       // → float
vector2Distance(a, b)  // → float
vector2DotProduct(a, b)  // → float
vector2Lerp(a, b, t)
vector2Rotate(v, angle)
vector2MoveTowards(v, target, maxDist)
vector2Reflect(v, normal)
vector2Clamp(v, min, max)
vector2Equals(a, b)    // → bool
```

### Vector3 — `{x, y, z}`

```candy
vector3Zero() / vector3One()
vector3Add(a, b) / vector3Subtract(a, b)
vector3Scale(v, f) / vector3Multiply(a, b) / vector3Divide(a, b)
vector3Normalize(v)
vector3CrossProduct(a, b)   // perpendicular vector
vector3DotProduct(a, b)     // → float
vector3Length(v)            // → float
vector3Distance(a, b)       // → float
vector3Lerp(a, b, t)
vector3RotateByAxisAngle(v, axis, angle)
vector3Transform(v, matrix)
vector3Min(a, b) / vector3Max(a, b)
vector3Equals(a, b)         // → bool
```

### Matrix — `{m0, m1, … m15}`

4×4 matrices used for 3D transforms:

```candy
matrixIdentity()
matrixTranslate(x, y, z)
matrixScale(x, y, z)
matrixRotate(axis, angle)     // axis is {x,y,z}, angle is radians
matrixRotateX(angle) / matrixRotateY(angle) / matrixRotateZ(angle)
matrixMultiply(a, b)
matrixTranspose(m)
matrixInvert(m)
matrixLookAt(eye, target, up)          // view matrix
matrixPerspective(fovY, aspect, near, far)
matrixOrtho(left, right, bottom, top, near, far)
matrixDecompose(m)  // → {translation:{x,y,z}, rotation:{x,y,z,w}, scale:{x,y,z}}
```

### Quaternion — `{x, y, z, w}`

Quaternions represent rotations without gimbal lock:

```candy
quaternionIdentity()
quaternionFromAxisAngle(axis, angle)
quaternionFromEuler(pitch, yaw, roll)
quaternionToEuler(q)           // → {x,y,z} in radians
quaternionMultiply(a, b)
quaternionNormalize(q)
quaternionSlerp(a, b, t)       // smooth rotation interpolation
quaternionToMatrix(q)          // → matrix map
quaternionEquals(a, b)         // → bool
```

---

## Easing functions

Easing functions smoothly interpolate a value from 0.0 to 1.0. Useful for animations and transitions.

```candy
// All take (currentTime, startValue, changeInValue, duration)
easeLinear(t, b, c, d)
easeSineIn(t, b, c, d)
easeSineOut(t, b, c, d)
easeSineInOut(t, b, c, d)
easeCubicIn(t, b, c, d)
easeCubicOut(t, b, c, d)
easeBounceOut(t, b, c, d)
```

---

## Physics

Candy includes a built-in 3D physics engine written in pure Go. It works on **Windows, macOS, and Linux** with no extra setup — just use the functions below.

Features: rigid bodies, gravity, AABB/sphere/capsule collision, impulse-based collision response, friction, restitution (bounciness), angular velocity, drag, sensors, raycasting.

### Quick start

```candy
// 1. Create a world (gravity x, y, z)
let world = physicsCreateWorld(0, -9.81, 0)

// 2. Add a static ground plane
let ground = physicsCreatePlane(world, 0)

// 3. Add a dynamic falling box: position x,y,z then half-extents w,h,d
let box = physicsCreateBox(world, 0, 10, 0,  1, 1, 1)

// 4. Step the simulation every frame
while !shouldClose() {
    physicsStep(world, getFrameTime())

    // Read position and draw
    let p = physicsGetPosition(world, box)
    clear("skyblue")
    beginMode3D(0, 8, 12,  0, 0, 0,  45)
        drawGrid(10, 1)
        drawCube(p.x, p.y, p.z, 2, 2, 2, "red")
    endMode3D()
    flip()
}

// 5. Clean up
physicsDestroyWorld(world)
```

### Motion types

Pass the motion type as the 8th argument to `physicsCreateBox`, `physicsCreateSphere`, and `physicsCreateCapsule`:

| Constant | Value | Meaning |
|---|---|---|
| `PHYSICS_STATIC()` | `0` | Never moves, infinite mass — use for floors, walls |
| `PHYSICS_DYNAMIC()` | `1` | Fully simulated by forces and gravity (default) |
| `PHYSICS_KINEMATIC()` | `2` | Moved by you, not by forces — use for platforms, elevators |

```candy
let floor  = physicsCreateBox(world, 0, 0, 0,  10, 0.5, 10,  PHYSICS_STATIC())
let player = physicsCreateCapsule(world, 0, 5, 0,  0.5, 1,  PHYSICS_KINEMATIC())
let ball   = physicsCreateSphere(world, 0, 10, 0,  1)  // dynamic by default
```

### World management

```candy
let world = physicsCreateWorld(0, -9.81, 0)   // create world, returns worldId
physicsDestroyWorld(world)                     // free world and all its bodies
physicsStep(world, dt)                         // advance simulation by dt seconds
physicsSetGravity(world, 0, -20, 0)            // change gravity
let g = physicsGetGravity(world)               // → {x, y, z}
```

### Creating bodies

```candy
// physicsCreateBox(world, x, y, z, halfW, halfH, halfD, [motionType], [isSensor])
let box = physicsCreateBox(world, 0, 5, 0,  1, 1, 1)

// physicsCreateSphere(world, x, y, z, radius, [motionType], [isSensor])
let ball = physicsCreateSphere(world, 0, 5, 0,  0.5)

// physicsCreateCapsule(world, x, y, z, radius, halfHeight, [motionType], [isSensor])
let cap = physicsCreateCapsule(world, 0, 5, 0,  0.4, 1)

// physicsCreatePlane(world, [yOffset])  — infinite static ground
let ground = physicsCreatePlane(world, 0)

// Destroy a body
physicsDestroyBody(world, box)
```

> **isSensor** — pass `true` to make a body a trigger zone. It detects overlaps but doesn't physically push other bodies. Use for pickup areas, checkpoints, damage zones.

### Reading and writing body state

```candy
let pos = physicsGetPosition(world, body)        // → {x, y, z}
physicsSetPosition(world, body, 3, 0, -1)

let vel = physicsGetVelocity(world, body)         // → {x, y, z}
physicsSetVelocity(world, body, 0, 5, 0)

let avel = physicsGetAngularVelocity(world, body) // → {x, y, z}
physicsSetAngularVelocity(world, body, 0, 3, 0)

let rot = physicsGetRotation(world, body)         // → {x, y, z, w} quaternion
```

### Forces and impulses

```candy
// Force accumulates over the frame and clears each step — good for continuous thrust
physicsApplyForce(world, body, 0, 200, 0)

// Impulse is an instant velocity change — good for jumps, explosions, hits
physicsApplyImpulse(world, body, 0, 8, 0)

// Torque spins a body (angular force)
physicsApplyTorque(world, body, 0, 1, 0)
```

### Material properties

```candy
physicsSetMass(world, body, 5.0)         // kg — 0 makes it static
physicsGetMass(world, body)              // → float

physicsSetRestitution(world, body, 0.8)  // 0=no bounce, 1=perfect bounce
physicsSetFriction(world, body, 0.5)     // 0=ice, 1=rubber

physicsSetLinearDrag(world, body, 0.02)  // air resistance on velocity
physicsSetAngularDrag(world, body, 0.1)  // spin damping
```

### Activating and tagging bodies

```candy
physicsSetActive(world, body, false)   // deactivate — body is skipped in simulation
physicsIsActive(world, body)           // → bool

physicsSetUserData(world, body, 42)    // attach any integer tag (e.g. enemy type, score value)
physicsGetUserData(world, body)        // → int
```

### Contacts (collision events)

After each `physicsStep`, call this to get everything that collided that frame:

```candy
let contacts = physicsGetContacts(world)
let ci = 0
while ci < len(contacts) {
    let c = contacts[ci]
    // c.bodyA, c.bodyB  — the two body IDs that collided
    // c.normalX, c.normalY, c.normalZ  — collision normal
    // c.depth   — penetration depth
    // c.pointX, c.pointY, c.pointZ  — contact point
    ci = ci + 1
}
```

### Raycasting

Cast a ray and get a sorted list of all bodies it hits:

```candy
// physicsCastRay(world, ox, oy, oz, dx, dy, dz, maxDistance)
let hits = physicsCastRay(world,  0, 10, 0,  0, -1, 0,  100)
// hits is an array of: {bodyId, distance, pointX, pointY, pointZ, normalX, normalY, normalZ}

// Closest hit only (or null if nothing hit)
let hit = physicsCastRayFirst(world,  0, 10, 0,  0, -1, 0,  100)
if hit != null {
    drawText(10, 10, "Hit body: " + hit.bodyId, 20, "white")
}
```

---

## Example programs

The `compiler/scratch/` folder contains ready-to-run examples:

| File | What it shows |
|---|---|
| `raylib_input_demo.candy` | Keyboard and mouse input |
| `raylib_texture_demo.candy` | Loading and drawing textures |
| `raylib_text_demo.candy` | Custom fonts and text drawing |
| `raylib_image_demo.candy` | Image loading and CPU manipulation |
| `raylib_audio_demo.candy` | Playing sounds and music |
| `raylib_3d_demo.candy` | 3D shapes and camera |
| `raylib_billboard_demo.candy` | Billboards in 3D space |
| `raylib_showcase_3d.candy` | Full 3D scene showcase |
| `raylib_window_cursor_demo.candy` | Window state and cursor control |
| `physics_demo.candy` | Rigid body physics — boxes, spheres, bouncing |

---

## Common mistakes

| Problem | Fix |
|---|---|
| Window closes instantly | Wrap your drawing code in `while !shouldClose() { ... flip() }` |
| Nothing appears on screen | Make sure you call `clear()` and `flip()` every frame |
| Sound/music doesn't play | Call `initAudioDevice()` before loading any audio |
| Music cuts out or stutters | Call `updateMusicStream(bgm)` every single frame inside the loop |
| "Invalid handle" error | You may have called `unloadTexture` / `unloadSound` etc. more than once, or used a handle after unloading it |
| Image loads but draws black | Make sure you uploaded it to a texture: `loadTextureFromImage(imageId)` |
| 3D scene is empty | All 3D drawing must be inside `beginMode3D(...)` / `endMode3D()` |
| Text/model missing | Double-check the file path — paths are relative to where you run Candy from |

---

## ENET networking module

Candy includes an `enet` module with host/peer/packet APIs and normalized events.

### Core API

- `enet.init()` / `enet.deinit()` / `enet.version()` / `enet.backend()`
- `enet.address(host, port)` returns `{host, port}`
- `enet.host_create(addressOrNull, peerLimit, channelLimit, inBandwidth, outBandwidth)` returns `hostId`
- `enet.host_destroy(hostId)`
- `enet.host_service(hostId, timeoutMs)` returns event map
- `enet.host_flush(hostId)`
- `enet.host_bandwidth_limit(hostId, inBandwidth, outBandwidth)`
- `enet.host_channel_limit(hostId, channels)`
- `enet.host_compress_range_coder(hostId)` and alias `enet.set_range_coder(hostId)`
- `enet.host_connect(hostId, address, channels, data)` returns `peerId`
- `enet.peer_disconnect(peerId, data)` / `enet.peer_disconnect_now(peerId, data)`
- `enet.peer_ping(peerId)` / `enet.peer_timeout(peerId, limit, min, max)` / `enet.peer_reset(peerId)`
- `enet.packet_create(dataString, flags)` returns `packetId`
- `enet.packet_destroy(packetId)`
- `enet.peer_send(peerId, channel, packetId)`

### Event shape from `host_service`

`{"type","hostId","peerId","channel","data","packet","address"}`

- `packet`: `{"id","flags","data"}`
- `address`: `{"host","port"}`

### Constants

- Event types: `enet.EVENT_NONE`, `enet.EVENT_CONNECT`, `enet.EVENT_DISCONNECT`, `enet.EVENT_RECEIVE`
- Packet flags: `enet.PACKET_RELIABLE`, `enet.PACKET_UNSEQUENCED`, `enet.PACKET_NO_ALLOCATE`, `enet.PACKET_UNRELIABLE`

### Demos

- `compiler/scratch/enet_server_demo.candy`
- `compiler/scratch/enet_client_demo.candy`
