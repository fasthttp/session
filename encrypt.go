package session

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

// NewEncrypt return new encrypt instance
func NewEncrypt() *Encrypt {
	return new(Encrypt)
}

// GOBEncode gob encode
func (s *Encrypt) GOBEncode(data *Dict) ([]byte, error) {
	if len(data.D) == 0 {
		return []byte(""), nil
	}

	for _, kv := range data.D {
		gob.Register(kv)
	}
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(data)
	if err != nil {
		return []byte(""), err
	}

	return buf.Bytes(), nil
}

// GOBDecode gob decode data to map
func (s *Encrypt) GOBDecode(data []byte) (*Dict, error) {
	d := new(Dict)

	if len(data) == 0 {
		return d, nil
	}

	buf := bytes.NewBuffer(data)
	err := gob.NewDecoder(buf).Decode(d)
	if err != nil {
		return d, err
	}
	return d, nil
}

// Base64Encode base64 encode
func (s *Encrypt) Base64Encode(data *Dict) ([]byte, error) {
	var coder = base64.NewEncoding(base64Table)
	b, err := s.GOBEncode(data)
	if err != nil {
		return []byte{}, err
	}
	return []byte(coder.EncodeToString(b)), nil
}

// Base64Decode base64 decode
func (s *Encrypt) Base64Decode(data []byte) (*Dict, error) {
	var coder = base64.NewEncoding(base64Table)
	b, err := coder.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	return s.GOBDecode(b)
}
