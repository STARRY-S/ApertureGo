package utils

import "errors"

var (
	ErrInvalidDataType  = errors.New("invalid data type")
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrInvalidPointer   = errors.New("invalid pointer")
	ErrInvalidFilePath  = errors.New("invalid file path")
	ErrReInitialize     = errors.New("re-initialize the initialized resouce")
	ErrEmptyFile        = errors.New("file is empty")
	ErrPositionExceed   = errors.New("position exceeded of maximum value")
)
