# raygui — Immediate-Mode GUI

Candy exposes the full [raygui](https://github.com/raysan5/raygui) library as built-in functions.  
All `gui*` functions must be called **between** `beginDrawing()` / `endDrawing()`.

---

## Quick-start example

```candy
initWindow(800, 600, "My UI")

var clicked = false

while !windowShouldClose() {
    beginDrawing()
    clearBackground("darkgray")

    if guiButton(10, 10, 120, 40, "Click me!") {
        clicked = true
    }

    if clicked {
        guiLabel(10, 60, 200, 30, "Button was clicked!")
    }

    endDrawing()
}
```

---

## Coordinate conventions

Every control that draws to the screen accepts its bounds either as  
**4 flat numbers** `(x, y, width, height)` **or** a **map** `{x, y, width, height}`.

```candy
guiLabel(10, 20, 200, 30, "Hello")          # flat form
guiLabel({x:10, y:20, width:200, height:30}, "Hello")  # map form
```

---

## Global state

### `guiEnable()`
Re-enable all GUI controls after a `guiDisable()` call.

### `guiDisable()`
Disable all controls globally (they will be drawn in a grayed-out state and won't respond to input).

### `guiLock()`
Lock the GUI — controls are drawn but inputs are ignored.

### `guiUnlock()`
Unlock the GUI.

### `guiIsLocked() → bool`
Returns `true` when the GUI is locked.

### `guiSetAlpha(alpha)`
Set global alpha (transparency) for all controls. `alpha` is a float in `0.0 .. 1.0`.

### `guiSetState(state)` / `guiGetState() → int`
Manually override the current GUI state.  
State constants (from raygui):

| Value | Meaning  |
|-------|----------|
| 0     | NORMAL   |
| 1     | FOCUSED  |
| 2     | PRESSED  |
| 3     | DISABLED |

### `guiSetFont(fontId)` / `guiGetFont() → fontId`
Set / get the font used by all GUI controls.  
`fontId` is a handle returned by `loadFont()` or `loadFontEx()`.

```candy
var fid = loadFont("assets/myFont.ttf")
guiSetFont(fid)
```

---

## Style

### `guiSetStyle(control, property, value)`
Set a style property for a specific control.

### `guiGetStyle(control, property) → int`
Get a style property value.

### `guiLoadStyle(fileName)`
Load a `.rgs` style file from disk.

### `guiLoadStyleDefault()`
Reset all styles to the built-in default.

### `guiSetTooltip(text)` / `guiEnableTooltip()` / `guiDisableTooltip()`
Configure the tooltip that appears on hover.

```candy
guiEnableTooltip()
guiSetTooltip("This button does something awesome")
guiButton(10, 10, 120, 30, "Hover me")
```

### `guiSetIconScale(scale)`
Set the pixel scale multiplier for built-in icons.

---

## Basic controls

### `guiLabel(x, y, w, h, text)`
Draw a non-interactive text label.

### `guiButton(x, y, w, h, text) → bool`
Draw a button. Returns `true` on the frame it is clicked.

```candy
if guiButton(10, 10, 100, 30, "Start") {
    startGame()
}
```

### `guiLabelButton(x, y, w, h, text) → bool`
Like `guiButton` but styled as a hyperlink label.

### `guiToggle(x, y, w, h, text, active) → bool`
A two-state toggle button. Pass the current state as `active`; the returned value is the new state.

```candy
var muted = false
muted = guiToggle(10, 50, 80, 30, "Mute", muted)
```

### `guiToggleGroup(x, y, w, h, text, active) → int`
Mutually-exclusive toggle buttons. `text` is semicolon-separated item labels.  
Returns the index of the currently selected item.

```candy
var tab = 0
tab = guiToggleGroup(10, 10, 300, 30, "Play;Options;Credits", tab)
```

### `guiToggleSlider(x, y, w, h, text, active) → int`
Like `guiToggleGroup` but rendered as a sliding selector.

### `guiCheckBox(x, y, w, h, text, checked) → bool`
A checkbox. Returns the new checked state.

```candy
var fullscreen = false
fullscreen = guiCheckBox(10, 90, 20, 20, "Fullscreen", fullscreen)
```

### `guiComboBox(x, y, w, h, text, active) → int`
A combo-box (closed dropdown). `text` is semicolon-separated items.  
Returns the index of the selected item.

```candy
var quality = 0
quality = guiComboBox(10, 120, 160, 30, "Low;Medium;High", quality)
```

---

## Stateful controls (active handles)

Some controls need persistent state across frames. Use `guiNewActive()` to allocate a handle:

### `guiNewActive() → handle`
Allocate an integer state slot.

### `guiGetActive(handle) → int`
Read the current value.

### `guiSetActive(handle, value)`
Write a value.

### `guiDropdownBox(x, y, w, h, text, activeHandle, editMode) → bool`
An openable dropdown. `editMode` should be `true` while the list is open.  
Returns `true` when the selection changes.

```candy
var ddActive = guiNewActive()
var ddOpen = false

# inside draw loop:
if guiDropdownBox(10, 60, 160, 30, "Option A;Option B;Option C", ddActive, ddOpen) {
    ddOpen = !ddOpen
}
var selected = guiGetActive(ddActive)
```

### `guiSpinner(x, y, w, h, text, activeHandle, min, max, editMode) → bool`
An integer spinner. Returns `true` while the user is editing.

```candy
var spinHandle = guiNewActive()
guiSetActive(spinHandle, 5)      # starting value

var editing = false
if guiSpinner(10, 100, 120, 30, "Count", spinHandle, 0, 20, editing) {
    editing = !editing
}
var count = guiGetActive(spinHandle)
```

### `guiValueBox(x, y, w, h, text, activeHandle, min, max, editMode) → bool`
Like `guiSpinner` but displayed as a typed number field.

---

## Sliders

### `guiSlider(x, y, w, h, textLeft, textRight, value, min, max) → float`
Returns the new value after the user drags the thumb.

```candy
var vol = 0.5
vol = guiSlider(10, 140, 200, 20, "0", "1", vol, 0.0, 1.0)
```

### `guiSliderBar(x, y, w, h, textLeft, textRight, value, min, max) → float`
Like `guiSlider` but filled from the left edge.

### `guiProgressBar(x, y, w, h, textLeft, textRight, value, min, max) → float`
Read-only progress indicator. Returns `value` unchanged.

```candy
guiProgressBar(10, 170, 200, 20, "0%", "100%", loadProgress, 0.0, 1.0)
```

### `guiScrollBar(x, y, w, h, value, min, max) → int`
A standalone integer scroll bar.

---

## Text input

Text box content is managed through string-slot handles.

### `guiNewTextBox(initialText?) → handle`
Create a text buffer, optionally pre-filled.

### `guiGetText(handle) → string`
Read the current text.

### `guiSetText(handle, text)`
Overwrite the text.

### `guiTextBox(x, y, w, h, textHandle, maxLen, editMode) → bool`
Renders an editable text field. Returns `true` while focus is active.

```candy
var nameBox = guiNewTextBox("Player")
var editing = false

# inside draw loop:
if guiTextBox(10, 200, 200, 30, nameBox, 32, editing) {
    editing = !editing
}
var playerName = guiGetText(nameBox)
```

---

## Containers and layout

### `guiWindowBox(x, y, w, h, title) → bool`
Draws a draggable window panel. Returns `true` when the close button is clicked.

```candy
var open = true

if open {
    if guiWindowBox(100, 100, 300, 200, "Settings") {
        open = false
    }
    # draw controls inside here
}
```

### `guiGroupBox(x, y, w, h, text)`
Draw a labeled rectangle outline — useful for grouping controls visually.

### `guiLine(x, y, w, h, text)`
Draw a horizontal separator line with an optional label.

### `guiPanel(x, y, w, h, text)`
Draw a filled panel background.

### `guiStatusBar(x, y, w, h, text)`
Draw a status-bar strip (typically at the bottom of a window).

### `guiDummyRec(x, y, w, h, text)`
Draw a placeholder rectangle — useful during layout.

### `guiScrollPanel(x, y, w, h, text, contentW, contentH, scrollHandle) → {x,y,width,height}`
Render a scrollable viewport over content of size `contentW × contentH`.  
Returns the visible rectangle.  
Use `guiNewScroll()` / `guiGetScroll()` to manage the scroll position.

```candy
var scroll = guiNewScroll()

var view = guiScrollPanel(10, 10, 400, 300, "Items", 400, 800, scroll)
# draw items clipped to 'view'
```

### `guiNewScroll() → handle`
Allocate a scroll-position state slot.

### `guiGetScroll(handle) → {x, y}`
Read the current scroll offset.

---

## List and grid

### `guiListView(x, y, w, h, text, scrollHandle, active) → int`
A scrollable list. `text` is semicolon-separated items.  
Returns the currently selected item index.

```candy
var listScroll = guiNewActive()
var selected = 0
selected = guiListView(10, 10, 200, 150, "Sword;Shield;Potion;Arrow", listScroll, selected)
```

### `guiGrid(x, y, w, h, text, spacing, subdivs) → {x, y}`
Draw a grid overlay. Returns the cell position the mouse is currently hovering over (`{-1,-1}` when outside).

```candy
var cell = guiGrid(0, 0, 400, 400, "", 32, 1)
if cell.x >= 0 {
    drawText("Hovering: " + cell.x + ", " + cell.y, 10, 410, 16, "white")
}
```

### `guiTabBar(x, y, w, h, tabs, active) → int`
A tab bar. `tabs` is a Candy array of strings.  
Returns the index of the active tab.

```candy
var tabs = ["General", "Controls", "Graphics"]
var activeTab = 0
activeTab = guiTabBar(10, 10, 400, 30, tabs, activeTab)
```

---

## Dialogs

### `guiMessageBox(x, y, w, h, title, message, buttons) → int`
A modal message box. `buttons` is semicolon-separated (e.g. `"OK"` or `"Yes;No"`).  
Returns `-1` if no button clicked, or the button index otherwise.

```candy
var result = guiMessageBox(200, 150, 300, 150, "Quit?", "Are you sure?", "Yes;No")
if result == 0 {
    closeWindow()
}
```

### `guiTextInputBox(x, y, w, h, title, message, buttons, textHandle) → int`
A dialog with an embedded text field. Works like `guiMessageBox` but also captures typed input.

```candy
var inputBox = guiNewTextBox("")
var btn = guiTextInputBox(150, 100, 400, 200, "Name", "Enter your name:", "OK;Cancel", inputBox)
if btn == 0 {
    var name = guiGetText(inputBox)
}
```

---

## Color pickers

### `guiColorPicker(x, y, w, h, text, color) → {r,g,b,a}`
Full HSV color picker widget. Pass the current color; returns the newly selected color.

```candy
var col = {r:255, g:0, b:128, a:255}
col = guiColorPicker(10, 10, 200, 200, "Pick color", col)
drawRectangle(220, 10, 50, 50, col)
```

### `guiColorPanel(x, y, w, h, text, color) → {r,g,b,a}`
Color panel (SV square only, without hue/alpha sliders).

### `guiColorBarAlpha(x, y, w, h, text, alpha) → float`
Alpha slider. Returns the new alpha value (0.0 – 1.0).

### `guiColorBarHue(x, y, w, h, text, value) → float`
Hue slider. Returns the new hue value (0 – 360).

---

## Drawing utilities

### `guiDrawIcon(iconId, x, y, pixelSize, color)`
Draw one of the 200+ built-in raygui icons.  
Common icon IDs:

| ID  | Icon            |
|-----|-----------------|
| 1   | FILE_NEW        |
| 2   | FILE_SAVE       |
| 5   | FILE_OPEN       |
| 131 | ARROW_LEFT      |
| 132 | ARROW_RIGHT     |
| 133 | ARROW_UP        |
| 134 | ARROW_DOWN      |

```candy
guiDrawIcon(1, 10, 10, 2, "white")
```

### `guiIconText(iconId, text) → string`
Prepend an icon code to a text string for use in control labels.

```candy
guiButton(10, 10, 120, 30, guiIconText(1, "New File"))
```

### `guiLoadIcons(fileName, loadNames)`
Load custom icons from a `.rgi` file.

### `guiDrawRectangle(x, y, w, h, borderWidth, borderColor, fillColor)`
Draw a styled rectangle (raygui style-aware).

### `guiDrawText(text, x, y, w, h, alignment, color)`
Draw text inside a rectangle. Alignment: `0`=left, `1`=center, `2`=right.

### `guiGetTextBounds(control, x, y, w, h) → {x,y,width,height}`
Get the text bounds rectangle for a control type.

### `guiGetTextWidth(text) → int`
Measure text width using the current GUI font.

### `guiFade(color, alpha) → {r,g,b,a}`
Return a copy of `color` with its alpha multiplied by `alpha`.

### `guiGetColor(control, property) → {r,g,b,a}`
Get a style color property as a color map.

---

## Icon ID quick reference

raygui ships with 256 icons. You can browse them at:  
<https://raylibtech.itch.io/rguiicons>

Some frequently-used ones:

| ID  | Name             | ID  | Name           |
|-----|------------------|-----|----------------|
| 1   | FILE_NEW         | 131 | ARROW_LEFT      |
| 2   | FILE_SAVE        | 132 | ARROW_RIGHT     |
| 5   | FILE_OPEN        | 133 | ARROW_UP        |
| 9   | TRASH            | 134 | ARROW_DOWN      |
| 18  | INFO             | 140 | CURSOR_HAND     |
| 19  | WARNING          | 191 | MUSIC           |
| 20  | HELP             | 220 | PLAYER_PLAY     |
| 100 | UNDO             | 222 | PLAYER_STOP     |
| 101 | REDO             | 232 | VOLUME_UP       |

---

## Full working example — Settings screen

```candy
initWindow(800, 600, "Settings")
setTargetFPS(60)

var vol = 0.7
var fxVol = 1.0
var muted = false
var quality = 1
var nameBox = guiNewTextBox("Hero")
var nameEditing = false
var ddActive = guiNewActive()
var ddOpen = false
var windowOpen = true
var activeTab = 0

while !windowShouldClose() {
    beginDrawing()
    clearBackground("darkgray")

    if windowOpen {
        if guiWindowBox(100, 50, 600, 480, "Settings") {
            windowOpen = false
        }

        activeTab = guiTabBar(110, 90, 580, 32, ["Audio", "Display", "Player"], activeTab)

        if activeTab == 0 {
            guiLabel(120, 140, 120, 24, "Master Vol")
            vol = guiSliderBar(250, 140, 200, 24, "0", "1", vol, 0.0, 1.0)

            guiLabel(120, 175, 120, 24, "FX Vol")
            fxVol = guiSliderBar(250, 175, 200, 24, "0", "1", fxVol, 0.0, 1.0)

            muted = guiCheckBox(250, 210, 20, 20, "Mute all", muted)
        }

        if activeTab == 1 {
            guiLabel(120, 140, 120, 24, "Quality")
            quality = guiComboBox(250, 140, 180, 30, "Low;Medium;High;Ultra", quality)

            guiLabel(120, 180, 200, 24, "Resolution")
            if guiDropdownBox(250, 180, 200, 30, "1280x720;1920x1080;2560x1440", ddActive, ddOpen) {
                ddOpen = !ddOpen
            }
        }

        if activeTab == 2 {
            guiLabel(120, 140, 120, 24, "Name")
            if guiTextBox(250, 140, 200, 30, nameBox, 20, nameEditing) {
                nameEditing = !nameEditing
            }
        }

        if guiButton(490, 490, 100, 32, "Apply") {
            # save settings...
        }
    }

    endDrawing()
}
```

---

## Notes

- All `gui*` functions require the `raylib` build tag (they are part of `candy_raylib`).
- State handles (`guiNewActive`, `guiNewTextBox`, `guiNewScroll`) allocate memory that persists for the lifetime of the process. Create them once, outside the game loop.
- raygui draws everything immediately — there are no retained-mode objects.
- Load custom `.rgs` style files with `guiLoadStyle("mytheme.rgs")` to completely change the look.
