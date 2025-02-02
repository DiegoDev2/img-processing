// +build !noasm !appengine
// Code generated by asm2asm, DO NOT EDIT.

package avx2

import (
	`github.com/bytedance/sonic/loader`
)

const (
    _entry__lspace = 32
)

const (
    _stack__lspace = 8
)

const (
    _size__lspace = 220
)

var (
    _pcsp__lspace = [][2]uint32{
        {1, 0},
        {184, 8},
        {188, 0},
        {204, 8},
        {208, 0},
        {215, 8},
        {219, 0},
    }
)

var _cfunc_lspace = []loader.CFunc{
    {"_lspace_entry", 0,  _entry__lspace, 0, nil},
    {"_lspace", _entry__lspace, _size__lspace, _stack__lspace, _pcsp__lspace},
}
