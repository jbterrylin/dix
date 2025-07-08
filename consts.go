package dix

import (
	"errors"
)

var ErrValueIsNil = errors.New("value is nil")
var ErrInvalidKey = errors.New("invalid key")
var ErrValueNotFound = errors.New("value not found")
var ErrInjectStructMustBePointerStruct = errors.New("inject struct must be pointer struct")
var ErrFieldCannotBeSet = errors.New("field cannot be set")
var ErrInjectFuncMustBeFunc = errors.New("inject func must be func")
var ErrInvalidVariable = errors.New("invalid variable")
var ErrTypeMismatch = errors.New("type mismatch")
var ErrRefCounterBelowZero = errors.New("ref counter below zero")

var DefaultValueKey ValueKey = ""
var DefaultProviderKey ProviderKey = ""
