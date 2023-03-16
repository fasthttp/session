package session

import (
	"reflect"
	"testing"
)

func getSRC() Dict {
	src := newDictValue()

	src.KV["k1"] = "1"
	src.KV["k2"] = "2"

	return src
}

func getDST() Dict {
	return newDictValue()
}

func TestMSGPEncodeDecode(t *testing.T) {
	src := getSRC()
	dst := getDST()

	b1, err := MSGPEncode(src)
	if err != nil {
		t.Fatal(err)
	}

	if err := MSGPDecode(&dst, b1); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(src, dst) {
		t.Errorf("The results of 'src' and 'dst' must be equals, src = %v; dst = %v", src, dst)
	}
}

func TestBase64EncodeDecode(t *testing.T) {
	src := getSRC()
	dst := getDST()

	b1, err := Base64Encode(src)
	if err != nil {
		t.Fatal(err)
	}

	if err := Base64Decode(&dst, b1); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(src, dst) {
		t.Errorf("The results of 'src' and 'dst' must be equals, src = %v; dst = %v", src, dst)
	}
}

func BenchmarkMSGPEncode(b *testing.B) {
	src := getSRC()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := MSGPEncode(src); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMSGPDecode(b *testing.B) {
	src := getSRC()
	dst := getDST()

	srcBytes, _ := MSGPEncode(src)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := MSGPDecode(&dst, srcBytes); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBase64Encode(b *testing.B) {
	src := getSRC()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Base64Encode(src); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBase64Decode(b *testing.B) {
	src := getSRC()
	dst := getDST()

	srcBytes, _ := Base64Encode(src)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := Base64Decode(&dst, srcBytes); err != nil {
			b.Fatal(err)
		}
	}
}
