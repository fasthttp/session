package session

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"sync"
)

var b64Encoding = base64.NewEncoding(base64Table)

// NewEncrypt return new encrypt instance
func NewEncrypt() *Encrypt {
	e := new(Encrypt)
	e.gobEncoderPool = sync.Pool{
		New: func() interface{} {
			ge := new(gobEncoder)

			ge.buff = bytes.NewBuffer(nil)
			ge.encoder = gob.NewEncoder(ge.buff)

			return ge
		},
	}
	e.gobDecoderPool = sync.Pool{
		New: func() interface{} {
			gd := new(gobDecoder)

			gd.reader = bytes.NewReader(nil)
			gd.decoder = gob.NewDecoder(gd.reader)

			return gd
		},
	}

	return e
}

func (e *Encrypt) acquireGobEncoder() *gobEncoder {
	return e.gobEncoderPool.Get().(*gobEncoder)
}

func (e *Encrypt) releaseGobEncoder(ge *gobEncoder) {
	ge.buff.Reset()
	e.gobEncoderPool.Put(ge)
}

func (e *Encrypt) acquireGobDecoder() *gobDecoder {
	return e.gobDecoderPool.Get().(*gobDecoder)
}

func (e *Encrypt) releaseGobDecoder(gd *gobDecoder) {
	e.gobDecoderPool.Put(gd)
}

// GOBEncode gob encode
func (e *Encrypt) GOBEncode(src *Dict) ([]byte, error) {
	if len(src.D) == 0 {
		return nil, nil
	}

	ge := e.acquireGobEncoder()
	// defer e.releaseGobEncoder(ge)

	err := ge.encoder.Encode(src)
	if err != nil {
		return nil, err
	}

	a := ge.buff.Bytes()
	b := make([]byte, len(a))
	copy(b, a)

	return b, nil
}

// GOBDecode gob decode data to Dict
func (e *Encrypt) GOBDecode(src []byte) (*Dict, error) {
	dst := new(Dict)

	if len(src) == 0 {
		return dst, nil
	}

	gd := e.acquireGobDecoder()
	// defer e.releaseGobDecoder(gd)

	gd.reader.Reset(src)

	err := gd.decoder.Decode(dst)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// Base64Encode base64 encode
func (e *Encrypt) Base64Encode(src *Dict) ([]byte, error) {
	srcBytes, err := e.GOBEncode(src)
	if err != nil {
		return nil, err
	}

	dst := make([]byte, b64Encoding.EncodedLen(len(srcBytes)))
	b64Encoding.Encode(dst, srcBytes)

	return dst, nil
}

// Base64Decode base64 decode
func (e *Encrypt) Base64Decode(src []byte) (*Dict, error) {
	dst := make([]byte, b64Encoding.DecodedLen(len(src)))
	n, err := b64Encoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}

	return e.GOBDecode(dst[:n])
}
