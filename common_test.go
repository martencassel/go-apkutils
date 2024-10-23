package apkutils

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"testing"

	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

func TestTarFile(t *testing.T) {
	t.Run("Test write a .tar.gz file", func(t *testing.T) {
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		w.WriteString("hello-world-1")
		_, _, err := TarGzip("hello1", buf.Bytes(), true)
		if err != nil {
			t.Fatal("Error writing tar file:", err)
		}
		buf = bytes.Buffer{}
		w = bufio.NewWriter(&buf)
		w.WriteString("hello-world-2")
		_, result, err := TarGzip("hello2", buf.Bytes(), false)
		if err != nil {
			t.Fatal("Error writing tar file:", err)
		}
		dgst := digest.FromBytes(result)
		ioutil.WriteFile("test.tar.gz", result, 0644)
		assert.True(t, dgst.String() == "sha256:54d818a850ebeec73a3b453812c02fc8a498a8bd0a096ed6822a87d88e462ee2")
	})

}

func TestReadGzipHeader(t *testing.T) {
	t.Run("Test ReadGzipHeader returns true when all Ids and deflate flags are present", func(t *testing.T) {
		buff := []byte{
			31,
			139,
			8,
		}
		want := true

		got := ReadGzipHeader(buff)
		if got != want {
			t.Errorf("Invalid response from ReadGzipHeader, got: %v, want: %v", got, want)
		}
	})
	t.Run("Test ReadGzipHeader with buffer size less than 3", func(t *testing.T) {
		// This test covers if there is a panic
		buff := []byte{
			7,
			31,
		}
		want := false

		got := ReadGzipHeader(buff)
		if got != want {
			t.Errorf("Invalid response from ReadGzipHeader, got: %v, want: %v", got, want)
		}
	})
}
