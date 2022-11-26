package camera

import (
	"math"

	"github.com/engoengine/glm"
)

// CameraObj has the data required by a Euler angle Camera
type CameraObj struct {
	id   uint32
	name string

	position glm.Vec3
	front    glm.Vec3
	up       glm.Vec3
	right    glm.Vec3
	worldUp  glm.Vec3

	yaw         float32
	pitch       float32
	speed       float32
	sensitivity float32
	zoom        float32
}

const (
	DirectionForward = iota + 1000
	DirectionBackwoard
	DirectionLeft
	DirectionRight
	DirectionUp
	DirectionDown
)

var (
	defaultCameraPosition = glm.Vec3{0.0, 0.0, 0.0}
	defaultCameraRight    = glm.Vec3{0.0, 0.0, 0.0}
	defaultCameraUp       = glm.Vec3{0.0, 1.0, 0.0}

	// default set camera direction to (1,0,0) (yaw = 0, pitch = 0)
	defaultCameraFront = glm.Vec3{1.0, 0.0, 0.0}

	defaultCameraYaw         = 0.0
	defaultCameraPitch       = 0.0
	defaultCameraSpeed       = 1.0
	defaultCameraSensitivity = 0.04
	defaultCameraZoom        = 65.0
)

func (c *CameraObj) Init() error {
	c.position = defaultCameraPosition
	c.right = defaultCameraRight
	c.up = defaultCameraUp
	c.front = defaultCameraFront
	c.yaw = float32(defaultCameraYaw)
	c.pitch = float32(defaultCameraPitch)
	c.speed = float32(defaultCameraSpeed)
	c.sensitivity = float32(defaultCameraSensitivity)
	c.zoom = float32(defaultCameraZoom)

	return nil
}

func (c *CameraObj) GetViewMatrix() glm.Mat4 {
	target := c.position.Add(&c.front)
	return glm.LookAtV(&c.position, &target, &c.up)
}

func (c *CameraObj) GetPosition() glm.Vec3 {
	return c.position
}

func (c *CameraObj) SetPosition(pos glm.Vec3) {
	c.position = pos
}

func (c *CameraObj) GetZoom() float32 {
	return c.zoom
}

func (c *CameraObj) GetFront() glm.Vec3 {
	return c.front
}

func (c *CameraObj) SetUp(up glm.Vec3) {
	c.up = up
}

func (c *CameraObj) SetYaw(yaw float32) {
	c.yaw = yaw
	c.updateVectors()
}

func (c *CameraObj) SetPitch(pitch float32) {
	c.pitch = pitch
	c.updateVectors()
}

func (c *CameraObj) SetSensitivity(s float32) {
	c.sensitivity = s
}

func (c *CameraObj) SetSpeed(s float32) {
	c.speed = s
}

func (c *CameraObj) SetZoom(z float32) {
	c.zoom = z
}

func (c *CameraObj) ProcessMovement(dt float32, dir int, speedUp float32) {
	velocity := c.speed * dt * speedUp
	switch dir {
	case DirectionForward:
		tmp := c.front.Mul(velocity)
		c.position = c.position.Add(&tmp)
	case DirectionBackwoard:
		tmp := c.front.Mul(velocity)
		c.position = c.position.Sub(&tmp)
	case DirectionLeft:
		tmp := c.front.Cross(&c.up)
		tmp.Normalize()
		tmp = tmp.Mul(velocity)
		c.position = c.position.Sub(&tmp)
	case DirectionRight:
		tmp := c.front.Cross(&c.up)
		tmp.Normalize()
		tmp = tmp.Mul(velocity)
		c.position = c.position.Add(&tmp)
	case DirectionUp:
		tmp := c.up.Mul(velocity)
		c.position = c.position.Add(&tmp)
	case DirectionDown:
		tmp := c.up.Mul(velocity)
		c.position = c.position.Sub(&tmp)
	}
}

func (c *CameraObj) ProcessMouseMove(xOffset, yOffset float32, pitch bool) {
	xOffset *= c.sensitivity
	yOffset *= c.sensitivity
	c.yaw += xOffset
	c.pitch += yOffset
	if pitch {
		if c.pitch > 89.9 {
			c.pitch = 89.9
		}
		if c.pitch < -89.9 {
			c.pitch = -89.9
		}
	}
	c.updateVectors()
}

func (c *CameraObj) ProcessScroll(yOffset float32) {
	c.zoom -= yOffset
	// TODO: limit max zoom & min zoom
}

func (c *CameraObj) GetID() uint32 {
	return c.id
}

func (c *CameraObj) GetName() string {
	return c.name
}

func (c *CameraObj) SetName(name string) {
	c.name = name
}

func (c *CameraObj) updateVectors() {
	sinYaw := math.Sin(float64(glm.DegToRad(c.yaw)))
	cosYaw := math.Cos(float64(glm.DegToRad(c.yaw)))
	sinPitch := math.Sin(float64(glm.DegToRad(c.pitch)))
	cosPitch := math.Sin(float64(glm.DegToRad(c.pitch)))

	c.front = glm.NormalizeVec3(glm.Vec3{
		float32(cosYaw * cosPitch),
		float32(sinPitch),
		float32(sinYaw * cosPitch),
	})
}

func NewCameraObj() (*CameraObj, error) {
	c := &CameraObj{}
	c.Init()
	return c, nil
}
