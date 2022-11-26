package window

import (
	"testing"

	"github.com/STARRY-S/aperture"
	"github.com/STARRY-S/aperture/renderer"
	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	wObj := WindowObj{}
	var w aperture.Window
	w = &wObj
	if w.GetFrameCount() != 0 {
		t.Errorf("GetFrameCount should be 0")
	}
}

func TestNewDestroy(t *testing.T) {
	// create a window means create a OpenGL Context, so we need to initialize
	// renderer (Initialize GLFW) before create window.
	renderer.InitAll()
	win, err := NewWindowObj(&WindowInitParam{})
	if err != nil {
		t.Errorf(err.Error())
	}
	if win.glfwWindow == nil {
		t.Errorf("glfwWindow is nil")
	}
	win.Destroy()
	if win.glfwWindow != nil {
		t.Errorf("destroy failed")
	}
	renderer.TerminateAll()
}

func TestSetterGetter(t *testing.T) {
	renderer.InitAll()
	win, err := NewWindowObj(&WindowInitParam{})
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	win.SetName("AAAA")
	assert.Equal(t, "AAAA", win.GetName(), "name should be 'AAAA'")
	win.SetTitle("BBB")
	assert.Equal(t, "BBB", win.GetTitle(), "title should be 'BBB'")
	win.SetSize(100, 200)
	width, height := win.GetSize()
	assert.Equal(t, int32(100), width, "width should be int32(100)")
	assert.Equal(t, int32(200), height, "height should be int32(200)")
	win.SetVisible(true)
	assert.Equal(t, true, win.GetVisible(), "visible should be true")
	win.SetVisible(false)
	assert.Equal(t, false, win.GetVisible(), "visible should be false")
	win.SetResizable(false)
	assert.Equal(t, false, win.GetResizable(), "resizable should be false")
	win.SetResizable(true)
	assert.Equal(t, true, win.GetResizable(), "resizable should be true")
	win.SetRenderFunc(nil)
	if win.GetRenderFunc() != nil {
		t.Errorf("SetRenderFunc failed, render func should be nil")
	}
	win.SetClearColor([4]float32{0.1, 0.2, 0.3, 1.0})
	assert.Equal(t, [4]float32{0.1, 0.2, 0.3, 1.0}, win.GetClearColor(),
		"clear color should be [ 0.1, 0.2, 0.3, 1.0 ]")

	win.Close()
	assert.Equal(t, true, win.IsClosed(), "window should be closed")

	win.Destroy()
	renderer.TerminateAll()
}
