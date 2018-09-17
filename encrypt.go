package session

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"

	"github.com/savsgio/dictpool"
)

// session Encrypt tool
// - json
// - gob
// - base64

const (
	// BASE64TABLE base64 table
	BASE64TABLE = "1234567890poiuytreqwasdfghjklmnbvcxzQWERTYUIOPLKJHGFDSAZXCVBNM-_"
)

// Encrypt encrypt struct
type Encrypt struct{}

// NewEncrypt return new encrypt instance
func NewEncrypt() *Encrypt {
	return &Encrypt{}
}

// JSONEncode json encode
func (s *Encrypt) JSONEncode(data map[string]interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// JSONDecode json decode
func (s *Encrypt) JSONDecode(data []byte) (map[string]interface{}, error) {
	tempValue := make(map[string]interface{})
	err := json.Unmarshal(data, &tempValue)
	if err != nil {
		return tempValue, err
	}
	return tempValue, nil
}

// GOBEncode gob encode
func (s *Encrypt) GOBEncode(data *dictpool.Dict) ([]byte, error) {
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
func (s *Encrypt) GOBDecode(data []byte) (*dictpool.Dict, error) {
	d := dictpool.AcquireDict()

	if len(data) == 0 {
		return d, nil
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(d)
	if err != nil {
		return d, err
	}
	return d, nil
}

// Base64Encode base64 encode
func (s *Encrypt) Base64Encode(data *dictpool.Dict) ([]byte, error) {
	var coder = base64.NewEncoding(BASE64TABLE)
	b, err := s.GOBEncode(data)
	if err != nil {
		return []byte{}, err
	}
	return []byte(coder.EncodeToString(b)), nil
}

// Base64Decode base64 decode
func (s *Encrypt) Base64Decode(data []byte) (*dictpool.Dict, error) {
	var coder = base64.NewEncoding(BASE64TABLE)
	b, err := coder.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	return s.GOBDecode(b)
}
