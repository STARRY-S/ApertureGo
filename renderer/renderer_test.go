package renderer_test

import (
	"math"
	"testing"

	"github.com/STARRY-S/aperture/renderer"
)

func TestInit(t *testing.T) {
	p := renderer.RendererInitParam{
		Name:       "TestRenderer",
		Resizable:  true,
		Visiable:   true,
		FullScreen: false,
	}
	r := renderer.RendererObj{}
	err := r.Init(p)
	if err != nil {
		t.Errorf(err.Error())
	}
	r.Release()
}

func TestRenderer(t *testing.T) {
	r := renderer.RendererObj{}
	p := renderer.RendererInitParam{
		Name:       "TestRenderer",
		Width:      800,
		Height:     400,
		Resizable:  true,
		Visiable:   true,
		FullScreen: false,
	}
	err := r.Init(p)
	if err != nil {
		t.Errorf(err.Error())
	}

	// renderFunc is the function to check the timer in window and fps
	start := r.GetWindow(0).GetCFT()
	counter := 0
	renderFunc := func() {
		win := r.GetWindow(0)
		// Stop rendering after 1 second
		if r.GetWindow(0).GetCFT()-start > 1.0 {
			counter++
			t.Logf("-------- second %v ---------", counter)
			t.Logf("fps: %v", win.GetFPS())
			t.Logf("frame count: %v", win.GetFrameCount())
			if math.IsInf(win.GetFPS(), 0) || math.IsNaN(win.GetFPS()) || win.GetFPS() > 9999 || win.GetFPS() < 1 {
				t.Errorf("invalid fps")
			}
			if win.GetFrameCount() > 9999 && win.GetFrameCount() < 1 {
				t.Errorf("invalid frame count: %v", win.GetFrameCount())
			}
			start = r.GetWindow(0).GetCFT()
		}
		// render 5 seconds
		if counter >= 5 {
			win.ShouldClose(true)
		}
	}
	win := r.GetWindow(0)
	win.SetRenderFunc(renderFunc)
	// test set title
	win.SetTitle("Test")
	if title := win.GetTitle(); title != "Test" {
		t.Errorf("SetTitle failed, expect 'Test', got %v", title)
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
}
