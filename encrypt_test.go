package session

import (
	"testing"
)

func BenchmarkGOBEncode(b *testing.B) {
	e := NewEncrypt()
	d := new(Dict)

	d.Set("k1", 1)
	d.Set("k2", 2)
	d.Set("k3", 3)
	d.Set("k4", 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.GOBEncode(d)
	}
}

func BenchmarkGOBDecode(b *testing.B) {
	e := NewEncrypt()
	d := new(Dict)

	d.Set("k1", 1)
	d.Set("k2", 2)
	d.Set("k3", 3)
	d.Set("k4", 4)

	data, _ := e.GOBEncode(d)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.GOBDecode(data)
	}
}

func BenchmarkBase64Encode(b *testing.B) {
	e := NewEncrypt()
	d := new(Dict)

	d.Set("k1", 1)
	d.Set("k2", 2)
	d.Set("k3", 3)
	d.Set("k4", 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Base64Encode(d)
	}
}

func BenchmarkBase64Decode(b *testing.B) {
	e := NewEncrypt()
	d := new(Dict)

	d.Set("k1", 1)
	d.Set("k2", 2)
	d.Set("k3", 3)
	d.Set("k4", 4)

	data, _ := e.Base64Encode(d)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Base64Decode(data)
	}
}
