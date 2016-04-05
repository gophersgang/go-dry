package dry

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

type MyString struct {
	str string
}

func (t MyString) String() string {
	return t.str
}

type MyError struct {
	str string
}

func (t MyError) Error() string {
	return t.str
}

func Test_BytesReader(t *testing.T) {
	expected := []byte("hello")
	testBytesReaderFn := func(input interface{}) {
		result := make([]byte, 5)
		returnedIoReader := BytesReader(input)
		n, _ := returnedIoReader.Read(result)
		if n != 5 {
			t.FailNow()
		}
		for i, _ := range result {
			if result[i] != expected[i] {
				t.FailNow()
			}
		}
		n, err := returnedIoReader.Read(result)
		if n != 0 || err != io.EOF {
			t.FailNow()
		}
	}

	testBytesReaderFn(strings.NewReader("hello"))

	bytesInput := []byte("hello")
	testBytesReaderFn(bytesInput)

	testBytesReaderFn("hello")

	myStr := MyString{"hello"}
	testBytesReaderFn(myStr)

	myErr := MyError{"hello"}
	testBytesReaderFn(myErr)
}

func testCompressDecompress(t *testing.T,
	compressFunc func([]byte) []byte,
	decompressFunc func([]byte) []byte) {
	testFn := func(testData []byte) {
		compressedData := compressFunc(testData)
		uncompressedData := decompressFunc(compressedData)
		if !bytes.Equal(testData, uncompressedData) {
			t.FailNow()
		}
	}

	go testFn([]byte("hello123"))
	go testFn([]byte("gopher456"))
	go testFn([]byte("dry789"))
}

func Test_BytesDeflateInflate(t *testing.T) {
	testCompressDecompress(t, BytesDeflate, BytesInflate)
}

func Test_BytesGzipUnGzip(t *testing.T) {
	testCompressDecompress(t, BytesGzip, BytesUnGzip)
}

func Test_BytesMap(t *testing.T) {
	upper := func(b byte) byte {
		return b - ('a' - 'A')
	}
	result := BytesMap(upper, []byte("hello"))
	correct := []byte("HELLO")
	if len(result) != len(correct) {
		t.Fail()
	}
	for i, _ := range result {
		if result[i] != correct[i] {
			t.Fail()
		}
	}
}

func Test_BytesFilter(t *testing.T) {
	azFunc := func(b byte) bool {
		return b >= 'A' && b <= 'Z'
	}
	result := BytesFilter(azFunc, []byte{1, 2, 3, 'A', 'f', 'R', 123})
	correct := []byte{'A', 'R'}
	if len(result) != len(correct) {
		t.Fail()
	}
	for i, _ := range result {
		if result[i] != correct[i] {
			t.Fail()
		}
	}
}
