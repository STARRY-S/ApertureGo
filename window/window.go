// Package renderer has the built-in implemented window struct.
package window

import (
	"fmt"

	ap "github.com/STARRY-S/aperture"
	"github.com/STARRY-S/aperture/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/sirupsen/logrus"
)

// WindowObj implements the Window interface.
type WindowObj struct {
	name string

	// cft is the time of current frame (in second)
	cft float64
	// lft is the time of last frame (in second)
	lft float64
	// dt is the duration from last frame to current frame
	dt float64
	// frameCount is the total frame count of this window
	frameCount uint64

	width           int32
	height          int32
	title           string
	renderFunc      ap.RenderFunc
	backgroundColor [4]float32

	shaders  []ap.Shader
	textures []ap.Texture

	glfwWindow *glfw.Window

	rendering   bool
	initialized bool
}

// RendererInitParam is used for customize the parameters when init window.
type WindowInitParam struct {
	Name           string
	Width          int
	Height         int
	PosX           int
	PosY           int
	Title          string
	Func           ap.RenderFunc
	FullScreenMode bool

	// BackgroundColor is the RGBA value of the parameter of glClearColor func.
	BackgroundColor [4]float32
}

const (
	defaultWindowWidth  = 720
	defaultWindowHeight = 480
	defaultWindowTitle  = "Window"
)

func (w *WindowObj) Init(initParam interface{}) error {
	if w == nil {
		return fmt.Errorf("Init: %w", utils.ErrInvalidPointer)
	}

	if initParam == nil {
		initParam = WindowInitParam{}
	}
	p, ok := initParam.(WindowInitParam)
	if !ok {
		return fmt.Errorf("Init: %w", utils.ErrInvalidDataType)
	}

	if w.initialized {
		return fmt.Errorf("Init: %w", utils.ErrReInitialize)
	}

	var err error

	// Handle parameters
	if p.Width <= 0 {
		p.Width = defaultWindowWidth
	}
	if p.Height <= 0 {
		p.Height = defaultWindowHeight
	}
	if p.Title == "" {
		p.Title = defaultWindowTitle
	}
	if p.Func == nil {
		p.Func = defaultRenderFunc
	}

	// Generate GLFW Window.
	monitors := glfw.GetMonitors()
	if len(monitors) != 0 {
		mWidth := monitors[0].GetVideoMode().Width
		mHeight := monitors[0].GetVideoMode().Height

		// calculate window position
		if p.PosX == 0 {
			p.PosX = mWidth/2 - p.Width/2
		}
		if p.PosY == 0 {
			p.PosY = mHeight/2 - p.Height/2
		}

		if p.FullScreenMode {
			w.glfwWindow, err = glfw.CreateWindow(
				mWidth,
				mHeight,
				p.Title,
				monitors[0],
				nil)
			if err != nil {
				return fmt.Errorf("Init: %w", err)
			}
		} else {
			w.glfwWindow, err = glfw.CreateWindow(
				p.Width,
				p.Height,
				p.Title,
				nil,
				nil)
			if err != nil {
				return fmt.Errorf("Init: %w", err)
			}
			// Set window pos to the center of the monitor.
			w.glfwWindow.SetPos(p.PosX, p.PosY)
		}
	} else {
		logrus.Warnln("Init: failed to get monitors")
		if p.FullScreenMode {
			logrus.Errorln("Init: unable to set fullscreen mode")
		}
		w.glfwWindow, err = glfw.CreateWindow(
			p.Width,
			p.Height,
			p.Title,
			nil,
			nil)
	}

	if err != nil {
		return err
	}

	if w.glfwWindow == nil {
		return fmt.Errorf("failed to generate GLFW window")
	}

	w.width = int32(p.Width)
	w.height = int32(p.Height)
	w.title = p.Title
	w.renderFunc = p.Func
	w.rendering = true
	w.backgroundColor = p.BackgroundColor

	w.glfwWindow.MakeContextCurrent()
	// Call gl.Init only under the presence of an active OpenGL context,
	// i.e., after MakeContextCurrent.
	if err := gl.Init(); err != nil {
		return fmt.Errorf("Init: %w", err)
	}

	w.initialized = true

	return nil
}

func (w *WindowObj) SetTitle(s string) {
	w.title = s
	w.glfwWindow.SetTitle(s)
}

func (w *WindowObj) GetTitle() string {
	return w.title
}

