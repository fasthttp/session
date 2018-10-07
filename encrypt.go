package session

import (
	"encoding/base64"
)

var b64Encoding = base64.NewEncoding(base64Table)

// NewEncrypt return new encrypt instance
func NewEncrypt() *Encrypt {
	return new(Encrypt)
}

// MSGPEncode MessagePack encode
func (e *Encrypt) MSGPEncode(src *Dict) ([]byte, error) {
	if len(src.D) == 0 {
		return nil, nil
	}

	dst, err := src.MarshalMsg(nil)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// MSGPDecode MessagePack decode
func (e *Encrypt) MSGPDecode(src []byte) (*Dict, error) {
	dst := new(Dict)

	if len(src) == 0 {
		return dst, nil
	}

	_, err := dst.UnmarshalMsg(src)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// Base64Encode base64 encode
func (e *Encrypt) Base64Encode(src *Dict) ([]byte, error) {
	srcBytes, err := e.MSGPEncode(src)
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

	return e.MSGPDecode(dst[:n])
}
