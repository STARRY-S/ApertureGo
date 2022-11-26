// Package aperture defines the render related interfaces.
package aperture

// Renderer interface defines methods required by a renderer.
type Renderer interface {
	// Init method initializes the renderer.
	Init(interface{}) error

	// AppendWindow adds a window object to the renderer.
	AppendWindow(Window) error
	// GetWindow gets the nth window stored in renderer.
	GetWindow(int) Window
	// GetWindowNum gets total number of windows stored in renderer.
	GetWindowNum() int

	// Render method calls the main loop function of OpenGL,
	// it will render all visible windows by calling window.Flush method.
	Render() error

	// Release method releases the resource of the renderer.
	Release()
}

// RenderFunc is the function for customize rendering stuff in main loop,
// This function is implemented to be called in in Window.Flush() method.
type RenderFunc func()

// Window interface defines the methods required by a window,
// each window contains a different OpenGL Context, and the resources of each
// OpenGL Context is not shared.
type Window interface {
	// Init initializes the window.
	Init(interface{}) error

	// SetTitle sets the title of the window.
	SetTitle(string)
	// GetTitle gets the current title of the window.
	GetTitle() string

	// SetSize sets the window size: width, height.
	SetSize(int32, int32)
	// GetSize gets the window size: width, height.
	GetSize() (int32, int32)

	// SetVisible sets the window is visible or not.
	SetVisible(bool)
	// GetVisible gets the window is visible or not.
	GetVisible() bool

	// SetResizable sets the window is resizable or not.
	SetResizable(bool)
	// GetResizable gets the window is resizable or not.
	GetResizable() bool

	// GetFPS gets the current frame-per-second (FPS) value in float64 type.
	GetFPS() float64
	// GetFrameCount gets the total rendered frame count of this window.
	GetFrameCount() uint64

	// MakeContextCurrent makes the OpenGL Context of this window as the
	// current context in use. Call this method before do OpenGL render stuffs.
	MakeContextCurrent()

	SetRenderFunc(RenderFunc)
	GetRenderFunc() RenderFunc

	// Flush renders one frame of window (by calling RenderFunc function)
	// and updates the current status of window.
	Flush()

	// Close closes the window, by default the window should in open status,
	// you need to use this method to close the window, the window closed means
	// we finished all render stuffs of this window, and the window will not
	// re-open again.
	Close()
	// IsClosed gets the window is closed or not, renderer should use this
	// method to stop the main render loop function.
	IsClosed() bool

	// Destroy method is used to destroies all resources of this window.
	Destroy()
}

type Namer interface {
	GetName()
	SetName()
}
