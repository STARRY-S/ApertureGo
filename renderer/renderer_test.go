package renderer_test

import (
	"math"
	"testing"
	"time"

	ap "github.com/STARRY-S/aperture"
	"github.com/STARRY-S/aperture/renderer"
	"github.com/STARRY-S/aperture/window"
)

func TestInterface(t *testing.T) {
	rObj := renderer.RendererObj{}
	var r ap.Renderer
	r = &rObj
	if r.GetWindowNum() != 0 {
		t.Errorf("GetWindowNum should be 0")
	}
}

func TestInit(t *testing.T) {
	renderer.InitAll()
	r, err := renderer.NewRendererObj(&renderer.RendererInitParam{
		Name:      "TestRenderer",
		Resizable: true,
		Visiable:  true,
	})
	if err != nil {
		t.Errorf(err.Error())
	}
	r.Release()
	renderer.TerminateAll()
}

// TestRendererWindow tests the built-in renderer and window
// struct and methods, this function will create one demo window in renderer
// and run 5 seconds.
func TestRendererWindow(t *testing.T) {
	renderer.InitAll()
	r, err := renderer.NewRendererObj(&renderer.RendererInitParam{
		Name:      "TestRenderer",
		Resizable: true,
		Visiable:  true,
	})
	if err != nil {
		t.Errorf(err.Error())
	}
	w, err := window.NewWindowObj(&window.WindowInitParam{
		BackgroundColor: [4]float32{0.2, 0.3, 0.3, 1.0},
	})
	r.AppendWindow(w)

	oneSec := time.After(1 * time.Second)
	counter := 0
	renderFunc := func() {
		win := r.GetWindow(0)
		select {
		// after one second
		case <-oneSec:
			counter++
			t.Logf("-------- second %v ---------", counter)
			t.Logf("fps: %.2f", win.GetFPS())
			t.Logf("frame count: %v", win.GetFrameCount())
			if math.IsInf(win.GetFPS(), 0) || math.IsNaN(win.GetFPS()) ||
				win.GetFPS() > 10e6 || win.GetFPS() < 1 {
				t.Errorf("invalid fps: %.2f", win.GetFPS())
			}
			if win.GetFrameCount() > 10e6 && win.GetFrameCount() < 1 {
				t.Errorf("invalid frame count: %v", win.GetFrameCount())
			}
			if counter == 1 {
				win.Close()
				return
			}
			oneSec = time.After(1 * time.Second)
		default:
			// do nothing, prevent block when trying to read from channel
		}
	}
	win := r.GetWindow(0)
	win.SetRenderFunc(renderFunc)
	// test set title
	win.SetTitle("Renderer & Window Test")
	if title := win.GetTitle(); title != "Renderer & Window Test" {
		t.Errorf("SetTitle failed, got %v", title)
	}
	// test resizable
	win.SetResizable(false)
	if b := win.GetResizable(); b != false {
		t.Errorf("SetResizable failed")
	}
	// test SetSize
	win.SetSize(1280, 720)
	if w, h := win.GetSize(); w != 1280 || h != 720 {
		t.Errorf("SetSize failed")
	}

	// test main render loop
	err = r.Render()
	if err != nil {
		t.Errorf(err.Error())
	}

	r.Release()
	renderer.TerminateAll()
}

// TestRendererMultiWindow tests the built-in renderer and window
// struct and methods, this function will create multi demo windows in renderer
// and run 5 seconds.
func TestRendererMultiWindow(t *testing.T) {
	t.Log("Running multi-windows test")
	renderer.InitAll()
	r, err := renderer.NewRendererObj(&renderer.RendererInitParam{
		Name:      "TestRenderer",
		Resizable: true,
		Visiable:  true,
	})
	if err != nil {
		t.Errorf(err.Error())
	}
	w, err := window.NewWindowObj(&window.WindowInitParam{
		Width:           400,
		Height:          400,
		PosX:            300,
		PosY:            300,
		BackgroundColor: [4]float32{0.2, 0.3, 0.3, 1.0},
	})
	r.AppendWindow(w)

	w, err = window.NewWindowObj(&window.WindowInitParam{
		Width:           400,
		Height:          400,
		PosX:            800,
		PosY:            300,
		BackgroundColor: [4]float32{0.4, 0.3, 0.3, 1.0},
	})
	r.AppendWindow(w)

	w, err = window.NewWindowObj(&window.WindowInitParam{
		Width:           400,
		Height:          400,
		PosX:            1300,
		PosY:            300,
		BackgroundColor: [4]float32{0.6, 0.4, 0.3, 1.0},
	})
	r.AppendWindow(w)

	oneSec1 := time.After(1 * time.Second)
	counter1 := 0
	renderFunc1 := func() {
		win1 := r.GetWindow(0)

		select {
		// after one second
		case <-oneSec1:
			counter1++
			t.Logf("-------- Window 1 second %v ---------", counter1)
			t.Logf("\tfps: %.2f", win1.GetFPS())
			t.Logf("\tframe count: %v", win1.GetFrameCount())
			if counter1 == 3 {
				win1.Close()
				return
			}
			oneSec1 = time.After(1 * time.Second)
		default:
			// do nothing, prevent block when trying to read from channel
		}
	}

	oneSec2 := time.After(1 * time.Second)
	counter2 := 0
	renderFunc2 := func() {
		win2 := r.GetWindow(1)

		select {
		// after one second
		case <-oneSec2:
			counter2++
			t.Logf("-------- Window 2: second %v ---------", counter2)
			t.Logf("\tfps: %.2f", win2.GetFPS())
			t.Logf("\tframe count: %v", win2.GetFrameCount())
			if counter2 == 2 {
				win2.Close()
				return
			}
			oneSec2 = time.After(1 * time.Second)
		default:
			// do nothing, prevent block when trying to read from channel
		}
	}

	oneSec3 := time.After(1 * time.Second)
	counter3 := 0
	renderFunc3 := func() {
		win3 := r.GetWindow(2)

		select {
		// after one second
		case <-oneSec3:
			counter3++
			t.Logf("-------- Window 2: second %v ---------", counter3)
			t.Logf("\tfps: %.2f", win3.GetFPS())
			t.Logf("\tframe count: %v", win3.GetFrameCount())
			if counter3 == 1 {
				win3.Close()
				return
			}
			oneSec3 = time.After(1 * time.Second)
		default:
			// do nothing, prevent block when trying to read from channel
		}
	}
	r.GetWindow(0).SetRenderFunc(renderFunc1)
	r.GetWindow(0).SetTitle("Window 1")

	r.GetWindow(1).SetRenderFunc(renderFunc2)
	r.GetWindow(1).SetTitle("Window 2")

	r.GetWindow(2).SetRenderFunc(renderFunc3)
	r.GetWindow(2).SetTitle("Window 3")

	// test main render loop
	err = r.Render()
	if err != nil {
		t.Errorf(err.Error())
	}

	r.Release()
	renderer.TerminateAll()
}
