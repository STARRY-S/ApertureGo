package shader_test

import (
	"testing"
	"time"

	"github.com/STARRY-S/aperture/renderer"
	"github.com/STARRY-S/aperture/shader"
	"github.com/STARRY-S/aperture/window"
	"github.com/engoengine/glm"
)

const (
	TestVertexShader   = "test/vertex.glsl"
	TestFragmentShader = "test/fragment.glsl"
	TestGeometryShader = "test/geometry.glsl"
)

func TestLoadAndSet(t *testing.T) {
	renderer.InitAll()

	// Initialize a invisible OpenGL Context
	r, err := renderer.NewRendererObj(&renderer.RendererInitParam{
		Name:      "TestRenderer",
		Resizable: false,
		Visiable:  false,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	w, err := window.NewWindowObj(&window.WindowInitParam{})
	err = r.AppendWindow(w)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Load shader program
	s := shader.ShaderObj{}
	err = s.Load(TestVertexShader, TestFragmentShader, TestGeometryShader)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.GetID() > 0 {
		t.Logf("shader program ID: %v", s.GetID())
	} else {
		t.Fatalf("shader ID is 0, test failed")
	}

	// Set mat4
	second := time.After(1 * time.Second)
	renderFunc := func() {
		select {
		case <-second:
			r.GetWindow(0).Close()
		default:
			// prevent block
		}
		matTest := glm.Ident4()
		err = s.Set("view", matTest)
		if err != nil {
			t.Errorf("error occured when render frame count %v", r.GetWindow(0).GetFrameCount())
			t.Errorf(err.Error())
		}
	}
	r.GetWindow(0).SetRenderFunc(renderFunc)
	r.Render()
	r.Release()

	renderer.TerminateAll()
}
