package candy_stdlib

func init() {
	Modules["candy.proc"] = `
object Noise {
    fun perlin(x, y = 0.0, z = 0.0) {
        return noise(x, y, z)
    }
    
    fun fractal(x, y = 0.0, z = 0.0, octaves = 4, persistence = 0.5) {
        var total = 0.0
        var freq = 1.0
        var amp = 1.0
        var maxVal = 0.0
        for i in 0..octaves {
            total = total + noise(x * freq, y * freq, z * freq) * amp
            maxVal = maxVal + amp
            amp = amp * persistence
            freq = freq * 2.0
        }
        return total / maxVal
    }
}

class DungeonGenerator {
    var width = 50
    var height = 50
    var grid = []
    
    fun init(w = 50, h = 50) {
        width = w
        height = h
        var wh = width * height
        grid = array(wh)
        for i in 0..wh { grid[i] = 1; } // 1 = wall
    }
    
    fun generate() {
        // Simple BSP-like room carving
        carveRoom(2, 2, width - 4, height - 4)
        for i in 0..10 {
            var rw = rand.float(4, 10)
            var rh = rand.float(4, 10)
            var rx = rand.float(1, width - rw - 1)
            var ry = rand.float(1, height - rh - 1)
            carveRoom(rx, ry, rw, rh)
        }
    }
    
    fun carveRoom(x, y, w, h) {
        for iy in y..(y + h) {
            for ix in x..(x + w) {
                grid[iy * width + ix] = 0 // 0 = floor
            }
        }
    }
    
    fun getTile(x, y) {
        if x < 0 or x >= width or y < 0 or y >= height { return 1; }
        return grid[y * width + x]
    }
}
`
}
