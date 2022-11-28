// Package aperture defines the render related interfaces.
package aperture

import "github.com/engoengine/glm"

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

	SetClearColor([4]float32)
	GetClearColor() [4]float32

	AppendShader(Shader)
	GetShader(int) Shader
	GetShaderNum() int

	AppendTexture(Texture)
	GetTexture(int) Texture
	GetTextureNum() int

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

// Shader interface defines the methods required by the shader
type Shader interface {
	// Load loads the shader program from file.
	//
	// vs is the file path of vertex shader;
	// fs is the file path of fragment shader;
	// gs is the file path of geometry shader.
	Load(vs, fs, gs string) error

	// Load loads the shader program from memory.
	//
	// vs is the string value of vertex shader;
	// fs is the string value of fragment shader;
	// gs is the string value of geometry shader.
	LoadMemory(vs, fs, gs string) error

	// Set method sets uniform data to shader program.
	//
	// This method supports following data types:
	// int, int32, uint, uint32, float32, vec2 ([2]float), vec3, vec4
	// mat2 ([4]float), mat3 ([9]float), mat4, mat2x3, mat3x4.
	Set(string, interface{}) error

	// GetID gets the shader program ID
	GetID() uint32
}

// Texture interface defines the methods required by a texture
type Texture interface {
	// Load method loads the image file and convert it to OpenGL texture ID
	Load(string) error

	// LoadMemory loads the image from memory data,
	// others same as the Load method.
	LoadMemory(width, height int, data *[]byte) error

	// GetID gets the ID of texture
	GetID() uint32

	// SetFileName sets the image file name of texture
	SetFileName(string)
	// GetFileName gets the image file name of texture
	GetFileName() string
}

// Camera interface defines the methods required by a Euler angle camera.
type Camera interface {
	// Init initialize the default camera.
	Init() error

	// GetViewMatrix calculate the current view matrix of the camera.
	GetViewMatrix() glm.Mat4

	// GetPosition gets the camera position in the world coordinates.
	GetPosition() glm.Vec3
	// SetPosition sets the camera position in the world coordinates.
	SetPosition(glm.Vec3)

	// GetZoom gets the zoom value of the camera.
	GetZoom() float32
	// SetZoom sets the zoom value of the camera.
	SetZoom(float32)

	// GetFront gets the front vector of the camera.
	//
	// If you want to set the front vector of the camera, please use
	// SetYaw and SetPitch method instead create a SetFront method.
	GetFront() glm.Vec3

	// SetUp sets the world up vector of the camera.
	SetUp(glm.Vec3)

	// SetYaw sets the yaw value of the Euler coordinates.
	SetYaw(float32)
	// SetPitch sets the pitch value of the Euler coordinates.
	SetPitch(float32)

	// SetSensitivity sets the sensitivity when moving speed of the camera.
	SetSensitivity(float32)
	// SetSpeed sets the moving speed of the camera.
	SetSpeed(float32)

	// ProcessMovement handles the movements of the camera,
	// the first param `dt` is the delta time of the frame,
	// the second param `dir` is the value of the moving direction,
	// third param `speedUp` is the speed-up value of the camera.
	ProcessMovement(dt float32, dir int, speedUp float32)

	// ProcessMouseMove handles the mouse/touch-screen movements,
	// the `xOffset` and `yOffset` are the offsets of the mouse.
	ProcessMouseMove(xOffset, yOffset float32, pitch bool)

	// ProcessScroll process the scroll of the mouse.
	ProcessScroll(yOffset float32)

	// GetID gets the ID of the camera.
	GetID() uint32
}

type Mesh interface {
	// TODO:
}

type Model interface {
	// TODO:
}

// Namer interface defines the methods required by a namer
type Namer interface {
	// GetName gets the name of this object
	GetName() string
	// SetName sets the name of this object
	SetName()
}
