package gov7

/*

// include v7.c because go doesn't compile .c file in a subdirectory

#include <stdlib.h>
#include "v7/v7.c"

#cgo CFLAGS: -DV7_NO_FS -DV7_LARGE_AST

*/
import "C"
import "unsafe"

import (
	"errors"
	"math"
)

type V7 C.struct_v7

type Val C.v7_val_t

func New() *V7 {
	return (*V7)(C.v7_create())
}

func (v7 *V7) Destroy() {
	C.v7_destroy((*C.struct_v7)(v7))
}

func (v7 *V7) ToJSON(v Val, size int) string {
	buflen := C.size_t(size)
	buf := (*C.char)(C.malloc(buflen))
	defer C.free(unsafe.Pointer(buf))

	C.v7_to_json((*C.struct_v7)(v7), C.v7_val_t(v), buf, buflen)
	return C.GoString(buf)
}

func (v7 *V7) Exec(code string) (Val, error) {
	var result C.v7_val_t

	cs := C.CString(code)
	defer C.free(unsafe.Pointer(cs))

	e := C.v7_exec((*C.struct_v7)(v7), &result, cs)

	v := Val(result)

	// DEBUG
	//fmt.Println(v7.ToJSON(v, 2048))

	switch e {
	case C.V7_OK:
		return v, nil
	case C.V7_SYNTAX_ERROR:
		// the document says the result is undefiend, but actually it's an error object
		return v, errors.New("parse error: " + v7.ToJSON(v, 2048))
	case C.V7_EXEC_EXCEPTION:
		return v, errors.New("runtime error: " + v7.ToJSON(v, 2048))
	case C.V7_STACK_OVERFLOW:
		return v, errors.New("stack overflow")
	}
	return v, errors.New("unknown error")
}

func (v7 *V7) IsUndefined(v Val) bool {
	return C.v7_is_undefined(C.v7_val_t(v)) != 0
}

func (v7 *V7) IsNumber(v Val) bool {
	return C.v7_is_number((C.v7_val_t)(v)) != 0
}

func (v7 *V7) IsString(v Val) bool {
	return C.v7_is_string((C.v7_val_t)(v)) != 0
}

func (v7 *V7) ToNumber(v Val) (float64, error) {
	if !v7.IsNumber(v) {
		return math.NaN(), errors.New("value is not a number")
	}
	return float64(C.v7_to_number(C.v7_val_t(v))), nil
}

func (v7 *V7) ToString(v Val) (string, error) {
	if !v7.IsString(v) {
		return "", errors.New("value is not a string")
	}
	var l *C.size_t = new(C.size_t)
	cv := C.v7_val_t(v)
	return C.GoString(C.v7_to_string((*C.struct_v7)(v7), &cv, l)), nil
}

func (v7 *V7) GetGlobalObject() Val {
	return Val(C.v7_get_global_object((*C.struct_v7)(v7)))
}

func (v7 *V7) Get(obj Val, name string) Val {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))

	return Val(C.v7_get((*C.struct_v7)(v7), C.v7_val_t(obj), cs, C.size_t(len(name))))
}

const (
	PROPERTY_READ_ONLY   uint = 1
	PROPERTY_DONT_ENUM   uint = 2
	PROPERTY_DONT_DELETE uint = 4
	PROPERTY_HIDDEN      uint = 8
	PROPERTY_GETTER      uint = 16
	PROPERTY_SETTER      uint = 32
)

func (v7 *V7) Set(obj Val, name string, attrs uint, v Val) error {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))

	if 0 != C.v7_set((*C.struct_v7)(v7), C.v7_val_t(obj), cs, C.size_t(len(name)), C.uint(attrs), C.v7_val_t(v)) {
		return errors.New("failed to set " + name)
	}
	return nil
}

func (v7 *V7) CreateUndefined() Val {
	return Val(C.v7_create_undefined())
}

/*
 * CreateString copies the Go string twice
 */
func (v7 *V7) CreateString(str string) Val {
	cs := C.CString(str)
	defer C.free(unsafe.Pointer(cs))

	return Val(C.v7_create_string((*C.struct_v7)(v7), cs, C.size_t(len(str)), 1 /*copy*/))
}

func (v7 *V7) CreateArray() Val {
	return Val(C.v7_create_array((*C.struct_v7)(v7)))
}

func (v7 *V7) ArrayPush(ary Val, v Val) {
	C.v7_array_push((*C.struct_v7)(v7), C.v7_val_t(ary), C.v7_val_t(v))
}

/*
 * Unlike Exec, there is no way to tell if Apply threw an error
 */
func (v7 *V7) Apply(f Val, thisObj Val, args Val) Val {
	return Val(C.v7_apply((*C.struct_v7)(v7), C.v7_val_t(f), C.v7_val_t(thisObj), C.v7_val_t(args)))
}
