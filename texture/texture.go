package texture

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/STARRY-S/aperture/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type TextureObj struct {
	fileName string
	rgba     [4]float32
	id       uint32
}

// Load texture image from file
func (t *TextureObj) Load(f string) error {
	rgba, err := loadImage(f)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	id := newTextureRGBA(rgba)
	if id == 0 {
		return fmt.Errorf("failed to load texture from file %v", f)
	}
	t.id = id
	t.fileName = f
	return nil
}

// Load texture image from memory
func (t *TextureObj) LoadMemory(width, height int, data *[]byte) error {
	if width <= 0 || height <= 0 {
		return utils.ErrInvalidParameter
	}
	id := newTextureMemory(width, height, data)
	if id != 0 {
		return fmt.Errorf("failed to load texture from memory")
	}
	return nil
}

func (t TextureObj) GetID() uint32 {
	return t.id
}

func (t *TextureObj) SetFileName(name string) {
	t.fileName = name
}

func (t *TextureObj) GetFileName() string {
	return t.fileName
}

func loadImage(file string) (*image.RGBA, error) {
	imgFile, err := os.Open(file)
	defer imgFile.Close()
	if err != nil {
		return nil, fmt.Errorf("image %q not found on disk:\n%v", file, err)
	}
	imgFile.Seek(0, 0)
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba, nil
}

func newTextureRGBA(rgba *image.RGBA) uint32 {
	return newTextureMemory(rgba.Rect.Size().X, rgba.Rect.Size().X, &rgba.Pix)
}

func newTextureMemory(width, height int, data *[]byte) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(width),
		int32(height),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(*data))

	return texture
}
