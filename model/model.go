package model

import (
	"github.com/STARRY-S/aperture"
	"github.com/engoengine/glm"
)

type ModelObj struct {
	name string

	id uint32

	position glm.Vec3 // position of the model
	scale    glm.Vec3 // scale vector of the model

	rotateAngle float32
	rotateAxis  glm.Vec3

	textures []aperture.Texture
	meshes   []aperture.Mesh

	filename string
}

func (m *ModelObj) Init(filename string) error {
	m.position = [3]float32{0, 0, 0}
	m.scale = [3]float32{0, 0, 0}
	m.rotateAxis = [3]float32{0, 0, 0}
	m.rotateAngle = 0

	m.filename = filename

	return nil
}

func (m *ModelObj) Render() error {
	return nil
}

func (m *ModelObj) SetPos(p glm.Vec3) {
	m.position = p
}

func (m *ModelObj) SetScale(s glm.Vec3) {
	m.scale = s
}

func (m *ModelObj) SetRotate(axis glm.Vec3, angle float32) {
	m.rotateAxis = axis
	m.rotateAngle = angle
}

func (m *ModelObj) Release() {
	// for t := range m.textures {
	// 	// TODO: release texture
	// }
}
