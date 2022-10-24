package session

import (
	"encoding/base64"
)

var b64Encoding = base64.StdEncoding

// MSGPEncode MessagePack encode
func MSGPEncode(src Dict) ([]byte, error) {
	if len(src.KV) == 0 {
		return nil, nil
	}

	dst, err := src.MarshalMsg(nil)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// MSGPDecode MessagePack decode
func MSGPDecode(dst *Dict, src []byte) error {
	for k := range dst.KV {
		delete(dst.KV, k)
	}

	if len(src) == 0 {
		return nil
	}

	_, err := dst.UnmarshalMsg(src)

	return err
}

// Base64Encode base64 encode
func Base64Encode(src Dict) ([]byte, error) {
	srcBytes, err := MSGPEncode(src)
	if err != nil {
		return nil, err
	}

	dst := make([]byte, b64Encoding.EncodedLen(len(srcBytes)))
	b64Encoding.Encode(dst, srcBytes)

	return dst, nil
}

// Base64Decode base64 decode
func Base64Decode(dst *Dict, src []byte) error {
	tmp := make([]byte, b64Encoding.DecodedLen(len(src)))
	n, err := b64Encoding.Decode(tmp, src)
	if err != nil {
		return err
	}

	return MSGPDecode(dst, tmp[:n])
}
