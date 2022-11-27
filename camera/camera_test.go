package camera_test

import (
	"testing"

	"github.com/STARRY-S/aperture"
	"github.com/STARRY-S/aperture/camera"
)

func TestInterface(t *testing.T) {
	c, err := camera.NewCameraObj()
	if err != nil {
		t.Error(err.Error())
	}
	var cam aperture.Camera
	cam = c
	if cam.GetID() != 0 {
		t.Error("camera id should be 0")
	}
}
