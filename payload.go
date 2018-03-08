package webwire

import (
	"bytes"
	"fmt"
	"unicode/utf16"
	"unicode/utf8"
)

// Payload represents an encoded message payload
type Payload struct {
	Encoding PayloadEncoding
	Data     []byte
}

// Utf8 returns a UTF8 representation of the payload data
func (pld *Payload) Utf8() (string, error) {
	if pld.Encoding == EncodingUtf16 {
		if len(pld.Data)%2 != 0 {
			return "", fmt.Errorf("Cannot convert invalid UTF16 payload data to UTF8")
		}
		u16str := make([]uint16, 1)
		// TODO: check whether the builtin append would be more efficient than bytes.Buffer
		utf8str := &bytes.Buffer{}
		utf8buf := make([]byte, 4)
		for i := 0; i < len(pld.Data); i += 2 {
			u16str[0] = uint16(pld.Data[i]) + (uint16(pld.Data[i+1]) << 8)
			rn := utf16.Decode(u16str)
			rnSize := utf8.EncodeRune(utf8buf, rn[0])
			utf8str.Write(utf8buf[:rnSize])
		}
		return utf8str.String(), nil
	}

	// Binary and UTF8 encoded payloads should pass through untouched
	return string(pld.Data), nil
}
