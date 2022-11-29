package library

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -laperture -L${SRCDIR}/c_library/build
// #include "c_library/aperture.h"
import "C"
import "fmt"

func GetErrorLength() int {
	return C.AP_ERROR_LENGTH
}

func GetUnreleasedPointerNum() int {
	return int(C.ap_memory_unreleased_num())
}

func ReleaseAllMemory() {
	C.ap_memory_release()
}

func LoadModel(file string) error {
	cmodel := C.struct_AP_Model{}
	ret := C.ap_model_init_ptr(&cmodel, C.CString(file))
	C.AP_CHECK(ret)

	cmesh := cmodel.mesh
	fmt.Println("texture_length:", cmodel.texture_length)
	fmt.Println("mesh_length:", cmodel.mesh_length)
	fmt.Println("directory:", C.GoString(cmodel.directory))
	fmt.Println("cmesh.vertices_length:", cmesh.vertices_length)
	fmt.Println("cmesh.indices_length:", cmesh.indices_length)
	fmt.Println("cmesh.texture_length:", cmesh.texture_length)

	return nil
}
