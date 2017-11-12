package main

import (
  "os"
  "io/ioutil"
  "fmt"
  "unicode/utf16"
  "unicode/utf8"
  "bytes"
)

func main() {
  file, _ := os.Open("/home/vvlad/Documents/EVE/logs/Chatlogs/Local_20171027_144535.txt")
  x, _ := ioutil.ReadAll(file)
  defer file.Close()
  s, _ := DecodeUTF16(x)
  fmt.Println(s)
}

func DecodeUTF16(b []byte) (string, error) {

	if len(b)%2 != 0 {
		return "", fmt.Errorf("Must have even length byte slice")
	}

	u16s := make([]uint16, 1)

	ret := &bytes.Buffer{}

	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String(), nil
}
