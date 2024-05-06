package typescript

type Options struct {
}

// Renderer provides a go renderer.
type Renderer struct {
	opts Options
}

// NewRenderer creates a new Renderer instance.
func NewRenderer(opts Options) *Renderer {
	return &Renderer{opts: opts}
}
