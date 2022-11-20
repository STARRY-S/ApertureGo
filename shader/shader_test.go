package shader_test

import (
	"testing"

	"github.com/STARRY-S/aperture/renderer"
	"github.com/STARRY-S/aperture/shader"
	"github.com/engoengine/glm"
)

const (
	TestVertexShader   = "test/vertex.glsl"
	TestFragmentShader = "test/fragment.glsl"
	TestGeometryShader = "test/geometry.glsl"
)

func TestLoadAndSet(t *testing.T) {
	// Initialize a invisible OpenGL Context
	r := renderer.RendererObj{}
	p := renderer.RendererInitParam{
		Name:       "TestRenderer",
		Width:      1,
		Height:     1,
		Resizable:  false,
		Visiable:   false,
		FullScreen: false,
	}
	r.Init(p)

	// Load shader program
	s := shader.ShaderObj{}
	err := s.Load(TestVertexShader, TestFragmentShader, TestGeometryShader)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.GetID() > 0 {
		t.Logf("shader program: %v", s.GetID())
	} else {
		t.Errorf("shader ID is 0, test failed")
	}

	// Set mat4
	start := r.GetWindow(0).GetCFT()
	renderFunc := func() {
		matTest := glm.Ident4()
		err = s.Set("view", matTest)
		if err != nil {
			t.Errorf("error occured when render frame count %v", r.GetWindow(0).GetFrameCount())
			t.Errorf(err.Error())
		}
		if r.GetWindow(0).GetCFT()-start > 1 {
			r.GetWindow(0).ShouldClose(true)
		}
	}
	r.GetWindow(0).SetRenderFunc(renderFunc)
	r.Render()
	r.Release()
}
