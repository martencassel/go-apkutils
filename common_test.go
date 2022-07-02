package apkutils

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestTarFile(t *testing.T) {
	t.Run("Test TarFile 1", func(t *testing.T) {
		var buf bytes.Buffer
		tw := tar.NewWriter(&buf)
		hdr := &tar.Header{
			Name: "foobar.txt",
			Mode: 0600,
			Size: int64(len("Hello World")),
		}
		tw.WriteHeader(hdr)
		tw.Write([]byte("Hello World"))
		tw.Close()

		if int(hdr.Size) != len("Hello World") {
			t.Errorf("Expected %d bytes, got %d", hdr.Size, buf.Len())
		}
		fmt.Println(buf.Bytes())
		var buf2 bytes.Buffer
		io.Copy(&buf2, &buf)
		if buf2.Len() != 2048 {
			t.Errorf("Expected %d bytes, got %d", buf.Len(), buf2.Len())
		}
		fmt.Println(buf2.Bytes())
	})

}
