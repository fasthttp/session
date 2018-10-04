package session

import (
	"testing"
)

func BenchmarkGOBEncode(b *testing.B) {
	e := NewEncrypt()
	src := new(Dict)

	src.Set("k1", 1)
	src.Set("k2", 2)
	src.Set("k3", 3)
	src.Set("k4", 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.GOBEncode(src)
	}
}

func BenchmarkGOBDecode(b *testing.B) {
	e := NewEncrypt()
	src := new(Dict)

	src.Set("k1", 1)
	src.Set("k2", 2)
	src.Set("k3", 3)
	src.Set("k4", 4)

	srcBytes, _ := e.GOBEncode(src)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.GOBDecode(srcBytes)
	}
}

func BenchmarkBase64Encode(b *testing.B) {
	e := NewEncrypt()
	src := new(Dict)

	src.Set("k1", 1)
	src.Set("k2", 2)
	src.Set("k3", 3)
	src.Set("k4", 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Base64Encode(src)
	}
}

func BenchmarkBase64Decode(b *testing.B) {
	e := NewEncrypt()
	src := new(Dict)

	src.Set("k1", 1)
	src.Set("k2", 2)
	src.Set("k3", 3)
	src.Set("k4", 4)

	srcBytes, _ := e.Base64Encode(src)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Base64Decode(srcBytes)
	}
}
