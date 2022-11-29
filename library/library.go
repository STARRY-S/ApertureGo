package library

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -laperture -L${SRCDIR}/c_library/build
// #include "c_library/aperture.h"
import "C"

func testPrintAllErrorMessage() {
	for i := 0; i < C.AP_ERROR_LENGTH; i++ {
		C.AP_CHECK(C.int(i))
	}
}

func C_LoadModel() {

}
