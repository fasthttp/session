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
	e.gobEncodingPool = sync.Pool{
		New: func() interface{} {
			ge := new(gobEncoding)
			ge.buff = bytes.NewBuffer(nil)
			ge.encoder = gob.NewEncoder(ge.buff)
			ge.decoder = gob.NewDecoder(ge.buff)

			return ge
		},
	}

	return e
}

func (e *Encrypt) acquireGobEncoding() *gobEncoding {
	return e.gobEncodingPool.Get().(*gobEncoding)
}

func (e *Encrypt) releaseGobEncoding(ge *gobEncoding) {
	ge.buff.Reset()
	e.gobEncodingPool.Put(ge)
}

// GOBEncode gob encode
func (e *Encrypt) GOBEncode(src *Dict) ([]byte, error) {
	if len(src.D) == 0 {
		return nil, nil
	}

	ge := e.acquireGobEncoding()
	defer e.releaseGobEncoding(ge)

	err := ge.encoder.Encode(src)
	if err != nil {
		return nil, err
	}

	return ge.buff.Bytes(), nil
}

// GOBDecode gob decode data to Dict
func (e *Encrypt) GOBDecode(src []byte) (*Dict, error) {
	dst := new(Dict)

	if len(src) == 0 {
		return dst, nil
	}

	ge := e.acquireGobEncoding()
	defer e.releaseGobEncoding(ge)

	ge.buff.Write(src)

	err := ge.decoder.Decode(dst)
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
