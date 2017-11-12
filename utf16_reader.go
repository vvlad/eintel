package eintel

import (
	"bytes"
	"fmt"
	"os"
	"unicode/utf16"
	"unicode/utf8"
)

type UTF16Reader struct {
	r *os.File
}

func NewUTF16Reader(file *os.File) UTF16Reader {
	return UTF16Reader{
		r: file,
	}
}

func (r UTF16Reader) ReadLine() (line string, err error) {
	b := make([]byte, 1)
	ret := bytes.Buffer{}

	for {
		_, err := r.Read(b)
		if err != nil {
			return "", err
		}

		if b[0] == '\n' {
			line := ret.String()
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			return line, nil
		}
		ret.Write(b)
	}
}

func (r UTF16Reader) Read(p []byte) (n int, err error) {
	b := make([]byte, len(p)*2)

	_, err = r.r.Read(b)
	if err != nil {
		return 0, err
	}

	rb, err := DecodeUTF16(b)
	if rb[0] == 239 || rb[0] == 191 || rb[0] == 189 {
		rb = rb[1:len(rb)]
	}
	copy(p, rb)
	return len(rb), err
}

func DecodeUTF16(b []byte) ([]byte, error) {

	ret := &bytes.Buffer{}
	if len(b)%2 != 0 {
		return ret.Bytes(), fmt.Errorf("Must have even length byte slice")
	}

	u16s := make([]uint16, 2)

	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.Bytes(), nil
}
