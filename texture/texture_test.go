package texture_test

import (
	"testing"

	"github.com/STARRY-S/aperture"
	"github.com/STARRY-S/aperture/renderer"
	"github.com/STARRY-S/aperture/texture"
	"github.com/STARRY-S/aperture/window"
)

const (
	TestTextureJPG = "test/test.jpg"
	TestTexturePNG = "test/test.jpg"
)

func TestInterface(t *testing.T) {
	tObj := texture.TextureObj{}
	var tex aperture.Texture
	tex = &tObj
	if tex.GetID() != 0 {
		t.Errorf("GetID should be 0")
	}
}

func TestLoadJPG(t *testing.T) {
	renderer.InitAll()
	r, err := renderer.NewRendererObj(&renderer.RendererInitParam{
		Name:      "TestRenderer",
		Resizable: false,
		Visiable:  false,
	})
	if err != nil {
		t.Fatalf("failed to init renderer")
	}
	win, err := window.NewWindowObj(&window.WindowInitParam{})
	r.AppendWindow(win)

	tex := texture.TextureObj{}
	err = tex.Load(TestTextureJPG)
	if err != nil {
		t.Errorf("failed to load texture %s:", TestTextureJPG)
		t.Errorf(err.Error())
		return
	}

	if tex.GetID() == 0 {
		t.Errorf("texture id is 0")
	} else {
		t.Logf("load texture id: %v", tex.GetID())
	}

	if tex.GetFileName() != TestTextureJPG {
		t.Errorf("texture file name: %s, not equal to %s",
			tex.GetFileName(), TestTextureJPG)
	}

	r.Release()
	renderer.TerminateAll()
}

func TestLoadPNG(t *testing.T) {
	renderer.InitAll()
	r, err := renderer.NewRendererObj(&renderer.RendererInitParam{
		Name:      "TestRenderer",
		Resizable: false,
		Visiable:  false,
	})
	if err != nil {
		t.Errorf("failed to init renderer")
		return
	}
	win, err := window.NewWindowObj(&window.WindowInitParam{})
	r.AppendWindow(win)

	tex := texture.TextureObj{}
	err = tex.Load(TestTexturePNG)
	if err != nil {
		t.Errorf("failed to load texture %s:", TestTexturePNG)
		t.Errorf(err.Error())
		return
	}

	if tex.GetID() == 0 {
		t.Fatalf("texture id is 0")
	} else {
		t.Logf("load texture id: %v", tex.GetID())
	}

	if tex.GetFileName() != TestTexturePNG {
		t.Errorf("texture file name: %s, not equal to %s",
			tex.GetFileName(), TestTexturePNG)
	}

	r.Release()
	renderer.TerminateAll()
}
