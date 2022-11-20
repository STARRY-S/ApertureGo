package renderer

import (
	"fmt"

	"github.com/STARRY-S/aperture/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/sirupsen/logrus"
)

// RenderFunc is the function for customize rendering stuff in main loop
type RenderFunc func()

// The renderer struct object
type RendererObj struct {
	name         string
	viewDistance int32
	windows      []WindowObj // Each window is a different OpenGL Context
}

// The parameters for init renderer
type RendererInitParam struct {
	Name         string // Renderer name, reserved
	Width        int32  // First generated window width
	Height       int32  // First generated window Height
	ViewDistance int32  // View distance
	VersionMajor int    // ContextVersionMajor
	VersionMinor int    // ContextVersionMinor
	Resizable    bool   // All window resizable (cant modified after set)
	Visiable     bool   // All window visiable (cant modified after set)
	FullScreen   bool   // Fullscreen mode
}

// The window struct object
type WindowObj struct {
	// times (in seconds)
	cft float64
	lft float64
	dt  float64

	// total frame count of this window
	frameCount uint64

	width      int32        // buffer width
	height     int32        // buffer height
	title      string       // title
	renderFunc RenderFunc   // render function
	glfwWindow *glfw.Window // GLFW window
}

const (
	defaultWidth        = 720        // Default window width
	defaultHeight       = 480        // Default window height
	defaultTitle        = "Window"   // Default window title
	defaultName         = "Renderer" // Default name of renderer
	defaultViewDistance = 128        // Default view distance (16 * 8)
)

// Init initialize the renderer: setup OpenGL Context, compile built-in shaders
func (r *RendererObj) Init(p RendererInitParam) error {
	if r == nil {
		return utils.ErrInvalidPointer
	}

	if len(r.windows) > 0 {
		return utils.ErrReInitialize
	}

	// Set Renderer Name
	if p.Name != "" {
		r.name = p.Name
	} else {
		r.name = defaultName
	}

	// Set View Distance
	if p.ViewDistance > 0 {
		r.viewDistance = p.ViewDistance
	} else {
		r.viewDistance = defaultViewDistance
	}

	// Setup OpenGL Context
	err := glfw.Init()
	if err != nil {
		return err
	}
	if p.VersionMajor > 0 && p.VersionMinor > 0 {
		glfw.WindowHint(glfw.ContextVersionMajor, p.VersionMajor)
		glfw.WindowHint(glfw.ContextVersionMinor, p.VersionMinor)
	} else {
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
	}
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	if p.Visiable {
		glfw.WindowHint(glfw.Visible, glfw.True)
	} else {
		glfw.WindowHint(glfw.Visible, glfw.False)
	}
	if p.Resizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}

	if p.Width > 0 && p.Height > 0 {
		err = r.generateWindow(p.Width, p.Height, defaultTitle, p.FullScreen)
	} else {
		err = r.generateWindow(defaultWidth, defaultHeight, defaultTitle, p.FullScreen)
	}
	if err != nil {
		return err
	}
	r.GetWindow(0).MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return err
	}
	// TODO: Compile built-in shaders

	return nil
}

func (r RendererObj) GetName() string {
	return r.name
}

func (r *RendererObj) SetName(name string) {
	if r == nil {
		return
	}
	r.name = name
}

func (r RendererObj) GetViewDistance() int32 {
	return r.viewDistance
}

func (r *RendererObj) SetViewDistance(d int32) {
	if r == nil {
		return
	}
	r.viewDistance = d
}

// GetWindow gets the offset window of the renderer
func (r RendererObj) GetWindow(offset int32) *WindowObj {
	return &r.windows[offset]
}

