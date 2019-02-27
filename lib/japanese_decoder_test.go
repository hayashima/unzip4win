package unzip4win

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestDecodeString_UTF8(t *testing.T) {
	content := loadTestFile("utf8.txt")
	actual, _ := decodeString(content)
	expected := "UTF-8の日本語です。"
	if actual != expected {
		t.Errorf("{%v} is expected. but actual is {%v}", expected, actual)
	}
}

func TestDecodeString_ShiftJIS(t *testing.T) {
	content := loadTestFile("sjis.txt")
	actual, _ := decodeString(content)
	expected := "SJISの日本語です。"
	if actual != expected {
		t.Errorf("{%v} is expected. but actual is {%v}", expected, actual)
	}
}

func loadTestFile(filename string) string {
	c, _ := ioutil.ReadFile(filepath.Join("..", "_tests", filename))
	return *(*string)(unsafe.Pointer(&c))
}
