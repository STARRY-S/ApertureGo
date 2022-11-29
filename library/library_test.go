package library

import "testing"

func TestCgo(t *testing.T) {
	t.Logf("library supported error length: %v", GetErrorLength())

}

func TestLoadModel(t *testing.T) {
	LoadModel("test/cube.obj")
}

func TestCMemoryRelease(t *testing.T) {
	n := GetUnreleasedPointerNum()
	t.Logf("unreleased pointer num: %d", n)
	ReleaseAllMemory()
	n = GetUnreleasedPointerNum()
	if n != 0 {
		t.Error("ReleaseAllMemory failed")
	}
}
