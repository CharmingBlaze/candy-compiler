package candy_stdlib

func init() {
	Modules["candy.vfx"] = `
class PostProcess {
    var enabled = false
    var effects = []
    
    fun init() {
        // Setup internal render textures if supported by runtime
    }
    
    fun enable(effect) {
        enabled = true
        effects.add(effect)
    }
    
    fun disable(effect) {
        effects.remove(effect)
        if len(effects) == 0 { enabled = false; }
    }
    
    fun begin() {
        if !enabled { return; }
        // beginDrawingToRenderTexture
    }
    
    fun finish() {
        if !enabled { return; }
        // endDrawingToRenderTexture
        // for e in effects { drawWithShader(e); }
    }
}

object Effects {
    var Bloom = "bloom"
    var Blur = "blur"
    var GrayScale = "grayscale"
    var Sepia = "sepia"
    var Vignette = "vignette"
}
`
}