// generateWindow create a new window for renderer
func (r *RendererObj) generateWindow(width, height int32, title string, fullscreen bool) error {
	var glfwWindow *glfw.Window
	var err error
	win := WindowObj{
		width:      width,
		height:     height,
		renderFunc: defaultRenderFunc,
	}

	// make the window position to the center of the monitor
	monitors := glfw.GetMonitors()
	if len(monitors) != 0 {
		mWidth := monitors[0].GetVideoMode().Width
		mHeight := monitors[0].GetVideoMode().Height
		if fullscreen {
			glfwWindow, err = glfw.CreateWindow(
				int(mWidth),
				int(mHeight),
				title,
				monitors[0],
				nil)
		} else {
			glfwWindow, err = glfw.CreateWindow(
				int(width),
				int(height),
				title,
				nil,
				nil)
			glfwWindow.SetPos(mWidth/2-int(width)/2, mHeight/2-int(height)/2)
		}
	} else {
		logrus.Warnln("failed to get monitors")
		if fullscreen {
			logrus.Errorln("unable to set fullscreen mode")
		}
		glfwWindow, err = glfw.CreateWindow(
			int(width),
			int(height),
			title,
			nil,
			nil)
	}

	if err != nil {
		return err
	}

	win.glfwWindow = glfwWindow
	r.windows = append(r.windows, win)
	return nil
}

// Render calls the main loop function of renderer
func (r RendererObj) Render() error {
	if len(r.windows) == 0 {
		return fmt.Errorf("Renderer [%s] not initialized", r.name)
	}
	win := r.GetWindow(0)
	if win.renderFunc == nil {
		logrus.Errorln("Render function not set, set to defaultRenderFunc")
		win.renderFunc = defaultRenderFunc
	}

	for !win.glfwWindow.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(1.0, 1.0, 1.0, 1.0)

		// flush window, render one frame
		win.flush()

		win.glfwWindow.SwapBuffers()
		glfw.PollEvents()
	}
	return nil
}

// Release releases the resources of the Renderer
func (r RendererObj) Release() {
	// Destroy windows
	for _, win := range r.windows {
		win.Destroy()
	}
	r.windows = []WindowObj{}
	glfw.Terminate()
}

func (w WindowObj) MakeContextCurrent() {
	w.glfwWindow.MakeContextCurrent()
}

func (w *WindowObj) SetTitle(s string) {
	w.title = s
	w.glfwWindow.SetTitle(s)
}

func (w WindowObj) GetTitle() string {
	return w.title
}

// GetSize gets the width and height of the window
func (w WindowObj) GetSize() (width, height int32) {
	return w.width, w.height
}

// SetSize sets the width and height of the window
func (w *WindowObj) SetSize(width, height int32) {
	w.width = width
	w.height = height
	w.glfwWindow.SetSize(int(width), int(height))
}

// ShouldClose sets whether the window should be closed
func (w *WindowObj) ShouldClose(b bool) {
	w.glfwWindow.SetShouldClose(b)
}

// SetRenderFunc sets the render function in main loop
func (w *WindowObj) SetRenderFunc(f RenderFunc) {
	w.renderFunc = f
}

// GetFrameCount gets the total frame count of this window
func (w WindowObj) GetFrameCount() uint64 {
	return w.frameCount
}

// Set the window is visiable or not (show or hide)
func (w WindowObj) SetVisible(b bool) {
	if b {
		w.glfwWindow.Show()
	} else {
		w.glfwWindow.Hide()
	}
}

func (w WindowObj) GetVisiable() bool {
	i := w.glfwWindow.GetAttrib(glfw.Visible)
	if i == glfw.True {
		return true
	}
	return false
}

func (w WindowObj) SetResizable(b bool) {
	if b {
		w.glfwWindow.SetAttrib(glfw.Resizable, glfw.True)
	} else {
		w.glfwWindow.SetAttrib(glfw.Resizable, glfw.False)
	}
}

func (w WindowObj) GetResizable() bool {
	i := w.glfwWindow.GetAttrib(glfw.Resizable)
	if i == glfw.True {
		return true
	}
	return false
}

func (w WindowObj) GetFPS() float64 {
	return float64(1.0 / w.dt)
}

func (w WindowObj) GetCFT() float64 {
	return w.cft
}

func (w WindowObj) GetDT() float64 {
	return w.dt
}

func (w *WindowObj) flush() {
	w.frameCount++
	w.cft = glfw.GetTime()
	w.dt = w.cft - w.lft
	w.lft = w.cft
	// main render function
	w.renderFunc()
}

// Destroy destroies the GLFW window
func (w *WindowObj) Destroy() {
	w.glfwWindow.Destroy()
	w.glfwWindow = nil
}

// defaultRenderFunc is the default Render function of renderer
func defaultRenderFunc() {
	// do nothing
	return
}
