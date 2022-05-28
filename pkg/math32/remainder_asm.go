//go:build !noasm && (amd64 || s390x || arm || ppc64le || 386 || wasm)
// +build !noasm
// +build amd64 s390x arm ppc64le 386 wasm

package math32

const haveArchRemainder = true

func archRemainder(x, y float32) float32
