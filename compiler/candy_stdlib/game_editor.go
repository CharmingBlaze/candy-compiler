package candy_stdlib

func init() {
	Modules["candy.editor"] = `
object Editor {
    var enabled = false
    var watchedEntities = []

    fun enable() {
        enabled = true
        print("Editor enabled")
    }

    fun watch(entity) {
        watchedEntities.add(entity)
    }

    fun draw() {
        if !enabled { return; }
        
        drawRectangle(10, 10, 300, 400, color.rgba(0, 0, 0, 150))
        drawText("CANDY RUNTIME INSPECTOR", 20, 20, 20, color.gold)
        
        var y = 60
        for e in watchedEntities {
            drawText("Entity: {e.name}", 30, y, 18, color.white)
            y = y + 25
            if e.position != null {
                drawText("  pos: {e.position.x}, {e.position.y}", 40, y, 16, color.gray)
                y = y + 20
            }
        }
    }
}

object Gizmos {
    fun drawLine(start, end, col = "red") {
        drawLine3D(start, end, col)
    }
    
    fun drawSphere(pos, radius, col = "green") {
        drawSphereWires(pos, radius, 8, 8, col)
    }
}
`
}
