package apkutils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
)

func ReadGzipHeader(buf []byte) bool {
	if buf[0] != GzipID1 || buf[1] != GzipID2 || buf[2] != GzipDeflate {
		return false
	}
	return true
}

func TarGzip(filename string, b []byte, writeEOFTar bool) (int, []byte, error) {
	nRead := len(b)
	var buf bytes.Buffer
	//	io.Copy(&buf, bytes.NewReader(b))
	gz := gzip.NewWriter(&buf)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	if writeEOFTar {
		defer tw.Close()
	}
	tw.WriteHeader(&tar.Header{
		Name: filename,
		Size: int64(nRead),
	})
	n, err := tw.Write(b)
	fmt.Printf("Wrote %d bytes\n", n)
	if err != nil {
		return 0, nil, err
	}
	tw.Close()
	gz.Close()
	ret := buf.Bytes()
	return n, ret, nil
}