// SetSize sets the width and height of the window.
func (w *WindowObj) SetSize(width, height int32) {
	w.width = width
	w.height = height
	w.glfwWindow.SetSize(int(width), int(height))
}

// GetSize gets the width and height of the window.
func (w *WindowObj) GetSize() (width, height int32) {
	return w.width, w.height
}

// Set the window is visiable or not (show or hide).
func (w *WindowObj) SetVisible(b bool) {
	if b {
		w.glfwWindow.Show()
	} else {
		w.glfwWindow.Hide()
	}
	w.rendering = b
}

func (w *WindowObj) GetVisible() bool {
	return w.glfwWindow.GetAttrib(glfw.Visible) == glfw.True
}

func (w *WindowObj) SetResizable(b bool) {
	if b {
		w.glfwWindow.SetAttrib(glfw.Resizable, glfw.True)
	} else {
		w.glfwWindow.SetAttrib(glfw.Resizable, glfw.False)
	}
}

func (w *WindowObj) GetResizable() bool {
	return w.glfwWindow.GetAttrib(glfw.Resizable) == glfw.True
}

func (w *WindowObj) GetFPS() float64 {
	return float64(1.0 / w.dt)
}

// func (w *WindowObj) GetCFT() float64 {
// 	return w.cft
// }

// func (w *WindowObj) GetDT() float64 {
// 	return w.dt
// }

// GetFrameCount gets the total frame count of this window.
func (w *WindowObj) GetFrameCount() uint64 {
	return w.frameCount
}

func (w *WindowObj) MakeContextCurrent() {
	w.glfwWindow.MakeContextCurrent()
}

// SetRenderFunc sets the render function in main loop.
func (w *WindowObj) SetRenderFunc(f ap.RenderFunc) {
	w.renderFunc = f
}

func (w *WindowObj) GetRenderFunc() ap.RenderFunc {
	return w.renderFunc
}

func (w *WindowObj) SetClearColor(rgba [4]float32) {
	w.backgroundColor = rgba
}

func (w *WindowObj) GetClearColor() [4]float32 {
	return w.backgroundColor
}

func (w *WindowObj) AppendShader(s ap.Shader) {
	if s == nil {
		return
	}
	w.shaders = append(w.shaders, s)
}

func (w *WindowObj) GetShader(pos int) ap.Shader {
	return w.shaders[pos]
}

func (w *WindowObj) GetShaderNum() int {
	return len(w.shaders)
}

func (w *WindowObj) AppendTexture(tex ap.Texture) {
	if tex == nil {
		return
	}
	w.textures = append(w.textures, tex)
}

func (w *WindowObj) GetTexture(pos int) ap.Texture {
	return w.textures[pos]
}

func (w *WindowObj) GetTextureNum() int {
	return len(w.textures)
}

func (w *WindowObj) Flush() {
	// update fps
	w.frameCount++
	w.cft = glfw.GetTime()
	w.dt = w.cft - w.lft
	w.lft = w.cft

	w.glfwWindow.MakeContextCurrent()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(
		w.backgroundColor[0],
		w.backgroundColor[1],
		w.backgroundColor[2],
		w.backgroundColor[3])

	// main render function
	if w.renderFunc == nil {
		logrus.Warnln("Flush: render function is nil, set back to default.")
		w.renderFunc = defaultRenderFunc
	}
	w.renderFunc()

	// V-Sync
	w.glfwWindow.SwapBuffers()
	glfw.PollEvents()
}

func (w *WindowObj) Close() {
	w.glfwWindow.SetShouldClose(true)
	w.SetVisible(false)
	w.rendering = false
}

func (w *WindowObj) IsClosed() bool {
	return w.glfwWindow.ShouldClose()
}

func (w *WindowObj) Destroy() {
	w.glfwWindow.Destroy()
	w.glfwWindow = nil
	w.initialized = false
}

func (w *WindowObj) GetName() string {
	return w.name
}

func (w *WindowObj) SetName(name string) {
	w.name = name
}

// defaultRenderFunc is the default Render function of the window,
// this function will do nothing.
func defaultRenderFunc() {
	ctx := glfw.GetCurrentContext()
	if ctx != nil {
		ctx.SwapBuffers()
	}
	glfw.PollEvents()
	return
}

func NewWindowObj(p *WindowInitParam) (*WindowObj, error) {
	win := WindowObj{}
	err := win.Init(*p)
	if err != nil {
		return nil, err
	}
	return &win, nil
}
