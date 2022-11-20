package texture_test

import (
	"testing"

	"github.com/STARRY-S/aperture/renderer"
	"github.com/STARRY-S/aperture/texture"
)

const (
	TestTextureJPG = "test/test.jpg"
	TestTexturePNG = "test/test.jpg"
)

func TestLoadJPG(t *testing.T) {
	r := renderer.RendererObj{}
	// invisible opengl context
	p := renderer.RendererInitParam{
		Name:       "TestRenderer",
		Width:      1,
		Height:     1,
		Resizable:  false,
		Visiable:   false,
		FullScreen: false,
	}
	err := r.Init(p)
	if err != nil {
		t.Errorf("failed to init renderer")
		return
	}

	// set window invisible
	r.GetWindow(0).SetVisible(false)

	tex := texture.TextureObj{}
	err = tex.Load(TestTextureJPG)
	if err != nil {
		t.Errorf("failed to load texture %s:", TestTextureJPG)
		t.Errorf(err.Error())
		return
	}

	if tex.GetID() == 0 {
		t.Errorf("texture id is 0")
		return
	} else {
		t.Logf("load texture id: %v", tex.GetID())
	}

	if tex.GetFileName() != TestTextureJPG {
		t.Errorf("texture file name: %s, not equal to %s", tex.GetFileName(), TestTextureJPG)
	}

	r.Release()
}

func TestLoadPNG(t *testing.T) {
	r := renderer.RendererObj{}
	// invisible opengl context
	p := renderer.RendererInitParam{
		Name:       "TestRenderer",
		Width:      0,
		Height:     0,
		Resizable:  false,
		Visiable:   false,
		FullScreen: false,
	}
	err := r.Init(p)
	if err != nil {
		t.Errorf("failed to init renderer")
		return
	}

	// set window invisible
	r.GetWindow(0).SetVisible(false)

	tex := texture.TextureObj{}
	err = tex.Load(TestTexturePNG)
	if err != nil {
		t.Errorf("failed to load texture %s:", TestTexturePNG)
		t.Errorf(err.Error())
		return
	}

	if tex.GetID() == 0 {
		t.Errorf("texture id is 0")
		return
	} else {
		t.Logf("load texture id: %v", tex.GetID())
	}

	if tex.GetFileName() != TestTexturePNG {
		t.Errorf("texture file name: %s, not equal to %s", tex.GetFileName(), TestTexturePNG)
	}

	r.Release()
}
