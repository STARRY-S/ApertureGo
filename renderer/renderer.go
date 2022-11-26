// Package renderer has the built-in implemented renderer struct object.
package renderer

import (
	"fmt"

	ap "github.com/STARRY-S/aperture"
	"github.com/STARRY-S/aperture/utils"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type RendererObj struct {
	name         string
	viewDistance int32

	// windows stores multiple windows of the renderer,
	// however each window is a single OpenGL Context,
	// which means the GL resources of each window is not shared.
	windows []ap.Window

	allWindowsClosed bool
	initialized      bool
}

// RendererInitParam is used for customize the parameters when init renderer.
type RendererInitParam struct {
	Name         string
	ViewDistance int
	VersionMajor int  // ContextVersionMajor
	VersionMinor int  // ContextVersionMinor
	Resizable    bool // All window resizable (cant modified after set)
	Visiable     bool // All window visiable (cant modified after set)
}

const (
	defaultRendererName = "Renderer"
	defaultViewDistance = 128
	defaultVersionMajor = 4
	defaultVersionMinor = 1
)

func InitAll() error {
	return glfw.Init()
}

func TerminateAll() {
	glfw.Terminate()
}

// Init initialize the renderer: setup OpenGL Context, compile built-in shaders.
func (r *RendererObj) Init(initParam interface{}) error {
	if r == nil {
		return fmt.Errorf("Init: %w", utils.ErrInvalidPointer)
	}

	if initParam == nil {
		initParam = RendererInitParam{}
	}
	p, ok := initParam.(RendererInitParam)
	if !ok {
		return fmt.Errorf("Init: %w", utils.ErrInvalidDataType)
	}

	if r.initialized {
		return fmt.Errorf("Init: %w", utils.ErrReInitialize)
	}

	// Handle parameters
	if p.Name == "" {
		p.Name = defaultRendererName
	}
	if p.ViewDistance <= 0 {
		p.ViewDistance = defaultViewDistance
	}
	if p.VersionMajor <= 0 {
		p.VersionMajor = defaultVersionMajor
	}
	if p.VersionMinor <= 0 {
		p.VersionMinor = defaultVersionMinor
	}

	// Setup OpenGL Context
	glfw.WindowHint(glfw.ContextVersionMajor, p.VersionMajor)
	glfw.WindowHint(glfw.ContextVersionMinor, p.VersionMinor)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	if p.Resizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}
	if p.Visiable {
		glfw.WindowHint(glfw.Visible, glfw.True)
	} else {
		glfw.WindowHint(glfw.Visible, glfw.False)
	}

	r.name = p.Name
	r.viewDistance = int32(p.ViewDistance)
	r.initialized = true

	return nil
}

func (r *RendererObj) AppendWindow(win ap.Window) error {
	if win == nil {
		return fmt.Errorf("AppendWindow: %w", utils.ErrInvalidParameter)
	}

	r.windows = append(r.windows, win)
	return nil
}

func (r *RendererObj) GetWindow(pos int) ap.Window {
	return r.windows[pos]
}

func (r *RendererObj) GetWindowNum() int {
	return len(r.windows)
}

func (r *RendererObj) GetViewDistance() int32 {
	return r.viewDistance
}

func (r *RendererObj) SetViewDistance(d int32) {
	if r == nil {
		return
	}
	r.viewDistance = d
}

// Render renders all windows in renderer by calling Window.Flush() method.
// If all windows are closed, this method will return.
func (r *RendererObj) Render() error {
	if len(r.windows) == 0 {
		return fmt.Errorf("Renderer [%s] not initialized", r.name)
	}

	for !r.allWindowsClosed {
		allWindowsClosed := true
		for _, win := range r.windows {
			// skip the closed window
			if win.IsClosed() {
				continue
			} else {
				allWindowsClosed = false
			}

			win.Flush()
		}

		r.allWindowsClosed = allWindowsClosed
	}
	return nil
}

// Release releases the resources of the Renderer
func (r RendererObj) Release() {
	// Destroy windows
	for _, win := range r.windows {
		win.Destroy()
	}
	r.windows = []ap.Window{}
	// glfw.Terminate()
}

func (r *RendererObj) GetName() string {
	return r.name
}

func (r *RendererObj) SetName(name string) {
	if r == nil {
		return
	}
	r.name = name
}

func NewRendererObj(p *RendererInitParam) (*RendererObj, error) {
	r := RendererObj{}
	err := r.Init(*p)
	if err != nil {
		return nil, fmt.Errorf("NewRendererObj: %w", err)
	}
	return &r, nil
}
