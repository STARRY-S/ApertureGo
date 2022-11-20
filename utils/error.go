package utils

import "fmt"

var (
	ErrInvalidParameter = fmt.Errorf("invalid parameter")
	ErrInvalidPointer   = fmt.Errorf("invalid pointer")
	ErrInvalidFilePath  = fmt.Errorf("invalid file path")
	ErrReInitialize     = fmt.Errorf("resource already initialized")
	ErrEmptyFile        = fmt.Errorf("empty file")
)
