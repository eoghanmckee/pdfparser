package pdf

import (
	"bytes"
	"strings"
)

var FontDefault *Font = &Font{map[int16]string{}, 1}

type Font struct {
	Cmap map[int16]string
	Width int
}

func NewFont(d Dictionary) *Font {
	cmap_string, _ := d.GetStream("ToUnicode")
	cmap := []byte(cmap_string)

	// create new font object
	font := &Font{map[int16]string{}, 1}

	// create parser for parsing cmap
	parser := NewParser(bytes.NewReader(cmap))

	for {
		// read next command
		command, operands, err := parser.ReadCommand()
		if err == ErrorRead {
			break
		}

		if command == KEYWORD_BEGIN_BF_RANGE {
			count, _ := operands.GetInt(len(operands) - 1)
			for i := 0; i < count; i++ {
				start_b, err := parser.ReadHexString(noDecryptor)
				if err != nil {
					break
				}
				font.Width = len([]byte(start_b))
				start := BytesToInt16([]byte(start_b))

				end_b, err := parser.ReadHexString(noDecryptor)
				if err != nil {
					break
				}
				end := BytesToInt16([]byte(end_b))

				value, err := parser.ReadHexString(noDecryptor)
				if err != nil {
					break
				}

				for i := start; i <= end; i++ {
					font.Cmap[i] = string(value)
				}
			}
		} else if command == KEYWORD_BEGIN_BF_CHAR {
			count, _ := operands.GetInt(len(operands) - 1)
			for i := 0; i < count; i++ {
				key_b, err := parser.ReadHexString(noDecryptor)
				if err != nil {
					break
				}
				font.Width = len([]byte(key_b))
				key := BytesToInt16([]byte(key_b))

				value, err := parser.ReadHexString(noDecryptor)
				if err != nil {
					break
				}

				font.Cmap[key] = string(value)
			}
		}
	}

	return font
}

func (font *Font) Decode(b []byte) string {
	var s strings.Builder
	for i := 0; i + font.Width <= len(b); i += font.Width {
		bs := b[i:i + font.Width]
		k := BytesToInt16(bs)
		if v, ok := font.Cmap[k]; ok {
			s.WriteString(v)
		} else {
			s.WriteString(string(bs))
		}
	}
	return s.String()
}
