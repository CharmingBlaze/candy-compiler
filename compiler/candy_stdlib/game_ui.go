package candy_stdlib

func init() {
	Modules["candy.ui"] = `
class UIElement {
    var position = vec2(0, 0)
    var size = vec2(100, 30)
    var visible = true
    var anchor = 0 // TOP_LEFT
    var margin = vec2(0, 0)

    fun draw() {}
    fun hide() { visible = false; }
    fun show() { visible = true; }
}

class Canvas {
    var elements = []

    fun add(el) {
        elements.add(el)
    }

    fun draw() {
        for el in elements {
            el.draw()
        }
    }
}

class Button extends UIElement {
    var text = "Button"
    var onClick = null

    fun init(props = {}) {
        this.position = props.position
        if this.position == null { this.position = vec2(0,0); }
        this.size = props.size
        if this.size == null { this.size = vec2(100,30); }
        this.text = props.text
        if this.text == null { this.text = "Button"; }
        this.onClick = props.onClick
    }

    fun draw() {
        if !visible { return; }
        // drawRectangle(position.x, position.y, size.x, size.y, Color.GRAY)
        drawText(text, position.x + 10, position.y + 10, 20, COLOR_WHITE)
        
        if isMouseButtonPressed(0) {
            var m = getMousePosition()
            if m.x >= position.x and m.x <= position.x + size.x and m.y >= position.y and m.y <= position.y + size.y {
                if onClick != null { onClick(); }
            }
        }
    }
}

class Label extends UIElement {
    var text = "Label"
    var fontSize = 20
    var color = COLOR_WHITE

    fun init(props = {}) {
        this.position = props.position
        if this.position == null { this.position = vec2(0,0); }
        this.text = props.text
        if this.text == null { this.text = ""; }
        this.fontSize = props.fontSize
        if this.fontSize == null { this.fontSize = 20; }
        this.color = props.color
        if this.color == null { this.color = COLOR_WHITE; }
    }

    fun draw() {
        if !visible { return; }
        drawText(text, position.x, position.y, fontSize, color)
    }
}

class Panel extends UIElement {
    var color = COLOR_GRAY
    var elements = []

    fun init(props = {}) {
        this.position = props.position
        if this.position == null { this.position = vec2(0,0); }
        this.size = props.size
        if this.size == null { this.size = vec2(400,300); }
        this.color = props.color
        if this.color == null { this.color = COLOR_GRAY; }
    }

    fun add(el) {
        elements.add(el)
    }

    fun draw() {
        if !visible { return; }
        // drawRectangle(position.x, position.y, size.x, size.y, color)
        for el in elements {
            el.draw()
        }
    }
}

class ProgressBar extends UIElement {
    var value = 0.5
    var color = COLOR_GREEN

    fun draw() {
        if !visible { return; }
        // drawProgress
    }
}

class Slider extends UIElement {
    var min = 0.0
    var max = 100.0
    var value = 50.0
    var onChange = null

    fun draw() {
        if !visible { return; }
        // drawSlider
    }
}

class TextInput extends UIElement {
    var text = ""
    var placeholder = ""
}

class Checkbox extends UIElement {
    var checked = false
    var onChange = null
}

class Image extends UIElement {
    var texture = null
}

class VBoxLayout {
    var spacing = 10
    var padding = 20
    var elements = []

    fun add(el) {
        elements.add(el)
    }

    fun draw() {
        // layout logic
    }
}

class HBoxLayout extends VBoxLayout {}

class GridLayout {
    var columns = 2
    var spacing = 10
    var elements = []

    fun add(el) { elements.add(el); }
}

object Anchor {
    var TOP_LEFT = 0
    var TOP_RIGHT = 1
    var BOTTOM_LEFT = 2
    var BOTTOM_RIGHT = 3
    var CENTER = 4
}
`
}
