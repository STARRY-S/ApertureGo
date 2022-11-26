package utils_test

import (
	"errors"
	"testing"

	"github.com/STARRY-S/aperture/utils"
)

func TestError(t *testing.T) {
	errInvalidPointer := utils.ErrInvalidPointer
	if errInvalidPointer == utils.ErrInvalidPointer {
		t.Log("errInvalidPointer is equals to utils.ErrInvalidPointer")
	} else {
		t.Errorf("errInvalidPointer does not equals to utils.ErrInvalidPointer")
	}

	errInvalidPointer2 := errors.New(utils.ErrInvalidPointer.Error())
	if errInvalidPointer2 == utils.ErrInvalidPointer {
		t.Errorf("errInvalidPointer2 is equals to utils.ErrInvalidPointer")
	} else {
		t.Log("errInvalidPointer2 does not equals to utils.ErrInvalidPointer")
	}
}
